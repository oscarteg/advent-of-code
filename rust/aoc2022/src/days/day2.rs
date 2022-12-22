#![allow(dead_code)]
#![allow(unused_variables)]

use std::{cmp::Ordering, str::FromStr};

#[derive(Debug, PartialEq, Copy, Clone)]
enum Move {
    Rock = 1,
    Paper = 2,
    Scissors = 3,
}

impl PartialOrd for Move {
    fn partial_cmp(&self, other: &Move) -> Option<Ordering> {
        if self == &Move::Scissors && other == &Self::Rock {
            Some(Ordering::Less)
        } else if self == &Move::Rock && other == &Self::Scissors {
            Some(Ordering::Greater)
        } else {
            Some((*self as u8).cmp(&(*other as u8)))
        }
    }
}

// Implement FromStr for Move to be able to parse the input to a Move
impl FromStr for Move {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "A" | "X" => Ok(Move::Rock),
            "B" | "Y" => Ok(Move::Paper),
            "C" | "Z" => Ok(Move::Scissors),
            _ => Err("Not a known move".to_string()),
        }
    }
}

pub fn part_one(input: &str) -> Option<i32> {
    let result = input
        .lines()
        .map(|l| {
            let mut moves = l.split_whitespace().map(|m| m.parse::<Move>().unwrap());
            let a = moves.next().unwrap();
            let b = moves.next().unwrap();

            println!("{:?} vs {:?}", a, b);

            match a.partial_cmp(&b) {
                Some(Ordering::Greater) => b as i32,
                Some(Ordering::Equal) => 3 + b as i32,
                Some(Ordering::Less) => 6 + b as i32,
                _ => panic!("Error for some reason"),
            }
        });


    Some(result.sum())
}

pub fn part_two(input: &str) -> Option<i32> {
    let result: i32 = input
        .lines()
        .map(|l| {
            let mut moves = l.split_whitespace().map(|m| m.parse::<Move>().unwrap());
            let a = moves.next().unwrap();
            let b = moves.next().unwrap();

            match b {
                Move::Rock => {
                    let our_move = match a {
                        Move::Rock => Move::Scissors,
                        Move::Paper => Move::Rock,
                        Move::Scissors => Move::Paper,
                    };
                    our_move as i32
                }
                Move::Paper => 3 + a as i32,
                Move::Scissors => {
                    let our_move = match a {
                        Move::Rock => Move::Paper,
                        Move::Paper => Move::Scissors,
                        Move::Scissors => Move::Rock,
                    };
                    6 + our_move as i32
                }
            }
        })
        .sum();

    Some(result)
}

#[cfg(test)]
mod tests {
    use crate::utils::read_file;

    use super::*;

    const INPUT: &str = "A Y
B X
C Z";

    #[test]
    pub fn test_part_1() {
        // Example
        assert_eq!(part_one(INPUT), Some(15));

        // Input
        let file: String = read_file("input/day2.txt");
        assert_eq!(part_one(file.as_str()), Some(13809));
    }

    #[test]
    fn test_part_2() {
        // Example
        assert_eq!(part_two(INPUT), Some(12));

        // Input
        let file: String = read_file("input/day2.txt");
        assert_eq!(part_two(file.as_str()), Some(12316));
    }
}
