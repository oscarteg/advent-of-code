use anyhow::{bail, ensure};
use std::collections::{HashMap, HashSet};

type Pos = (i32, i32);
type Part1 = i32;
type Part2 = i32;

struct State {
    elves: HashSet<Pos>,
    direction_idx: usize,
}

impl State {
    fn from_input(input: &str) -> anyhow::Result<Self> {
        let mut elves = HashSet::new();
        for (y, line) in input.lines().enumerate() {
            for (x, ch) in line.chars().enumerate() {
                match ch {
                    '#' => {
                        elves.insert((x as i32, y as i32));
                    }
                    '.' => {}
                    _ => bail!("Invalid character: {}", ch),
                }
            }
        }
        ensure!(!elves.is_empty(), "No elves found");
        Ok(State {
            elves,
            direction_idx: 0,
        })
    }

    fn has_neighbors(&self, pos: Pos) -> bool {
        for dx in -1..=1 {
            for dy in -1..=1 {
                if (dx, dy) == (0, 0) {
                    continue;
                }
                if self.elves.contains(&(pos.0 + dx, pos.1 + dy)) {
                    return true;
                }
            }
        }
        false
    }

    fn can_move_north(&self, pos: Pos) -> bool {
        !self.elves.contains(&(pos.0 - 1, pos.1 - 1))
            && !self.elves.contains(&(pos.0, pos.1 - 1))
            && !self.elves.contains(&(pos.0 + 1, pos.1 - 1))
    }

    fn can_move_south(&self, pos: Pos) -> bool {
        !self.elves.contains(&(pos.0 - 1, pos.1 + 1))
            && !self.elves.contains(&(pos.0, pos.1 + 1))
            && !self.elves.contains(&(pos.0 + 1, pos.1 + 1))
    }

    fn can_move_west(&self, pos: Pos) -> bool {
        !self.elves.contains(&(pos.0 - 1, pos.1 - 1))
            && !self.elves.contains(&(pos.0 - 1, pos.1))
            && !self.elves.contains(&(pos.0 - 1, pos.1 + 1))
    }

    fn can_move_east(&self, pos: Pos) -> bool {
        !self.elves.contains(&(pos.0 + 1, pos.1 - 1))
            && !self.elves.contains(&(pos.0 + 1, pos.1))
            && !self.elves.contains(&(pos.0 + 1, pos.1 + 1))
    }

    fn propose_move(&self, pos: Pos) -> Option<Pos> {
        if !self.has_neighbors(pos) {
            return None;
        }

        for i in 0..4 {
            let dir = (self.direction_idx + i) % 4;
            let new_pos = match dir {
                0 if self.can_move_north(pos) => (pos.0, pos.1 - 1),
                1 if self.can_move_south(pos) => (pos.0, pos.1 + 1),
                2 if self.can_move_west(pos) => (pos.0 - 1, pos.1),
                3 if self.can_move_east(pos) => (pos.0 + 1, pos.1),
                _ => continue,
            };
            return Some(new_pos);
        }
        None
    }

    fn round(&mut self) -> bool {
        let mut proposals: HashMap<Pos, Vec<Pos>> = HashMap::new();

        for &elf in &self.elves {
            if let Some(new_pos) = self.propose_move(elf) {
                proposals.entry(new_pos).or_default().push(elf);
            }
        }

        let mut moved = false;
        for (new_pos, elves) in proposals {
            if elves.len() == 1 {
                let old_pos = elves[0];
                self.elves.remove(&old_pos);
                self.elves.insert(new_pos);
                moved = true;
            }
        }

        self.direction_idx = (self.direction_idx + 1) % 4;
        moved
    }

    fn count_empty_tiles(&self) -> i32 {
        let xs: Vec<i32> = self.elves.iter().map(|p| p.0).collect();
        let ys: Vec<i32> = self.elves.iter().map(|p| p.1).collect();

        let min_x = *xs.iter().min().unwrap();
        let max_x = *xs.iter().max().unwrap();
        let min_y = *ys.iter().min().unwrap();
        let max_y = *ys.iter().max().unwrap();

        let total_area = (max_x - min_x + 1) * (max_y - min_y + 1);
        total_area - self.elves.len() as i32
    }
}

pub fn solve(input: &str) -> anyhow::Result<(Part1, Part2)> {
    let mut state = State::from_input(input)?;

    let mut part2 = 0;
    for round_num in 1..=10 {
        if !state.round() {
            part2 = round_num;
            break;
        }
    }
    let part1 = state.count_empty_tiles();

    // Part 2: if no movement was found in first 10 rounds, continue
    if part2 == 0 {
        let mut round_num = 10;
        loop {
            round_num += 1;
            if !state.round() {
                part2 = round_num;
                break;
            }
        }
    }

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_part1_example() {
        let input = include_str!("../../input/day23-example.txt");
        let (part1, _) = solve(input).unwrap();
        assert_eq!(part1, 110);
    }

    #[test]
    fn test_part2_example() {
        let input = include_str!("../../input/day23-example.txt");
        let (_, part2) = solve(input).unwrap();
        assert_eq!(part2, 20);
    }

    #[test]
    fn test_single_elf() {
        let input = "#";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 0);
        assert_eq!(part2, 1);
    }

    #[test]
    fn test_malformed_input() {
        let input = "#\n$\n#";
        assert!(solve(input).is_err());
    }

    #[test]
    fn test_empty_input() {
        let input = "";
        assert!(solve(input).is_err());
    }
}
