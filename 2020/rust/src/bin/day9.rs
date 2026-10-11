#[path = "../days/day9.rs"]
mod day9;

use anyhow::Context;
use std::io::{self, Read};

fn main() -> anyhow::Result<()> {
    let mut input = String::new();
    io::stdin().read_to_string(&mut input)?;
    let (part1, part2) = match std::env::args().nth(1) {
        Some(size) => day9::solve_with(
            &input,
            size.parse().context("preamble size must be a number")?,
        )?,
        None => day9::solve(&input)?,
    };
    println!("Part 1: {part1}\nPart 2: {part2}");
    Ok(())
}
