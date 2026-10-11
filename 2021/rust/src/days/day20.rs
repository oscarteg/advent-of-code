use std::collections::HashSet;

use anyhow::{bail, ensure, Context, Result};

type Point = (i64, i64);

#[derive(Clone)]
struct Image {
    lit: HashSet<Point>,
    min: Point,
    max: Point,
}

impl Image {
    fn parse(input: &str) -> Result<Self> {
        let rows: Vec<&str> = input.lines().collect();
        ensure!(!rows.is_empty(), "missing image rows");
        let width = rows[0].len();
        ensure!(width > 0, "image rows must not be empty");
        let mut lit = HashSet::new();
        for (y, row) in rows.iter().enumerate() {
            ensure!(row.len() == width, "image rows must have equal length");
            for (x, ch) in row.chars().enumerate() {
                match ch {
                    '#' => {
                        lit.insert((x as i64, y as i64));
                    }
                    '.' => {}
                    other => bail!("invalid image pixel '{other}'"),
                }
            }
        }
        Ok(Self {
            lit,
            min: (0, 0),
            max: (width as i64 - 1, rows.len() as i64 - 1),
        })
    }

    fn pixel_at(&self, x: i64, y: i64, background: bool) -> bool {
        if x < self.min.0 || x > self.max.0 || y < self.min.1 || y > self.max.1 {
            background
        } else {
            self.lit.contains(&(x, y))
        }
    }

    fn enhance(&self, algorithm: &Algorithm, background: bool) -> (Self, bool) {
        let (min_x, min_y) = (self.min.0 - 1, self.min.1 - 1);
        let (max_x, max_y) = (self.max.0 + 1, self.max.1 + 1);
        let mut lit = HashSet::new();
        for y in min_y..=max_y {
            for x in min_x..=max_x {
                let mut index = 0usize;
                for dy in -1..=1 {
                    for dx in -1..=1 {
                        index <<= 1;
                        if self.pixel_at(x + dx, y + dy, background) {
                            index |= 1;
                        }
                    }
                }
                if algorithm.enhances(index) {
                    lit.insert((x, y));
                }
            }
        }
        let next_background = algorithm.enhances(if background { 511 } else { 0 });
        (
            Self {
                lit,
                min: (min_x, min_y),
                max: (max_x, max_y),
            },
            next_background,
        )
    }
}

struct Algorithm {
    rules: [bool; 512],
}

impl Algorithm {
    fn parse(line: &str) -> Result<Self> {
        ensure!(
            line.len() == 512,
            "image enhancement algorithm must have 512 entries, got {}",
            line.len()
        );
        let mut rules = [false; 512];
        for (i, ch) in line.chars().enumerate() {
            rules[i] = match ch {
                '#' => true,
                '.' => false,
                other => bail!("invalid algorithm entry '{other}'"),
            };
        }
        Ok(Self { rules })
    }

    fn enhances(&self, index: usize) -> bool {
        self.rules[index]
    }
}

fn parse_input(input: &str) -> Result<(Algorithm, Image)> {
    let normalized = input.replace("\r\n", "\n");
    let (algorithm_part, image_part) = normalized
        .split_once("\n\n")
        .context("missing blank line between algorithm and image")?;
    let algorithm: String = algorithm_part
        .chars()
        .filter(|c| !c.is_whitespace())
        .collect();
    let algorithm = Algorithm::parse(&algorithm)?;
    let image = Image::parse(image_part.trim_end_matches('\n'))?;
    Ok((algorithm, image))
}

fn enhance_n_times(image: Image, algorithm: &Algorithm, steps: usize) -> usize {
    let mut image = image;
    let mut background = false;
    for _ in 0..steps {
        let (next_image, next_background) = image.enhance(algorithm, background);
        image = next_image;
        background = next_background;
    }
    image.lit.len()
}

fn solve_with(input: &str, part1_steps: usize, part2_steps: usize) -> Result<(usize, usize)> {
    let (algorithm, image) = parse_input(input)?;
    let part1 = enhance_n_times(image.clone(), &algorithm, part1_steps);
    let part2 = enhance_n_times(image, &algorithm, part2_steps);
    Ok((part1, part2))
}

pub fn solve(input: &str) -> Result<(usize, usize)> {
    solve_with(input, 2, 50)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day20-example.txt");
        assert_eq!(solve(input).unwrap(), (35, 3351));
    }

    #[test]
    fn single_step_matches_known_intermediate_count() {
        let input = include_str!("../../input/day20-example.txt");
        // After the first enhancement the example has 24 lit pixels.
        assert_eq!(solve_with(input, 1, 1).unwrap(), (24, 24));
    }

    #[test]
    fn missing_blank_line_is_rejected() {
        assert!(solve("...\n#..").is_err());
    }

    #[test]
    fn wrong_algorithm_length_is_rejected() {
        assert!(solve("...\n\n#..").is_err());
    }

    #[test]
    fn invalid_algorithm_character_is_rejected() {
        let mut algorithm = "#".repeat(511);
        algorithm.push('x');
        let input = format!("{algorithm}\n\n#..\n.#.\n..#");
        assert!(solve(&input).is_err());
    }

    #[test]
    fn invalid_image_character_is_rejected() {
        let algorithm = "#".repeat(512);
        let input = format!("{algorithm}\n\n#.x\n.#.\n..#");
        assert!(solve(&input).is_err());
    }

    #[test]
    fn uneven_image_rows_are_rejected() {
        let algorithm = "#".repeat(512);
        let input = format!("{algorithm}\n\n#.\n.#.");
        assert!(solve(&input).is_err());
    }
}
