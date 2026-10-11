use miette::Diagnostic;
use thiserror::Error;

#[derive(Error, Diagnostic, Debug)]
pub enum AocError {
    #[error("invalid input: {0}")]
    InvalidInput(&'static str),

    #[error("invalid hand format")]
    InvalidHandFormat,

    #[error(transparent)]
    #[diagnostic(code(aoc::parse_error))]
    ParseError(#[from] std::string::ParseError),

    #[error(transparent)]
    #[diagnostic(code(aoc::parse_int_error))]
    ParseIntError(#[from] std::num::ParseIntError),
}
