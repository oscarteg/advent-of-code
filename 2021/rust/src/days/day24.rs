use anyhow::{bail, ensure, Context, Result};

#[derive(Clone, Copy, PartialEq, Eq, Debug)]
enum Register {
    W,
    X,
    Y,
    Z,
}

impl Register {
    fn parse(token: &str) -> Result<Self> {
        match token {
            "w" => Ok(Self::W),
            "x" => Ok(Self::X),
            "y" => Ok(Self::Y),
            "z" => Ok(Self::Z),
            other => bail!("unknown register `{other}`"),
        }
    }

    fn index(self) -> usize {
        match self {
            Self::W => 0,
            Self::X => 1,
            Self::Y => 2,
            Self::Z => 3,
        }
    }
}

#[derive(Clone, Copy, PartialEq, Eq, Debug)]
enum Operand {
    Register(Register),
    Literal(i64),
}

impl Operand {
    fn parse(token: &str) -> Result<Self> {
        match token {
            "w" | "x" | "y" | "z" => Ok(Self::Register(Register::parse(token)?)),
            literal => Ok(Self::Literal(
                literal
                    .parse()
                    .with_context(|| format!("operand `{literal}`"))?,
            )),
        }
    }
}

/// The five binary operations; `inp` is [`Instruction::Input`] instead.
#[derive(Clone, Copy, PartialEq, Eq, Debug)]
enum Operation {
    Add,
    Mul,
    Div,
    Mod,
    Eql,
}

#[derive(Clone, Copy, PartialEq, Eq, Debug)]
enum Instruction {
    Input(Register),
    Binary {
        operation: Operation,
        target: Register,
        operand: Operand,
    },
}

impl Instruction {
    fn parse(line: &str) -> Result<Self> {
        let mut tokens = line.split_whitespace();
        let opcode = tokens.next().context("empty instruction")?;
        let target = Register::parse(
            tokens
                .next()
                .with_context(|| format!("`{opcode}` without a target register"))?,
        )?;
        let instruction = if opcode == "inp" {
            Self::Input(target)
        } else {
            let operand = Operand::parse(
                tokens
                    .next()
                    .with_context(|| format!("`{opcode}` without a second operand"))?,
            )?;
            let operation = match opcode {
                "add" => Operation::Add,
                "mul" => Operation::Mul,
                "div" => Operation::Div,
                "mod" => Operation::Mod,
                "eql" => Operation::Eql,
                other => bail!("unknown instruction `{other}`"),
            };
            Self::Binary {
                operation,
                target,
                operand,
            }
        };
        ensure!(tokens.next().is_none(), "trailing tokens in `{line}`");
        Ok(instruction)
    }
}

fn parse_program(input: &str) -> Result<Vec<Instruction>> {
    input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(Instruction::parse)
        .collect()
}

/// Feeds `inputs` to the `inp` instructions in order and returns `[w, x, y, z]`.
fn run(program: &[Instruction], inputs: &[i64]) -> Result<[i64; 4]> {
    let mut registers = [0i64; 4];
    let mut inputs = inputs.iter();
    for &instruction in program {
        let (operation, target, operand) = match instruction {
            Instruction::Input(target) => {
                registers[target.index()] = *inputs
                    .next()
                    .context("the program reads more inputs than supplied")?;
                continue;
            }
            Instruction::Binary {
                operation,
                target,
                operand,
            } => (operation, target, operand),
        };
        let left = registers[target.index()];
        let right = match operand {
            Operand::Register(register) => registers[register.index()],
            Operand::Literal(literal) => literal,
        };
        registers[target.index()] = match operation {
            Operation::Add => left.checked_add(right).context("ALU register overflow")?,
            Operation::Mul => left.checked_mul(right).context("ALU register overflow")?,
            Operation::Div => {
                ensure!(right != 0, "`div {left} {right}` crashes the ALU");
                left.checked_div(right).context("ALU register overflow")?
            }
            Operation::Mod => {
                ensure!(
                    left >= 0 && right > 0,
                    "`mod {left} {right}` crashes the ALU"
                );
                left % right
            }
            Operation::Eql => i64::from(left == right),
        };
    }
    Ok(registers)
}

const BLOCK_LEN: usize = 18;

/// A MONAD block reads a digit, divides `z` by `div`, compares `z % 26 + add_x`
/// with that digit and, when they differ, pushes `digit + add_y` onto `z`.
#[derive(Clone, Copy, PartialEq, Eq, Debug)]
struct Block {
    div: i64,
    add_x: i64,
    add_y: i64,
}

impl Block {
    /// The block as MONAD spells it out, used to reject other program shapes.
    fn instructions(self) -> [Instruction; BLOCK_LEN] {
        use Operation::{Add, Div, Eql, Mod, Mul};
        use Register::{W, X, Y, Z};
        let register = |operation, target, source| Instruction::Binary {
            operation,
            target,
            operand: Operand::Register(source),
        };
        let literal = |operation, target, literal| Instruction::Binary {
            operation,
            target,
            operand: Operand::Literal(literal),
        };
        [
            Instruction::Input(W),
            literal(Mul, X, 0),
            register(Add, X, Z),
            literal(Mod, X, 26),
            literal(Div, Z, self.div),
            literal(Add, X, self.add_x),
            register(Eql, X, W),
            literal(Eql, X, 0),
            literal(Mul, Y, 0),
            literal(Add, Y, 25),
            register(Mul, Y, X),
            literal(Add, Y, 1),
            register(Mul, Z, Y),
            literal(Mul, Y, 0),
            register(Add, Y, W),
            literal(Add, Y, self.add_y),
            register(Mul, Y, X),
            register(Add, Z, Y),
        ]
    }

    fn parse(instructions: &[Instruction]) -> Result<Self> {
        ensure!(
            instructions.len() == BLOCK_LEN,
            "a MONAD block has {BLOCK_LEN} instructions, found {}",
            instructions.len()
        );
        let block = Self {
            div: literal_operand(instructions[4])?,
            add_x: literal_operand(instructions[5])?,
            add_y: literal_operand(instructions[15])?,
        };
        ensure!(
            matches!(block.div, 1 | 26),
            "a MONAD block divides z by 1 or 26, found {}",
            block.div
        );
        ensure!(
            block.instructions().as_slice() == instructions,
            "the block does not follow the MONAD template"
        );
        Ok(block)
    }
}

fn literal_operand(instruction: Instruction) -> Result<i64> {
    match instruction {
        Instruction::Binary {
            operand: Operand::Literal(literal),
            ..
        } => Ok(literal),
        other => bail!("expected a literal operand, found {other:?}"),
    }
}

/// A `div 26` block pops the value an earlier `div 1` block pushed, which forces
/// `digits[pop] == digits[push] + offset`.
#[derive(Clone, Copy, Debug)]
struct Constraint {
    push: usize,
    pop: usize,
    offset: i64,
}

#[derive(Clone, Copy)]
enum Extreme {
    Largest,
    Smallest,
}

struct Monad {
    program: Vec<Instruction>,
    blocks: Vec<Block>,
}

impl Monad {
    fn parse(input: &str) -> Result<Self> {
        let program = parse_program(input)?;
        ensure!(!program.is_empty(), "the program is empty");
        ensure!(
            program.len() % BLOCK_LEN == 0,
            "MONAD repeats blocks of {BLOCK_LEN} instructions, found {} in total",
            program.len()
        );
        let blocks = program
            .chunks(BLOCK_LEN)
            .map(Block::parse)
            .collect::<Result<Vec<_>>>()?;
        Ok(Self { program, blocks })
    }

    fn constraints(&self) -> Result<Vec<Constraint>> {
        let mut pushed: Vec<(usize, i64)> = Vec::new();
        let mut constraints = Vec::with_capacity(self.blocks.len() / 2);
        for (index, block) in self.blocks.iter().enumerate() {
            if block.div == 1 {
                ensure!(
                    block.add_x > 9,
                    "block {index} pushes, so `add x {}` must exceed any digit",
                    block.add_x
                );
                pushed.push((index, block.add_y));
                continue;
            }
            let (push, add_y) = pushed
                .pop()
                .with_context(|| format!("block {index} pops while z is empty"))?;
            constraints.push(Constraint {
                push,
                pop: index,
                offset: add_y + block.add_x,
            });
        }
        ensure!(
            pushed.is_empty(),
            "{} pushed values are never popped, so z cannot end at 0",
            pushed.len()
        );
        Ok(constraints)
    }

    /// Picks the extreme digit for every constrained pair, then checks the
    /// resulting model number with the interpreter.
    fn model_number(&self, extreme: Extreme) -> Result<u64> {
        let mut digits = vec![0i64; self.blocks.len()];
        for Constraint { push, pop, offset } in self.constraints()? {
            ensure!(
                offset.abs() <= 8,
                "digits {push} and {pop} must differ by {offset}, which digits 1..9 cannot do"
            );
            let (push_digit, pop_digit) = match extreme {
                Extreme::Largest => {
                    let pop_digit = 9.min(9 + offset);
                    (pop_digit - offset, pop_digit)
                }
                Extreme::Smallest => {
                    let push_digit = 1.max(1 - offset);
                    (push_digit, push_digit + offset)
                }
            };
            digits[push] = push_digit;
            digits[pop] = pop_digit;
        }
        let z = run(&self.program, &digits)?[Register::Z.index()];
        ensure!(z == 0, "the chosen model number leaves z = {z}");
        digits.iter().try_fold(0u64, |number, &digit| {
            Ok(number * 10 + u64::try_from(digit)?)
        })
    }
}

pub fn solve(input: &str) -> Result<(u64, u64)> {
    let monad = Monad::parse(input)?;
    Ok((
        monad.model_number(Extreme::Largest)?,
        monad.model_number(Extreme::Smallest)?,
    ))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day24-example.txt");

    fn program_text(blocks: &[(i64, i64, i64)]) -> String {
        blocks
            .iter()
            .map(|(div, add_x, add_y)| {
                format!(
                    "inp w\nmul x 0\nadd x z\nmod x 26\ndiv z {div}\nadd x {add_x}\neql x w\n\
                     eql x 0\nmul y 0\nadd y 25\nmul y x\nadd y 1\nmul z y\nmul y 0\nadd y w\n\
                     add y {add_y}\nmul y x\nadd z y\n"
                )
            })
            .collect()
    }

    fn digits_of(number: u64) -> Result<Vec<i64>> {
        number
            .to_string()
            .chars()
            .map(|digit| digit.to_digit(10).map(i64::from).context("not a digit"))
            .collect()
    }

    #[test]
    fn negates_its_input() -> Result<()> {
        let program = parse_program("inp x\nmul x -1\n")?;
        assert_eq!(run(&program, &[7])?[Register::X.index()], -7);
        Ok(())
    }

    #[test]
    fn detects_three_times_larger_inputs() -> Result<()> {
        let program = parse_program("inp z\ninp x\nmul z 3\neql z x\n")?;
        assert_eq!(run(&program, &[3, 9])?[Register::Z.index()], 1);
        assert_eq!(run(&program, &[3, 8])?[Register::Z.index()], 0);
        Ok(())
    }

    #[test]
    fn splits_its_input_into_bits() -> Result<()> {
        let program = parse_program(
            "inp w\nadd z w\nmod z 2\ndiv w 2\nadd y w\nmod y 2\ndiv w 2\nadd x w\nmod x 2\n\
             div w 2\nmod w 2\n",
        )?;
        assert_eq!(run(&program, &[11])?, [1, 0, 1, 1]);
        Ok(())
    }

    #[test]
    fn solves_the_example_monad() -> Result<()> {
        assert_eq!(solve(EXAMPLE)?, (99_999_596_891_916, 49_411_151_121_911));
        Ok(())
    }

    #[test]
    fn agrees_with_brute_force_on_a_four_block_monad() -> Result<()> {
        let input = program_text(&[(1, 13, 6), (1, 11, 9), (26, -4, 7), (26, -8, 5)]);
        let monad = Monad::parse(&input)?;
        let mut valid = Vec::new();
        for number in 1_111..=9_999u64 {
            let digits = digits_of(number)?;
            if !digits.contains(&0) && run(&monad.program, &digits)?[Register::Z.index()] == 0 {
                valid.push(number);
            }
        }
        let largest = *valid.last().context("no valid model number")?;
        let smallest = *valid.first().context("no valid model number")?;
        assert_eq!(solve(&input)?, (largest, smallest));
        Ok(())
    }

    #[test]
    fn a_pair_eight_apart_forces_both_digits() -> Result<()> {
        assert_eq!(solve(&program_text(&[(1, 13, 0), (26, 8, 0)]))?, (19, 19));
        Ok(())
    }

    #[test]
    fn rejects_unsatisfiable_digit_pairs() {
        assert!(solve(&program_text(&[(1, 13, 0), (26, 9, 0)])).is_err());
    }

    #[test]
    fn rejects_programs_that_are_not_monad_shaped() {
        let monad = program_text(&[(1, 13, 6), (26, -4, 7)]);
        assert!(solve(&monad.replace("mod x 26", "mod x 25")).is_err());
        assert!(solve(&program_text(&[(1, 13, 6), (1, 11, 9)])).is_err());
        assert!(solve(&program_text(&[(26, -4, 7), (1, 13, 6)])).is_err());
    }

    #[test]
    fn rejects_malformed_input() {
        for input in [
            "",
            "inp q\n",
            "nop x 1\n",
            "inp\n",
            "add x\n",
            "mul x y z\n",
        ] {
            assert!(solve(input).is_err(), "expected an error for {input:?}");
        }
    }

    #[test]
    fn rejects_instructions_that_crash_the_alu() -> Result<()> {
        assert!(run(&parse_program("inp x\ndiv x 0\n")?, &[1]).is_err());
        assert!(run(&parse_program("inp x\nmul x -1\nmod x 2\n")?, &[1]).is_err());
        assert!(run(&parse_program("inp x\n")?, &[]).is_err());
        Ok(())
    }
}
