use anyhow::{ensure, Result};

fn parse(input: &str) -> Result<Vec<u64>> {
    input
        .lines()
        .filter(|line| !line.trim().is_empty())
        .map(|line| {
            line.trim()
                .parse::<u64>()
                .map_err(|_| anyhow::anyhow!("invalid number: {line:?}"))
        })
        .collect()
}

fn is_valid(number: u64, preamble: &[u64]) -> bool {
    for i in 0..preamble.len() {
        for j in (i + 1)..preamble.len() {
            if preamble[i] + preamble[j] == number {
                return true;
            }
        }
    }
    false
}

fn find_invalid(numbers: &[u64], preamble_size: usize) -> Result<u64> {
    ensure!(
        numbers.len() > preamble_size,
        "input has {} numbers, preamble size is {}",
        numbers.len(),
        preamble_size
    );

    for i in preamble_size..numbers.len() {
        if !is_valid(numbers[i], &numbers[i - preamble_size..i]) {
            return Ok(numbers[i]);
        }
    }

    anyhow::bail!("no invalid number found")
}

fn find_weakness(numbers: &[u64], target: u64) -> Result<u64> {
    for start in 0..numbers.len() {
        let mut sum = numbers[start];
        for end in (start + 1)..numbers.len() {
            sum += numbers[end];
            if sum == target {
                let range = &numbers[start..=end];
                let min = range.iter().min().copied().unwrap_or(0);
                let max = range.iter().max().copied().unwrap_or(0);
                return Ok(min + max);
            }
            if sum > target {
                break;
            }
        }
    }

    anyhow::bail!("no contiguous sequence found that sums to {target}")
}

pub const PREAMBLE_SIZE: usize = 25;

pub fn solve(input: &str) -> Result<(u64, u64)> {
    solve_with(input, PREAMBLE_SIZE)
}

pub fn solve_with(input: &str, preamble_size: usize) -> Result<(u64, u64)> {
    let numbers = parse(input)?;
    let part1 = find_invalid(&numbers, preamble_size)?;
    let part2 = find_weakness(&numbers, part1)?;
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/2020/day9-example.txt");

    #[test]
    fn sample_part1() {
        let (part1, _) = solve_with(EXAMPLE, 5).unwrap();
        assert_eq!(part1, 127);
    }

    #[test]
    fn sample_part2() {
        let (_, part2) = solve_with(EXAMPLE, 5).unwrap();
        assert_eq!(part2, 62);
    }

    #[test]
    fn edge_case_minimal_sequence() {
        let input = "1\n2\n3\n6";
        let (part1, part2) = solve_with(input, 2).unwrap();
        assert_eq!(part1, 6);
        assert_eq!(part2, 4);
    }

    #[test]
    fn invalid_input_empty() {
        let result = parse("");
        assert!(result.is_ok());
        assert_eq!(result.unwrap().len(), 0);
    }

    #[test]
    fn invalid_input_not_numbers() {
        let result = parse("1\nabc\n3");
        assert!(result.is_err());
    }

    #[test]
    fn invalid_input_preamble_too_large() {
        let input = "1\n2\n3";
        let result = solve_with(input, 5);
        assert!(result.is_err());
    }

    #[test]
    fn invalid_input_no_weakness() {
        let input = "1\n2\n3\n5\n8\n13\n21";
        let result = solve_with(input, 2);
        assert!(result.is_err());
    }
}
