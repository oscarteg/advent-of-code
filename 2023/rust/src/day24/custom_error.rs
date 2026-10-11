use miette::Diagnostic;
use thiserror::Error;

#[derive(Error, Diagnostic, Debug)]
pub enum AocError {
    #[error("invalid hailstones: {0}")]
    #[diagnostic(code(aoc::invalid_input))]
    InvalidInput(&'static str),

    #[error("exact arithmetic overflowed")]
    #[diagnostic(code(aoc::overflow))]
    Overflow,

    #[error(transparent)]
    #[diagnostic(code(aoc::parse_int_error))]
    ParseIntError(#[from] std::num::ParseIntError),
}
