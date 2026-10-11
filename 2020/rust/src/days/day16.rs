use anyhow::{ensure, Context, Result};
use std::ops::RangeInclusive;

/// Fields whose values form the part 2 product.
const DEPARTURE_PREFIX: &str = "departure";

struct Rule {
    name: String,
    ranges: Vec<RangeInclusive<u64>>,
}

impl Rule {
    fn accepts(&self, value: u64) -> bool {
        self.ranges.iter().any(|range| range.contains(&value))
    }
}

/// The three sections of the notes: the field rules, your own ticket and the nearby tickets.
struct Notes {
    rules: Vec<Rule>,
    ticket: Vec<u64>,
    nearby: Vec<Vec<u64>>,
}

impl Notes {
    fn accepted_by_any_rule(&self, value: u64) -> bool {
        self.rules.iter().any(|rule| rule.accepts(value))
    }

    /// Sums every nearby value that no rule accepts.
    fn error_rate(&self) -> u64 {
        self.nearby
            .iter()
            .flatten()
            .filter(|&&value| !self.accepted_by_any_rule(value))
            .sum()
    }

    /// Matches every ticket column to a field name.
    ///
    /// A column starts out with every field that accepts its value on all valid tickets. The
    /// column left with a single candidate is assigned and that field is struck from the other
    /// columns, which is repeated until each column has its field.
    fn field_order(&self) -> Result<Vec<&str>> {
        let valid: Vec<&[u64]> = self
            .nearby
            .iter()
            .filter(|ticket| ticket.iter().all(|&value| self.accepted_by_any_rule(value)))
            .map(Vec::as_slice)
            .collect();
        let mut candidates: Vec<Vec<usize>> = (0..self.rules.len())
            .map(|column| {
                (0..self.rules.len())
                    .filter(|&field| {
                        valid
                            .iter()
                            .all(|ticket| self.rules[field].accepts(ticket[column]))
                    })
                    .collect()
            })
            .collect();

        let mut order = vec![""; self.rules.len()];
        for _ in 0..self.rules.len() {
            let column = candidates
                .iter()
                .position(|fields| fields.len() == 1)
                .context("the notes do not pin down a single field order")?;
            let field = candidates[column][0];
            order[column] = self.rules[field].name.as_str();
            for fields in &mut candidates {
                fields.retain(|&candidate| candidate != field);
            }
        }
        Ok(order)
    }
}

fn parse_range(text: &str) -> Result<RangeInclusive<u64>> {
    let (start, end) = text
        .trim()
        .split_once('-')
        .with_context(|| format!("invalid range: {text}"))?;
    let start: u64 = start
        .parse()
        .with_context(|| format!("invalid range start: {start}"))?;
    let end: u64 = end
        .parse()
        .with_context(|| format!("invalid range end: {end}"))?;
    ensure!(start <= end, "range {start}-{end} is empty");
    Ok(start..=end)
}

fn parse_rule(line: &str) -> Result<Rule> {
    let (name, ranges) = line
        .trim()
        .split_once(": ")
        .with_context(|| format!("invalid field rule: {line}"))?;
    let ranges = ranges
        .split(" or ")
        .map(parse_range)
        .collect::<Result<_>>()?;
    Ok(Rule {
        name: name.to_owned(),
        ranges,
    })
}

fn parse_ticket(line: &str) -> Result<Vec<u64>> {
    line.split(',')
        .map(|value| {
            value
                .trim()
                .parse()
                .with_context(|| format!("invalid ticket value: {value}"))
        })
        .collect()
}

fn parse_tickets(section: &str, header: &str) -> Result<Vec<Vec<u64>>> {
    let mut lines = section
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty());
    let found = lines
        .next()
        .with_context(|| format!("missing the {header} header"))?;
    ensure!(
        found == header,
        "expected the {header} header, found {found}"
    );
    lines.map(parse_ticket).collect()
}

fn parse_notes(input: &str) -> Result<Notes> {
    let mut sections = input.trim().split("\n\n");
    let rule_section = sections.next().context("missing the field rules")?;
    let own_section = sections.next().context("missing your ticket")?;
    let nearby_section = sections.next().context("missing the nearby tickets")?;
    ensure!(sections.next().is_none(), "unexpected fourth section");

    let rules: Vec<Rule> = rule_section
        .lines()
        .map(parse_rule)
        .collect::<Result<_>>()?;
    ensure!(!rules.is_empty(), "missing the field rules");

    let own = parse_tickets(own_section, "your ticket:")?;
    ensure!(
        own.len() == 1,
        "expected exactly one ticket of your own, found {}",
        own.len()
    );
    let ticket = own.into_iter().next().context("missing your ticket")?;

    let nearby = parse_tickets(nearby_section, "nearby tickets:")?;
    ensure!(!nearby.is_empty(), "missing the nearby tickets");

    for ticket in std::iter::once(&ticket).chain(&nearby) {
        ensure!(
            ticket.len() == rules.len(),
            "ticket has {} values but there are {} fields",
            ticket.len(),
            rules.len()
        );
    }

    Ok(Notes {
        rules,
        ticket,
        nearby,
    })
}

pub fn solve(input: &str) -> Result<(u64, u64)> {
    let notes = parse_notes(input)?;
    let order = notes.field_order()?;
    // Notes without a departure field (such as the examples) leave an empty product of 1.
    let departure_product: u64 = order
        .iter()
        .zip(&notes.ticket)
        .filter(|(name, _)| name.starts_with(DEPARTURE_PREFIX))
        .map(|(_, &value)| value)
        .product();
    Ok((notes.error_rate(), departure_product))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/2020/day16-example.txt");

    /// The second example of the puzzle: it has no departure fields, but its field order is
    /// fully determined.
    const FIELD_ORDER_EXAMPLE: &str = "\
class: 0-1 or 4-19
row: 0-5 or 8-19
seat: 0-13 or 16-19

your ticket:
11,12,13

nearby tickets:
3,9,18
15,1,5
5,14,9
";

    #[test]
    fn official_example() {
        assert_eq!(solve(EXAMPLE).unwrap(), (71, 1));
    }

    #[test]
    fn second_example_resolves_the_field_order() {
        let notes = parse_notes(FIELD_ORDER_EXAMPLE).unwrap();
        let fields: Vec<(&str, u64)> = notes
            .field_order()
            .unwrap()
            .into_iter()
            .zip(notes.ticket.iter().copied())
            .collect();
        assert_eq!(fields, [("row", 11), ("class", 12), ("seat", 13)]);
    }

    #[test]
    fn departure_fields_are_multiplied() {
        let input = "\
departure gate: 1-3
departure track: 1-6
seat: 1-9

your ticket:
2,5,7

nearby tickets:
1,4,9
3,6,8
";
        assert_eq!(solve(input).unwrap(), (0, 2 * 5));
    }

    #[test]
    fn ambiguous_field_order_is_rejected() {
        let input = "\
first: 1-2
second: 1-2

your ticket:
1,2

nearby tickets:
1,2
";
        assert!(solve(input).is_err());
    }

    #[test]
    fn malformed_input_is_rejected() {
        let valid = "class: 1-3\n\nyour ticket:\n1\n\nnearby tickets:\n2\n";
        assert!(solve(valid).is_ok());
        assert!(solve("").is_err());
        assert!(solve("class: 1-3\n\nyour ticket:\n1\n").is_err());
        assert!(solve("class 1-3\n\nyour ticket:\n1\n\nnearby tickets:\n2\n").is_err());
        assert!(solve("class: 5-1\n\nyour ticket:\n1\n\nnearby tickets:\n2\n").is_err());
        assert!(solve("class: 1-3\n\nthe ticket:\n1\n\nnearby tickets:\n2\n").is_err());
        assert!(solve("class: 1-3\n\nyour ticket:\n1\n\nnearby tickets:\nx\n").is_err());
        assert!(solve("class: 1-3\n\nyour ticket:\n1,2\n\nnearby tickets:\n2,3\n").is_err());
    }
}
