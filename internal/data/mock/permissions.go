package mock

import "github.com/IfedayoAwe/greenlight/internal/data"

type MockPermissionModel struct{}

func (m MockPermissionModel) GetAllForUser(userID int64) (data.Permissions, error) {
	return nil, nil
}

func (m MockPermissionModel) AddForUser(userID int64, codes ...string) error {
	return nil
}
