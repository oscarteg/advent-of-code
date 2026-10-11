use anyhow::{bail, ensure, Context, Result};
use std::cmp::Reverse;
use std::collections::{BinaryHeap, HashMap};

const HALL_LEN: usize = 11;
const ROOM_COUNT: usize = 4;
/// Slots per room: the two rows of the diagram plus the two part two unfolds.
const MAX_DEPTH: usize = 4;
/// Hallway cells in front of the doors; amphipods pass through but never stop there.
const DOORS: [usize; ROOM_COUNT] = [2, 4, 6, 8];
/// The room rows part two reveals between the two rows of the diagram.
const FOLDED_ROWS: [[Amphipod; ROOM_COUNT]; 2] = [
    [
        Amphipod::Desert,
        Amphipod::Copper,
        Amphipod::Bronze,
        Amphipod::Amber,
    ],
    [
        Amphipod::Desert,
        Amphipod::Bronze,
        Amphipod::Amber,
        Amphipod::Copper,
    ],
];

#[derive(Clone, Copy, Debug, PartialEq, Eq, PartialOrd, Ord, Hash)]
enum Amphipod {
    Amber,
    Bronze,
    Copper,
    Desert,
}

impl Amphipod {
    const ALL: [Self; ROOM_COUNT] = [Self::Amber, Self::Bronze, Self::Copper, Self::Desert];

    fn home(self) -> usize {
        match self {
            Self::Amber => 0,
            Self::Bronze => 1,
            Self::Copper => 2,
            Self::Desert => 3,
        }
    }

    fn step_cost(self) -> usize {
        match self {
            Self::Amber => 1,
            Self::Bronze => 10,
            Self::Copper => 100,
            Self::Desert => 1000,
        }
    }
}

impl TryFrom<char> for Amphipod {
    type Error = anyhow::Error;

    fn try_from(cell: char) -> Result<Self> {
        match cell {
            'A' => Ok(Self::Amber),
            'B' => Ok(Self::Bronze),
            'C' => Ok(Self::Copper),
            'D' => Ok(Self::Desert),
            other => bail!("{other:?} is not an amphipod"),
        }
    }
}

/// A burrow state: the hallway plus the first `depth` slots of every room.
#[derive(Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Hash)]
struct Burrow {
    hall: [Option<Amphipod>; HALL_LEN],
    /// `rooms[room][0]` is the slot next to the hallway; rooms fill from the back.
    rooms: [[Option<Amphipod>; MAX_DEPTH]; ROOM_COUNT],
    depth: usize,
}

impl Burrow {
    /// Reads the diagram: a wall row, an empty hallway, then one amphipod per room
    /// and row, left to right.
    fn parse(input: &str) -> Result<Self> {
        let mut rows = input.lines().map(str::trim).filter(|row| !row.is_empty());
        let wall = rows.next().context("missing burrow diagram")?;
        ensure!(
            wall.chars().all(|cell| cell == '#'),
            "the diagram starts with a wall row"
        );
        let hallway = rows.next().context("missing hallway row")?;
        ensure!(
            hallway.len() == HALL_LEN + 2 && hallway.matches('.').count() == HALL_LEN,
            "the hallway must be {HALL_LEN} empty cells wide"
        );
        let mut rooms = [[None; MAX_DEPTH]; ROOM_COUNT];
        let mut depth = 0;
        for row in rows {
            let cells: Vec<char> = row.chars().filter(|cell| *cell != '#').collect();
            if cells.is_empty() {
                break; // the closing wall row
            }
            ensure!(
                depth < MAX_DEPTH,
                "a burrow has at most {MAX_DEPTH} room rows"
            );
            let width = cells.len();
            ensure!(
                width == ROOM_COUNT,
                "a room row holds {ROOM_COUNT} amphipods, found {width}"
            );
            for (room, &cell) in cells.iter().enumerate() {
                rooms[room][depth] = Some(Amphipod::try_from(cell)?);
            }
            depth += 1;
        }
        ensure!(depth > 0, "the burrow has no room rows");
        let burrow = Self {
            hall: [None; HALL_LEN],
            rooms,
            depth,
        };
        for kind in Amphipod::ALL {
            let found = burrow
                .rooms
                .iter()
                .flatten()
                .filter(|slot| **slot == Some(kind))
                .count();
            ensure!(
                found == depth,
                "the burrow needs {depth} {kind:?} amphipods, found {found}"
            );
        }
        Ok(burrow)
    }

    /// Part two inserts two more room rows between the two rows of the diagram.
    fn unfold(&self) -> Result<Self> {
        ensure!(self.depth == 2, "only a two-row burrow can be unfolded");
        let mut unfolded = Self {
            depth: MAX_DEPTH,
            ..*self
        };
        for (room, slots) in unfolded.rooms.iter_mut().enumerate() {
            *slots = [
                self.rooms[room][0],
                Some(FOLDED_ROWS[0][room]),
                Some(FOLDED_ROWS[1][room]),
                self.rooms[room][1],
            ];
        }
        Ok(unfolded)
    }

    fn topmost(&self, room: usize) -> Option<(usize, Amphipod)> {
        (0..self.depth).find_map(|slot| self.rooms[room][slot].map(|amphipod| (slot, amphipod)))
    }

    /// True when `room` holds only its own kind: nobody has to leave it, and that
    /// kind is allowed to move in.
    fn settled(&self, room: usize) -> bool {
        self.rooms[room][..self.depth]
            .iter()
            .flatten()
            .all(|amphipod| amphipod.home() == room)
    }

    fn is_organized(&self) -> bool {
        (0..ROOM_COUNT).all(|room| self.rooms[room][0].is_some() && self.settled(room))
    }

    /// True when every hallway cell between `from` and `to` is free, `from` excluded.
    fn hall_free(&self, from: usize, to: usize) -> bool {
        (from.min(to)..=from.max(to))
            .filter(|cell| *cell != from)
            .all(|cell| self.hall[cell].is_none())
    }

    /// Every state one move away, with the energy that move costs.
    ///
    /// Walking an amphipod home from the hallway is always part of an optimal
    /// solution: the move costs the same whenever it happens, it frees a hallway
    /// cell and it fills a slot that amphipod has to occupy anyway. Taking that
    /// move on its own keeps the search small.
    fn moves(&self) -> Vec<(usize, Self)> {
        if let Some(move_home) = self.move_home() {
            return vec![move_home];
        }
        self.moves_into_hallway()
    }

    /// The first amphipod that can walk from the hallway straight into its room.
    fn move_home(&self) -> Option<(usize, Self)> {
        self.hall.iter().enumerate().find_map(|(cell, occupant)| {
            let amphipod = (*occupant)?;
            let room = amphipod.home();
            let door = DOORS[room];
            let free = self.topmost(room).map_or(self.depth, |(slot, _)| slot);
            let slot = free.checked_sub(1)?;
            (self.settled(room) && self.hall_free(cell, door)).then(|| {
                let mut next = *self;
                next.hall[cell] = None;
                next.rooms[room][slot] = Some(amphipod);
                let steps = cell.abs_diff(door) + slot + 1;
                (amphipod.step_cost() * steps, next)
            })
        })
    }

    /// Every move that walks the amphipod nearest the hallway out of a room that
    /// still holds a stranger.
    fn moves_into_hallway(&self) -> Vec<(usize, Self)> {
        let mut moves = Vec::new();
        for (room, &door) in DOORS.iter().enumerate() {
            if self.settled(room) {
                continue;
            }
            let Some((slot, amphipod)) = self.topmost(room) else {
                continue;
            };
            for cell in (0..HALL_LEN).filter(|cell| !DOORS.contains(cell)) {
                if !self.hall_free(door, cell) {
                    continue;
                }
                let mut next = *self;
                next.rooms[room][slot] = None;
                next.hall[cell] = Some(amphipod);
                let steps = slot + 1 + door.abs_diff(cell);
                moves.push((amphipod.step_cost() * steps, next));
            }
        }
        moves
    }
}

/// Least energy that organizes the burrow, by Dijkstra over burrow states.
fn least_energy(start: Burrow) -> Result<usize> {
    let mut best = HashMap::from([(start, 0)]);
    let mut queue = BinaryHeap::from([Reverse((0, start))]);
    while let Some(Reverse((energy, burrow))) = queue.pop() {
        if burrow.is_organized() {
            return Ok(energy);
        }
        if best.get(&burrow).is_some_and(|&shortest| shortest < energy) {
            continue;
        }
        for (cost, next) in burrow.moves() {
            let next_energy = energy + cost;
            if best
                .get(&next)
                .is_none_or(|&shortest| next_energy < shortest)
            {
                best.insert(next, next_energy);
                queue.push(Reverse((next_energy, next)));
            }
        }
    }
    bail!("the amphipods cannot be organized")
}

/// Returns the least energy for the diagram and for the unfolded burrow.
pub fn solve(input: &str) -> Result<(usize, usize)> {
    let burrow = Burrow::parse(input)?;
    Ok((least_energy(burrow)?, least_energy(burrow.unfold()?)?))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day23-example.txt");
    const ORGANIZED: &str =
        "#############\n#...........#\n###A#B#C#D###\n  #A#B#C#D#\n  #########\n";

    #[test]
    fn official_example() {
        assert_eq!(solve(EXAMPLE).unwrap(), (12521, 44169));
    }

    #[test]
    fn an_organized_burrow_needs_no_energy() {
        let (part1, part2) = solve(ORGANIZED).unwrap();
        assert_eq!(part1, 0);
        assert!(part2 > 0, "the unfolded rows start out scrambled");
    }

    #[test]
    fn swapped_rooms_step_aside() {
        // Amber leaves room B for hallway cell 5 (2 steps), Bronze crosses into room B
        // along its shortest path (4 steps) and Amber walks on into room A (4 steps):
        // 6 energy for Amber, 40 for Bronze.
        let input = "#############\n#...........#\n###B#A#C#D###\n  #########\n";
        let burrow = Burrow::parse(input).unwrap();
        assert_eq!(least_energy(burrow).unwrap(), 46);
    }

    #[test]
    fn windows_newlines_are_accepted() {
        let burrow = Burrow::parse(&ORGANIZED.replace('\n', "\r\n")).unwrap();
        assert_eq!(least_energy(burrow).unwrap(), 0);
    }

    #[test]
    fn only_a_two_row_diagram_can_be_unfolded() {
        let input = "#############\n#...........#\n###A#B#C#D###\n  #########\n";
        assert!(Burrow::parse(input).is_ok());
        assert!(solve(input).is_err());
    }

    #[test]
    fn malformed_input_is_rejected() {
        let inputs = [
            "",
            "#############\n",
            "#############\n#..B........#\n###A#B#C#D###\n  #A#B#C#D#\n  #########\n",
            "#############\n#...........#\n###A#B#X#D###\n  #A#B#C#D#\n  #########\n",
            "#############\n#...........#\n###A#B#C###\n  #A#B#C#D#\n  #########\n",
            "#############\n#...........#\n###A#B#C#D###\n  #A#B#C#A#\n  #########\n",
            "#############\n#...........#\n  #########\n",
            "#############\n#...........#\n###A#B#C#D###\n  #A#B#C#D#\n  #A#B#C#D#\n  #A#B#C#D#\n  #A#B#C#D#\n  #########\n",
        ];
        for input in inputs {
            assert!(solve(input).is_err(), "accepted {input:?}");
        }
    }
}
