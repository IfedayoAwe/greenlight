package mock

import (
	"time"

	"github.com/IfedayoAwe/greenlight/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func pass(password string) []byte {
	pass, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	return pass
}

var (
	MockUser = &data.User{
		ID:        1,
		Name:      "Olalekan Ifedayo Awe",
		Email:     "olalekanawe99@gmail.com",
		CreatedAt: time.Now(),
		Activated: true,
		Admin:     true,
		Version:   1,
		Password: data.Password{
			Hash: pass("1234567890"),
		},
	}
	MockUser2 = &data.User{
		ID:        2,
		Name:      "Ayo Awe",
		Email:     "ayo@gmail.com",
		CreatedAt: time.Now(),
		Activated: false,
		Admin:     false,
		Version:   1,
	}
	MockUser3 = &data.User{
		ID:        3,
		Name:      "Vicky Awe",
		Email:     "vicky@gmail.com",
		CreatedAt: time.Now(),
		Activated: true,
		Admin:     false,
		Version:   1,
	}
	MockUser4 = &data.User{
		ID:        4,
		Name:      "Mummy Awe",
		Email:     "mummy@gmail.com",
		CreatedAt: time.Now(),
		Activated: true,
		Admin:     false,
		Version:   1,
	}
)

var ()

type MockUserModel struct{}

func (m MockUserModel) Insert(user *data.User) error {
	switch user.Email {
	case "olalekanawe99@gmail.com", "ayo@gmail.com", "vicky@gmail.com", "mummy@gmail.com":
		return data.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m MockUserModel) GetByEmail(email string) (*data.User, error) {
	switch email {
	case "olalekanawe99@gmail.com":
		return MockUser, nil
	case "ayo@gmail.com":
		return MockUser2, nil
	case "vicky@gmail.com":
		return MockUser3, nil
	case "mummy@gmail.com":
		return MockUser4, nil
	default:
		return nil, data.ErrInvalidCredentials
	}
}

func (m MockUserModel) Update(user *data.User) error {
	return nil
}

func (m MockUserModel) GetForToken(tokenScope, tokenPlaintext string) (*data.User, error) {
	switch tokenPlaintext {
	case "HTE34GKUHNDUSJ3QRUT6IKWKRI":
		return MockUser, nil
	case "HTE34GKUHNDUSJ3QRUT6IKWKRJ":
		return MockUser2, nil
	case "HTE34GKUHNDUSJ3QRUT6IKWKRL":
		return MockUser3, nil
	case "HTE34GKUHNDUSJ3QRUT6IKWKRM":
		return MockUser4, nil
	default:
		return nil, data.ErrRecordNotFound
	}
}

func (m MockUserModel) ChangePassword(id int64, newPassword string) error {
	return nil
}

func (m MockUserModel) Delete(id int64) error {
	return nil
}
