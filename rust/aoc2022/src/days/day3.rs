#![allow(dead_code)]
#![allow(unused_variables)]

use std::collections::{HashMap, HashSet};

pub fn part_one(input: &str) -> Option<u32> {
    let result = input.lines();

    let mut sum = 0;
    for rucksack in result {
        // Create hashmap for each rucksack
        let mut counts = HashMap::new();

        // Iterate over the characters in the rucksack
        // Count the number of occurrences of each character in the first compartment
        for c in rucksack[..rucksack.len() / 2].chars() {
            *counts.entry(c).or_insert(0) += 1;
        }

        // Find the character that appears in both compartments and add its priority to the sum
        for c in rucksack[rucksack.len() / 2..].chars() {
            if counts.contains_key(&c) {
                sum += get_priority(c);
                break;
            }
        }
    }

    Some(sum)
}

/// Returns the priority of a character
fn get_priority(c: char) -> u32 {
    // Check if the character is a lowercase letter
    if c.is_lowercase() {
        // Convert the character to a number
        (c as u8 - 'a' as u8 + 1) as u32
    } else {
        (c as u8 - 'A' as u8 + 1) as u32 + 26
    }
}

// Returns a hashmap of the priority of all possible characters
fn all_priority() -> HashMap<char, usize> {
    ('a'..='z')
        .chain('A'..='Z')
        .enumerate()
        .map(|(idx, c)| (c, idx + 1))
        .collect::<HashMap<char, usize>>()
}

pub fn part_two(input: &str) -> Option<u32> {
    let result = input.lines();
    let mut sum: u32 = 0;

    // Split the result in multiple groups of 3 lines
    for rucksack in result.collect::<Vec<&str>>().chunks(3) {
        // Create hashmap for each rucksack
        let mut counts = HashSet::new();

        rucksack.iter().for_each(|x| counts.extend(x.chars()));
        rucksack
            .iter()
            .for_each(|x| counts.retain(|y| x.contains(*y)));

        sum += get_priority(*counts.iter().next().unwrap()) as u32;
    }

    Some(sum)
}

#[cfg(test)]
mod tests {
    use crate::utils::read_file;

    use super::*;

    const INPUT: &str = "vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
PmmdzqPrVvPwwTWBwg
wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn
ttgJtRGJQctTZtZT
CrZsJsPPZsGzwwsLwLmpwMDw";

    #[test]
    fn test_part_1() {
        // Example
        assert_eq!(part_one(INPUT), Some(157));

        // Real input
        let file: String = read_file("input/day3.txt");
        assert_eq!(part_one(file.as_str()), Some(7845));
    }

    #[test]
    fn test_part_2() {
        // Example
        assert_eq!(part_two(INPUT), Some(70));

        // Real input
        let file: String = read_file("input/day3.txt");
        assert_eq!(part_two(file.as_str()), Some(2790));
    }
}
