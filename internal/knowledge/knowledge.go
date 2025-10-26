package knowledge

// KnowledgeBit represents a concept/tidbit to remember
type KnowledgeBit struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Code         string `json:"code,omitempty"`       // For Java tidbits (unchanged)
	HTMLCode     string `json:"htmlCode,omitempty"`   // HTML code examples for jQuery tidbits
	JQueryCode   string `json:"jqueryCode,omitempty"` // jQuery code for tidbits
	ModernCode   string `json:"modernCode,omitempty"` // Modern JS content
	Category     string `json:"category"`
}

// GetKnowledgeBits returns all embedded knowledge tidbits (Java)
func GetKnowledgeBits() []KnowledgeBit {
	return getJavaBits()
}

// GetJQueryBits returns jQuery knowledge tidbits
func GetJQueryBits() []KnowledgeBit {
	return getJQueryBits()
}
