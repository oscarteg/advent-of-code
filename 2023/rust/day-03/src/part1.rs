use std::collections::{BTreeMap, HashSet};

use crate::custom_error::AocError;

#[derive(Debug)]
enum Value {
    Symbol(char),
    Number(u32),
    Empty,
}

#[derive(Debug)]
struct Position(i32, i32, bool);

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    // Create a map with all the values
    let map = input
        .lines()
        .enumerate()
        .flat_map(|(y, line)| {
            line.chars().enumerate().map(move |(x, character)| {
                ((x as i32, y as i32), {
                    match character {
                        '.' => Value::Empty,
                        character if character.is_ascii_digit() => {
                            Value::Number(character.to_digit(10).unwrap())
                        }
                        _ => Value::Symbol(character),
                    }
                })
            })
        })
        .collect::<BTreeMap<(i32, i32), Value>>();

    // let positions = [
    //     (-1, -1),
    //     (0, 1),
    //     (1, 1),
    //     (-1, 0),
    //     (1, 0),
    //     (-1, 1),
    //     (1, -1),
    //     (0, -1),
    // ];
    //
    // let mut positions_to_process = Vec::new();
    // // Loop over the input
    // for ((x, y), value) in &map {
    //     match value {
    //         Value::Symbol(_) => {
    //             // If we have a symbol, we need to check if we have 8 adjescant numbers
    //             // If we do, we can replace the symbol with the sum of the numbers
    //             // If we don't, we can leave the symbol as is
    //
    //             for (outer_x, outer_y) in positions.iter() {
    //                 let outer_position = (x + outer_x, y + outer_y);
    //                 if let Some(Value::Number(_)) = map.get(&outer_position) {
    //                     // We have a number, so we can add it to the list
    //                     // Check the adjescant numbers
    //
    //                     let nr = get_full_number(outer_position, &map);
    //                     map.insert(outer_position, Value::Empty);
    //                     positions_to_process.push(nr);
    //                 }
    //             }
    //         }
    //         _ => {}
    //     }
    // }
    //

    let positions = [
        (-1, -1),
        (0, 1),
        (1, 1),
        (-1, 0),
        (1, 0),
        (-1, 1),
        (1, -1),
        (0, -1),
    ];

    let mut sum = 0u32;
    let mut processed_starts = HashSet::new();

    // Loop over the input
    for ((x, y), value) in &map {
        // If we have a symbol, we need to check if the 8 adjescant positions
        if let Value::Symbol(_) = value {
            for (dx, dy) in positions.iter() {
                // Check if the position is a number
                let adj_pos = (x + dx, y + dy);
                if let Some(Value::Number(_)) = map.get(&adj_pos) {
                    // Check if we already processed the first index of this number
                    let start_pos = find_starting_position(adj_pos, &map);
                    if processed_starts.insert(start_pos) {
                        let number = get_full_number(adj_pos, &map);
                        sum += number;
                    }
                }
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

    // Go the the left
    let (mut x, y) = pos;
    while let Some(Value::Number(_)) = map.get(&(x - 1, y)) {
        x -= 1;
    }

    // Check right (including the starting position)
    while let Some(Value::Number(digit)) = map.get(&(x, y)) {
        number = number * 10 + digit;
        x += 1;
    }

    number
}

///
/// -1-1  0-1  1-1
/// -1 0  0-0  1 0
/// -1 1  0 1  1 1

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = "467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..";
        assert_eq!("4361", process(input)?);
        Ok(())
    }
}
