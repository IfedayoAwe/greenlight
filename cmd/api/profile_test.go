package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestUserProfile(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Success", func(t *testing.T) {
		// Open the image file
		file, err := os.Open("../../images/Default.jpg")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		// Create a new multipart request
		buf := &bytes.Buffer{}
		writer := multipart.NewWriter(buf)

		// Add the image file to the request
		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			t.Fatal(err)
		}

		// Close the multipart request
		err = writer.Close()
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(http.MethodPut, ts.URL+"/v1/users/profile", buf)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI")
		req.Header.Add("Content-Type", writer.FormDataContentType())

		code, _, body := ts.do(t, req)
		if code != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, code)
		}

		if !bytes.Contains(body, []byte("profile")) {
			t.Errorf("want body to contain %q", []byte("profile"))
		}
	})
}

func TestGetUserProfile(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Success", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, ts.URL+"/v1/user/profile", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI")

		code, _, body := ts.do(t, req)
		if code != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, code)
		}

		if !bytes.Contains(body, []byte("olalekanawe99@gmail.com")) {
			t.Errorf("want body to contain %q", []byte("olalekanawe99@gmail.com"))
		}
	})
}
