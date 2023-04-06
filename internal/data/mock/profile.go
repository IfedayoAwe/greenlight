package mock

import "github.com/IfedayoAwe/greenlight/internal/data"

type MockProfileModel struct{}

func (m MockProfileModel) Insert(profile *data.UserProfile) error {
	return nil
}

func (m MockProfileModel) Update(profile *data.UserProfile) error {
	return nil
}

func (m MockProfileModel) InsertProfilePic(userID int64) error {
	return nil
}

func (m MockProfileModel) Get(userID int64) (*data.UserProfile, error) {
	return nil, nil
}

func (m MockProfileModel) DeletOldPicture(imagePath string) error {
	return nil
}
