package mock

import "github.com/IfedayoAwe/greenlight/internal/data"

var mockPermissions1 = &data.Permissions{"movies:read", "movies:write"}
var mockPermissions2 = &data.Permissions{"movies:read"}

type MockPermissionModel struct{}

func (m MockPermissionModel) GetAllForUser(userID int64) (data.Permissions, error) {
	switch userID {
	case 1, 4:
		return *mockPermissions1, nil
	default:
		return *mockPermissions2, nil
	}

}

func (m MockPermissionModel) AddForUser(userID int64, codes ...string) error {
	return nil
}
