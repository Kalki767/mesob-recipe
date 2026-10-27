package recipe

import (
	"errors"
	"fmt"
	"strings"
)

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
	foundRecipe, ok := r.Recipes[id]
	if !ok {
		return nil, errors.New("recipe was not found")
	}

	return foundRecipe, nil
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

func (r *RecipeStore) filter(match func(*Recipe) bool) []int {
	recipeIds := r.List()
	matchedRecipes := []int{}
	for _, id := range recipeIds {
		recipe := r.Recipes[id]
		if match(recipe) {
			matchedRecipes = append(matchedRecipes, id)
		}
	}
	return matchedRecipes
}

func (r *RecipeStore) FindByIngredient(name string) ([]int, error) {
	name = strings.ToLower(name)
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("ingredient is empty")
	}

	return r.filter(func(recipe *Recipe) bool {
		for _, ingredient := range recipe.Ingredients {
			if strings.Contains(ingredient.Name, name) {
				return true
			}
		}
		return false
	}), nil
}

func (r *RecipeStore) FindByMaxMinutes(maxMinutes int) ([]int, error) {
	if maxMinutes <= 0 {
		return nil, fmt.Errorf("max minutes must be positive, got: %d", maxMinutes)
	}
	return r.filter(func(recipe *Recipe) bool {
		return recipe.TotalMinutes <= maxMinutes
	}), nil

}

func (r *RecipeStore) FindByTitle(title string) ([]int, error) {
	title = strings.ToLower(title)
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("title is empty")
	}

	return r.filter(func(recipe *Recipe) bool {
		storedTitle := strings.ToLower(recipe.Title)
		return strings.Contains(storedTitle, title)
	}), nil
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