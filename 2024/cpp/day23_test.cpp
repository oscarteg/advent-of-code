#include "check.hpp"
#include "day23.hpp"

#include <string>

int main() {
    const std::string example = R"(kh-tc
qp-kh
de-cg
ka-co
yn-aq
qp-ub
cg-tb
vc-aq
tb-ka
wh-tc
yn-cg
kh-ub
ta-co
de-co
tc-td
tb-wq
wh-td
ta-ka
td-qp
aq-cg
wq-ub
ub-vc
de-ta
wq-aq
wq-vc
wh-yn
ka-de
kh-ta
co-tc
wh-qp
tb-vc
td-yn)";

    const auto answers = day23::solve(example);
    check::equal(answers.part1, 7LL);
    check::equal(answers.part2, std::string("co,de,ka,ta"));

    const std::string single = "a-b\nb-c\na-c";
    const auto single_ans = day23::solve(single);
    check::equal(single_ans.part1, 0LL);
    check::equal(single_ans.part2, std::string("a,b,c"));

    const std::string with_t = "a-b\nb-t\na-t";
    const auto t_ans = day23::solve(with_t);
    check::equal(t_ans.part1, 1LL);
    check::equal(t_ans.part2, std::string("a,b,t"));

    check::throws([]() { day23::solve("invalid"); });

    return check::exit_code();
}
