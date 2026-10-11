#pragma once

#include <charconv>
#include <stdexcept>
#include <string>
#include <string_view>
#include <unordered_map>
#include <vector>

namespace day11 {

struct Answers {
    long long part1;
    long long part2;
};

namespace detail {

inline long long parseNumber(std::string_view text) {
    long long value = 0;
    const auto* begin = text.data();
    const auto* end = text.data() + text.size();
    const auto result = std::from_chars(begin, end, value);
    if (result.ec != std::errc{} || result.ptr != end || value < 0) {
        throw std::invalid_argument("invalid stone value: " + std::string{text});
    }
    return value;
}

inline std::vector<long long> parse(std::string_view input) {
    std::vector<long long> stones;
    std::size_t pos = 0;
    while (pos < input.size()) {
        while (pos < input.size() && (input[pos] == ' ' || input[pos] == '\n' || input[pos] == '\r' ||
                                       input[pos] == '\t')) {
            ++pos;
        }
        if (pos >= input.size()) break;
        std::size_t next = pos;
        while (next < input.size() && input[next] != ' ' && input[next] != '\n' && input[next] != '\r' &&
               input[next] != '\t') {
            ++next;
        }
        stones.push_back(parseNumber(input.substr(pos, next - pos)));
        pos = next;
    }
    if (stones.empty()) {
        throw std::invalid_argument("input has no stones");
    }
    return stones;
}

inline int digitCount(long long value) {
    int digits = 1;
    while (value >= 10) {
        value /= 10;
        ++digits;
    }
    return digits;
}

// Applies the blink rule once to every distinct stone value, keeping a
// count per value instead of materializing the (exponentially growing)
// sequence of stones.
inline std::unordered_map<long long, long long> blink(const std::unordered_map<long long, long long>& stones) {
    std::unordered_map<long long, long long> next;
    next.reserve(stones.size() * 2);
    for (const auto& [value, count] : stones) {
        if (value == 0) {
            next[1] += count;
            continue;
        }
        const int digits = digitCount(value);
        if (digits % 2 == 0) {
            long long divisor = 1;
            for (int i = 0; i < digits / 2; ++i) divisor *= 10;
            next[value / divisor] += count;
            next[value % divisor] += count;
        } else {
            next[value * 2024] += count;
        }
    }
    return next;
}

inline long long totalStones(const std::unordered_map<long long, long long>& stones) {
    long long total = 0;
    for (const auto& [value, count] : stones) total += count;
    return total;
}

}  // namespace detail

// Runs up to max(blinks1, blinks2) blinks, returning the total stone count
// observed right after blinks1 and blinks2 blinks respectively.
inline Answers solveWith(std::string_view input, int blinks1, int blinks2) {
    const auto parsed = detail::parse(input);

    std::unordered_map<long long, long long> stones;
    for (long long value : parsed) stones[value] += 1;

    long long part1 = 0;
    long long part2 = 0;
    const int totalBlinks = blinks1 > blinks2 ? blinks1 : blinks2;
    for (int blink = 1; blink <= totalBlinks; ++blink) {
        stones = detail::blink(stones);
        if (blink == blinks1) part1 = detail::totalStones(stones);
        if (blink == blinks2) part2 = detail::totalStones(stones);
    }
    return {.part1 = part1, .part2 = part2};
}

inline Answers solve(std::string_view input) {
    return solveWith(input, 25, 75);
}

}  // namespace day11
