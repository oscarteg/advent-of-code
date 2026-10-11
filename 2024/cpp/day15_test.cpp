#include "check.hpp"
#include "day15.hpp"

namespace {

constexpr std::string_view smallExample = R"(########
#..O.O.#
##@.O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

<^^>>>vv<v>>v<<
)";

constexpr std::string_view largeExample = R"(##########
#..O..O.O#
#......O.#
#.OO..O.O#
#..O@..O.#
#O#..O...#
#O..O..O.#
#.OO.O.OO#
#....O...#
##########

<vv>^<v^>v>^vv^v>v<>v^v<v<^vv<<<^><<><>>v<vvv<>^v^>^<<<><<v<<<v^vv^v>^
vvv<<^>^v^^><<>>><>^<<><^vv^^<>vvv<>><^^v>^>vv<>v<<<<v<^v>^<^^>>>^<v<v
><>vv>v^v^<>><>>>><^^>vv>v<^^^>>v^v^<^^>v^^>v^<^v>v<>>v^v^<v>v^^<^^vv<
<<v<^>>^^^^>>>v^<>vvv^><v<<<>^^^vv^<vvv>^>v<^^^^v<>^>vvvv><>>v^<<^^^^^
^><^><>>><>^^<<^^v>>><^<v>^<vv>>v>>>^v><>^v><<<<v>>v<v<v>vvv>^<><<>^><
^>><>^v<><^vvv<^^<><v<<<<<><^v<<<><<<^^<v<^^^><^>>^<v^><<<^>>^v<v^v<v^
>^>>^v>vv>^<<^v<>><<><<v<<v><>v<^vv<<<>^^v^>^^>>><<^v>>v^v><^^>>^<>vv^
<><^^>^^^<><vvvvv^v<v<<>^v<v>v<<^><<><<><<<^^<<<^<<>><<><^^^>^^<>^>v<>
^^>vv<^v^v<vv>^<><v<^v>^^^>>>^^vvv^>vvv<>>>^<^>>>>>^<<^v>^vvv<>^<><<v>
v^^>>><<^^<>>^v^<v^vv<>v^<<>^<^v^v><^<<<><<^<v><v<>vv>>v><v^<vv<>v^<<^
)";

constexpr std::string_view verticalPushExample = R"(#######
#...#.#
#.....#
#..OO@#
#..O..#
#.....#
#######

<vv<<^^<<^^
)";

constexpr std::string_view blockedPush = R"(#####
#.O@#
#####

<<
)";

}  // namespace

int main() {
    // The puzzle states only the narrow sum for the smaller map.
    const auto small = day15::solve(smallExample);
    check::equal(small.part1, 2028);

    const auto large = day15::solve(largeExample);
    check::equal(large.part1, 10092);
    check::equal(large.part2, 9021);

    // The puzzle states only the widened sum for this map; its boxes end up
    // half-aligned, so a vertical push moves a growing fan of boxes.
    const auto vertical = day15::solve(verticalPushExample);
    check::equal(vertical.part2, 618);

    // Edge case: the second push runs into the wall and leaves the narrow map
    // untouched, while the widened map has one floor tile more and takes both.
    const auto blocked = day15::solve(blockedPush);
    check::equal(blocked.part1, 101);
    check::equal(blocked.part2, 102);

    check::throws([] { day15::solve(""); });
    check::throws([] { day15::solve("#####\n###\n\n<\n"); });
    check::throws([] { day15::solve("#####\n#.X@#\n#####\n\n<\n"); });
    check::throws([] { day15::solve("#####\n#...#\n#####\n\n<\n"); });
    check::throws([] { day15::solve("#####\n#@@.#\n#####\n\n<\n"); });
    check::throws([] { day15::solve("#####\n#.@..\n#####\n\n<\n"); });
    check::throws([] { day15::solve("#####\n#.@.#\n#####\n"); });
    check::throws([] { day15::solve("#####\n#.@.#\n#####\n\n<x\n"); });

    return check::exit_code();
}
