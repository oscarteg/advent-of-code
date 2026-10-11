use crate::day02::custom_error::AocError;

/// Number of cubes of each color shown in a single handful.
#[derive(Debug, Default, Clone, Copy, PartialEq, Eq)]
pub struct Reveal {
    pub red: u32,
    pub green: u32,
    pub blue: u32,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Game {
    pub id: u32,
    pub reveals: Vec<Reveal>,
}

impl Game {
    /// Whether every reveal fits within the given cube counts.
    pub fn is_possible(&self, red: u32, green: u32, blue: u32) -> bool {
        self.reveals
            .iter()
            .all(|reveal| reveal.red <= red && reveal.green <= green && reveal.blue <= blue)
    }

    /// The minimum number of cubes of each color needed for every reveal,
    /// multiplied together.
    pub fn power(&self) -> u32 {
        let min = self
            .reveals
            .iter()
            .fold(Reveal::default(), |acc, reveal| Reveal {
                red: acc.red.max(reveal.red),
                green: acc.green.max(reveal.green),
                blue: acc.blue.max(reveal.blue),
            });
        min.red * min.green * min.blue
    }
}

/// Parses one cube count like "3 blue" into a `(count, color)` pair applied to `reveal`.
fn apply_cube(reveal: &mut Reveal, cube: &str) -> Result<(), AocError> {
    let (count, color) = cube
        .trim()
        .split_once(' ')
        .ok_or(AocError::InvalidInput("cube needs a count and a color"))?;
    let count: u32 = count.parse()?;
    match color {
        "red" => reveal.red = count,
        "green" => reveal.green = count,
        "blue" => reveal.blue = count,
        _ => return Err(AocError::InvalidInput("unknown cube color")),
    }
    Ok(())
}

fn parse_game(line: &str) -> Result<Game, AocError> {
    let (header, body) = line
        .split_once(':')
        .ok_or(AocError::InvalidInput("missing ':' separating game id"))?;
    let id: u32 = header
        .trim()
        .strip_prefix("Game ")
        .ok_or(AocError::InvalidInput("missing 'Game ' prefix"))?
        .parse()?;
    let reveals = body
        .split(';')
        .map(|subset| {
            let mut reveal = Reveal::default();
            for cube in subset.split(',') {
                apply_cube(&mut reveal, cube)?;
            }
            Ok(reveal)
        })
        .collect::<Result<Vec<Reveal>, AocError>>()?;
    if reveals.is_empty() {
        return Err(AocError::InvalidInput("game has no reveals"));
    }
    Ok(Game { id, reveals })
}

/// Parses every non-empty line of the puzzle input into a [`Game`].
pub fn parse(input: &str) -> Result<Vec<Game>, AocError> {
    let games = input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(parse_game)
        .collect::<Result<Vec<Game>, AocError>>()?;
    if games.is_empty() {
        return Err(AocError::InvalidInput("no games found"));
    }
    Ok(games)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parses_game_with_multiple_reveals() {
        let game = parse_game("Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green").unwrap();
        assert_eq!(game.id, 3);
        assert_eq!(
            game.reveals,
            vec![
                Reveal {
                    green: 8,
                    blue: 6,
                    red: 20
                },
                Reveal {
                    blue: 5,
                    red: 4,
                    green: 13
                },
            ]
        );
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in [
            "",
            "Game 1",
            "Game: 3 red",
            "Foo 1: 3 red",
            "Game x: 3 red",
            "Game 1: 3 red; ",
            "Game 1: three red",
            "Game 1: 3 purple",
            "Game 1: 3",
        ] {
            assert!(parse(input).is_err(), "{input}");
        }
    }
}
