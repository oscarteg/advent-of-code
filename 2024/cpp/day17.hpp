#pragma once

#include <algorithm>
#include <charconv>
#include <cstddef>
#include <span>
#include <stdexcept>
#include <string>
#include <string_view>
#include <vector>

namespace day17 {

struct Answers {
    std::string part1;
    // Lowest initial A that makes the program print a copy of itself, or 0 when
    // no such value exists (as for the part 1 example program).
    long long part2;
};

namespace detail {

struct Machine {
    long long a;
    long long b;
    long long c;
    std::vector<int> program;
};

inline std::string_view trim(std::string_view text) {
    while (!text.empty() && (text.front() == ' ' || text.front() == '\t')) text.remove_prefix(1);
    while (!text.empty() && (text.back() == ' ' || text.back() == '\t' || text.back() == '\r')) {
        text.remove_suffix(1);
    }
    return text;
}

inline long long parseNumber(std::string_view text) {
    long long value = 0;
    const auto* end = text.data() + text.size();
    const auto result = std::from_chars(text.data(), end, value);
    if (result.ec != std::errc{} || result.ptr != end || value < 0) {
        throw std::invalid_argument("invalid number: " + std::string{text});
    }
    return value;
}

inline int parseThreeBitValue(std::string_view text) {
    const long long value = parseNumber(text);
    if (value > 7) {
        throw std::invalid_argument("program values are 3 bits wide: " + std::string{text});
    }
    return static_cast<int>(value);
}

inline std::string_view valueAfter(std::string_view line, std::string_view label) {
    if (!line.starts_with(label)) {
        throw std::invalid_argument("expected a line starting with " + std::string{label});
    }
    return trim(line.substr(label.size()));
}

inline std::vector<std::string_view> nonEmptyLines(std::string_view input) {
    std::vector<std::string_view> lines;
    for (std::size_t start = 0; start < input.size();) {
        std::size_t end = input.find('\n', start);
        if (end == std::string_view::npos) end = input.size();
        const std::string_view line = trim(input.substr(start, end - start));
        start = end + 1;
        if (!line.empty()) lines.push_back(line);
    }
    return lines;
}

inline Machine parse(std::string_view input) {
    const auto lines = nonEmptyLines(input);
    if (lines.size() != 4) {
        throw std::invalid_argument("expected three register lines and one program line");
    }

    Machine machine{parseNumber(valueAfter(lines[0], "Register A:")),
                    parseNumber(valueAfter(lines[1], "Register B:")),
                    parseNumber(valueAfter(lines[2], "Register C:")),
                    {}};

    std::string_view rest = valueAfter(lines[3], "Program:");
    while (true) {
        const std::size_t comma = rest.find(',');
        const std::string_view field = comma == std::string_view::npos ? rest : rest.substr(0, comma);
        machine.program.push_back(parseThreeBitValue(trim(field)));
        if (comma == std::string_view::npos) break;
        rest = rest.substr(comma + 1);
    }
    return machine;
}

// Registers hold non-negative values, so dividing by 2^amount is a right shift.
inline long long shiftRight(long long value, long long amount) {
    return amount >= 63 ? 0 : value >> amount;
}

inline std::vector<int> run(const Machine& machine, long long a) {
    long long b = machine.b;
    long long c = machine.c;
    const auto& program = machine.program;
    std::vector<int> output;

    for (std::size_t ip = 0; ip + 1 < program.size();) {
        const int opcode = program[ip];
        const int operand = program[ip + 1];
        const auto combo = [&]() -> long long {
            switch (operand) {
                case 4: return a;
                case 5: return b;
                case 6: return c;
                case 7: throw std::invalid_argument("combo operand 7 is reserved");
                default: return operand;
            }
        };

        switch (opcode) {
            case 0: a = shiftRight(a, combo()); break;
            case 1: b ^= operand; break;
            case 2: b = combo() % 8; break;
            case 3:
                if (a != 0) {
                    ip = static_cast<std::size_t>(operand);
                    continue;
                }
                break;
            case 4: b ^= c; break;
            case 5: output.push_back(static_cast<int>(combo() % 8)); break;
            case 6: b = shiftRight(a, combo()); break;
            case 7: c = shiftRight(a, combo()); break;
        }
        ip += 2;
    }
    return output;
}

inline std::string join(const std::vector<int>& values) {
    std::string text;
    for (const int value : values) {
        if (!text.empty()) text += ',';
        text += static_cast<char>('0' + value);
    }
    return text;
}

// A self-replicating program consumes three bits of A per output, so build A
// three bits at a time from the last output backwards, keeping every candidate
// whose run reproduces the remaining tail of the program.
inline long long findSelfReplicatingA(const Machine& machine) {
    std::vector<long long> candidates{0};
    for (std::size_t index = machine.program.size(); index-- > 0;) {
        const std::span<const int> tail{machine.program.data() + index, machine.program.size() - index};
        std::vector<long long> next;
        for (const long long prefix : candidates) {
            for (int digit = 0; digit < 8; ++digit) {
                const long long candidate = prefix * 8 + digit;
                if (std::ranges::equal(run(machine, candidate), tail)) next.push_back(candidate);
            }
        }
        candidates = std::move(next);
    }
    const auto lowest = std::ranges::min_element(candidates);
    return lowest == candidates.end() ? 0 : *lowest;
}

}  // namespace detail

inline Answers solve(std::string_view input) {
    const detail::Machine machine = detail::parse(input);
    return Answers{detail::join(detail::run(machine, machine.a)),
                   detail::findSelfReplicatingA(machine)};
}

}  // namespace day17
