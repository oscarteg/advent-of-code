pub fn solve_part_01(input: &[String]) -> u32 {
    println!("{:?}", input);
    return 2224913600;
}

pub fn solve_part_02(input: &[String]) -> u32 {
    return 2224913600;
}

#[cfg(test)]
mod tests {
    use crate::utils::{clean_input, input_vec};

    use super::*;

    /// Test example data on part 1
    #[test]
    fn sample_01() {
        let data = "
ecl:gry pid:860033327 eyr:2020 hcl:#fffffd
byr:1937 iyr:2017 cid:147 hgt:183cm

iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884
hcl:#cfa07d byr:1929

hcl:#ae17e1 iyr:2013
eyr:2024
ecl:brn pid:760753108 byr:1931
hgt:179cm

hcl:#cfa07d eyr:2025 pid:166559648
iyr:2011 ecl:brn hgt:59in";

        assert_eq!(solve_part_01(&input_vec(data)), 2)
    }

    /// Test example data on part 2
    #[test]
    fn sample_02() {
        let data = "";

        assert_eq!(solve_part_02(&input_vec(data)), 336)
    }
}
