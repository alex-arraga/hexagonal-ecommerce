package security

import "golang.org/x/crypto/bcrypt"

type Hasher struct{}

// Hash hashes input password using bcrypt
func (h *Hasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// ComparePassword compares input password with hashed password
func (h *Hasher) ComparePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
