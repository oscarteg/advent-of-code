#pragma once

#include <algorithm>
#include <charconv>
#include <cstdint>
#include <stdexcept>
#include <string>
#include <string_view>
#include <unordered_map>
#include <vector>

namespace day24 {

struct Answers {
    long long part1;
    std::string part2;
};

namespace detail {

enum class Op : std::uint8_t { And, Or, Xor };

using WireId = int;

constexpr WireId kNoWire = -1;
constexpr std::size_t kNoGate = static_cast<std::size_t>(-1);
constexpr int kMaxSwaps = 4;
constexpr int kMaxBits = 63;

struct Gate {
    WireId left;
    WireId right;
    Op op;
};

struct Circuit {
    std::vector<std::string> names;
    std::unordered_map<std::string, WireId> ids;
    std::vector<signed char> initialValues;  // per wire: -1 unknown, else 0 or 1
    std::vector<std::size_t> driverByWire;   // per wire: gate index, or kNoGate
    std::vector<Gate> gates;
    std::vector<WireId> outputs;             // outputs[i] is the wire gate i drives
    int bits = 0;                            // number of x (and y) input wires
};

inline std::string_view trim(std::string_view text) {
    while (!text.empty() && (text.front() == ' ' || text.front() == '\t')) text.remove_prefix(1);
    while (!text.empty() && (text.back() == ' ' || text.back() == '\t' || text.back() == '\r')) {
        text.remove_suffix(1);
    }
    return text;
}

inline std::vector<std::string_view> splitLines(std::string_view input) {
    std::vector<std::string_view> lines;
    while (!input.empty()) {
        const std::size_t end = input.find('\n');
        if (end == std::string_view::npos) {
            lines.push_back(input);
            break;
        }
        lines.push_back(input.substr(0, end));
        input.remove_prefix(end + 1);
    }
    return lines;
}

inline std::vector<std::string_view> splitWords(std::string_view line) {
    std::vector<std::string_view> words;
    std::size_t position = 0;
    while (position < line.size()) {
        while (position < line.size() && line[position] == ' ') ++position;
        const std::size_t start = position;
        while (position < line.size() && line[position] != ' ') ++position;
        if (position > start) words.push_back(line.substr(start, position - start));
    }
    return words;
}

inline Op parseOp(std::string_view text) {
    if (text == "AND") return Op::And;
    if (text == "OR") return Op::Or;
    if (text == "XOR") return Op::Xor;
    throw std::invalid_argument("unknown gate operator: " + std::string{text});
}

inline std::string wireName(char prefix, int bit) {
    std::string name(1, prefix);
    if (bit < 10) name += '0';
    name += std::to_string(bit);
    return name;
}

inline int parseBitIndex(const std::string& wire) {
    const std::string_view digits{wire.data() + 1, wire.size() - 1};
    int bit = 0;
    const auto result = std::from_chars(digits.data(), digits.data() + digits.size(), bit);
    if (result.ec != std::errc{} || result.ptr != digits.data() + digits.size() || bit > kMaxBits) {
        throw std::invalid_argument("wire has no usable bit index: " + wire);
    }
    return bit;
}

inline WireId intern(Circuit& circuit, std::string_view name) {
    const auto [entry, inserted] =
        circuit.ids.try_emplace(std::string{name}, static_cast<WireId>(circuit.names.size()));
    if (inserted) {
        circuit.names.emplace_back(name);
        circuit.initialValues.push_back(-1);
        circuit.driverByWire.push_back(kNoGate);
    }
    return entry->second;
}

inline Circuit parse(std::string_view input) {
    Circuit circuit;

    for (const std::string_view raw : splitLines(input)) {
        const std::string_view line = trim(raw);
        if (line.empty()) continue;

        if (line.find("->") != std::string_view::npos) {
            const auto words = splitWords(line);
            if (words.size() != 5 || words[3] != "->") {
                throw std::invalid_argument("malformed gate: " + std::string{line});
            }
            const Op op = parseOp(words[1]);
            const WireId left = intern(circuit, words[0]);
            const WireId right = intern(circuit, words[2]);
            const WireId output = intern(circuit, words[4]);
            if (circuit.driverByWire[output] != kNoGate) {
                throw std::invalid_argument("wire is driven twice: " + circuit.names[output]);
            }
            circuit.driverByWire[output] = circuit.gates.size();
            circuit.gates.push_back({.left = left, .right = right, .op = op});
            circuit.outputs.push_back(output);
            continue;
        }

        const std::size_t colon = line.find(':');
        if (colon == std::string_view::npos) {
            throw std::invalid_argument("malformed line: " + std::string{line});
        }
        const std::string_view name = trim(line.substr(0, colon));
        const std::string_view value = trim(line.substr(colon + 1));
        if (name.empty() || (value != "0" && value != "1")) {
            throw std::invalid_argument("malformed initial value: " + std::string{line});
        }
        circuit.initialValues[intern(circuit, name)] = value == "1" ? 1 : 0;
    }

    const bool hasOutputWire = std::ranges::any_of(
        circuit.outputs, [&](WireId wire) { return circuit.names[wire].starts_with('z'); });
    if (!hasOutputWire) throw std::invalid_argument("input has no z output wires");

    while (circuit.bits < kMaxBits && circuit.ids.contains(wireName('x', circuit.bits))) {
        ++circuit.bits;
    }
    return circuit;
}

struct Evaluator {
    const Circuit& circuit;
    std::vector<signed char> values;
    std::vector<bool> visiting;

    bool valueOf(WireId wire) {
        if (values[wire] >= 0) return values[wire] == 1;
        const std::size_t driver = circuit.driverByWire[wire];
        if (driver == kNoGate) {
            throw std::invalid_argument("wire never receives a value: " + circuit.names[wire]);
        }
        if (visiting[wire]) {
            throw std::invalid_argument("circuit has a cycle at wire: " + circuit.names[wire]);
        }
        visiting[wire] = true;
        const Gate& gate = circuit.gates[driver];
        const bool left = valueOf(gate.left);
        const bool right = valueOf(gate.right);
        visiting[wire] = false;

        const bool value = gate.op == Op::And ? left && right
                           : gate.op == Op::Or ? left || right
                                               : left != right;
        values[wire] = value ? 1 : 0;
        return value;
    }
};

inline long long simulate(const Circuit& circuit) {
    Evaluator evaluator{.circuit = circuit,
                        .values = circuit.initialValues,
                        .visiting = std::vector<bool>(circuit.names.size(), false)};

    long long result = 0;
    for (const WireId wire : circuit.outputs) {
        const std::string& name = circuit.names[wire];
        if (!name.starts_with('z')) continue;
        if (evaluator.valueOf(wire)) result |= 1LL << parseBitIndex(name);
    }
    return result;
}

inline std::uint64_t gateKey(WireId left, WireId right, Op op) {
    const auto low = static_cast<std::uint64_t>(std::min(left, right));
    const auto high = static_cast<std::uint64_t>(std::max(left, right));
    return (low << 32) | (high << 2) | static_cast<std::uint64_t>(op);
}

using GateLookup = std::unordered_map<std::uint64_t, std::size_t>;

inline GateLookup buildLookup(const Circuit& circuit) {
    GateLookup lookup;
    for (std::size_t index = 0; index < circuit.gates.size(); ++index) {
        const Gate& gate = circuit.gates[index];
        lookup.try_emplace(gateKey(gate.left, gate.right, gate.op), index);
    }
    return lookup;
}

// A correct ripple-carry adder computes z0 = x0 XOR y0 and carry0 = x0 AND y0, then for
// every higher bit z_i = (x_i XOR y_i) XOR carry_(i-1) and
// carry_i = (x_i AND y_i) OR ((x_i XOR y_i) AND carry_(i-1)), with the last carry on the
// top z wire. Returns the lowest bit that does not match, or bits + 1 when all of them do.
inline int firstBrokenBit(const Circuit& circuit, const GateLookup& lookup,
                          const std::vector<WireId>& outputs) {
    const auto wireId = [&circuit](char prefix, int bit) {
        const auto entry = circuit.ids.find(wireName(prefix, bit));
        return entry == circuit.ids.end() ? kNoWire : entry->second;
    };
    const auto driven = [&](WireId left, WireId right, Op op) {
        if (left == kNoWire || right == kNoWire) return kNoWire;
        const auto entry = lookup.find(gateKey(left, right, op));
        return entry == lookup.end() ? kNoWire : outputs[entry->second];
    };

    WireId carry = kNoWire;
    for (int bit = 0; bit < circuit.bits; ++bit) {
        const WireId output = wireId('z', bit);
        const WireId sum = driven(wireId('x', bit), wireId('y', bit), Op::Xor);
        const WireId generated = driven(wireId('x', bit), wireId('y', bit), Op::And);
        if (output == kNoWire || sum == kNoWire || generated == kNoWire) return bit;

        if (bit == 0) {
            if (sum != output) return 0;
            carry = generated;
            continue;
        }
        if (driven(sum, carry, Op::Xor) != output) return bit;
        const WireId propagated = driven(sum, carry, Op::And);
        if (propagated == kNoWire) return bit;
        const WireId nextCarry = driven(propagated, generated, Op::Or);
        if (nextCarry == kNoWire) return bit;
        carry = nextCarry;
    }

    const WireId topBit = wireId('z', circuit.bits);
    if (topBit == kNoWire || carry != topBit) return circuit.bits;
    return circuit.bits + 1;
}

// Depth-first search over output swaps, keeping only swaps that push the first broken bit
// higher; the puzzle input needs exactly four of them.
inline bool repair(const Circuit& circuit, const GateLookup& lookup, std::vector<WireId>& outputs,
                   int remaining, std::vector<WireId>& swapped) {
    const int broken = firstBrokenBit(circuit, lookup, outputs);
    if (broken > circuit.bits) return true;
    if (remaining == 0) return false;

    for (std::size_t a = 0; a + 1 < outputs.size(); ++a) {
        for (std::size_t b = a + 1; b < outputs.size(); ++b) {
            std::swap(outputs[a], outputs[b]);
            if (firstBrokenBit(circuit, lookup, outputs) > broken) {
                swapped.push_back(outputs[a]);
                swapped.push_back(outputs[b]);
                if (repair(circuit, lookup, outputs, remaining - 1, swapped)) return true;
                swapped.pop_back();
                swapped.pop_back();
            }
            std::swap(outputs[a], outputs[b]);
        }
    }
    return false;
}

// Empty when the circuit is not an adder that four swaps can repair, which includes both
// part 1 examples.
inline std::string findSwappedWires(const Circuit& circuit) {
    const GateLookup lookup = buildLookup(circuit);
    std::vector<WireId> outputs = circuit.outputs;
    std::vector<WireId> swapped;
    if (!repair(circuit, lookup, outputs, kMaxSwaps, swapped)) return {};

    std::vector<std::string> names;
    names.reserve(swapped.size());
    for (const WireId wire : swapped) names.push_back(circuit.names[wire]);
    std::ranges::sort(names);

    std::string joined;
    for (const std::string& name : names) {
        if (!joined.empty()) joined += ',';
        joined += name;
    }
    return joined;
}

}  // namespace detail

inline Answers solve(std::string_view input) {
    const detail::Circuit circuit = detail::parse(input);
    return {.part1 = detail::simulate(circuit), .part2 = detail::findSwappedWires(circuit)};
}

}  // namespace day24
