package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Kalki767/mesob-recipe/internal/recipe"
)

type Server struct {
	store *recipe.RecipeStore
}

func NewServer(store *recipe.RecipeStore) *Server{
	return &Server{store: store}
}

func (s *Server) getRecipe(w http.ResponseWriter, r *http.Request){
	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil{
		http.Error(w, "invalid id format. Id should be a number", http.StatusBadRequest)
		return
	}

	view, err := s.store.Get(id)
	if err != nil {
		http.Error(w, "recipe doesn't exist", http.StatusNotFound)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(view)
	if err != nil {
		log.Printf("encode recipe %d: %v", id, err)
	}

}

func (s *Server) Routes() http.Handler{
	mux := http.NewServeMux()

	mux.HandleFunc("GET /recipes/{id}", s.getRecipe)

	return mux
}