use crate::day18::custom_error::AocError;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub(crate) enum Direction {
    Up,
    Down,
    Left,
    Right,
}

impl Direction {
    fn from_letter(letter: char) -> Result<Self, AocError> {
        match letter {
            'U' => Ok(Self::Up),
            'D' => Ok(Self::Down),
            'L' => Ok(Self::Left),
            'R' => Ok(Self::Right),
            _ => Err(AocError::InvalidInput("unknown direction letter")),
        }
    }

    fn from_digit(digit: char) -> Result<Self, AocError> {
        match digit {
            '0' => Ok(Self::Right),
            '1' => Ok(Self::Down),
            '2' => Ok(Self::Left),
            '3' => Ok(Self::Up),
            _ => Err(AocError::InvalidInput("unknown direction digit")),
        }
    }

    pub(crate) fn offset(self) -> (i64, i64) {
        match self {
            Self::Up => (0, -1),
            Self::Down => (0, 1),
            Self::Left => (-1, 0),
            Self::Right => (1, 0),
        }
    }
}

#[derive(Debug, Clone, Copy)]
pub(crate) struct Step {
    pub direction: Direction,
    pub distance: i64,
}

#[derive(Debug)]
pub(crate) struct DigPlan {
    pub steps: Vec<Step>,
}

impl DigPlan {
    pub fn parse(input: &str, use_colour: bool) -> Result<Self, AocError> {
        let steps = input
            .lines()
            .map(str::trim)
            .filter(|line| !line.is_empty())
            .map(|line| Self::parse_line(line, use_colour))
            .collect::<Result<Vec<Step>, AocError>>()?;
        if steps.is_empty() {
            return Err(AocError::InvalidInput("empty dig plan"));
        }
        Ok(Self { steps })
    }

    fn parse_line(line: &str, use_colour: bool) -> Result<Step, AocError> {
        let mut fields = line.split_whitespace();
        let letter = fields
            .next()
            .ok_or(AocError::InvalidInput("missing direction"))?;
        let distance = fields
            .next()
            .ok_or(AocError::InvalidInput("missing distance"))?;
        let colour = fields
            .next()
            .ok_or(AocError::InvalidInput("missing colour"))?;
        if fields.next().is_some() {
            return Err(AocError::InvalidInput("too many fields"));
        }

        if use_colour {
            Self::parse_colour(colour)
        } else {
            let mut letters = letter.chars();
            let direction = letters
                .next()
                .ok_or(AocError::InvalidInput("missing direction letter"))?;
            if letters.next().is_some() {
                return Err(AocError::InvalidInput("direction must be one letter"));
            }
            Ok(Step {
                direction: Direction::from_letter(direction)?,
                distance: distance.parse()?,
            })
        }
    }

    fn parse_colour(colour: &str) -> Result<Step, AocError> {
        let hex = colour
            .strip_prefix("(#")
            .and_then(|rest| rest.strip_suffix(')'))
            .ok_or(AocError::InvalidInput("malformed colour"))?;
        if hex.len() != 6 {
            return Err(AocError::InvalidInput("colour must be six hex digits"));
        }
        let distance = i64::from_str_radix(&hex[..5], 16)
            .map_err(|_| AocError::InvalidInput("colour distance is not hexadecimal"))?;
        let direction_digit = hex
            .chars()
            .nth(5)
            .ok_or(AocError::InvalidInput("colour missing direction digit"))?;
        Ok(Step {
            direction: Direction::from_digit(direction_digit)?,
            distance,
        })
    }
}

/// Signed area of the lagoon's trench polygon via the shoelace formula, then
/// Pick's theorem (`A = inside + boundary / 2 - 1`) to recover every dug
/// tile, including the boundary itself.
pub(crate) fn lagoon_area(steps: &[Step]) -> i64 {
    let (mut x, mut y) = (0i64, 0i64);
    let mut shoelace_sum = 0i64;
    let mut boundary = 0i64;
    for step in steps {
        let (dx, dy) = step.direction.offset();
        let (next_x, next_y) = (x + dx * step.distance, y + dy * step.distance);
        shoelace_sum += x * next_y - next_x * y;
        boundary += step.distance;
        x = next_x;
        y = next_y;
    }
    let area = shoelace_sum.abs() / 2;
    area + boundary / 2 + 1
}
