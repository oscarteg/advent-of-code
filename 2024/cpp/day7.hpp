#pragma once

#include <charconv>
#include <stdexcept>
#include <string>
#include <string_view>
#include <vector>

namespace day7 {

struct Answers {
    long long part1;
    long long part2;
};

struct Equation {
    long long target;
    std::vector<long long> operands;
};

namespace detail {

inline long long parseNumber(std::string_view text) {
    long long value = 0;
    const auto* begin = text.data();
    const auto* end = text.data() + text.size();
    const auto result = std::from_chars(begin, end, value);
    if (result.ec != std::errc{} || result.ptr != end) {
        throw std::invalid_argument("invalid number: " + std::string{text});
    }
    return value;
}

inline std::vector<Equation> parse(std::string_view input) {
    std::vector<Equation> equations;
    std::size_t lineStart = 0;
    while (lineStart < input.size()) {
        std::size_t lineEnd = input.find('\n', lineStart);
        if (lineEnd == std::string_view::npos) lineEnd = input.size();
        std::string_view line = input.substr(lineStart, lineEnd - lineStart);
        while (!line.empty() && (line.back() == '\r' || line.back() == ' ')) line.remove_suffix(1);
        lineStart = lineEnd + 1;
        if (line.empty()) continue;

        const auto colon = line.find(':');
        if (colon == std::string_view::npos) {
            throw std::invalid_argument("missing ':' in line: " + std::string{line});
        }
        Equation equation{.target = parseNumber(line.substr(0, colon)), .operands = {}};

        std::string_view rest = line.substr(colon + 1);
        std::size_t pos = 0;
        while (pos < rest.size()) {
            while (pos < rest.size() && rest[pos] == ' ') ++pos;
            if (pos >= rest.size()) break;
            std::size_t next = rest.find(' ', pos);
            if (next == std::string_view::npos) next = rest.size();
            equation.operands.push_back(parseNumber(rest.substr(pos, next - pos)));
            pos = next;
        }
        if (equation.operands.empty()) {
            throw std::invalid_argument("equation has no operands: " + std::string{line});
        }
        equations.push_back(std::move(equation));
    }
    if (equations.empty()) {
        throw std::invalid_argument("input has no equations");
    }
    return equations;
}

inline long long concatenate(long long left, long long right) {
    long long multiplier = 1;
    while (multiplier <= right) multiplier *= 10;
    return left * multiplier + right;
}

inline bool canReach(long long target, const std::vector<long long>& operands, std::size_t index,
                      long long accumulated, bool allowConcatenation) {
    if (accumulated > target) return false;
    if (index == operands.size()) return accumulated == target;

    const long long next = operands[index];
    if (canReach(target, operands, index + 1, accumulated + next, allowConcatenation)) return true;
    if (canReach(target, operands, index + 1, accumulated * next, allowConcatenation)) return true;
    if (allowConcatenation &&
        canReach(target, operands, index + 1, concatenate(accumulated, next), allowConcatenation)) {
        return true;
    }
    return false;
}

inline bool isSolvable(const Equation& equation, bool allowConcatenation) {
    return canReach(equation.target, equation.operands, 1, equation.operands.front(), allowConcatenation);
}

}  // namespace detail

inline Answers solve(std::string_view input) {
    const auto equations = detail::parse(input);

    long long part1 = 0;
    long long part2 = 0;
    for (const auto& equation : equations) {
        if (detail::isSolvable(equation, false)) {
            part1 += equation.target;
            part2 += equation.target;
        } else if (detail::isSolvable(equation, true)) {
            part2 += equation.target;
        }
    }
    return {.part1 = part1, .part2 = part2};
}

}  // namespace day7
