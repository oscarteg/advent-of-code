use anyhow::{anyhow, ensure, Result};
use std::collections::HashSet;

struct Entry {
    signals: Vec<u8>,
    outputs: Vec<u8>,
}

impl Entry {
    fn parse(line: &str) -> Result<Self> {
        let (signals, outputs) = line
            .split_once('|')
            .ok_or_else(|| anyhow!("entry must contain a '|' separator"))?;
        ensure!(
            !outputs.contains('|'),
            "entry must contain exactly one '|' separator"
        );

        let signals = signals
            .split_whitespace()
            .map(parse_pattern)
            .collect::<Result<Vec<_>>>()?;
        let outputs = outputs
            .split_whitespace()
            .map(parse_pattern)
            .collect::<Result<Vec<_>>>()?;

        ensure!(
            signals.len() == 10,
            "an entry must contain exactly ten signal patterns"
        );
        ensure!(
            outputs.len() == 4,
            "an entry must contain exactly four output patterns"
        );

        let unique_signals: HashSet<u8> = signals.iter().copied().collect();
        ensure!(
            unique_signals.len() == signals.len(),
            "signal patterns must be unique"
        );

        Ok(Self { signals, outputs })
    }

    fn decoded_output(&self) -> Result<u64> {
        let mut digits = [0_u8; 10];
        for &signal in &self.signals {
            match signal.count_ones() {
                2 => digits[1] = signal,
                3 => digits[7] = signal,
                4 => digits[4] = signal,
                7 => digits[8] = signal,
                _ => {}
            }
        }
        ensure!(
            digits[1] != 0 && digits[4] != 0 && digits[7] != 0 && digits[8] != 0,
            "entry must include the unique-length digits"
        );

        for &signal in self
            .signals
            .iter()
            .filter(|signal| signal.count_ones() == 6)
        {
            if contains(signal, digits[4]) {
                digits[9] = signal;
            } else if contains(signal, digits[1]) {
                digits[0] = signal;
            } else {
                digits[6] = signal;
            }
        }
        ensure!(
            digits[0] != 0 && digits[6] != 0 && digits[9] != 0,
            "entry must identify every six-segment digit"
        );

        for &signal in self
            .signals
            .iter()
            .filter(|signal| signal.count_ones() == 5)
        {
            if contains(signal, digits[1]) {
                digits[3] = signal;
            } else if contains(digits[6], signal) {
                digits[5] = signal;
            } else {
                digits[2] = signal;
            }
        }
        ensure!(
            digits[2] != 0 && digits[3] != 0 && digits[5] != 0,
            "entry must identify every five-segment digit"
        );

        self.outputs.iter().try_fold(0_u64, |value, &output| {
            Ok(value * 10 + digit_for_pattern(output, &digits)?)
        })
    }
}

fn parse_pattern(pattern: &str) -> Result<u8> {
    ensure!(!pattern.is_empty(), "segment pattern must not be empty");

    pattern.bytes().try_fold(0_u8, |mask, segment| {
        ensure!(
            (b'a'..=b'g').contains(&segment),
            "segment pattern may only contain letters a through g"
        );
        let segment_mask = 1_u8 << (segment - b'a');
        ensure!(
            mask & segment_mask == 0,
            "segment pattern may not repeat a segment"
        );
        Ok(mask | segment_mask)
    })
}

fn contains(pattern: u8, subset: u8) -> bool {
    pattern & subset == subset
}

fn digit_for_pattern(pattern: u8, digits: &[u8; 10]) -> Result<u64> {
    let digit = digits
        .iter()
        .position(|&digit| digit == pattern)
        .ok_or_else(|| anyhow!("output pattern does not match a decoded digit"))?;
    Ok(u64::try_from(digit)?)
}

pub fn solve(input: &str) -> Result<(usize, u64)> {
    let entries = input
        .lines()
        .filter(|line| !line.trim().is_empty())
        .map(Entry::parse)
        .collect::<Result<Vec<_>>>()?;
    ensure!(!entries.is_empty(), "input must contain at least one entry");

    let part_one = entries
        .iter()
        .flat_map(|entry| &entry.outputs)
        .filter(|output| matches!(output.count_ones(), 2 | 3 | 4 | 7))
        .count();
    let part_two = entries
        .iter()
        .map(Entry::decoded_output)
        .sum::<Result<u64>>()?;

    Ok((part_one, part_two))
}

#[cfg(test)]
mod tests {
    use super::solve;

    const EXAMPLE: &str = include_str!("../../input/day8-example.txt");

    #[test]
    fn solve_returns_official_example_answers() {
        let answers = solve(EXAMPLE).unwrap();

        assert_eq!(answers, (26, 61_229));
    }

    #[test]
    fn solve_decodes_all_digits_from_one_entry() {
        let input =
            "cf acf bcdf abcdefg abcefg acdeg acdfg abdfg abdefg abcdfg | abcdefg abcefg acdfg cf";

        let answers = solve(input).unwrap();

        assert_eq!(answers, (2, 8_031));
    }

    #[test]
    fn solve_rejects_malformed_entries() {
        let missing_separator =
            "cf acf bcdf abcdefg abcefg acdeg acdfg abdfg abdefg abcdfg abcdefg abcefg acdfg cf";
        let duplicate_signal =
            "cf cf bcdf abcdefg abcefg acdeg acdfg abdfg abdefg abcdfg | abcdefg abcefg acdfg cf";

        assert!(solve(missing_separator).is_err());
        assert!(solve(duplicate_signal).is_err());
    }
}
