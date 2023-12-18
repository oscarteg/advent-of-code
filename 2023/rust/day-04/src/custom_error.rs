use miette::Diagnostic;
use thiserror::Error;

#[derive(Error, Diagnostic, Debug)]
pub enum AocError {
    #[error(transparent)]
    #[diagnostic(code(aoc::io_error))]
    IoError(#[from] std::io::Error),

    #[error(transparent)]
    #[diagnostic(code(aoc::parse_error))]
    ParseError(#[from] std::string::ParseError),

    #[error(transparent)]
    #[diagnostic(code(aoc::parse_int_error))]
    ParseIntError(#[from] std::num::ParseIntError),
}
