use crate::day20::{custom_error::AocError, pulses::Network};

const BUTTON_PRESSES: usize = 1000;

/// Pushes the button 1000 times and multiplies the total low pulses by the
/// total high pulses sent (including the button's own low pulse each press).
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let mut network = Network::parse(input)?;
    let mut low = 0u64;
    let mut high = 0u64;
    for _ in 0..BUTTON_PRESSES {
        for (_, _, pulse_high) in network.press_button() {
            if pulse_high {
                high += 1;
            } else {
                low += 1;
            }
        }
    }
    Ok((low * high).to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[test]
    fn test_process_first_example() -> miette::Result<()> {
        let input = include_str!("../../examples/day20.txt");
        assert_eq!("32000000", process(input)?);
        Ok(())
    }

    #[test]
    fn test_process_second_example() -> miette::Result<()> {
        let input = include_str!("../../examples/day20-2.txt");
        assert_eq!("11687500", process(input)?);
        Ok(())
    }

    #[rstest]
    #[case("broadcaster -> a\n%a -> b\n&b -> output")]
    fn single_flip_flop_chain_runs_without_error(#[case] input: &str) {
        assert!(process(input).is_ok());
    }

    #[test]
    fn missing_arrow_is_malformed_input() {
        let input = "broadcaster a, b";
        assert!(process(input).is_err());
    }

    #[test]
    fn empty_input_is_malformed() {
        assert!(process("").is_err());
    }
}
