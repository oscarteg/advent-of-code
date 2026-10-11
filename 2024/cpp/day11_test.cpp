#include "check.hpp"
#include "day11.hpp"

int main() {
    constexpr std::string_view example = "125 17\n";

    const auto afterSix = day11::solveWith(example, 6, 6);
    check::equal(afterSix.part1, 22);

    const auto official = day11::solveWith(example, 25, 75);
    check::equal(official.part1, 55312);
    check::equal(official.part2, 65601038650482LL);

    const auto answers = day11::solve(example);
    check::equal(answers.part1, 55312);
    check::equal(answers.part2, 65601038650482LL);

    // Edge case: a single zero stone becomes 1 after one blink.
    const auto zero = day11::solveWith("0\n", 1, 1);
    check::equal(zero.part1, 1);

    check::throws([] { day11::solve(""); });
    check::throws([] { day11::solve("abc\n"); });
    check::throws([] { day11::solve("-5\n"); });

    return check::exit_code();
}
