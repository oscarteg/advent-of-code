use anyhow::{bail, ensure, Context, Result};

#[derive(Clone, Copy, PartialEq, Eq, Debug)]
enum Token {
    Num(i64),
    Plus,
    Star,
    LParen,
    RParen,
}

fn tokenize(line: &str) -> Result<Vec<Token>> {
    let mut tokens = Vec::new();
    let mut chars = line.chars().peekable();

    while let Some(&c) = chars.peek() {
        match c {
            ' ' => {
                chars.next();
            }
            '+' => {
                tokens.push(Token::Plus);
                chars.next();
            }
            '*' => {
                tokens.push(Token::Star);
                chars.next();
            }
            '(' => {
                tokens.push(Token::LParen);
                chars.next();
            }
            ')' => {
                tokens.push(Token::RParen);
                chars.next();
            }
            d if d.is_ascii_digit() => {
                let mut digits = String::new();
                while let Some(&d2) = chars.peek() {
                    if d2.is_ascii_digit() {
                        digits.push(d2);
                        chars.next();
                    } else {
                        break;
                    }
                }
                let num = digits
                    .parse()
                    .with_context(|| format!("invalid number {digits:?} in {line:?}"))?;
                tokens.push(Token::Num(num));
            }
            other => bail!("unexpected character {other:?} in {line:?}"),
        }
    }

    ensure!(!tokens.is_empty(), "empty expression in {line:?}");
    Ok(tokens)
}

/// Recursive-descent parser. When `addition_binds_tighter` is set, `+`
/// binds more tightly than `*` (part 2); otherwise both are left-to-right
/// with equal precedence (part 1).
struct Parser<'a> {
    tokens: &'a [Token],
    pos: usize,
    addition_binds_tighter: bool,
}

impl<'a> Parser<'a> {
    fn peek(&self) -> Option<Token> {
        self.tokens.get(self.pos).copied()
    }

    fn next_token(&mut self) -> Result<Token> {
        let token = self
            .peek()
            .context("unexpected end of expression, expected a token")?;
        self.pos += 1;
        Ok(token)
    }

    fn parse_expr(&mut self) -> Result<i64> {
        if self.addition_binds_tighter {
            let mut value = self.parse_sum()?;
            while self.peek() == Some(Token::Star) {
                self.pos += 1;
                value *= self.parse_sum()?;
            }
            Ok(value)
        } else {
            self.parse_flat()
        }
    }

    /// Chain of `+` only, used when addition binds tighter.
    fn parse_sum(&mut self) -> Result<i64> {
        let mut value = self.parse_factor()?;
        while self.peek() == Some(Token::Plus) {
            self.pos += 1;
            value += self.parse_factor()?;
        }
        Ok(value)
    }

    /// Left-to-right chain of `+`/`*` with equal precedence, used for part 1.
    fn parse_flat(&mut self) -> Result<i64> {
        let mut value = self.parse_factor()?;
        loop {
            match self.peek() {
                Some(Token::Plus) => {
                    self.pos += 1;
                    value += self.parse_factor()?;
                }
                Some(Token::Star) => {
                    self.pos += 1;
                    value *= self.parse_factor()?;
                }
                _ => break,
            }
        }
        Ok(value)
    }

    fn parse_factor(&mut self) -> Result<i64> {
        match self.next_token()? {
            Token::Num(n) => Ok(n),
            Token::LParen => {
                let value = self.parse_expr()?;
                ensure!(
                    self.next_token()? == Token::RParen,
                    "expected closing parenthesis"
                );
                Ok(value)
            }
            other => bail!("unexpected token {other:?}, expected a number or '('"),
        }
    }
}

fn eval(tokens: &[Token], addition_binds_tighter: bool) -> Result<i64> {
    let mut parser = Parser {
        tokens,
        pos: 0,
        addition_binds_tighter,
    };
    let value = parser.parse_expr()?;
    ensure!(
        parser.pos == tokens.len(),
        "trailing tokens after a complete expression"
    );
    Ok(value)
}

fn lines(input: &str) -> Result<Vec<&str>> {
    let lines: Vec<&str> = input
        .lines()
        .map(str::trim)
        .filter(|l| !l.is_empty())
        .collect();
    ensure!(!lines.is_empty(), "no expressions found in input");
    Ok(lines)
}

pub fn solve(input: &str) -> Result<(i64, i64)> {
    let mut part1 = 0i64;
    let mut part2 = 0i64;

    for line in lines(input)? {
        let tokens = tokenize(line)?;
        part1 += eval(&tokens, false)?;
        part2 += eval(&tokens, true)?;
    }

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/2020/day18-example.txt");

    fn line_part1(line: &str) -> i64 {
        eval(&tokenize(line).unwrap(), false).unwrap()
    }

    fn line_part2(line: &str) -> i64 {
        eval(&tokenize(line).unwrap(), true).unwrap()
    }

    #[test]
    fn part1_per_line_examples() {
        let lines: Vec<&str> = EXAMPLE.lines().collect();
        let expected = [71, 51, 26, 437, 12240, 13632];
        for (line, want) in lines.iter().zip(expected) {
            assert_eq!(line_part1(line), want, "line {line:?}");
        }
    }

    #[test]
    fn part2_per_line_examples() {
        let lines: Vec<&str> = EXAMPLE.lines().collect();
        let expected = [231, 51, 46, 1445, 669060, 23340];
        for (line, want) in lines.iter().zip(expected) {
            assert_eq!(line_part2(line), want, "line {line:?}");
        }
    }

    #[test]
    fn sample_totals() {
        let (part1, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 71 + 51 + 26 + 437 + 12240 + 13632);
        assert_eq!(part2, 231 + 51 + 46 + 1445 + 669060 + 23340);
    }

    #[test]
    fn single_number_is_itself() {
        assert_eq!(line_part1("42"), 42);
        assert_eq!(line_part2("42"), 42);
    }

    #[test]
    fn malformed_unknown_character_is_error() {
        assert!(tokenize("1 + a").is_err());
    }

    #[test]
    fn malformed_trailing_operator_is_error() {
        let tokens = tokenize("1 +").unwrap();
        assert!(eval(&tokens, false).is_err());
    }

    #[test]
    fn malformed_unmatched_parenthesis_is_error() {
        let tokens = tokenize("(1 + 2").unwrap();
        assert!(eval(&tokens, false).is_err());
    }

    #[test]
    fn malformed_trailing_tokens_is_error() {
        let tokens = tokenize("(1 + 2))").unwrap();
        assert!(eval(&tokens, false).is_err());
    }

    #[test]
    fn empty_input_is_error() {
        assert!(solve("").is_err());
    }
}
