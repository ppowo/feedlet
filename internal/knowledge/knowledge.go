package knowledge

import (
	"crypto/rand"
	"math/big"
)

// KnowledgeBit represents an interview question with its search term
type KnowledgeBit struct {
	Question   string // Interview-style question displayed to user
	SearchTerm string // Clean concept for DuckDuckGo search
}

// GetAllBits returns all knowledge bits from all languages combined
func GetAllBits() []KnowledgeBit {
	var all []KnowledgeBit
	all = append(all, getJavaBits()...)
	all = append(all, getJQueryBits()...)
	all = append(all, getJSPBits()...)
	return all
}

// GetRandomBit returns a single random knowledge bit
func GetRandomBit() KnowledgeBit {
	bits := GetAllBits()
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(bits))))
	if err != nil {
		return bits[0]
	}
	return bits[n.Int64()]
}
