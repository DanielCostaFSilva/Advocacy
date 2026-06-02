package hasher

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{cost: bcrypt.DefaultCost}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type BcryptVerifier struct{}

func NewBcryptVerifier() *BcryptVerifier {
	return &BcryptVerifier{}
}

func (v *BcryptVerifier) Compare(plainPassword, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainPassword))
}

