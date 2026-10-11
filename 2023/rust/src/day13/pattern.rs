use crate::day13::custom_error::AocError;

/// A grid of rocks (`true`) and ash (`.` -> `false`).
#[derive(Debug)]
pub(crate) struct Pattern {
    rows: Vec<Vec<bool>>,
}

impl Pattern {
    fn parse(block: &str) -> Result<Self, AocError> {
        let rows = block
            .lines()
            .map(|line| {
                line.chars()
                    .map(|c| match c {
                        '#' => Ok(true),
                        '.' => Ok(false),
                        _ => Err(AocError::InvalidInput("pattern contains unknown symbol")),
                    })
                    .collect::<Result<Vec<_>, _>>()
            })
            .collect::<Result<Vec<_>, _>>()?;
        if rows.is_empty() {
            return Err(AocError::InvalidInput("pattern has no rows"));
        }
        let width = rows[0].len();
        if width == 0 || rows.iter().any(|row| row.len() != width) {
            return Err(AocError::InvalidInput(
                "pattern rows have inconsistent width",
            ));
        }
        Ok(Self { rows })
    }

    fn width(&self) -> usize {
        self.rows[0].len()
    }

    fn column(&self, index: usize) -> Vec<bool> {
        self.rows.iter().map(|row| row[index]).collect()
    }

    fn diff_count(a: &[bool], b: &[bool]) -> usize {
        a.iter().zip(b).filter(|(x, y)| x != y).count()
    }

    /// Position of the horizontal mirror line whose reflected rows differ by
    /// exactly `smudges` cells in total, expressed as the row count above it.
    fn horizontal_reflection(&self, smudges: usize) -> Option<usize> {
        (1..self.rows.len()).find(|&line| {
            let span = line.min(self.rows.len() - line);
            let diff: usize = (0..span)
                .map(|offset| {
                    Self::diff_count(&self.rows[line - 1 - offset], &self.rows[line + offset])
                })
                .sum();
            diff == smudges
        })
    }

    /// Position of the vertical mirror line, expressed as the column count to its left.
    fn vertical_reflection(&self, smudges: usize) -> Option<usize> {
        let width = self.width();
        (1..width).find(|&line| {
            let span = line.min(width - line);
            let diff: usize = (0..span)
                .map(|offset| {
                    Self::diff_count(&self.column(line - 1 - offset), &self.column(line + offset))
                })
                .sum();
            diff == smudges
        })
    }

    /// Summary value for this pattern's mirror line: `100 * rows above` for a
    /// horizontal mirror, or `columns to the left` for a vertical one.
    pub fn summarize(&self, smudges: usize) -> Result<usize, AocError> {
        if let Some(line) = self.horizontal_reflection(smudges) {
            return Ok(100 * line);
        }
        if let Some(line) = self.vertical_reflection(smudges) {
            return Ok(line);
        }
        Err(AocError::InvalidInput("pattern has no mirror line"))
    }
}

pub(crate) fn parse_patterns(input: &str) -> Result<Vec<Pattern>, AocError> {
    let patterns = input
        .split("\n\n")
        .map(str::trim)
        .filter(|block| !block.is_empty())
        .map(Pattern::parse)
        .collect::<Result<Vec<_>, _>>()?;
    if patterns.is_empty() {
        return Err(AocError::InvalidInput("notes contain no patterns"));
    }
    Ok(patterns)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn rejects_unknown_symbols() {
        assert!(Pattern::parse("#.x\n...").is_err());
    }

    #[test]
    fn rejects_inconsistent_widths() {
        assert!(Pattern::parse("##\n#").is_err());
    }

    #[test]
    fn rejects_empty_notes() {
        assert!(parse_patterns("").is_err());
        assert!(parse_patterns("\n\n").is_err());
    }
}
