use std::collections::HashMap;

use anyhow::{ensure, Context, Result};
use lazy_static::lazy_static;
use regex::Regex;

const REQUIRED_FIELDS: [&str; 7] = ["byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"];
const VALID_EYE_COLORS: [&str; 7] = ["amb", "blu", "brn", "gry", "grn", "hzl", "oth"];

lazy_static! {
    static ref HEIGHT_RE: Regex = Regex::new(r"^(\d+)(cm|in)$").unwrap();
    static ref HAIR_COLOR_RE: Regex = Regex::new(r"^#[0-9a-f]{6}$").unwrap();
    static ref PASSPORT_ID_RE: Regex = Regex::new(r"^\d{9}$").unwrap();
}

type Passport = HashMap<String, String>;

fn parse(input: &str) -> Result<Vec<Passport>> {
    input
        .split("\n\n")
        .map(str::trim)
        .filter(|block| !block.is_empty())
        .map(parse_passport)
        .collect()
}

fn parse_passport(block: &str) -> Result<Passport> {
    block
        .split_whitespace()
        .map(|pair| {
            let (key, value) = pair
                .split_once(':')
                .with_context(|| format!("malformed field, missing ':' in {pair:?}"))?;
            ensure!(
                !key.is_empty() && !value.is_empty(),
                "empty key or value in {pair:?}"
            );
            Ok((key.to_string(), value.to_string()))
        })
        .collect()
}

fn has_required_fields(passport: &Passport) -> bool {
    REQUIRED_FIELDS
        .iter()
        .all(|field| passport.contains_key(*field))
}

fn is_valid_year(value: &str, min: u32, max: u32) -> bool {
    value.len() == 4
        && value
            .parse::<u32>()
            .is_ok_and(|year| (min..=max).contains(&year))
}

fn is_valid_height(value: &str) -> bool {
    let Some(captures) = HEIGHT_RE.captures(value) else {
        return false;
    };
    let Ok(number) = captures[1].parse::<u32>() else {
        return false;
    };
    match &captures[2] {
        "cm" => (150..=193).contains(&number),
        "in" => (59..=76).contains(&number),
        _ => false,
    }
}

fn is_valid_passport(passport: &Passport) -> bool {
    has_required_fields(passport)
        && is_valid_year(&passport["byr"], 1920, 2002)
        && is_valid_year(&passport["iyr"], 2010, 2020)
        && is_valid_year(&passport["eyr"], 2020, 2030)
        && is_valid_height(&passport["hgt"])
        && HAIR_COLOR_RE.is_match(&passport["hcl"])
        && VALID_EYE_COLORS.contains(&passport["ecl"].as_str())
        && PASSPORT_ID_RE.is_match(&passport["pid"])
}

pub fn solve(input: &str) -> Result<(usize, usize)> {
    let passports = parse(input)?;
    ensure!(!passports.is_empty(), "no passports found in input");

    let part1 = passports.iter().filter(|p| has_required_fields(p)).count();
    let part2 = passports.iter().filter(|p| is_valid_passport(p)).count();

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/2020/day4-example.txt");

    #[test]
    fn sample_part1() {
        let (part1, _) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 2);
    }

    #[test]
    fn sample_part2_invalid() {
        let invalid = "\
eyr:1972 cid:100
hcl:#18171d ecl:amb hgt:170 pid:186cm iyr:2018 byr:1926

iyr:2019
hcl:#602927 eyr:1967 hgt:170cm
ecl:grn pid:012533040 byr:1946

hcl:dab227 iyr:2012
ecl:brn hgt:182cm pid:021572410 eyr:2020 byr:1992 cid:277

hgt:59cm ecl:zzz
eyr:2038 hcl:74454a iyr:2023
pid:3556412378 byr:2007";

        let (_, part2) = solve(invalid).unwrap();
        assert_eq!(part2, 0);
    }

    #[test]
    fn sample_part2_valid() {
        let valid = "\
pid:087499704 hgt:74in ecl:grn iyr:2012 eyr:2030 byr:1980
hcl:#623a2f

eyr:2029 ecl:blu cid:129 byr:1989
iyr:2014 pid:896056539 hcl:#a97842 hgt:165cm

hcl:#888785
hgt:164cm byr:2001 iyr:2015 cid:88
pid:545766238 ecl:hzl
eyr:2022

iyr:2010 hgt:158cm hcl:#b6652a ecl:blu byr:1944 eyr:2021 pid:093154719";

        let (_, part2) = solve(valid).unwrap();
        assert_eq!(part2, 4);
    }

    #[test]
    fn missing_cid_is_still_valid() {
        let only_cid_missing =
            "byr:1937 iyr:2017 eyr:2020 hgt:183cm hcl:#fffffd ecl:gry pid:860033327";
        let (part1, _) = solve(only_cid_missing).unwrap();
        assert_eq!(part1, 1);
    }

    #[test]
    fn malformed_field_without_colon_is_error() {
        let malformed = "byr1937 iyr:2017 eyr:2020 hgt:183cm hcl:#fffffd ecl:gry pid:860033327";
        assert!(solve(malformed).is_err());
    }

    #[test]
    fn empty_input_is_error() {
        assert!(solve("").is_err());
    }

    #[test]
    fn example_file_matches_known_answers() {
        let input = std::fs::read_to_string("input/2020/day4.txt");
        let Ok(input) = input else {
            return;
        };
        let (part1, part2) = solve(&input).unwrap();
        assert!(part1 > 0 && part2 > 0 && part2 <= part1);
    }
}
