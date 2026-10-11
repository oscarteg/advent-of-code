use std::collections::HashSet;

pub fn solve(input: &str) -> anyhow::Result<(usize, usize)> {
    let signal = input.trim();
    anyhow::ensure!(!signal.is_empty(), "signal cannot be empty");

    let part1 = find_marker(signal, 4)?;
    let part2 = find_marker(signal, 14)?;

    Ok((part1, part2))
}

fn find_marker(signal: &str, window_size: usize) -> anyhow::Result<usize> {
    let chars: Vec<char> = signal.chars().collect();
    anyhow::ensure!(
        chars.len() >= window_size,
        "signal length must be at least {window_size}"
    );

    for i in window_size..=chars.len() {
        let window = &chars[i - window_size..i];
        let unique: HashSet<char> = window.iter().copied().collect();
        if unique.len() == window_size {
            return Ok(i);
        }
    }

    anyhow::bail!("no marker found")
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_part1_examples() {
        assert_eq!(find_marker("mjqjpqmgbljsphdztnvjfqwrcgsmlb", 4).unwrap(), 7);
        assert_eq!(find_marker("bvwbjplbgvbhsrlpgdmjqwftvncz", 4).unwrap(), 5);
        assert_eq!(find_marker("nppdvjthqldpwncqszvftbrmjlhg", 4).unwrap(), 6);
        assert_eq!(
            find_marker("nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 4).unwrap(),
            10
        );
        assert_eq!(
            find_marker("zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 4).unwrap(),
            11
        );
    }

    #[test]
    fn test_part2_examples() {
        assert_eq!(
            find_marker("mjqjpqmgbljsphdztnvjfqwrcgsmlb", 14).unwrap(),
            19
        );
        assert_eq!(find_marker("bvwbjplbgvbhsrlpgdmjqwftvncz", 14).unwrap(), 23);
        assert_eq!(find_marker("nppdvjthqldpwncqszvftbrmjlhg", 14).unwrap(), 23);
        assert_eq!(
            find_marker("nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 14).unwrap(),
            29
        );
        assert_eq!(
            find_marker("zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 14).unwrap(),
            26
        );
    }

    #[test]
    fn test_example_from_file() {
        let input = include_str!("../../input/day6-example.txt");
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 7);
        assert_eq!(part2, 19);
    }

    #[test]
    fn test_edge_cases() {
        assert_eq!(find_marker("abcd", 4).unwrap(), 4);
        assert_eq!(find_marker("abcde", 4).unwrap(), 4);
    }

    #[test]
    fn test_malformed_input() {
        assert!(find_marker("ab", 4).is_err());
        assert!(find_marker("", 4).is_err());
    }

    #[test]
    fn test_no_marker_found() {
        assert!(find_marker("aaaaaaaaaa", 4).is_err());
    }
}
