package source

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	sucuriChallengePattern = regexp.MustCompile(`S='([^']+)'`)
	sucuriValueExprPattern = regexp.MustCompile(`(?s)\b([A-Za-z_$][A-Za-z0-9_$]*)\s*=(.*?);\s*document\.cookie\s*=`)
	sucuriCookiePattern    = regexp.MustCompile(`(?s)document\.cookie\s*=(.*?);\s*location\.reload\(\);`)
	fromCharCodePattern    = regexp.MustCompile(`^String\.fromCharCode\((\d+)\)$`)
)

func parseSucuriChallengeCookie(body []byte) (*http.Cookie, bool, error) {
	text := string(body)
	if !strings.Contains(text, "sucuri_cloudproxy_js") && !sucuriChallengePattern.MatchString(text) {
		return nil, false, nil
	}

	match := sucuriChallengePattern.FindStringSubmatch(text)
	if len(match) != 2 {
		return nil, true, fmt.Errorf("missing encoded challenge payload")
	}

	decoded, err := base64.StdEncoding.DecodeString(padBase64(match[1]))
	if err != nil {
		return nil, true, fmt.Errorf("decode challenge payload: %w", err)
	}

	script := string(decoded)
	valueMatch := sucuriValueExprPattern.FindStringSubmatch(script)
	if len(valueMatch) != 3 {
		return nil, true, fmt.Errorf("missing challenge value expression")
	}

	valueName := valueMatch[1]
	value, err := evalJSSimpleConcatExpr(valueMatch[2], nil)
	if err != nil {
		return nil, true, fmt.Errorf("evaluate challenge value: %w", err)
	}

	cookieMatch := sucuriCookiePattern.FindStringSubmatch(script)
	if len(cookieMatch) != 2 {
		return nil, true, fmt.Errorf("missing challenge cookie expression")
	}

	cookieLine, err := evalJSSimpleConcatExpr(cookieMatch[1], map[string]string{valueName: value})
	if err != nil {
		return nil, true, fmt.Errorf("evaluate challenge cookie: %w", err)
	}

	cookie, err := parseCookieAssignment(cookieLine)
	if err != nil {
		return nil, true, err
	}

	return cookie, true, nil
}

func padBase64(value string) string {
	if mod := len(value) % 4; mod != 0 {
		value += strings.Repeat("=", 4-mod)
	}
	return value
}

func parseCookieAssignment(raw string) (*http.Cookie, error) {
	parts := strings.Split(raw, ";")
	pair := strings.TrimSpace(parts[0])
	name, value, ok := strings.Cut(pair, "=")
	if !ok || strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("invalid challenge cookie assignment")
	}

	cookie := &http.Cookie{
		Name:  strings.TrimSpace(name),
		Value: strings.TrimSpace(value),
		Path:  "/",
	}

	for _, attr := range parts[1:] {
		attr = strings.TrimSpace(attr)
		lower := strings.ToLower(attr)
		switch {
		case strings.HasPrefix(lower, "path="):
			cookie.Path = strings.TrimSpace(attr[len("path="):])
		case strings.HasPrefix(lower, "max-age="):
			maxAge, err := strconv.Atoi(strings.TrimSpace(attr[len("max-age="):]))
			if err == nil {
				cookie.MaxAge = maxAge
			}
		}
	}

	return cookie, nil
}

func evalJSSimpleConcatExpr(expr string, vars map[string]string) (string, error) {
	parts, err := splitTopLevelConcat(expr)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	for _, part := range parts {
		value, err := evalJSSimpleTerm(part, vars)
		if err != nil {
			return "", err
		}
		b.WriteString(value)
	}

	return b.String(), nil
}

func splitTopLevelConcat(expr string) ([]string, error) {
	var parts []string
	start := 0
	depth := 0
	var quote byte
	escaped := false

	for i := 0; i < len(expr); i++ {
		ch := expr[i]
		if quote != 0 {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == quote {
				quote = 0
			}
			continue
		}

		switch ch {
		case '\'', '"':
			quote = ch
		case '(':
			depth++
		case ')':
			if depth == 0 {
				return nil, fmt.Errorf("unbalanced parentheses in expression")
			}
			depth--
		case '+':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(expr[start:i]))
				start = i + 1
			}
		}
	}

	if quote != 0 || depth != 0 {
		return nil, fmt.Errorf("unterminated expression")
	}

	parts = append(parts, strings.TrimSpace(expr[start:]))
	return parts, nil
}

func evalJSSimpleTerm(term string, vars map[string]string) (string, error) {
	term = strings.TrimSpace(term)
	if term == "" {
		return "", nil
	}

	if term[0] == '\'' || term[0] == '"' {
		value, err := unquoteJSString(term)
		if err != nil {
			return "", fmt.Errorf("invalid string literal %q: %w", term, err)
		}
		return value, nil
	}

	if match := fromCharCodePattern.FindStringSubmatch(term); len(match) == 2 {
		codePoint, err := strconv.Atoi(match[1])
		if err != nil {
			return "", fmt.Errorf("invalid char code %q: %w", term, err)
		}
		return string(rune(codePoint)), nil
	}

	if vars != nil {
		if value, ok := vars[term]; ok {
			return value, nil
		}
	}

	return "", fmt.Errorf("unsupported expression term %q", term)
}

func unquoteJSString(literal string) (string, error) {
	if len(literal) < 2 || literal[0] != literal[len(literal)-1] {
		return "", fmt.Errorf("unterminated string")
	}

	quote := literal[0]
	remaining := literal[1 : len(literal)-1]
	var b strings.Builder
	for len(remaining) > 0 {
		r, _, tail, err := strconv.UnquoteChar(remaining, quote)
		if err != nil {
			return "", err
		}
		b.WriteRune(r)
		remaining = tail
	}

	return b.String(), nil
}
