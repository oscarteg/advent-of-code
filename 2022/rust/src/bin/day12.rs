#[path = "../days/day12.rs"]
mod day12;

use std::io::{self, Read};

fn main() -> anyhow::Result<()> {
    let mut input = String::new();
    io::stdin().read_to_string(&mut input)?;
    let (part1, part2) = day12::solve(&input)?;
    println!("Part 1: {part1}\nPart 2: {part2}");
    Ok(())
}
