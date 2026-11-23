package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Kalki767/mesob-recipe/internal/api"
	"github.com/Kalki767/mesob-recipe/internal/recipe"
)

// ============================================================================
// Main Execution
// ============================================================================

func main() {
	recipeStore := recipe.NewRecipeStore()
	recipe1, err := recipe.NewRecipe(
		"Doro Wot",
		25005,
		[]recipe.Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 7},
		},
		"Habesha food",
		[]recipe.Ingredient{
			{Name: "  Onion  ", Amount: 3, Unit: ""},
		},
	)
	if err == nil {
		err = recipe1.AddIngredient(recipe.Ingredient{Name: "  oil  ", Amount: 1, Unit: "Cups"})
		if err != nil {
			fmt.Printf("couldn't add ingredient, found %s\n", err)
		}
		recipeStore.Add(recipe1)
	} else {
		fmt.Printf("couldn't create recipe: found %s\n", err)
	}

	recipe2, err := recipe.NewRecipe(
		"Shiro Wot",
		25005,
		[]recipe.Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 2},
		},
		"Habesha food",
		[]recipe.Ingredient{
			{Name: "onion", Amount: 2},
		},
	)
	if err == nil {
		recipeStore.Add(recipe2)
	} else {
		fmt.Printf("couldn't create recipe: found %s\n", err)
	}

	recipe3, err := recipe.NewRecipe(
		"Lasagna",
		25005,
		[]recipe.Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 4},
			{Task: "Mix it with egg", Minutes: 5},
		},
		"Foreign food",
		[]recipe.Ingredient{
			{Name: "onion", Amount: 2},
			{Name: "salt", ToTaste: true},
		},
	)
	if err == nil {
		recipeStore.Add(recipe3)
	} else {
		fmt.Printf("couldn't create recipe: found %s\n", err)
	}

	err = recipeStore.Delete(2)
	if err != nil {
		fmt.Printf("couldn't delete recipe: %s\n", err)
	}

	server := api.NewServer(recipeStore)

	err = http.ListenAndServe(":8081", server.Routes())

	if err != nil {
		log.Fatalf("failed to start server, found %s", err)
	}


}
