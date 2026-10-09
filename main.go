package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
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
	recipes := make(map[int]*Recipe)
	recipes[1] = recipe1
	recipes[2] = recipe2
	recipes[3] = recipe3

	ids := []int{}

	count := make(map[string]int)

	for i, recipe := range recipes {
		count[recipe.Category] += 1
		ids = append(ids, i)
	}

	slices.Sort(ids)

	for _, id := range ids {
		recipe := recipes[id]
		recipe.DisplayRecipe()
	}
	for key, value := range count {
		fmt.Printf("Category count of %s: %d\n", key, value)
	}
	recipe, ok := recipes[3]
	if ok {
		fmt.Println("Recipe found")
		recipe.DisplayRecipe()
	} else {
		fmt.Println("Recipe not found")
	}

	delete(recipes, 3)

	recipe, ok = recipes[3]
	if ok {
		fmt.Println("Recipe found")
		recipe.DisplayRecipe()
	} else {
		fmt.Println("Recipe not found")
	}
}
