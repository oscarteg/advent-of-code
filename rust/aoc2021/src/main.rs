mod days;

fn main() {
    println!("Hello, world!");
}

#[cfg(test)]
pub mod utils {
    use std::fmt::Debug;
    use std::fs;
    use std::str::FromStr;

    /// Reads an file as a string
    pub fn read_file(filename: &str) -> String {
        fs::read_to_string(filename).expect("File does not exist")
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
