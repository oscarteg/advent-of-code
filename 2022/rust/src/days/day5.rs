use anyhow::{bail, ensure, Context, Result};
use std::collections::VecDeque;

type Stack = VecDeque<char>;

#[derive(Debug, Clone, Copy)]
struct Move {
    count: usize,
    from: usize,
    to: usize,
}

#[derive(Debug, Clone)]
struct Drawing {
    stacks: Vec<Stack>,
    moves: Vec<Move>,
}

/// Parses the crate drawing and move list. Crate letters sit at fixed
/// columns: index 1, 5, 9, ... (every 4th column), so stacks are read by
/// character position rather than by splitting on whitespace.
fn parse_input(input: &str) -> Result<Drawing> {
    let (stacks_input, moves_input) = input
        .split_once("\n\n")
        .context("missing blank line separating stacks from moves")?;

    let mut lines = stacks_input.lines().collect::<Vec<_>>();
    let labels = lines.pop().context("missing stack label row")?;
    let num_stacks = labels.split_whitespace().count();
    ensure!(num_stacks > 0, "no stack labels found");

    let mut stacks = vec![Stack::new(); num_stacks];
    for line in lines.iter().rev() {
        for (index, stack) in stacks.iter_mut().enumerate() {
            let column = 1 + index * 4;
            if let Some(ch) = line.chars().nth(column) {
                if ch != ' ' {
                    stack.push_back(ch);
                }
            }
        }
    }

    let moves = moves_input
        .lines()
        .filter(|line| !line.trim().is_empty())
        .map(parse_move)
        .collect::<Result<Vec<_>>>()?;

    Ok(Drawing { stacks, moves })
}

fn parse_move(line: &str) -> Result<Move> {
    let parts = line.split_whitespace().collect::<Vec<_>>();
    ensure!(
        parts.len() == 6 && parts[0] == "move" && parts[2] == "from" && parts[4] == "to",
        "malformed move line: {line:?}"
    );
    let count = parts[1]
        .parse()
        .with_context(|| format!("invalid move count in {line:?}"))?;
    let from: usize = parts[3]
        .parse()
        .with_context(|| format!("invalid source stack in {line:?}"))?;
    let to: usize = parts[5]
        .parse()
        .with_context(|| format!("invalid destination stack in {line:?}"))?;
    ensure!(
        from >= 1 && to >= 1,
        "stack numbers are 1-indexed: {line:?}"
    );
    Ok(Move {
        count,
        from: from - 1,
        to: to - 1,
    })
}

fn apply_moves(mut stacks: Vec<Stack>, moves: &[Move], preserve_order: bool) -> Result<Vec<Stack>> {
    for mv in moves {
        ensure!(
            mv.from < stacks.len() && mv.to < stacks.len(),
            "move references unknown stack"
        );
        let from_len = stacks[mv.from].len();
        ensure!(mv.count <= from_len, "not enough crates to move: {mv:?}");

        let split_at = from_len - mv.count;
        let mut moved: Vec<char> = stacks[mv.from].split_off(split_at).into();
        if !preserve_order {
            moved.reverse();
        }
        stacks[mv.to].extend(moved);
    }
    Ok(stacks)
}

fn top_crates(stacks: &[Stack]) -> String {
    stacks.iter().filter_map(|stack| stack.back()).collect()
}

fn solve_with(input: &str, preserve_order_part2: bool) -> Result<(String, String)> {
    let drawing = parse_input(input)?;
    let part1 = top_crates(&apply_moves(drawing.stacks.clone(), &drawing.moves, false)?);
    let part2 = top_crates(&apply_moves(
        drawing.stacks,
        &drawing.moves,
        preserve_order_part2,
    )?);
    Ok((part1, part2))
}

pub fn solve(input: &str) -> Result<(String, String)> {
    if parse_input(input)?.moves.is_empty() && input.trim().is_empty() {
        bail!("empty input");
    }
    solve_with(input, true)
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "    [D]    
[N] [C]    
[Z] [M] [P]
 1   2   3 

move 1 from 2 to 1
move 3 from 1 to 3
move 2 from 2 to 1
move 1 from 1 to 2";

    #[test]
    fn test_part_1_example() {
        let (part1, _) = solve_with(EXAMPLE, true).unwrap();
        assert_eq!(part1, "CMZ");
    }

    #[test]
    fn test_part_2_example() {
        let (_, part2) = solve_with(EXAMPLE, true).unwrap();
        assert_eq!(part2, "MCD");
    }

    #[test]
    fn test_part_1_real_input() {
        let file = std::fs::read_to_string("input/day5.txt").unwrap();
        let (part1, _) = solve(&file).unwrap();
        assert_eq!(part1, "ZSQVCCJLL");
    }

    #[test]
    fn test_single_crate_single_move() {
        let input = "[A]\n 1 \n\nmove 1 from 1 to 1";
        let (part1, part2) = solve_with(input, true).unwrap();
        assert_eq!(part1, "A");
        assert_eq!(part2, "A");
    }

    #[test]
    fn test_malformed_input_missing_blank_line() {
        let input = "[A]\n 1 \nmove 1 from 1 to 1";
        assert!(parse_input(input).is_err());
    }

    #[test]
    fn test_malformed_move_line() {
        let input = "[A]\n 1 \n\nmove one from 1 to 1";
        assert!(parse_input(input).is_err());
    }

    #[test]
    fn test_move_too_many_crates() {
        let input = "[A]\n 1 \n\nmove 5 from 1 to 1";
        let drawing = parse_input(input).unwrap();
        assert!(apply_moves(drawing.stacks, &drawing.moves, false).is_err());
    }
}
