use crate::day01::custom_error::AocError;

#[tracing::instrument]
fn process_line(line: &str) -> miette::Result<u32, AocError> {
    let digits = line
        .chars()
        .filter(|c| c.is_ascii_digit())
        .collect::<Vec<_>>();

    let first = digits
        .first()
        .ok_or_else(|| AocError::NoDigitsFound(line.to_string()))?;
    let last = digits.last().unwrap();

    let d = format!("{}{}", first, last);
    d.parse::<u32>().map_err(|_| AocError::ParseError(d))
}

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let mut result = 0u32;
    for line in input.lines() {
        let digit = process_line(line)?;
        result += digit;
    }
    Ok(result.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    use rstest::rstest;

    #[rstest]
    #[case("1abc2", 12)]
    #[case("pqr3stu8vwx", 38)]
    #[case("a1b2c3d4e5f", 15)]
    #[case("treb7uchet", 77)]
    fn line_test(#[case] line: &str, #[case] expected: u32) {
        assert_eq!(expected, process_line(line).unwrap())
    }

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = "1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet";
        assert_eq!("142", process(input)?);
        Ok(())
    }

    #[test]
    fn test_single_digit() -> miette::Result<()> {
        let input = "5";
        assert_eq!("55", process(input)?);
        Ok(())
    }

    #[test]
    fn test_many_digits() -> miette::Result<()> {
        let input = "1234567";
        assert_eq!("17", process(input)?);
        Ok(())
    }

    #[test]
    fn test_no_digits() -> miette::Result<()> {
        let input = "abc";
        assert!(process(input).is_err());
        Ok(())
    }

    #[test]
    fn test_empty_line_with_other_lines() -> miette::Result<()> {
        let input = "1abc2\nabc\npqr3stu8vwx";
        assert!(process(input).is_err());
        Ok(())
    }
}
