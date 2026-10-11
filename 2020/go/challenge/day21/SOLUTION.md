# Day 21: Allergen Assessment

## Part A
For each allergen, intersect the ingredient sets of every food that lists it, leaving
the set of ingredients that could carry that allergen. Count ingredients across all
foods that never appear in any allergen's candidate set.

Complexity: O(F * I), F foods, I ingredients per food.

## Part B
Repeatedly find an allergen whose candidate set has exactly one ingredient left, assign
it, and remove that ingredient from every other allergen's candidate set, until all
allergens are resolved. Sort assignments by allergen name and join ingredients with
commas.

Complexity: O(A^2), A allergens (small).

## Commands
```bash
go test -race -count=1 ./challenge/day21
go run . 21 < ../rust/input/2020/day21-example.txt
```
