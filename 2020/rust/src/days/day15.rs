use anyhow::{ensure, Context, Result};

/// Parses the starting numbers from a comma-separated list.
fn parse_starting_numbers(input: &str) -> Result<Vec<u32>> {
    let numbers = input
        .trim()
        .split(',')
        .map(|value| {
            value
                .trim()
                .parse::<u32>()
                .context("invalid starting number")
        })
        .collect::<Result<Vec<u32>>>()?;
    ensure!(!numbers.is_empty(), "missing starting numbers");
    Ok(numbers)
}

/// Plays the memory game until `turn` numbers have been spoken and returns the last one.
///
/// A flat array maps each spoken number to the last turn it was said on (0 means "never
/// seen"), which keeps each step O(1) instead of scanning history.
fn play(starting: &[u32], turn: usize) -> Result<u32> {
    ensure!(
        turn >= starting.len(),
        "turn must reach past the starting numbers"
    );
    let max_starting = starting.iter().copied().max().unwrap_or(0) as usize;
    let mut last_seen = vec![0u32; turn.max(max_starting + 1)];
    for (index, &number) in starting.iter().enumerate().take(starting.len() - 1) {
        last_seen[number as usize] = (index + 1) as u32;
    }
    let mut last = *starting.last().context("missing starting numbers")?;
    for current_turn in starting.len()..turn {
        let seen_at = last_seen[last as usize];
        let next = if seen_at == 0 {
            0
        } else {
            current_turn as u32 - seen_at
        };
        last_seen[last as usize] = current_turn as u32;
        last = next;
    }
    Ok(last)
}

pub fn solve(input: &str) -> Result<(u32, u32)> {
    solve_with(input, 2020, 30_000_000)
}

/// Runs the game for both requested turn counts, letting tests pick smaller targets.
pub fn solve_with(input: &str, part1_turn: usize, part2_turn: usize) -> Result<(u32, u32)> {
    let starting = parse_starting_numbers(input)?;
    let part1 = play(&starting, part1_turn)?;
    let part2 = play(&starting, part2_turn)?;
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example_part1() {
        let input = include_str!("../../input/2020/day15-example.txt");
        let (part1, _) = solve_with(input, 2020, 2020).unwrap();
        assert_eq!(part1, 436);
    }

    #[test]
    fn official_example_part2() {
        let input = include_str!("../../input/2020/day15-example.txt");
        let (_, part2) = solve_with(input, 30, 30_000_000).unwrap();
        assert_eq!(part2, 175_594);
    }

    #[test]
    fn more_part1_examples() {
        let cases = [
            ("1,3,2", 1),
            ("2,1,3", 10),
            ("1,2,3", 27),
            ("2,3,1", 78),
            ("3,2,1", 438),
            ("3,1,2", 1836),
        ];
        for (input, expected) in cases {
            let (part1, _) = solve_with(input, 2020, 2020).unwrap();
            assert_eq!(part1, expected, "input {input}");
        }
    }

    #[test]
    fn turn_equal_to_starting_length_returns_last_starting_number() {
        let (part1, _) = solve_with("0,3,6", 3, 3).unwrap();
        assert_eq!(part1, 6);
    }

    #[test]
    fn malformed_input_is_rejected() {
        assert!(solve_with("1,2,a", 10, 10).is_err());
        assert!(solve_with("", 10, 10).is_err());
    }
}
