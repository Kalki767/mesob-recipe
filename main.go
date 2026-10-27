package main

import (
	"fmt"

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
		err = recipe2.AddIngredient(recipe.Ingredient{Name: "Onion", Amount: 0, Unit: ""})
		if err != nil {
			fmt.Printf("couldn't add ingredient, found %s\n", err)
		}
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
	recipe4, err := recipe.NewRecipe(
		"Cake",
		25005,
		[]recipe.Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 9},
		},
		"Foreign food",
		[]recipe.Ingredient{},
	)
	if err == nil {
		recipeStore.Add(recipe4)
	} else {
		fmt.Printf("couldn't create recipe: found %s\n", err)
	}

	ids := recipeStore.List()
	for _, id := range ids {
		currentRecipe, _ := recipeStore.Get(id)
		fmt.Printf("Recipe #%d. ", id)
		currentRecipe.DisplayRecipe()
	}

	onionRecipes, err := recipeStore.FindByIngredient("Onion")
	if err == nil {
		fmt.Println(onionRecipes)
	} else {
		fmt.Printf("couldn't find ingredient %s\n", err)
	}

	saltRecipes, err := recipeStore.FindByIngredient("   salt   ")
	if err == nil {
		fmt.Println(saltRecipes)
	} else {
		fmt.Printf("couldn't find ingredient %s\n", err)
	}

	garlicRecipes, err := recipeStore.FindByIngredient(" garlic")
	if err == nil {
		fmt.Println(garlicRecipes)
	} else {
		fmt.Printf("couldn't find ingredient %s\n", err)
	}

	emptyRecipes, err := recipeStore.FindByIngredient("   ")
	if err == nil {
		fmt.Println(emptyRecipes)
	} else {
		fmt.Printf("couldn't find ingredient %s\n", err)
	}

	recipes, err := recipeStore.FindByMaxMinutes(19)
	if err != nil {
		fmt.Printf("couldn't find recipe: %s\n", err)
	} else {
		fmt.Println(recipes)
	}

	recipes, err = recipeStore.FindByMaxMinutes(20)
	if err != nil {
		fmt.Printf("couldn't find recipe: %s\n", err)
	} else {
		fmt.Println(recipes)
	}

	recipes, err = recipeStore.FindByMaxMinutes(15)
	if err != nil {
		fmt.Printf("couldn't find recipe: %s\n", err)
	} else {
		fmt.Println(recipes)
	}

	recipes, err = recipeStore.FindByMaxMinutes(-10)
	if err != nil {
		fmt.Printf("couldn't find recipe: %s\n", err)
	} else {
		fmt.Println(recipes)
	}

	doroRecipes, err := recipeStore.FindByTitle("doro")
	if err != nil {
		fmt.Printf("couldn't find recipe: %s\n", err)
	} else {
		fmt.Println(doroRecipes)
	}

	doroWotRecipes, err := recipeStore.FindByTitle("WOT")
	if err != nil {
		fmt.Printf("couldn't find recipe: %s\n", err)
	} else {
		fmt.Println(doroWotRecipes)
	}

	lasagnaRecipes, err := recipeStore.FindByTitle("  lasagna   ")
	if err != nil {
		fmt.Printf("couldn't find recipe: %s\n", err)
	} else {
		fmt.Println(lasagnaRecipes)
	}

	firfirRecipes, err := recipeStore.FindByTitle("firfir")
	if err != nil {
		fmt.Printf("couldn't find recipe: %s\n", err)
	} else {
		fmt.Println(firfirRecipes)
	}

	emptyTitle, err := recipeStore.FindByTitle("    ")
	if err != nil {
		fmt.Printf("couldn't find recipe: %s\n", err)
	} else {
		fmt.Println(emptyTitle)
	}
	err = recipeStore.Delete(5)
	if err != nil {
		fmt.Printf("couldn't delete recipe %s\n", err)
	}

	categoryCount := recipeStore.CountByCategory()
	for category, count := range categoryCount {
		fmt.Printf("Category count of %s is %d\n", category, count)
	}
}
