package day21

import (
	"fmt"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

type food struct {
	ingredients []string
	allergens   []string
}

func parseFoods(lines []string) ([]food, error) {
	foods := make([]food, 0, len(lines))
	for _, line := range lines {
		f, err := parseFood(line)
		if err != nil {
			return nil, err
		}
		foods = append(foods, f)
	}
	return foods, nil
}

func parseFood(line string) (food, error) {
	ingredientsPart, allergensPart, found := strings.Cut(line, " (contains ")
	if !found {
		return food{}, fmt.Errorf("missing '(contains ...)' section in %q", line)
	}
	allergensPart, ok := strings.CutSuffix(allergensPart, ")")
	if !ok {
		return food{}, fmt.Errorf("allergen list not closed with ')' in %q", line)
	}

	ingredients := strings.Fields(ingredientsPart)
	if len(ingredients) == 0 {
		return food{}, fmt.Errorf("no ingredients found in %q", line)
	}

	allergens := strings.Split(allergensPart, ", ")
	for _, allergen := range allergens {
		if allergen == "" {
			return food{}, fmt.Errorf("empty allergen name in %q", line)
		}
	}

	return food{ingredients: ingredients, allergens: allergens}, nil
}

func candidatesByAllergen(foods []food) map[string]map[string]bool {
	candidates := make(map[string]map[string]bool)

	for _, f := range foods {
		ingredientSet := make(map[string]bool, len(f.ingredients))
		for _, ingredient := range f.ingredients {
			ingredientSet[ingredient] = true
		}

		for _, allergen := range f.allergens {
			existing, ok := candidates[allergen]
			if !ok {
				copied := make(map[string]bool, len(ingredientSet))
				for ingredient := range ingredientSet {
					copied[ingredient] = true
				}
				candidates[allergen] = copied
				continue
			}
			for ingredient := range existing {
				if !ingredientSet[ingredient] {
					delete(existing, ingredient)
				}
			}
		}
	}

	return candidates
}

func A(input *challenge.Input) (int, error) {
	foods, err := parseFoods(input.LinesSlice())
	if err != nil {
		return 0, err
	}
	if len(foods) == 0 {
		return 0, fmt.Errorf("no foods found in input")
	}

	candidates := candidatesByAllergen(foods)
	resolved, err := resolveAllergens(candidates)
	if err != nil {
		return 0, err
	}

	dangerous := make(map[string]bool, len(resolved))
	for _, ingredient := range resolved {
		dangerous[ingredient] = true
	}

	count := 0
	for _, f := range foods {
		for _, ingredient := range f.ingredients {
			if !dangerous[ingredient] {
				count++
			}
		}
	}

	return count, nil
}
