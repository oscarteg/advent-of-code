use anyhow::{bail, ensure, Context, Result};

/// A single navigation instruction: an action letter paired with a value.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum Instruction {
    North(i32),
    South(i32),
    East(i32),
    West(i32),
    Left(i32),
    Right(i32),
    Forward(i32),
}

fn parse_instructions(input: &str) -> Result<Vec<Instruction>> {
    let instructions: Result<Vec<Instruction>> = input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(parse_instruction)
        .collect();

    let instructions = instructions?;
    ensure!(!instructions.is_empty(), "input contains no instructions");
    Ok(instructions)
}

fn parse_instruction(line: &str) -> Result<Instruction> {
    ensure!(line.len() > 1, "instruction line too short: {line:?}");
    let (action, value) = line.split_at(1);
    let value: i32 = value
        .parse()
        .with_context(|| format!("invalid instruction value in line: {line:?}"))?;

    match action {
        "N" => Ok(Instruction::North(value)),
        "S" => Ok(Instruction::South(value)),
        "E" => Ok(Instruction::East(value)),
        "W" => Ok(Instruction::West(value)),
        "L" => {
            ensure!(
                value % 90 == 0,
                "left turn must be a multiple of 90 degrees: {line:?}"
            );
            Ok(Instruction::Left(value))
        }
        "R" => {
            ensure!(
                value % 90 == 0,
                "right turn must be a multiple of 90 degrees: {line:?}"
            );
            Ok(Instruction::Right(value))
        }
        "F" => Ok(Instruction::Forward(value)),
        other => bail!("unknown instruction action {other:?} in line: {line:?}"),
    }
}

/// Rotate a vector `(x, y)` counter-clockwise by `degrees` (a multiple of 90).
fn rotate(x: i32, y: i32, degrees: i32) -> (i32, i32) {
    let steps = (degrees / 90).rem_euclid(4);
    let mut point = (x, y);
    for _ in 0..steps {
        point = (-point.1, point.0);
    }
    point
}

fn part1(instructions: &[Instruction]) -> i32 {
    // Ship starts at the origin facing east.
    let (mut x, mut y) = (0, 0);
    let (mut dx, mut dy) = (1, 0);

    for instruction in instructions {
        match *instruction {
            Instruction::North(v) => y += v,
            Instruction::South(v) => y -= v,
            Instruction::East(v) => x += v,
            Instruction::West(v) => x -= v,
            Instruction::Left(v) => (dx, dy) = rotate(dx, dy, v),
            Instruction::Right(v) => (dx, dy) = rotate(dx, dy, -v),
            Instruction::Forward(v) => {
                x += dx * v;
                y += dy * v;
            }
        }
    }

    x.abs() + y.abs()
}

fn part2(instructions: &[Instruction]) -> i32 {
    // Ship starts at the origin, waypoint starts 10 east and 1 north of it.
    let (mut x, mut y) = (0, 0);
    let (mut wx, mut wy) = (10, 1);

    for instruction in instructions {
        match *instruction {
            Instruction::North(v) => wy += v,
            Instruction::South(v) => wy -= v,
            Instruction::East(v) => wx += v,
            Instruction::West(v) => wx -= v,
            Instruction::Left(v) => (wx, wy) = rotate(wx, wy, v),
            Instruction::Right(v) => (wx, wy) = rotate(wx, wy, -v),
            Instruction::Forward(v) => {
                x += wx * v;
                y += wy * v;
            }
        }
    }

    x.abs() + y.abs()
}

pub fn solve(input: &str) -> Result<(i32, i32)> {
    let instructions = parse_instructions(input)?;
    Ok((part1(&instructions), part2(&instructions)))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/2020/day12-example.txt");

    #[test]
    fn example_part1() {
        let instructions = parse_instructions(EXAMPLE).unwrap();
        assert_eq!(part1(&instructions), 25);
    }

    #[test]
    fn example_part2() {
        let instructions = parse_instructions(EXAMPLE).unwrap();
        assert_eq!(part2(&instructions), 286);
    }

    #[test]
    fn solve_returns_both_parts() {
        let (p1, p2) = solve(EXAMPLE).unwrap();
        assert_eq!(p1, 25);
        assert_eq!(p2, 286);
    }

    #[test]
    fn edge_case_full_rotation_returns_to_start_heading() {
        // Four right turns of 90 degrees bring the ship back to its original heading.
        let instructions = parse_instructions("R90\nR90\nR90\nR90\nF5").unwrap();
        // Facing east again, so moving forward only changes x.
        assert_eq!(part1(&instructions), 5);
    }

    #[test]
    fn malformed_action_is_rejected() {
        assert!(parse_instructions("Q10").is_err());
    }

    #[test]
    fn malformed_value_is_rejected() {
        assert!(parse_instructions("N10x").is_err());
    }

    #[test]
    fn non_right_angle_turn_is_rejected() {
        assert!(parse_instructions("L45").is_err());
    }

    #[test]
    fn empty_input_is_rejected() {
        assert!(parse_instructions("").is_err());
    }
}
