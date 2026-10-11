use crate::day20::custom_error::AocError;
use std::collections::{HashMap, VecDeque};

#[derive(Debug, Clone)]
enum Kind {
    Broadcaster,
    FlipFlop(bool),
    Conjunction(HashMap<String, bool>),
}

#[derive(Debug, Clone)]
struct ModuleDef {
    kind: Kind,
    destinations: Vec<String>,
}

/// A graph of pulse-propagation modules, parsed from puzzle input.
#[derive(Debug, Clone)]
pub struct Network {
    modules: HashMap<String, ModuleDef>,
}

/// One delivered pulse: `source` sent a high (`true`) or low pulse to `target`.
pub type Pulse = (String, String, bool);

impl Network {
    pub fn parse(input: &str) -> Result<Self, AocError> {
        let mut modules = HashMap::new();
        for line in input.lines().map(str::trim).filter(|line| !line.is_empty()) {
            let (name_part, dest_part) = line
                .split_once("->")
                .ok_or(AocError::InvalidInput("module line is missing '->'"))?;
            let name_part = name_part.trim();
            let (name, kind) = if let Some(rest) = name_part.strip_prefix('%') {
                (rest, Kind::FlipFlop(false))
            } else if let Some(rest) = name_part.strip_prefix('&') {
                (rest, Kind::Conjunction(HashMap::new()))
            } else {
                (name_part, Kind::Broadcaster)
            };
            if name.is_empty() {
                return Err(AocError::InvalidInput("module is missing a name"));
            }
            let destinations: Vec<String> = dest_part
                .split(',')
                .map(|d| d.trim().to_string())
                .filter(|d| !d.is_empty())
                .collect();
            if destinations.is_empty() {
                return Err(AocError::InvalidInput("module has no destinations"));
            }
            modules.insert(name.to_string(), ModuleDef { kind, destinations });
        }
        if modules.is_empty() {
            return Err(AocError::InvalidInput("module configuration is empty"));
        }

        let edges: Vec<(String, String)> = modules
            .iter()
            .flat_map(|(name, def)| {
                def.destinations
                    .iter()
                    .map(move |dest| (name.clone(), dest.clone()))
            })
            .collect();
        for (source, target) in edges {
            if let Some(def) = modules.get_mut(&target) {
                if let Kind::Conjunction(memory) = &mut def.kind {
                    memory.insert(source, false);
                }
            }
        }

        Ok(Network { modules })
    }

    /// Pushes the button once and returns every pulse sent, in delivery order,
    /// including the initial low pulse the button sends to `broadcaster`.
    pub fn press_button(&mut self) -> Vec<Pulse> {
        let mut log = Vec::new();
        let mut queue = VecDeque::new();
        queue.push_back(("button".to_string(), "broadcaster".to_string(), false));

        while let Some((source, target, high)) = queue.pop_front() {
            log.push((source.clone(), target.clone(), high));
            let Some(def) = self.modules.get_mut(&target) else {
                continue;
            };
            let send = match &mut def.kind {
                Kind::Broadcaster => Some(high),
                Kind::FlipFlop(on) => {
                    if high {
                        None
                    } else {
                        *on = !*on;
                        Some(*on)
                    }
                }
                Kind::Conjunction(memory) => {
                    memory.insert(source, high);
                    Some(!memory.values().all(|&remembered| remembered))
                }
            };
            if let Some(pulse) = send {
                for destination in &def.destinations {
                    queue.push_back((target.clone(), destination.clone(), pulse));
                }
            }
        }

        log
    }

    /// The single module name that is referenced as a destination but never
    /// defined with its own line (e.g. `rx` in the real puzzle input, or
    /// `output` in the official examples). Errors if there is none or more
    /// than one such name.
    pub fn implicit_sink(&self) -> Result<String, AocError> {
        let mut undefined: Vec<&String> = self
            .modules
            .values()
            .flat_map(|def| &def.destinations)
            .filter(|dest| !self.modules.contains_key(dest.as_str()))
            .collect();
        undefined.sort_unstable();
        undefined.dedup();
        match undefined.as_slice() {
            [single] => Ok(single.to_string()),
            [] => Err(AocError::InvalidInput("no undefined sink module found")),
            _ => Err(AocError::InvalidInput(
                "more than one undefined sink module",
            )),
        }
    }

    /// Names of every module that lists `target` as a destination.
    pub fn inputs_of(&self, target: &str) -> Vec<String> {
        self.modules
            .iter()
            .filter(|(_, def)| def.destinations.iter().any(|dest| dest == target))
            .map(|(name, _)| name.clone())
            .collect()
    }
}
