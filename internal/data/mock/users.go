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

type MockUserModel struct{}

func (m MockUserModel) Insert(user *data.User) error {
	return nil
}

func (m MockUserModel) GetByEmail(email string) (*data.User, error) {
	switch {
	case email == "olalekanawe99@gmail.com":
		return mockUser, nil
	default:
		return nil, data.ErrInvalidCredentials
	}
}

func (m MockUserModel) Update(user *data.User) error {
	return nil
}

func (m MockUserModel) GetForToken(tokenScope, tokenPlaintext string) (*data.User, error) {
	return mockUser, nil
}

func (m MockUserModel) ChangePassword(id int64, newPassword string) error {
	return nil
}

func (m MockUserModel) Delete(id int64) error {
	return nil
}
