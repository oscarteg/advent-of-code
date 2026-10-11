use crate::day15::custom_error::AocError;
use crate::day15::hashmap::{parse_instruction, steps, HashMap};

/// Runs the initialization sequence through the HASHMAP and sums the resulting focusing power.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let mut map = HashMap::new();
    for step in steps(input) {
        map.apply(parse_instruction(step)?);
    }
    Ok(map.focusing_power().to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day15.txt");
        assert_eq!("145", process(input)?);
        Ok(())
    }

    #[test]
    fn reinserting_a_label_replaces_focal_length() {
        assert_eq!("2", process("rn=1,rn=2").unwrap());
    }

    #[test]
    fn removing_a_label_clears_its_slot() {
        assert_eq!("0", process("rn=1,rn-").unwrap());
    }

    #[test]
    fn empty_input_has_no_focusing_power() {
        assert_eq!("0", process("").unwrap());
    }

    #[test]
    fn malformed_step_is_rejected() {
        assert!(process("rn=abc").is_err());
    }
}
