use crate::day09::custom_error::AocError;

/// A single history of readings reported by one OASIS sensor.
#[derive(Debug, PartialEq, Eq)]
pub(crate) struct History(Vec<i64>);

impl History {
    fn parse(line: &str) -> Result<Self, AocError> {
        let values = line
            .split_whitespace()
            .map(str::parse)
            .collect::<Result<Vec<i64>, _>>()?;
        if values.is_empty() {
            return Err(AocError::InvalidInput("history has no values"));
        }
        Ok(Self(values))
    }

    /// Extrapolates the next value by recursing on successive differences.
    pub fn extrapolate_next(&self) -> i64 {
        Self::next(&self.0)
    }

    /// Extrapolates the previous value the same way, from the front.
    pub fn extrapolate_previous(&self) -> i64 {
        Self::previous(&self.0)
    }

    fn differences(values: &[i64]) -> Vec<i64> {
        values.windows(2).map(|pair| pair[1] - pair[0]).collect()
    }

    fn next(values: &[i64]) -> i64 {
        if values.iter().all(|&value| value == 0) {
            return 0;
        }
        values.last().copied().unwrap_or(0) + Self::next(&Self::differences(values))
    }

    fn previous(values: &[i64]) -> i64 {
        if values.iter().all(|&value| value == 0) {
            return 0;
        }
        values.first().copied().unwrap_or(0) - Self::previous(&Self::differences(values))
    }
}

pub(crate) fn parse(input: &str) -> Result<Vec<History>, AocError> {
    let histories = input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(History::parse)
        .collect::<Result<Vec<_>, _>>()?;
    if histories.is_empty() {
        return Err(AocError::InvalidInput("report has no histories"));
    }
    Ok(histories)
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case(&[0, 3, 6, 9, 12, 15], 18)]
    #[case(&[1, 3, 6, 10, 15, 21], 28)]
    #[case(&[10, 13, 16, 21, 30, 45], 68)]
    #[case(&[5], 5)]
    fn extrapolates_next(#[case] values: &[i64], #[case] expected: i64) {
        assert_eq!(History(values.to_vec()).extrapolate_next(), expected);
    }

    #[rstest]
    #[case(&[0, 3, 6, 9, 12, 15], -3)]
    #[case(&[1, 3, 6, 10, 15, 21], 0)]
    #[case(&[10, 13, 16, 21, 30, 45], 5)]
    fn extrapolates_previous(#[case] values: &[i64], #[case] expected: i64) {
        assert_eq!(History(values.to_vec()).extrapolate_previous(), expected);
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in ["", "   \n  ", "1 2 x"] {
            assert!(parse(input).is_err(), "{input}");
        }
    }

    #[test]
    fn blank_lines_between_histories_are_ignored() {
        assert_eq!(parse("1 2\n\n3 4").unwrap().len(), 2);
    }
}
