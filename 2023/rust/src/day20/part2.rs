use crate::day20::{custom_error::AocError, pulses::Network};
use std::collections::HashMap;

const MAX_BUTTON_PRESSES: u64 = 1_000_000;

/// The puzzle's `rx` is never defined with its own line, only referenced as
/// a destination; it is fed by a single conjunction whose inputs each send a
/// high pulse on a fixed cycle. The button count at which `rx` first goes
/// low is the least common multiple of each input's first high-pulse press.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let mut network = Network::parse(input)?;
    let target = network.implicit_sink()?;
    presses_until_target_goes_low(&mut network, &target)
}

fn gcd(a: u64, b: u64) -> u64 {
    if b == 0 {
        a
    } else {
        gcd(b, a % b)
    }
}

fn lcm(a: u64, b: u64) -> u64 {
    a / gcd(a, b) * b
}

fn presses_until_target_goes_low(
    network: &mut Network,
    target: &str,
) -> miette::Result<String, AocError> {
    let feeders = network.inputs_of(target);
    let feeder = match feeders.as_slice() {
        [single] => single,
        [] => return Err(AocError::InvalidInput("target module has no input")),
        _ => {
            return Err(AocError::InvalidInput(
                "target module needs exactly one input",
            ))
        }
    };
    let watched = network.inputs_of(feeder);
    if watched.is_empty() {
        return Err(AocError::InvalidInput("feeder module has no inputs"));
    }

    let mut first_high: HashMap<String, u64> = HashMap::new();
    for press in 1..=MAX_BUTTON_PRESSES {
        if first_high.len() == watched.len() {
            break;
        }
        for (source, destination, high) in network.press_button() {
            if high && destination == *feeder {
                first_high.entry(source).or_insert(press);
            }
        }
    }
    if first_high.len() != watched.len() {
        return Err(AocError::InvalidInput("no high-pulse cycle found"));
    }

    let answer = first_high.values().fold(1u64, |acc, &n| lcm(acc, n));
    Ok(answer.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    /// Synthetic module graph (not a puzzle example): two binary counters of
    /// period 2 and 4 feed single-input conjunctions `inv1`/`inv2`, which
    /// feed `con`, the sole input of the undefined sink `rx`. `inv1` first
    /// sends a high pulse at press 2, `inv2` at press 4, so the answer is
    /// lcm(2, 4) = 4.
    const SYNTHETIC_GRAPH: &str = "broadcaster -> a, b
%a -> inv1
%b -> c
%c -> inv2
&inv1 -> con
&inv2 -> con
&con -> rx";

    #[test]
    fn test_process_synthetic_cycle() {
        assert_eq!(process(SYNTHETIC_GRAPH).unwrap(), "4");
    }

    #[test]
    fn no_undefined_sink_is_rejected() {
        let input = "broadcaster -> a\n%a -> broadcaster";
        assert!(process(input).is_err());
    }

    #[test]
    fn ambiguous_sink_is_rejected() {
        let input = "broadcaster -> a, b\n%a -> rx\n%b -> rx2";
        assert!(process(input).is_err());
    }

    #[test]
    fn target_with_two_inputs_is_rejected() {
        let input = "broadcaster -> a, b\n%a -> rx\n%b -> rx";
        assert!(process(input).is_err());
    }

    #[test]
    fn missing_arrow_is_malformed_input() {
        assert!(process("broadcaster a, b").is_err());
    }
}
