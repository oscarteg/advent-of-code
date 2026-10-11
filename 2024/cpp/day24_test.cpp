#include "check.hpp"
#include "day24.hpp"

#include <string>
#include <utility>
#include <vector>

namespace {

using Swap = std::pair<std::string, std::string>;

std::string wire(char prefix, int bit) {
    return std::string(1, prefix) + (bit < 10 ? "0" : "") + std::to_string(bit);
}

std::string swapped(const std::string& name, const std::vector<Swap>& swaps) {
    for (const auto& [left, right] : swaps) {
        if (name == left) return right;
        if (name == right) return left;
    }
    return name;
}

// Builds a correct ripple-carry adder for x + y, then renames the gate outputs listed in
// `swaps` so the circuit is wrong in exactly the way part 2 describes.
std::string adderInput(int bits, long long x, long long y, const std::vector<Swap>& swaps) {
    std::string text;
    for (int bit = 0; bit < bits; ++bit) {
        text += wire('x', bit) + ": " + std::to_string((x >> bit) & 1) + '\n';
    }
    for (int bit = 0; bit < bits; ++bit) {
        text += wire('y', bit) + ": " + std::to_string((y >> bit) & 1) + '\n';
    }
    text += '\n';

    const auto gate = [&](const std::string& left, const std::string& op, const std::string& right,
                          const std::string& output) {
        text += left + ' ' + op + ' ' + right + " -> " + swapped(output, swaps) + '\n';
    };
    gate(wire('x', 0), "XOR", wire('y', 0), wire('z', 0));
    gate(wire('x', 0), "AND", wire('y', 0), bits == 1 ? wire('z', 1) : wire('c', 0));
    for (int bit = 1; bit < bits; ++bit) {
        const std::string carryIn = wire('c', bit - 1);
        gate(wire('x', bit), "XOR", wire('y', bit), wire('s', bit));
        gate(wire('x', bit), "AND", wire('y', bit), wire('a', bit));
        gate(wire('s', bit), "XOR", carryIn, wire('z', bit));
        gate(wire('s', bit), "AND", carryIn, wire('p', bit));
        gate(wire('p', bit), "OR", wire('a', bit),
             bit + 1 == bits ? wire('z', bits) : wire('c', bit));
    }
    return text;
}

}  // namespace

int main() {
    const auto small = day24::solve(R"(x00: 1
x01: 1
x02: 1
y00: 0
y01: 1
y02: 0

x00 AND y00 -> z00
x01 XOR y01 -> z01
x02 OR y02 -> z02
)");
    check::equal(small.part1, 4LL);
    check::equal(small.part2, std::string{});

    const auto larger = day24::solve(R"(x00: 1
x01: 0
x02: 1
x03: 1
x04: 0
y00: 1
y01: 1
y02: 1
y03: 1
y04: 1

ntg XOR fgs -> mjb
y02 OR x01 -> tnw
kwq OR kpj -> z05
x00 OR x03 -> fst
tgd XOR rvg -> z01
vdt OR tnw -> bfw
bfw AND frj -> z10
ffh OR nrd -> bqk
y00 AND y03 -> djm
y03 OR y00 -> psh
bqk OR frj -> z08
tnw OR fst -> frj
gnj AND tgd -> z11
bfw XOR mjb -> z00
x03 OR x00 -> vdt
gnj AND wpb -> z02
x04 AND y00 -> kjc
djm OR pbm -> qhw
nrd AND vdt -> hwm
kjc AND fst -> rvg
y04 OR y02 -> fgs
y01 AND x02 -> pbm
ntg OR kjc -> kwq
psh XOR fgs -> tgd
qhw XOR tgd -> z09
pbm OR djm -> kpj
x03 XOR y03 -> ffh
x00 XOR y04 -> ntg
bfw OR bqk -> z06
nrd XOR fgs -> wpb
frj XOR qhw -> z04
bqk AND frj -> z07
y03 OR x01 -> nrd
hwm AND bqk -> z03
tgd XOR rvg -> z12
tnw OR pbm -> gnj
)");
    check::equal(larger.part1, 2024LL);

    const auto correct = day24::solve(adderInput(8, 0b1011'0101, 0b0110'1110, {}));
    check::equal(correct.part1, 0b1011'0101LL + 0b0110'1110LL);
    check::equal(correct.part2, std::string{});

    const std::vector<Swap> swaps{{"s01", "a01"}, {"z03", "a03"}, {"z05", "c05"}, {"s06", "a06"}};
    const auto crossed = day24::solve(adderInput(8, 0b1011'0101, 0b0110'1110, swaps));
    check::equal(crossed.part2, std::string{"a01,a03,a06,c05,s01,s06,z03,z05"});

    // Edge case: the smallest possible adder, a half adder whose carry is the top bit.
    const auto halfAdder = day24::solve(adderInput(1, 1, 1, {}));
    check::equal(halfAdder.part1, 2LL);
    check::equal(halfAdder.part2, std::string{});

    check::throws([] { return day24::solve(""); });
    check::throws([] { return day24::solve("x00: 1\ny00: 1\n"); });
    check::throws([] { return day24::solve("x00: 2\ny00: 1\n\nx00 AND y00 -> z00\n"); });
    check::throws([] { return day24::solve("x00: 1\ny00: 1\n\nx00 NAND y00 -> z00\n"); });
    check::throws([] { return day24::solve("x00: 1\ny00: 1\n\nx00 AND y00 z00\n"); });
    check::throws([] { return day24::solve("x00: 1\ny00: 1\n\nx00 XOR -> z00\n"); });
    check::throws([] { return day24::solve("x00: 1\n\nx00 AND y00 -> z00\n"); });
    check::throws([] {
        return day24::solve("x00: 1\ny00: 1\n\nx00 AND y00 -> z00\nx00 OR y00 -> z00\n");
    });
    check::throws([] { return day24::solve("y00: 1\n\nz01 AND y00 -> z00\nz00 OR y00 -> z01\n"); });

    return check::exit_code();
}
