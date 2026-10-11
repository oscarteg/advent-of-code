pub fn solve(input: &str) -> anyhow::Result<(usize, usize)> {
    let depths: Vec<u32> = input
        .lines()
        .filter(|l| !l.is_empty())
        .map(|l| {
            l.trim()
                .parse::<u32>()
                .map_err(|_| anyhow::anyhow!("invalid depth: {}", l))
        })
        .collect::<anyhow::Result<_>>()?;

    anyhow::ensure!(!depths.is_empty(), "no depth measurements found");

    let part1 = depths.windows(2).filter(|w| w[0] < w[1]).count();
    let part2 = depths.windows(4).filter(|w| w[0] < w[3]).count();

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "199
200
208
210
200
207
240
269
260
263";

    #[test]
    fn test_example_part1() {
        let (part1, _) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 7);
    }

    #[test]
    fn test_example_part2() {
        let (_, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part2, 5);
    }

    #[test]
    fn test_real_input() {
        let input = std::fs::read_to_string("input/day1.txt").unwrap();
        let (part1, part2) = solve(&input).unwrap();
        assert_eq!(part1, 1692);
        assert_eq!(part2, 1724);
    }

    #[test]
    fn test_single_depth() {
        let (part1, part2) = solve("100").unwrap();
        assert_eq!(part1, 0);
        assert_eq!(part2, 0);
    }

    #[test]
    fn test_two_depths() {
        let (part1, part2) = solve("100\n200").unwrap();
        assert_eq!(part1, 1);
        assert_eq!(part2, 0);
    }

    #[test]
    fn test_malformed_input() {
        let result = solve("100\nabc\n200");
        assert!(result.is_err());
    }

    #[test]
    fn test_empty_input() {
        let result = solve("");
        assert!(result.is_err());
    }
}
