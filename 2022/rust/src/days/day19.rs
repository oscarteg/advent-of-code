use anyhow::{ensure, Context, Result};

const ORE: usize = 0;
const CLAY: usize = 1;
const OBSIDIAN: usize = 2;
const GEODE: usize = 3;
const RESOURCES: usize = 3;

const PART1_MINUTES: u32 = 24;
const PART2_MINUTES: u32 = 32;
const PART2_BLUEPRINTS: usize = 3;

/// Reads a blueprint line literal by literal so malformed input is rejected
/// instead of silently misparsed.
struct Scanner<'a> {
    rest: &'a str,
}

impl<'a> Scanner<'a> {
    fn new(text: &'a str) -> Self {
        Self { rest: text }
    }

    fn literal(&mut self, expected: &str) -> Result<()> {
        self.rest = self
            .rest
            .trim_start()
            .strip_prefix(expected)
            .with_context(|| format!("expected {expected:?}"))?;
        Ok(())
    }

    fn number(&mut self) -> Result<u32> {
        let rest = self.rest.trim_start();
        let digits = rest.len() - rest.trim_start_matches(|c: char| c.is_ascii_digit()).len();
        ensure!(digits > 0, "expected a number");
        self.rest = &rest[digits..];
        rest[..digits].parse().context("number too large")
    }

    fn finish(&self) -> Result<()> {
        let trailing = self.rest.trim();
        ensure!(trailing.is_empty(), "unexpected trailing text {trailing:?}");
        Ok(())
    }
}

/// Robot costs in ore, clay and obsidian, indexed by robot type.
struct Blueprint {
    id: u32,
    costs: [[u32; RESOURCES]; 4],
    /// Highest useful number of ore, clay and obsidian robots: no recipe can
    /// spend more of a resource per minute than its largest cost.
    caps: [u32; RESOURCES],
}

impl Blueprint {
    fn parse(line: &str) -> Result<Self> {
        let mut scanner = Scanner::new(line);
        scanner.literal("Blueprint")?;
        let id = scanner.number()?;
        scanner.literal(":")?;
        scanner.literal("Each ore robot costs")?;
        let ore_robot = scanner.number()?;
        scanner.literal("ore.")?;
        scanner.literal("Each clay robot costs")?;
        let clay_robot = scanner.number()?;
        scanner.literal("ore.")?;
        scanner.literal("Each obsidian robot costs")?;
        let obsidian_robot_ore = scanner.number()?;
        scanner.literal("ore and")?;
        let obsidian_robot_clay = scanner.number()?;
        scanner.literal("clay.")?;
        scanner.literal("Each geode robot costs")?;
        let geode_robot_ore = scanner.number()?;
        scanner.literal("ore and")?;
        let geode_robot_obsidian = scanner.number()?;
        scanner.literal("obsidian.")?;
        scanner.finish()?;

        Ok(Self {
            id,
            costs: [
                [ore_robot, 0, 0],
                [clay_robot, 0, 0],
                [obsidian_robot_ore, obsidian_robot_clay, 0],
                [geode_robot_ore, 0, geode_robot_obsidian],
            ],
            caps: [
                ore_robot
                    .max(clay_robot)
                    .max(obsidian_robot_ore)
                    .max(geode_robot_ore),
                obsidian_robot_clay,
                geode_robot_obsidian,
            ],
        })
    }

    fn max_geodes(&self, minutes: u32) -> u32 {
        let start = State {
            minutes_left: minutes,
            resources: [0; RESOURCES],
            robots: [1, 0, 0],
            geodes: 0,
        };
        let mut best = 0;
        self.search(&start, &mut best);
        best
    }

    /// Depth-first search over the next robot to build, geode robots first,
    /// dropping every branch whose optimistic bound cannot beat the best
    /// schedule found so far.
    fn search(&self, state: &State, best: &mut u32) {
        *best = (*best).max(state.geodes);
        if self.optimistic_geodes(state) <= *best {
            return;
        }
        for robot in [GEODE, OBSIDIAN, CLAY, ORE] {
            if robot != GEODE && state.robots[robot] >= self.caps[robot] {
                continue;
            }
            if let Some(next) = state.build(robot, &self.costs[robot]) {
                self.search(&next, best);
            }
        }
    }

    /// Upper bound on the geodes reachable from `state`, from a relaxed game
    /// that ignores ore costs and lets one robot of every type be built per
    /// minute. Building everything affordable as early as possible is optimal
    /// there (clay is only spent on obsidian robots, obsidian only on geode
    /// robots), so no real schedule can beat the result.
    fn optimistic_geodes(&self, state: &State) -> u32 {
        let obsidian_price = self.costs[GEODE][OBSIDIAN];
        let clay_price = self.costs[OBSIDIAN][CLAY];
        let mut clay = state.resources[CLAY];
        let mut obsidian = state.resources[OBSIDIAN];
        let mut obsidian_robots = state.robots[OBSIDIAN];
        let mut geodes = state.geodes;

        let clay_rates = state.robots[CLAY]..;
        for (clay_robots, remaining) in clay_rates.zip((1..=state.minutes_left).rev()) {
            let mut built_obsidian_robot = 0;
            if obsidian >= obsidian_price {
                obsidian -= obsidian_price;
                geodes += remaining - 1;
            }
            if clay >= clay_price {
                clay -= clay_price;
                built_obsidian_robot = 1;
            }
            clay += clay_robots;
            obsidian += obsidian_robots;
            obsidian_robots += built_obsidian_robot;
        }
        geodes
    }
}

/// A moment in a blueprint's schedule. Geode robots are not tracked: the
/// geodes they will have cracked once the clock runs out are added to
/// `geodes` the minute they are finished.
struct State {
    minutes_left: u32,
    resources: [u32; RESOURCES],
    robots: [u32; RESOURCES],
    geodes: u32,
}

impl State {
    /// Waits until `robot` is affordable and builds it, or returns `None` when
    /// it can never pay off before the clock runs out.
    fn build(&self, robot: usize, cost: &[u32; RESOURCES]) -> Option<Self> {
        let mut wait = 0;
        for ((&price, &amount), &rate) in cost.iter().zip(&self.resources).zip(&self.robots) {
            let missing = price.saturating_sub(amount);
            if missing == 0 {
                continue;
            }
            if rate == 0 {
                return None;
            }
            wait = wait.max(missing.div_ceil(rate));
        }

        let elapsed = wait + 1;
        if elapsed >= self.minutes_left {
            return None;
        }

        let mut next = Self {
            minutes_left: self.minutes_left - elapsed,
            resources: std::array::from_fn(|resource| {
                self.resources[resource] + self.robots[resource] * elapsed - cost[resource]
            }),
            robots: self.robots,
            geodes: self.geodes,
        };
        if robot == GEODE {
            next.geodes += next.minutes_left;
        } else {
            next.robots[robot] += 1;
        }
        Some(next)
    }
}

pub fn solve(input: &str) -> Result<(u32, u32)> {
    let mut blueprints = Vec::new();
    for (number, line) in input.lines().enumerate() {
        let line = line.trim();
        if line.is_empty() {
            continue;
        }
        blueprints.push(Blueprint::parse(line).with_context(|| format!("line {}", number + 1))?);
    }
    ensure!(!blueprints.is_empty(), "input has no blueprints");

    let part1 = blueprints
        .iter()
        .map(|blueprint| blueprint.id * blueprint.max_geodes(PART1_MINUTES))
        .sum();
    let part2 = blueprints
        .iter()
        .take(PART2_BLUEPRINTS)
        .map(|blueprint| blueprint.max_geodes(PART2_MINUTES))
        .product();
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const SECOND_BLUEPRINT: &str = "Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian.";

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day19-example.txt");
        assert_eq!(solve(input).unwrap(), (33, 3472));
    }

    #[test]
    fn quality_level_uses_the_blueprint_id() {
        // The example's second blueprint cracks 12 geodes in 24 minutes and 62
        // in 32 minutes, so on its own it scores 2 * 12 and 62.
        assert_eq!(solve(SECOND_BLUEPRINT).unwrap(), (24, 62));
    }

    #[test]
    fn unaffordable_blueprint_cracks_nothing() {
        let input = "Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 50 clay. Each geode robot costs 90 ore and 90 obsidian.";
        assert_eq!(solve(input).unwrap(), (0, 0));
    }

    #[test]
    fn blank_input_is_rejected() {
        assert!(solve("\n  \n").is_err());
    }

    #[test]
    fn truncated_blueprint_is_rejected() {
        assert!(solve("Blueprint 1: Each ore robot costs 4 ore.").is_err());
    }

    #[test]
    fn non_numeric_id_is_rejected() {
        assert!(solve(&SECOND_BLUEPRINT.replace("Blueprint 2", "Blueprint two")).is_err());
    }

    #[test]
    fn trailing_text_is_rejected() {
        assert!(solve(&format!("{SECOND_BLUEPRINT} and a pony.")).is_err());
    }
}
