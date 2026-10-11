use std::collections::{BTreeMap, HashSet};

use crate::day03::custom_error::AocError;

#[derive(Debug)]
enum Value {
    Symbol(char),
    Number(u32),
    Empty,
}

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    crate::day03::ensure_rectangular(input)?;

    let map = input
        .lines()
        .enumerate()
        .flat_map(|(y, line)| {
            line.chars().enumerate().map(move |(x, character)| {
                ((x as i32, y as i32), {
                    match character {
                        '.' => Value::Empty,
                        c if c.is_ascii_digit() => Value::Number(c.to_digit(10).unwrap()),
                        _ => Value::Symbol(character),
                    }
                })
            })
        })
        .collect::<BTreeMap<(i32, i32), Value>>();

    let positions = [
        (-1, -1),
        (0, -1),
        (1, -1),
        (-1, 0),
        (1, 0),
        (-1, 1),
        (0, 1),
        (1, 1),
    ];

    let mut sum = 0u32;

    for ((x, y), value) in &map {
        if let Value::Symbol('*') = value {
            let mut adjacent_numbers = Vec::new();
            let mut processed_starts = HashSet::new();

            for (dx, dy) in positions.iter() {
                let adj_pos = (x + dx, y + dy);
                if let Some(Value::Number(_)) = map.get(&adj_pos) {
                    let start_pos = find_starting_position(adj_pos, &map);
                    if processed_starts.insert(start_pos) {
                        let number = get_full_number(adj_pos, &map);
                        adjacent_numbers.push(number);
                    }
                }
            }

            if adjacent_numbers.len() == 2 {
                sum += adjacent_numbers[0] * adjacent_numbers[1];
            }
        }
    }

    Ok(sum.to_string())
}

fn find_starting_position(mut pos: (i32, i32), map: &BTreeMap<(i32, i32), Value>) -> (i32, i32) {
    while let Some(Value::Number(_)) = map.get(&(pos.0 - 1, pos.1)) {
        pos.0 -= 1;
    }
    pos
}

fn get_full_number(pos: (i32, i32), map: &BTreeMap<(i32, i32), Value>) -> u32 {
    let mut number = 0;
    let (mut x, y) = pos;

    while let Some(Value::Number(_)) = map.get(&(x - 1, y)) {
        x -= 1;
    }

    while let Some(Value::Number(digit)) = map.get(&(x, y)) {
        number = number * 10 + digit;
        x += 1;
    }

    number
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day03.txt");
        assert_eq!("467835", process(input)?);
        Ok(())
    }

    #[test]
    fn test_edge_case_no_gears() -> miette::Result<()> {
        let input = "123.456\n.......";
        assert_eq!("0", process(input)?);
        Ok(())
    }

    #[test]
    fn test_edge_case_gear_adjacent_to_one_number() -> miette::Result<()> {
        let input = "..*.\n123.\n....";
        assert_eq!("0", process(input)?);
        Ok(())
    }

    #[test]
    fn test_malformed_ragged_rows() {
        assert!(process("467..\n...*").is_err());
    }
}
