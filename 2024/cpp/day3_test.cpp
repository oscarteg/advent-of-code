#include "check.hpp"
#include "day3.hpp"

int main() {
    const auto part1_example =
        day3::solve("xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))");
    check::equal(part1_example.part1, 161LL);

    const auto part2_example =
        day3::solve("xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))");
    check::equal(part2_example.part1, 161LL);
    check::equal(part2_example.part2, 48LL);

    const auto across_lines = day3::solve("don't()mul(2,3)\nmul(4,5)do()mul(1,1)");
    check::equal(across_lines.part1, 27LL);
    check::equal(across_lines.part2, 1LL);

    const auto malformed = day3::solve("mul(1234,5) mul(2, 3) mul(4*5) mul(6,7");
    check::equal(malformed.part1, 0LL);

    check::throws([] { day3::solve(" \n"); });
    return check::exit_code();
}
