use std::collections::HashMap;

struct Move(i32, i32, i32);

#[derive(Debug)]
struct Crate<'a> {
    pub crates: HashMap<&'a str, &'a str>,
}

fn part_one(input: &str) -> Option<&str> {
    let (crates, procedure) = input.trim().split_once("\n\n").unwrap();

    // Split

    Some("CMZ")
}

fn part_two(input: &str) -> Option<&str> {
    None
}

#[cfg(test)]
mod tests {
    use crate::utils::read_file;

    use super::*;

    const INPUT: &str = "[D]    
[N] [C]    
[Z] [M] [P]
 1   2   3 

move 1 from 2 to 1
move 3 from 1 to 3
move 2 from 2 to 1
move 1 from 1 to 2";

    #[test]
    fn test_part_1() {
        // Example
        assert_eq!(part_one(INPUT), Some("CMZ"));

        // Real input
        let file: String = read_file("input/day5.txt");
        assert_eq!(part_one(file.as_str()), Some(607));
    }

    #[test]
    fn test_part_2() {
        // Example
        assert_eq!(part_two(INPUT), Some(4));

        // Real input
        let file: String = read_file("input/day5.txt");
        assert_eq!(part_two(file.as_str()), Some(914));
    }
}
