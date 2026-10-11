#include "check.hpp"
#include "day21.hpp"

int main() {
    constexpr std::string_view example = R"(029A
980A
179A
456A
379A
)";

    const auto answers = day21::solve(example);
    check::equal(answers.part1, 126384);
    check::equal(answers.part2, 154115708116294LL);

    // Published shortest sequence lengths for the five example codes.
    check::equal(day21::solve("029A\n").part1, 68 * 29);
    check::equal(day21::solve("980A\n").part1, 60 * 980);
    check::equal(day21::solve("179A\n").part1, 68 * 179);
    check::equal(day21::solve("456A\n").part1, 64 * 456);
    check::equal(day21::solve("379A\n").part1, 64 * 379);

    // Leading zeroes are not part of the numeric part, so this code has complexity zero.
    const auto zeroes = day21::solve("000A\n");
    check::equal(zeroes.part1, 0);
    check::equal(zeroes.part2, 0);

    // Codes are independent, so a list is worth the sum of its codes.
    const auto two = day21::solve("029A\n980A\n");
    check::equal(two.part1, day21::solve("029A\n").part1 + day21::solve("980A\n").part1);
    check::equal(two.part2, day21::solve("029A\n").part2 + day21::solve("980A\n").part2);

    // Carriage returns and blank lines are ignored.
    check::equal(day21::solve("029A\r\n\r\n980A\r\n").part1, two.part1);

    check::throws([] { return day21::solve(""); });            // no codes
    check::throws([] { return day21::solve("   \n"); });       // still no codes
    check::throws([] { return day21::solve("029\n"); });       // no A to submit the code
    check::throws([] { return day21::solve("02XA\n"); });      // not a digit
    check::throws([] { return day21::solve("A\n"); });         // no digits
    check::throws([] { return day21::solve("-29A\n"); });      // signed numeric part
    check::throws([] { return day21::solve("029a\n"); });      // lowercase submit key

    return check::exit_code();
}
