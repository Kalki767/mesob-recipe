package api

import (
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

func TestRoutes(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		wantStatus   int
		wantContains []string // each must appear in the body
		wantMissing  []string // none may appear in the body
		wantBody     string   // if set, the body must be exactly this
	}{
		// GET /recipes/{id}
		{name: "get existing recipe", method: "GET", path: "/recipes/1",
			wantStatus: http.StatusOK, wantContains: []string{`"id":1`, `"steps"`}},
		{name: "get missing recipe", method: "GET", path: "/recipes/99",
			wantStatus: http.StatusNotFound},
		{name: "get with non-number id", method: "GET", path: "/recipes/abc",
			wantStatus: http.StatusBadRequest},
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
			wantStatus: http.StatusBadRequest},
		{name: "title of only spaces", method: "GET", path: "/recipes?title=%20%20",
			wantStatus: http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := newTestHandler(t) // fresh store per case
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			body := rec.Body.String()
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body: %q)", rec.Code, tc.wantStatus, body)
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

func TestJSONResponsesSetContentType(t *testing.T) {
	handler := newTestHandler(t)

	for _, path := range []string{"/recipes", "/recipes/1", "/recipes?title=doro"} {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))

		if got := rec.Header().Get("Content-Type"); got != "application/json" {
			t.Errorf("%s: Content-Type = %q, want application/json", path, got)
		}
	}
}