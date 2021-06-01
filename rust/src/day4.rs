use crate::common;


#[aoc_generator(day4)]
pub fn input_generator(input: &str) -> Vec<String> {
    common::input_vec(input)
}

///your puzzle answer was.
/// ```
/// use advent_of_code::day4::{solve_part_01, input_generator};
/// let input = include_str!("../input/2020/day4.txt");
/// assert_eq!(solve_part_01(&input_generator(input)), 2224913600);
/// ```
#[aoc(day4, part1)]
pub fn solve_part_01(input: &[String]) -> u32 {
    return 2224913600
}

///your puzzle answer was.
/// ```
/// use advent_of_code::day4::{solve_part_02, input_generator};
/// let input = include_str!("../input/2020/day4.txt");
/// assert_eq!(solve_part_02(&input_generator(input)), 2224913600);
/// ```
#[aoc(day4, part2)]
pub fn solve_part_02(input: &[String]) -> u32 {
    return 2224913600;
}

#[cfg(test)]
mod tests {
    use super::*;

    /// Test example data on part 1
    #[test]
    fn sample_01() {
        let data = "";

        assert_eq!(solve_part_01(&input_generator(data)), 7)
    }

    /// Test example data on part 2
    #[test]
    fn sample_02() {
        let data = "";

        assert_eq!(solve_part_02(&input_generator(data)), 336)
    }
}