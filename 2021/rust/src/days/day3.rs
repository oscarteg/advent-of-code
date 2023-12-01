#![allow(dead_code)]
#![allow(unused_variables)]

pub fn part_a<'a>(lines: Vec<&str>) -> (u32, u32, u32) {
    let max = (lines.len() + 1) / 2;

    let bits = lines[0].len();

    // To make sure we only perform operations on the bits we want.
    let bitmask = (1 << bits) - 1;

    // Bit shift gamma as required
    let mut gamma = 0;

    lines
        .into_iter()
        .fold(vec![0; bits], |count, reading| {
            // Count all occurrences of 1 at each bit position within a reading
            // It returns 1 array with an incremented value when 1 the bit position has value 1
            count
                .into_iter()
                .enumerate()
                .map(|(index, number)| {
                    // Compare
                    return if reading.as_bytes()[index] == b'1' {
                        number + 1
                    } else {
                        number
                    };
                })
                .collect()
        })
        .into_iter()
        .enumerate()
        .for_each(|(i, bit)| {
            // If the number of bits is more than half of all lines than 1 occures more than 0
            if bit >= max {
                gamma += 1 << (bits - (i + 1))
            }
        });

    // let epsilon = !gamma & bitmask;
    let epsilon = !gamma & bitmask;

    (gamma, epsilon, gamma * epsilon)
}

pub fn part_b(readings: Vec<&str>) -> (u32, u32, u32) {

    let bits = readings[0].len();
    let max = (lines.len() + 1) / 2;

    let bits = lines[0].len();

    // To make sure we only perform operations on the bits we want.
    let bitmask = (1 << bits) - 1;

    // Bit shift gamma as required
    let mut gamma = 0;

    let oxygen_rating = (0..bits)
        .scan(readings.clone(), |oxygen, i| {
            let max = (oxygen.len() as f32 / 2.0).ceil() as usize;

            let mut sig = b'0';
            let count = oxygen.iter()
                .filter(|o| (**o).as_bytes()[i] == b'1')
                .count();
            if count >= max {
                sig = b'1';
            }

            oxygen.drain_filter(|o| (*o).as_bytes()[i] != sig);
            oxygen.first().copied()
        })
        .last()
        .unwrap();

    let co2_rating = (0..bits)
        .scan(readings, |co2, i| {
            let max = (co2.len() as f32 / 2.0).ceil() as usize;

            let mut sig = b'1';
            let count = co2.iter()
                .filter(|c| (**c).as_bytes()[i] == b'1')
                .count();
            if count >= max {
                sig = b'0';
            }

            co2.drain_filter(|c| (*c).as_bytes()[i] != sig);
            co2.first().copied()
        })
        .last()
        .unwrap();

    let oxygen = usize::from_str_radix(oxygen_rating, 2).unwrap();
    let co2 = usize::from_str_radix(co2_rating, 2).unwrap();

    (oxygen as u32, co2 as u32, (oxygen * co2) as u32)
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

        let (oxygen, co2, oxygen_co2) = part_b(input);

        assert_eq!(oxygen, 23);
        assert_eq!(co2, 10);
        assert_eq!(oxygen_co2, 230);
    }

    #[test]
    fn test_part_b() {
        let file = read_file("input/day3.txt");
        let input = clean_input(file.as_str()).collect();

        let (oxygen, co2, oxygen_co2) = part_b(input);

        assert_eq!(oxygen, 2349);
        assert_eq!(co2, 1190);
        assert_eq!(oxygen_co2, 2795310);
    }
}
