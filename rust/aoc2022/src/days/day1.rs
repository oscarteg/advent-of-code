fn part_one(input: &str) -> Option<i32> {
    let elfs = input.split("\n\n");

    let mut max_calories = Some(0);
    for elf in elfs {
        let calories = elf
            .split('\n')
            .map(|food| food.parse::<i32>().unwrap_or(0))
            .reduce(|a, b| a + b);
        if calories.is_some() && calories > max_calories {
            max_calories = calories;
        }
    }
    max_calories
}

fn part_two(input: &str) -> Option<i32> {
    let elfs = input.split("\n\n");

    let mut calories_arr: Vec<i32> = vec![];
    for elf in elfs {
        let calories = elf
            .split('\n')
            .map(|food| food.parse::<i32>().unwrap_or(0))
            .reduce(|a, b| a + b);
        if calories.is_some() {
            calories_arr.push(calories.unwrap_or(0))
        }
    }
    calories_arr.sort_by(|a, b| b.cmp(a));
    let result: i32 = calories_arr[0..3].iter().sum();
    Some(result)
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
    pub fn test_part_1() {
        let file: String = read_file("input/day1.txt");
        assert_eq!(part_one(INPUT), Some(24000));
        assert_eq!(part_two(file.as_str()), Some(203905));
    }

    #[test]
    fn test_part_2() {
        let file: String = read_file("input/day1.txt");
        assert_eq!(part_two(INPUT), Some(45000));
        assert_eq!(part_two(file.as_str()), Some(203905));
    }
}
