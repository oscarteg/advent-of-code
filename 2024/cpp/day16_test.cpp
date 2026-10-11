#include "check.hpp"
#include "day16.hpp"

namespace {

constexpr std::string_view kFirstExample = R"(###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############
)";

constexpr std::string_view kSecondExample = R"(#################
#...#...#...#..E#
#.#.#.#.#.#.#.#.#
#.#.#.#...#...#.#
#.#.#.#.###.#.#.#
#...#.#.#.....#.#
#.#.#.#.#.#####.#
#.#...#.#.#.....#
#.#.#####.#.###.#
#.#.#.......#...#
#.#.###.#####.###
#.#.#...#.....#.#
#.#.#.#####.###.#
#.#.#.........#.#
#.#.#.#########.#
#S#.............#
#################
)";

// Facing east at the start, so the end due north costs one turn plus two steps.
constexpr std::string_view kCorridor = R"(###
#E#
#.#
#S#
###
)";

constexpr std::string_view kAdjacentEnd = R"(####
#SE#
####
)";

constexpr std::string_view kRaggedRows = R"(####
#SE#
###
)";

constexpr std::string_view kUnknownTile = R"(####
#SX#
####
)";

constexpr std::string_view kWithoutStart = R"(####
#.E#
####
)";

constexpr std::string_view kWithoutEnd = R"(####
#S.#
####
)";

constexpr std::string_view kWalledOffEnd = R"(#####
#S#E#
#####
)";

constexpr std::string_view kEmptyInput = "";

}  // namespace

int main() {
    const auto first = day16::solve(kFirstExample);
    check::equal(first.part1, 7036);
    check::equal(first.part2, 45);

    const auto second = day16::solve(kSecondExample);
    check::equal(second.part1, 11048);
    check::equal(second.part2, 64);

    const auto corridor = day16::solve(kCorridor);
    check::equal(corridor.part1, 1002);
    check::equal(corridor.part2, 3);

    const auto adjacent = day16::solve(kAdjacentEnd);
    check::equal(adjacent.part1, 1);
    check::equal(adjacent.part2, 2);

    check::throws([] { day16::solve(kRaggedRows); });
    check::throws([] { day16::solve(kUnknownTile); });
    check::throws([] { day16::solve(kWithoutStart); });
    check::throws([] { day16::solve(kWithoutEnd); });
    check::throws([] { day16::solve(kWalledOffEnd); });
    check::throws([] { day16::solve(kEmptyInput); });

    return check::exit_code();
}
