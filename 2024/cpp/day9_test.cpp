#include "check.hpp"
#include "day9.hpp"

int main() {
    constexpr std::string_view example = "2333133121414131402";

    const auto answers = day9::solve(example);
    check::equal(answers.part1, 1928);
    check::equal(answers.part2, 2858);

    // Edge case: a single file with no free space at all.
    const auto single = day9::solve("5");
    check::equal(single.part1, 0);
    check::equal(single.part2, 0);

    check::throws([] { day9::solve(""); });
    check::throws([] { day9::solve("12a45"); });

    return check::exit_code();
}
