package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Email string `json:"email,omitempty"`
	Pwd   string `json:"password,omitempty"`
}

func (s *User) HashPassword() (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.Pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	s.Pwd = string(hashedPassword)
	return s, nil
}
