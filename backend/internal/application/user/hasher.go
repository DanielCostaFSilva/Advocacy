package user

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type MockPasswordHasher struct {
	HashFunc func(password string) (string, error)
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	if m.HashFunc != nil {
		return m.HashFunc(password)
	}
	return "", nil
}
