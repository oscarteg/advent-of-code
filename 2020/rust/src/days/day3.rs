fn slope_finder(grid: &[&str], down: usize, right: usize) -> u32 {
    let mut trees = 0;
    let mut row = 0;
    let mut col = 0;

    while row < grid.len() {
        let line = grid[row];
        let width = line.len();
        if line.as_bytes()[col % width] == b'#' {
            trees += 1;
        }
        row += down;
        col += right;
    }

    trees
}

pub fn solve(input: &str) -> anyhow::Result<(u32, u32)> {
    let grid: Vec<&str> = input.lines().filter(|line| !line.is_empty()).collect();

    anyhow::ensure!(!grid.is_empty(), "empty grid");
    for line in &grid {
        anyhow::ensure!(
            !line.is_empty() && line.chars().all(|c| c == '.' || c == '#'),
            "malformed line: {}",
            line
        );
    }

    let part1 = slope_finder(&grid, 1, 3);

    let slopes = [(1, 1), (1, 3), (1, 5), (1, 7), (2, 1)];
    let part2 = slopes.iter().fold(1u32, |acc, (down, right)| {
        acc * slope_finder(&grid, *down, *right)
    });

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "..##.......
#...#...#..
.#....#..#.
..#.#...#.#
.#...##..#.
..#.##.....
.#.#.#....#
.#........#
#.##...#...
#...##....#
.#..#...#.#";

    #[test]
    fn example_part1() {
        let (part1, _) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 7);
    }

    #[test]
    fn example_part2() {
        let (_, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part2, 336);
    }

    #[test]
    fn empty_input() {
        assert!(solve("").is_err());
    }

    #[test]
    fn invalid_chars() {
        assert!(solve("..#\n.x#").is_err());
    }

    #[test]
    fn single_line() {
        let (part1, part2) = solve("#..").unwrap();
        assert_eq!(part1, 1);
        assert_eq!(part2, 1);
    }
}
