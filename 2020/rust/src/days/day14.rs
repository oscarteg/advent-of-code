use regex::Regex;
use std::collections::HashMap;

const MAX_FLOATING_BITS: usize = 20;

pub fn solve(input: &str) -> anyhow::Result<(u64, u64)> {
    let part1 = part1(input)?;
    let part2 = part2(input)?;
    Ok((part1, part2))
}

fn part1(input: &str) -> anyhow::Result<u64> {
    let mut memory: HashMap<u64, u64> = HashMap::new();
    let mut mask = String::new();

    let mask_regex = Regex::new(r"^mask = ([01X]+)$")?;
    let mem_regex = Regex::new(r"^mem\[(\d+)\] = (\d+)$")?;

    for line in input.lines() {
        let line = line.trim();
        if let Some(caps) = mask_regex.captures(line) {
            mask = caps[1].to_string();
        } else if let Some(caps) = mem_regex.captures(line) {
            let addr: u64 = caps[1].parse()?;
            let value: u64 = caps[2].parse()?;
            let masked_value = apply_mask_part1(value, &mask);
            memory.insert(addr, masked_value);
        }
    }

    Ok(memory.values().sum())
}

fn apply_mask_part1(value: u64, mask: &str) -> u64 {
    let mut result = value;
    for (i, bit) in mask.chars().rev().enumerate() {
        match bit {
            '0' => result &= !(1u64 << i),
            '1' => result |= 1u64 << i,
            'X' => {}
            _ => {}
        }
    }
    result
}

fn part2(input: &str) -> anyhow::Result<u64> {
    let mut memory: HashMap<u64, u64> = HashMap::new();
    let mut mask = String::new();

    let mask_regex = Regex::new(r"^mask = ([01X]+)$")?;
    let mem_regex = Regex::new(r"^mem\[(\d+)\] = (\d+)$")?;

    for line in input.lines() {
        let line = line.trim();
        if let Some(caps) = mask_regex.captures(line) {
            mask = caps[1].to_string();
        } else if let Some(caps) = mem_regex.captures(line) {
            let addr: u64 = caps[1].parse()?;
            let value: u64 = caps[2].parse()?;
            let addresses = decode_addresses_part2(addr, &mask)?;
            for decoded_addr in addresses {
                memory.insert(decoded_addr, value);
            }
        }
    }

    Ok(memory.values().sum())
}

fn decode_addresses_part2(address: u64, mask: &str) -> anyhow::Result<Vec<u64>> {
    let mut floating_positions = Vec::new();
    let mut base_addr = address;

    for (i, bit) in mask.chars().rev().enumerate() {
        match bit {
            '0' => {}
            '1' => base_addr |= 1u64 << i,
            'X' => {
                floating_positions.push(i);
                base_addr &= !(1u64 << i);
            }
            _ => {}
        }
    }

    let num_floating = floating_positions.len();
    anyhow::ensure!(
        num_floating <= MAX_FLOATING_BITS,
        "Mask has {} floating bits, which exceeds maximum of {}",
        num_floating,
        MAX_FLOATING_BITS
    );

    let mut addresses = Vec::new();

    for variation in 0..(1u64 << num_floating) {
        let mut addr = base_addr;
        for (idx, &pos) in floating_positions.iter().enumerate() {
            if (variation >> idx) & 1 == 1 {
                addr |= 1u64 << pos;
            }
        }
        addresses.push(addr);
    }

    Ok(addresses)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_part1_example() {
        let input = "mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X
mem[8] = 11
mem[7] = 101
mem[8] = 0";
        let result = part1(input).unwrap();
        assert_eq!(result, 165);
    }

    #[test]
    fn test_part1_mask_application() {
        let masked = apply_mask_part1(11, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X");
        assert_eq!(masked, 73);
    }

    #[test]
    fn test_part2_example() {
        let input = include_str!("../../input/2020/day14-example.txt");
        let result = part2(input).unwrap();
        assert_eq!(result, 208);
    }

    #[test]
    fn test_part2_address_decoding() {
        let addresses = decode_addresses_part2(27, "000000000000000000000000000000X1101X").unwrap();
        assert_eq!(addresses.len(), 4);
        assert!(addresses.contains(&26));
        assert!(addresses.contains(&27));
        assert!(addresses.contains(&58));
        assert!(addresses.contains(&59));
    }

    #[test]
    fn test_part2_too_many_floating_bits() {
        let mask = "X".repeat(25);
        let result = decode_addresses_part2(0, &mask);
        assert!(result.is_err());
        assert!(result.unwrap_err().to_string().contains("exceeds maximum"));
    }

    #[test]
    fn test_malformed_mask() {
        let input = "invalid mask line";
        let result = part1(input);
        assert_eq!(result.unwrap(), 0);
    }
}
