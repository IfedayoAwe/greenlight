package mock

import "github.com/IfedayoAwe/greenlight/internal/data"

func NewMockModels() data.Models {
	return data.Models{
		Movies:       &MockMovieModel{},
		Users:        &MockUserModel{},
		Tokens:       &MockTokenModel{},
		UsersProfile: &MockProfileModel{},
		Permissions:  &MockPermissionModel{},
	}
}
