use anyhow::{ensure, Context, Result};
use std::collections::HashMap;

/// A player's position (1-10) and score.
#[derive(Clone, Copy, PartialEq, Eq, Hash)]
struct Player {
    position: u8,
    score: u16,
}

fn parse_player(line: &str, expected: u32) -> Result<Player> {
    let (label, position) = line
        .split_once(": ")
        .with_context(|| format!("malformed starting position line: {line}"))?;
    ensure!(
        label.trim() == format!("Player {expected} starting position"),
        "expected player {expected}, got: {label}"
    );
    let position: u8 = position
        .trim()
        .parse()
        .context("invalid starting position")?;
    ensure!(
        (1..=10).contains(&position),
        "starting position must be between 1 and 10"
    );
    Ok(Player { position, score: 0 })
}

fn parse(input: &str) -> Result<(Player, Player)> {
    let mut lines = input.lines().filter(|line| !line.trim().is_empty());
    let player1 = parse_player(lines.next().context("missing player 1 line")?, 1)?;
    let player2 = parse_player(lines.next().context("missing player 2 line")?, 2)?;
    ensure!(lines.next().is_none(), "unexpected extra input");
    Ok((player1, player2))
}

fn advance(position: u8, roll: u32) -> u8 {
    (((position as u32 - 1) + roll) % 10 + 1) as u8
}

/// Part 1: deterministic 100-sided die, three rolls per turn.
fn play_deterministic(mut players: [Player; 2]) -> u64 {
    let mut die = (1..=100).cycle();
    let mut rolls = 0u64;
    let mut turn = 0usize;
    loop {
        let total: u32 = (0..3).map(|_| die.next().unwrap()).sum();
        rolls += 3;
        let player = &mut players[turn];
        player.position = advance(player.position, total);
        player.score += player.position as u16;
        if player.score >= 1000 {
            let loser = &players[1 - turn];
            return loser.score as u64 * rolls;
        }
        turn = 1 - turn;
    }
}

/// Multiplicities of each possible sum of three 3-sided dirac dice rolls.
const ROLL_SUMS: [(u32, u64); 7] = [(3, 1), (4, 3), (5, 6), (6, 7), (7, 6), (8, 3), (9, 1)];

type Memo = HashMap<(Player, Player), (u64, u64)>;

/// Returns (wins for the player to move, wins for the other player).
fn count_wins(current: Player, other: Player, memo: &mut Memo) -> (u64, u64) {
    if other.score >= 21 {
        return (0, 1);
    }
    if let Some(&result) = memo.get(&(current, other)) {
        return result;
    }
    let mut wins_current = 0u64;
    let mut wins_other = 0u64;
    for (roll, multiplicity) in ROLL_SUMS {
        let position = advance(current.position, roll);
        let score = current.score + position as u16;
        let moved = Player { position, score };
        let (other_wins, current_wins) = count_wins(other, moved, memo);
        wins_current += multiplicity * current_wins;
        wins_other += multiplicity * other_wins;
    }
    memo.insert((current, other), (wins_current, wins_other));
    (wins_current, wins_other)
}

/// Part 2: quantum dirac dice, every universe explored, first to 21 wins.
fn play_dirac(players: [Player; 2]) -> u64 {
    let mut memo = Memo::new();
    let (wins1, wins2) = count_wins(players[0], players[1], &mut memo);
    wins1.max(wins2)
}

pub fn solve(input: &str) -> Result<(u64, u64)> {
    let (player1, player2) = parse(input)?;
    let part1 = play_deterministic([player1, player2]);
    let part2 = play_dirac([player1, player2]);
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day21-example.txt");
        assert_eq!(solve(input).unwrap(), (739785, 444356092776315));
    }

    #[test]
    fn advance_wraps_around_the_track() {
        assert_eq!(advance(7, 5), 2);
        assert_eq!(advance(10, 1), 1);
    }

    #[test]
    fn starting_position_must_be_on_the_track() {
        let input = "Player 1 starting position: 0\nPlayer 2 starting position: 8\n";
        assert!(solve(input).is_err());
    }

    #[test]
    fn malformed_input_is_rejected() {
        assert!(solve("not a valid game").is_err());
        assert!(solve("Player 1 starting position: 4\n").is_err());
    }
}
