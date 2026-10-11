#include "check.hpp"
#include "day5.hpp"

int main() {
    {
        const std::string example = R"(47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47)";

        auto answers = day5::solve(example);
        check::equal(answers.part1, 143LL);
        check::equal(answers.part2, 123LL);
    }

    {
        const std::string single = "1|2\n\n1\n";
        auto answers = day5::solve(single);
        check::equal(answers.part1, 1LL);
    }

    {
        const std::string pair = "1|2\n\n1,2\n";
        auto answers = day5::solve(pair);
        check::equal(answers.part1, 2LL);
    }

    check::throws([] { day5::solve("1 2\n\n1,2\n"); });

    check::throws([] { day5::solve("1|abc\n\n1,2\n"); });

    return check::exit_code();
}
