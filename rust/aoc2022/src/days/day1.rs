fn max_calorie_elf(input: &str) -> String {
    let result = input
        .split("\n\n")
        .map(|l| {
            l.lines()
                .map(|item| item.parse::<u32>().unwrap())
                .sum::<u32>()
        })
        .max()
        .unwrap();

    result.to_string()
}

#[cfg(test)]
mod tests {
    use crate::utils::read_file;

    use super::*;

    const INPUT: &str = "1000
2000
3000

4000

5000
6000

7000
8000
9000

10000";

    #[test]
    pub fn part_1() {
        let part_1 = max_calorie_elf(INPUT);
        assert_eq!(part_1, "24000");

        let file = read_file("input/day1.txt");
        let part_2 = max_calorie_elf(file.as_str());
        assert_eq!(part_2, "24000");
    }

    #[test]
    fn part_2() {}
}
