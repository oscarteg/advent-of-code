use crate::custom_error::AocError;

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let result = input.lines().fold(0, |acc, line| {
        let digits = line
            .chars()
            .filter(|c| c.is_ascii_digit())
            .collect::<Vec<_>>();

        let first = digits.first().unwrap();
        let last = digits.last().unwrap();

        let d = format!("{}{}", first, last);

        let digit = d.parse::<u32>().unwrap();

        acc + digit
    });

    Ok(result.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

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
