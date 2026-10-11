#pragma once

#include <algorithm>
#include <charconv>
#include <sstream>
#include <stdexcept>
#include <string>
#include <string_view>
#include <unordered_map>
#include <vector>

namespace day1 {

struct Answers {
    long long part1;
    long long part2;
};

namespace detail {

inline std::pair<std::vector<int>, std::vector<int>> parse_lists(std::string_view input) {
    std::vector<int> left, right;
    size_t pos = 0;

    while (pos < input.size()) {
        size_t newline = input.find('\n', pos);
        std::string_view line = newline == std::string_view::npos ? input.substr(pos)
                                                                    : input.substr(pos, newline - pos);
        pos = newline == std::string_view::npos ? input.size() : newline + 1;

        if (line.empty() || line == "\r") continue;
        if (!line.empty() && line.back() == '\r') line.remove_suffix(1);

        std::istringstream iss{std::string{line}};
        std::string left_str, right_str;

        if (!(iss >> left_str >> right_str)) {
            throw std::invalid_argument("invalid line format");
        }

        int l = 0, r = 0;
        auto [l_ptr, l_ec] = std::from_chars(left_str.data(), left_str.data() + left_str.size(), l);
        if (l_ec != std::errc{}) {
            throw std::invalid_argument("invalid left number: " + left_str);
        }

        auto [r_ptr, r_ec] = std::from_chars(right_str.data(), right_str.data() + right_str.size(), r);
        if (r_ec != std::errc{}) {
            throw std::invalid_argument("invalid right number: " + right_str);
        }

        left.push_back(l);
        right.push_back(r);
    }

    if (left.empty()) {
        throw std::invalid_argument("empty input");
    }

    return {left, right};
}

}  // namespace detail

inline Answers solve(std::string_view input) {
    auto [left, right] = detail::parse_lists(input);

    std::sort(left.begin(), left.end());
    std::sort(right.begin(), right.end());

    long long distance = 0;
    for (size_t i = 0; i < left.size(); ++i) {
        distance += std::abs(left[i] - right[i]);
    }

    std::unordered_map<int, int> right_counts;
    for (int r : right) {
        right_counts[r]++;
    }

    long long score = 0;
    for (int l : left) {
        score += static_cast<long long>(l) * right_counts[l];
    }

    return {distance, score};
}

}  // namespace day1
