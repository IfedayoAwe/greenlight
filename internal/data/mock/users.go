package mock

import (
	"time"

	"github.com/IfedayoAwe/greenlight/internal/data"
)

var mockUser = &data.User{
	ID:        1,
	Name:      "Olalekan Ifedayo Awe",
	Email:     "olalekanawe99@gmail.com",
	CreatedAt: time.Now(),
	Activated: true,
	Admin:     true,
	Version:   1,
}

var mockUser2 = &data.User{
	ID:        2,
	Name:      "Ayo Awe",
	Email:     "ayo@gmail.com",
	CreatedAt: time.Now(),
	Activated: false,
	Admin:     false,
	Version:   1,
}

type MockUserModel struct{}

func (m MockUserModel) Insert(user *data.User) error {
	return nil
}

func (m MockUserModel) GetByEmail(email string) (*data.User, error) {
	switch {
	case email == "olalekanawe99@gmail.com":
		return mockUser, nil
	case email == "ayo@gmail.com":
		return mockUser2, nil
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
		return mockUser, nil
	case "HTE34GKUHNDUSJ3QRUT6IKWKRJ":
		return mockUser2, nil
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
