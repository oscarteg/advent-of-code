use anyhow::{ensure, Context, Result};
use std::collections::HashMap;

const START_VALVE: &str = "AA";
const MINUTES_ALONE: u32 = 30;
const MINUTES_WITH_ELEPHANT: u32 = 26;
/// Bounds the bitmask table below; real scans stay far under this.
const MAX_OPENABLE_VALVES: usize = 20;
/// Travel time standing in for "no tunnel route at all".
const UNREACHABLE: u32 = u32::MAX / 4;

/// One line of the scan: a valve, its flow rate and the valves its tunnels reach.
#[derive(Debug)]
struct ScannedValve<'a> {
    name: &'a str,
    flow: u32,
    tunnels: Vec<&'a str>,
}

/// Splits the valve list off its prefix, which is singular for a lone tunnel.
fn strip_tunnel_prefix(text: &str) -> Option<&str> {
    text.strip_prefix("tunnels lead to valves ")
        .or_else(|| text.strip_prefix("tunnel leads to valve "))
}

fn parse_valve(line: &str) -> Result<ScannedValve<'_>> {
    let rest = line
        .strip_prefix("Valve ")
        .with_context(|| format!("malformed valve: {line}"))?;
    let (name, rest) = rest
        .split_once(" has flow rate=")
        .with_context(|| format!("missing flow rate: {line}"))?;
    let (flow, tunnels) = rest
        .split_once("; ")
        .with_context(|| format!("missing tunnel list: {line}"))?;
    let tunnels: Vec<&str> = strip_tunnel_prefix(tunnels)
        .with_context(|| format!("malformed tunnel list: {line}"))?
        .split(", ")
        .collect();
    ensure!(!name.is_empty(), "valve without a name: {line}");
    ensure!(
        tunnels.iter().all(|tunnel| !tunnel.is_empty()),
        "nameless valve in tunnel list: {line}"
    );
    Ok(ScannedValve {
        name,
        flow: flow
            .parse()
            .with_context(|| format!("invalid flow rate: {flow}"))?,
        tunnels,
    })
}

fn parse(input: &str) -> Result<Vec<ScannedValve<'_>>> {
    let valves = input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(parse_valve)
        .collect::<Result<Vec<_>>>()?;
    ensure!(!valves.is_empty(), "input contains no valves");
    Ok(valves)
}

/// Shortest walking time between every pair of valves, by Floyd-Warshall.
fn travel_times(
    valves: &[ScannedValve<'_>],
    index: &HashMap<&str, usize>,
) -> Result<Vec<Vec<u32>>> {
    let count = valves.len();
    let mut minutes = vec![vec![UNREACHABLE; count]; count];
    for (from, valve) in valves.iter().enumerate() {
        for tunnel in &valve.tunnels {
            let to = *index.get(tunnel).with_context(|| {
                format!("tunnel from {} leads to unknown valve {tunnel}", valve.name)
            })?;
            minutes[from][to] = 1;
        }
        minutes[from][from] = 0;
    }
    for via in 0..count {
        let through_via = minutes[via].clone();
        for row in minutes.iter_mut() {
            let hop = row[via];
            if hop >= UNREACHABLE {
                continue;
            }
            for (minutes_to, &via_to) in row.iter_mut().zip(&through_via) {
                let detour = hop + via_to;
                if detour < *minutes_to {
                    *minutes_to = detour;
                }
            }
        }
    }
    Ok(minutes)
}

/// The valves worth opening and the walking times the search needs; every index
/// is a bit position in the sets of opened valves.
struct Network {
    flows: Vec<u32>,
    travel: Vec<Vec<u32>>,
    from_start: Vec<u32>,
}

impl Network {
    fn build(valves: &[ScannedValve<'_>]) -> Result<Self> {
        let mut index = HashMap::with_capacity(valves.len());
        for (position, valve) in valves.iter().enumerate() {
            ensure!(
                index.insert(valve.name, position).is_none(),
                "valve {} is listed twice",
                valve.name
            );
        }
        let start = *index
            .get(START_VALVE)
            .with_context(|| format!("scan has no {START_VALVE} valve to start from"))?;
        let minutes = travel_times(valves, &index)?;

        let openable: Vec<usize> = (0..valves.len())
            .filter(|&position| valves[position].flow > 0)
            .collect();
        ensure!(
            openable.len() <= MAX_OPENABLE_VALVES,
            "{} valves with flow is more than this solver handles ({MAX_OPENABLE_VALVES})",
            openable.len()
        );
        Ok(Self {
            flows: openable.iter().map(|&from| valves[from].flow).collect(),
            travel: openable
                .iter()
                .map(|&from| openable.iter().map(|&to| minutes[from][to]).collect())
                .collect(),
            from_start: openable.iter().map(|&to| minutes[start][to]).collect(),
        })
    }
}

/// Depth-first walk over the valves still closed, recording for every set of
/// opened valves the most pressure that set can release.
fn explore(
    network: &Network,
    best: &mut [u32],
    valve: usize,
    minutes_left: u32,
    opened: usize,
    released: u32,
) {
    best[opened] = best[opened].max(released);
    for (next, &flow) in network.flows.iter().enumerate() {
        if opened & (1 << next) != 0 {
            continue;
        }
        let Some(left) = minutes_left
            .checked_sub(network.travel[valve][next] + 1)
            .filter(|&left| left > 0)
        else {
            continue;
        };
        explore(
            network,
            best,
            next,
            left,
            opened | 1 << next,
            released + left * flow,
        );
    }
}

/// Most pressure released per set of opened valves within the minutes given.
fn best_per_opened_set(network: &Network, minutes: u32) -> Vec<u32> {
    let mut best = vec![0; 1 << network.flows.len()];
    for (valve, &flow) in network.flows.iter().enumerate() {
        if let Some(left) = minutes
            .checked_sub(network.from_start[valve] + 1)
            .filter(|&left| left > 0)
        {
            explore(network, &mut best, valve, left, 1 << valve, left * flow);
        }
    }
    best
}

/// Turns `best[set]` into the best over every subset of `set`, so two disjoint
/// routes can be combined by looking up a set and its complement.
fn maximize_over_subsets(best: &mut [u32], valve_count: usize) {
    for bit in 0..valve_count {
        for set in 0..best.len() {
            if set & (1 << bit) != 0 {
                best[set] = best[set].max(best[set ^ (1 << bit)]);
            }
        }
    }
}

pub fn solve(input: &str) -> Result<(u32, u32)> {
    let network = Network::build(&parse(input)?)?;

    let alone = best_per_opened_set(&network, MINUTES_ALONE)
        .into_iter()
        .max()
        .unwrap_or(0);

    let mut shared = best_per_opened_set(&network, MINUTES_WITH_ELEPHANT);
    maximize_over_subsets(&mut shared, network.flows.len());
    let every_valve = shared.len() - 1;
    let with_elephant = shared
        .iter()
        .enumerate()
        .map(|(mine, released)| released + shared[every_valve ^ mine])
        .max()
        .unwrap_or(0);

    Ok((alone, with_elephant))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day16-example.txt");

    #[test]
    fn official_example() {
        assert_eq!(solve(EXAMPLE).unwrap(), (1651, 1707));
    }

    #[test]
    fn lone_valve_is_opened_as_early_as_possible() {
        // A minute to walk to BB and a minute to open it leaves 28 of the 30
        // minutes at flow 10, and 24 of the 26 when sharing the work.
        let input = "\
Valve AA has flow rate=0; tunnel leads to valve BB
Valve BB has flow rate=10; tunnel leads to valve AA
";
        assert_eq!(solve(input).unwrap(), (280, 240));
    }

    #[test]
    fn valves_without_flow_release_nothing() {
        let input = "\
Valve AA has flow rate=0; tunnels lead to valves BB, CC
Valve BB has flow rate=0; tunnel leads to valve AA
Valve CC has flow rate=0; tunnel leads to valve AA
";
        assert_eq!(solve(input).unwrap(), (0, 0));
    }

    #[test]
    fn valves_without_a_tunnel_route_stay_closed() {
        let input = "\
Valve AA has flow rate=0; tunnel leads to valve BB
Valve BB has flow rate=0; tunnel leads to valve AA
Valve CC has flow rate=5; tunnel leads to valve DD
Valve DD has flow rate=0; tunnel leads to valve CC
";
        assert_eq!(solve(input).unwrap(), (0, 0));
    }

    #[test]
    fn empty_input_is_rejected() {
        assert!(solve("\n\n").is_err());
    }

    #[test]
    fn malformed_scans_are_rejected() {
        assert!(solve("Valve AA has flow rate=0\n").is_err());
        assert!(solve("Valve AA has flow rate=x; tunnel leads to valve BB\n").is_err());
        assert!(solve("Valve AA has flow rate=0; leads to BB\n").is_err());
        assert!(solve("AA has flow rate=0; tunnel leads to valve AA\n").is_err());
    }

    #[test]
    fn inconsistent_scans_are_rejected() {
        assert!(solve("Valve AA has flow rate=0; tunnel leads to valve ZZ\n").is_err());
        assert!(solve("Valve BB has flow rate=1; tunnel leads to valve BB\n").is_err());
        assert!(solve(
            "Valve AA has flow rate=0; tunnel leads to valve AA\n\
             Valve AA has flow rate=1; tunnel leads to valve AA\n"
        )
        .is_err());
    }
}
