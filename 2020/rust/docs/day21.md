# Day 21: Allergen Assessment

## Part 1
For each allergen, intersect the ingredient sets of every food that lists it; this leaves
the set of ingredients that could possibly carry that allergen. Union all remaining
candidate ingredients across allergens, then count how many times ingredients *not* in
that union appear across all foods' ingredient lists.

Complexity: O(F * I) to build candidate sets, where F is the number of foods and I the
average ingredients per food.

## Part 2
Resolve allergens one at a time: while any allergen's candidate set has exactly one
ingredient left, assign it and remove that ingredient from every other allergen's
candidate set. This narrows every allergen to a single ingredient since each ingredient
maps to at most one allergen. Sort assignments alphabetically by allergen name and join
the ingredients with commas.

Complexity: O(A^2) in the worst case, where A is the number of allergens (small).

## Commands
```bash
cargo test --bin day21
cargo run -q --bin day21 < input/2020/day21-example.txt
```
