#include "check.hpp"
#include "day25.hpp"

int main() {
    const std::string example =
        "#####\n.####\n.####\n.####\n.#.#.\n.#...\n.....\n"
        "\n"
        "#####\n##.##\n.#.##\n...##\n...#.\n...#.\n.....\n"
        "\n"
        ".....\n#....\n#....\n#...#\n#.#.#\n#.###\n#####\n"
        "\n"
        ".....\n.....\n#.#..\n###..\n###.#\n###.#\n#####\n"
        "\n"
        ".....\n.....\n.....\n#....\n#.#..\n#.#.#\n#####\n";
    check::equal(day25::solve(example).part1, 3);

    // Every lock paired with itself as a key never fits: a lock's top row
    // mirrors a key's bottom row, so the shared full row guarantees overlap.
    const std::string noFit =
        "#####\n#####\n#####\n#####\n#####\n#####\n.....\n"
        "\n"
        ".....\n#####\n#####\n#####\n#####\n#####\n#####\n";
    check::equal(day25::solve(noFit).part1, 0);

    check::throws([] { day25::solve("not a schematic"); });
    check::throws([] { day25::solve("#####\n.####\n.....\n"); });

    return check::exit_code();
}
