use anyhow::{bail, ensure, Context, Result};
use std::collections::HashSet;

const REAL_ROW: i64 = 2_000_000;
const REAL_SEARCH_LIMIT: i64 = 4_000_000;
const TUNING_MULTIPLIER: i64 = 4_000_000;

/// A sensor, its closest beacon and the Manhattan distance between the two.
#[derive(Debug, Clone, Copy)]
struct Reading {
    sensor: (i64, i64),
    beacon: (i64, i64),
    radius: i64,
}

impl Reading {
    /// The inclusive x-range this sensor rules out on `row`, if it reaches it.
    fn row_span(&self, row: i64) -> Option<(i64, i64)> {
        let spread = self.radius - (self.sensor.1 - row).abs();
        (spread >= 0).then_some((self.sensor.0 - spread, self.sensor.0 + spread))
    }
}

fn parse_point(text: &str) -> Result<(i64, i64)> {
    let (x, y) = text
        .split_once(", ")
        .with_context(|| format!("malformed coordinate pair: {text}"))?;
    let x = x
        .strip_prefix("x=")
        .with_context(|| format!("missing x= in: {x}"))?;
    let y = y
        .strip_prefix("y=")
        .with_context(|| format!("missing y= in: {y}"))?;
    Ok((
        x.parse().with_context(|| format!("invalid x: {x}"))?,
        y.parse().with_context(|| format!("invalid y: {y}"))?,
    ))
}

fn parse(input: &str) -> Result<Vec<Reading>> {
    let readings = input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(|line| {
            let (sensor, beacon) = line
                .split_once(": closest beacon is at ")
                .with_context(|| format!("malformed reading: {line}"))?;
            let sensor = parse_point(
                sensor
                    .strip_prefix("Sensor at ")
                    .with_context(|| format!("malformed sensor: {sensor}"))?,
            )?;
            let beacon = parse_point(beacon)?;
            Ok(Reading {
                sensor,
                beacon,
                radius: (sensor.0 - beacon.0).abs() + (sensor.1 - beacon.1).abs(),
            })
        })
        .collect::<Result<Vec<_>>>()?;
    ensure!(!readings.is_empty(), "input contains no sensor readings");
    Ok(readings)
}

/// Disjoint inclusive x-ranges covering `row`, sorted by start and with touching
/// ranges joined, so that any boundary left between them is a real gap.
fn merged_row_spans(readings: &[Reading], row: i64) -> Vec<(i64, i64)> {
    let mut spans: Vec<(i64, i64)> = readings
        .iter()
        .filter_map(|reading| reading.row_span(row))
        .collect();
    spans.sort_unstable();

    let mut merged: Vec<(i64, i64)> = Vec::with_capacity(spans.len());
    for (start, end) in spans {
        match merged.last_mut() {
            Some(last) if start <= last.1 + 1 => last.1 = last.1.max(end),
            _ => merged.push((start, end)),
        }
    }
    merged
}

/// Positions on `row` that cannot hold a beacon: every covered position minus the
/// distinct known beacons there, which each sit inside their own sensor's range.
fn excluded_positions(readings: &[Reading], row: i64) -> Result<i64> {
    let covered: i64 = merged_row_spans(readings, row)
        .iter()
        .map(|(start, end)| end - start + 1)
        .sum();
    let beacons: HashSet<i64> = readings
        .iter()
        .filter(|reading| reading.beacon.1 == row)
        .map(|reading| reading.beacon.0)
        .collect();
    Ok(covered - i64::try_from(beacons.len())?)
}

/// The smallest x in `0..=limit` on `row` that no sensor covers.
fn uncovered_x(readings: &[Reading], row: i64, limit: i64) -> Option<i64> {
    let mut x = 0;
    for (start, end) in merged_row_spans(readings, row) {
        if start > x {
            return Some(x);
        }
        x = x.max(end + 1);
        if x > limit {
            return None;
        }
    }
    (x <= limit).then_some(x)
}

/// Tuning frequency of the one position in the search square no sensor covers.
fn tuning_frequency(readings: &[Reading], limit: i64) -> Result<i64> {
    ensure!(limit >= 0, "search limit must not be negative");
    for row in 0..=limit {
        if let Some(x) = uncovered_x(readings, row, limit) {
            return Ok(x * TUNING_MULTIPLIER + row);
        }
    }
    bail!("no position in the search area is free of sensor coverage")
}

pub fn solve_with(input: &str, row: i64, search_limit: i64) -> Result<(i64, i64)> {
    let readings = parse(input)?;
    Ok((
        excluded_positions(&readings, row)?,
        tuning_frequency(&readings, search_limit)?,
    ))
}

pub fn solve(input: &str) -> Result<(i64, i64)> {
    solve_with(input, REAL_ROW, REAL_SEARCH_LIMIT)
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day15-example.txt");
    const EXAMPLE_ROW: i64 = 10;
    const EXAMPLE_SEARCH_LIMIT: i64 = 20;

    #[test]
    fn official_example() {
        assert_eq!(
            solve_with(EXAMPLE, EXAMPLE_ROW, EXAMPLE_SEARCH_LIMIT).unwrap(),
            (26, 56_000_011)
        );
    }

    #[test]
    fn solve_uses_the_real_puzzle_row() {
        let (part1, _) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 0, "no example sensor reaches row 2000000");
    }

    #[test]
    fn single_sensor_rules_out_its_diamond_row() {
        // Row 0 covers x in -2..=2; the known beacon at x=2 stays a valid position,
        // and row 1 leaves x=2 uncovered.
        let input = "Sensor at x=0, y=0: closest beacon is at x=2, y=0\n";
        assert_eq!(solve_with(input, 0, 2).unwrap(), (4, 8_000_001));
    }

    #[test]
    fn row_outside_every_diamond_rules_out_nothing() {
        let (part1, _) = solve_with(EXAMPLE, 1_000, EXAMPLE_SEARCH_LIMIT).unwrap();
        assert_eq!(part1, 0);
    }

    #[test]
    fn fully_covered_search_area_has_no_distress_beacon() {
        let input = "Sensor at x=0, y=0: closest beacon is at x=10, y=0\n";
        assert!(solve_with(input, 0, 2).is_err());
    }

    #[test]
    fn empty_input_is_rejected() {
        assert!(solve("\n\n").is_err());
    }

    #[test]
    fn malformed_readings_are_rejected() {
        assert!(solve_with("Sensor at x=2, y=18\n", 10, 20).is_err());
        assert!(solve_with(
            "Sensor at x=a, y=18: closest beacon is at x=-2, y=15\n",
            10,
            20
        )
        .is_err());
        assert!(solve_with("at x=2, y=18: closest beacon is at x=-2, y=15\n", 10, 20).is_err());
    }
}
