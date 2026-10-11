use anyhow::{bail, Context, Result};

fn snafu_digit_value(digit: char) -> Result<i64> {
    match digit {
        '2' => Ok(2),
        '1' => Ok(1),
        '0' => Ok(0),
        '-' => Ok(-1),
        '=' => Ok(-2),
        other => bail!("invalid SNAFU digit: {other:?}"),
    }
}

fn parse_snafu(line: &str) -> Result<i64> {
    if line.is_empty() {
        bail!("SNAFU number must not be empty");
    }
    let mut value = 0i64;
    for digit in line.chars() {
        value = value * 5 + snafu_digit_value(digit)?;
    }
    Ok(value)
}

fn parse(input: &str) -> Result<Vec<i64>> {
    input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(|line| parse_snafu(line).with_context(|| format!("malformed SNAFU number: {line:?}")))
        .collect()
}

fn to_snafu(mut value: i64) -> String {
    if value == 0 {
        return "0".to_string();
    }
    let mut digits = Vec::new();
    while value != 0 {
        let remainder = value.rem_euclid(5);
        let (digit, carry) = match remainder {
            0 => ('0', 0),
            1 => ('1', 0),
            2 => ('2', 0),
            3 => ('=', 1),
            4 => ('-', 1),
            _ => unreachable!("remainder is always in 0..5"),
        };
        digits.push(digit);
        value = value / 5 + carry;
    }
    digits.iter().rev().collect()
}

pub fn solve(input: &str) -> Result<String> {
    let numbers = parse(input)?;
    ensure_not_empty(&numbers)?;
    let sum: i64 = numbers.iter().sum();
    Ok(to_snafu(sum))
}

fn ensure_not_empty(numbers: &[i64]) -> Result<()> {
    if numbers.is_empty() {
        bail!("input must contain at least one SNAFU number");
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day25-example.txt");

    #[test]
    fn parses_and_sums_the_example() {
        assert_eq!(solve(EXAMPLE).unwrap(), "2=-1=0");
    }

    #[test]
    fn converts_single_digits() {
        assert_eq!(parse_snafu("1").unwrap(), 1);
        assert_eq!(parse_snafu("2").unwrap(), 2);
        assert_eq!(parse_snafu("1=").unwrap(), 3);
        assert_eq!(parse_snafu("1-").unwrap(), 4);
        assert_eq!(parse_snafu("10").unwrap(), 5);
    }

    #[test]
    fn round_trips_to_snafu() {
        for value in [1i64, 2, 3, 4, 5, 10, 2022, 4890] {
            let snafu = to_snafu(value);
            assert_eq!(parse_snafu(&snafu).unwrap(), value);
        }
    }

    #[test]
    fn handles_a_single_number() {
        assert_eq!(solve("1=-0-2").unwrap(), "1=-0-2");
    }

    #[test]
    fn rejects_empty_input() {
        assert!(solve("").is_err());
        assert!(solve("   \n  \n").is_err());
    }

    #[test]
    fn rejects_invalid_digits() {
        assert!(solve("1=-0-9").is_err());
        assert!(solve("1a2").is_err());
    }
}
