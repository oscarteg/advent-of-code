use crate::day25::{custom_error::AocError, wiring::Diagram};

/// Disconnects the three wires of the diagram's minimum cut and multiplies the
/// sizes of the two groups that remain.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let diagram = Diagram::parse(input)?;
    let (left, right) = diagram.minimum_cut();
    Ok((left * right).to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day25.txt");
        assert_eq!("54", process(input)?);
        Ok(())
    }

    #[rstest]
    // Two components joined by a single wire.
    #[case("jqt: rhn", "1")]
    // Two triangles joined by a single wire: 3 * 3.
    #[case("a: b c\nb: c\nc: d\nd: e f\ne: f", "9")]
    // Two five-component cliques joined by three wires: 5 * 5. Every cheaper
    // split would have to cut a clique, which costs at least four wires.
    #[case(
        "a: b c d e f\nb: c d e g\nc: d e h\nd: e\nf: g h i j\ng: h i j\nh: i j\ni: j",
        "25"
    )]
    fn small_diagrams(#[case] input: &str, #[case] expected: &str) {
        assert_eq!(expected, process(input).unwrap());
    }

    #[rstest]
    #[case("")]
    #[case("jqt rhn xhk")]
    #[case("jqt:")]
    #[case(": rhn")]
    #[case("jqt: jqt")]
    fn malformed_input_returns_error(#[case] input: &str) {
        assert!(process(input).is_err(), "{input}");
    }
}
