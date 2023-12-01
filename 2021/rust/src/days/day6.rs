use std::collections::VecDeque;

pub fn get_counts(data: Vec<usize>, days: u32) -> u64 {
    let mut counts = VecDeque::from(vec![0; 9]);

    data.iter().for_each(|i| counts[*i] += 1);

    (0..days).for_each(|_| {
        let new_babies = counts.pop_front().unwrap();
        counts[6] += new_babies;
        counts.push_back(new_babies);
    });

    counts.iter().sum()
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::utils::{clean_input, read_file};

    #[test]
    pub fn test_sonar_sweep_example() {}
}
