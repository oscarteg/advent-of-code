use crate::day05::custom_error::AocError;
use std::ops::Range;

#[derive(Debug)]
pub(crate) struct Mapping {
    pub source: Range<u64>,
    pub destination: u64,
}

#[derive(Debug)]
pub(crate) struct Almanac {
    pub seeds: Vec<u64>,
    pub maps: Vec<Vec<Mapping>>,
}

impl Almanac {
    pub fn parse(input: &str) -> Result<Self, AocError> {
        let mut lines = input.lines().map(str::trim).filter(|line| !line.is_empty());
        let seeds = lines
            .next()
            .and_then(|line| line.strip_prefix("seeds:"))
            .ok_or(AocError::InvalidInput("missing seeds"))?
            .split_whitespace()
            .map(str::parse)
            .collect::<Result<Vec<u64>, _>>()?;
        if seeds.is_empty() {
            return Err(AocError::InvalidInput("empty seed list"));
        }
        let mut maps: Vec<Vec<Mapping>> = Vec::new();
        for line in lines {
            if line.ends_with(" map:") {
                maps.push(Vec::new());
                continue;
            }
            let values = line
                .split_whitespace()
                .map(str::parse)
                .collect::<Result<Vec<u64>, _>>()?;
            let [destination, source, length] = values.as_slice() else {
                return Err(AocError::InvalidInput("mapping requires three numbers"));
            };
            let end = source
                .checked_add(*length)
                .ok_or(AocError::InvalidInput("source range overflows"))?;
            destination
                .checked_add(*length)
                .ok_or(AocError::InvalidInput("destination range overflows"))?;
            maps.last_mut()
                .ok_or(AocError::InvalidInput("mapping has no header"))?
                .push(Mapping {
                    source: *source..end,
                    destination: *destination,
                });
        }
        for map in &mut maps {
            map.sort_by_key(|entry| entry.source.start);
            if map
                .windows(2)
                .any(|pair| pair[0].source.end > pair[1].source.start)
            {
                return Err(AocError::InvalidInput("overlapping source ranges"));
            }
        }
        Ok(Self { seeds, maps })
    }

    pub fn location(&self, seed: u64) -> u64 {
        self.maps.iter().fold(seed, |value, map| {
            match map.iter().find(|entry| entry.source.contains(&value)) {
                Some(entry) => entry.destination + (value - entry.source.start),
                None => value,
            }
        })
    }

    pub fn locations(&self, seeds: Vec<Range<u64>>) -> Vec<Range<u64>> {
        self.maps.iter().fold(seeds, |ranges, map| {
            let mut mapped = Vec::new();
            for range in ranges {
                let mut cursor = range.start;
                for entry in map {
                    if cursor >= range.end {
                        break;
                    }
                    if entry.source.end <= cursor {
                        continue;
                    }
                    if entry.source.start >= range.end {
                        break;
                    }
                    if cursor < entry.source.start {
                        mapped.push(cursor..entry.source.start);
                        cursor = entry.source.start;
                    }
                    let end = range.end.min(entry.source.end);
                    let start = entry.destination + (cursor - entry.source.start);
                    mapped.push(start..start + (end - cursor));
                    cursor = end;
                }
                if cursor < range.end {
                    mapped.push(cursor..range.end);
                }
            }
            mapped
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn interval_mapping_matches_individual_seeds() {
        let almanac =
            Almanac::parse("seeds: 0\n\nfirst map:\n20 3 4\n1 10 3\n\nsecond map:\n50 21 2\n0 7 3")
                .unwrap();
        for start in 0..16 {
            for end in start + 1..=16 {
                let mut expected: Vec<_> =
                    (start..end).map(|seed| almanac.location(seed)).collect();
                let mut actual: Vec<_> = almanac
                    .locations(std::iter::once(start..end).collect())
                    .into_iter()
                    .flatten()
                    .collect();
                expected.sort_unstable();
                actual.sort_unstable();
                assert_eq!(actual, expected, "seed range {start}..{end}");
            }
        }
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in [
            "",
            "seeds:",
            "seeds: x",
            "seeds: 1\n2 3 4",
            "seeds: 1\n\nx map:\n2 3",
            "seeds: 1\n\nx map:\n0 2 3\n9 4 2",
            "seeds: 1\n\nx map:\n0 18446744073709551615 2",
        ] {
            assert!(Almanac::parse(input).is_err(), "{input}");
        }
    }
}
