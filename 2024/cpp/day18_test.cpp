#include "check.hpp"
#include "day18.hpp"

int main() {
    constexpr std::string_view example =
        "5,4\n4,2\n4,5\n3,0\n2,1\n6,3\n2,4\n1,5\n0,6\n3,3\n2,6\n5,1\n1,2\n5,5\n2,5\n6,5\n1,4\n0,4\n6,4\n1,1\n6,1\n1,0\n"
        "0,5\n1,6\n2,0\n";

    const auto answers = day18::solve(example, 7, 12);
    check::equal(answers.part1, 22);
    check::equal(answers.part2, "6,1");

    constexpr std::string_view twoByteGrid = "1,0\n0,1\n";
    const auto trivial = day18::solve(twoByteGrid, 2, 0);
    check::equal(trivial.part1, 2);
    check::equal(trivial.part2, "0,1");

    check::throws([] { day18::solve("not a coordinate\n", 7, 1); });
    check::throws([&] { day18::solve(example, 7, 1000); });

    return check::exit_code();
}
