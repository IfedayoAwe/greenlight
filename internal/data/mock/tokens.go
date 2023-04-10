package mock

import (
	"net/http"
	"time"

	"github.com/IfedayoAwe/greenlight/internal/data"
)

type MockTokenModel struct{}

func (m MockTokenModel) Insert(token *data.Token) error {
	return nil
}

func (m MockTokenModel) DeleteAllForUser(scope string, userID int64, userIP *string) error {
	return nil
}

func (m MockTokenModel) New(userID int64, ttl time.Duration, scope string, r *http.Request) (*data.Token, error) {
	token := data.Token{}
	return &token, nil
}
