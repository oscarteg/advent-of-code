#pragma once

#include <array>
#include <charconv>
#include <cstdint>
#include <stdexcept>
#include <string>
#include <string_view>
#include <unordered_map>
#include <vector>

namespace day22 {

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
    if (result.ec != std::errc{} || result.ptr != end) {
        throw std::invalid_argument("invalid number: " + std::string{text});
    }
    return value;
}

inline std::vector<long long> parse(std::string_view input) {
    std::vector<long long> seeds;
    std::size_t lineStart = 0;
    while (lineStart < input.size()) {
        std::size_t lineEnd = input.find('\n', lineStart);
        if (lineEnd == std::string_view::npos) lineEnd = input.size();
        std::string_view line = input.substr(lineStart, lineEnd - lineStart);
        while (!line.empty() && (line.back() == '\r' || line.back() == ' ')) line.remove_suffix(1);
        lineStart = lineEnd + 1;
        if (line.empty()) continue;
        seeds.push_back(parseNumber(line));
    }
    if (seeds.empty()) {
        throw std::invalid_argument("input has no seeds");
    }
    return seeds;
}

constexpr long long kModulus = 16777216;

inline long long nextSecret(long long secret) {
    secret = ((secret * 64) ^ secret) % kModulus;
    secret = ((secret / 32) ^ secret) % kModulus;
    secret = ((secret * 2048) ^ secret) % kModulus;
    return secret;
}

// Encodes four price deltas, each in [-9, 9], into a single non-negative key.
inline int encodeSequence(const std::array<int, 4>& deltas) {
    int key = 0;
    for (const int delta : deltas) key = key * 19 + (delta + 9);
    return key;
}

}  // namespace detail

inline Answers solve(std::string_view input) {
    const auto seeds = detail::parse(input);

    long long part1 = 0;
    std::unordered_map<int, long long> totalBananasBySequence;

    for (const long long seed : seeds) {
        long long secret = seed;
        std::array<int, 4> deltas{};
        int previousPrice = static_cast<int>(secret % 10);
        std::unordered_map<int, bool> seenBySequence;

        for (int step = 1; step <= 2000; ++step) {
            secret = detail::nextSecret(secret);
            const int price = static_cast<int>(secret % 10);
            for (int i = 0; i < 3; ++i) deltas[i] = deltas[i + 1];
            deltas[3] = price - previousPrice;
            previousPrice = price;

            if (step >= 4) {
                const int key = detail::encodeSequence(deltas);
                if (seenBySequence.try_emplace(key, true).second) {
                    totalBananasBySequence[key] += price;
                }
            }
        }
        part1 += secret;
    }

    long long part2 = 0;
    for (const auto& [sequence, bananas] : totalBananasBySequence) {
        (void)sequence;
        if (bananas > part2) part2 = bananas;
    }

    return {.part1 = part1, .part2 = part2};
}

}  // namespace day22
