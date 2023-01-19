#![allow(dead_code)]
#![allow(unused_variables)]

pub fn part_one(input: &str) -> Option<u32> {
    let mut out = 0;

    let foo = input
        .trim()
        .lines()
        .map(|l| l.split_once(",").unwrap())
        .map(|(a, b)| (split_range(a), split_range(b)));

    for (x, y) in foo {
        out += ((x.0 >= y.0 && x.1 <= y.1) || (y.0 >= x.0 && y.1 <= x.1)) as u32;
    }

    Some(out)
}

fn split_range(range: &str) -> (u32, u32) {
    let (a, b) = range.split_once("-").unwrap();

    (a.parse::<u32>().unwrap(), b.parse::<u32>().unwrap())
}

pub fn part_two(input: &str) -> Option<u32> {
    let mut out = 0;
    let foo = input
        .trim()
        .lines()
        .map(|l| l.split_once(",").unwrap())
        .map(|(a, b)| (split_range(a), split_range(b)));

    for (x, y) in foo {
        out += (x.0.max(y.0) <= x.1.min(y.1)) as u32;
    }

    Some(out)
}

#[cfg(test)]
mod tests {
    use crate::utils::read_file;

    use super::*;

    const INPUT: &str = "2-4,6-8
2-3,4-5
5-7,7-9
2-8,3-7
6-6,4-6
2-6,4-8";

    #[test]
    fn test_part_1() {
        // Example
        assert_eq!(part_one(INPUT), Some(2));

        // Real input
        let file: String = read_file("input/day4.txt");
        assert_eq!(part_one(file.as_str()), Some(605));
    }

    #[test]
    fn test_part_2() {
        // Example
        assert_eq!(part_two(INPUT), Some(4));

        // Real input
        let file: String = read_file("input/day4.txt");
        assert_eq!(part_two(file.as_str()), Some(914));
    }
}
