package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	type user struct {
		Name     string
		Email    string
		Password string
		Role     string
	}

	user1 := struct {
		Name     string
		Password string
	}{"Olalekan", "1234567890"}
	user2 := struct {
		Name  string
		Email string
	}{"Olalekan", "olalekan@gmail.com"}
	user3 := struct {
		Email    string
		Password string
	}{"olalekan@gmail.com", "1234567890"}
	user4 := user{"Olalekan", "olalekan99@gmail.com", "1234567890", "viewer"}
	user5 := user{"Olalekan Awe", "olalekanawe99@gmail.com", "1234567890", "viewer"}
	user6 := user{"", "", "", ""}
	user7 := user{"Olalekan Awe", "ola.com", "1234567890", "viewer"}
	user8 := user{"Olalekan Awe", "olalekan99@gmail.com", "123", "viewer"}
	user9 := struct{ Foo string }{"1234567890"}
	user10 := user{"Olalekan", "olalekan99@gmail.com", "1234567890", "contributor"}

	tests := []struct {
		name     string
		user     interface{}
		wantCode int
		wantBody []byte
	}{
		{"Email", user1, http.StatusUnprocessableEntity, []byte("\"email\": \"must be provided\"")},
		{"Password", user2, http.StatusUnprocessableEntity, []byte("\"password\": \"must be provided\"")},
		{"Name", user3, http.StatusUnprocessableEntity, []byte("\"name\": \"must be provided\"")},
		{"NewUserViewer", user4, http.StatusAccepted, []byte("olalekan99@gmail.com")},
		{"DuplicateUser", user5, http.StatusUnprocessableEntity, []byte("a user with this email address already exists")},
		{"EmptyParameters", user6, http.StatusUnprocessableEntity, []byte("must be provided")},
		{"invalidEmail", user7, http.StatusUnprocessableEntity, []byte("must be a valid email address")},
		{"invalidPassword", user8, http.StatusUnprocessableEntity, []byte("\"password\": \"must be at least 10 bytes long\"")},
		{"UnknownKey", user9, http.StatusBadRequest, []byte("body contains unknown key")},
		{"NewUserContributor", user10, http.StatusAccepted, []byte("olalekan99@gmail.com")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			payload, err := json.Marshal(tt.user)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/users", bytes.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Add("Content-Type", "application/json")

			code, header, body := ts.do(t, req)
			if contentType := header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("want %q; got %q", "application/json", contentType)
			}

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}

		})
	}

}

func TestActivateUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	testToken := struct{ Token string }{"HTE34GKUHNDUSJ3QRUT6IKWKRJ"}
	testToken2 := struct{ Token string }{"HTE34GKUHNDUSJ3QRUT6IKWKRG"}
	testToken3 := struct{ Foo string }{"HTE34GKUHNDUSJ3QRUT6IKWKRJ"}
	testToken4 := struct{ Token string }{"$#@USJ3QRUT6IKWKRG"}
	testToken5 := struct{ Token string }{""}

	tests := []struct {
		name     string
		token    interface{}
		wantCode int
		wantBody []byte
	}{
		{"ValidToken", testToken, http.StatusOK, []byte("ayo@gmail.com")},
		{"InvalidToken", testToken2, http.StatusUnprocessableEntity, []byte("invalid or expired activation token")},
		{"BadRequest", testToken3, http.StatusBadRequest, []byte("body contains unknown key")},
		{"ShortToken", testToken4, http.StatusUnprocessableEntity, []byte("must be 26 bytes long")},
		{"EmptyToken", testToken5, http.StatusUnprocessableEntity, []byte("must be provided")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			payload, err := json.Marshal(tt.token)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPut, ts.URL+"/v1/users/activated", bytes.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Add("Content-Type", "application/json")

			code, header, body := ts.do(t, req)
			if contentType := header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("want %q; got %q", "application/json", contentType)
			}

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}

		})
	}

}

func TestChangePassword(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name            string
		token           string
		wantCode        int
		wantBody        []byte
		currentpassword string
		newpassword     string
		confirmpassword string
	}{
		{"ValidPassword", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusOK, []byte("token"), "1234567890", "pa5555word", "pa5555word"},
		{"InvalidPassword", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusUnauthorized, []byte("invalid password"), "12345678", "pa5555word", "pa5555word"},
		{"ShortPassword", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusUnprocessableEntity, []byte("must be at least 10 bytes long"), "1234567890", "paword", "paword"},
		{"NoCurrentPassword", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusUnprocessableEntity, []byte("Current Password field cannot be empty"), "", "pa5555word", "pa5555word"},
		{"NoNewPassword", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusUnprocessableEntity, []byte("must be provided"), "1234567890", "", ""},
		{"PasswordOverflow", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusUnprocessableEntity, []byte("must not be more than 72 bytes long"), "1234567890", "pa5555worddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd", "pa5555worddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"},
		{"PasswordMismatch", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusUnprocessableEntity, []byte("New Password and password confirmation do not match"), "1234567890", "pa5555word", "paword"},
		{"NoConfirmPassword", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusUnprocessableEntity, []byte("Confirm password field cannot be empty"), "1234567890", "pa5555word", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			password := struct {
				CurrentPassword string
				Password        string
				ConfirmPassword string
			}{
				CurrentPassword: tt.currentpassword,
				Password:        tt.newpassword,
				ConfirmPassword: tt.confirmpassword,
			}

			payload, err := json.Marshal(password)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPut, ts.URL+"/v1/users/change-password", bytes.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", tt.token)
			req.Header.Add("Content-Type", "application/json")

			code, header, body := ts.do(t, req)
			if contentType := header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("want %q; got %q", "application/json", contentType)
			}

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}

		})
	}
}

func TestResetUserPassword(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		token    string
		wantCode int
		wantBody []byte
		password string
	}{
		{"Sucessful", "HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusOK, []byte("your password was successfully reset"), "pa5555word"},
		{"EmptyToken", "", http.StatusUnprocessableEntity, []byte("must be provided"), "pa5555word"},
		{"BadToken", "NDUSJ3QRUT6IKWK", http.StatusUnprocessableEntity, []byte("must be 26 bytes long"), "pa5555word"},
		{"InvalidToken", "HTE34GKUHNDUSJ3QRUT6IKWKRX", http.StatusUnprocessableEntity, []byte("invalid or expired password reset token"), "pa5555word"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			password := struct {
				Password string
				Token    string
			}{
				Password: tt.password,
				Token:    tt.token,
			}

			payload, err := json.Marshal(password)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPut, ts.URL+"/v1/users/password", bytes.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Add("Content-Type", "application/json")

			code, header, body := ts.do(t, req)
			if contentType := header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("want %q; got %q", "application/json", contentType)
			}

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}

		})
	}

}

func TestUpdateUserDetails(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	user1 := struct{ Name string }{"Jerry"}
	user2 := struct{ Email string }{"jerry@gmail.com"}
	user3 := struct{ Foo string }{"jerry@gmail.com"}
	user4 := struct{ Email string }{"foo@gmail.com"}
	user5 := struct{ Name string }{""}

	tests := []struct {
		name     string
		token    string
		wantCode int
		wantBody []byte
		user     interface{}
	}{
		{"Name", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusOK, []byte("Jerry"), user1},
		{"Email", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusOK, []byte("jerry@gmail.com"), user2},
		{"Foo", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusBadRequest, []byte("body contains unknown key"), user3},
		{"DuplicateEmail", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusUnprocessableEntity, []byte("a user with this email address already exists"), user4},
		{"EmptyName", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusUnprocessableEntity, []byte("must be provided"), user5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			payload, err := json.Marshal(tt.user)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPatch, ts.URL+"/v1/users/update-details", bytes.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", tt.token)
			req.Header.Add("Content-Type", "application/json")

			code, header, body := ts.do(t, req)
			if contentType := header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("want %q; got %q", "application/json", contentType)
			}

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}

		})
	}
}

func TestUserLogout(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		token    string
		wantCode int
		wantBody []byte
	}{
		{"Sucessful", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusOK, []byte("user sucessfully logged out")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodDelete, ts.URL+"/v1/users/logout", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", tt.token)
			req.Header.Add("Content-Type", "application/json")

			code, header, body := ts.do(t, req)
			if contentType := header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("want %q; got %q", "application/json", contentType)
			}

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}

		})
	}

}

func TestDeleteUserAccount(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		token    string
		wantCode int
		wantBody []byte
	}{
		{"Sucessful", "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", http.StatusOK, []byte("user account successfully deleted")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodDelete, ts.URL+"/v1/users/delete", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", tt.token)
			req.Header.Add("Content-Type", "application/json")

			code, header, body := ts.do(t, req)
			if contentType := header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("want %q; got %q", "application/json", contentType)
			}

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}

		})
	}

}
