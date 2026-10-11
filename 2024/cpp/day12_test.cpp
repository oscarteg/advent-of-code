#include "check.hpp"
#include "day12.hpp"

namespace {

constexpr std::string_view kFiveRegions = R"(AAAA
BBCD
BBCC
EEEC
)";

constexpr std::string_view kNestedRegions = R"(OOOOO
OXOXO
OOOOO
OXOXO
OOOOO
)";

constexpr std::string_view kLargeExample = R"(RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE
)";

// Part 2 only: the puzzle states the side-counting price for these two maps.
constexpr std::string_view kEShape = R"(EEEEE
EXXXX
EEEEE
EXXXX
EEEEE
)";

constexpr std::string_view kDiagonallyTouching = R"(AAAAAA
AAABBA
AAABBA
ABBAAA
ABBAAA
AAAAAA
)";

constexpr std::string_view kSinglePlot = "A\n";

constexpr std::string_view kRaggedGarden = R"(AAAA
BB
)";

constexpr std::string_view kNonLetterPlot = R"(AA
A1
)";

constexpr std::string_view kEmptyInput = "";

}  // namespace

int main() {
    const auto fiveRegions = day12::solve(kFiveRegions);
    check::equal(fiveRegions.part1, 140);
    check::equal(fiveRegions.part2, 80);

    const auto nested = day12::solve(kNestedRegions);
    check::equal(nested.part1, 772);
    check::equal(nested.part2, 436);

    const auto large = day12::solve(kLargeExample);
    check::equal(large.part1, 1930);
    check::equal(large.part2, 1206);

    check::equal(day12::solve(kEShape).part2, 236);
    check::equal(day12::solve(kDiagonallyTouching).part2, 368);

    const auto single = day12::solve(kSinglePlot);
    check::equal(single.part1, 4);
    check::equal(single.part2, 4);

    check::throws([] { day12::solve(kRaggedGarden); });
    check::throws([] { day12::solve(kNonLetterPlot); });
    check::throws([] { day12::solve(kEmptyInput); });

    return check::exit_code();
}
