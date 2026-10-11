use anyhow::{bail, Result};

#[derive(Clone, Copy, Debug, PartialEq, Eq)]
enum Bracket {
    Paren,
    Square,
    Curly,
    Angle,
}

impl Bracket {
    fn from_open(ch: char) -> Option<Self> {
        match ch {
            '(' => Some(Bracket::Paren),
            '[' => Some(Bracket::Square),
            '{' => Some(Bracket::Curly),
            '<' => Some(Bracket::Angle),
            _ => None,
        }
    }

    fn from_close(ch: char) -> Option<Self> {
        match ch {
            ')' => Some(Bracket::Paren),
            ']' => Some(Bracket::Square),
            '}' => Some(Bracket::Curly),
            '>' => Some(Bracket::Angle),
            _ => None,
        }
    }

    fn error_score(self) -> u64 {
        match self {
            Bracket::Paren => 3,
            Bracket::Square => 57,
            Bracket::Curly => 1197,
            Bracket::Angle => 25137,
        }
    }

    fn completion_score(self) -> u64 {
        match self {
            Bracket::Paren => 1,
            Bracket::Square => 2,
            Bracket::Curly => 3,
            Bracket::Angle => 4,
        }
    }
}

#[derive(Debug, PartialEq, Eq)]
enum LineStatus {
    Corrupted(u64),           // First mismatched bracket score
    Incomplete(Vec<Bracket>), // Stack of unclosed brackets
}

fn analyze_line(line: &str) -> Result<LineStatus> {
    let mut stack: Vec<Bracket> = Vec::new();

    for ch in line.chars() {
        if let Some(bracket) = Bracket::from_open(ch) {
            stack.push(bracket);
        } else if let Some(bracket) = Bracket::from_close(ch) {
            match stack.pop() {
                Some(open) if open == bracket => {}
                Some(_) => {
                    // Corrupted: wrong closing bracket
                    return Ok(LineStatus::Corrupted(bracket.error_score()));
                }
                None => {
                    // Extra closing bracket
                    return Ok(LineStatus::Corrupted(bracket.error_score()));
                }
            }
        } else if !ch.is_whitespace() {
            bail!("invalid character: {}", ch);
        }
    }

    if stack.is_empty() {
        bail!("complete line should not appear");
    }
    Ok(LineStatus::Incomplete(stack))
}

pub fn solve(input: &str) -> Result<(u64, u64)> {
    let lines: Vec<&str> = input.lines().filter(|l| !l.trim().is_empty()).collect();

    let mut part1 = 0u64;
    let mut completion_scores = Vec::new();

    for line in lines {
        match analyze_line(line)? {
            LineStatus::Corrupted(score) => {
                part1 += score;
            }
            LineStatus::Incomplete(mut stack) => {
                // Calculate completion score by closing all unclosed brackets
                let mut score = 0u64;
                while let Some(bracket) = stack.pop() {
                    score = score * 5 + bracket.completion_score();
                }
                completion_scores.push(score);
            }
        }
    }

    completion_scores.sort_unstable();
    let Some(&part2) = completion_scores.get(completion_scores.len() / 2) else {
        bail!("no incomplete lines to score");
    };

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day10-example.txt");

    #[test]
    fn test_example() {
        assert_eq!(solve(EXAMPLE).unwrap(), (26397, 288957));
    }

    #[test]
    fn test_single_corrupted() {
        assert_eq!(
            analyze_line("{([(<{}[<>[]}>{[]{[(<()>").unwrap(),
            LineStatus::Corrupted(1197)
        );
    }

    #[test]
    fn test_incomplete_line() {
        let result = analyze_line("[({(<(())[]>[[{[]{<()<>>").unwrap();
        // Should be incomplete
        assert!(matches!(result, LineStatus::Incomplete(_)));
    }

    #[test]
    fn test_malformed_empty_stack() {
        let result = analyze_line(")");
        let status = result.unwrap();
        assert_eq!(status, LineStatus::Corrupted(3));
    }

    #[test]
    fn test_invalid_character_is_rejected() {
        assert!(solve("[(x)]").is_err());
    }

    #[test]
    fn test_input_without_incomplete_lines_is_rejected() {
        assert!(solve("{([(<{}[<>[]}>{[]{[(<()>").is_err());
    }
}
