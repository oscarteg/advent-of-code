use miette::Diagnostic;
use thiserror::Error;

#[derive(Error, Diagnostic, Debug)]
pub enum AocError {
    #[error("invalid trail map: {0}")]
    #[diagnostic(code(aoc::invalid_input))]
    InvalidInput(&'static str),

    #[error("invalid trail map: unknown tile {0:?}")]
    #[diagnostic(code(aoc::invalid_input))]
    UnknownTile(char),

    #[error("no hike reaches the exit")]
    #[diagnostic(code(aoc::no_hike))]
    NoHike,

    #[error(transparent)]
    #[diagnostic(code(aoc::io_error))]
    IoError(#[from] std::io::Error),
}
