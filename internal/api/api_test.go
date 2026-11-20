package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kalki767/mesob-recipe/internal/recipe"
)


func TestGetRecipeFound(t *testing.T){
	req := httptest.NewRequest("GET", "/recipes/1", nil)
	rec := httptest.NewRecorder()

	Dororecipe, err := recipe.NewRecipe("  Doro Wot  ", 100,
		[]recipe.Step{{Task: "  Chop the onion  ", Minutes: 10}, {Task: "Cook", Minutes: 12}},
		"Habesha food",
		[]recipe.Ingredient{{Name: "  Onion  ", Amount: 3}},
	)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	store := recipe.NewRecipeStore()
	store.Add(Dororecipe)
	handler := NewServer(store).Routes()

	handler.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("wrong status code expected %d found %d", http.StatusOK, rec.Code)
	} else if !strings.Contains(rec.Body.String(), `"id":1`) {
		t.Errorf("expected id 1 but found %s", rec.Body.String())
	}
	
}

func TestRecipeNotFound(t *testing.T){
	req := httptest.NewRequest("GET", "/recipes/2", nil)
	rec := httptest.NewRecorder()

	Dororecipe, err := recipe.NewRecipe("  Doro Wot  ", 100,
		[]recipe.Step{{Task: "  Chop the onion  ", Minutes: 10}, {Task: "Cook", Minutes: 12}},
		"Habesha food",
		[]recipe.Ingredient{{Name: "  Onion  ", Amount: 3}},
	)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	store := recipe.NewRecipeStore()
	store.Add(Dororecipe)
	handler := NewServer(store).Routes()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound{
		t.Fatalf("wrong status code expected %d found %d", http.StatusNotFound, rec.Code)
	}
	
}

func TestInvalidId(t *testing.T){
	req := httptest.NewRequest("GET", "/recipes/abc", nil)
	rec := httptest.NewRecorder()

	Dororecipe, err := recipe.NewRecipe("  Doro Wot  ", 100,
		[]recipe.Step{{Task: "  Chop the onion  ", Minutes: 10}, {Task: "Cook", Minutes: 12}},
		"Habesha food",
		[]recipe.Ingredient{{Name: "  Onion  ", Amount: 3}},
	)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	store := recipe.NewRecipeStore()
	store.Add(Dororecipe)
	handler := NewServer(store).Routes()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("wrong status code expected %d found %d", http.StatusBadRequest, rec.Code)
	}
	
}