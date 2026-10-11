#include "check.hpp"
#include "day13.hpp"

#include <stdexcept>

int main() {
    {
        const auto result = day13::solve(R"(Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279
)");
        check::equal(result.part1, 480LL, std::source_location::current());
    }

    {
        const auto result = day13::solve(R"(Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279
)");
        check::equal(result.part2, 875318608908LL, std::source_location::current());
    }

    {
        const auto result = day13::solve(R"(Button A: X+1, Y+1
Button B: X+2, Y+2
Prize: X=3, Y=4
)");
        check::equal(result.part1, 0LL, std::source_location::current());
    }

    {
        const auto result = day13::solve(R"(Button A: X+1, Y+0
Button B: X+0, Y+1
Prize: X=3, Y=2
)");
        check::equal(result.part1, 11LL, std::source_location::current());
    }

    try {
        day13::solve(R"(Button A: X+1, Y+1
Prize: X=3, Y=4
)");
        return 1;
    } catch (const std::invalid_argument&) {
    }

    try {
        day13::solve("");
        return 1;
    } catch (const std::invalid_argument&) {
    }

    return check::exit_code();
}
