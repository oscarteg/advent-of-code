use crate::day12::custom_error::AocError;

#[derive(Clone, Copy, Debug, PartialEq, Eq)]
enum Spring {
    Operational,
    Damaged,
    Unknown,
}

/// One row of the condition records: the spring states and the sizes of the
/// contiguous groups of damaged springs.
#[derive(Debug)]
pub(crate) struct Record {
    springs: Vec<Spring>,
    groups: Vec<usize>,
}

impl Record {
    fn parse(line: &str) -> Result<Self, AocError> {
        let mut fields = line.split_whitespace();
        let (Some(springs), Some(groups), None) = (fields.next(), fields.next(), fields.next())
        else {
            return Err(AocError::InvalidInput("row needs springs and group sizes"));
        };
        let springs = springs
            .chars()
            .map(|symbol| match symbol {
                '.' => Ok(Spring::Operational),
                '#' => Ok(Spring::Damaged),
                '?' => Ok(Spring::Unknown),
                _ => Err(AocError::InvalidInput("unknown spring symbol")),
            })
            .collect::<Result<Vec<_>, _>>()?;
        let groups = groups
            .split(',')
            .map(|size| match size.parse::<usize>()? {
                0 => Err(AocError::InvalidInput("group sizes must be positive")),
                size => Ok(size),
            })
            .collect::<Result<Vec<_>, _>>()?;
        Ok(Self { springs, groups })
    }

    /// Repeats the row `copies` times, joining the spring states with `?`.
    pub fn unfold(&self, copies: usize) -> Self {
        let mut springs = Vec::with_capacity((self.springs.len() + 1) * copies);
        for copy in 0..copies {
            if copy > 0 {
                springs.push(Spring::Unknown);
            }
            springs.extend_from_slice(&self.springs);
        }
        Self {
            springs,
            groups: self.groups.repeat(copies),
        }
    }

    /// Counts the ways the unknown springs can be resolved so that the damaged
    /// groups match the listed sizes.
    pub fn arrangements(&self) -> u64 {
        Counter {
            springs: &self.springs,
            groups: &self.groups,
            cache: vec![None; (self.springs.len() + 1) * (self.groups.len() + 1)],
        }
        .count(0, 0)
    }
}

pub(crate) fn parse_records(input: &str) -> Result<Vec<Record>, AocError> {
    let records = input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(Record::parse)
        .collect::<Result<Vec<_>, _>>()?;
    if records.is_empty() {
        return Err(AocError::InvalidInput("no condition records"));
    }
    Ok(records)
}

/// Memoised walk over the row: every state is a spring position paired with the
/// group that still has to be placed, so each pair is counted only once.
struct Counter<'a> {
    springs: &'a [Spring],
    groups: &'a [usize],
    cache: Vec<Option<u64>>,
}

impl Counter<'_> {
    fn count(&mut self, spring: usize, group: usize) -> u64 {
        let key = spring * (self.groups.len() + 1) + group;
        if let Some(cached) = self.cache[key] {
            return cached;
        }
        let arrangements = self.count_uncached(spring, group);
        self.cache[key] = Some(arrangements);
        arrangements
    }

    fn count_uncached(&mut self, spring: usize, group: usize) -> u64 {
        let Some(&state) = self.springs.get(spring) else {
            return u64::from(group == self.groups.len());
        };
        let mut arrangements = 0;
        if state != Spring::Damaged {
            arrangements += self.count(spring + 1, group);
        }
        if state != Spring::Operational {
            arrangements += self.count_group_at(spring, group);
        }
        arrangements
    }

    /// Counts the arrangements that start the remaining group at `spring`. The
    /// group needs that many non-operational springs and an operational one (or
    /// the end of the row) right behind it.
    fn count_group_at(&mut self, spring: usize, group: usize) -> u64 {
        let Some(&size) = self.groups.get(group) else {
            return 0;
        };
        let end = spring + size;
        if end > self.springs.len()
            || self.springs[spring..end].contains(&Spring::Operational)
            || self.springs.get(end) == Some(&Spring::Damaged)
        {
            return 0;
        }
        self.count((end + 1).min(self.springs.len()), group + 1)
    }
}
