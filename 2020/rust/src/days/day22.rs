use anyhow::{bail, ensure, Context};
use std::collections::{HashSet, VecDeque};

type Deck = VecDeque<u32>;

/// Parses the two decks from the puzzle input: two "Player N:" blocks
/// separated by a blank line, each followed by one card number per line.
fn parse(input: &str) -> anyhow::Result<(Deck, Deck)> {
    let mut blocks = input.trim().split("\n\n");
    let first = blocks.next().context("missing player 1 block")?;
    let second = blocks.next().context("missing player 2 block")?;
    ensure!(
        blocks.next().is_none(),
        "expected exactly two player blocks"
    );

    Ok((parse_deck(first)?, parse_deck(second)?))
}

fn parse_deck(block: &str) -> anyhow::Result<Deck> {
    let mut lines = block.lines();
    let header = lines.next().context("missing player header")?;
    ensure!(
        header.starts_with("Player") && header.trim_end().ends_with(':'),
        "malformed player header: {header}"
    );

    lines
        .map(|line| {
            line.trim()
                .parse::<u32>()
                .with_context(|| format!("invalid card: {line}"))
        })
        .collect()
}

fn score(deck: &Deck) -> u64 {
    deck.iter()
        .rev()
        .enumerate()
        .map(|(i, card)| u64::from(*card) * (i as u64 + 1))
        .sum()
}

/// Plays standard Combat to completion and returns the winner's score.
fn play_combat(mut player1: Deck, mut player2: Deck) -> u64 {
    while !player1.is_empty() && !player2.is_empty() {
        let (card1, card2) = (player1.pop_front().unwrap(), player2.pop_front().unwrap());
        if card1 > card2 {
            player1.push_back(card1);
            player1.push_back(card2);
        } else {
            player2.push_back(card2);
            player2.push_back(card1);
        }
    }

    score(if player1.is_empty() {
        &player2
    } else {
        &player1
    })
}

enum Winner {
    Player1,
    Player2,
}

/// Plays Recursive Combat to completion, returning the winner and their final deck.
fn play_recursive_combat(mut player1: Deck, mut player2: Deck) -> (Winner, Deck) {
    let mut seen = HashSet::new();

    loop {
        if player1.is_empty() {
            return (Winner::Player2, player2);
        }
        if player2.is_empty() {
            return (Winner::Player1, player1);
        }

        let state = (player1.clone(), player2.clone());
        if !seen.insert(state) {
            return (Winner::Player1, player1);
        }

        let (card1, card2) = (player1.pop_front().unwrap(), player2.pop_front().unwrap());
        let round_winner = if (card1 as usize) <= player1.len() && (card2 as usize) <= player2.len()
        {
            let sub1 = player1.iter().take(card1 as usize).copied().collect();
            let sub2 = player2.iter().take(card2 as usize).copied().collect();
            play_recursive_combat(sub1, sub2).0
        } else if card1 > card2 {
            Winner::Player1
        } else {
            Winner::Player2
        };

        match round_winner {
            Winner::Player1 => {
                player1.push_back(card1);
                player1.push_back(card2);
            }
            Winner::Player2 => {
                player2.push_back(card2);
                player2.push_back(card1);
            }
        }
    }
}

pub fn solve(input: &str) -> anyhow::Result<(u64, u64)> {
    let (player1, player2) = parse(input)?;
    if player1.is_empty() || player2.is_empty() {
        bail!("both players need at least one card");
    }

    let part1 = play_combat(player1.clone(), player2.clone());
    let (_, final_deck) = play_recursive_combat(player1, player2);
    let part2 = score(&final_deck);

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/2020/day22-example.txt");

    #[test]
    fn part1_example() {
        let (part1, _) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 306);
    }

    #[test]
    fn part2_example() {
        let (_, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part2, 291);
    }

    #[test]
    fn part2_terminates_on_infinite_game() {
        // The official "infinite game" edge case from the puzzle description: plain
        // Combat cycles forever on this deck, so only Recursive Combat (with its
        // seen-states rule) is exercised here.
        let (player1, player2) = parse("Player 1:\n43\n19\n\nPlayer 2:\n2\n29\n14").unwrap();
        let (_, final_deck) = play_recursive_combat(player1, player2);
        assert_eq!(score(&final_deck), 105);
    }

    #[test]
    fn malformed_input_missing_block() {
        assert!(solve("Player 1:\n9\n2\n6\n3\n1").is_err());
    }

    #[test]
    fn malformed_input_bad_card() {
        let bad = "Player 1:\n9\nabc\n\nPlayer 2:\n5\n4";
        assert!(solve(bad).is_err());
    }

    #[test]
    fn malformed_input_empty_deck() {
        let bad = "Player 1:\n\nPlayer 2:\n5\n4";
        assert!(solve(bad).is_err());
    }
}
