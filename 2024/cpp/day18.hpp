#pragma once

#include <charconv>
#include <queue>
#include <stdexcept>
#include <string>
#include <string_view>
#include <vector>

namespace day18 {

struct Answers {
    long long part1;
    std::string part2;
};

struct Point {
    int x;
    int y;
};

namespace detail {

inline int parseNumber(std::string_view text) {
    int value = 0;
    const auto* begin = text.data();
    const auto* end = text.data() + text.size();
    const auto result = std::from_chars(begin, end, value);
    if (result.ec != std::errc{} || result.ptr != end) {
        throw std::invalid_argument("invalid number: " + std::string{text});
    }
    return value;
}

inline std::vector<Point> parse(std::string_view input) {
    std::vector<Point> points;
    std::size_t lineStart = 0;
    while (lineStart < input.size()) {
        std::size_t lineEnd = input.find('\n', lineStart);
        if (lineEnd == std::string_view::npos) lineEnd = input.size();
        std::string_view line = input.substr(lineStart, lineEnd - lineStart);
        while (!line.empty() && (line.back() == '\r' || line.back() == ' ')) line.remove_suffix(1);
        lineStart = lineEnd + 1;
        if (line.empty()) continue;

        const auto comma = line.find(',');
        if (comma == std::string_view::npos) {
            throw std::invalid_argument("missing ',' in line: " + std::string{line});
        }
        points.push_back(Point{.x = parseNumber(line.substr(0, comma)), .y = parseNumber(line.substr(comma + 1))});
    }
    if (points.empty()) {
        throw std::invalid_argument("input has no byte coordinates");
    }
    return points;
}

// Shortest path from (0,0) to (size-1,size-1) avoiding corrupted cells, or -1 if unreachable.
inline long long shortestPath(const std::vector<Point>& corrupted, int size) {
    std::vector<std::vector<bool>> blocked(static_cast<std::size_t>(size),
                                            std::vector<bool>(static_cast<std::size_t>(size), false));
    for (const auto& p : corrupted) {
        if (p.x >= 0 && p.x < size && p.y >= 0 && p.y < size) {
            blocked[static_cast<std::size_t>(p.y)][static_cast<std::size_t>(p.x)] = true;
        }
    }

    std::vector<std::vector<int>> distance(static_cast<std::size_t>(size),
                                            std::vector<int>(static_cast<std::size_t>(size), -1));
    std::queue<Point> frontier;
    distance[0][0] = 0;
    frontier.push(Point{0, 0});

    constexpr int dx[] = {1, -1, 0, 0};
    constexpr int dy[] = {0, 0, 1, -1};

    while (!frontier.empty()) {
        const Point current = frontier.front();
        frontier.pop();
        if (current.x == size - 1 && current.y == size - 1) {
            return distance[static_cast<std::size_t>(current.y)][static_cast<std::size_t>(current.x)];
        }
        for (int direction = 0; direction < 4; ++direction) {
            const int nx = current.x + dx[direction];
            const int ny = current.y + dy[direction];
            if (nx < 0 || nx >= size || ny < 0 || ny >= size) continue;
            if (blocked[static_cast<std::size_t>(ny)][static_cast<std::size_t>(nx)]) continue;
            if (distance[static_cast<std::size_t>(ny)][static_cast<std::size_t>(nx)] != -1) continue;
            distance[static_cast<std::size_t>(ny)][static_cast<std::size_t>(nx)] =
                distance[static_cast<std::size_t>(current.y)][static_cast<std::size_t>(current.x)] + 1;
            frontier.push(Point{nx, ny});
        }
    }
    return -1;
}

// Binary search over the byte count for the first one that blocks the exit.
inline Point firstBlockingByte(const std::vector<Point>& points, int size) {
    std::size_t low = 0;
    std::size_t high = points.size();
    while (low < high) {
        const std::size_t mid = low + (high - low) / 2;
        const std::vector<Point> prefix(points.begin(), points.begin() + static_cast<long>(mid) + 1);
        if (shortestPath(prefix, size) == -1) {
            high = mid;
        } else {
            low = mid + 1;
        }
    }
    if (low >= points.size()) {
        throw std::invalid_argument("no byte blocks the exit");
    }
    return points[low];
}

}  // namespace detail

inline Answers solve(std::string_view input, int size = 7, std::size_t fallenBytes = 12) {
    const auto points = detail::parse(input);
    if (fallenBytes > points.size()) {
        throw std::invalid_argument("fewer byte coordinates than requested fallen bytes");
    }

    const std::vector<Point> prefix(points.begin(), points.begin() + static_cast<long>(fallenBytes));
    const long long part1 = detail::shortestPath(prefix, size);
    if (part1 == -1) {
        throw std::invalid_argument("exit unreachable after the given number of fallen bytes");
    }

    const Point blocker = detail::firstBlockingByte(points, size);
    return Answers{.part1 = part1,
                    .part2 = std::to_string(blocker.x) + "," + std::to_string(blocker.y)};
}

}  // namespace day18
