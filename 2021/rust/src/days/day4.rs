use anyhow::{anyhow, bail, ensure, Context, Result};

struct Board {
    numbers: [u32; 25],
    marked: [bool; 25],
    won: bool,
}

impl Board {
    fn parse(input: &str) -> Result<Self> {
        let rows: Vec<_> = input.lines().collect();
        ensure!(rows.len() == 5, "a bingo board must have five rows");
        let mut numbers = Vec::with_capacity(25);
        for row in rows {
            let values = row
                .split_whitespace()
                .map(str::parse)
                .collect::<Result<Vec<u32>, _>>()?;
            ensure!(values.len() == 5, "a bingo row must have five numbers");
            numbers.extend(values);
        }
        Ok(Self {
            numbers: numbers
                .try_into()
                .map_err(|_| anyhow!("a bingo board must have 25 numbers"))?,
            marked: [false; 25],
            won: false,
        })
    }

    fn mark(&mut self, draw: u32) -> Result<Option<u64>> {
        if self.won {
            return Ok(None);
        }
        for (number, marked) in self.numbers.iter().zip(&mut self.marked) {
            *marked |= *number == draw;
        }
        self.won = (0..5).any(|line| {
            (0..5).all(|column| self.marked[line * 5 + column])
                || (0..5).all(|row| self.marked[row * 5 + line])
        });
        if !self.won {
            return Ok(None);
        }
        let unmarked: u64 = self
            .numbers
            .iter()
            .zip(&self.marked)
            .filter(|(_, marked)| !**marked)
            .map(|(&number, _)| u64::from(number))
            .sum();
        Ok(Some(
            unmarked
                .checked_mul(u64::from(draw))
                .context("bingo score overflows")?,
        ))
    }
}

/// Returns the first and last winning scores.
pub fn solve(input: &str) -> Result<(u64, u64)> {
    let normalized = input.replace("\r\n", "\n");
    let mut sections = normalized.trim().split("\n\n");
    let draws = sections
        .next()
        .context("missing draws")?
        .split(',')
        .map(|value| value.trim().parse())
        .collect::<Result<Vec<u32>, _>>()?;
    let mut boards = sections.map(Board::parse).collect::<Result<Vec<_>>>()?;
    ensure!(!boards.is_empty(), "missing bingo boards");
    let mut remaining = boards.len();
    let mut first = None;
    for draw in draws {
        for board in &mut boards {
            if let Some(score) = board.mark(draw)? {
                let first_score = *first.get_or_insert(score);
                remaining -= 1;
                if remaining == 0 {
                    return Ok((first_score, score));
                }
            }
        }
    }
    bail!("draws end before every board wins")
}

#[cfg(test)]
mod tests {
    use super::*;

    const BOARD: &str = "0 1 2 3 4\n5 6 7 8 9\n10 11 12 13 14\n15 16 17 18 19\n20 21 22 23 24";

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day4-example.txt");
        assert_eq!(solve(input).unwrap(), (4512, 1924));
    }

    #[test]
    fn repeated_draws_do_not_change_the_score() {
        let input = format!("0,0,1,2,3,4\n\n{BOARD}");
        assert_eq!(solve(&input).unwrap(), (1160, 1160));
    }

    #[test]
    fn completed_columns_win() {
        let input = format!("0,5,10,15,20\n\n{BOARD}");
        assert_eq!(solve(&input).unwrap(), (5000, 5000));
    }

    #[test]
    fn diagonals_do_not_win() {
        let input = format!("0,6,12,18,24\n\n{BOARD}");
        assert!(solve(&input).is_err());
    }

    #[test]
    fn windows_newlines_are_accepted() {
        let input = format!("0,1,2,3,4\n\n{BOARD}").replace('\n', "\r\n");
        assert_eq!(solve(&input).unwrap(), (1160, 1160));
    }

    #[test]
    fn malformed_boards_are_rejected() {
        let input = "1,2,3,4,5\n\n1 2 3 4 5";
        assert!(solve(input).is_err());
    }
}
