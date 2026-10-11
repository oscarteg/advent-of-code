package day21

import (
	"fmt"
	"sort"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func resolveAllergens(candidates map[string]map[string]bool) (map[string]string, error) {
	resolved := make(map[string]string)

	for len(candidates) > 0 {
		allergen, ingredient, found := singleCandidate(candidates)
		if !found {
			return nil, fmt.Errorf("ambiguous allergen assignment: no allergen has a single candidate left")
		}

		resolved[allergen] = ingredient
		delete(candidates, allergen)
		for _, set := range candidates {
			delete(set, ingredient)
		}
	}

	return resolved, nil
}

func singleCandidate(candidates map[string]map[string]bool) (allergen, ingredient string, found bool) {
	allergens := make([]string, 0, len(candidates))
	for a := range candidates {
		allergens = append(allergens, a)
	}
	sort.Strings(allergens)

	for _, a := range allergens {
		if set := candidates[a]; len(set) == 1 {
			for ingredient := range set {
				return a, ingredient, true
			}
		}
	}
	return "", "", false
}

func B(input *challenge.Input) (string, error) {
	foods, err := parseFoods(input.LinesSlice())
	if err != nil {
		return "", err
	}
	if len(foods) == 0 {
		return "", fmt.Errorf("no foods found in input")
	}

	candidates := candidatesByAllergen(foods)
	resolved, err := resolveAllergens(candidates)
	if err != nil {
		return "", err
	}

	allergens := make([]string, 0, len(resolved))
	for allergen := range resolved {
		allergens = append(allergens, allergen)
	}
	sort.Strings(allergens)

	ingredients := make([]string, 0, len(allergens))
	for _, allergen := range allergens {
		ingredients = append(ingredients, resolved[allergen])
	}

	return strings.Join(ingredients, ","), nil
}
