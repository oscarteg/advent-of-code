#include "check.hpp"
#include "day7.hpp"

int main() {
    constexpr std::string_view example =
        "190: 10 19\n"
        "3267: 81 40 27\n"
        "83: 17 5\n"
        "156: 15 6\n"
        "7290: 6 8 6 15\n"
        "161011: 16 10 13\n"
        "192: 17 8 14\n"
        "21037: 9 7 18 13\n"
        "292: 11 6 16 20\n";

    const auto answers = day7::solve(example);
    check::equal(answers.part1, 3749);
    check::equal(answers.part2, 11387);

    // Edge case: single operand, equation is trivially solvable.
    const auto single = day7::solve("5: 5\n");
    check::equal(single.part1, 5);
    check::equal(single.part2, 5);

    // Edge case: no combination of operators reaches the target.
    const auto unsolvable = day7::solve("10: 2 2\n");
    check::equal(unsolvable.part1, 0);
    check::equal(unsolvable.part2, 0);

    check::throws([] { day7::solve(""); });
    check::throws([] { day7::solve("190 10 19\n"); });
    check::throws([] { day7::solve("190:\n"); });
    check::throws([] { day7::solve("abc: 1 2\n"); });

    return check::exit_code();
}
