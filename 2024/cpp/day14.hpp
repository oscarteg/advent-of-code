#pragma once

#include <algorithm>
#include <regex>
#include <stdexcept>
#include <string>
#include <string_view>
#include <unordered_set>
#include <utility>
#include <vector>

namespace day14 {

struct Robot {
    long long px, py, vx, vy;
};

struct Answers {
    long long part1;
    long long part2;
};

using Position = std::pair<long long, long long>;

inline std::vector<Robot> parse(std::string_view input) {
    std::vector<Robot> robots;
    const std::regex robot_regex(R"(p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+))", std::regex::optimize);

    const std::string text(input);
    for (std::sregex_iterator it(text.begin(), text.end(), robot_regex), end; it != end; ++it) {
        const auto& match = *it;
        robots.push_back({std::stoll(match[1]), std::stoll(match[2]), std::stoll(match[3]),
                           std::stoll(match[4])});
    }

    if (robots.empty()) {
        throw std::invalid_argument("no valid robots found");
    }

    return robots;
}

inline long long wrap(long long value, long long modulus) {
    const long long remainder = value % modulus;
    return remainder < 0 ? remainder + modulus : remainder;
}

inline Position position_at(const Robot& robot, long long seconds, long long width,
                              long long height) {
    return {wrap(robot.px + robot.vx * seconds, width), wrap(robot.py + robot.vy * seconds, height)};
}

inline long long safety_factor(const std::vector<Robot>& robots, long long width, long long height,
                                long long seconds) {
    const long long mid_x = width / 2;
    const long long mid_y = height / 2;
    long long quadrants[4] = {0, 0, 0, 0};

    for (const auto& robot : robots) {
        const auto [x, y] = position_at(robot, seconds, width, height);
        if (x == mid_x || y == mid_y) continue;
        const int index = (x > mid_x ? 1 : 0) + (y > mid_y ? 2 : 0);
        ++quadrants[index];
    }

    return quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3];
}

inline bool positions_unique(const std::vector<Position>& positions, long long width) {
    std::unordered_set<long long> seen;
    seen.reserve(positions.size());
    for (const auto& [x, y] : positions) {
        if (!seen.insert(y * width + x).second) return false;
    }
    return true;
}

// Robots form a recognizable picture when no two of them overlap. Positions
// repeat with period width*height, so that bounds the search.
inline long long find_picture_second(const std::vector<Robot>& robots, long long width,
                                      long long height) {
    const long long period = width * height;
    std::vector<Position> positions;
    positions.reserve(robots.size());

    for (long long t = 0; t < period; ++t) {
        positions.clear();
        for (const auto& robot : robots) {
            positions.push_back(position_at(robot, t, width, height));
        }
        if (positions_unique(positions, width)) {
            return t;
        }
    }

    throw std::invalid_argument("no non-overlapping arrangement found within one period");
}

inline Answers compute(const std::vector<Robot>& robots, long long width, long long height) {
    return {safety_factor(robots, width, height, 100), find_picture_second(robots, width, height)};
}

inline Answers solve_with(std::string_view input, long long width, long long height) {
    return compute(parse(input), width, height);
}

// The input never states the grid size: the puzzle's example uses 11x7 and
// every real input uses 101x103, so pick between those two.
inline std::pair<long long, long long> infer_dimensions(const std::vector<Robot>& robots) {
    const bool fits_example = std::ranges::all_of(
        robots, [](const Robot& robot) { return robot.px < 11 && robot.py < 7; });
    return fits_example ? std::pair{11LL, 7LL} : std::pair{101LL, 103LL};
}

inline Answers solve(std::string_view input) {
    const auto robots = parse(input);
    const auto [width, height] = infer_dimensions(robots);
    return compute(robots, width, height);
}

}
