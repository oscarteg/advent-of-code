#include "check.hpp"
#include "day2.hpp"

int main() {
    const std::string example =
        "7 6 4 2 1\n"
        "1 2 7 8 9\n"
        "9 7 6 2 1\n"
        "1 3 2 4 5\n"
        "8 6 4 4 1\n"
        "1 3 6 7 9\n";

    const auto answers = day2::solve(example);
    check::equal(answers.part1, 2);
    check::equal(answers.part2, 4);

    // Edge case: a single-level report is trivially safe for both parts.
    const auto single = day2::solve("5\n");
    check::equal(single.part1, 1);
    check::equal(single.part2, 1);

    check::throws([] { day2::solve("1 two 3\n"); });
    check::throws([] { day2::solve("   \n"); });

    return check::exit_code();
}
