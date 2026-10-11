use anyhow::{bail, ensure, Context, Result};

/// A pair nested inside four pairs has its regular numbers at this depth.
const EXPLODE_DEPTH: u32 = 5;

/// A regular number together with the number of pairs it is nested in.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
struct Element {
    value: u64,
    depth: u32,
}

/// A snailfish number flattened to its regular numbers in reading order; the
/// bracket structure follows from the depths alone.
#[derive(Clone, Debug, PartialEq, Eq)]
struct Snailfish {
    elements: Vec<Element>,
}

impl Snailfish {
    fn parse(line: &str) -> Result<Self> {
        let mut parser = Parser {
            bytes: line.as_bytes(),
            position: 0,
        };
        let mut elements = Vec::new();
        parser
            .pair(1, &mut elements)
            .with_context(|| format!("invalid snailfish number {line:?}"))?;
        ensure!(
            parser.position == parser.bytes.len(),
            "trailing input in {line:?}"
        );
        Ok(Self { elements })
    }

    /// Pairs both numbers and reduces the sum.
    fn add(&self, other: &Self) -> Self {
        let mut sum = Self {
            elements: self
                .elements
                .iter()
                .chain(&other.elements)
                .map(|element| Element {
                    value: element.value,
                    depth: element.depth + 1,
                })
                .collect(),
        };
        while sum.explode() || sum.split() {}
        sum
    }

    /// Explodes the leftmost pair nested inside four pairs, which in this
    /// representation is the leftmost couple of neighbours sharing a depth of
    /// five or more.
    fn explode(&mut self) -> bool {
        let Some(index) = self
            .elements
            .windows(2)
            .position(|pair| pair[0].depth >= EXPLODE_DEPTH && pair[0].depth == pair[1].depth)
        else {
            return false;
        };
        let (left, right) = (self.elements[index], self.elements[index + 1]);
        if index > 0 {
            self.elements[index - 1].value += left.value;
        }
        if let Some(neighbour) = self.elements.get_mut(index + 2) {
            neighbour.value += right.value;
        }
        self.elements[index] = Element {
            value: 0,
            depth: left.depth - 1,
        };
        self.elements.remove(index + 1);
        true
    }

    /// Splits the leftmost regular number of ten or more, rounding the left
    /// half down and the right half up.
    fn split(&mut self) -> bool {
        let Some(index) = self.elements.iter().position(|element| element.value >= 10) else {
            return false;
        };
        let Element { value, depth } = self.elements[index];
        self.elements[index] = Element {
            value: value / 2,
            depth: depth + 1,
        };
        self.elements.insert(
            index + 1,
            Element {
                value: value.div_ceil(2),
                depth: depth + 1,
            },
        );
        true
    }

    /// Folds neighbours of equal depth into 3 * left + 2 * right until the
    /// whole number collapses into a single value.
    fn magnitude(&self) -> Result<u64> {
        let mut stack: Vec<Element> = Vec::with_capacity(self.elements.len());
        for &element in &self.elements {
            let mut folded = element;
            while let Some(&top) = stack.last() {
                if top.depth != folded.depth {
                    break;
                }
                stack.truncate(stack.len() - 1);
                folded = Element {
                    value: 3 * top.value + 2 * folded.value,
                    depth: folded.depth - 1,
                };
            }
            stack.push(folded);
        }
        match stack.as_slice() {
            [Element { value, depth: 0 }] => Ok(*value),
            _ => bail!("a snailfish number must be a single pair"),
        }
    }
}

struct Parser<'a> {
    bytes: &'a [u8],
    position: usize,
}

impl Parser<'_> {
    fn pair(&mut self, depth: u32, elements: &mut Vec<Element>) -> Result<()> {
        self.expect(b'[')?;
        self.element(depth, elements)?;
        self.expect(b',')?;
        self.element(depth, elements)?;
        self.expect(b']')
    }

    fn element(&mut self, depth: u32, elements: &mut Vec<Element>) -> Result<()> {
        if self.peek() == Some(b'[') {
            return self.pair(depth + 1, elements);
        }
        let start = self.position;
        while self.peek().is_some_and(|byte| byte.is_ascii_digit()) {
            self.position += 1;
        }
        ensure!(
            self.position > start,
            "expected a pair or a regular number at offset {start}"
        );
        let digits = std::str::from_utf8(&self.bytes[start..self.position])?;
        elements.push(Element {
            value: digits
                .parse()
                .with_context(|| format!("regular number {digits:?}"))?,
            depth,
        });
        Ok(())
    }

    fn peek(&self) -> Option<u8> {
        self.bytes.get(self.position).copied()
    }

    fn expect(&mut self, byte: u8) -> Result<()> {
        ensure!(
            self.peek() == Some(byte),
            "expected {:?} at offset {}",
            char::from(byte),
            self.position
        );
        self.position += 1;
        Ok(())
    }
}

fn parse_numbers(input: &str) -> Result<Vec<Snailfish>> {
    let numbers = input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(Snailfish::parse)
        .collect::<Result<Vec<_>>>()?;
    ensure!(!numbers.is_empty(), "no snailfish numbers in the input");
    Ok(numbers)
}

/// Returns the magnitude of the homework sum and the largest magnitude of any
/// sum of two different numbers.
pub fn solve(input: &str) -> Result<(u64, u64)> {
    let numbers = parse_numbers(input)?;
    ensure!(
        numbers.len() >= 2,
        "part two needs at least two snailfish numbers"
    );
    let homework = numbers[1..]
        .iter()
        .fold(numbers[0].clone(), |sum, number| sum.add(number))
        .magnitude()?;
    let mut best = 0;
    for (index, left) in numbers.iter().enumerate() {
        for right in &numbers[index + 1..] {
            for sum in [left.add(right), right.add(left)] {
                best = best.max(sum.magnitude()?);
            }
        }
    }
    Ok((homework, best))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day18-example.txt");

    const LARGER_EXAMPLE: &str = "\
[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]
[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]
[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]
[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]
[7,[5,[[3,8],[1,4]]]]
[[2,[2,2]],[8,[8,1]]]
[2,9]
[1,[[[9,3],9],[[9,0],[0,7]]]]
[[[5,[7,4]],7],1]
[[[[4,2],2],6],[8,7]]";

    #[test]
    fn official_example() -> Result<()> {
        assert_eq!(solve(EXAMPLE)?, (4140, 3993));
        Ok(())
    }

    #[test]
    fn official_magnitudes() -> Result<()> {
        for (line, magnitude) in [
            ("[[1,2],[[3,4],5]]", 143),
            ("[[[[0,7],4],[[7,8],[6,0]]],[8,1]]", 1384),
            ("[[[[1,1],[2,2]],[3,3]],[4,4]]", 445),
            ("[[[[3,0],[5,3]],[4,4]],[5,5]]", 791),
            ("[[[[5,0],[7,4]],[5,5]],[6,6]]", 1137),
            (
                "[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]",
                3488,
            ),
        ] {
            assert_eq!(Snailfish::parse(line)?.magnitude()?, magnitude, "{line}");
        }
        Ok(())
    }

    #[test]
    fn official_sums() -> Result<()> {
        assert_eq!(solve("[1,1]\n[2,2]\n[3,3]\n[4,4]")?.0, 445);
        assert_eq!(solve("[1,1]\n[2,2]\n[3,3]\n[4,4]\n[5,5]")?.0, 791);
        assert_eq!(solve("[1,1]\n[2,2]\n[3,3]\n[4,4]\n[5,5]\n[6,6]")?.0, 1137);
        assert_eq!(solve(LARGER_EXAMPLE)?.0, 3488);
        Ok(())
    }

    #[test]
    fn addition_reduces_the_sum() -> Result<()> {
        let sum =
            Snailfish::parse("[[[[4,3],4],4],[7,[[8,4],9]]]")?.add(&Snailfish::parse("[1,1]")?);
        assert_eq!(sum, Snailfish::parse("[[[[0,7],4],[[7,8],[6,0]]],[8,1]]")?);
        Ok(())
    }

    #[test]
    fn single_explodes_feed_their_neighbours() -> Result<()> {
        for (before, after) in [
            ("[[[[[9,8],1],2],3],4]", "[[[[0,9],2],3],4]"),
            ("[7,[6,[5,[4,[3,2]]]]]", "[7,[6,[5,[7,0]]]]"),
            ("[[6,[5,[4,[3,2]]]],1]", "[[6,[5,[7,0]]],3]"),
            (
                "[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]",
                "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]",
            ),
            (
                "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]",
                "[[3,[2,[8,0]]],[9,[5,[7,0]]]]",
            ),
        ] {
            let mut number = Snailfish::parse(before)?;
            assert!(number.explode(), "{before} should explode");
            assert_eq!(number, Snailfish::parse(after)?, "{before}");
        }
        Ok(())
    }

    #[test]
    fn splits_round_down_then_up() -> Result<()> {
        let mut number = Snailfish::parse("[11,0]")?;
        assert!(number.split());
        assert_eq!(number, Snailfish::parse("[[5,6],0]")?);
        assert!(!number.split(), "nothing left to split");
        Ok(())
    }

    #[test]
    fn addition_order_matters_for_part_two() -> Result<()> {
        // [[9,1],[1,9]] has magnitude 129, the reverse order only 121.
        assert_eq!(solve("[9,1]\n[1,9]")?, (129, 129));
        Ok(())
    }

    #[test]
    fn windows_newlines_are_accepted() -> Result<()> {
        assert_eq!(solve("[9,1]\r\n[1,9]\r\n")?, (129, 129));
        Ok(())
    }

    #[test]
    fn a_single_number_has_nothing_to_pair_with() {
        assert!(solve("[1,2]").is_err());
    }

    #[test]
    fn malformed_input_is_rejected() {
        for input in [
            "",
            "   \n\n",
            "5\n[1,2]",
            "[1,2\n[3,4]",
            "[1,2]]\n[3,4]",
            "[1,2]x\n[3,4]",
            "[1,2,3]\n[3,4]",
            "[]\n[1,2]",
            "[1,2]\n[3,[4]]",
            "[1,-2]\n[3,4]",
            "[99999999999999999999,2]\n[3,4]",
        ] {
            assert!(solve(input).is_err(), "{input:?} should be rejected");
        }
    }
}
