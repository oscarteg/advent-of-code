use crate::day25::custom_error::AocError;
use std::collections::HashMap;

/// A wiring diagram: components connected by undirected wires, stored as paired
/// unit-capacity arcs so that a maximum flow run yields a minimum cut. Arc `i`
/// and arc `i ^ 1` are the two directions of the same wire; `outgoing[c]` holds
/// the arcs leaving component `c` and `heads[a]` the component arc `a` enters.
#[derive(Debug)]
pub(crate) struct Diagram {
    outgoing: Vec<Vec<usize>>,
    heads: Vec<usize>,
}

impl Diagram {
    pub fn parse(input: &str) -> Result<Self, AocError> {
        let mut ids = HashMap::new();
        let mut diagram = Self {
            outgoing: Vec::new(),
            heads: Vec::new(),
        };
        for line in input.lines().map(str::trim).filter(|line| !line.is_empty()) {
            let (name, connections) = line
                .split_once(':')
                .ok_or(AocError::InvalidInput("connection line without ':'"))?;
            let component = diagram.component_id(&mut ids, name.trim())?;
            let mut wires = 0;
            for neighbour in connections.split_whitespace() {
                let neighbour = diagram.component_id(&mut ids, neighbour)?;
                if neighbour == component {
                    return Err(AocError::InvalidInput("component wired to itself"));
                }
                diagram.connect(component, neighbour);
                wires += 1;
            }
            if wires == 0 {
                return Err(AocError::InvalidInput("component without connections"));
            }
        }
        if diagram.components() < 2 {
            return Err(AocError::InvalidInput("fewer than two components"));
        }
        Ok(diagram)
    }

    pub fn components(&self) -> usize {
        self.outgoing.len()
    }

    /// Sizes of the two groups left by the cheapest set of wires that splits the
    /// diagram. Keeping the cheapest maximum flow from one fixed component to
    /// every other one yields the global minimum cut, because the fixed
    /// component always lies on one of its two sides.
    pub fn minimum_cut(&self) -> (usize, usize) {
        let mut cheapest = usize::MAX;
        let mut source_side = 1;
        for sink in 1..self.components() {
            if let Some((wires, reachable)) = self.flow_under(sink, cheapest) {
                cheapest = wires;
                source_side = reachable;
            }
        }
        (source_side, self.components() - source_side)
    }

    fn component_id<'a>(
        &mut self,
        ids: &mut HashMap<&'a str, usize>,
        name: &'a str,
    ) -> Result<usize, AocError> {
        if name.is_empty() {
            return Err(AocError::InvalidInput("component without a name"));
        }
        Ok(*ids.entry(name).or_insert_with(|| {
            self.outgoing.push(Vec::new());
            self.outgoing.len() - 1
        }))
    }

    fn connect(&mut self, left: usize, right: usize) {
        self.outgoing[left].push(self.heads.len());
        self.heads.push(right);
        self.outgoing[right].push(self.heads.len());
        self.heads.push(left);
    }

    /// Pushes unit flow from component 0 to `sink` along augmenting paths and
    /// returns the flow plus the number of components on the source side of the
    /// resulting cut, or `None` once the flow can no longer stay below `limit`.
    fn flow_under(&self, sink: usize, limit: usize) -> Option<(usize, usize)> {
        let mut capacity = vec![1u8; self.heads.len()];
        let mut flow = 0;
        loop {
            let mut incoming = vec![usize::MAX; self.components()];
            let mut reached = vec![false; self.components()];
            let mut queue = vec![0];
            reached[0] = true;
            let mut next = 0;
            while let Some(&component) = queue.get(next) {
                next += 1;
                for &arc in &self.outgoing[component] {
                    let head = self.heads[arc];
                    if capacity[arc] == 0 || reached[head] {
                        continue;
                    }
                    reached[head] = true;
                    incoming[head] = arc;
                    queue.push(head);
                }
            }
            if !reached[sink] {
                return Some((flow, queue.len()));
            }
            if flow + 1 >= limit {
                return None;
            }
            let mut component = sink;
            while component != 0 {
                let arc = incoming[component];
                capacity[arc] -= 1;
                capacity[arc ^ 1] += 1;
                component = self.heads[arc ^ 1];
            }
            flow += 1;
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn example_splits_into_nine_and_six_components() {
        let diagram = Diagram::parse(include_str!("../../examples/day25.txt")).unwrap();
        assert_eq!(diagram.components(), 15);
        let (left, right) = diagram.minimum_cut();
        assert_eq!((left.max(right), left.min(right)), (9, 6));
    }

    #[test]
    fn wires_are_undirected() {
        let forwards = Diagram::parse("a: b c\nb: c").unwrap();
        let backwards = Diagram::parse("c: a b\nb: a").unwrap();
        assert_eq!(forwards.minimum_cut(), (1, 2));
        assert_eq!(backwards.minimum_cut(), (1, 2));
    }
}
