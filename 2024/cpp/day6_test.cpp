#include "check.hpp"
#include "day6.hpp"

int main() {
    const std::string example =
        "....#.....\n"
        ".........#\n"
        "..........\n"
        "..#.......\n"
        ".......#..\n"
        "..........\n"
        ".#..^.....\n"
        "........#.\n"
        "#.........\n"
        "......#...\n";
    const auto answers = day6::solve(example);
    check::equal(answers.part1, 41);
    check::equal(answers.part2, 6);

    const auto straight = day6::solve("...\n.^.\n...\n");
    check::equal(straight.part1, 2);
    check::equal(straight.part2, 0);

    check::throws([] { day6::solve(""); });
    check::throws([] { day6::solve("..\n...\n"); });
    check::throws([] { day6::solve("..\n.?\n"); });
    check::throws([] { day6::solve("...\n...\n"); });
    check::throws([] { day6::solve("^.\n.^\n"); });

    return check::exit_code();
}
