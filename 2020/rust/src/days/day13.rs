use anyhow::{ensure, Context, Result};

/// A bus schedule: the earliest timestamp to depart from, and the bus IDs by
/// offset from that timestamp (`None` marks the "x" placeholders).
struct Schedule {
    earliest: i64,
    buses: Vec<(i64, i64)>,
}

impl Schedule {
    fn parse(input: &str) -> Result<Self> {
        let mut lines = input.lines();
        let earliest: i64 = lines
            .next()
            .context("missing earliest departure timestamp")?
            .trim()
            .parse()
            .context("earliest departure timestamp must be an integer")?;
        let bus_line = lines.next().context("missing bus id list")?.trim();
        ensure!(!bus_line.is_empty(), "bus id list must not be empty");

        let buses = bus_line
            .split(',')
            .enumerate()
            .filter(|(_, id)| *id != "x")
            .map(|(offset, id)| {
                let id: i64 = id
                    .parse()
                    .with_context(|| format!("invalid bus id: {id}"))?;
                ensure!(id > 0, "bus id must be positive");
                Ok((offset as i64, id))
            })
            .collect::<Result<Vec<_>>>()?;
        ensure!(
            !buses.is_empty(),
            "bus id list must contain at least one bus"
        );

        Ok(Self { earliest, buses })
    }
}

/// Part 1: the first bus to depart at or after `earliest`, times the wait.
fn earliest_bus_wait_product(schedule: &Schedule) -> i64 {
    let (id, wait) = schedule
        .buses
        .iter()
        .map(|&(_, id)| {
            let wait = (-schedule.earliest).rem_euclid(id);
            (id, wait)
        })
        .min_by_key(|&(_, wait)| wait)
        .expect("buses is non-empty by construction");
    id * wait
}

/// Part 2: the earliest timestamp `t` such that each bus departs `offset`
/// minutes after `t`, found incrementally via the Chinese remainder theorem.
/// `t` must satisfy `t + offset = 0 (mod id)` for every bus. Buses are folded
/// in one at a time, keeping `t` aligned to the combined period (the lcm of
/// the ids seen so far, which is their product since AoC bus ids are prime).
fn earliest_sequential_timestamp(schedule: &Schedule) -> i64 {
    let mut timestamp = 0i64;
    let mut period = 1i64;
    for &(offset, id) in &schedule.buses {
        while (timestamp + offset) % id != 0 {
            timestamp += period;
        }
        period *= id;
    }
    timestamp
}

pub fn solve(input: &str) -> Result<(i64, i64)> {
    let schedule = Schedule::parse(input)?;
    Ok((
        earliest_bus_wait_product(&schedule),
        earliest_sequential_timestamp(&schedule),
    ))
}

#[cfg(test)]
mod tests {
    use super::*;

    fn part2(bus_line: &str) -> i64 {
        let input = format!("0\n{bus_line}");
        solve(&input).unwrap().1
    }

    #[test]
    fn official_example() {
        let input = include_str!("../../input/2020/day13-example.txt");
        assert_eq!(solve(input).unwrap(), (295, 1068781));
    }

    #[test]
    fn additional_part2_examples() {
        assert_eq!(part2("17,x,13,19"), 3417);
        assert_eq!(part2("67,7,59,61"), 754018);
        assert_eq!(part2("67,x,7,59,61"), 779210);
        assert_eq!(part2("67,7,x,59,61"), 1261476);
        assert_eq!(part2("1789,37,47,1889"), 1202161486);
    }

    #[test]
    fn single_bus_waits_zero_or_the_full_period() {
        assert_eq!(solve("0\n5").unwrap(), (0, 0));
        assert_eq!(solve("3\n5").unwrap(), (10, 0));
    }

    #[test]
    fn missing_bus_line_is_rejected() {
        assert!(Schedule::parse("939").is_err());
    }

    #[test]
    fn all_x_bus_line_is_rejected() {
        assert!(Schedule::parse("939\nx,x,x").is_err());
    }

    #[test]
    fn non_numeric_bus_id_is_rejected() {
        assert!(Schedule::parse("939\n7,abc,59").is_err());
    }
}
