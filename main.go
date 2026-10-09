package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Step struct {
	Task    string
	Minutes int
}

type Recipe struct {
	Title         string
	PriceInSantim int
	Steps         []Step
	TotalMinutes  int
	Category      string
}

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

const santimPerBirr = 100

func (r *Recipe) calculateMinutes() int {
	minutes := 0
	for _, step := range r.Steps {
		minutes += step.Minutes
	}
	return minutes
}

func formatPrice(priceInSantim int) string {
	return fmt.Sprintf("%d.%02d", priceInSantim/santimPerBirr, priceInSantim%santimPerBirr)
}

func splitDuration(minutes int) (int, int) {
	return minutes / 60, minutes % 60
}

func NewRecipe(title string, priceInSantim int, steps []Step, category string) (*Recipe, error) {
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
	recipe := &Recipe{Title: title, PriceInSantim: priceInSantim, Steps: steps}
	recipe.TotalMinutes = recipe.calculateMinutes()
	recipe.Category = category
	return recipe, nil
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

func CheckRecipeError(err error) {
	if err != nil {
		fmt.Println("couldn't create a recipe :", err)
		os.Exit(1)
	}
}

func (r *Recipe) DisplayRecipe() {
	fmt.Printf("%s costs %s birr\n", r.Title, formatPrice(r.PriceInSantim))
	for i, step := range r.Steps {
		fmt.Printf("%d. %s (%d min)\n", i+1, step.Task, step.Minutes)
	}
	fmt.Printf("Category: %s\n", r.Category)
	hours, remainingMinutes := splitDuration(r.TotalMinutes)
	fmt.Printf("Total: %d hour %d minutes\n", hours, remainingMinutes)
}
func main() {
	recipe1, err := NewRecipe(
		"Doro Wot",
		25005,
		[]Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 5},
		},
		"Habesha food",
	)

	CheckRecipeError(err)

	recipe2, err := NewRecipe(
		"Shiro Wot",
		25005,
		[]Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 5},
		},
		"Habesha food",
	)
	CheckRecipeError(err)
	recipe3, err := NewRecipe(
		"Lazagna",
		25005,
		[]Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 5},
		},
		"Foreign food",
	)
	CheckRecipeError(err)
	recipeStore := NewRecipeStore()
	recipeStore.Add(recipe1)
	recipeStore.Add(recipe2)
	recipeStore.Add(recipe3)

	recipeStore.Delete(2)
	recipe4, err := NewRecipe(
		"Cake",
		25005,
		[]Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 5},
		},
		"Foreign food",
	)
	CheckRecipeError(err)
	recipeStore.Add(recipe4)
	ids := recipeStore.List()
	for _, id := range ids {
		recipe, _ := recipeStore.Get(id)
		fmt.Printf("Recipe #%d. ", id)
		recipe.DisplayRecipe()
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
