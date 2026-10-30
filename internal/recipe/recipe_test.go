package recipe

import (
	"slices"
	"testing"
)

// Fresh test data for every call. NewRecipe trims and stores the slices it is
// given, so sharing one slice between rows would let one row's changes leak
// into another.
func validSteps() []Step {
	return []Step{
		{Task: "Chop the onion", Minutes: 10},
		{Task: "Cook it with oil", Minutes: 5},
	}
}

func validIngredients() []Ingredient {
	return []Ingredient{
		{Name: "onion", Amount: 2},
	}
}

func TestNewRecipe(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		price       int
		steps       []Step
		category    string
		ingredients []Ingredient
		wantErr     bool
	}{
		{"valid recipe", "Doro Wot", 25005, validSteps(), "Habesha food", validIngredients(), false},
		{"free recipe is allowed", "Shiro Wot", 0, validSteps(), "Habesha food", validIngredients(), false},
		{"spaces around text are fine", "  Doro Wot  ", 100, validSteps(), "  Habesha food ", validIngredients(), false},

		{"empty title", "", 100, validSteps(), "Habesha food", validIngredients(), true},
		{"title of only spaces", "   ", 100, validSteps(), "Habesha food", validIngredients(), true},
		{"negative price", "Doro Wot", -1, validSteps(), "Habesha food", validIngredients(), true},
		{"empty category", "Doro Wot", 100, validSteps(), "", validIngredients(), true},
		{"category of only spaces", "Doro Wot", 100, validSteps(), "   ", validIngredients(), true},

		{"step with empty task", "Doro Wot", 100, []Step{{Task: "  ", Minutes: 5}}, "Habesha food", validIngredients(), true},
		{"step with zero minutes", "Doro Wot", 100, []Step{{Task: "Chop", Minutes: 0}}, "Habesha food", validIngredients(), true},
		{"step with negative minutes", "Doro Wot", 100, []Step{{Task: "Chop", Minutes: -5}}, "Habesha food", validIngredients(), true},

		{"no ingredients", "Doro Wot", 100, validSteps(), "Habesha food", []Ingredient{}, true},
		{"nil ingredients", "Doro Wot", 100, validSteps(), "Habesha food", nil, true},
		{"ingredient with empty name", "Doro Wot", 100, validSteps(), "Habesha food", []Ingredient{{Name: " ", Amount: 1}}, true},
		{"ingredient with negative amount", "Doro Wot", 100, validSteps(), "Habesha food", []Ingredient{{Name: "salt", Amount: -1}}, true},
		{"ingredient with no amount and not to taste", "Doro Wot", 100, validSteps(), "Habesha food", []Ingredient{{Name: "salt"}}, true},
		{"ingredient with amount and to taste", "Doro Wot", 100, validSteps(), "Habesha food", []Ingredient{{Name: "salt", Amount: 1, ToTaste: true}}, true},
		{"ingredient to taste", "Doro Wot", 100, validSteps(), "Habesha food", []Ingredient{{Name: "salt", ToTaste: true}}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewRecipe(tc.title, tc.price, tc.steps, tc.category, tc.ingredients)
			if tc.wantErr && err == nil {
				t.Errorf("expected an error, got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestNewRecipeCleansInput(t *testing.T) {
	r, err := NewRecipe("  Doro Wot  ", 100,
		[]Step{{Task: "  Chop the onion  ", Minutes: 10}, {Task: "Cook", Minutes: 12}},
		"Habesha food",
		[]Ingredient{{Name: "  Onion  ", Amount: 3}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.title != "Doro Wot" {
		t.Errorf("title = %q, want %q", r.title, "Doro Wot")
	}
	if r.steps[0].Task != "Chop the onion" {
		t.Errorf("step task = %q, want %q", r.steps[0].Task, "Chop the onion")
	}
	if r.ingredients[0].Name != "onion" {
		t.Errorf("ingredient name = %q, want %q", r.ingredients[0].Name, "onion")
	}
	if r.totalMinutes != 22 {
		t.Errorf("totalMinutes = %d, want 22", r.totalMinutes)
	}
}

func TestAddStep(t *testing.T) {
	tests := []struct {
		name    string
		step    Step
		wantErr bool
	}{
		{"valid step", Step{Task: "Serve", Minutes: 3}, false},
		{"empty task", Step{Task: "   ", Minutes: 3}, true},
		{"zero minutes", Step{Task: "Serve", Minutes: 0}, true},
		{"negative minutes", Step{Task: "Serve", Minutes: -3}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewRecipe("Doro Wot", 100, validSteps(), "Habesha food", validIngredients())
			if err != nil {
				t.Fatalf("setup: %v", err)
			}
			stepsBefore, minutesBefore := len(r.steps), r.totalMinutes

			err = r.AddStep(tc.step)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got none")
				}
				// A rejected step must leave the recipe untouched.
				if len(r.steps) != stepsBefore || r.totalMinutes != minutesBefore {
					t.Errorf("recipe changed after a rejected step: %d steps, %d min (before: %d steps, %d min)",
						len(r.steps), r.totalMinutes, stepsBefore, minutesBefore)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(r.steps) != stepsBefore+1 {
				t.Errorf("got %d steps, want %d", len(r.steps), stepsBefore+1)
			}
			if r.totalMinutes != minutesBefore+tc.step.Minutes {
				t.Errorf("totalMinutes = %d, want %d", r.totalMinutes, minutesBefore+tc.step.Minutes)
			}
		})
	}
}

func TestAddIngredient(t *testing.T) {
	tests := []struct {
		name       string
		ingredient Ingredient
		wantErr    bool
		wantName   string
	}{
		{"valid ingredient is trimmed and lowercased", Ingredient{Name: "  Garlic ", Amount: 2}, false, "garlic"},
		{"to taste", Ingredient{Name: "Salt", ToTaste: true}, false, "salt"},
		{"empty name", Ingredient{Name: "  ", Amount: 2}, true, ""},
		{"negative amount", Ingredient{Name: "salt", Amount: -2}, true, ""},
		{"no amount and not to taste", Ingredient{Name: "salt"}, true, ""},
		{"amount and to taste", Ingredient{Name: "salt", Amount: 1, ToTaste: true}, true, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewRecipe("Doro Wot", 100, validSteps(), "Habesha food", validIngredients())
			if err != nil {
				t.Fatalf("setup: %v", err)
			}
			countBefore := len(r.ingredients)

			err = r.AddIngredient(tc.ingredient)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got none")
				}
				if len(r.ingredients) != countBefore {
					t.Errorf("ingredient was added despite the error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(r.ingredients) != countBefore+1 {
				t.Fatalf("got %d ingredients, want %d", len(r.ingredients), countBefore+1)
			}
			if got := r.ingredients[len(r.ingredients)-1].Name; got != tc.wantName {
				t.Errorf("stored name = %q, want %q", got, tc.wantName)
			}
		})
	}
}

// newTestStore builds a store with three known recipes:
//
//	1: Doro Wot  (22 min, Habesha food) onion, oil
//	2: Shiro Wot (15 min, Habesha food) chickpea flour, onion
//	3: Lasagna   (19 min, Foreign food) pasta, salt (to taste)
func newTestStore(t *testing.T) *RecipeStore {
	t.Helper() // failures point at the test that called this, not at this line

	add := func(store *RecipeStore, title string, minutes []int, category string, ingredients []Ingredient) {
		steps := []Step{}
		for i, m := range minutes {
			steps = append(steps, Step{Task: "step " + string(rune('A'+i)), Minutes: m})
		}
		r, err := NewRecipe(title, 100, steps, category, ingredients)
		if err != nil {
			t.Fatalf("setup %q: %v", title, err)
		}
		store.Add(r)
	}

	store := NewRecipeStore()
	add(store, "Doro Wot", []int{10, 5, 7}, "Habesha food",
		[]Ingredient{{Name: "Onion", Amount: 3}, {Name: "oil", Amount: 1, Unit: "cup"}})
	add(store, "Shiro Wot", []int{10, 5}, "Habesha food",
		[]Ingredient{{Name: "Chickpea flour", Amount: 2, Unit: "cup"}, {Name: "onion", Amount: 1}})
	add(store, "Lasagna", []int{10, 9}, "Foreign food",
		[]Ingredient{{Name: "pasta", Amount: 500, Unit: "g"}, {Name: "salt", ToTaste: true}})
	return store
}

func TestFindByIngredient(t *testing.T) {
	store := newTestStore(t)

	tests := []struct {
		name    string
		query   string
		want    []int
		wantErr bool
	}{
		{"exact name", "oil", []int{1}, false},
		{"several recipes, in ID order", "onion", []int{1, 2}, false},
		{"case-insensitive", "ONION", []int{1, 2}, false},
		{"surrounding spaces ignored", "  salt  ", []int{3}, false},
		{"partial match", "chick", []int{2}, false},
		{"no match is empty, not nil", "garlic", []int{}, false},
		{"blank query", "   ", nil, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := store.FindByIngredient(tc.query)
			checkIDs(t, got, err, tc.want, tc.wantErr)
		})
	}
}

func TestFindByMaxMinutes(t *testing.T) {
	store := newTestStore(t)

	tests := []struct {
		name    string
		max     int
		want    []int
		wantErr bool
	}{
		{"exactly a recipe's total is included", 19, []int{2, 3}, false},
		{"everything fits", 60, []int{1, 2, 3}, false},
		{"only the shortest", 15, []int{2}, false},
		{"nothing fits", 10, []int{}, false},
		{"zero is invalid", 0, nil, true},
		{"negative is invalid", -10, nil, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := store.FindByMaxMinutes(tc.max)
			checkIDs(t, got, err, tc.want, tc.wantErr)
		})
	}
}

func TestFindByTitle(t *testing.T) {
	store := newTestStore(t)

	tests := []struct {
		name    string
		query   string
		want    []int
		wantErr bool
	}{
		{"partial word", "doro", []int{1}, false},
		{"matches several, in ID order", "WOT", []int{1, 2}, false},
		{"surrounding spaces ignored", "  lasagna  ", []int{3}, false},
		{"no match", "firfir", []int{}, false},
		{"blank query", "   ", nil, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := store.FindByTitle(tc.query)
			checkIDs(t, got, err, tc.want, tc.wantErr)
		})
	}
}

// checkIDs is the shared assert for the Find* tests.
func checkIDs(t *testing.T, got []int, err error, want []int, wantErr bool) {
	t.Helper()
	if wantErr {
		if err == nil {
			t.Errorf("expected an error, got none")
		}
		return
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Errorf("got nil, want a non-nil slice (it becomes null instead of [] in JSON)")
	}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestStoreIDsAreNeverReused(t *testing.T) {
	store := newTestStore(t)

	if err := store.Delete(2); err != nil {
		t.Fatalf("delete 2: %v", err)
	}
	r, err := NewRecipe("Cake", 100, validSteps(), "Dessert", validIngredients())
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	if id := store.Add(r); id != 4 {
		t.Errorf("new recipe got ID %d, want 4 (IDs must never be reused)", id)
	}
	if got := store.List(); !slices.Equal(got, []int{1, 3, 4}) {
		t.Errorf("List() = %v, want [1 3 4]", got)
	}
}

func TestStoreGetAndDeleteMissing(t *testing.T) {
	store := newTestStore(t)

	if _, err := store.Get(99); err == nil {
		t.Errorf("Get(99): expected an error, got none")
	}
	if err := store.Delete(99); err == nil {
		t.Errorf("Delete(99): expected an error, got none")
	}
	if err := store.Delete(1); err != nil {
		t.Fatalf("Delete(1): %v", err)
	}
	if _, err := store.Get(1); err == nil {
		t.Errorf("Get(1) after delete: expected an error, got none")
	}
}

func TestCountByCategory(t *testing.T) {
	store := newTestStore(t)

	got := store.CountByCategory()
	want := map[string]int{"Habesha food": 2, "Foreign food": 1}

	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for category, count := range want {
		if got[category] != count {
			t.Errorf("%q: got %d, want %d", category, got[category], count)
		}
	}
}