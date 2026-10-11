use anyhow::{anyhow, Context, Result};

pub fn solve(input: &str) -> Result<(i32, String)> {
    let instructions = parse_input(input)?;
    let mut x_values = vec![1];
    let mut x = 1;

    for instruction in &instructions {
        match instruction {
            Instruction::Noop => {
                x_values.push(x);
            }
            Instruction::Addx(value) => {
                x_values.push(x);
                x_values.push(x);
                x += value;
            }
        }
    }

    let part1 = calculate_signal_strength(&x_values);
    let part2 = render_crt(&x_values);

    Ok((part1, part2))
}

#[derive(Debug, Clone)]
enum Instruction {
    Noop,
    Addx(i32),
}

fn parse_input(input: &str) -> Result<Vec<Instruction>> {
    let mut instructions = Vec::new();
    for line in input.lines() {
        let line = line.trim();
        if line.is_empty() {
            continue;
        }
        if line == "noop" {
            instructions.push(Instruction::Noop);
        } else if let Some(rest) = line.strip_prefix("addx ") {
            let value = rest
                .parse::<i32>()
                .context(format!("Invalid addx value: {}", rest))?;
            instructions.push(Instruction::Addx(value));
        } else {
            return Err(anyhow!("Invalid instruction: {}", line));
        }
    }
    Ok(instructions)
}

fn calculate_signal_strength(x_values: &[i32]) -> i32 {
    let cycles = [20, 60, 100, 140, 180, 220];
    cycles
        .iter()
        .filter_map(|&cycle| {
            if cycle < x_values.len() {
                Some((cycle as i32) * x_values[cycle])
            } else {
                None
            }
        })
        .sum()
}

fn render_crt(x_values: &[i32]) -> String {
    let mut result = String::new();
    for row in 0..6 {
        if row > 0 {
            result.push('\n');
        }
        for col in 0..40 {
            let cycle = row * 40 + col + 1;
            let pixel_x = col as i32;
            if cycle < x_values.len() {
                let sprite_x = x_values[cycle];
                if (sprite_x - 1..=sprite_x + 1).contains(&pixel_x) {
                    result.push('#');
                } else {
                    result.push('.');
                }
            } else {
                result.push('.');
            }
        }
    }
    result
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day10-example.txt");

    #[test]
    fn test_example_part1() -> Result<()> {
        let (part1, _part2) = solve(EXAMPLE)?;
        assert_eq!(part1, 13140);
        Ok(())
    }

    #[test]
    fn test_example_part2() -> Result<()> {
        let (_part1, part2) = solve(EXAMPLE)?;
        let expected = "##..##..##..##..##..##..##..##..##..##..\n\
                        ###...###...###...###...###...###...###.\n\
                        ####....####....####....####....####....\n\
                        #####.....#####.....#####.....#####.....\n\
                        ######......######......######......####\n\
                        #######.......#######.......#######.....";
        assert_eq!(part2, expected);
        Ok(())
    }

    #[test]
    fn test_noop_only() -> Result<()> {
        let (part1, _part2) = solve("noop\nnoop\nnoop")?;
        assert_eq!(part1, 0);
        Ok(())
    }

    #[test]
    fn test_addx() -> Result<()> {
        let input = "addx 3\naddx -5";
        let (part1, _part2) = solve(input)?;
        assert_eq!(part1, 0);
        Ok(())
    }

    #[test]
    fn test_malformed_addx() {
        let result = solve("addx abc");
        assert!(result.is_err());
    }

    #[test]
    fn test_malformed_instruction() {
        let result = solve("invalid");
        assert!(result.is_err());
    }
}
