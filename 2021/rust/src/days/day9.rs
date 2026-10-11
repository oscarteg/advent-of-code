use anyhow::{ensure, Context, Result};

struct HeightMap {
    heights: Vec<u8>,
    width: usize,
    height: usize,
}

impl HeightMap {
    fn parse(input: &str) -> Result<Self> {
        let rows: Vec<&str> = input.lines().filter(|line| !line.is_empty()).collect();
        ensure!(!rows.is_empty(), "height map must have at least one row");
        let width = rows[0].len();
        ensure!(width > 0, "height map rows must not be empty");
        let mut heights = Vec::with_capacity(width * rows.len());
        for row in &rows {
            ensure!(row.len() == width, "all rows must have the same width");
            for ch in row.chars() {
                let digit = ch
                    .to_digit(10)
                    .with_context(|| format!("invalid height digit '{ch}'"))?;
                heights.push(u8::try_from(digit).context("height digit out of range")?);
            }
        }
        Ok(Self {
            heights,
            width,
            height: rows.len(),
        })
    }

    fn at(&self, x: usize, y: usize) -> u8 {
        self.heights[y * self.width + x]
    }

    fn neighbors(&self, x: usize, y: usize) -> impl Iterator<Item = (usize, usize)> + '_ {
        let width = self.width;
        let height = self.height;
        [(0_i32, -1), (0, 1), (-1, 0), (1, 0)]
            .into_iter()
            .filter_map(move |(dx, dy)| {
                let nx = x as i32 + dx;
                let ny = y as i32 + dy;
                (nx >= 0 && ny >= 0 && (nx as usize) < width && (ny as usize) < height)
                    .then_some((nx as usize, ny as usize))
            })
    }

    fn is_low_point(&self, x: usize, y: usize) -> bool {
        let here = self.at(x, y);
        self.neighbors(x, y).all(|(nx, ny)| self.at(nx, ny) > here)
    }

    fn low_points(&self) -> impl Iterator<Item = (usize, usize)> + '_ {
        (0..self.height)
            .flat_map(move |y| (0..self.width).map(move |x| (x, y)))
            .filter(move |&(x, y)| self.is_low_point(x, y))
    }

    fn basin_size(&self, start: (usize, usize)) -> usize {
        let mut visited = vec![false; self.heights.len()];
        let mut stack = vec![start];
        let mut size = 0;
        while let Some((x, y)) = stack.pop() {
            let index = y * self.width + x;
            if visited[index] || self.at(x, y) == 9 {
                continue;
            }
            visited[index] = true;
            size += 1;
            stack.extend(self.neighbors(x, y));
        }
        size
    }
}

pub fn solve(input: &str) -> Result<(u32, usize)> {
    let map = HeightMap::parse(input)?;
    let low_points: Vec<_> = map.low_points().collect();

    let risk_level: u32 = low_points
        .iter()
        .map(|&(x, y)| u32::from(map.at(x, y)) + 1)
        .sum();

    let mut basin_sizes: Vec<usize> = low_points
        .iter()
        .map(|&point| map.basin_size(point))
        .collect();
    basin_sizes.sort_unstable_by(|a, b| b.cmp(a));
    ensure!(basin_sizes.len() >= 3, "fewer than three basins found");
    let largest_basins: usize = basin_sizes[..3].iter().product();

    Ok((risk_level, largest_basins))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day9-example.txt");
        assert_eq!(solve(input).unwrap(), (15, 1134));
    }

    #[test]
    fn flat_grid_with_no_low_points_is_rejected() {
        let input = "11\n11";
        assert!(solve(input).is_err());
    }

    #[test]
    fn three_small_basins_are_multiplied() {
        let input = "109\n919\n091";
        assert_eq!(solve(input).unwrap(), (4, 3));
    }

    #[test]
    fn malformed_rows_are_rejected() {
        assert!(HeightMap::parse("123\n12").is_err());
        assert!(HeightMap::parse("12a\n123").is_err());
        assert!(HeightMap::parse("").is_err());
    }

    #[test]
    fn malformed_input_is_rejected_by_solve() {
        assert!(solve("12a").is_err());
    }
}
