use anyhow::{ensure, Context};

#[derive(Debug, Clone, Copy)]
struct Seat {
    row: u32,
    col: u32,
}

impl Seat {
    fn id(&self) -> u32 {
        self.row * 8 + self.col
    }
}

fn parse_seat(line: &str) -> anyhow::Result<Seat> {
    ensure!(line.len() == 10, "boarding pass must be 10 characters");

    let (row_part, col_part) = line.split_at(7);

    for c in row_part.chars() {
        ensure!(c == 'F' || c == 'B', "row part must contain only F or B");
    }
    for c in col_part.chars() {
        ensure!(c == 'L' || c == 'R', "column part must contain only L or R");
    }

    let row = row_part
        .chars()
        .fold(0u32, |acc, c| acc * 2 + if c == 'B' { 1 } else { 0 });

    let col = col_part
        .chars()
        .fold(0u32, |acc, c| acc * 2 + if c == 'R' { 1 } else { 0 });

    ensure!(row < 128, "row must be < 128, got {}", row);
    ensure!(col < 8, "column must be < 8, got {}", col);

    Ok(Seat { row, col })
}

fn seat_ids(input: &str) -> anyhow::Result<Vec<u32>> {
    let ids = input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(|line| parse_seat(line).map(|seat| seat.id()))
        .collect::<anyhow::Result<Vec<_>>>()?;
    ensure!(!ids.is_empty(), "no boarding passes found");
    Ok(ids)
}

fn free_seat(sorted_ids: &[u32]) -> Option<u32> {
    sorted_ids
        .windows(2)
        .find(|pair| pair[1] - pair[0] == 2)
        .map(|pair| pair[0] + 1)
}

pub fn solve(input: &str) -> anyhow::Result<(u32, u32)> {
    let mut ids = seat_ids(input)?;
    ids.sort_unstable();
    let highest = *ids.last().context("no boarding passes found")?;
    let free = free_seat(&ids).context("no free seat between two occupied seats")?;
    Ok((highest, free))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_examples() {
        let tests = vec![
            ("FBFBBFFRLR", 44, 5, 357),
            ("BFFFBBFRRR", 70, 7, 567),
            ("FFFBBBFRRR", 14, 7, 119),
            ("BBFFBBFRLL", 102, 4, 820),
        ];

        for (input, expected_row, expected_col, expected_id) in tests {
            let seat = parse_seat(input).unwrap();
            assert_eq!(seat.row, expected_row);
            assert_eq!(seat.col, expected_col);
            assert_eq!(seat.id(), expected_id);
        }
    }

    #[test]
    fn test_part1_example() {
        let input = include_str!("../../input/2020/day5-example.txt");
        let ids = seat_ids(input).unwrap();
        assert_eq!(ids.into_iter().max(), Some(820));
    }

    #[test]
    fn test_part2_missing_seat() {
        let input = "FFFBBFFRLL\nFFFBBFFRLR\nFFFBBFFRRR";
        let (_, part2) = solve(input).unwrap();
        assert_eq!(part2, 102);
    }

    #[test]
    fn test_malformed_input_wrong_length() {
        let result = parse_seat("FBFBBFFRL");
        assert!(result.is_err());
    }

    #[test]
    fn test_malformed_input_invalid_chars() {
        let result = parse_seat("FBFBBFFXLR");
        assert!(result.is_err());
    }

    #[test]
    fn test_empty_input() {
        let result = solve("");
        assert!(result.is_err());
    }

    #[test]
    fn test_no_free_seat_is_an_error() {
        assert_eq!(seat_ids("FBFBBFFRLR").unwrap(), vec![357]);
        assert!(solve("FBFBBFFRLR").is_err());
    }
}
