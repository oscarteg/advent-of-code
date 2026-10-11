use anyhow::{bail, ensure, Context, Result};

/// A bit cursor over the hexadecimal digits of a transmission.
struct Bits {
    nibbles: Vec<u8>,
    position: usize,
}

impl Bits {
    fn parse(input: &str) -> Result<Self> {
        let digits = input.trim();
        ensure!(!digits.is_empty(), "the transmission is empty");
        let nibbles = digits
            .chars()
            .map(|digit| {
                let value = digit
                    .to_digit(16)
                    .with_context(|| format!("{digit:?} is not a hexadecimal digit"))?;
                Ok(u8::try_from(value)?)
            })
            .collect::<Result<_>>()?;
        Ok(Self {
            nibbles,
            position: 0,
        })
    }

    fn len(&self) -> usize {
        self.nibbles.len() * 4
    }

    fn bit(&self, index: usize) -> u64 {
        u64::from((self.nibbles[index / 4] >> (3 - index % 4)) & 1)
    }

    fn take(&mut self, count: usize) -> Result<u64> {
        let end = self.position + count;
        ensure!(end <= self.len(), "the transmission ends inside a packet");
        let value = (self.position..end).fold(0, |value, index| (value << 1) | self.bit(index));
        self.position = end;
        Ok(value)
    }

    /// Only the zero bits that pad the final nibble may follow the outermost packet.
    fn rest_is_padding(&self) -> bool {
        (self.position..self.len()).all(|index| self.bit(index) == 0)
    }
}

#[derive(Clone, Copy)]
enum Operator {
    Sum,
    Product,
    Minimum,
    Maximum,
    GreaterThan,
    LessThan,
    EqualTo,
}

impl Operator {
    /// Type ID 4 marks a literal value; every other type ID is an operator.
    fn from_type_id(type_id: u64) -> Option<Self> {
        match type_id {
            0 => Some(Self::Sum),
            1 => Some(Self::Product),
            2 => Some(Self::Minimum),
            3 => Some(Self::Maximum),
            5 => Some(Self::GreaterThan),
            6 => Some(Self::LessThan),
            7 => Some(Self::EqualTo),
            _ => None,
        }
    }

    fn apply(self, operands: &[u64]) -> Result<u64> {
        let pair = || -> Result<(u64, u64)> {
            match operands {
                [left, right] => Ok((*left, *right)),
                _ => bail!("a comparison needs exactly two sub-packets"),
            }
        };
        match self {
            Self::Sum => operands.iter().try_fold(0u64, |total, &value| {
                total.checked_add(value).context("a sum overflows")
            }),
            Self::Product => operands.iter().try_fold(1u64, |total, &value| {
                total.checked_mul(value).context("a product overflows")
            }),
            Self::Minimum => operands
                .iter()
                .copied()
                .min()
                .context("a minimum needs at least one sub-packet"),
            Self::Maximum => operands
                .iter()
                .copied()
                .max()
                .context("a maximum needs at least one sub-packet"),
            Self::GreaterThan => {
                let (left, right) = pair()?;
                Ok(u64::from(left > right))
            }
            Self::LessThan => {
                let (left, right) = pair()?;
                Ok(u64::from(left < right))
            }
            Self::EqualTo => {
                let (left, right) = pair()?;
                Ok(u64::from(left == right))
            }
        }
    }
}

enum Body {
    Literal(u64),
    Operator {
        operator: Operator,
        operands: Vec<Packet>,
    },
}

/// A decoded packet: either a literal value or an operator over its sub-packets.
struct Packet {
    version: u64,
    body: Body,
}

impl Packet {
    fn parse(bits: &mut Bits) -> Result<Self> {
        let version = bits.take(3)?;
        let type_id = bits.take(3)?;
        let body = match Operator::from_type_id(type_id) {
            Some(operator) => Body::Operator {
                operator,
                operands: Self::parse_operands(bits)?,
            },
            None => Body::Literal(Self::parse_literal(bits)?),
        };
        Ok(Self { version, body })
    }

    /// A literal is a chain of five-bit groups, each but the last flagged with a leading bit.
    fn parse_literal(bits: &mut Bits) -> Result<u64> {
        let mut value = 0u64;
        loop {
            let group = bits.take(5)?;
            value = value.checked_mul(16).context("a literal value overflows")? | (group & 0b1111);
            if group & 0b1_0000 == 0 {
                return Ok(value);
            }
        }
    }

    fn parse_operands(bits: &mut Bits) -> Result<Vec<Self>> {
        if bits.take(1)? == 1 {
            let count = usize::try_from(bits.take(11)?)?;
            return (0..count).map(|_| Self::parse(bits)).collect();
        }
        let length = usize::try_from(bits.take(15)?)?;
        let end = bits.position + length;
        ensure!(end <= bits.len(), "sub-packets run past the transmission");
        let mut operands = Vec::new();
        while bits.position < end {
            operands.push(Self::parse(bits)?);
        }
        ensure!(
            bits.position == end,
            "sub-packets overrun their declared length"
        );
        Ok(operands)
    }

    fn version_sum(&self) -> u64 {
        match &self.body {
            Body::Literal(_) => self.version,
            Body::Operator { operands, .. } => {
                self.version + operands.iter().map(Self::version_sum).sum::<u64>()
            }
        }
    }

    fn value(&self) -> Result<u64> {
        match &self.body {
            Body::Literal(value) => Ok(*value),
            Body::Operator { operator, operands } => {
                let values = operands
                    .iter()
                    .map(Self::value)
                    .collect::<Result<Vec<_>>>()?;
                operator.apply(&values)
            }
        }
    }
}

/// Returns the sum of every packet version and the value of the outermost packet.
pub fn solve(input: &str) -> Result<(u64, u64)> {
    let mut bits = Bits::parse(input)?;
    let outermost = Packet::parse(&mut bits)?;
    ensure!(
        bits.rest_is_padding(),
        "the transmission has data after the outermost packet"
    );
    Ok((outermost.version_sum(), outermost.value()?))
}

#[cfg(test)]
mod tests {
    use super::*;

    fn version_sum(transmission: &str) -> u64 {
        solve(transmission).unwrap().0
    }

    fn value(transmission: &str) -> u64 {
        solve(transmission).unwrap().1
    }

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day16-example.txt");
        assert_eq!(solve(input).unwrap(), (6, 2021));
    }

    #[test]
    fn official_version_sums() {
        for (transmission, expected) in [
            ("8A004A801A8002F478", 16),
            ("620080001611562C8802118E34", 12),
            ("C0015000016115A2E0802F182340", 23),
            ("A0016C880162017C3686B18A3D4780", 31),
        ] {
            assert_eq!(version_sum(transmission), expected, "{transmission}");
        }
    }

    #[test]
    fn official_values() {
        for (transmission, expected) in [
            ("C200B40A82", 3),
            ("04005AC33890", 54),
            ("880086C3E88112", 7),
            ("CE00C43D881120", 9),
            ("D8005AC2A8F0", 1),
            ("F600BC2D8F", 0),
            ("9C005AC2F8F0", 0),
            ("9C0141080250320F1802104A08", 1),
        ] {
            assert_eq!(value(transmission), expected, "{transmission}");
        }
    }

    #[test]
    fn both_length_types_are_decoded() {
        // 38006F45291200 is "10 is less than 20", EE00D40C823060 the maximum of 1, 2 and 3.
        assert_eq!(value("38006F45291200"), 1);
        assert_eq!(value("EE00D40C823060"), 3);
    }

    #[test]
    fn lowercase_digits_and_surrounding_whitespace_are_accepted() {
        assert_eq!(solve("  d2fe28\n").unwrap(), (6, 2021));
    }

    #[test]
    fn malformed_transmissions_are_rejected() {
        for transmission in ["", "\n", "D2FG28", "D2", "38006F452912", "D2FE28F"] {
            assert!(solve(transmission).is_err(), "{transmission:?}");
        }
    }
}
