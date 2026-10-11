use crate::day11::custom_error::AocError;

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let answer = solve_with_expansion(input, 2)?;
    Ok(answer.to_string())
}

fn parse_galaxies(input: &str) -> Result<Vec<(usize, usize)>, AocError> {
    let mut galaxies = Vec::new();
    for (row, line) in input.lines().enumerate() {
        for (col, ch) in line.chars().enumerate() {
            if ch == '#' {
                galaxies.push((row, col));
            }
        }
    }
    if galaxies.is_empty() {
        return Err(AocError::InvalidInput("no galaxies found"));
    }
    Ok(galaxies)
}

pub fn solve_with_expansion(input: &str, expansion_factor: usize) -> Result<u64, AocError> {
    let galaxies = parse_galaxies(input)?;

    // Find empty rows and columns
    let max_row = galaxies.iter().map(|(r, _)| *r).max().unwrap_or(0);
    let max_col = galaxies.iter().map(|(_, c)| *c).max().unwrap_or(0);

    let mut empty_rows = vec![true; max_row + 1];
    let mut empty_cols = vec![true; max_col + 1];

    for (row, col) in &galaxies {
        empty_rows[*row] = false;
        empty_cols[*col] = false;
    }

    // Calculate sum of distances with expansion
    let mut total_distance = 0u64;
    for i in 0..galaxies.len() {
        for j in (i + 1)..galaxies.len() {
            let (r1, c1) = galaxies[i];
            let (r2, c2) = galaxies[j];

            let row_start = r1.min(r2);
            let row_end = r1.max(r2);
            let col_start = c1.min(c2);
            let col_end = c1.max(c2);

            // Count empty rows and columns crossed
            let empty_rows_crossed = empty_rows[row_start..row_end]
                .iter()
                .filter(|&&e| e)
                .count();
            let empty_cols_crossed = empty_cols[col_start..col_end]
                .iter()
                .filter(|&&e| e)
                .count();

            // Calculate distance with expansion
            let base_distance = (row_end - row_start) + (col_end - col_start);
            let extra_distance = (empty_rows_crossed + empty_cols_crossed) * (expansion_factor - 1);
            let distance = base_distance + extra_distance;

            total_distance += distance as u64;
        }
    }

    Ok(total_distance)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day11.txt");
        assert_eq!("374", process(input)?);
        Ok(())
    }

    #[test]
    fn test_expansion_factor_10() -> miette::Result<()> {
        let input = include_str!("../../examples/day11.txt");
        let result = solve_with_expansion(input, 10)?;
        assert_eq!(result, 1030);
        Ok(())
    }

    #[test]
    fn test_expansion_factor_100() -> miette::Result<()> {
        let input = include_str!("../../examples/day11.txt");
        let result = solve_with_expansion(input, 100)?;
        assert_eq!(result, 8410);
        Ok(())
    }

    #[test]
    fn test_two_galaxies_with_empty_col() -> miette::Result<()> {
        let input = "#.#";
        let result = solve_with_expansion(input, 2)?;
        assert_eq!(result, 3); // base distance 2 + 1 empty col crossed
        Ok(())
    }

    #[test]
    fn test_malformed_input() {
        assert!(process("").is_err());
    }
}
