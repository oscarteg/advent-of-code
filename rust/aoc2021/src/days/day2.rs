use std::str::FromStr;

type Error = String;

pub enum Direction {
    Forward(u32),
    Down(u32),
    Up(u32),
}

impl FromStr for Direction {
    type Err = Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let (action, value) = s.split_once(' ').ok_or(format!("Wrong format: `{}`", s))?;
        let value = value.parse::<u32>().map_err(|e| e.to_string())?;
        match action {
            "down" => Ok(self::Direction::Down(value)),
            "forward" => Ok(self::Direction::Forward(value)),
            "up" => Ok(self::Direction::Up(value)),
            _ => Err(format!("Wrong direction: {}", s)),
        }
    }
}

pub fn part_a<'a>(directions: impl Iterator<Item = &'a Direction>) -> (u32, u32) {
    // Fold = Reduce, Accumulator are the coordinates
    let (x, y) = directions.fold((0, 0), |(x, y), dir| match dir {
        Direction::Forward(v) => (x + v, y),
        Direction::Down(v) => (x, y + v),
        // saturating_sub makes sure you cant go past the minimum of 0
        Direction::Up(v) => (x, y.saturating_sub(*v)),
    });
    (x, y)
}

pub fn part_b<'a>(directions: impl Iterator<Item = &'a Direction>) -> (u32, u32) {
    // Fold = Reduce, Accumulator are the coordinates
    let (x, y, _) = directions.fold((0, 0, 0), |(x, y, aim), dir| match dir {
        Direction::Forward(v) => (x + v, y + aim * v, aim),
        Direction::Down(v) => (x, y, aim + v),
        Direction::Up(v) => (x, y, aim.saturating_sub(*v)),
    });
    (x, y)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::utils::{clean_input, read_file};

    const INPUT: &'static str = "
forward 5
down 5
forward 8
up 3
down 8
forward 2";

    #[test]
    pub fn test_part_a_example() {
        let input = clean_input(INPUT);

        let directions: Vec<Direction> =
            input.map(|str| Direction::from_str(str).unwrap()).collect();

        let (x, y) = part_a(directions.iter());

        assert_eq!(x * y, 150);
    }

    #[test]
    fn test_part_a() {
        let file = read_file("input/day2.txt");

        let input = clean_input(file.as_str());

        let directions: Vec<Direction> =
            input.map(|str| Direction::from_str(str).unwrap()).collect();

        let (x, y) = part_a(directions.iter());

        assert_eq!(x * y, 2039256);
    }

    #[test]
    fn test_part_b_example() {
        let input = clean_input(INPUT);

        let directions: Vec<Direction> =
            input.map(|str| Direction::from_str(str).unwrap()).collect();

        let (x, y) = part_b(directions.iter());

        assert_eq!(x * y, 900);
    }

    #[test]
    fn test_part_b() {
        let file = read_file("input/day2.txt");

        let input = clean_input(file.as_str());

        let directions: Vec<Direction> =
            input.map(|str| Direction::from_str(str).unwrap()).collect();

        let (x, y) = part_b(directions.iter());

        assert_eq!(x * y, 1856459736);
    }
}
