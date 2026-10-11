#include "check.hpp"
#include "day17.hpp"

namespace {

constexpr auto kExample = R"(Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0
)";

// Part 2's example: this program prints itself when A starts at 117440.
constexpr auto kSelfReplicating = R"(Register A: 2024
Register B: 0
Register C: 0

Program: 0,3,5,4,3,0
)";

}  // namespace

int main() {
    const auto example = day17::solve(kExample);
    check::equal(example.part1, "4,6,3,5,6,3,5,2,1,0");
    check::equal(example.part2, 0LL);

    const auto selfReplicating = day17::solve(kSelfReplicating);
    check::equal(selfReplicating.part1, "5,7,3,0");
    check::equal(selfReplicating.part2, 117440LL);

    // Single-instruction examples from the puzzle's instruction list.
    check::equal(day17::solve("Register A: 10\nRegister B: 0\nRegister C: 0\n\nProgram: 5,0,5,1,5,4\n").part1,
                 "0,1,2");
    check::equal(day17::solve("Register A: 2024\nRegister B: 0\nRegister C: 0\n\nProgram: 0,1,5,4,3,0\n").part1,
                 "4,2,5,6,7,7,7,7,3,1,0");

    // Edge case: a program without an out instruction prints nothing.
    check::equal(day17::solve("Register A: 0\nRegister B: 29\nRegister C: 0\n\nProgram: 1,7\n").part1, "");

    check::throws([] { return day17::solve(""); });
    check::throws([] { return day17::solve("Register A: 1\nRegister B: 0\nRegister C: 0\n"); });
    check::throws([] { return day17::solve("Register A: x\nRegister B: 0\nRegister C: 0\n\nProgram: 0,1\n"); });
    check::throws([] { return day17::solve("Register A: 1\nRegister B: 0\nRegister C: 0\n\nProgram: 0,8\n"); });

    return check::exit_code();
}
