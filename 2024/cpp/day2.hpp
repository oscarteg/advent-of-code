#pragma once

#include <charconv>
#include <cmath>
#include <sstream>
#include <stdexcept>
#include <string>
#include <string_view>
#include <vector>

namespace day2 {

struct Answers {
    long long part1;
    long long part2;
};

namespace detail {

inline std::vector<int> parse_report(std::string_view line) {
    std::vector<int> levels;
    std::istringstream iss{std::string{line}};
    std::string token;
    while (iss >> token) {
        int value = 0;
        auto [ptr, ec] = std::from_chars(token.data(), token.data() + token.size(), value);
        if (ec != std::errc{} || ptr != token.data() + token.size()) {
            throw std::invalid_argument("invalid level in report: " + token);
        }
        levels.push_back(value);
    }
    if (levels.empty()) {
        throw std::invalid_argument("empty report");
    }
    return levels;
}

inline bool is_safe(const std::vector<int>& levels) {
    if (levels.size() < 2) return true;
    const bool increasing = levels[1] > levels[0];
    for (size_t i = 1; i < levels.size(); ++i) {
        const int diff = levels[i] - levels[i - 1];
        const int magnitude = std::abs(diff);
        if (magnitude < 1 || magnitude > 3) return false;
        if (increasing && diff < 0) return false;
        if (!increasing && diff > 0) return false;
    }
    return true;
}

inline bool is_safe_with_dampener(const std::vector<int>& levels) {
    if (is_safe(levels)) return true;
    for (size_t i = 0; i < levels.size(); ++i) {
        std::vector<int> without = levels;
        without.erase(without.begin() + static_cast<long>(i));
        if (is_safe(without)) return true;
    }
    return false;
}

}  // namespace detail

inline Answers solve(std::string_view input) {
    long long safe_count = 0;
    long long safe_with_dampener_count = 0;

    size_t pos = 0;
    while (pos < input.size()) {
        size_t newline = input.find('\n', pos);
        std::string_view line = newline == std::string_view::npos ? input.substr(pos)
                                                                    : input.substr(pos, newline - pos);
        pos = newline == std::string_view::npos ? input.size() : newline + 1;
        if (line.empty() || line == "\r") continue;
        if (!line.empty() && line.back() == '\r') line.remove_suffix(1);

        const auto levels = detail::parse_report(line);
        if (detail::is_safe(levels)) ++safe_count;
        if (detail::is_safe_with_dampener(levels)) ++safe_with_dampener_count;
    }

    return {safe_count, safe_with_dampener_count};
}

}  // namespace day2
