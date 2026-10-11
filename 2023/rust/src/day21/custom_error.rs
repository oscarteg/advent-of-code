use miette::Diagnostic;
use thiserror::Error;

#[derive(Error, Diagnostic, Debug)]
pub enum AocError {
    #[error("invalid map: {0}")]
    #[diagnostic(code(aoc::invalid_input))]
    InvalidInput(&'static str),

    #[error(transparent)]
    #[diagnostic(code(aoc::io_error))]
    IoError(#[from] std::io::Error),
}
