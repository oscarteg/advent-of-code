#include "check.hpp"
#include "day19.hpp"

int main() {
    {
        const std::string_view example = R"(r, wr, b, g, bwu, rb, gb, br
brwrr
bggr
gbbr
rrbgbr
ubwu
bwurrg
brgr
bbrgwb)";
        auto answers = day19::solve(example);
        check::equal(answers.part1, 6LL);
        check::equal(answers.part2, 16LL);
    }

    {
        const std::string_view single_pattern = R"(a
a
aa
aaa)";
        auto answers = day19::solve(single_pattern);
        check::equal(answers.part1, 3LL);
        check::equal(answers.part2, 3LL);
    }

    {
        const std::string_view impossible = R"(a, b
c)";
        auto answers = day19::solve(impossible);
        check::equal(answers.part1, 0LL);
        check::equal(answers.part2, 0LL);
    }

    check::throws([]() { day19::solve(""); });

    return check::exit_code();
}
