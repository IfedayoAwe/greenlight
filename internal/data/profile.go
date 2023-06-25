package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/nfnt/resize"

	"github.com/IfedayoAwe/greenlight/internal/validator"
)

type UserProfile struct {
	ImagePath string `json:"image_path"`
	UserID    int64  `json:"user_id"`
}

const maxFileSize = 2 * 1024 * 1024 // 2MB
const allowedExtensions = ".jpg,.jpeg,.png"

func validateAspectRatio(v *validator.Validator, img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	v.Check(width == height, "image", "aspect ratio is not 1:1")
	if img.Bounds().Dx() > 300 || img.Bounds().Dy() > 300 {
		img = resize.Resize(300, 300, img, resize.Lanczos3)
		return img
	}
	return img
}

func ValidateProfilePicture(v *validator.Validator, fileHeaderSize int64, ext string, file multipart.File) (image.Image, error) {
	fileExtensionErrorMessage := fmt.Sprintf("Invalid file type. Allowed file types are %s", allowedExtensions)
	fileSizeErrorMessage := fmt.Sprintf("File size exceeds the limit of %d bytes", maxFileSize)
	v.Check(fileHeaderSize <= maxFileSize, "file", fileSizeErrorMessage)

	var img image.Image
	var err error

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
		img = validateAspectRatio(v, img)

	case ".png":
		img, err = png.Decode(file)
		if err != nil {
			return nil, err
		}
		img = validateAspectRatio(v, img)

	default:
		v.AddError("file", fileExtensionErrorMessage)
	}

	return img, nil
}

func copyDefaultImage(destPath string) error {
	srcPath := "images/Default.jpg"

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcData, err := io.ReadAll(srcFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(destPath, srcData, 0644)
	if err != nil {
		return err
	}
	return nil
}

type ProfileModel struct {
	DB *sql.DB
}

var (
	ErrDuplicateProfile = errors.New("duplicate profile")
)

func (p ProfileModel) Insert(profile *UserProfile) error {
	query := `
	INSERT INTO users_profile (user_id, image_path)
	VALUES ($1, $2)`

	args := []interface{}{profile.UserID, profile.ImagePath}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_profile_user_id_key"`:
			return ErrDuplicateProfile
		default:
			return err
		}
	}
	return nil

}

func (p ProfileModel) InsertProfilePic(userID int64) error {
	fileName := fmt.Sprintf("%d%d%s", userID, time.Now().UnixNano(), ".jpg")
	filePath := filepath.Join("images/profile", fileName)
	newFilePath := filepath.Join("profile", fileName)

	err := copyDefaultImage(filePath)
	if err != nil {
		return err
	}

	userProfile := &UserProfile{
		ImagePath: newFilePath,
		UserID:    userID,
	}

	err = p.Insert(userProfile)

	return err

}

func (p ProfileModel) Update(profile *UserProfile) error {
	query := `
	UPDATE users_profile 
	SET image_path = $1
	WHERE user_id = $2`

	args := []interface{}{profile.ImagePath, profile.UserID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, args...)
	return err
}

func (p ProfileModel) Get(userID int64) (*UserProfile, error) {
	query := `
	SELECT image_path
	FROM users_profile
	WHERE user_id = $1`

	var userProfile UserProfile

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, userID).Scan(
		&userProfile.ImagePath,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &userProfile, nil
}

func (p ProfileModel) DeletOldPicture(imagePath string) error {
	oldPath := filepath.Join("images/", imagePath)
	err := os.Remove(oldPath)
	if err != nil {
		return err
	}
	return nil
}
