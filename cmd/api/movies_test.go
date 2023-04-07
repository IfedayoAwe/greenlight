package main

// func TestShowMovies(t *testing.T) {
// 	app := newTestApplication(t)
// 	ts := newTestServer(t, app.routes())
// 	defer ts.Close()

// 	t.Run("Unauthenticated", func(t *testing.T) {
// 		code, header, _ := ts.get(t, "/v1/movies/1")
// 		contentType := header.Get("Content-Type")
// 		if contentType != "application/json" {
// 			t.Errorf("want %q; got %q", "application/json", contentType)
// 		}

// 		if code != http.StatusUnauthorized {
// 			t.Errorf("want %d; got %d", http.StatusUnauthorized, code)
// 		}
// 	})

// 	t.Run("Authenticated", func(t *testing.T) {
// 		// Authenticate the user...
// 		code, header, _ := ts.get(t, "/v1/tokens/authentication")
// 		contentType := header.Get("Content-Type")
// 		if contentType != "application/json" {
// 			t.Errorf("want %q; got %q", "application/json", contentType)
// 		}

// 		if code != http.StatusCreated {
// 			t.Errorf("want %d; got %d", http.StatusUnauthorized, code)
// 		}

// 		form := url.Values{}
// 		form.Add("email", "olalekanawe99@gmail.com")
// 		form.Add("password", "")
// 		ts.postForm(t, "/user/login", form)
// 		// Then check that the authenticated user is shown the create snippet form.
// 		code, _, body := ts.get(t, "/snippet/create")
// 		if code != 200 {
// 			t.Errorf("want %d; got %d", 200, code)
// 		}

// 		formTag := "<form action='/snippet/create' method='POST'>"
// 		if !bytes.Contains(body, []byte(formTag)) {
// 			t.Errorf("want body %s to contain %q", body, formTag)
// 		}
// 	})

// }
