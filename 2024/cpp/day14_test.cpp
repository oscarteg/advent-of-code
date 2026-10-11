#include "check.hpp"
#include "day14.hpp"

#include <stdexcept>

int main() {
    {
        const auto result = day14::solve_with(R"(p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3
)",
                                               11, 7);
        check::equal(result.part1, 12LL, std::source_location::current());
    }

    {
        const auto robots = day14::parse("p=0,0 v=0,0\np=1,0 v=0,0\n");
        check::equal(day14::find_picture_second(robots, 3, 3), 0LL, std::source_location::current());
    }

    {
        const auto robots = day14::parse("p=0,0 v=1,0\np=0,0 v=2,0\n");
        check::equal(day14::find_picture_second(robots, 5, 1), 1LL, std::source_location::current());
    }

    {
        const auto robots = day14::parse("p=0,0 v=0,0\n");
        check::equal(day14::positions_unique({{0, 0}, {0, 0}}, 1), false, std::source_location::current());
        check::equal(day14::positions_unique({{0, 0}, {1, 0}}, 1), true, std::source_location::current());
        (void)robots;
    }

    check::throws([] { day14::solve(""); });
    check::throws([] { day14::solve("garbage line\n"); });

    return check::exit_code();
}
