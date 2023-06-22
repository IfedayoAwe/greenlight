package main

import (
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/IfedayoAwe/greenlight/internal/data"
	"github.com/IfedayoAwe/greenlight/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func (app *application) userProfileHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" || !strings.HasPrefix(contentType, "multipart/form-data") {
		app.badRequestResponse(w, r, fmt.Errorf("Content-Type is not multipart/form-data"))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024) //2MB

	err := r.ParseMultipartForm(2 * 1024 * 1024) //2MB
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)

	v := validator.New()

	img, err := data.ValidateProfilePicture(v, fileHeader.Size, ext, file)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	fileName := fmt.Sprintf("%d%d%s", user.ID, time.Now().UnixNano(), ext)
	filePath := filepath.Join("images/profile", fileName)

	if app.config.profile.enabled {
		destFile, err := os.Create(filePath)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		defer destFile.Close()

		err = jpeg.Encode(destFile, img, &jpeg.Options{Quality: 90})
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	userProfile, err := app.models.UsersProfile.Get(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.UsersProfile.DeletOldPicture(userProfile.ImagePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	newFilePath := filepath.Join("profile", fileName)

	userProfile = &data.UserProfile{
		ImagePath: newFilePath,
		UserID:    user.ID,
	}

	err = app.models.UsersProfile.Update(userProfile)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": newFilePath}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	userProfile, err := app.models.UsersProfile.Get(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	userProfiledetails := struct {
		ID        int64
		CreatedAt time.Time
		Name      string
		Email     string
		ImagePath string
	}{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		Name:      user.Name,
		Email:     user.Email,
		ImagePath: userProfile.ImagePath,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": userProfiledetails}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	filePath := httprouter.ParamsFromContext(r.Context()).ByName("filepath")
	path := filepath.Join("./images/profile", filePath)
	http.ServeFile(w, r, filepath.FromSlash(filepath.Clean(path)))
}
