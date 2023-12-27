mod days;

fn main() {
    println!("Hello, world!");
}

pub mod utils {
    use anyhow::Result;
    use std::fmt::Debug;
    use std::fs::{self, read_to_string};
    use std::str::FromStr;

    /// Reads an file as a string
    pub fn read_file(filename: &str) -> String {
        fs::read_to_string(filename).expect("File does not exist")
    }

    pub fn read_one_per_line<T>(path: &str) -> Result<Vec<T>>
    where
        T: FromStr,
    {
        Ok(read_to_string(path)?
            .split("\n")
            .filter_map(|line| line.parse::<T>().ok())
            .collect())
    }

    pub fn clean_input(input: &str) -> impl Iterator<Item = &str> {
        input.lines().map(|l| l.trim()).filter(|l| !l.is_empty())
    }

    /// Cleans the input and returns a Vec of the lines
    pub fn input_vec<T>(input: &str) -> Vec<T>
    where
        T: FromStr,
        <T as FromStr>::Err: Debug,
    {
        clean_input(input).map(|l| l.parse().unwrap()).collect()
    }
}
