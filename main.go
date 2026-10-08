package main

import (
	"fmt"
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

func (r *Recipe) AddStep(step Step) {
	r.Steps = append(r.Steps, step)
	r.TotalMinutes += step.Minutes
}

func NewRecipe(title string, priceinsSantim int, steps []Step) *Recipe {
	recipe := &Recipe{Title: title, PriceInSantim: priceinsSantim, Steps: steps}
	recipe.TotalMinutes = recipe.calculateMinutes()
	return recipe
}

func main() {
	recipe := NewRecipe(
		"Doro wot",
		25005,
		[]Step{
			{Task: "Chop the onion", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 5},
		},
	)

	fmt.Printf("%s costs %s birr\n", recipe.Title, formatPrice(recipe.PriceInSantim))

	recipe.AddStep(Step{Task: "Serve with injera", Minutes: 2})

	for i, s := range recipe.Steps {
		fmt.Printf("%d. %s (%d min)\n", i+1, s.Task, s.Minutes)
	}

	hours, remainingMinutes := splitDuration(recipe.TotalMinutes)
	fmt.Printf("Total: %d hour %d minutes\n", hours, remainingMinutes)
}
