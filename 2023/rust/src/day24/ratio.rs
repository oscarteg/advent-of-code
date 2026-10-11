use crate::day24::custom_error::AocError;

/// Greatest common divisor by Euclid's algorithm. The sign follows the first
/// argument, which callers either normalise away or do not care about.
pub fn gcd(first: i128, second: i128) -> i128 {
    let (mut first, mut second) = (first, second);
    while second != 0 {
        (first, second) = (second, first % second);
    }
    first
}

/// An exact rational number in lowest terms with a positive denominator. Every
/// operation is checked, so an overflow becomes an error instead of a silently
/// wrong answer.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct Ratio {
    numerator: i128,
    denominator: i128,
}

impl Ratio {
    pub fn new(numerator: i128, denominator: i128) -> Result<Self, AocError> {
        if denominator == 0 {
            return Err(AocError::InvalidInput("division by zero"));
        }
        let divisor = gcd(numerator, denominator);
        let reduce = |value: i128| value.checked_div(divisor).ok_or(AocError::Overflow);
        let (numerator, denominator) = (reduce(numerator)?, reduce(denominator)?);
        if denominator < 0 {
            return Ok(Self {
                numerator: numerator.checked_neg().ok_or(AocError::Overflow)?,
                denominator: denominator.checked_neg().ok_or(AocError::Overflow)?,
            });
        }
        Ok(Self {
            numerator,
            denominator,
        })
    }

    pub const fn from_integer(value: i128) -> Self {
        Self {
            numerator: value,
            denominator: 1,
        }
    }

    pub fn checked_add(self, other: Self) -> Result<Self, AocError> {
        self.combine(other, i128::checked_add)
    }

    pub fn checked_sub(self, other: Self) -> Result<Self, AocError> {
        self.combine(other, i128::checked_sub)
    }

    pub fn checked_mul(self, other: Self) -> Result<Self, AocError> {
        let numerator = self
            .numerator
            .checked_mul(other.numerator)
            .ok_or(AocError::Overflow)?;
        let denominator = self
            .denominator
            .checked_mul(other.denominator)
            .ok_or(AocError::Overflow)?;
        Self::new(numerator, denominator)
    }

    pub fn checked_div(self, other: Self) -> Result<Self, AocError> {
        let numerator = self
            .numerator
            .checked_mul(other.denominator)
            .ok_or(AocError::Overflow)?;
        let denominator = self
            .denominator
            .checked_mul(other.numerator)
            .ok_or(AocError::Overflow)?;
        Self::new(numerator, denominator)
    }

    pub const fn is_zero(self) -> bool {
        self.numerator == 0
    }

    /// The value as an integer, or `None` when it is a true fraction.
    pub const fn to_integer(self) -> Option<i128> {
        if self.denominator == 1 {
            Some(self.numerator)
        } else {
            None
        }
    }

    /// Brings both values over the common denominator and joins the numerators.
    fn combine(self, other: Self, join: fn(i128, i128) -> Option<i128>) -> Result<Self, AocError> {
        let numerator = self
            .numerator
            .checked_mul(other.denominator)
            .zip(other.numerator.checked_mul(self.denominator))
            .and_then(|(left, right)| join(left, right))
            .ok_or(AocError::Overflow)?;
        let denominator = self
            .denominator
            .checked_mul(other.denominator)
            .ok_or(AocError::Overflow)?;
        Self::new(numerator, denominator)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case(4, 6, 2, 3)]
    #[case(2, -4, -1, 2)]
    #[case(0, -7, 0, 1)]
    #[case(-9, -3, 3, 1)]
    fn fractions_are_reduced(
        #[case] numerator: i128,
        #[case] denominator: i128,
        #[case] expected_numerator: i128,
        #[case] expected_denominator: i128,
    ) {
        assert_eq!(
            Ratio::new(expected_numerator, expected_denominator).unwrap(),
            Ratio::new(numerator, denominator).unwrap()
        );
    }

    #[test]
    fn arithmetic_stays_exact() -> Result<(), AocError> {
        let third = Ratio::new(1, 3)?;
        let sixth = Ratio::new(1, 6)?;
        assert_eq!(Ratio::new(1, 2)?, third.checked_add(sixth)?);
        assert_eq!(sixth, third.checked_sub(sixth)?);
        assert_eq!(
            Some(1),
            third.checked_mul(Ratio::from_integer(3))?.to_integer()
        );
        assert_eq!(Some(2), third.checked_div(sixth)?.to_integer());
        assert_eq!(None, third.to_integer());
        assert!(Ratio::from_integer(0).is_zero());
        Ok(())
    }

    #[test]
    fn invalid_arithmetic_returns_errors() {
        assert!(Ratio::new(1, 0).is_err());
        assert!(Ratio::from_integer(1)
            .checked_div(Ratio::from_integer(0))
            .is_err());
        assert!(Ratio::from_integer(i128::MAX)
            .checked_mul(Ratio::from_integer(2))
            .is_err());
        assert!(Ratio::from_integer(i128::MAX)
            .checked_add(Ratio::from_integer(1))
            .is_err());
    }

    #[rstest]
    #[case(12, 8, 4)]
    #[case(0, 5, 5)]
    #[case(7, 0, 7)]
    #[case(17, 5, 1)]
    fn greatest_common_divisors(#[case] first: i128, #[case] second: i128, #[case] expected: i128) {
        assert_eq!(expected, gcd(first, second));
    }
}
