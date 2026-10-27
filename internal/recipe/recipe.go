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

func validateStep(step Step) error {
	if step.Task == "" {
		return errors.New("step task is empty")
	}

	if step.Minutes <= 0 {
		return fmt.Errorf("step minutes must be positive, got %d", step.Minutes)
	}

	return nil
}

func cleanSteps(steps []Step) error {
	for i := range steps {
		steps[i].Task = strings.TrimSpace(steps[i].Task)
		if err := validateStep(steps[i]); err != nil {
			return err
		}
	}

	return nil
}

// ============================================================================
// Recipe: Model, Constructor & Methods
// ============================================================================

type Recipe struct {
	title         string
	priceInSantim int
	steps         []Step
	totalMinutes  int
	category      string
	ingredients   []Ingredient
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

	if err := cleanSteps(steps); err != nil {
		return nil, err
	}

	if len(ingredients) == 0 {
		return nil, errors.New("ingredients cannot be empty")
	}
	if err := cleanIngredients(ingredients); err != nil {
		return nil, err
	}
	recipe := &Recipe{title: title, priceInSantim: priceInSantim, steps: steps, ingredients: ingredients}
	recipe.totalMinutes = recipe.calculateMinutes()
	recipe.category = category
	return recipe, nil
}

func cleanIngredients(ingredients []Ingredient) error {
	for i := range ingredients {
		ingredients[i].Name = strings.TrimSpace(ingredients[i].Name)
		if err := validateIngredient(ingredients[i]); err != nil {
			return err
		}
		ingredients[i].Name = strings.ToLower(ingredients[i].Name)
	}
	return nil
}

func validateIngredient(ingredient Ingredient) error {
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
	err := validateIngredient(ingredient)
	ingredient.Name = strings.ToLower(ingredient.Name)
	if err != nil {
		return err
	}
	r.ingredients = append(r.ingredients, ingredient)
	return nil
}

func (r *Recipe) AddStep(step Step) error {
	step.Task = strings.TrimSpace(step.Task)
	if err := validateStep(step); err != nil {
		return err
	}
	r.steps = append(r.steps, step)
	r.totalMinutes += step.Minutes
	return nil
}

func (r *Recipe) calculateMinutes() int {
	minutes := 0
	for _, step := range r.steps {
		minutes += step.Minutes
	}
	return minutes
}

func (r *Recipe) DisplayRecipe() {
	fmt.Printf("%s costs %s birr\n", r.title, formatPrice(r.priceInSantim))
	for i, ingredient := range r.ingredients {
		if !ingredient.ToTaste {
			fmt.Printf("%d: %.02f %s %s\n", i+1, ingredient.Amount, ingredient.Unit, ingredient.Name)
		} else {
			fmt.Printf("%d: Add an amount you prefer %s\n", i+1, ingredient.Name)
		}
	}
	for i, step := range r.steps {
		fmt.Printf("%d. %s (%d min)\n", i+1, step.Task, step.Minutes)
	}
	fmt.Printf("Category: %s\n", r.category)
	hours, remainingMinutes := splitDuration(r.totalMinutes)
	fmt.Printf("Total: %d hour %d minutes\n", hours, remainingMinutes)
}

// ============================================================================
// RecipeStore: Storage & Operations
// ============================================================================

type RecipeStore struct {
	nextID  int
	recipes map[int]*Recipe
}

func NewRecipeStore() *RecipeStore {
	return &RecipeStore{nextID: 1, recipes: make(map[int]*Recipe)}
}

func (r *RecipeStore) Add(recipe *Recipe) int {
	r.recipes[r.nextID] = recipe
	r.nextID += 1
	return r.nextID - 1
}

func (r *RecipeStore) Get(id int) (*Recipe, error) {
	foundRecipe, ok := r.recipes[id]
	if !ok {
		return nil, errors.New("recipe was not found")
	}

	return foundRecipe, nil
}

func (r *RecipeStore) Delete(id int) error {
	if _, ok := r.recipes[id]; !ok {
		return errors.New("can't delete non existing recipe")
	}
	delete(r.recipes, id)
	return nil
}

func (r *RecipeStore) List() []int {
	recipes := []int{}
	start := 1
	for start <= r.nextID {
		if _, ok := r.recipes[start]; ok {
			recipes = append(recipes, start)
		}
		start += 1
	}

	return recipes
}

func (r *RecipeStore) CountByCategory() map[string]int {
	count := make(map[string]int)

	for _, recipe := range r.recipes {
		count[recipe.category] += 1
	}
	return count
}

func (r *RecipeStore) filter(match func(*Recipe) bool) []int {
	recipeIds := r.List()
	matchedRecipes := []int{}
	for _, id := range recipeIds {
		recipe := r.recipes[id]
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
		for _, ingredient := range recipe.ingredients {
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
		return recipe.totalMinutes <= maxMinutes
	}), nil

}

func (r *RecipeStore) FindByTitle(title string) ([]int, error) {
	title = strings.ToLower(title)
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("title is empty")
	}

	return r.filter(func(recipe *Recipe) bool {
		storedTitle := strings.ToLower(recipe.title)
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