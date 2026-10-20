package main

// func newTestHandler() http.Handler {
// 	svc := model.NewService(model.NewStore())
// 	log := slog.New(slog.NewTextHandler(io.Discard, nil))
// 	return NewServerHTTP(svc, log).routes()
// }

// func postRegister(h http.Handler, body string) httptest.ResponseRecorder {
// 	rec := httptest.NewRecorder()
// 	h.ServeHTTP(rec, httptest.NewRequest("POST", "/register", strings.NewReader(body)))
// 	return *rec
// }

// func TestRegister_Created(t *testing.T) {
// 	rec := postRegister(newTestHandler(), `{"email":"a@b.com","password":"supersecret"}`)

// 	if rec.Code != http.StatusCreated {
// 		t.Fatalf("got %d, want %d, body: %s", rec.Code, http.StatusCreated, rec.Body)
// 	}
// 	body := strings.ToLower(rec.Body.String())
// 	if strings.Contains(body, "supersecret") || strings.Contains(body, "hash") {
// 		t.Errorf("в ответе утечка пароля или хэша: %s", body)
// 	}
// }
