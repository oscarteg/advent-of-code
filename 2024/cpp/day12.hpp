#pragma once

#include <array>
#include <cctype>
#include <cstddef>
#include <sstream>
#include <stdexcept>
#include <string>
#include <string_view>
#include <vector>

namespace day12 {

struct Answers {
    long long part1;
    long long part2;
};

namespace detail {

struct Plot {
    int row;
    int col;
};

// What one region costs to fence: part 1 prices area * perimeter, part 2 area * sides.
struct Region {
    long long area = 0;
    long long perimeter = 0;
    long long sides = 0;
};

inline constexpr std::array<Plot, 4> kNeighbourSteps = {{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}};
inline constexpr std::array<int, 2> kDiagonalSteps = {-1, 1};

using Garden = std::vector<std::string>;

inline Garden parseGarden(std::string_view input) {
    Garden garden;
    std::istringstream stream{std::string{input}};
    std::string line;
    while (std::getline(stream, line)) {
        if (!line.empty() && line.back() == '\r') line.pop_back();
        if (line.empty()) continue;
        for (const char plant : line) {
            if (!std::isalpha(static_cast<unsigned char>(plant))) {
                throw std::invalid_argument("day12: garden plot is not a plant letter");
            }
        }
        garden.emplace_back(std::move(line));
    }
    if (garden.empty()) throw std::invalid_argument("day12: input has no garden rows");
    const std::size_t width = garden.front().size();
    for (const auto& row : garden) {
        if (row.size() != width) {
            throw std::invalid_argument("day12: garden rows have inconsistent width");
        }
    }
    return garden;
}

inline bool hasPlant(const Garden& garden, Plot at, char plant) {
    return at.row >= 0 && at.col >= 0 && at.row < static_cast<int>(garden.size()) &&
           at.col < static_cast<int>(garden.front().size()) && garden[at.row][at.col] == plant;
}

// A region has as many sides as corners: a corner is convex when both plots
// adjacent to a diagonal lie outside the region, and concave when both lie
// inside while the diagonal plot itself does not.
inline long long cornersAt(const Garden& garden, Plot at) {
    const char plant = garden[at.row][at.col];
    long long corners = 0;
    for (const int rowStep : kDiagonalSteps) {
        for (const int colStep : kDiagonalSteps) {
            const bool vertical = hasPlant(garden, {at.row + rowStep, at.col}, plant);
            const bool horizontal = hasPlant(garden, {at.row, at.col + colStep}, plant);
            const bool diagonal = hasPlant(garden, {at.row + rowStep, at.col + colStep}, plant);
            if (!vertical && !horizontal) ++corners;
            if (vertical && horizontal && !diagonal) ++corners;
        }
    }
    return corners;
}

// Flood-fills the region containing start, marking every plot it reaches visited.
inline Region measureRegion(const Garden& garden, Plot start, std::vector<bool>& visited) {
    const int width = static_cast<int>(garden.front().size());
    const char plant = garden[start.row][start.col];
    std::vector<Plot> queue{start};
    visited[start.row * width + start.col] = true;

    Region region;
    while (!queue.empty()) {
        const Plot current = queue.back();
        queue.pop_back();
        ++region.area;
        region.sides += cornersAt(garden, current);
        for (const Plot step : kNeighbourSteps) {
            const Plot next{current.row + step.row, current.col + step.col};
            if (!hasPlant(garden, next, plant)) {
                ++region.perimeter;
                continue;
            }
            if (visited[next.row * width + next.col]) continue;
            visited[next.row * width + next.col] = true;
            queue.push_back(next);
        }
    }
    return region;
}

}  // namespace detail

inline Answers solve(std::string_view input) {
    const auto garden = detail::parseGarden(input);
    const int height = static_cast<int>(garden.size());
    const int width = static_cast<int>(garden.front().size());
    std::vector<bool> visited(static_cast<std::size_t>(height) * static_cast<std::size_t>(width),
                              false);

    Answers answers{0, 0};
    for (int row = 0; row < height; ++row) {
        for (int col = 0; col < width; ++col) {
            if (visited[row * width + col]) continue;
            const auto region = detail::measureRegion(garden, {row, col}, visited);
            answers.part1 += region.area * region.perimeter;
            answers.part2 += region.area * region.sides;
        }
    }
    return answers;
}

}  // namespace day12
