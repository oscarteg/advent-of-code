#include "check.hpp"
#include "day20.hpp"

int main() {
    constexpr std::string_view example = R"(###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############
)";

    // Published counts for 2-picosecond cheats on the example map.
    check::equal(day20::solve(example, 64).part1, 1);
    check::equal(day20::solve(example, 40).part1, 2);
    check::equal(day20::solve(example, 20).part1, 5);
    check::equal(day20::solve(example, 2).part1, 44);

    // Published counts for 20-picosecond cheats on the same map.
    check::equal(day20::solve(example, 76).part2, 3);
    check::equal(day20::solve(example, 74).part2, 7);
    check::equal(day20::solve(example, 72).part2, 29);
    check::equal(day20::solve(example, 50).part2, 285);

    // The best cheat on the example saves 76, so the real threshold counts nothing.
    const auto atRealThreshold = day20::solve(example);
    check::equal(atRealThreshold.part1, 0);
    check::equal(atRealThreshold.part2, 0);

    // A U-shaped track: cheating straight through the middle wall column saves 4
    // picoseconds from the start and 2 from the cell below it.
    constexpr std::string_view corridor = R"(#####
#S#E#
#.#.#
#...#
#####
)";
    const auto corridorCheats = day20::solve(corridor, 2);
    check::equal(corridorCheats.part1, 2);
    check::equal(corridorCheats.part2, 4);
    const auto bestCorridorCheats = day20::solve(corridor, 4);
    check::equal(bestCorridorCheats.part1, 1);
    check::equal(bestCorridorCheats.part2, 1);

    check::throws([] { return day20::solve(""); });
    check::throws([] { return day20::solve("###\n#.#\n###\n"); });            // no start or end
    check::throws([] { return day20::solve("#####\n#S#E#\n#####\n"); });      // end unreachable
    check::throws([] { return day20::solve("#####\n#SSE#\n#####\n"); });      // two starts
    check::throws([] { return day20::solve("#####\n#S.E#\n####\n"); });       // ragged rows
    check::throws([] { return day20::solve("#####\n#S?E#\n#####\n"); });      // unknown character

    return check::exit_code();
}
