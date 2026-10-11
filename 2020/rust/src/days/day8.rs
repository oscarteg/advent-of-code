use anyhow::{bail, ensure, Context, Result};

#[derive(Clone, Copy, Debug, PartialEq, Eq)]
enum Instruction {
    Acc(i64),
    Jmp(i64),
    Nop(i64),
}

impl Instruction {
    fn parse(line: &str) -> Result<Self> {
        let (op, arg) = line
            .split_once(' ')
            .with_context(|| format!("malformed instruction: {line:?}"))?;
        let arg: i64 = arg
            .parse()
            .with_context(|| format!("malformed argument: {arg:?}"))?;
        match op {
            "acc" => Ok(Instruction::Acc(arg)),
            "jmp" => Ok(Instruction::Jmp(arg)),
            "nop" => Ok(Instruction::Nop(arg)),
            other => bail!("unknown operation: {other:?}"),
        }
    }

    /// Returns the jmp/nop-swapped variant, or `None` for `acc` instructions.
    fn swapped(self) -> Option<Self> {
        match self {
            Instruction::Jmp(arg) => Some(Instruction::Nop(arg)),
            Instruction::Nop(arg) => Some(Instruction::Jmp(arg)),
            Instruction::Acc(_) => None,
        }
    }
}

fn parse(input: &str) -> Result<Vec<Instruction>> {
    let program = input
        .lines()
        .filter(|line| !line.trim().is_empty())
        .map(Instruction::parse)
        .collect::<Result<Vec<_>>>()?;
    ensure!(!program.is_empty(), "program must not be empty");
    Ok(program)
}

enum Outcome {
    Terminated(i64),
    Looped(i64),
}

/// Executes the program, stopping either when it runs off the end
/// (terminated) or when it is about to re-run an already-visited
/// instruction (looped).
fn run(program: &[Instruction]) -> Outcome {
    let mut visited = vec![false; program.len()];
    let mut pc: i64 = 0;
    let mut acc: i64 = 0;
    while (0..program.len() as i64).contains(&pc) {
        let index = pc as usize;
        if visited[index] {
            return Outcome::Looped(acc);
        }
        visited[index] = true;
        match program[index] {
            Instruction::Acc(arg) => {
                acc += arg;
                pc += 1;
            }
            Instruction::Jmp(arg) => pc += arg,
            Instruction::Nop(_) => pc += 1,
        }
    }
    Outcome::Terminated(acc)
}

/// Fixes exactly one `jmp`/`nop` instruction so the program terminates,
/// returning the accumulator once it does.
fn fix_and_run(program: &[Instruction]) -> Result<i64> {
    for index in 0..program.len() {
        let Some(replacement) = program[index].swapped() else {
            continue;
        };
        let mut patched = program.to_vec();
        patched[index] = replacement;
        if let Outcome::Terminated(acc) = run(&patched) {
            return Ok(acc);
        }
    }
    bail!("no single jmp/nop swap makes the program terminate")
}

pub fn solve(input: &str) -> Result<(i64, i64)> {
    let program = parse(input)?;
    let part1 = match run(&program) {
        Outcome::Looped(acc) => acc,
        Outcome::Terminated(acc) => acc,
    };
    let part2 = fix_and_run(&program)?;
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/2020/day8-example.txt");
        assert_eq!(solve(input).unwrap(), (5, 8));
    }

    #[test]
    fn a_program_with_no_jmp_or_nop_cannot_be_fixed() {
        let input = "acc +1\nacc +2\nacc +3";
        assert!(solve(input).is_err());
    }

    #[test]
    fn negative_jumps_can_loop_back_past_the_start() {
        let input = "jmp +1\nacc +5\njmp -2";
        let outcome = run(&parse(input).unwrap());
        assert!(matches!(outcome, Outcome::Looped(5)));
    }

    #[test]
    fn unknown_operations_are_rejected() {
        assert!(solve("foo +1").is_err());
    }

    #[test]
    fn missing_arguments_are_rejected() {
        assert!(solve("acc").is_err());
    }

    #[test]
    fn non_numeric_arguments_are_rejected() {
        assert!(solve("acc +x").is_err());
    }

    #[test]
    fn empty_programs_are_rejected() {
        assert!(solve("").is_err());
    }
}
