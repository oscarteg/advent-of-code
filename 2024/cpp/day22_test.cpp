#include "check.hpp"
#include "day22.hpp"

int main() {
    const auto example = day22::solve("1\n10\n100\n2024\n");
    check::equal(example.part1, 37327623LL);

    const auto bananas = day22::solve("1\n2\n3\n2024\n");
    check::equal(bananas.part2, 23LL);

    // Edge case: a single seed's 2000th secret matches the official per-seed example.
    const auto single = day22::solve("1\n");
    check::equal(single.part1, 8685429LL);

    check::throws([] { day22::solve(""); });
    check::throws([] { day22::solve("not-a-number\n"); });

    return check::exit_code();
}
