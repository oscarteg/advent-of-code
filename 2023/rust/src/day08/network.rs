use crate::day08::custom_error::AocError;
use std::collections::HashMap;

#[derive(Clone, Copy, Debug)]
enum Instruction {
    Left,
    Right,
}

#[derive(Debug)]
pub(crate) struct Network<'a> {
    instructions: Vec<Instruction>,
    nodes: HashMap<&'a str, (&'a str, &'a str)>,
}

impl<'a> Network<'a> {
    pub fn parse(input: &'a str) -> Result<Self, AocError> {
        let mut lines = input.lines().map(str::trim).filter(|line| !line.is_empty());
        let instructions = lines
            .next()
            .ok_or(AocError::InvalidInput("missing instructions"))?
            .chars()
            .map(|step| match step {
                'L' => Ok(Instruction::Left),
                'R' => Ok(Instruction::Right),
                _ => Err(AocError::InvalidInput("instructions are L or R")),
            })
            .collect::<Result<Vec<_>, _>>()?;
        let mut nodes = HashMap::new();
        for line in lines {
            let (name, targets) = line
                .split_once(" = ")
                .ok_or(AocError::InvalidInput("node needs `name = (left, right)`"))?;
            let (left, right) = targets
                .strip_prefix('(')
                .and_then(|inner| inner.strip_suffix(')'))
                .and_then(|pair| pair.split_once(", "))
                .ok_or(AocError::InvalidInput("node needs `name = (left, right)`"))?;
            if nodes.insert(name, (left, right)).is_some() {
                return Err(AocError::InvalidInput("duplicate node"));
            }
        }
        if nodes.is_empty() {
            return Err(AocError::InvalidInput("missing nodes"));
        }
        if nodes
            .values()
            .flat_map(|&(left, right)| [left, right])
            .any(|target| !nodes.contains_key(target))
        {
            return Err(AocError::InvalidInput("node points at an unknown node"));
        }
        Ok(Self {
            instructions,
            nodes,
        })
    }

    /// The nodes the ghosts start from.
    pub fn starts(&self) -> impl Iterator<Item = &'a str> + '_ {
        self.nodes
            .keys()
            .copied()
            .filter(|name| name.ends_with('A'))
    }

    /// Counts the steps from `start` to the first node `at_end` accepts.
    ///
    /// Within `nodes * instructions` steps a walk repeats its (node,
    /// instruction) state, so a longer walk can never arrive.
    pub fn steps(&self, start: &'a str, at_end: impl Fn(&str) -> bool) -> Result<u64, AocError> {
        if !self.nodes.contains_key(start) {
            return Err(AocError::InvalidInput("unknown start node"));
        }
        let mut current = start;
        for steps in 0..self.nodes.len() * self.instructions.len() {
            if at_end(current) {
                return Ok(steps as u64);
            }
            let &(left, right) = self
                .nodes
                .get(current)
                .ok_or(AocError::InvalidInput("unknown node"))?;
            current = match self.instructions[steps % self.instructions.len()] {
                Instruction::Left => left,
                Instruction::Right => right,
            };
        }
        Err(AocError::InvalidInput("no destination is reachable"))
    }
}
