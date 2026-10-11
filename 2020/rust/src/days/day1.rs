use std::collections::HashSet;

/// Day 1: Report Repair
/// Find two or three entries that sum to 2020 and multiply them.
pub fn solve(input: &str) -> anyhow::Result<(u32, u32)> {
    let mut entries = HashSet::new();

    for line in input.lines() {
        let line = line.trim();
        if line.is_empty() {
            continue;
        }
        let num: u32 = line
            .parse()
            .map_err(|_| anyhow::anyhow!("Invalid number: {}", line))?;
        entries.insert(num);
    }

    anyhow::ensure!(!entries.is_empty(), "No entries found");

    let part1 = solve_part1(&entries)?;
    let part2 = solve_part2(&entries)?;

    Ok((part1, part2))
}

fn solve_part1(entries: &HashSet<u32>) -> anyhow::Result<u32> {
    for &x in entries {
        let y = 2020u32.saturating_sub(x);
        if y > 0 && entries.contains(&y) && x != y {
            return Ok(x * y);
        }
    }
    anyhow::bail!("No solution found for part 1")
}

fn solve_part2(entries: &HashSet<u32>) -> anyhow::Result<u32> {
    let entries_vec: Vec<u32> = entries.iter().copied().collect();

    for (i, &x) in entries_vec.iter().enumerate() {
        for &y in &entries_vec[i..] {
            if x + y >= 2020 {
                continue;
            }
            let z = 2020 - x - y;
            if entries.contains(&z) && z != x && z != y {
                return Ok(x * y * z);
            }
        }
    }
    anyhow::bail!("No solution found for part 2")
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "1721
979
366
299
675
1456";

    #[test]
    fn test_example() {
        let (part1, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 514579);
        assert_eq!(part2, 241861950);
    }

    #[test]
    fn test_small_example_part1() {
        let input = "1010\n1010";
        let result = solve(input);
        assert!(result.is_err());
    }

    #[test]
    fn test_empty_input() {
        let result = solve("");
        assert!(result.is_err());
    }

    #[test]
    fn test_invalid_number() {
        let result = solve("123\nabc\n456");
        assert!(result.is_err());
    }

    #[test]
    fn test_insufficient_entries() {
        let result = solve("1000");
        assert!(result.is_err());
    }

    #[test]
    fn test_no_solution() {
        let result = solve("1\n2\n3\n4\n5");
        assert!(result.is_err());
    }

    #[test]
    fn test_whitespace_handling() {
        let input = "  1721  \n  979  \n  366  \n  299  \n  675  \n  1456  ";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 514579);
        assert_eq!(part2, 241861950);
    }
}
