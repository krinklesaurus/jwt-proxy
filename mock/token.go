package mock

import "github.com/krinklesaurus/jwt_proxy"

func NewTokenizer() app.Tokenizer {
	return &tokenizer{}
}

type tokenizer struct {
}

func (t *tokenizer) Serialize(claims map[string]interface{}) ([]byte, error) {
	return []byte("MOCK_JWT_TOKEN"), nil
}
