package mock

import "github.com/krinklesaurus/jwt_proxy"

func NewMockTokenizer() app.Tokenizer {
	return &MockTokenizer{}
}

type MockTokenizer struct {
}

func (t *MockTokenizer) Serialize(claims map[string]interface{}) ([]byte, error) {
	return []byte("MOCK_JWT_TOKEN"), nil
}
