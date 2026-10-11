use miette::Diagnostic;
use thiserror::Error;

#[derive(Error, Diagnostic, Debug)]
pub enum AocError {
    #[error("no digits found in line: {0}")]
    #[diagnostic(code(aoc::no_digits))]
    NoDigitsFound(String),

    #[error("failed to parse number: {0}")]
    #[diagnostic(code(aoc::parse_error))]
    ParseError(String),

    #[error(transparent)]
    #[diagnostic(code(aoc::io_error))]
    IoError(#[from] std::io::Error),
}
