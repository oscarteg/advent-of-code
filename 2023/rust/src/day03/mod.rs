pub mod custom_error;

pub mod part1;
pub mod part2;

use custom_error::AocError;

pub fn ensure_rectangular(input: &str) -> Result<(), AocError> {
    let input = input.trim_end();
    let width = input.lines().next().map_or(0, str::len);
    match input.lines().position(|line| line.len() != width) {
        Some(index) => Err(AocError::RaggedRow {
            line: index + 1,
            width,
        }),
        None => Ok(()),
    }
}
