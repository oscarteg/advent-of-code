use anyhow::{anyhow, Result};
use std::collections::{HashSet, VecDeque};

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
struct Cube(i32, i32, i32);

impl Cube {
    fn neighbors(&self) -> [Cube; 6] {
        let Cube(x, y, z) = self;
        [
            Cube(x + 1, *y, *z),
            Cube(x - 1, *y, *z),
            Cube(*x, y + 1, *z),
            Cube(*x, y - 1, *z),
            Cube(*x, *y, z + 1),
            Cube(*x, *y, z - 1),
        ]
    }
}

fn parse_input(input: &str) -> Result<Vec<Cube>> {
    let mut cubes = Vec::new();
    for line in input.trim().lines() {
        let parts: Vec<&str> = line.split(',').collect();
        if parts.len() != 3 {
            return Err(anyhow!("Invalid cube format: {}", line));
        }
        let x = parts[0]
            .parse::<i32>()
            .map_err(|_| anyhow!("Invalid x coordinate: {}", parts[0]))?;
        let y = parts[1]
            .parse::<i32>()
            .map_err(|_| anyhow!("Invalid y coordinate: {}", parts[1]))?;
        let z = parts[2]
            .parse::<i32>()
            .map_err(|_| anyhow!("Invalid z coordinate: {}", parts[2]))?;
        cubes.push(Cube(x, y, z));
    }
    if cubes.is_empty() {
        return Err(anyhow!("No cubes in input"));
    }
    Ok(cubes)
}

fn surface_area(cubes: &[Cube], set: &HashSet<Cube>) -> usize {
    let mut area = 0;
    for &cube in cubes {
        for neighbor in cube.neighbors() {
            if !set.contains(&neighbor) {
                area += 1;
            }
        }
    }
    area
}

fn external_surface_area(cubes: &[Cube], set: &HashSet<Cube>) -> usize {
    // Find bounding box
    let (mut min_x, mut max_x) = (i32::MAX, i32::MIN);
    let (mut min_y, mut max_y) = (i32::MAX, i32::MIN);
    let (mut min_z, mut max_z) = (i32::MAX, i32::MIN);

    for &Cube(x, y, z) in cubes {
        min_x = min_x.min(x);
        max_x = max_x.max(x);
        min_y = min_y.min(y);
        max_y = max_y.max(y);
        min_z = min_z.min(z);
        max_z = max_z.max(z);
    }

    // Expand bounds by 1 to ensure we start from outside
    min_x -= 1;
    max_x += 1;
    min_y -= 1;
    max_y += 1;
    min_z -= 1;
    max_z += 1;

    // Flood fill from outside to find all external space
    let mut external = HashSet::new();
    let mut queue = VecDeque::new();
    queue.push_back(Cube(min_x, min_y, min_z));
    external.insert(Cube(min_x, min_y, min_z));

    while let Some(cube) = queue.pop_front() {
        for neighbor in cube.neighbors() {
            let Cube(x, y, z) = neighbor;
            if x >= min_x
                && x <= max_x
                && y >= min_y
                && y <= max_y
                && z >= min_z
                && z <= max_z
                && !set.contains(&neighbor)
                && !external.contains(&neighbor)
            {
                external.insert(neighbor);
                queue.push_back(neighbor);
            }
        }
    }

    // Count faces adjacent to external space
    let mut area = 0;
    for &cube in cubes {
        for neighbor in cube.neighbors() {
            if external.contains(&neighbor) {
                area += 1;
            }
        }
    }
    area
}

pub fn solve(input: &str) -> Result<(usize, usize)> {
    let cubes = parse_input(input)?;
    let set: HashSet<Cube> = cubes.iter().copied().collect();

    let part1 = surface_area(&cubes, &set);
    let part2 = external_surface_area(&cubes, &set);

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day18-example.txt");

    #[test]
    fn test_tiny_example() {
        let input = "1,1,1\n2,1,1";
        let (part1, _) = solve(input).unwrap();
        assert_eq!(part1, 10);
    }

    #[test]
    fn test_example() {
        let (part1, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 64);
        assert_eq!(part2, 58);
    }

    #[test]
    fn test_malformed_input() {
        assert!(solve("1,1").is_err());
        assert!(solve("1,1,a").is_err());
        assert!(solve("").is_err());
    }
}
