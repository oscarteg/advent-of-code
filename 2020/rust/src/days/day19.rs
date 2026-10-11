use anyhow::{ensure, Context, Result};
use std::collections::HashMap;
use std::str::FromStr;

const ROOT: u32 = 0;

/// End positions reachable per (rule, start position), for a single message.
type Cache = HashMap<(u32, usize), Vec<usize>>;

#[derive(Clone)]
enum Rule {
    Literal(u8),
    Alternatives(Vec<Vec<u32>>),
}

impl FromStr for Rule {
    type Err = anyhow::Error;

    fn from_str(body: &str) -> Result<Self> {
        if let Some(literal) = body
            .strip_prefix('"')
            .and_then(|rest| rest.strip_suffix('"'))
        {
            let bytes = literal.as_bytes();
            ensure!(bytes.len() == 1, "literal must be one character: {body:?}");
            return Ok(Rule::Literal(bytes[0]));
        }
        let alternatives = body
            .split('|')
            .map(|alternative| {
                let ids = alternative
                    .split_whitespace()
                    .map(|id| {
                        id.parse::<u32>()
                            .with_context(|| format!("invalid rule id {id:?}"))
                    })
                    .collect::<Result<Vec<_>>>()?;
                ensure!(!ids.is_empty(), "empty alternative in {body:?}");
                Ok(ids)
            })
            .collect::<Result<Vec<_>>>()?;
        Ok(Rule::Alternatives(alternatives))
    }
}

struct Grammar {
    rules: HashMap<u32, Rule>,
}

impl Grammar {
    fn parse(lines: &[&str]) -> Result<Self> {
        let mut rules = HashMap::new();
        for line in lines {
            let (id, body) = line
                .split_once(':')
                .with_context(|| format!("rule without ':': {line:?}"))?;
            let id: u32 = id
                .trim()
                .parse()
                .with_context(|| format!("invalid rule id in {line:?}"))?;
            let rule: Rule = body
                .trim()
                .parse()
                .with_context(|| format!("invalid rule {id}"))?;
            ensure!(!rules.contains_key(&id), "duplicate rule {id}");
            rules.insert(id, rule);
        }
        let grammar = Grammar { rules };
        grammar.validate()?;
        Ok(grammar)
    }

    /// Matching resolves rule ids without failing, so every reference must exist.
    fn validate(&self) -> Result<()> {
        ensure!(self.rules.contains_key(&ROOT), "rule {ROOT} is missing");
        for (id, rule) in &self.rules {
            let Rule::Alternatives(alternatives) = rule else {
                continue;
            };
            for referenced in alternatives.iter().flatten() {
                ensure!(
                    self.rules.contains_key(referenced),
                    "rule {id} references unknown rule {referenced}"
                );
            }
        }
        Ok(())
    }

    /// Part 2 makes rules 8 and 11 recursive. A rule set without them (the part 1
    /// example) is left alone, so its answer stays the same.
    fn with_loops(&self) -> Result<Self> {
        let mut rules = self.rules.clone();
        for (id, body) in [(8, "42 | 42 8"), (11, "42 31 | 42 11 31")] {
            if let std::collections::hash_map::Entry::Occupied(mut e) = rules.entry(id) {
                e.insert(body.parse()?);
            }
        }
        let grammar = Grammar { rules };
        grammar.validate()?;
        Ok(grammar)
    }

    fn count_matching(&self, messages: &[&str]) -> usize {
        messages
            .iter()
            .filter(|message| self.matches(message))
            .count()
    }

    fn matches(&self, message: &str) -> bool {
        let message = message.as_bytes();
        let mut cache = Cache::new();
        self.rule_ends(ROOT, message, 0, &mut cache)
            .contains(&message.len())
    }

    /// Every position the rule can end at when it starts at `start`. Alternatives of
    /// different lengths make this a set instead of a single position, which is what
    /// makes the recursive part 2 rules work without special cases.
    fn rule_ends(&self, id: u32, message: &[u8], start: usize, cache: &mut Cache) -> Vec<usize> {
        if let Some(ends) = cache.get(&(id, start)) {
            return ends.clone();
        }
        let ends = match self.rules.get(&id) {
            None => Vec::new(),
            Some(Rule::Literal(byte)) => match message.get(start) {
                Some(found) if found == byte => vec![start + 1],
                _ => Vec::new(),
            },
            Some(Rule::Alternatives(alternatives)) => {
                let mut ends = Vec::new();
                for alternative in alternatives {
                    ends.extend(self.sequence_ends(alternative, message, start, cache));
                }
                sorted_unique(ends)
            }
        };
        cache.insert((id, start), ends.clone());
        ends
    }

    fn sequence_ends(
        &self,
        ids: &[u32],
        message: &[u8],
        start: usize,
        cache: &mut Cache,
    ) -> Vec<usize> {
        let mut positions = vec![start];
        for &id in ids {
            let mut next = Vec::new();
            for &position in &positions {
                next.extend(self.rule_ends(id, message, position, cache));
            }
            if next.is_empty() {
                return Vec::new();
            }
            positions = sorted_unique(next);
        }
        positions
    }
}

fn sorted_unique(mut positions: Vec<usize>) -> Vec<usize> {
    positions.sort_unstable();
    positions.dedup();
    positions
}

pub fn solve(input: &str) -> Result<(usize, usize)> {
    let mut lines = input.lines().map(str::trim);
    let rules: Vec<&str> = lines.by_ref().take_while(|line| !line.is_empty()).collect();
    let messages: Vec<&str> = lines.filter(|line| !line.is_empty()).collect();
    ensure!(!messages.is_empty(), "no messages after the rules");

    let grammar = Grammar::parse(&rules)?;
    let part1 = grammar.count_matching(&messages);
    let part2 = grammar.with_loops()?.count_matching(&messages);
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/2020/day19-example.txt");

    /// The larger official example; its rules 8 and 11 loop in part 2.
    const LOOP_EXAMPLE: &str = r#"42: 9 14 | 10 1
9: 14 27 | 1 26
10: 23 14 | 28 1
1: "a"
11: 42 31
5: 1 14 | 15 1
19: 14 1 | 14 14
12: 24 14 | 19 1
16: 15 1 | 14 14
31: 14 17 | 1 13
6: 14 14 | 1 14
2: 1 24 | 14 4
0: 8 11
13: 14 3 | 1 12
15: 1 | 14
17: 14 2 | 1 7
23: 25 1 | 22 14
28: 16 1
4: 1 1
20: 14 14 | 1 15
3: 5 14 | 16 1
27: 1 6 | 14 18
14: "b"
21: 14 1 | 1 14
25: 1 1 | 1 14
22: 14 14
8: 42
26: 14 22 | 1 20
18: 15 15
7: 14 5 | 1 21
24: 14 1

abbbbbabbbaaaababbaabbbbabababbbabbbbbbabaaaa
bbabbbbaabaabba
babbbbaabbbbbabbbbbbaabaaabaaa
aaabbbbbbaaaabaababaabababbabaaabbababababaaa
bbbbbbbaaaabbbbaaabbabaaa
bbbababbbbaaaaaaaabbababaaababaabab
ababaaaaaabaaab
ababaaaaabbbaba
baabbaaaabbaaaababbaababb
abbbbabbbbaaaababbbbbbaaaababb
aaaaabbaabaaaaababaa
aaaabbaaaabbaaa
aaaabbaabbaaaaaaabbbabbbaaabbaabaab
babaaabbbaaabaababbaabababaaab
aabbbbbaabbbaaaaaabbbbbababaaaaabbaaabba
"#;

    #[test]
    fn part1_example() {
        assert_eq!(solve(EXAMPLE).unwrap().0, 2);
    }

    #[test]
    fn part2_leaves_a_rule_set_without_loops_alone() {
        assert_eq!(solve(EXAMPLE).unwrap().1, 2);
    }

    #[test]
    fn loop_example_counts_both_parts() {
        assert_eq!(solve(LOOP_EXAMPLE).unwrap(), (3, 12));
    }

    #[test]
    fn alternatives_of_different_lengths_both_count() {
        let input = "0: 1 | 1 1\n1: \"a\"\n\na\naa\naaa\n";
        assert_eq!(solve(input).unwrap(), (2, 2));
    }

    #[test]
    fn leftover_characters_do_not_match() {
        let input = "0: 1 1\n1: \"a\"\n\naa\naaa\n";
        assert_eq!(solve(input).unwrap().0, 1);
    }

    #[test]
    fn loop_rules_without_the_rules_they_need_are_rejected() {
        let input = "0: 8 11\n8: 9\n11: 9 9\n9: \"a\"\n\naaa\n";
        assert!(solve(input).is_err());
    }

    #[test]
    fn malformed_input_is_rejected() {
        for input in [
            "",
            "0: 4 1 5\n\nab",
            "0: \"ab\"\n\nab",
            "0 1\n1: \"a\"\n\na",
            "x: \"a\"\n\na",
            "0: 1\n1: \"a\"\n1: \"b\"\n\na",
            "1: \"a\"\n\na",
            "0: \"a\"\n",
            "0: | 1\n1: \"a\"\n\na",
        ] {
            assert!(solve(input).is_err(), "expected an error for {input:?}");
        }
    }
}
