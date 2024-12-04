#![allow(dead_code, unused_imports, unused_variables)]
use crate::custom_error::AocError;
use miette::{IntoDiagnostic, Result};
use std::{collections::HashMap, str::FromStr};
use tracing::info;

struct Process {
    destination: i32,
    source: i32,
    length: i32,
}

impl FromStr for Process {
    type Err = AocError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let iter = s.split_whitespace();
        Ok(Self {
            destination: iter.next().unwrap().parse::<i32>(),
            source: iter.next().unwrap().parse::<i32>(),
            length: iter.next().unwrap().parse::<i32>(),
        })
    }
}

struct Partial {
    mapping: HashMap<u32, u32>,
}

fn parse_number(input: &str) -> IResult<&str, u32> {
    map_res(digit1, |c: &str| c.parse::<u32>())(input)
}

fn parse_seeds(input: &str) -> IResult<&str, AocError> {
    println!("input: {} {}", x, y);
    Ok(())
}

// fn parse_row(input: &str) -> IResult<&str, i32> {
//     return input;
// }

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<()> {
    let p = input.split_terminator("\n\n").collect::<Vec<&str>>();
    let mut seeds = parse_seeds(p[0]);

    println!("seeds: {:?}", seeds);
    info!(?seeds);

    // todo!("day 01 - part 1");
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = "seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4";
        // assert_eq!("", process(input)?);
        Ok(())
    }
}
