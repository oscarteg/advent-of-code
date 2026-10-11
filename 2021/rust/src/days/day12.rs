use std::collections::{HashMap, HashSet};

pub fn solve(input: &str) -> anyhow::Result<(u32, u32)> {
    let graph = parse_graph(input)?;
    let part1 = count_paths(&graph, false);
    let part2 = count_paths(&graph, true);
    Ok((part1, part2))
}

fn parse_graph(input: &str) -> anyhow::Result<HashMap<String, Vec<String>>> {
    let mut graph: HashMap<String, Vec<String>> = HashMap::new();

    for line in input.lines() {
        let line = line.trim();
        if line.is_empty() {
            continue;
        }

        let parts: Vec<&str> = line.split('-').collect();
        anyhow::ensure!(
            parts.len() == 2,
            "expected edge format 'a-b', got '{}'",
            line
        );
        let (a, b) = (parts[0].to_string(), parts[1].to_string());

        graph.entry(a.clone()).or_default().push(b.clone());
        graph.entry(b).or_default().push(a);
    }

    Ok(graph)
}

fn is_small_cave(cave: &str) -> bool {
    cave.chars().all(|c| c.is_lowercase())
}

fn count_paths(graph: &HashMap<String, Vec<String>>, allow_double: bool) -> u32 {
    dfs(graph, "start", HashSet::new(), allow_double, allow_double)
}

fn dfs(
    graph: &HashMap<String, Vec<String>>,
    current: &str,
    mut visited: HashSet<String>,
    allow_double: bool,
    double_available: bool,
) -> u32 {
    if current == "end" {
        return 1;
    }

    if is_small_cave(current) {
        visited.insert(current.to_string());
    }

    let neighbors = graph.get(current).cloned().unwrap_or_default();
    let mut count = 0;

    for neighbor in neighbors {
        let already_visited = visited.contains(&neighbor);

        if !already_visited {
            count += dfs(
                graph,
                &neighbor,
                visited.clone(),
                allow_double,
                double_available,
            );
        } else if allow_double && double_available && neighbor != "start" {
            count += dfs(graph, &neighbor, visited.clone(), allow_double, false);
        }
    }

    count
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_small_example() {
        let input = include_str!("../../input/day12-example.txt");
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 10);
        assert_eq!(part2, 36);
    }

    #[test]
    fn test_medium_example() {
        let input = "dc-end
HN-start
start-kj
dc-start
dc-HN
LN-dc
HN-end
kj-sa
kj-HN
kj-dc";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 19);
        assert_eq!(part2, 103);
    }

    #[test]
    fn test_large_example() {
        let input = "fs-end
he-DX
fs-he
start-DX
pj-DX
end-zg
zg-sl
zg-pj
pj-he
RW-he
fs-DX
pj-RW
zg-RW
start-pj
he-WI
zg-he
pj-fs
start-RW";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 226);
        assert_eq!(part2, 3509);
    }

    #[test]
    fn test_single_path() {
        let input = "start-A\nA-end";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 1);
        assert_eq!(part2, 1);
    }

    #[test]
    fn test_invalid_edge_format() {
        let input = "start-A-B";
        assert!(solve(input).is_err());
    }

    #[test]
    fn test_invalid_edge_no_hyphen() {
        let input = "startA";
        assert!(solve(input).is_err());
    }
}
