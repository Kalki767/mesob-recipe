package main

import (
	"errors"
	"fmt"
	"strings"
)

// ============================================================================
// Constants
// ============================================================================

const santimPerBirr = 100

// ============================================================================
// Ingredient & Step: Models & Validation
// ============================================================================

type Ingredient struct {
	Name    string
	Amount  float64
	Unit    string
	ToTaste bool
}

type Step struct {
	Task    string
	Minutes int
}

func ValidateStep(step Step) error {
	if step.Task == "" {
		return errors.New("step task is empty")
	}

	if step.Minutes <= 0 {
		return fmt.Errorf("step minutes must be positive, got %d", step.Minutes)
	}

	return nil
}

func CleanSteps(steps []Step) error {
	for i := range steps {
		steps[i].Task = strings.TrimSpace(steps[i].Task)
		if err := ValidateStep(steps[i]); err != nil {
			return err
		}
	}

	return nil
}

// ============================================================================
// Recipe: Model, Constructor & Methods
// ============================================================================

type Recipe struct {
	Title         string
	PriceInSantim int
	Steps         []Step
	TotalMinutes  int
	Category      string
	Ingredients   []Ingredient
}

func NewRecipe(title string, priceInSantim int, steps []Step, category string, ingredients []Ingredient) (*Recipe, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	if priceInSantim < 0 {
		return nil, fmt.Errorf("price cannot be negative, got %d", priceInSantim)
	}

	category = strings.TrimSpace(category)
	if category == "" {
		return nil, errors.New("category cannot be empty")
	}

	if err := CleanSteps(steps); err != nil {
		return nil, err
	}

	if len(ingredients) == 0 {
		return nil, errors.New("ingredients cannot be empty")
	}
	if err := CleanIngredients(ingredients); err != nil {
		return nil, err
	}
	recipe := &Recipe{Title: title, PriceInSantim: priceInSantim, Steps: steps, Ingredients: ingredients}
	recipe.TotalMinutes = recipe.calculateMinutes()
	recipe.Category = category
	return recipe, nil
}

func CleanIngredients(ingredients []Ingredient) error {
	for i := range ingredients {
		ingredients[i].Name = strings.TrimSpace(ingredients[i].Name)
		if err := ValidateIngredient(ingredients[i]); err != nil {
			return err
		}
		ingredients[i].Name = strings.ToLower(ingredients[i].Name)
	}
	return nil
}

func ValidateIngredient(ingredient Ingredient) error {
	if ingredient.Name == "" {
		return errors.New("empty ingredient name is not allowed")
	}
	if ingredient.Amount < 0 {
		return fmt.Errorf("amount cannot be negative, got %.02f", ingredient.Amount)
	}

	if ingredient.Amount == 0 && !ingredient.ToTaste {
		return errors.New("ingredient needs amount or to taste")
	}

	if ingredient.Amount > 0 && ingredient.ToTaste {
		return errors.New("ingredient cannot have both amount and to taste")
	}
	return nil
}
func (r *Recipe) AddIngredient(ingredient Ingredient) error {
	ingredient.Name = strings.TrimSpace(ingredient.Name)
	err := ValidateIngredient(ingredient)
	ingredient.Name = strings.ToLower(ingredient.Name)
	if err != nil {
		return err
	}
	r.Ingredients = append(r.Ingredients, ingredient)
	return nil
}

func (r *Recipe) AddStep(step Step) error {
	step.Task = strings.TrimSpace(step.Task)
	if err := ValidateStep(step); err != nil {
		return err
	}
	r.Steps = append(r.Steps, step)
	r.TotalMinutes += step.Minutes
	return nil
}

func (r *Recipe) calculateMinutes() int {
	minutes := 0
	for _, step := range r.Steps {
		minutes += step.Minutes
	}
	return minutes
}

func (r *Recipe) DisplayRecipe() {
	fmt.Printf("%s costs %s birr\n", r.Title, formatPrice(r.PriceInSantim))
	for i, ingredient := range r.Ingredients {
		if !ingredient.ToTaste {
			fmt.Printf("%d: %.02f %s %s\n", i+1, ingredient.Amount, ingredient.Unit, ingredient.Name)
		} else {
			fmt.Printf("%d: Add an amount you prefer %s\n", i+1, ingredient.Name)
		}
	}
	for i, step := range r.Steps {
		fmt.Printf("%d. %s (%d min)\n", i+1, step.Task, step.Minutes)
	}
	fmt.Printf("Category: %s\n", r.Category)
	hours, remainingMinutes := splitDuration(r.TotalMinutes)
	fmt.Printf("Total: %d hour %d minutes\n", hours, remainingMinutes)
}

// ============================================================================
// RecipeStore: Storage & Operations
// ============================================================================

type RecipeStore struct {
	NextID  int
	Recipes map[int]*Recipe
}

func NewRecipeStore() *RecipeStore {
	return &RecipeStore{NextID: 1, Recipes: make(map[int]*Recipe)}
}

func (r *RecipeStore) Add(recipe *Recipe) int {
	r.Recipes[r.NextID] = recipe
	r.NextID += 1
	return r.NextID - 1
}

func (r *RecipeStore) Get(id int) (*Recipe, error) {
	recipe, ok := r.Recipes[id]
	if !ok {
		return nil, errors.New("recipe was not found")
	}

	return recipe, nil
}

func (r *RecipeStore) Delete(id int) error {
	if _, ok := r.Recipes[id]; !ok {
		return errors.New("can't delete non existing recipe")
	}
	delete(r.Recipes, id)
	return nil
}

func (r *RecipeStore) List() []int {
	recipes := []int{}
	start := 1
	for start <= r.NextID {
		if _, ok := r.Recipes[start]; ok {
			recipes = append(recipes, start)
		}
		start += 1
	}

	return recipes
}

func (r *RecipeStore) CountByCategory() map[string]int {
	count := make(map[string]int)

	for _, recipe := range r.Recipes {
		count[recipe.Category] += 1
	}
	return count
}

func (r *RecipeStore) FindByIngredient(name string) ([]int, error) {
	name = strings.ToLower(name)
	name = strings.TrimSpace(name)
	recipes := r.List()
	var matchedRecipes []int = []int{}
	if name == "" {
		return nil, errors.New("ingredient is empty")
	}
	for _, id := range recipes {
		recipe := r.Recipes[id]
		for _, ingredient := range recipe.Ingredients {
			if strings.Contains(ingredient.Name, name) {
				matchedRecipes = append(matchedRecipes, id)
				break
			}
		}
	}

	return matchedRecipes, nil
}

// ============================================================================
// General Utilities & Formatting Helpers
// ============================================================================

func formatPrice(priceInSantim int) string {
	return fmt.Sprintf("%d.%02d", priceInSantim/santimPerBirr, priceInSantim%santimPerBirr)
}

func splitDuration(minutes int) (int, int) {
	return minutes / 60, minutes % 60
}

// ============================================================================
// Main Execution
// ============================================================================

func main() {
	recipeStore := NewRecipeStore()
	recipe1, err := NewRecipe(
		"Doro Wot",
		25005,
		[]Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 5},
		},
		"Habesha food",
		[]Ingredient{
			{Name: "  Onion  ", Amount: 3, Unit: ""},
		},
	)

	if err == nil {
		err = recipe1.AddIngredient(Ingredient{Name: "  oil  ", Amount: 1, Unit: "Cups"})
		if err != nil {
			fmt.Printf("couldn't add ingredient, found %s\n", err)
		}
		recipeStore.Add(recipe1)
	} else {
		fmt.Printf("couldn't create recipe: found %s\n", err)
	}

	recipe2, err := NewRecipe(
		"Shiro Wot",
		25005,
		[]Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 5},
		},
		"Habesha food",
		[]Ingredient{
			{Name: "onion", Amount: 2},
		},
	)
	if err == nil {
		err = recipe2.AddIngredient(Ingredient{Name: "Onion", Amount: 0, Unit: ""})
		if err != nil {
			fmt.Printf("couldn't add ingredient, found %s\n", err)
		}
		recipeStore.Add(recipe2)
	} else {
		fmt.Printf("couldn't create recipe: found %s\n", err)
	}

	recipe3, err := NewRecipe(
		"Lazagna",
		25005,
		[]Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 5},
		},
		"Foreign food",
		[]Ingredient{
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
	recipe4, err := NewRecipe(
		"Cake",
		25005,
		[]Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 5},
		},
		"Foreign food",
		[]Ingredient{},
	)
	if err == nil {
		recipeStore.Add(recipe4)
	} else {
		fmt.Printf("couldn't create recipe: found %s\n", err)
	}

	ids := recipeStore.List()
	for _, id := range ids {
		recipe, _ := recipeStore.Get(id)
		fmt.Printf("Recipe #%d. ", id)
		recipe.DisplayRecipe()
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

	err = recipeStore.Delete(5)
	if err != nil {
		fmt.Printf("couldn't delete recipe %s\n", err)
	}

	categoryCount := recipeStore.CountByCategory()
	for category, count := range categoryCount {
		fmt.Printf("Category count of %s is %d\n", category, count)
	}
}
