package core

// Tokenizer creates a byte array from an input map.
type Tokenizer interface {
	Serialize(claims map[string]interface{}) ([]byte, error)
}
