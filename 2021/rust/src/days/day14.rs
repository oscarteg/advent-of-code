use anyhow::{ensure, Context, Result};
use std::collections::HashMap;

struct Polymer {
    template: String,
    rules: HashMap<String, char>,
}

impl Polymer {
    fn parse(input: &str) -> Result<Self> {
        let parts: Vec<&str> = input.trim().split("\n\n").collect();
        ensure!(parts.len() == 2, "expected template and rules sections");

        let template = parts[0].to_string();
        ensure!(!template.is_empty(), "template cannot be empty");

        let mut rules = HashMap::new();
        for line in parts[1].lines() {
            let pair_and_char: Vec<&str> = line.split(" -> ").collect();
            ensure!(pair_and_char.len() == 2, "invalid rule format: {}", line);
            let pair = pair_and_char[0];
            let insertion = pair_and_char[1]
                .chars()
                .next()
                .context("missing insertion char")?;
            ensure!(pair.len() == 2, "pair must be exactly 2 characters");
            rules.insert(pair.to_string(), insertion);
        }
        ensure!(!rules.is_empty(), "no insertion rules found");

        Ok(Polymer { template, rules })
    }

    fn solve(&self, steps: usize) -> Result<u64> {
        let mut pair_counts: HashMap<String, u64> = HashMap::new();
        for i in 0..self.template.len() - 1 {
            let pair = self.template[i..i + 2].to_string();
            *pair_counts.entry(pair).or_insert(0) += 1;
        }

        let first_char = self.template.chars().next().context("empty template")?;

        for _ in 0..steps {
            let mut new_pair_counts: HashMap<String, u64> = HashMap::new();
            for (pair, count) in &pair_counts {
                if let Some(&insertion) = self.rules.get(pair) {
                    let chars: Vec<char> = pair.chars().collect();
                    let left_pair = format!("{}{}", chars[0], insertion);
                    let right_pair = format!("{}{}", insertion, chars[1]);

                    *new_pair_counts.entry(left_pair).or_insert(0) += count;
                    *new_pair_counts.entry(right_pair).or_insert(0) += count;
                } else {
                    *new_pair_counts.entry(pair.clone()).or_insert(0) += count;
                }
            }
            pair_counts = new_pair_counts;
        }

        let mut char_counts: HashMap<char, u64> = HashMap::new();
        char_counts.insert(first_char, 1);

        for (pair, count) in &pair_counts {
            let second_char = pair.chars().nth(1).context("pair too short")?;
            *char_counts.entry(second_char).or_insert(0) += count;
        }

        let max_count = char_counts.values().max().context("no characters found")?;
        let min_count = char_counts.values().min().context("no characters found")?;
        Ok(max_count - min_count)
    }
}

pub fn solve(input: &str) -> Result<(u64, u64)> {
    let polymer = Polymer::parse(input)?;
    let part1 = polymer.solve(10)?;
    let part2 = polymer.solve(40)?;
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day14-example.txt");
        assert_eq!(solve(input).unwrap(), (1588, 2188189693529));
    }

    #[test]
    fn edge_case() {
        let input = "AB\n\nAA -> X\nAB -> Y\nAX -> Z\nAY -> W\nAZ -> Q\nAW -> K\nYA -> M\nYB -> V";
        let result = solve(input);
        assert!(result.is_ok());
    }

    #[test]
    fn missing_template() {
        let input = "\n\nAB -> C";
        assert!(solve(input).is_err());
    }

    #[test]
    fn missing_rules() {
        let input = "ABC";
        assert!(solve(input).is_err());
    }

    #[test]
    fn invalid_rule_format() {
        let input = "ABC\n\nAB : C";
        assert!(solve(input).is_err());
    }
}
