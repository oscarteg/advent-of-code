use miette::{Diagnostic, SourceSpan};
use thiserror::Error;

#[derive(Error, Diagnostic, Debug)]
pub enum AocError {
    #[error(transparent)]
    #[diagnostic(code(aoc::io_error))]
    IoError(#[from] std::io::Error),

    #[error("Failed to parse '{input}' as a color")]
    #[diagnostic(code(aoc::parse_error))]
    ParseError {
        input: String,
        #[source_code]
        source: String, // optional, use if you have source code to highlight
        #[label("this is not a valid color")]
        span: SourceSpan, // optional, use if you have a specific span to highlight
    },
}
