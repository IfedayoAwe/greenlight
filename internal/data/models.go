package data

import (
	"database/sql"
	"errors"
	"net/http"
	"time"
)

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrEditConflict       = errors.New("edit conflict")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Models struct {
	Movies interface {
		Insert(movie *Movie) error
		Get(id int64) (*Movie, error)
		Update(movie *Movie) error
		Delete(id int64) error
		GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error)
	}
	Tokens interface {
		Insert(token *Token) error
		DeleteAllForUser(scope string, userID int64, userIP *string) error
		New(userID int64, ttl time.Duration, scope string, r *http.Request) (*Token, error)
	}
	Users interface {
		Insert(user *User) error
		GetByEmail(email string) (*User, error)
		Update(user *User) error
		GetForToken(tokenScope, tokenPlaintext string) (*User, error)
		ChangePassword(id int64, newPassword string) error
		Delete(id int64) error
	}
	Permissions interface {
		GetAllForUser(userID int64) (Permissions, error)
		AddForUser(userID int64, codes ...string) error
	}
	UsersProfile interface {
		Insert(profile *UserProfile) error
		Update(profile *UserProfile) error
		InsertProfilePic(userID int64) error
		Get(userID int64) (*UserProfile, error)
		DeletOldPicture(imagePath string) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:       MovieModel{DB: db},
		Users:        UserModel{DB: db},
		Tokens:       TokenModel{DB: db},
		Permissions:  PermissionModel{DB: db},
		UsersProfile: ProfileModel{DB: db},
	}
}
