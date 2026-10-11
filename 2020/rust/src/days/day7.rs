use anyhow::{Context, Result};
use lazy_static::lazy_static;
use regex::Regex;
use std::collections::{HashMap, HashSet};

type P1 = usize;
type P2 = usize;

#[derive(Clone, Debug)]
struct Rule {
    bag_type: String,
    contents: Vec<(usize, String)>,
}

fn parse_rules(input: &str) -> Result<Vec<Rule>> {
    lazy_static! {
        static ref BAG_REGEX: Regex = Regex::new(r"^(\w+ \w+) bags contain (.+)$").unwrap();
        static ref CONTENT_REGEX: Regex = Regex::new(r"(\d+) (\w+ \w+) bags?").unwrap();
    }

    let mut rules = Vec::new();
    for line in input.lines() {
        let line = line.trim();
        if line.is_empty() {
            continue;
        }

        let caps = BAG_REGEX
            .captures(line)
            .context("Failed to parse bag rule")?;
        let bag_type = caps.get(1).unwrap().as_str().to_string();
        let contents_str = caps.get(2).unwrap().as_str();

        let mut contents = Vec::new();
        if contents_str != "no other bags." {
            for cap in CONTENT_REGEX.captures_iter(contents_str) {
                let count: usize = cap
                    .get(1)
                    .unwrap()
                    .as_str()
                    .parse()
                    .context("Failed to parse bag count")?;
                let contained_type = cap.get(2).unwrap().as_str().to_string();
                contents.push((count, contained_type));
            }
        }

        rules.push(Rule { bag_type, contents });
    }

    anyhow::ensure!(!rules.is_empty(), "No rules parsed");
    Ok(rules)
}

fn build_reverse_index(rules: &[Rule]) -> HashMap<String, HashSet<String>> {
    let mut can_contain: HashMap<String, HashSet<String>> = HashMap::new();

    for rule in rules {
        for (_count, contained) in &rule.contents {
            can_contain
                .entry(contained.clone())
                .or_default()
                .insert(rule.bag_type.clone());
        }
    }

    can_contain
}

fn count_containers(bag_type: &str, can_contain: &HashMap<String, HashSet<String>>) -> usize {
    let mut visited = HashSet::new();
    let mut queue = vec![bag_type.to_string()];
    visited.insert(bag_type.to_string());

    while let Some(current) = queue.pop() {
        if let Some(containers) = can_contain.get(&current) {
            for container in containers {
                if !visited.contains(container) {
                    visited.insert(container.clone());
                    queue.push(container.clone());
                }
            }
        }
    }

    visited.len() - 1
}

fn part1(rules: &[Rule]) -> P1 {
    let can_contain = build_reverse_index(rules);
    count_containers("shiny gold", &can_contain)
}

fn count_contained_bags(
    bag_type: &str,
    rules_map: &HashMap<String, Vec<(usize, String)>>,
    memo: &mut HashMap<String, usize>,
) -> usize {
    if let Some(&cached) = memo.get(bag_type) {
        return cached;
    }

    memo.insert(bag_type.to_string(), 0);

    let result = if let Some(contents) = rules_map.get(bag_type) {
        contents
            .iter()
            .map(|(count, contained)| {
                count * (1 + count_contained_bags(contained, rules_map, memo))
            })
            .sum()
    } else {
        0
    };

    memo.insert(bag_type.to_string(), result);
    result
}

fn part2(rules: &[Rule]) -> P2 {
    let rules_map: HashMap<String, Vec<(usize, String)>> = rules
        .iter()
        .map(|r| (r.bag_type.clone(), r.contents.clone()))
        .collect();

    let mut memo = HashMap::new();
    count_contained_bags("shiny gold", &rules_map, &mut memo)
}

pub fn solve(input: &str) -> Result<(P1, P2)> {
    let rules = parse_rules(input)?;
    Ok((part1(&rules), part2(&rules)))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE1: &str = include_str!("../../input/2020/day7-example.txt");

    #[test]
    fn test_example_part1() {
        let rules = parse_rules(EXAMPLE1).unwrap();
        assert_eq!(part1(&rules), 4);
    }

    #[test]
    fn test_example_part2() {
        let rules = parse_rules(EXAMPLE1).unwrap();
        assert_eq!(part2(&rules), 32);
    }

    #[test]
    fn test_example2_part2() {
        let example2 = "shiny gold bags contain 2 dark red bags.
dark red bags contain 2 dark orange bags.
dark orange bags contain 2 dark yellow bags.
dark yellow bags contain 2 dark green bags.
dark green bags contain 2 dark blue bags.
dark blue bags contain 2 dark violet bags.
dark violet bags contain no other bags.";
        let rules = parse_rules(example2).unwrap();
        assert_eq!(part2(&rules), 126);
    }

    #[test]
    fn test_parsing_error_on_empty() {
        let result = parse_rules("");
        assert!(result.is_err());
    }

    #[test]
    fn test_single_rule() {
        let input = "shiny gold bags contain 2 dark red bags.";
        let rules = parse_rules(input).unwrap();
        assert_eq!(rules.len(), 1);
        assert_eq!(rules[0].bag_type, "shiny gold");
        assert_eq!(rules[0].contents.len(), 1);
        assert_eq!(rules[0].contents[0], (2, "dark red".to_string()));
    }

    #[test]
    fn test_no_contents() {
        let input = "faded blue bags contain no other bags.";
        let rules = parse_rules(input).unwrap();
        assert_eq!(rules[0].contents.len(), 0);
    }
}
