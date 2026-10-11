use anyhow::{ensure, Context, Result};

fn parse(input: &str) -> Result<Vec<u8>> {
    let timers: Vec<u8> = input
        .trim()
        .split(',')
        .map(|value| value.trim().parse::<u8>().context("invalid fish timer"))
        .collect::<Result<_>>()?;
    ensure!(!timers.is_empty(), "missing fish timers");
    ensure!(
        timers.iter().all(|&timer| timer <= 8),
        "fish timer must be between 0 and 8"
    );
    Ok(timers)
}

// Tracks how many fish share each timer value instead of simulating every
// fish individually, since the population grows exponentially.
fn simulate(timers: &[u8], days: u32) -> u64 {
    let mut counts = [0u64; 9];
    for &timer in timers {
        counts[timer as usize] += 1;
    }
    for _ in 0..days {
        let spawning = counts[0];
        counts.rotate_left(1);
        counts[6] += spawning;
        counts[8] = spawning;
    }
    counts.iter().sum()
}

pub fn solve(input: &str) -> Result<(u64, u64)> {
    let timers = parse(input)?;
    Ok((simulate(&timers, 80), simulate(&timers, 256)))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day6-example.txt");
        assert_eq!(solve(input).unwrap(), (5934, 26984457539));
    }

    #[test]
    fn single_fish_about_to_spawn() {
        assert_eq!(simulate(&parse("0").unwrap(), 1), 2);
        assert_eq!(simulate(&parse("0").unwrap(), 10), 4);
    }

    #[test]
    fn zero_days_returns_input_count() {
        assert_eq!(simulate(&parse("3,4,3,1,2").unwrap(), 0), 5);
    }

    #[test]
    fn non_numeric_timer_is_rejected() {
        assert!(solve("3,4,x,1,2").is_err());
    }

    #[test]
    fn out_of_range_timer_is_rejected() {
        assert!(solve("3,4,9,1,2").is_err());
    }

    #[test]
    fn empty_input_is_rejected() {
        assert!(solve("").is_err());
    }
}
