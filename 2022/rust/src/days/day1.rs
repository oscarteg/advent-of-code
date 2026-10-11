use anyhow::{Context, Result};

struct ElfCalories(u32);

fn parse(input: &str) -> Result<Vec<ElfCalories>> {
    let input = input.trim();
    anyhow::ensure!(!input.is_empty(), "input is empty");

    input
        .split("\n\n")
        .map(|elf| {
            elf.lines()
                .map(|line| {
                    line.trim()
                        .parse::<u32>()
                        .with_context(|| format!("invalid calorie value: {line:?}"))
                })
                .try_fold(0u32, |total, calories| calories.map(|value| total + value))
                .map(ElfCalories)
        })
        .collect()
}

pub fn solve(input: &str) -> Result<(u32, u32)> {
    let elves = parse(input)?;
    anyhow::ensure!(!elves.is_empty(), "input has no elves");

    let mut totals: Vec<u32> = elves.into_iter().map(|elf| elf.0).collect();
    totals.sort_unstable_by(|a, b| b.cmp(a));

    let part1 = totals[0];
    let part2 = totals.iter().take(3).sum();
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::solve;

    const EXAMPLE: &str = include_str!("../../input/day1-example.txt");

    #[test]
    fn test_part_1_and_2_example() {
        let (part1, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 24000);
        assert_eq!(part2, 45000);
    }

    #[test]
    fn test_part_1_and_2_real_input() {
        let file = std::fs::read_to_string("input/day1.txt").unwrap();
        let (part1, part2) = solve(&file).unwrap();
        assert_eq!(part1, 70764);
        assert_eq!(part2, 203905);
    }

    #[test]
    fn test_single_elf() {
        let (part1, part2) = solve("100\n200").unwrap();
        assert_eq!(part1, 300);
        assert_eq!(part2, 300);
    }

    #[test]
    fn test_malformed_input_is_err() {
        assert!(solve("1000\nnot-a-number\n\n2000").is_err());
    }

    #[test]
    fn test_empty_input_is_err() {
        assert!(solve("").is_err());
    }
}
