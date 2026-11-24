package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kalki767/mesob-recipe/internal/recipe"
)

// newTestHandler builds a server over a fresh store with two known recipes:
//
//	1: Doro Wot  (Habesha food)
//	2: Shiro Wot (Habesha food)
//
// and returns its routes, so every request goes through the real mux.
func newTestHandler(t *testing.T) http.Handler {
	t.Helper()

	store := recipe.NewRecipeStore()
	for _, title := range []string{"Doro Wot", "Shiro Wot"} {
		r, err := recipe.NewRecipe(title, 100,
			[]recipe.Step{{Task: "Cook", Minutes: 10}},
			"Habesha food",
			[]recipe.Ingredient{{Name: "onion", Amount: 1}},
		)
		if err != nil {
			t.Fatalf("setup %q: %v", title, err)
		}
		store.Add(r)
	}
	return NewServer(store).Routes()
}

// send runs one request through the handler and returns what it wrote.
// An empty body means "no body".
func send(handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

// validTibs is a correct create-recipe body. The store already holds two
// recipes, so it gets ID 3.
const validTibs = `{"title":"Tibs","price_in_santim":30000,"category":"Habesha food",
	"steps":[{"task":"Fry the meat","minutes":15}],
	"ingredients":[{"name":"beef","amount":500,"unit":"g"},{"name":"salt","to_taste":true}]}`

func TestRoutes(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		body         string // request body; empty means none
		wantStatus   int
		wantLocation string   // if set, the Location header must be exactly this
		wantContains []string // each must appear in the response body
		wantMissing  []string // none may appear in the response body
		wantBody     string   // if set, the response body must be exactly this
	}{
		// GET /recipes/{id}
		{name: "get existing recipe", method: "GET", path: "/recipes/1",
			wantStatus: http.StatusOK, wantContains: []string{`"id":1`, `"steps"`}},
		{name: "get missing recipe", method: "GET", path: "/recipes/99",
			wantStatus: http.StatusNotFound, wantContains: []string{`"error"`}},
		{name: "get with non-number id", method: "GET", path: "/recipes/abc",
			wantStatus: http.StatusBadRequest, wantContains: []string{`"error"`}},
		{name: "wrong method on a recipe", method: "POST", path: "/recipes/1",
			wantStatus: http.StatusMethodNotAllowed},

		// GET /recipes
		{name: "list all", method: "GET", path: "/recipes",
			wantStatus: http.StatusOK, wantContains: []string{`"id":1`, `"id":2`}},
		{name: "list is summaries, not full recipes", method: "GET", path: "/recipes",
			wantStatus: http.StatusOK, wantMissing: []string{`"steps"`, `"ingredients"`}},
		{name: "search matches one", method: "GET", path: "/recipes?title=doro",
			wantStatus: http.StatusOK, wantContains: []string{`"id":1`}, wantMissing: []string{`"id":2`}},
		{name: "search matches several", method: "GET", path: "/recipes?title=WOT",
			wantStatus: http.StatusOK, wantContains: []string{`"id":1`, `"id":2`}},
		{name: "search with no match is [] not null", method: "GET", path: "/recipes?title=firfir",
			wantStatus: http.StatusOK, wantBody: "[]\n"},
		{name: "blank title", method: "GET", path: "/recipes?title=",
			wantStatus: http.StatusBadRequest, wantContains: []string{`"error"`}},
		{name: "title of only spaces", method: "GET", path: "/recipes?title=%20%20",
			wantStatus: http.StatusBadRequest, wantContains: []string{`"error"`}},

		// POST /recipes
		{name: "create valid recipe", method: "POST", path: "/recipes", body: validTibs,
			wantStatus: http.StatusCreated, wantLocation: "/recipes/3", wantBody: "{\"id\":3}\n"},
		{name: "create with broken JSON", method: "POST", path: "/recipes", body: `{"title": "Tibs",`,
			wantStatus: http.StatusBadRequest, wantContains: []string{`"error"`}},
		{name: "create with no body", method: "POST", path: "/recipes", body: "",
			wantStatus: http.StatusBadRequest, wantContains: []string{`"error"`}},
		{name: "create with price as text", method: "POST", path: "/recipes",
			body:       `{"title":"Tibs","price_in_santim":"abc","category":"x","steps":[{"task":"a","minutes":1}],"ingredients":[{"name":"b","amount":1}]}`,
			wantStatus: http.StatusBadRequest, wantContains: []string{`"error"`}},
		{name: "create with empty title", method: "POST", path: "/recipes",
			body:       `{"title":"","price_in_santim":100,"category":"x","steps":[{"task":"a","minutes":1}],"ingredients":[{"name":"b","amount":1}]}`,
			wantStatus: http.StatusBadRequest, wantContains: []string{`title cannot be empty`}},
		{name: "create with no ingredients", method: "POST", path: "/recipes",
			body:       `{"title":"Tibs","price_in_santim":100,"category":"x","steps":[{"task":"a","minutes":1}],"ingredients":[]}`,
			wantStatus: http.StatusBadRequest, wantContains: []string{`ingredients cannot be empty`}},
		{name: "create with ingredient both amount and to taste", method: "POST", path: "/recipes",
			body:       `{"title":"Tibs","price_in_santim":100,"category":"x","steps":[{"task":"a","minutes":1}],"ingredients":[{"name":"salt","amount":1,"to_taste":true}]}`,
			wantStatus: http.StatusBadRequest, wantContains: []string{`"error"`}},
		{name: "create with negative price", method: "POST", path: "/recipes",
			body:       `{"title":"Tibs","price_in_santim":-5,"category":"x","steps":[{"task":"a","minutes":1}],"ingredients":[{"name":"b","amount":1}]}`,
			wantStatus: http.StatusBadRequest, wantContains: []string{`"error"`}},
		{name: "wrong method on the list", method: "DELETE", path: "/recipes",
			wantStatus: http.StatusMethodNotAllowed},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := newTestHandler(t) // fresh store per case

			rec := send(handler, tc.method, tc.path, tc.body)

			body := rec.Body.String()
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body: %q)", rec.Code, tc.wantStatus, body)
			}
			if tc.wantLocation != "" {
				if got := rec.Header().Get("Location"); got != tc.wantLocation {
					t.Errorf("Location = %q, want %q", got, tc.wantLocation)
				}
			}
			if tc.wantBody != "" && body != tc.wantBody {
				t.Errorf("body = %q, want %q", body, tc.wantBody)
			}
			for _, s := range tc.wantContains {
				if !strings.Contains(body, s) {
					t.Errorf("body %q does not contain %s", body, s)
				}
			}
			for _, s := range tc.wantMissing {
				if strings.Contains(body, s) {
					t.Errorf("body %q should not contain %s", body, s)
				}
			}
		})
	}
}

// Every response the API itself writes (success or error) must be JSON.
// 405 is left out: the mux writes that one, not our code.
func TestResponsesAreJSON(t *testing.T) {
	tests := []struct {
		method, path, body string
	}{
		{"GET", "/recipes", ""},
		{"GET", "/recipes/1", ""},
		{"GET", "/recipes?title=doro", ""},
		{"POST", "/recipes", validTibs},
		{"GET", "/recipes/99", ""},        // 404
		{"GET", "/recipes/abc", ""},       // 400
		{"GET", "/recipes?title=", ""},    // 400
		{"POST", "/recipes", `{"title":`}, // 400
	}

	for _, tc := range tests {
		handler := newTestHandler(t)
		rec := send(handler, tc.method, tc.path, tc.body)

		if got := rec.Header().Get("Content-Type"); got != "application/json" {
			t.Errorf("%s %s: Content-Type = %q, want application/json", tc.method, tc.path, got)
		}
		if !json.Valid(rec.Body.Bytes()) {
			t.Errorf("%s %s: body is not valid JSON: %q", tc.method, tc.path, rec.Body.String())
		}
	}
}

// Error bodies must have the agreed shape: {"error": "<non-empty message>"}.
func TestErrorShape(t *testing.T) {
	handler := newTestHandler(t)
	rec := send(handler, "GET", "/recipes/99", "")

	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("error body is not a JSON object of strings: %v (body: %q)", err, rec.Body.String())
	}
	if got["error"] == "" {
		t.Errorf(`error body = %v, want a non-empty "error" field`, got)
	}
}

// A created recipe must really be in the store: fetch it with the URL from
// the Location header.
func TestCreatedRecipeCanBeFetched(t *testing.T) {
	handler := newTestHandler(t)

	created := send(handler, "POST", "/recipes", validTibs)
	if created.Code != http.StatusCreated {
		t.Fatalf("create: status = %d, want 201 (body: %q)", created.Code, created.Body.String())
	}
	location := created.Header().Get("Location")

	fetched := send(handler, "GET", location, "")
	if fetched.Code != http.StatusOK {
		t.Fatalf("GET %s: status = %d, want 200", location, fetched.Code)
	}
	for _, s := range []string{`"title":"Tibs"`, `"price_in_santim":30000`, `"total_minutes":15`, `"name":"beef"`} {
		if !strings.Contains(fetched.Body.String(), s) {
			t.Errorf("GET %s: body %q does not contain %s", location, fetched.Body.String(), s)
		}
	}
}