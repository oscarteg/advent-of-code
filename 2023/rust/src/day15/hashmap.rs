use crate::day15::custom_error::AocError;

/// Splits the comma-separated, newline-stripped initialization sequence into steps.
pub fn steps(input: &str) -> impl Iterator<Item = &str> {
    input.trim().split(',').filter(|step| !step.is_empty())
}

/// The AoC "HASH" algorithm: multiply by 17 and take the remainder mod 256 for every byte.
pub fn hash(step: &str) -> u32 {
    step.bytes()
        .fold(0u32, |acc, byte| (acc + byte as u32) * 17 % 256)
}

/// A single lens, as named in an initialization step.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Lens<'a> {
    pub label: &'a str,
    pub focal_length: u32,
}

/// A parsed initialization step: either remove a lens (`label-`) or insert one (`label=focal`).
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Instruction<'a> {
    Remove { label: &'a str },
    Insert(Lens<'a>),
}

impl<'a> Instruction<'a> {
    pub fn label(&self) -> &'a str {
        match self {
            Instruction::Remove { label } => label,
            Instruction::Insert(lens) => lens.label,
        }
    }
}

pub fn parse_instruction(step: &str) -> Result<Instruction<'_>, AocError> {
    if let Some(label) = step.strip_suffix('-') {
        if label.is_empty() {
            return Err(AocError::InvalidStep(step.to_string()));
        }
        return Ok(Instruction::Remove { label });
    }
    if let Some((label, focal)) = step.split_once('=') {
        if label.is_empty() {
            return Err(AocError::InvalidStep(step.to_string()));
        }
        let focal_length = focal
            .parse()
            .map_err(|_| AocError::InvalidStep(step.to_string()))?;
        return Ok(Instruction::Insert(Lens {
            label,
            focal_length,
        }));
    }
    Err(AocError::InvalidStep(step.to_string()))
}

/// The lens configuration HASHMAP: 256 boxes, each an ordered list of lenses.
#[derive(Debug, Default)]
pub struct HashMap<'a> {
    boxes: Vec<Vec<Lens<'a>>>,
}

impl<'a> HashMap<'a> {
    pub fn new() -> Self {
        Self {
            boxes: (0..256).map(|_| Vec::new()).collect(),
        }
    }

    pub fn apply(&mut self, instruction: Instruction<'a>) {
        let box_index = hash(instruction.label()) as usize;
        let lenses = &mut self.boxes[box_index];
        match instruction {
            Instruction::Remove { label } => {
                lenses.retain(|lens| lens.label != label);
            }
            Instruction::Insert(lens) => {
                if let Some(existing) = lenses.iter_mut().find(|l| l.label == lens.label) {
                    existing.focal_length = lens.focal_length;
                } else {
                    lenses.push(lens);
                }
            }
        }
    }

    /// Sum of focusing powers: (box index + 1) * (slot + 1) * focal length.
    pub fn focusing_power(&self) -> u64 {
        self.boxes
            .iter()
            .enumerate()
            .flat_map(|(box_index, lenses)| {
                lenses.iter().enumerate().map(move |(slot, lens)| {
                    (box_index as u64 + 1) * (slot as u64 + 1) * lens.focal_length as u64
                })
            })
            .sum()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case("HASH", 52)]
    #[case("rn=1", 30)]
    #[case("cm-", 253)]
    #[case("qp=3", 97)]
    fn hash_matches_known_values(#[case] step: &str, #[case] expected: u32) {
        assert_eq!(expected, hash(step));
    }

    #[test]
    fn parse_instruction_rejects_malformed_steps() {
        assert!(parse_instruction("-").is_err());
        assert!(parse_instruction("=5").is_err());
        assert!(parse_instruction("rn=abc").is_err());
        assert!(parse_instruction("rn").is_err());
    }
}
