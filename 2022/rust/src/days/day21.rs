//! Day 21: Monkey Math.
//!
//! Every monkey either yells a number or the result of an operation on two
//! other monkeys' numbers. Part 1 evaluates the riddle rooted at `root`; part 2
//! reads `root` as an equality and solves it for the number `humn` must yell.

use std::collections::HashMap;

use anyhow::{bail, ensure, Context, Result};

const ROOT: &str = "root";
const HUMAN: &str = "humn";

#[derive(Clone, Copy)]
enum Operator {
    Add,
    Subtract,
    Multiply,
    Divide,
}

impl TryFrom<&str> for Operator {
    type Error = anyhow::Error;

    fn try_from(token: &str) -> Result<Self> {
        Ok(match token {
            "+" => Self::Add,
            "-" => Self::Subtract,
            "*" => Self::Multiply,
            "/" => Self::Divide,
            other => bail!("unknown operator: {other}"),
        })
    }
}

impl Operator {
    fn apply(self, left: i64, right: i64) -> Result<i64> {
        match self {
            Self::Add => left.checked_add(right),
            Self::Subtract => left.checked_sub(right),
            Self::Multiply => left.checked_mul(right),
            Self::Divide => left.checked_div(right),
        }
        .with_context(|| format!("operation on {left} and {right} has no result"))
    }
}

#[derive(Clone, Copy)]
enum Job<'a> {
    Number(i64),
    Operation {
        left: &'a str,
        operator: Operator,
        right: &'a str,
    },
}

type Monkeys<'a> = HashMap<&'a str, Job<'a>>;

fn parse(input: &str) -> Result<Monkeys<'_>> {
    let mut monkeys = Monkeys::new();
    for line in input.lines().map(str::trim).filter(|line| !line.is_empty()) {
        let (name, job) = line
            .split_once(':')
            .with_context(|| format!("missing ':' in line: {line}"))?;
        let name = name.trim();
        ensure!(!name.is_empty(), "missing monkey name in line: {line}");

        let tokens: Vec<&str> = job.split_whitespace().collect();
        let job = match tokens.as_slice() {
            [number] => Job::Number(
                number
                    .parse()
                    .with_context(|| format!("invalid number for monkey {name}: {number}"))?,
            ),
            [left, operator, right] => Job::Operation {
                left,
                operator: Operator::try_from(*operator)
                    .with_context(|| format!("invalid job for monkey {name}"))?,
                right,
            },
            _ => bail!("invalid job in line: {line}"),
        };
        ensure!(
            monkeys.insert(name, job).is_none(),
            "duplicate monkey: {name}"
        );
    }
    ensure!(
        monkeys.contains_key(ROOT),
        "input must define monkey {ROOT}"
    );
    Ok(monkeys)
}

/// Yells `name`'s number. `budget` bounds the depth by one step per monkey, so
/// a riddle that refers back to itself fails instead of looping forever.
fn evaluate(monkeys: &Monkeys<'_>, name: &str, budget: usize) -> Result<i64> {
    let budget = budget
        .checked_sub(1)
        .with_context(|| format!("cycle while evaluating monkey {name}"))?;
    match *monkeys
        .get(name)
        .with_context(|| format!("unknown monkey: {name}"))?
    {
        Job::Number(number) => Ok(number),
        Job::Operation {
            left,
            operator,
            right,
        } => operator.apply(
            evaluate(monkeys, left, budget)?,
            evaluate(monkeys, right, budget)?,
        ),
    }
}

/// Maps every monkey to the monkey waiting for its number. Part 2 inverts a
/// single path, so a monkey heard by two monkeys is rejected.
fn listeners<'a>(monkeys: &Monkeys<'a>) -> Result<HashMap<&'a str, &'a str>> {
    let mut listeners = HashMap::new();
    for (&name, job) in monkeys {
        if let Job::Operation { left, right, .. } = *job {
            for operand in [left, right] {
                ensure!(
                    listeners.insert(operand, name).is_none(),
                    "monkey {operand} is heard by two monkeys"
                );
            }
        }
    }
    Ok(listeners)
}

fn path_to_human<'a>(monkeys: &Monkeys<'a>) -> Result<Vec<&'a str>> {
    let listeners = listeners(monkeys)?;
    ensure!(
        monkeys.contains_key(HUMAN),
        "input must define monkey {HUMAN}"
    );

    let mut path = vec![HUMAN];
    let mut current: &str = HUMAN;
    while current != ROOT {
        current = listeners
            .get(current)
            .with_context(|| format!("monkey {current} is not heard by {ROOT}"))?;
        path.push(current);
        ensure!(
            path.len() <= monkeys.len(),
            "cycle between {HUMAN} and {ROOT}"
        );
    }
    path.reverse();
    Ok(path)
}

/// Solves `target = left <operator> right` for the operand on `humn`'s side.
fn invert(operator: Operator, target: i64, known: i64, human_on_left: bool) -> Result<i64> {
    match (operator, human_on_left) {
        (Operator::Add, _) => target.checked_sub(known),
        (Operator::Subtract, true) => target.checked_add(known),
        (Operator::Subtract, false) => known.checked_sub(target),
        (Operator::Multiply, _) => exact_division(target, known),
        (Operator::Divide, true) => target.checked_mul(known),
        (Operator::Divide, false) => exact_division(known, target),
    }
    .with_context(|| format!("no integer operand solves {target} against {known}"))
}

/// Divides only when the division leaves no remainder; anything else means the
/// riddle has no integer answer.
fn exact_division(dividend: i64, divisor: i64) -> Option<i64> {
    match dividend.checked_rem(divisor) {
        Some(0) => dividend.checked_div(divisor),
        _ => None,
    }
}

/// Walks from `root` down to `humn`, inverting every operation on the way: the
/// branch without `humn` evaluates to a known number, which pins down what the
/// branch holding `humn` has to yell.
fn human_value(monkeys: &Monkeys<'_>) -> Result<i64> {
    let path = path_to_human(monkeys)?;
    let budget = monkeys.len();
    let mut target = 0;

    for (&parent, &child) in path.iter().zip(path.iter().skip(1)) {
        let Some(Job::Operation {
            left,
            operator,
            right,
        }) = monkeys.get(parent).copied()
        else {
            bail!("monkey {parent} must be an operation to reach {HUMAN}");
        };
        let human_on_left = left == child;
        let other = if human_on_left { right } else { left };
        let known = evaluate(monkeys, other, budget)?;

        // `root` compares its operands instead of combining them.
        target = if parent == ROOT {
            known
        } else {
            invert(operator, target, known, human_on_left)?
        };
    }
    Ok(target)
}

pub fn solve(input: &str) -> Result<(i64, i64)> {
    let monkeys = parse(input)?;
    Ok((
        evaluate(&monkeys, ROOT, monkeys.len())?,
        human_value(&monkeys)?,
    ))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day21-example.txt");

    #[test]
    fn solve_example() {
        let (part1, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 152);
        assert_eq!(part2, 301);
    }

    #[test]
    fn human_as_direct_operand_of_root() {
        let (part1, part2) = solve("root: humn + other\nother: 7\nhumn: 3\n").unwrap();
        assert_eq!(part1, 10);
        assert_eq!(part2, 7);
    }

    #[test]
    fn inverts_operations_with_human_on_the_right() {
        let input = "root: target + chain\n\
                     target: 10\n\
                     chain: base - inner\n\
                     base: 30\n\
                     inner: fixed / humn\n\
                     fixed: 40\n\
                     humn: 1\n";
        assert_eq!(solve(input).unwrap(), (0, 2));
    }

    #[test]
    fn inverts_multiplication_and_division_with_human_on_the_left() {
        let input = "root: a + b\na: 36\nb: c * three\nthree: 3\nc: humn / two\ntwo: 2\nhumn: 0\n";
        assert_eq!(solve(input).unwrap(), (36, 24));
    }

    #[test]
    fn ignores_blank_lines_and_surrounding_whitespace() {
        let input = "\n  root: humn + other  \n\nother: 7\n humn: 3\n\n";
        assert_eq!(solve(input).unwrap(), (10, 7));
    }

    #[test]
    fn malformed_input_errors() {
        for input in [
            "",
            "root 1\n",
            ": 1\n",
            "root: pppw ? sjmn\nhumn: 1\npppw: 1\nsjmn: 1\n",
            "root: humn + other\nother: twelve\nhumn: 1\n",
            "root: a + b + c\n",
            "root: humn + other\nother: 7\nother: 8\nhumn: 1\n",
            "other: 7\n",
        ] {
            assert!(solve(input).is_err(), "expected an error for {input:?}");
        }
    }

    #[test]
    fn unknown_or_missing_monkeys_error() {
        assert!(solve("root: humn + missing\nhumn: 1\n").is_err());
        assert!(solve("root: a + b\na: 1\nb: 2\n").is_err());
    }

    #[test]
    fn cycles_error_instead_of_looping() {
        assert!(solve("root: a + b\na: 1\nb: b - a\nhumn: 1\n").is_err());
    }

    #[test]
    fn riddle_without_integer_answer_errors() {
        assert!(
            solve("root: odd + scaled\nodd: 7\nscaled: humn * two\ntwo: 2\nhumn: 0\n").is_err()
        );
    }

    #[test]
    fn division_by_zero_errors() {
        assert!(solve("root: a + b\na: 1\nb: a / zero\nzero: 0\nhumn: 1\n").is_err());
    }
}
