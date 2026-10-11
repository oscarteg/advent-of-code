#include "check.hpp"
#include "day4.hpp"

namespace {

constexpr std::string_view kExample =
    "MMMSXXMASM\n"
    "MSAMXMSMSA\n"
    "AMXSXMAAMM\n"
    "MSAMASMSMX\n"
    "XMASAMXAMM\n"
    "XXAMMXXAMA\n"
    "SMSMSASXSS\n"
    "SAXAMASAAA\n"
    "MAMMMXMMMM\n"
    "MXMXAXMASX\n";

constexpr std::string_view kSingleDiagonal =
    "M.S\n"
    ".A.\n"
    "M.S\n";

constexpr std::string_view kRaggedGrid =
    "XMAS\n"
    "XM\n";

constexpr std::string_view kEmptyInput = "";

}  // namespace

int main() {
    const auto example = day4::solve(kExample);
    check::equal(example.part1, 18);
    check::equal(example.part2, 9);

    const auto single = day4::solve(kSingleDiagonal);
    check::equal(single.part2, 1);

    check::throws([] { day4::solve(kRaggedGrid); });
    check::throws([] { day4::solve(kEmptyInput); });

    return check::exit_code();
}
