use crate::day19::custom_error::AocError;
use std::collections::HashMap;

/// Every rating lies between these bounds, inclusive.
const MIN_RATING: u64 = 1;
const MAX_RATING: u64 = 4000;
/// Parts enter the system through this workflow.
const START: &str = "in";

#[derive(Debug, Clone, Copy)]
enum Category {
    X,
    M,
    A,
    S,
}

impl Category {
    fn parse(name: &str) -> Result<Self, AocError> {
        match name {
            "x" => Ok(Self::X),
            "m" => Ok(Self::M),
            "a" => Ok(Self::A),
            "s" => Ok(Self::S),
            _ => Err(AocError::InvalidInput("unknown rating category")),
        }
    }

    fn index(self) -> usize {
        match self {
            Self::X => 0,
            Self::M => 1,
            Self::A => 2,
            Self::S => 3,
        }
    }
}

/// The four ratings of a single machine part, indexed by [`Category`].
#[derive(Debug, Clone, Copy)]
pub(crate) struct Part([u64; 4]);

impl Part {
    fn parse(line: &str) -> Result<Self, AocError> {
        let body = line
            .strip_prefix('{')
            .and_then(|rest| rest.strip_suffix('}'))
            .ok_or(AocError::InvalidInput("part must be wrapped in braces"))?;
        let mut ratings = [None; 4];
        for assignment in body.split(',') {
            let (name, value) = assignment
                .split_once('=')
                .ok_or(AocError::InvalidInput("rating needs category=value"))?;
            let slot = &mut ratings[Category::parse(name.trim())?.index()];
            if slot.replace(parse_rating(value)?).is_some() {
                return Err(AocError::InvalidInput("duplicate rating category"));
            }
        }
        let [Some(x), Some(m), Some(a), Some(s)] = ratings else {
            return Err(AocError::InvalidInput("part misses a rating category"));
        };
        Ok(Self([x, m, a, s]))
    }

    pub fn total_rating(self) -> u64 {
        self.0.iter().sum()
    }

    fn rating(self, category: Category) -> u64 {
        self.0[category.index()]
    }
}

/// Ratings are parsed as [`u32`] so that range arithmetic cannot overflow.
fn parse_rating(text: &str) -> Result<u64, AocError> {
    Ok(u64::from(text.trim().parse::<u32>()?))
}

#[derive(Debug, Clone, Copy)]
enum Comparison {
    Less,
    Greater,
}

#[derive(Debug)]
struct Condition {
    category: Category,
    comparison: Comparison,
    threshold: u64,
}

impl Condition {
    fn parse(text: &str) -> Result<Self, AocError> {
        let operator = text
            .find(['<', '>'])
            .ok_or(AocError::InvalidInput("condition needs < or >"))?;
        let comparison = if text[operator..].starts_with('<') {
            Comparison::Less
        } else {
            Comparison::Greater
        };
        Ok(Self {
            category: Category::parse(text[..operator].trim())?,
            comparison,
            threshold: parse_rating(&text[operator + 1..])?,
        })
    }

    fn matches(&self, part: Part) -> bool {
        let rating = part.rating(self.category);
        match self.comparison {
            Comparison::Less => rating < self.threshold,
            Comparison::Greater => rating > self.threshold,
        }
    }

    /// Splits `space` into the ratings satisfying the condition and the rest.
    fn split(&self, space: Ratings) -> (Ratings, Ratings) {
        match self.comparison {
            Comparison::Less => (
                space.clamp(self.category, MIN_RATING, self.threshold - 1),
                space.clamp(self.category, self.threshold, MAX_RATING),
            ),
            Comparison::Greater => (
                space.clamp(self.category, self.threshold + 1, MAX_RATING),
                space.clamp(self.category, MIN_RATING, self.threshold),
            ),
        }
    }
}

/// An inclusive range of still-possible ratings; empty when `high < low`.
#[derive(Debug, Clone, Copy)]
struct Span {
    low: u64,
    high: u64,
}

impl Span {
    fn len(self) -> u64 {
        if self.high < self.low {
            0
        } else {
            self.high - self.low + 1
        }
    }
}

/// A search space of part ratings: one span per category.
#[derive(Debug, Clone, Copy)]
struct Ratings([Span; 4]);

impl Ratings {
    fn full() -> Self {
        Self(
            [Span {
                low: MIN_RATING,
                high: MAX_RATING,
            }; 4],
        )
    }

    fn combinations(self) -> u64 {
        self.0.iter().map(|span| span.len()).product()
    }

    fn clamp(self, category: Category, low: u64, high: u64) -> Self {
        let mut clamped = self;
        let span = &mut clamped.0[category.index()];
        span.low = span.low.max(low);
        span.high = span.high.min(high);
        clamped
    }
}

#[derive(Debug, Clone, Copy)]
enum Target<'a> {
    Accept,
    Reject,
    Workflow(&'a str),
}

impl<'a> Target<'a> {
    fn parse(name: &'a str) -> Result<Self, AocError> {
        match name {
            "" => Err(AocError::InvalidInput("rule has no target")),
            "A" => Ok(Self::Accept),
            "R" => Ok(Self::Reject),
            _ => Ok(Self::Workflow(name)),
        }
    }
}

/// A rule routes a part to `target`. Only the last rule of a workflow is
/// unconditional, so it always applies once reached.
#[derive(Debug)]
struct Rule<'a> {
    condition: Option<Condition>,
    target: Target<'a>,
}

impl<'a> Rule<'a> {
    fn parse(text: &'a str) -> Result<Self, AocError> {
        match text.split_once(':') {
            Some((condition, target)) => Ok(Self {
                condition: Some(Condition::parse(condition)?),
                target: Target::parse(target.trim())?,
            }),
            None => Ok(Self {
                condition: None,
                target: Target::parse(text)?,
            }),
        }
    }
}

fn parse_workflow(line: &str) -> Result<(&str, Vec<Rule<'_>>), AocError> {
    let (name, body) = line
        .split_once('{')
        .ok_or(AocError::InvalidInput("workflow needs a rule list"))?;
    if name.is_empty() {
        return Err(AocError::InvalidInput("workflow has no name"));
    }
    let rules = body
        .strip_suffix('}')
        .ok_or(AocError::InvalidInput("workflow rule list is unclosed"))?
        .split(',')
        .map(Rule::parse)
        .collect::<Result<Vec<_>, _>>()?;
    let (last, routed) = rules
        .split_last()
        .ok_or(AocError::InvalidInput("workflow has no rules"))?;
    if last.condition.is_some() || routed.iter().any(|rule| rule.condition.is_none()) {
        return Err(AocError::InvalidInput(
            "only the last rule may be unconditional",
        ));
    }
    Ok((name, rules))
}

#[derive(Debug)]
pub(crate) struct System<'a> {
    workflows: HashMap<&'a str, Vec<Rule<'a>>>,
    pub parts: Vec<Part>,
}

impl<'a> System<'a> {
    pub fn parse(input: &'a str) -> Result<Self, AocError> {
        let mut lines = input
            .lines()
            .map(str::trim)
            .skip_while(|line| line.is_empty());
        let mut workflows = HashMap::new();
        for line in lines.by_ref() {
            if line.is_empty() {
                break;
            }
            let (name, rules) = parse_workflow(line)?;
            if workflows.insert(name, rules).is_some() {
                return Err(AocError::InvalidInput("duplicate workflow name"));
            }
        }
        if !workflows.contains_key(START) {
            return Err(AocError::InvalidInput("missing workflow in"));
        }
        let parts = lines
            .filter(|line| !line.is_empty())
            .map(Part::parse)
            .collect::<Result<Vec<_>, _>>()?;
        if parts.is_empty() {
            return Err(AocError::InvalidInput("missing part ratings"));
        }
        let system = Self { workflows, parts };
        system.validate_targets()?;
        Ok(system)
    }

    /// Walks a part through the workflows until it is accepted or rejected.
    pub fn accepts(&self, part: Part) -> Result<bool, AocError> {
        let mut target = Target::Workflow(START);
        // Visiting more workflows than exist means the rules form a cycle.
        for _ in 0..=self.workflows.len() {
            let name = match target {
                Target::Accept => return Ok(true),
                Target::Reject => return Ok(false),
                Target::Workflow(name) => name,
            };
            target = self
                .rules(name)?
                .iter()
                .find(|rule| rule.condition.as_ref().is_none_or(|c| c.matches(part)))
                .ok_or(AocError::InvalidInput("workflow has no matching rule"))?
                .target;
        }
        Err(AocError::InvalidInput("workflows form a cycle"))
    }

    /// Counts the rating combinations the workflows accept.
    pub fn accepted_combinations(&self) -> Result<u64, AocError> {
        self.count_accepted(Target::Workflow(START), Ratings::full(), 0)
    }

    /// Pushes a whole search space through `target`. Every condition splits the
    /// space into the ratings it routes away and the ratings the next rule
    /// sees, so the accepted fragments are disjoint and can simply be summed.
    fn count_accepted(
        &self,
        target: Target<'_>,
        space: Ratings,
        depth: usize,
    ) -> Result<u64, AocError> {
        let combinations = space.combinations();
        if combinations == 0 {
            return Ok(0);
        }
        let name = match target {
            Target::Accept => return Ok(combinations),
            Target::Reject => return Ok(0),
            Target::Workflow(name) => name,
        };
        if depth > self.workflows.len() {
            return Err(AocError::InvalidInput("workflows form a cycle"));
        }
        let mut accepted = 0;
        let mut remaining = space;
        for rule in self.rules(name)? {
            let Some(condition) = &rule.condition else {
                return Ok(accepted + self.count_accepted(rule.target, remaining, depth + 1)?);
            };
            let (matched, rest) = condition.split(remaining);
            accepted += self.count_accepted(rule.target, matched, depth + 1)?;
            remaining = rest;
        }
        Ok(accepted)
    }

    fn rules(&self, name: &str) -> Result<&[Rule<'a>], AocError> {
        self.workflows
            .get(name)
            .map(Vec::as_slice)
            .ok_or(AocError::InvalidInput("unknown workflow"))
    }

    fn validate_targets(&self) -> Result<(), AocError> {
        let unknown = self
            .workflows
            .values()
            .flatten()
            .filter_map(|rule| match rule.target {
                Target::Workflow(name) => Some(name),
                Target::Accept | Target::Reject => None,
            })
            .any(|name| !self.workflows.contains_key(name));
        if unknown {
            return Err(AocError::InvalidInput("rule targets unknown workflow"));
        }
        Ok(())
    }
}
