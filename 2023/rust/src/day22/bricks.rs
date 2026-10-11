use crate::day22::custom_error::AocError;
use std::collections::{HashMap, VecDeque};

#[derive(Debug, Clone, Copy)]
struct Brick {
    /// Inclusive lower bound per axis.
    min: [u32; 3],
    /// Inclusive upper bound per axis.
    max: [u32; 3],
}

impl Brick {
    fn parse(line: &str) -> Result<Self, AocError> {
        let (start, end) = line.split_once('~').ok_or(AocError::InvalidInput(
            "brick needs two ends separated by '~'",
        ))?;
        let (start, end) = (parse_corner(start)?, parse_corner(end)?);
        let brick = Self {
            min: std::array::from_fn(|axis| start[axis].min(end[axis])),
            max: std::array::from_fn(|axis| start[axis].max(end[axis])),
        };
        if brick.min[2] == 0 {
            return Err(AocError::InvalidInput("bricks rest at z >= 1"));
        }
        if brick
            .min
            .iter()
            .zip(&brick.max)
            .filter(|(a, b)| a != b)
            .count()
            > 1
        {
            return Err(AocError::InvalidInput(
                "brick must be a straight line of cubes",
            ));
        }
        Ok(brick)
    }

    fn footprint(self) -> impl Iterator<Item = (u32, u32)> {
        (self.min[0]..=self.max[0])
            .flat_map(move |x| (self.min[1]..=self.max[1]).map(move |y| (x, y)))
    }

    fn height(self) -> u64 {
        u64::from(self.max[2] - self.min[2]) + 1
    }
}

fn parse_corner(input: &str) -> Result<[u32; 3], AocError> {
    let mut coordinates = input.trim().split(',');
    let mut corner = [0; 3];
    for axis in &mut corner {
        *axis = coordinates
            .next()
            .ok_or(AocError::InvalidInput("corner needs three coordinates"))?
            .parse()?;
    }
    if coordinates.next().is_some() {
        return Err(AocError::InvalidInput("corner needs three coordinates"));
    }
    Ok(corner)
}

pub(crate) struct Stack {
    /// `supports[i]` holds the bricks resting directly on brick `i`.
    supports: Vec<Vec<usize>>,
    /// `supported_by[i]` holds the bricks directly underneath brick `i`.
    supported_by: Vec<Vec<usize>>,
}

impl Stack {
    /// Drops every brick as far as it goes, lowest first, and records the
    /// support relations of the resulting stack.
    pub fn settle(input: &str) -> Result<Self, AocError> {
        let mut bricks = input
            .lines()
            .map(str::trim)
            .filter(|line| !line.is_empty())
            .map(Brick::parse)
            .collect::<Result<Vec<_>, _>>()?;
        if bricks.is_empty() {
            return Err(AocError::InvalidInput("empty snapshot"));
        }
        bricks.sort_unstable_by_key(|brick| brick.min[2]);

        // Per column: the highest occupied z and the brick owning that cube.
        let mut tops: HashMap<(u32, u32), (u64, usize)> = HashMap::new();
        let mut stack = Self {
            supports: vec![Vec::new(); bricks.len()],
            supported_by: vec![Vec::new(); bricks.len()],
        };
        for (index, brick) in bricks.iter().enumerate() {
            let columns: Vec<(u32, u32)> = brick.footprint().collect();
            let below: Vec<(u64, usize)> = columns
                .iter()
                .filter_map(|column| tops.get(column).copied())
                .collect();
            let resting = below.iter().map(|&(top, _)| top).max().unwrap_or(0);
            let mut supporters: Vec<usize> = below
                .iter()
                .filter(|&&(top, _)| top == resting)
                .map(|&(_, owner)| owner)
                .collect();
            supporters.sort_unstable();
            supporters.dedup();
            for &supporter in &supporters {
                stack.supports[supporter].push(index);
            }
            stack.supported_by[index] = supporters;
            let top = resting + brick.height();
            for column in columns {
                tops.insert(column, (top, index));
            }
        }
        Ok(stack)
    }

    pub fn safe_to_disintegrate(&self) -> usize {
        self.supports
            .iter()
            .filter(|resting| {
                resting
                    .iter()
                    .all(|&brick| self.supported_by[brick].len() > 1)
            })
            .count()
    }

    pub fn total_chain_reaction(&self) -> usize {
        (0..self.supports.len())
            .map(|index| self.chain_reaction(index))
            .sum()
    }

    /// Walks up the support graph: a brick falls once all of its supporters do.
    fn chain_reaction(&self, removed: usize) -> usize {
        let mut falling = vec![false; self.supports.len()];
        falling[removed] = true;
        let mut queue = VecDeque::from([removed]);
        let mut fallen = 0;
        while let Some(index) = queue.pop_front() {
            for &above in &self.supports[index] {
                let unsupported = self.supported_by[above]
                    .iter()
                    .all(|&supporter| falling[supporter]);
                if falling[above] || !unsupported {
                    continue;
                }
                falling[above] = true;
                fallen += 1;
                queue.push_back(above);
            }
        }
        fallen
    }
}
