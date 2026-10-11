//! Day 11: Monkey in the Middle.
//!
//! Each monkey inspects its items with a worry operation, tests divisibility,
//! and throws the item to another monkey. Part 1 relieves worry by dividing
//! by three after each inspection; part 2 instead keeps worry bounded by
//! reducing it modulo the product of every monkey's divisor.

use anyhow::{bail, ensure, Context, Result};
use std::collections::VecDeque;

const ROUNDS_PART1: usize = 20;
const ROUNDS_PART2: usize = 10_000;
const RELIEF_DIVISOR: u64 = 3;

#[derive(Clone, Copy)]
enum Operand {
    Old,
    Value(u64),
}

#[derive(Clone, Copy)]
enum Operator {
    Add,
    Multiply,
}

#[derive(Clone)]
struct Operation {
    operator: Operator,
    operand: Operand,
}

impl Operation {
    fn apply(&self, old: u64) -> u64 {
        let rhs = match self.operand {
            Operand::Old => old,
            Operand::Value(v) => v,
        };
        match self.operator {
            Operator::Add => old + rhs,
            Operator::Multiply => old * rhs,
        }
    }
}

#[derive(Clone)]
struct Monkey {
    items: VecDeque<u64>,
    operation: Operation,
    divisor: u64,
    if_true: usize,
    if_false: usize,
    inspections: u64,
}

fn parse_suffix_numbers(line: &str, prefix: &str) -> Result<Vec<u64>> {
    let rest = line
        .trim()
        .strip_prefix(prefix)
        .with_context(|| format!("expected line to start with '{prefix}': {line}"))?;
    rest.split(',')
        .map(|n| {
            n.trim()
                .parse::<u64>()
                .with_context(|| format!("invalid number in '{line}'"))
        })
        .collect()
}

fn parse_suffix_number(line: &str, prefix: &str) -> Result<u64> {
    let rest = line
        .trim()
        .strip_prefix(prefix)
        .with_context(|| format!("expected line to start with '{prefix}': {line}"))?;
    rest.trim()
        .parse::<u64>()
        .with_context(|| format!("invalid number in '{line}'"))
}

fn parse_operation(line: &str) -> Result<Operation> {
    let rest = line
        .trim()
        .strip_prefix("Operation: new = old ")
        .with_context(|| format!("invalid operation line: {line}"))?;
    let mut parts = rest.split_whitespace();
    let operator = match parts.next() {
        Some("+") => Operator::Add,
        Some("*") => Operator::Multiply,
        other => bail!("unsupported operator in operation line: {other:?}"),
    };
    let operand = match parts.next() {
        Some("old") => Operand::Old,
        Some(value) => Operand::Value(
            value
                .parse::<u64>()
                .with_context(|| format!("invalid operand in '{line}'"))?,
        ),
        None => bail!("missing operand in operation line: {line}"),
    };
    Ok(Operation { operator, operand })
}

fn parse_monkey(block: &str) -> Result<Monkey> {
    let lines: Vec<&str> = block
        .lines()
        .map(str::trim)
        .filter(|l| !l.is_empty())
        .collect();
    ensure!(lines.len() == 6, "monkey block needs six lines: {block}");

    let items = parse_suffix_numbers(lines[1], "Starting items:")?.into();
    let operation = parse_operation(lines[2])?;
    let divisor = parse_suffix_number(lines[3], "Test: divisible by")?;
    let if_true = parse_suffix_number(lines[4], "If true: throw to monkey")? as usize;
    let if_false = parse_suffix_number(lines[5], "If false: throw to monkey")? as usize;

    Ok(Monkey {
        items,
        operation,
        divisor,
        if_true,
        if_false,
        inspections: 0,
    })
}

fn parse(input: &str) -> Result<Vec<Monkey>> {
    let monkeys: Result<Vec<Monkey>> = input
        .split("\n\n")
        .map(str::trim)
        .filter(|block| !block.is_empty())
        .map(parse_monkey)
        .collect();
    let monkeys = monkeys?;
    ensure!(
        !monkeys.is_empty(),
        "input must contain at least one monkey"
    );
    for monkey in &monkeys {
        ensure!(monkey.divisor != 0, "monkey divisor must not be zero");
        ensure!(
            monkey.if_true < monkeys.len() && monkey.if_false < monkeys.len(),
            "throw target out of range"
        );
    }
    Ok(monkeys)
}

/// Runs the monkeys for `rounds`, reducing worry after each inspection via
/// `relief`, and returns the product of the two highest inspection counts.
fn run(monkeys: &mut [Monkey], rounds: usize, relief: impl Fn(u64) -> u64) -> u64 {
    for _ in 0..rounds {
        for i in 0..monkeys.len() {
            while let Some(item) = monkeys[i].items.pop_front() {
                monkeys[i].inspections += 1;
                let worry = relief(monkeys[i].operation.apply(item));
                let target = if worry.is_multiple_of(monkeys[i].divisor) {
                    monkeys[i].if_true
                } else {
                    monkeys[i].if_false
                };
                monkeys[target].items.push_back(worry);
            }
        }
    }
    let mut counts: Vec<u64> = monkeys.iter().map(|m| m.inspections).collect();
    counts.sort_unstable_by(|a, b| b.cmp(a));
    counts[0] * counts[1]
}

pub fn solve(input: &str) -> Result<(u64, u64)> {
    let mut part1_monkeys = parse(input)?;
    let part1 = run(&mut part1_monkeys, ROUNDS_PART1, |worry| {
        worry / RELIEF_DIVISOR
    });

    let mut part2_monkeys = parse(input)?;
    let modulus: u64 = part2_monkeys.iter().map(|m| m.divisor).product();
    let part2 = run(&mut part2_monkeys, ROUNDS_PART2, |worry| worry % modulus);

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day11-example.txt");

    #[test]
    fn part1_example() {
        let mut monkeys = parse(EXAMPLE).unwrap();
        let answer = run(&mut monkeys, ROUNDS_PART1, |worry| worry / RELIEF_DIVISOR);
        assert_eq!(answer, 10605);
    }

    #[test]
    fn part2_example() {
        let mut monkeys = parse(EXAMPLE).unwrap();
        let modulus: u64 = monkeys.iter().map(|m| m.divisor).product();
        let answer = run(&mut monkeys, ROUNDS_PART2, |worry| worry % modulus);
        assert_eq!(answer, 2713310158);
    }

    #[test]
    fn solve_matches_example() {
        let (part1, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 10605);
        assert_eq!(part2, 2713310158);
    }

    #[test]
    fn two_monkeys_swap_items_each_round() {
        let input = "Monkey 0:\n  Starting items: 1, 2\n  Operation: new = old + 1\n  Test: divisible by 2\n    If true: throw to monkey 1\n    If false: throw to monkey 1\n\nMonkey 1:\n  Starting items: 3\n  Operation: new = old + 1\n  Test: divisible by 2\n    If true: throw to monkey 0\n    If false: throw to monkey 0\n";
        let mut monkeys = parse(input).unwrap();
        let answer = run(&mut monkeys, 1, |worry| worry / RELIEF_DIVISOR);
        // Monkey 0 inspects 2 items, monkey 1 inspects the 3 items it now holds: 2 * 3 = 6.
        assert_eq!(answer, 6);
    }

    #[test]
    fn malformed_input_missing_field_is_error() {
        let input = "Monkey 0:\n  Starting items: 1\n  Operation: new = old + 1\n";
        assert!(parse(input).is_err());
    }

    #[test]
    fn malformed_input_bad_operator_is_error() {
        let input = "Monkey 0:\n  Starting items: 1\n  Operation: new = old ^ 1\n  Test: divisible by 2\n    If true: throw to monkey 0\n    If false: throw to monkey 0\n";
        assert!(parse(input).is_err());
    }

    #[test]
    fn malformed_input_throw_target_out_of_range_is_error() {
        let input = "Monkey 0:\n  Starting items: 1\n  Operation: new = old + 1\n  Test: divisible by 2\n    If true: throw to monkey 1\n    If false: throw to monkey 0\n";
        assert!(parse(input).is_err());
    }
}
