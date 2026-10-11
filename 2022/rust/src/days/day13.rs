use anyhow::{bail, Context, Result};
use std::cmp::Ordering;

#[derive(Debug, Clone, PartialEq, Eq)]
enum Packet {
    Int(i32),
    List(Vec<Packet>),
}

impl Packet {
    fn parse(s: &str) -> Result<Self> {
        let s = s.trim();
        if s.is_empty() {
            bail!("empty packet");
        }
        let (packet, rest) = Self::parse_packet(s)?;
        if !rest.is_empty() {
            bail!("trailing content after packet");
        }
        Ok(packet)
    }

    fn parse_packet(s: &str) -> Result<(Self, &str)> {
        if s.starts_with('[') {
            Self::parse_list(s)
        } else if s
            .chars()
            .next()
            .is_some_and(|c| c.is_ascii_digit() || c == '-')
        {
            Self::parse_int(s)
        } else {
            bail!("invalid packet start: {}", s.chars().next().unwrap_or('?'))
        }
    }

    fn parse_list(s: &str) -> Result<(Self, &str)> {
        if !s.starts_with('[') {
            bail!("expected '[' at start of list");
        }
        let mut s = &s[1..];
        let mut items = Vec::new();

        if let Some(stripped) = s.strip_prefix(']') {
            return Ok((Packet::List(items), stripped));
        }

        loop {
            let (packet, rest) = Self::parse_packet(s)?;
            items.push(packet);
            s = rest.trim_start();

            if s.starts_with(',') {
                s = s[1..].trim_start();
            } else if let Some(stripped) = s.strip_prefix(']') {
                return Ok((Packet::List(items), stripped));
            } else {
                bail!(
                    "expected ',' or ']' in list, got: {}",
                    s.chars().next().unwrap_or('?')
                );
            }
        }
    }

    fn parse_int(s: &str) -> Result<(Self, &str)> {
        let mut end = 0;
        let bytes = s.as_bytes();
        if end < bytes.len() && bytes[end] == b'-' {
            end += 1;
        }
        while end < bytes.len() && bytes[end].is_ascii_digit() {
            end += 1;
        }
        if end == 0 || (end == 1 && bytes[0] == b'-') {
            bail!("invalid integer");
        }
        let num_str = &s[..end];
        let num = num_str.parse::<i32>().context("number out of range")?;
        Ok((Packet::Int(num), &s[end..]))
    }

    fn compare(&self, other: &Packet) -> Ordering {
        match (self, other) {
            (Packet::Int(a), Packet::Int(b)) => a.cmp(b),
            (Packet::List(a), Packet::List(b)) => {
                for i in 0..a.len().min(b.len()) {
                    match a[i].compare(&b[i]) {
                        Ordering::Equal => continue,
                        other => return other,
                    }
                }
                a.len().cmp(&b.len())
            }
            (Packet::Int(a), Packet::List(b)) => {
                let a_list = Packet::List(vec![Packet::Int(*a)]);
                a_list.compare(&Packet::List(b.clone()))
            }
            (Packet::List(a), Packet::Int(b)) => {
                let b_list = Packet::List(vec![Packet::Int(*b)]);
                Packet::List(a.clone()).compare(&b_list)
            }
        }
    }
}

impl Ord for Packet {
    fn cmp(&self, other: &Self) -> Ordering {
        self.compare(other)
    }
}

impl PartialOrd for Packet {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

pub fn solve(input: &str) -> Result<(usize, usize)> {
    let part1 = solve_part1(input)?;
    let part2 = solve_part2(input)?;
    Ok((part1, part2))
}

fn solve_part1(input: &str) -> Result<usize> {
    let pairs: Vec<&str> = input.trim().split("\n\n").collect();
    let mut sum = 0;

    for (idx, pair) in pairs.iter().enumerate() {
        let lines: Vec<&str> = pair.trim().lines().collect();
        if lines.len() != 2 {
            bail!("expected 2 packets per pair, got {}", lines.len());
        }

        let left = Packet::parse(lines[0])?;
        let right = Packet::parse(lines[1])?;

        if left.cmp(&right) == Ordering::Less {
            sum += idx + 1;
        }
    }

    Ok(sum)
}

fn solve_part2(input: &str) -> Result<usize> {
    let mut packets = Vec::new();

    for line in input.trim().lines() {
        let line = line.trim();
        if !line.is_empty() {
            packets.push(Packet::parse(line)?);
        }
    }

    let divider1 = Packet::List(vec![Packet::List(vec![Packet::Int(2)])]);
    let divider2 = Packet::List(vec![Packet::List(vec![Packet::Int(6)])]);

    packets.push(divider1.clone());
    packets.push(divider2.clone());

    packets.sort();

    let pos1 = packets
        .iter()
        .position(|p| p == &divider1)
        .context("divider 1 not found")?;
    let pos2 = packets
        .iter()
        .position(|p| p == &divider2)
        .context("divider 2 not found")?;

    Ok((pos1 + 1) * (pos2 + 1))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_part1_example() {
        let input = include_str!("../../input/day13-example.txt");
        assert_eq!(solve_part1(input).unwrap(), 13);
    }

    #[test]
    fn test_part2_example() {
        let input = include_str!("../../input/day13-example.txt");
        assert_eq!(solve_part2(input).unwrap(), 140);
    }

    #[test]
    fn test_parse_int() {
        let (packet, rest) = Packet::parse_int("42abc").unwrap();
        assert_eq!(packet, Packet::Int(42));
        assert_eq!(rest, "abc");
    }

    #[test]
    fn test_parse_list() {
        let (packet, rest) = Packet::parse_list("[]rest").unwrap();
        assert_eq!(packet, Packet::List(vec![]));
        assert_eq!(rest, "rest");
    }

    #[test]
    fn test_parse_complex() {
        let packet = Packet::parse("[1,[2,3],4]").unwrap();
        match packet {
            Packet::List(ref items) => {
                assert_eq!(items.len(), 3);
            }
            _ => panic!("expected list"),
        }
    }

    #[test]
    fn test_compare_ints() {
        let a = Packet::Int(1);
        let b = Packet::Int(2);
        assert_eq!(a.cmp(&b), Ordering::Less);
    }

    #[test]
    fn test_compare_int_vs_list() {
        let a = Packet::Int(2);
        let b = Packet::List(vec![Packet::Int(3)]);
        assert_eq!(a.cmp(&b), Ordering::Less);
    }

    #[test]
    fn test_compare_lists() {
        let a = Packet::List(vec![Packet::Int(1)]);
        let b = Packet::List(vec![Packet::Int(2)]);
        assert_eq!(a.cmp(&b), Ordering::Less);
    }

    #[test]
    fn test_malformed_empty() {
        assert!(Packet::parse("").is_err());
    }

    #[test]
    fn test_malformed_invalid_start() {
        assert!(Packet::parse("abc").is_err());
    }

    #[test]
    fn test_malformed_unclosed_list() {
        assert!(Packet::parse("[1,2").is_err());
    }
}
