use aoc2023::day02::part1;
use aoc2023::day02::part2;
use miette::{Context, IntoDiagnostic};
use std::io::{self, Read};

fn main() -> miette::Result<()> {
    let mut input = String::new();
    io::stdin().read_to_string(&mut input).into_diagnostic()?;
    // The puzzle plays with 12 red, 13 green and 14 blue cubes.
    println!("Part 1: {}", part1::process(&input, 12, 13, 14).context("part 1")?);
    println!("Part 2: {}", part2::process(&input).context("part 2")?);
    Ok(())
}
