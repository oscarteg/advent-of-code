#include "check.hpp"
#include "day10.hpp"

int main() {
    std::string_view example = R"(89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732
)";

    {
        auto [p1, p2] = day10::solve(example);
        check::equal(p1, 36LL);
        check::equal(p2, 81LL);
    }

    {
        std::string_view single = "0123456789\n";
        auto [p1, p2] = day10::solve(single);
        check::equal(p1, 1LL);
        check::equal(p2, 1LL);
    }

    {
        std::string_view vertical = "0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n";
        auto [p1, p2] = day10::solve(vertical);
        check::equal(p1, 1LL);
        check::equal(p2, 1LL);
    }

    check::throws([] { day10::solve("0a1\n"); });
    check::throws([] { day10::solve("012\n01\n"); });
    check::throws([] { day10::solve(""); });

    return check::exit_code();
}
