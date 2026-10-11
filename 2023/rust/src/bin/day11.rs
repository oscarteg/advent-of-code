use aoc2023::day11::part1;
use aoc2023::day11::part2;
use miette::{Context, IntoDiagnostic};
use std::io::{self, Read};

fn main() -> miette::Result<()> {
    let mut input = String::new();
    io::stdin().read_to_string(&mut input).into_diagnostic()?;
    println!("Part 1: {}", part1::process(&input).context("part 1")?);
    println!("Part 2: {}", part2::process(&input).context("part 2")?);
    Ok(())
}
