#[cfg(test)]
mod tests {
    use anyhow::Error;
    use std::collections::{BTreeMap, VecDeque};

    use crate::utils::read_file;

    type Stack = VecDeque<char>;
    type Move = (u32, u32, u32);

    fn parse_input(input: &str) -> (BTreeMap<u32, Stack>, Vec<Move>) {
        let (stacks_input, moves_input) = input.trim().split_once("\n\n").unwrap();
        let mut stacks = BTreeMap::new();

        let lines = stacks_input.lines().collect::<Vec<_>>();

        let num_stacks = lines.last().unwrap().trim().split_whitespace().count();

        for stack_index in 0..num_stacks {
            let mut stack: VecDeque<char> = VecDeque::new();
            for line in &lines {
                let crate_str = line.split_whitespace().nth(stack_index);
                if let Some(crate_str) = crate_str {
                    if crate_str != " " && !crate_str.is_empty() {
                        println!("crate_str: {:?}", crate_str);
                        stack.push_front(
                            crate_str
                                .chars()
                                .nth(1)
                                .or_else(|| crate_str.chars().next())
                                .unwrap(),
                        );
                    }
                }
            }
            stacks.insert(stack_index as u32 + 1, stack);
        }

        println!("{:?}", stacks);

        // Parse the stacks and crates

        let moves = moves_input
            .lines()
            .map(|line| {
                let parts = line.split_whitespace().collect::<Vec<_>>();
                let num_crates = parts[1].parse::<u32>().unwrap();
                let from_stack = parts[3].parse::<u32>().unwrap();
                let to_stack = parts[5].parse::<u32>().unwrap();
                (num_crates, from_stack, to_stack)
            })
            .collect::<Vec<_>>();

        (stacks, moves)
    }

    fn simulate_moves(
        mut stacks: BTreeMap<u32, VecDeque<char>>,
        moves: Vec<(u32, u32, u32)>,
    ) -> BTreeMap<u32, VecDeque<char>> {
        for (num_crates, from_stack, to_stack) in moves {
            for _ in 0..num_crates {
                if let Some(crate_char) = stacks
                    .get_mut(&from_stack)
                    .and_then(|stack| stack.pop_front())
                {
                    stacks.entry(to_stack).or_default().push_front(crate_char);
                }
            }
        }
        stacks
    }

    fn get_top_crates(stacks: BTreeMap<u32, VecDeque<char>>) -> String {
        let mut result = String::new();
        let mut sorted_keys: Vec<_> = stacks.keys().cloned().collect();
        sorted_keys.sort();

        for key in sorted_keys {
            if let Some(stack) = stacks.get(&key) {
                if let Some(&crate_char) = stack.front() {
                    result.push(crate_char);
                }
            }
        }
        result
    }

    fn part_one(input: &str) -> Result<String, Error> {
        let (initial_stacks, moves) = parse_input(input);
        let final_stacks = simulate_moves(initial_stacks, moves);
        // println!("{:?}", final_stacks);
        let result = get_top_crates(final_stacks);
        Ok(result)
    }

    const INPUT: &str = "    [D]    
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
        assert_eq!(part_one(INPUT).unwrap(), "CMZ");

        // Real input
        let file: String = read_file("input/day5.txt");
        assert_eq!(part_one(file.as_str()).unwrap(), "JXFSZCNLHY");
    }
}
