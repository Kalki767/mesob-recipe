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

func NewRecipe(title string, priceInSantim int, steps []Step) (*Recipe, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	if priceInSantim < 0 {
		return nil, fmt.Errorf("price cannot be negative got %d", priceInSantim)
	}
	if err := CleanSteps(steps); err != nil {
		return nil, err
	}
	recipe := &Recipe{Title: title, PriceInSantim: priceInSantim, Steps: steps}
	recipe.TotalMinutes = recipe.calculateMinutes()
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
		return fmt.Errorf("step minutes must be positive got %d", step.Minutes)
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

func main() {
	recipe, err := NewRecipe(
		"  ",
		25005,
		[]Step{
			{Task: "  Chop the onion   ", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 5},
		},
	)

	if err != nil {
		fmt.Println("couldn't create a recipe :", err)
		os.Exit(1)
	}

	fmt.Printf("%s costs %s birr\n", recipe.Title, formatPrice(recipe.PriceInSantim))

	if err := recipe.AddStep(Step{Task: "  Serve with injera", Minutes: -5}); err != nil {
		fmt.Printf("%s\n", err)
	}

	if err := recipe.AddStep(Step{Task: "    ", Minutes: 2}); err != nil {
		fmt.Printf("%s\n", err)
	}

	for i, s := range recipe.Steps {
		fmt.Printf("%d. %s (%d min)\n", i+1, s.Task, s.Minutes)
	}

	hours, remainingMinutes := splitDuration(recipe.TotalMinutes)
	fmt.Printf("Total: %d hour %d minutes\n", hours, remainingMinutes)
}
