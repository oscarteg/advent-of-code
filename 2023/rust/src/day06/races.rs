use crate::day06::custom_error::AocError;

#[derive(Debug, PartialEq, Eq)]
pub struct Race {
    pub time: u64,
    pub distance: u64,
}

impl Race {
    /// Counts the integer hold times that beat the race's record distance.
    pub fn winning_hold_times(&self) -> u64 {
        (0..=self.time)
            .filter(|hold| hold * (self.time - hold) > self.distance)
            .count() as u64
    }
}

fn numbers_after(line: &str, prefix: &str) -> Result<Vec<u64>, AocError> {
    line.strip_prefix(prefix)
        .ok_or(AocError::InvalidInput("missing section prefix"))?
        .split_whitespace()
        .map(|field| field.parse().map_err(AocError::from))
        .collect()
}

fn two_lines(input: &str) -> Result<(&str, &str), AocError> {
    let mut lines = input.lines().map(str::trim).filter(|line| !line.is_empty());
    let times = lines
        .next()
        .ok_or(AocError::InvalidInput("missing time line"))?;
    let distances = lines
        .next()
        .ok_or(AocError::InvalidInput("missing distance line"))?;
    Ok((times, distances))
}

/// Parses one race per column of times and distances.
pub fn parse_races(input: &str) -> Result<Vec<Race>, AocError> {
    let (times, distances) = two_lines(input)?;
    let times = numbers_after(times, "Time:")?;
    let distances = numbers_after(distances, "Distance:")?;
    if times.is_empty() || times.len() != distances.len() {
        return Err(AocError::InvalidInput(
            "time and distance columns must match and be non-empty",
        ));
    }
    Ok(times
        .into_iter()
        .zip(distances)
        .map(|(time, distance)| Race { time, distance })
        .collect())
}

/// Parses a single race by joining the digits on each line into one number.
pub fn parse_single_race(input: &str) -> Result<Race, AocError> {
    let (times, distances) = two_lines(input)?;
    let join = |line: &str, prefix: &str| -> Result<u64, AocError> {
        line.strip_prefix(prefix)
            .ok_or(AocError::InvalidInput("missing section prefix"))?
            .chars()
            .filter(|c| !c.is_whitespace())
            .collect::<String>()
            .parse()
            .map_err(AocError::from)
    };
    Ok(Race {
        time: join(times, "Time:")?,
        distance: join(distances, "Distance:")?,
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn malformed_input_returns_errors() {
        for input in [
            "",
            "Time: 1",
            "Time: x\nDistance: 1",
            "Distance: 1\nTime: 1",
        ] {
            assert!(parse_races(input).is_err(), "{input}");
            assert!(parse_single_race(input).is_err(), "{input}");
        }
        assert!(parse_races("Time: 1 2\nDistance: 1").is_err());
    }
}
