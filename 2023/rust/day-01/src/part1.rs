use crate::custom_error::AocError;

#[tracing::instrument]
pub fn process_line(line: &str) -> u32 {
    let digits = line
        .chars()
        .filter(|c| c.is_ascii_digit())
        .collect::<Vec<_>>();

    let first = digits.first().unwrap();
    let last = digits.last().unwrap();

    let d = format!("{}{}", first, last);

    d.parse::<u32>().unwrap()
}

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let result = input.lines().fold(0, |acc, line| {
        let digit = process_line(line);

        acc + digit
    });

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
        assert_eq!(expected, process_line(line))
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
}
