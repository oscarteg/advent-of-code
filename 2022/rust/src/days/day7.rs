use anyhow::{ensure, Context, Result};
use std::collections::HashMap;

const DISK_SPACE: u64 = 70_000_000;
const REQUIRED_UNUSED: u64 = 30_000_000;
const SMALL_DIRECTORY_LIMIT: u64 = 100_000;

/// Walks the `cd`/`ls` terminal transcript and returns the total size of every
/// directory encountered, keyed by its absolute path.
fn directory_sizes(input: &str) -> Result<HashMap<String, u64>> {
    let mut sizes: HashMap<String, u64> = HashMap::new();
    let mut stack: Vec<String> = Vec::new();

    for line in input.lines().map(str::trim).filter(|line| !line.is_empty()) {
        if line == "$ ls" {
            continue;
        }
        if let Some(target) = line.strip_prefix("$ cd ") {
            match target {
                "/" => stack.clear(),
                ".." => {
                    ensure!(stack.pop().is_some(), "cd .. attempted above root");
                }
                name => stack.push(name.to_string()),
            }
            sizes.entry(path_key(&stack)).or_insert(0);
            continue;
        }
        if let Some(name) = line.strip_prefix("dir ") {
            ensure!(!name.is_empty(), "directory entry missing a name");
            continue;
        }

        let (size_str, name) = line
            .split_once(' ')
            .with_context(|| format!("malformed listing line: {line}"))?;
        ensure!(!name.is_empty(), "file entry missing a name");
        let size: u64 = size_str
            .parse()
            .with_context(|| format!("invalid file size: {size_str}"))?;

        for depth in 0..=stack.len() {
            *sizes.entry(path_key(&stack[..depth])).or_insert(0) += size;
        }
    }

    ensure!(!sizes.is_empty(), "transcript contains no directories");
    Ok(sizes)
}

fn path_key(stack: &[String]) -> String {
    if stack.is_empty() {
        "/".to_string()
    } else {
        format!("/{}", stack.join("/"))
    }
}

/// Returns the sum of sizes of directories at most `small_directory_limit`,
/// and the size of the smallest directory that frees enough space once
/// deleted from a disk of `disk_space` bytes with `required_unused` free.
pub fn solve_with(
    input: &str,
    disk_space: u64,
    required_unused: u64,
    small_directory_limit: u64,
) -> Result<(u64, u64)> {
    let sizes = directory_sizes(input)?;
    let part1: u64 = sizes
        .values()
        .filter(|&&size| size <= small_directory_limit)
        .sum();

    let used = *sizes.get("/").context("missing root directory")?;
    ensure!(used <= disk_space, "used space exceeds total disk space");
    let unused = disk_space - used;
    let need_to_free = required_unused.saturating_sub(unused);
    let part2 = sizes
        .values()
        .filter(|&&size| size >= need_to_free)
        .min()
        .copied()
        .context("no directory is large enough to delete")?;

    Ok((part1, part2))
}

pub fn solve(input: &str) -> Result<(u64, u64)> {
    solve_with(input, DISK_SPACE, REQUIRED_UNUSED, SMALL_DIRECTORY_LIMIT)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day7-example.txt");
        assert_eq!(solve(input).unwrap(), (95437, 24933642));
    }

    #[test]
    fn empty_directories_count_as_zero() {
        let input = "$ cd /\n$ ls\ndir a\n100000 big.txt\n$ cd a\n$ ls\n";
        let (part1, _) = solve(input).unwrap();
        assert_eq!(part1, 100_000);
    }

    #[test]
    fn cd_above_root_is_rejected() {
        let input = "$ cd /\n$ cd ..\n";
        assert!(solve(input).is_err());
    }

    #[test]
    fn malformed_listing_line_is_rejected() {
        let input = "$ cd /\n$ ls\nnotanumber\n";
        assert!(solve(input).is_err());
    }
}
