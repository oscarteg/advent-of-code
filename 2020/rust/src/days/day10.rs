use anyhow::{ensure, Context, Result};

/// Parses the adapter joltages, one integer per line.
fn parse(input: &str) -> Result<Vec<u32>> {
    input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(|line| {
            line.parse::<u32>()
                .with_context(|| format!("invalid joltage: {line:?}"))
        })
        .collect()
}

/// Sorts the adapters and appends the device joltage (highest + 3), with the
/// built-in charging outlet (0) implied as the starting point.
fn adapter_chain(adapters: &[u32]) -> Result<Vec<u32>> {
    ensure!(!adapters.is_empty(), "no adapters given");
    let mut chain = adapters.to_vec();
    chain.sort_unstable();
    chain.push(chain.last().unwrap() + 3);
    Ok(chain)
}

/// Counts the 1-jolt and 3-jolt gaps between consecutive adapters in the
/// chain (starting from the outlet at 0), and returns their product.
fn count_joltage_gaps(chain: &[u32]) -> Result<u32> {
    let mut previous = 0;
    let mut ones = 0;
    let mut threes = 0;
    for &joltage in chain {
        match joltage - previous {
            1 => ones += 1,
            2 => {}
            3 => threes += 1,
            gap => anyhow::bail!("unreachable gap of {gap} jolts after {previous}"),
        }
        previous = joltage;
    }
    Ok(ones * threes)
}

/// Counts the distinct ways to arrange the adapters (plus the outlet) into a
/// valid chain reaching the device, each adapter differing by at most 3 jolts
/// from the previous one. Dynamic programming over joltages: the number of
/// ways to reach a joltage is the sum of ways to reach the joltages within 3
/// below it.
fn count_arrangements(chain: &[u32]) -> u64 {
    let device = *chain.last().unwrap();
    let mut ways = vec![0u64; device as usize + 1];
    ways[0] = 1;
    for &joltage in chain {
        let joltage = joltage as usize;
        ways[joltage] = (1..=3)
            .filter_map(|delta| joltage.checked_sub(delta))
            .map(|lower| ways[lower])
            .sum();
    }
    ways[device as usize]
}

pub fn solve(input: &str) -> Result<(u32, u64)> {
    let adapters = parse(input)?;
    let chain = adapter_chain(&adapters)?;
    let part1 = count_joltage_gaps(&chain)?;
    let part2 = count_arrangements(&chain);
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const SMALL_EXAMPLE: &str = include_str!("../../input/2020/day10-example.txt");

    const LARGER_EXAMPLE: &str = "28
33
18
42
31
14
46
20
48
47
24
23
49
45
19
38
39
11
1
32
25
35
8
17
7
9
4
2
34
10
3";

    #[test]
    fn small_example_part1() {
        let (part1, _) = solve(SMALL_EXAMPLE).unwrap();
        assert_eq!(part1, 35);
    }

    #[test]
    fn small_example_part2() {
        let (_, part2) = solve(SMALL_EXAMPLE).unwrap();
        assert_eq!(part2, 8);
    }

    #[test]
    fn larger_example_part1() {
        let (part1, _) = solve(LARGER_EXAMPLE).unwrap();
        assert_eq!(part1, 220);
    }

    #[test]
    fn larger_example_part2() {
        let (_, part2) = solve(LARGER_EXAMPLE).unwrap();
        assert_eq!(part2, 19208);
    }

    #[test]
    fn single_adapter_edge_case() {
        // Outlet (0) -> adapter (1) -> device (4): one 1-jolt gap, one 3-jolt gap.
        let (part1, part2) = solve("1").unwrap();
        assert_eq!(part1, 1);
        assert_eq!(part2, 1);
    }

    #[test]
    fn malformed_input_is_rejected() {
        assert!(solve("not-a-number").is_err());
    }

    #[test]
    fn empty_input_is_rejected() {
        assert!(solve("").is_err());
    }
}
