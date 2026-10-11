use anyhow::{anyhow, Result};

/// Parse a range like "2-4" into (start, end)
fn parse_range(s: &str) -> Result<(u32, u32)> {
    let (start_str, end_str) = s
        .split_once('-')
        .ok_or_else(|| anyhow!("Invalid range format: {}", s))?;
    let start = start_str
        .parse::<u32>()
        .map_err(|_| anyhow!("Invalid start number: {}", start_str))?;
    let end = end_str
        .parse::<u32>()
        .map_err(|_| anyhow!("Invalid end number: {}", end_str))?;
    Ok((start, end))
}

/// Check if range a fully contains range b
fn contains(a: (u32, u32), b: (u32, u32)) -> bool {
    a.0 <= b.0 && a.1 >= b.1
}

/// Check if ranges overlap at all
fn overlaps(a: (u32, u32), b: (u32, u32)) -> bool {
    a.0.max(b.0) <= a.1.min(b.1)
}

pub fn solve(input: &str) -> Result<(u32, u32)> {
    let mut part1 = 0;
    let mut part2 = 0;

    for line in input.trim().lines() {
        let (a_str, b_str) = line
            .split_once(',')
            .ok_or_else(|| anyhow!("Invalid line format: {}", line))?;

        let a = parse_range(a_str)?;
        let b = parse_range(b_str)?;

        if contains(a, b) || contains(b, a) {
            part1 += 1;
        }

        if overlaps(a, b) {
            part2 += 1;
        }
    }

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day4-example.txt");

    #[test]
    fn test_parse_range() {
        assert_eq!(parse_range("2-4").unwrap(), (2, 4));
        assert_eq!(parse_range("6-8").unwrap(), (6, 8));
        assert!(parse_range("invalid").is_err());
        assert!(parse_range("2-a").is_err());
    }

    #[test]
    fn test_contains() {
        assert!(contains((2, 8), (3, 7)));
        assert!(!contains((3, 7), (2, 8)));
        assert!(contains((4, 6), (6, 6)));
        assert!(!contains((6, 6), (4, 6)));
    }

    #[test]
    fn test_overlaps() {
        assert!(overlaps((5, 7), (7, 9)));
        assert!(overlaps((2, 8), (3, 7)));
        assert!(overlaps((6, 6), (4, 6)));
        assert!(!overlaps((2, 4), (6, 8)));
    }

    #[test]
    fn test_example() {
        let (part1, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 2);
        assert_eq!(part2, 4);
    }

    #[test]
    fn test_real_input() {
        let file = std::fs::read_to_string("input/day4.txt").unwrap();
        let (part1, part2) = solve(&file).unwrap();
        assert_eq!(part1, 605);
        assert_eq!(part2, 914);
    }

    #[test]
    fn test_malformed_input() {
        assert!(solve("2-4").is_err());
        assert!(solve("2-4,invalid").is_err());
        assert!(solve("a-4,6-8").is_err());
    }
}
