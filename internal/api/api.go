package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Kalki767/mesob-recipe/internal/recipe"
)

type Server struct {
	store *recipe.RecipeStore
}

type Request struct {
	Title string `json:"title"`
	PriceInSantim int `json:"price_in_santim"`
	Steps         []recipe.Step `json:"steps"`
	Category      string `json:"category"`
	Ingredients   []recipe.Ingredient `json:"ingredients"`
}

func NewServer(store *recipe.RecipeStore) *Server{
	return &Server{store: store}
}

func (s *Server) getRecipe(w http.ResponseWriter, r *http.Request){
	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil{
		writeJSON(w, http.StatusBadRequest, map[string]string{"error":err.Error()})
		return
	}

	view, err := s.store.Get(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error":err.Error()})
		return

	}
	writeJSON(w, http.StatusOK, view)
}

func (s *Server) listRecipes (w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if !params.Has("title"){
		recipeSummary := s.store.List()
		writeJSON(w, http.StatusOK, recipeSummary)
		return
	} else {
		title := params.Get("title")
		recipeSummary, err := s.store.FindByTitle(title)
		if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error":err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, recipeSummary)
	}	
	
}

func (s *Server) postRecipe (w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error":err.Error()})
		return
	}

	newRecipe, err := recipe.NewRecipe(req.Title, req.PriceInSantim, req.Steps, req.Category, req.Ingredients)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error":err.Error()})
		return
	}

	id := s.store.Add(newRecipe)
	w.Header().Set("Location", fmt.Sprintf("/recipes/%d", id))
	data := map[string]int{"id" : id}
	writeJSON(w, http.StatusCreated, data)
	
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Printf("couldn't convert to json found %s", err)
	}
}
func (s *Server) Routes() http.Handler{
	mux := http.NewServeMux()

	mux.HandleFunc("GET /recipes/{id}", s.getRecipe)
	mux.HandleFunc("GET /recipes", s.listRecipes)
	mux.HandleFunc("POST /recipes", s.postRecipe)
	

	return mux
}