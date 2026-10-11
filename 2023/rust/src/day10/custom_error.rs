use miette::Diagnostic;
use thiserror::Error;

#[derive(Error, Diagnostic, Debug)]
pub enum AocError {
    #[error("invalid pipe maze: {0}")]
    InvalidInput(&'static str),

    #[error("invalid tile character: '{0}'")]
    InvalidTile(char),

    #[error(transparent)]
    #[diagnostic(code(aoc::io_error))]
    IoError(#[from] std::io::Error),
}
