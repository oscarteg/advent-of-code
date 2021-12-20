pub fn sonar_sweep(data: Vec<u64>) -> u64 {
    data.windows(2).filter(|x| x[0] < x[1]).count() as u64
}

pub fn sonar_sweep_window(data: Vec<u64>) -> u64 {
    data.windows(4).filter(|x| x[3] > x[0]).count() as u64
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::utils::{clean_input, read_file};

    const INPUT: &'static str = "
199
200
208
210
200
207
240
269
260
263";

    #[test]
    pub fn test_sonar_sweep_example() {
        let input = clean_input(INPUT);

        let vec = input.map(|l| l.parse::<u64>().unwrap());
        assert_eq!(sonar_sweep(vec.collect()), 7);
    }

    #[test]
    fn test_sonar_sweep() {
        let file = read_file("input/day1_b.txt");
        let input = clean_input(file.as_str());

        let vec = input.map(|l| l.parse::<u64>().unwrap());
        assert_eq!(sonar_sweep(vec.collect()), 1692);
    }

    #[test]
    fn test_sonar_sweep_window_example() {
        let input = clean_input(INPUT);

        let vec = input.map(|l| l.parse::<u64>().unwrap());
        assert_eq!(sonar_sweep_window(vec.collect::<Vec<u64>>()), 5);
    }

    #[test]
    fn test_sonar_sweep_window() {
        let file = read_file("input/day1_b.txt");
        let input = clean_input(file.as_str());
        let vec = input
            .map(|l| l.parse::<u64>().unwrap())
            .collect::<Vec<u64>>();
        assert_eq!(sonar_sweep_window(vec), 1724);
    }
}
