use anyhow::{ensure, Context, Result};

// Day 25: Combo Breaker
//
// The handshake uses repeated modular exponentiation (transform) with a fixed
// subject number (7) and modulus (20201227). Each device's public key reveals
// how many times it applied the transform (its "loop size"), found here by
// brute-force multiplication since the modulus is small. Once one loop size is
// known, applying it to the other device's public key yields the shared
// encryption key.

const SUBJECT_NUMBER: u64 = 7;
const MODULUS: u64 = 20_201_227;

struct PublicKeys {
    card: u64,
    door: u64,
}

fn parse(input: &str) -> Result<PublicKeys> {
    let mut keys = input.lines().map(str::trim).filter(|line| !line.is_empty());

    let card = keys
        .next()
        .context("missing card public key")?
        .parse()
        .context("card public key is not a number")?;
    let door = keys
        .next()
        .context("missing door public key")?
        .parse()
        .context("door public key is not a number")?;
    ensure!(keys.next().is_none(), "expected exactly two public keys");

    Ok(PublicKeys { card, door })
}

fn find_loop_size(public_key: u64) -> Result<u64> {
    let mut value = 1;
    for loop_size in 1.. {
        value = (value * SUBJECT_NUMBER) % MODULUS;
        if value == public_key {
            return Ok(loop_size);
        }
        if loop_size > MODULUS {
            break;
        }
    }
    anyhow::bail!("could not find a loop size for public key {public_key}")
}

fn transform(subject: u64, loop_size: u64) -> u64 {
    let mut value = 1;
    for _ in 0..loop_size {
        value = (value * subject) % MODULUS;
    }
    value
}

pub fn solve(input: &str) -> Result<u64> {
    let keys = parse(input)?;
    let loop_size = find_loop_size(keys.card)?;
    Ok(transform(keys.door, loop_size))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn sample() {
        let data = include_str!("../../input/2020/day25-example.txt");
        assert_eq!(solve(data).unwrap(), 14897079);
    }

    #[test]
    fn find_loop_size_matches_known_values() {
        assert_eq!(find_loop_size(5764801).unwrap(), 8);
        assert_eq!(find_loop_size(17807724).unwrap(), 11);
    }

    #[test]
    fn door_loop_size_gives_same_key() {
        let keys = parse("5764801\n17807724\n").unwrap();
        let loop_size = find_loop_size(keys.door).unwrap();
        assert_eq!(transform(keys.card, loop_size), 14897079);
    }

    #[test]
    fn missing_door_key_is_an_error() {
        assert!(parse("5764801").is_err());
    }

    #[test]
    fn non_numeric_key_is_an_error() {
        assert!(solve("not-a-number\n17807724\n").is_err());
    }

    #[test]
    fn extra_line_is_an_error() {
        assert!(parse("5764801\n17807724\n123\n").is_err());
    }
}
