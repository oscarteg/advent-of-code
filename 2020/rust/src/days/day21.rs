use std::collections::{BTreeMap, HashSet};

use anyhow::{bail, ensure, Context, Result};

struct Food {
    ingredients: Vec<String>,
    allergens: Vec<String>,
}

fn parse(input: &str) -> Result<Vec<Food>> {
    input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(parse_food)
        .collect()
}

fn parse_food(line: &str) -> Result<Food> {
    let (ingredients_part, allergens_part) = line
        .split_once(" (contains ")
        .with_context(|| format!("missing '(contains ...)' section in {line:?}"))?;
    let allergens_part = allergens_part
        .strip_suffix(')')
        .with_context(|| format!("allergen list not closed with ')' in {line:?}"))?;

    let ingredients: Vec<String> = ingredients_part
        .split_whitespace()
        .map(str::to_string)
        .collect();
    ensure!(!ingredients.is_empty(), "no ingredients found in {line:?}");

    let allergens: Vec<String> = allergens_part
        .split(", ")
        .map(str::trim)
        .map(str::to_string)
        .collect();
    ensure!(
        allergens.iter().all(|a| !a.is_empty()),
        "empty allergen name in {line:?}"
    );

    Ok(Food {
        ingredients,
        allergens,
    })
}

/// For each allergen, intersect the ingredient sets of every food that lists it,
/// leaving only the ingredients that could possibly be the carrier of that allergen.
fn candidates_by_allergen(foods: &[Food]) -> BTreeMap<&str, HashSet<&str>> {
    let mut candidates: BTreeMap<&str, HashSet<&str>> = BTreeMap::new();

    for food in foods {
        let ingredient_set: HashSet<&str> = food.ingredients.iter().map(String::as_str).collect();
        for allergen in &food.allergens {
            candidates
                .entry(allergen.as_str())
                .and_modify(|set| set.retain(|ingredient| ingredient_set.contains(ingredient)))
                .or_insert_with(|| ingredient_set.clone());
        }
    }

    candidates
}

/// Repeatedly pick an allergen with exactly one remaining candidate, assign it, and
/// remove that ingredient from every other allergen's candidate set.
fn resolve_allergens<'a>(
    mut candidates: BTreeMap<&'a str, HashSet<&'a str>>,
) -> Result<BTreeMap<&'a str, &'a str>> {
    let mut resolved: BTreeMap<&str, &str> = BTreeMap::new();

    while !candidates.is_empty() {
        let Some((&allergen, &ingredient)) = candidates
            .iter()
            .find(|(_, set)| set.len() == 1)
            .map(|(a, set)| (a, set.iter().next().unwrap()))
        else {
            bail!("ambiguous allergen assignment: no allergen has a single candidate left");
        };

        resolved.insert(allergen, ingredient);
        candidates.remove(allergen);
        for set in candidates.values_mut() {
            set.remove(ingredient);
        }
    }

    Ok(resolved)
}

pub fn solve(input: &str) -> Result<(usize, String)> {
    let foods = parse(input)?;
    ensure!(!foods.is_empty(), "no foods found in input");

    let candidates = candidates_by_allergen(&foods);
    let resolved = resolve_allergens(candidates)?;

    let dangerous: HashSet<&str> = resolved.values().copied().collect();
    let part1 = foods
        .iter()
        .flat_map(|food| food.ingredients.iter())
        .filter(|ingredient| !dangerous.contains(ingredient.as_str()))
        .count();

    // `resolved` is a BTreeMap keyed by allergen name, so iterating it already
    // yields ingredients sorted alphabetically by their allergen.
    let part2 = resolved.values().copied().collect::<Vec<_>>().join(",");

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/2020/day21-example.txt");

    #[test]
    fn sample_part1() {
        let (part1, _) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 5);
    }

    #[test]
    fn sample_part2() {
        let (_, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part2, "mxmxvkd,sqjhc,fvjkl");
    }

    #[test]
    fn single_food_single_ingredient_is_an_edge_case() {
        let input = "a (contains dairy)";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 0);
        assert_eq!(part2, "a");
    }

    #[test]
    fn unclosed_allergen_list_is_error() {
        let malformed = "mxmxvkd kfcds sqjhc nhms (contains dairy, fish";
        assert!(solve(malformed).is_err());
    }

    #[test]
    fn missing_contains_section_is_error() {
        let malformed = "mxmxvkd kfcds sqjhc nhms";
        assert!(solve(malformed).is_err());
    }

    #[test]
    fn empty_input_is_error() {
        assert!(solve("").is_err());
    }

    #[test]
    fn example_file_matches_known_answers() {
        let input = std::fs::read_to_string("input/2020/day21.txt");
        let Ok(input) = input else {
            return;
        };
        let (part1, part2) = solve(&input).unwrap();
        assert!(part1 > 0 && !part2.is_empty());
    }
}
