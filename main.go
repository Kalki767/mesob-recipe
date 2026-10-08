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

func totalMinutes(steps []Step) int {
	minutes := 0
	for _, step := range steps {
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

func addStep(recipe *Recipe, step Step) {
	recipe.Steps = append(recipe.Steps, step)
	recipe.TotalMinutes += step.Minutes
}

func main() {
	recipe := Recipe{
		Title:         "Doro wot",
		PriceInSantim: 25005,
		Steps: []Step{
			{Task: "Chop the onion", Minutes: 10},
			{Task: "Cook it with oil", Minutes: 5},
			{Task: "Mix it with egg", Minutes: 5},
		},
	}

	fmt.Printf("%s costs %s birr\n", recipe.Title, formatPrice(recipe.PriceInSantim))
	recipe.TotalMinutes = totalMinutes(recipe.Steps)

	addStep(&recipe, Step{Task: "Serve with injera", Minutes: 2})

	for i, s := range recipe.Steps {
		fmt.Printf("%d. %s (%d min)\n", i+1, s.Task, s.Minutes)
	}

	hours, remainingMinutes := splitDuration(recipe.TotalMinutes)
	fmt.Printf("Total: %d hour %d minutes\n", hours, remainingMinutes)
}
