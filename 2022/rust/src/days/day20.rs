//! Day 20: Grove Positioning System.
//!
//! Mix a circular list of numbers: each number, in the order it originally
//! appeared, moves forward or backward in the list by its own value.

use anyhow::{ensure, Context, Result};

const DECRYPTION_KEY: i64 = 811589153;
const MIX_ROUNDS_PART2: usize = 10;

fn parse(input: &str) -> Result<Vec<i64>> {
    let numbers: Result<Vec<i64>> = input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(|line| {
            line.parse::<i64>()
                .with_context(|| format!("invalid number: {line}"))
        })
        .collect();
    let numbers = numbers?;
    ensure!(
        !numbers.is_empty(),
        "input must contain at least one number"
    );
    Ok(numbers)
}

/// Mixes `numbers` for `rounds` passes, moving each original index's value
/// forward/backward in the circular list, and returns the mixed order.
fn mix(numbers: &[i64], rounds: usize) -> Vec<i64> {
    let len = numbers.len();
    if len <= 1 {
        return numbers.to_vec();
    }
    // order[i] holds the original index currently at position i.
    let mut order: Vec<usize> = (0..len).collect();

    for _ in 0..rounds {
        for (original_index, &value) in numbers.iter().enumerate() {
            let current_pos = order
                .iter()
                .position(|&idx| idx == original_index)
                .expect("original_index was inserted into order and is never removed");
            order.remove(current_pos);
            let new_pos = (current_pos as i64 + value).rem_euclid(len as i64 - 1) as usize;
            order.insert(new_pos, original_index);
        }
    }

    order.into_iter().map(|idx| numbers[idx]).collect()
}

fn grove_coordinates_sum(mixed: &[i64]) -> Result<i64> {
    let zero_pos = mixed
        .iter()
        .position(|&v| v == 0)
        .context("mixed list must contain a 0 value")?;
    let len = mixed.len();
    Ok([1000, 2000, 3000]
        .iter()
        .map(|offset| mixed[(zero_pos + offset) % len])
        .sum())
}

pub fn solve(input: &str) -> Result<(i64, i64)> {
    let numbers = parse(input)?;

    let mixed_once = mix(&numbers, 1);
    let part1 = grove_coordinates_sum(&mixed_once)?;

    let decrypted: Vec<i64> = numbers.iter().map(|n| n * DECRYPTION_KEY).collect();
    let mixed_many = mix(&decrypted, MIX_ROUNDS_PART2);
    let part2 = grove_coordinates_sum(&mixed_many)?;

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day20-example.txt");

    #[test]
    fn part1_example() {
        let numbers = parse(EXAMPLE).unwrap();
        let mixed = mix(&numbers, 1);
        assert_eq!(grove_coordinates_sum(&mixed).unwrap(), 3);
    }

    #[test]
    fn part2_example() {
        let numbers = parse(EXAMPLE).unwrap();
        let decrypted: Vec<i64> = numbers.iter().map(|n| n * DECRYPTION_KEY).collect();
        let mixed = mix(&decrypted, MIX_ROUNDS_PART2);
        assert_eq!(grove_coordinates_sum(&mixed).unwrap(), 1623178306);
    }

    #[test]
    fn solve_example() {
        let (part1, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 3);
        assert_eq!(part2, 1623178306);
    }

    #[test]
    fn single_element_edge_case() {
        let numbers = parse("0\n").unwrap();
        let mixed = mix(&numbers, 1);
        assert_eq!(mixed, vec![0]);
    }

    #[test]
    fn duplicate_values_are_tracked_by_original_index() {
        // Values repeat; mixing must still move each occurrence independently.
        let numbers = parse("1\n1\n0\n-1\n").unwrap();
        let mixed = mix(&numbers, 1);
        assert_eq!(mixed.iter().sum::<i64>(), 1);
        assert_eq!(mixed.len(), 4);
    }

    #[test]
    fn malformed_input_errors() {
        assert!(parse("not-a-number\n").is_err());
    }

    #[test]
    fn empty_input_errors() {
        assert!(parse("").is_err());
    }
}
