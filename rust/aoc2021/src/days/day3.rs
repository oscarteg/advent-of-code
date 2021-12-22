#![allow(dead_code)]
#![allow(unused_variables)]

pub fn part_a<'a>(lines: Vec<&str>) -> (u32, u32, u32) {
    let mut bits = Vec::<(i32, i32)>::new();

    for line in lines {
        line.chars().enumerate().for_each(|(i, b)| {
            if i >= bits.len() {
                bits.push((0, 0));
            }

            if b == '0' {
                bits[i].0 += 1;
            } else {
                bits[i].1 += 1;
            }
        });
    }

    let mut gamma = 0;
    let mut epsilon = 0;
    bits.iter().enumerate().for_each(|(i, (a, b))| {
        if a < b {
            gamma |= 1 << (bits.len() - i - 1)
        }
        if a > b {
            epsilon |= 1 << (bits.len() - i - 1)
        }
    });

    (gamma, epsilon, gamma * epsilon)
}

pub fn part_b(lines: Vec<&str>) -> (u32, u32, u32) {
    let mut oxygen: u32 = 0;
    let mut scrubber: u32 = 0;

    let mut curr_oxygen = String::new();
    let mut curr_scrubber = String::new();
    let l = lines.first().unwrap_or(&"").len();
    for i in 0..l {
        let mut count = (0, 0);
        for j in 0..lines.len() {
            if !lines[j].starts_with(&curr_oxygen) {
                continue;
            }

            if lines[j].chars().nth(i).unwrap() == '0' {
                count.0 += 1;
            } else {
                count.1 += 1;
            }
        }

        if count.0 <= count.1 {
            oxygen |= 1 << (l - i - 1);
            curr_oxygen += "1";
        } else {
            curr_oxygen += "0";
        }

        let mut count = (0, 0);
        for j in 0..lines.len() {
            if !lines[j].starts_with(&curr_scrubber) {
                continue;
            }

            if lines[j].chars().nth(i).unwrap() == '0' {
                count.0 += 1;
            } else {
                count.1 += 1;
            }
        }

        if count.1 == 0 || (count.0 > 0 && count.0 <= count.1) {
            curr_scrubber += "0";
        } else {
            scrubber |= 1 << (l - i - 1);
            curr_scrubber += "1";
        }
    }

    (oxygen, scrubber, oxygen * scrubber)
}

#[cfg(test)]
mod tests {
    #![allow(dead_code)]
    #![allow(unused_variables)]

    use crate::utils::{clean_input, read_file};

    use super::*;

    const INPUT: &'static str = "
00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010";

    #[test]
    pub fn test_part_a_example() {
        let input = clean_input(INPUT).collect();

        let (gamma, epsilon, gamma_epsilon) = part_a(input);

        assert_eq!(gamma, 22);
        assert_eq!(epsilon, 9);
        assert_eq!(gamma_epsilon, 198);
    }

    #[test]
    fn test_part_a() {
        let file = read_file("input/day3.txt");
        let input = clean_input(file.as_str()).collect();

        let (gamma, epsilon, gamma_epsilon) = part_a(input);

        assert_eq!(gamma, 3069);
        assert_eq!(epsilon, 1026);
        assert_eq!(gamma_epsilon, 3148794);
    }

    #[test]
    fn test_part_b_example() {
        let input: Vec<&str> = clean_input(INPUT).collect();

        let (gamma, epsilon, gamma_epsilon) = part_b(input);

        assert_eq!(gamma, 23);
        assert_eq!(epsilon, 10);
        assert_eq!(gamma_epsilon, 230);
    }

    #[test]
    fn test_part_b() {
        let file = read_file("input/day3.txt");
        let input = clean_input(file.as_str()).collect();

        let (gamma, epsilon, gamma_epsilon) = part_b(input);

        assert_eq!(gamma, 2349);
        assert_eq!(epsilon, 1190);
        assert_eq!(gamma_epsilon, 2795310);
    }
}
