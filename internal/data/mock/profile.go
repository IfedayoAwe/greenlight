package mock

import "github.com/IfedayoAwe/greenlight/internal/data"

var mockUserProfile = data.UserProfile{
	UserID:    1,
	ImagePath: "/mock/image.png",
}

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
	switch userID {
	case 1:
		return &mockUserProfile, nil
	default:
		return nil, data.ErrRecordNotFound
	}

}

func (m MockProfileModel) DeletOldPicture(imagePath string) error {
	return nil
}
