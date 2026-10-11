#[path = "../days/day25.rs"]
mod day25;

use std::io::{self, Read};

fn main() -> anyhow::Result<()> {
    let mut input = String::new();
    io::stdin().read_to_string(&mut input)?;
    let part1 = day25::solve(&input)?;
    println!("Part 1: {part1}");
    Ok(())
}
