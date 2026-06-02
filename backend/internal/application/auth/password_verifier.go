package auth

type MockPasswordVerifier struct {
	CompareFunc func(plainPassword, hash string) error
}

func (m *MockPasswordVerifier) Compare(plainPassword, hash string) error {
	if m.CompareFunc != nil {
		return m.CompareFunc(plainPassword, hash)
	}
	return nil
}
