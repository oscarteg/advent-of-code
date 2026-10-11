use anyhow::ensure;

pub fn solve(input: &str) -> anyhow::Result<(u32, u32)> {
    let lines: Vec<&str> = input.lines().filter(|l| !l.trim().is_empty()).collect();
    ensure!(!lines.is_empty(), "input must not be empty");

    let bits = lines[0].len();
    ensure!(bits > 0, "each line must have at least one bit");
    for (i, line) in lines.iter().enumerate() {
        ensure!(
            line.len() == bits,
            "line {}: all lines must have the same length",
            i
        );
        ensure!(
            line.chars().all(|c| c == '0' || c == '1'),
            "line {}: contains non-binary characters",
            i
        );
    }

    let part1 = power_consumption(&lines, bits)?;
    let part2 = life_support_rating(&lines, bits)?;

    Ok((part1, part2))
}

fn power_consumption(lines: &[&str], bits: usize) -> anyhow::Result<u32> {
    let max = lines.len().div_ceil(2);
    let bitmask = (1u32 << bits) - 1;

    let mut gamma = 0u32;

    lines
        .iter()
        .fold(vec![0usize; bits], |count, reading| {
            count
                .into_iter()
                .enumerate()
                .map(|(index, number)| {
                    if reading.as_bytes()[index] == b'1' {
                        number + 1
                    } else {
                        number
                    }
                })
                .collect()
        })
        .into_iter()
        .enumerate()
        .for_each(|(i, bit_count)| {
            if bit_count >= max {
                gamma += 1u32 << (bits - (i + 1));
            }
        });

    let epsilon = !gamma & bitmask;
    Ok(gamma * epsilon)
}

fn life_support_rating(readings: &[&str], bits: usize) -> anyhow::Result<u32> {
    let oxygen_rating = find_rating(readings, bits, true)?;
    let co2_rating = find_rating(readings, bits, false)?;

    let oxygen = u32::from_str_radix(oxygen_rating, 2)?;
    let co2 = u32::from_str_radix(co2_rating, 2)?;

    Ok(oxygen * co2)
}

fn find_rating<'a>(
    readings: &[&'a str],
    bits: usize,
    keep_most_common: bool,
) -> anyhow::Result<&'a str> {
    let mut candidates = readings.to_vec();

    for i in 0..bits {
        if candidates.len() == 1 {
            break;
        }

        let ones_count = candidates
            .iter()
            .filter(|&&s| s.as_bytes()[i] == b'1')
            .count();
        let zeros_count = candidates.len() - ones_count;
        let desired_bit = if (ones_count >= zeros_count) == keep_most_common {
            b'1'
        } else {
            b'0'
        };

        candidates.retain(|&s| s.as_bytes()[i] == desired_bit);
    }

    ensure!(!candidates.is_empty(), "no candidates found");
    Ok(candidates[0])
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010";

    #[test]
    fn test_example() {
        let (part1, part2) = solve(EXAMPLE).expect("valid input");
        assert_eq!(part1, 198, "example part 1");
        assert_eq!(part2, 230, "example part 2");
    }

    #[test]
    fn test_real_input() {
        let input = std::fs::read_to_string("input/day3.txt").expect("failed to read real input");
        let (part1, part2) = solve(&input).expect("valid input");
        assert_eq!(part1, 3148794, "real input part 1");
        assert_eq!(part2, 2795310, "real input part 2");
    }

    #[test]
    fn test_single_line() {
        let (part1, part2) = solve("10101").expect("valid input");
        assert_eq!(part1, 210);
        assert_eq!(part2, 441);
    }

    #[test]
    fn test_alternating_bits() {
        let input = "10\n01";
        let (part1, part2) = solve(input).expect("valid input");
        assert_eq!(part1, 0);
        assert_eq!(part2, 2);
    }

    #[test]
    fn test_malformed_empty() {
        assert!(solve("").is_err(), "empty input should fail");
    }

    #[test]
    fn test_malformed_non_binary() {
        assert!(solve("10102").is_err(), "non-binary character should fail");
    }

    #[test]
    fn test_malformed_uneven_lines() {
        assert!(solve("101\n10").is_err(), "uneven line lengths should fail");
    }
}
