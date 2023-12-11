use crate::custom_error::AocError;

#[tracing::instrument]
fn process_line(line: &str, red: u32, green: u32, blue: u32) -> (u32, bool) {
    let parts: Vec<&str> = line.splitn(2, ": ").collect();
    let game_number_part = parts[0];
    let game_number_str = game_number_part.trim_start_matches("Game ").trim();
    let game_number = game_number_str.parse::<u32>().unwrap();

    let mut subsets = parts[1].split("; ");
    let possible = subsets.all(|subset| {
        subset.split(", ").all(|game| {
            let mut parts = game.split_whitespace();
            let number = parts.next().unwrap().parse::<u32>().unwrap();
            let color = parts.next().unwrap();

            match color {
                "red" => number <= red,
                "green" => number <= green,
                "blue" => number <= blue,
                _ => unreachable!(),
            }
        })
    });

    (game_number, possible)
}

#[tracing::instrument]
pub fn process(input: &str, red: u32, green: u32, blue: u32) -> miette::Result<String, AocError> {
    let result = input.lines().fold(0, |acc, line| {
        let (game_id, possible) = process_line(line, red, green, blue);

        if possible {
            acc + game_id
        } else {
            acc
        }
    });

    Ok(result.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green";
        assert_eq!("8", process(input, 12, 13, 14)?);
        Ok(())
    }
}
