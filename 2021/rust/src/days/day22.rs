use anyhow::{bail, ensure, Context, Result};

const INITIALIZATION_RANGE: Interval = Interval {
    start: -50,
    end: 50,
};

/// Inclusive range of integer coordinates along one axis.
#[derive(Debug, Clone, Copy)]
struct Interval {
    start: i64,
    end: i64,
}

impl Interval {
    fn parse(field: &str, axis: &str) -> Result<Self> {
        let range = field
            .trim()
            .strip_prefix(axis)
            .and_then(|rest| rest.strip_prefix('='))
            .with_context(|| format!("expected `{axis}=start..end`, got `{field}`"))?;
        let (start, end) = range
            .split_once("..")
            .with_context(|| format!("range `{range}` must contain `..`"))?;
        let start: i64 = start
            .parse()
            .with_context(|| format!("invalid range start `{start}`"))?;
        let end: i64 = end
            .parse()
            .with_context(|| format!("invalid range end `{end}`"))?;
        ensure!(start <= end, "range {start}..{end} is empty");
        Ok(Self { start, end })
    }

    fn intersect(self, other: Self) -> Option<Self> {
        let start = self.start.max(other.start);
        let end = self.end.min(other.end);
        (start <= end).then_some(Self { start, end })
    }

    fn length(self) -> i128 {
        i128::from(self.end - self.start + 1)
    }
}

#[derive(Debug, Clone, Copy)]
struct Cuboid {
    x: Interval,
    y: Interval,
    z: Interval,
}

impl Cuboid {
    fn parse(ranges: &str) -> Result<Self> {
        let fields: Vec<&str> = ranges.split(',').collect();
        ensure!(
            fields.len() == 3,
            "expected x, y and z ranges, got `{ranges}`"
        );
        Ok(Self {
            x: Interval::parse(fields[0], "x")?,
            y: Interval::parse(fields[1], "y")?,
            z: Interval::parse(fields[2], "z")?,
        })
    }

    fn intersect(self, other: Self) -> Option<Self> {
        Some(Self {
            x: self.x.intersect(other.x)?,
            y: self.y.intersect(other.y)?,
            z: self.z.intersect(other.z)?,
        })
    }

    fn volume(self) -> i128 {
        self.x.length() * self.y.length() * self.z.length()
    }
}

#[derive(Debug, Clone, Copy)]
struct Step {
    on: bool,
    cuboid: Cuboid,
}

impl Step {
    fn parse(line: &str) -> Result<Self> {
        let (state, ranges) = line
            .trim()
            .split_once(' ')
            .with_context(|| format!("step `{line}` must look like `on x=0..1,y=0..1,z=0..1`"))?;
        let on = match state {
            "on" => true,
            "off" => false,
            other => bail!("unknown reboot state `{other}`"),
        };
        Ok(Self {
            on,
            cuboid: Cuboid::parse(ranges)?,
        })
    }
}

/// Counts the cubes left on after `steps`, clipped to `region` when given.
///
/// Each counted cuboid carries a sign. A new cuboid cancels its overlap with
/// every already counted cuboid by adding that overlap with the opposite sign,
/// so the signed volumes always sum to the number of cubes that are on.
fn count_on(steps: &[Step], region: Option<Cuboid>) -> Result<u64> {
    let mut counted: Vec<(Cuboid, i128)> = Vec::new();
    for step in steps {
        let cuboid = match region {
            Some(region) => match region.intersect(step.cuboid) {
                Some(clipped) => clipped,
                None => continue,
            },
            None => step.cuboid,
        };
        let corrections: Vec<(Cuboid, i128)> = counted
            .iter()
            .filter_map(|&(previous, sign)| Some((previous.intersect(cuboid)?, -sign)))
            .collect();
        counted.extend(corrections);
        if step.on {
            counted.push((cuboid, 1));
        }
    }
    let total: i128 = counted
        .iter()
        .map(|&(cuboid, sign)| sign * cuboid.volume())
        .sum();
    u64::try_from(total).context("cube count must not be negative")
}

/// Returns (part1, part2): cubes on inside the initialization region, and cubes
/// on in the whole reactor.
pub fn solve(input: &str) -> Result<(u64, u64)> {
    let steps: Vec<Step> = input
        .lines()
        .filter(|line| !line.trim().is_empty())
        .map(Step::parse)
        .collect::<Result<Vec<_>>>()?;
    let initialization = Cuboid {
        x: INITIALIZATION_RANGE,
        y: INITIALIZATION_RANGE,
        z: INITIALIZATION_RANGE,
    };
    Ok((
        count_on(&steps, Some(initialization))?,
        count_on(&steps, None)?,
    ))
}

#[cfg(test)]
mod tests {
    use super::*;

    const SMALL_EXAMPLE: &str = "\
on x=10..12,y=10..12,z=10..12
on x=11..13,y=11..13,z=11..13
off x=9..11,y=9..11,z=9..11
on x=10..10,y=10..10,z=10..10
";

    const LARGER_EXAMPLE: &str = "\
on x=-20..26,y=-36..17,z=-47..7
on x=-20..33,y=-21..23,z=-26..28
on x=-22..28,y=-29..23,z=-38..16
on x=-46..7,y=-6..46,z=-50..-1
on x=-49..1,y=-3..46,z=-24..28
on x=2..47,y=-22..22,z=-23..27
on x=-27..23,y=-28..26,z=-21..29
on x=-39..5,y=-6..47,z=-3..44
on x=-30..21,y=-8..43,z=-13..34
on x=-22..26,y=-27..20,z=-29..19
off x=-48..-32,y=26..41,z=-47..-37
on x=-12..35,y=6..50,z=-50..-2
off x=-48..-32,y=-32..-16,z=-15..-5
on x=-18..26,y=-33..15,z=-7..46
off x=-40..-22,y=-38..-28,z=23..41
on x=-16..35,y=-41..10,z=-47..6
off x=-32..-23,y=11..30,z=-14..3
on x=-49..-5,y=-3..45,z=-29..18
off x=18..30,y=-20..-8,z=-3..13
on x=-41..9,y=-7..43,z=-33..15
on x=-54112..-39298,y=-85059..-49293,z=-27449..7877
on x=967..23432,y=45373..81175,z=27513..53682
";

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day22-example.txt");
        assert_eq!(solve(input).unwrap(), (474_140, 2_758_514_936_282_235));
    }

    #[test]
    fn small_example() {
        assert_eq!(solve(SMALL_EXAMPLE).unwrap(), (39, 39));
    }

    #[test]
    fn larger_example_initialization_region() {
        let (part1, _) = solve(LARGER_EXAMPLE).unwrap();
        assert_eq!(part1, 590_784);
    }

    #[test]
    fn cuboid_outside_initialization_region_counts_only_in_part2() {
        assert_eq!(
            solve("on x=100..101,y=100..101,z=100..101").unwrap(),
            (0, 8)
        );
    }

    #[test]
    fn step_is_clipped_at_the_region_border() {
        assert_eq!(solve("on x=49..51,y=0..0,z=0..0").unwrap(), (2, 3));
    }

    #[test]
    fn overlapping_on_steps_count_every_cube_once() {
        let input = "on x=0..1,y=0..0,z=0..0\non x=1..2,y=0..0,z=0..0\n";
        assert_eq!(solve(input).unwrap(), (3, 3));
    }

    #[test]
    fn off_step_clears_earlier_cubes() {
        let input = "on x=0..1,y=0..1,z=0..1\noff x=0..1,y=0..1,z=0..1\n";
        assert_eq!(solve(input).unwrap(), (0, 0));
    }

    #[test]
    fn off_step_without_cubes_on_is_a_no_op() {
        assert_eq!(solve("off x=0..1,y=0..1,z=0..1").unwrap(), (0, 0));
    }

    #[test]
    fn empty_input() {
        assert_eq!(solve("").unwrap(), (0, 0));
    }

    #[test]
    fn malformed_missing_state() {
        assert!(solve("x=0..1,y=0..1,z=0..1").is_err());
    }

    #[test]
    fn malformed_unknown_state() {
        assert!(solve("maybe x=0..1,y=0..1,z=0..1").is_err());
    }

    #[test]
    fn malformed_missing_axis() {
        assert!(solve("on x=0..1,y=0..1").is_err());
    }

    #[test]
    fn malformed_axis_order() {
        assert!(solve("on y=0..1,x=0..1,z=0..1").is_err());
    }

    #[test]
    fn malformed_coordinate() {
        assert!(solve("on x=a..1,y=0..1,z=0..1").is_err());
    }

    #[test]
    fn malformed_reversed_range() {
        assert!(solve("on x=1..0,y=0..1,z=0..1").is_err());
    }
}
