package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

func TestShowMovies(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
		token    string
	}{
		{"Unauthenticated", "/v1/movies/1", http.StatusUnauthorized, nil, ""},
		{"Authenticated", "/v1/movies/1", http.StatusOK, []byte("Test Movie"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"NotFoundFoo", "/v1/movies/foo", http.StatusNotFound, nil, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"NotFound", "/v1/movies/2", http.StatusNotFound, nil, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"UnActivated", "/v1/movies/1", http.StatusForbidden, nil, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRJ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodGet, ts.URL+tt.urlPath, nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", tt.token)

			code, header, body := ts.do(t, req)
			contentType := header.Get("Content-Type")
			if contentType != "application/json" {
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

func TestCreateMovie(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		title    string
		year     int
		runtime  string
		genres   []string
		wantCode int
		wantBody []byte
		token    string
	}{
		{"Authenticated", "Mountain", 2003, "200 mins", []string{"Comedy", "Romance"}, http.StatusCreated, []byte("Mountain"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"UnAuthenticated", "Road", 2000, "300 mins", []string{"Drama", "Comedy"}, http.StatusUnauthorized, []byte("you must be authenticated to access this resource"), ""},
		{"UnActivated", "Mountain", 2004, "21 mins", []string{"Action", "Horror"}, http.StatusForbidden, []byte("your user account must be activated to access this resource"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRJ"},
		{"NotPermitted", "Top", 1991, "500 mins", []string{"Drama", "Romance"}, http.StatusForbidden, []byte("your user account is not permitted to access this resource"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRL"},
		{"NoTitle", "", 2003, "200 mins", []string{"Comedy", "Romance"}, http.StatusUnprocessableEntity, []byte("\"title\": \"must be provided\""), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"NoYear", "Mountain", 0, "200 mins", []string{"Comedy", "Romance"}, http.StatusUnprocessableEntity, []byte("\"year\": \"must be provided\""), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"LongTitle", "qwertyuioplkjhgfdsazxcvbnmklpoiuytrewqasdfghjklmnbvc", 2005, "200 mins", []string{"Comedy", "Romance"}, http.StatusUnprocessableEntity, []byte("\"title\": \"must not be more than 50 bytes long\""), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"EarlyYear", "Mountain", 1882, "200 mins", []string{"Comedy", "Romance"}, http.StatusUnprocessableEntity, []byte("\"year\": \"must be greater than 1888\""), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"FutureYear", "Mountain", 3000, "200 mins", []string{"Comedy", "Romance"}, http.StatusUnprocessableEntity, []byte("\"year\": \"must not be in the future\""), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"NoRuntime", "Mountain", 2000, "0 mins", []string{"Comedy", "Romance"}, http.StatusUnprocessableEntity, []byte("\"runtime\": \"must be provided\""), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"NegativeRuntime", "Mountain", 2000, "-1 mins", []string{"Comedy", "Romance"}, http.StatusUnprocessableEntity, []byte("\"runtime\": \"must be a positive integer\""), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"NoGenres", "Mountain", 2003, "200 mins", nil, http.StatusUnprocessableEntity, []byte("\"genres\": \"must be provided\""), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"EmptyGenres", "Mountain", 2003, "200 mins", []string{}, http.StatusUnprocessableEntity, []byte("\"genres\": \"must contain at least 1 genre\""), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"PlentyGenres", "Mountain", 2003, "200 mins", []string{"Horror", "Comedy", "Romance", "Action", "Drama", "SCI-FI"}, http.StatusUnprocessableEntity, []byte("\"genres\": \"must not contain more than 5 genres\""), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"DuplicateGenres", "Mountain", 2003, "200 mins", []string{"Horror", "Horror"}, http.StatusUnprocessableEntity, []byte("\"genres\": \"must not contain duplicate values\""), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"EmptyGenreString", "Mountain", 2003, "200 mins", []string{""}, http.StatusUnprocessableEntity, []byte("\"genres\": \"field must not be empty\""), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			movie := struct {
				Title   string
				Year    int
				Runtime string
				Genres  []string
			}{
				Title:   tt.title,
				Year:    tt.year,
				Runtime: tt.runtime,
				Genres:  tt.genres,
			}

			payload, err := json.Marshal(movie)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/movies", bytes.NewReader(payload))
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

func TestUpdateMovie(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	movie1 := struct{ Title string }{"Movies Test 1"}
	movie2 := struct{ Year int }{2005}
	movie3 := struct{ Runtime string }{"230 mins"}
	movie4 := struct{ Genres []string }{[]string{"Fantasy"}}
	movie5 := struct{ Year int }{20}

	tests := []struct {
		name     string
		movie    interface{}
		wantCode int
		wantBody []byte
		token    string
		urlPath  string
	}{
		{"Title", movie1, http.StatusOK, []byte("Movies Test 1"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies/1"},
		{"Year", movie2, http.StatusOK, []byte(strconv.Itoa(2005)), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies/1"},
		{"Runtime", movie3, http.StatusOK, []byte("230 mins"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies/1"},
		{"Genres", movie4, http.StatusOK, []byte("Fantasy"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies/1"},
		{"FailedValidation", movie5, http.StatusUnprocessableEntity, []byte("must be greater than 1888"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies/1"},
		{"NotExist", movie4, http.StatusNotFound, []byte("the requested resource could not be found"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies/5"},
		{"NotExistFoo", movie4, http.StatusNotFound, []byte("the requested resource could not be found"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies/foo"},
		{"NotPermitted", movie1, http.StatusForbidden, []byte("your user account is not permitted to access this resource"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRM", "/v1/movies/1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			payload, err := json.Marshal(tt.movie)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPatch, ts.URL+tt.urlPath, bytes.NewReader(payload))
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

func TestDeleteMovie(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		wantCode int
		wantBody []byte
		token    string
		urlPath  string
	}{
		{"Authenticated", http.StatusOK, []byte("movie successfully deleted"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies/1"},
		{"NotFound", http.StatusNotFound, []byte("the requested resource could not be found"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies/2"},
		{"NotFoundFoo", http.StatusNotFound, []byte("the requested resource could not be found"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies/foo"},
		{"NotPermitted", http.StatusForbidden, []byte("your user account is not permitted to access this resource"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRM", "/v1/movies/1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodDelete, ts.URL+tt.urlPath, nil)
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

func TestListMovie(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		wantCode int
		token    string
		urlPath  string
	}{
		{"AllMovies", http.StatusOK, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies"},
		{"FilterMovies", http.StatusOK, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies?title=testmovie&genres=comedy,action&page=1&page_size=5&sort=-year"},
		{"SortYear", http.StatusOK, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies?sort=year"},
		{"SortIDDesc", http.StatusOK, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies?sort=-id"},
		{"SortID", http.StatusOK, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies?sort=id"},
		{"SortTitle", http.StatusOK, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies?sort=title"},
		{"SortTitleDesc", http.StatusOK, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies?sort=-title"},
		{"SortRuntime", http.StatusOK, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies?sort=runtime"},
		{"SortRuntimeDesc", http.StatusOK, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies?sort=-runtime"},
		{"FailedValidation", http.StatusUnprocessableEntity, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI", "/v1/movies?sort=foo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, ts.URL+tt.urlPath, nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", tt.token)

			code, header, _ := ts.do(t, req)
			if contentType := header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("want %q; got %q", "application/json", contentType)
			}

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

		})
	}

}
