#include "check.hpp"
#include "day1.hpp"

int main() {
    const std::string example = R"(
3   4
4   3
2   5
1   3
3   9
3   3
)";

    const auto answers = day1::solve(example);
    check::equal(answers.part1, 11);
    check::equal(answers.part2, 31);

    check::throws([] { day1::solve("3\n"); });
    check::throws([] { day1::solve("abc 123\n"); });
    check::throws([] { day1::solve(""); });

    return check::exit_code();
}
