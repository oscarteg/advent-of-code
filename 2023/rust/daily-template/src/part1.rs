use crate::custom_error::AocError;
use nom::{
    bytes::complete::tag,
    character::complete::{self, digit1, line_ending, space0, space1},
    multi::{fold_many1, separated_list1},
    sequence::{delimited, separated_pair, terminated, tuple},
    IResult, Parser,
};

#[tracing::instrument]
pub fn process(_input: &str) -> miette::Result<String, AocError> {
    todo!("day 01 - part 1");
}

#[tracing::instrument]
pub fn process_line(_input: &str) -> miette::Result<String, AocError> {
    todo!("day 01 - part 1");
}

#[cfg(test)]
mod tests {
    use super::*;

    #[rstest]
    #[case("1abc2", 12)]
    #[case("pqr3stu8vwx", 38)]
    fn line_test(#[case] line: &str, #[case] expected: u32) {
        assert_eq!(expected, process_line(line))
    }

    #[test]
    fn test_process() -> miette::Result<()> {
        todo!("haven't built test yet");
        let input = "";
        assert_eq!("", process(input)?);
        Ok(())
    }
}
