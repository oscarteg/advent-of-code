use anyhow::{ensure, Context, Result};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
struct Rect {
    x_min: i32,
    x_max: i32,
    y_min: i32,
    y_max: i32,
}

impl Rect {
    fn parse(input: &str) -> Result<Self> {
        let line = input.trim();
        ensure!(
            line.starts_with("target area: "),
            "expected 'target area:' prefix"
        );
        let coords = &line[13..];
        let parts: Vec<&str> = coords.split(", ").collect();
        ensure!(parts.len() == 2, "expected 'x=..., y=...' format");

        let x_range = parts[0].strip_prefix("x=").context("expected x range")?;
        let y_range = parts[1].strip_prefix("y=").context("expected y range")?;

        let x_bounds: Vec<&str> = x_range.split("..").collect();
        let y_bounds: Vec<&str> = y_range.split("..").collect();

        ensure!(x_bounds.len() == 2, "expected x_min..x_max");
        ensure!(y_bounds.len() == 2, "expected y_min..y_max");

        let x_min: i32 = x_bounds[0].parse().context("x_min not an integer")?;
        let x_max: i32 = x_bounds[1].parse().context("x_max not an integer")?;
        let y_min: i32 = y_bounds[0].parse().context("y_min not an integer")?;
        let y_max: i32 = y_bounds[1].parse().context("y_max not an integer")?;

        Ok(Rect {
            x_min,
            x_max,
            y_min,
            y_max,
        })
    }

    fn contains(&self, x: i32, y: i32) -> bool {
        x >= self.x_min && x <= self.x_max && y >= self.y_min && y <= self.y_max
    }
}

fn simulate(mut vx: i32, mut vy: i32, target: Rect) -> (i32, bool) {
    let mut x = 0;
    let mut y = 0;
    let mut max_y = 0;

    for _ in 0..1000 {
        x += vx;
        y += vy;
        max_y = max_y.max(y);

        if target.contains(x, y) {
            return (max_y, true);
        }

        if x > target.x_max && (y < target.y_min || vy < 0) {
            return (max_y, false);
        }

        vx = (vx - 1).max(0);
        vy -= 1;

        if y < target.y_min - 1000 {
            return (max_y, false);
        }
    }

    (max_y, false)
}

/// Finds highest y-position and count of valid velocities to reach target area.
pub fn solve(input: &str) -> Result<(i32, i32)> {
    let target = Rect::parse(input)?;

    let mut max_height = 0;
    let mut valid_count = 0;

    let vx_min = 1;
    let vx_max = target.x_max + 1;
    let vy_min = target.y_min - 1;
    let vy_max = (-target.y_min).abs() * 2;

    for vx in vx_min..=vx_max {
        for vy in vy_min..=vy_max {
            let (max_y, hit) = simulate(vx, vy, target);
            if hit {
                valid_count += 1;
                max_height = max_height.max(max_y);
            }
        }
    }

    Ok((max_height, valid_count))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parse_target() {
        let target = Rect::parse("target area: x=20..30, y=-10..-5").unwrap();
        assert_eq!(target.x_min, 20);
        assert_eq!(target.x_max, 30);
        assert_eq!(target.y_min, -10);
        assert_eq!(target.y_max, -5);
    }

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day17-example.txt");
        assert_eq!(solve(input).unwrap(), (45, 112));
    }

    #[test]
    fn specific_trajectories() {
        let input = "target area: x=20..30, y=-10..-5";
        let (max_h, count) = solve(input).unwrap();
        assert_eq!(max_h, 45);
        assert_eq!(count, 112);
    }

    #[test]
    fn single_step() {
        let target = Rect::parse("target area: x=10..10, y=-1..-1").unwrap();
        let (_, hit) = simulate(10, -1, target);
        assert!(hit);
    }

    #[test]
    fn malformed_input_missing_prefix() {
        assert!(Rect::parse("x=20..30, y=-10..-5").is_err());
    }

    #[test]
    fn malformed_input_bad_format() {
        assert!(Rect::parse("target area: x=20..30").is_err());
    }

    #[test]
    fn malformed_input_non_integer() {
        assert!(Rect::parse("target area: x=a..30, y=-10..-5").is_err());
    }
}
