#pragma once

#include <sstream>
#include <stdexcept>
#include <string>
#include <string_view>
#include <vector>

namespace day4 {

struct Answers {
    long long part1;
    long long part2;
};

namespace detail {

struct Direction {
    int dx;
    int dy;
};

inline constexpr Direction kAllDirections[] = {
    {1, 0},  {-1, 0}, {0, 1},  {0, -1},
    {1, 1},  {1, -1}, {-1, 1}, {-1, -1},
};

inline std::vector<std::string> parseGrid(std::string_view input) {
    std::vector<std::string> grid;
    std::istringstream stream{std::string{input}};
    std::string line;
    while (std::getline(stream, line)) {
        if (!line.empty() && line.back() == '\r') line.pop_back();
        if (line.empty()) continue;
        grid.emplace_back(std::move(line));
    }
    if (grid.empty()) throw std::invalid_argument("day4: input has no grid rows");
    const std::size_t width = grid.front().size();
    for (const auto& row : grid) {
        if (row.size() != width) throw std::invalid_argument("day4: grid rows have inconsistent width");
    }
    return grid;
}

inline bool inBounds(const std::vector<std::string>& grid, int row, int col) {
    return row >= 0 && col >= 0 && row < static_cast<int>(grid.size()) &&
           col < static_cast<int>(grid.front().size());
}

inline bool matchesWord(const std::vector<std::string>& grid, int row, int col, Direction dir,
                         std::string_view word) {
    for (std::size_t i = 0; i < word.size(); ++i) {
        const int r = row + dir.dy * static_cast<int>(i);
        const int c = col + dir.dx * static_cast<int>(i);
        if (!inBounds(grid, r, c) || grid[r][c] != word[i]) return false;
    }
    return true;
}

inline long long countXmas(const std::vector<std::string>& grid) {
    long long count = 0;
    for (int row = 0; row < static_cast<int>(grid.size()); ++row) {
        for (int col = 0; col < static_cast<int>(grid.front().size()); ++col) {
            if (grid[row][col] != 'X') continue;
            for (const auto& dir : kAllDirections) {
                if (matchesWord(grid, row, col, dir, "XMAS")) ++count;
            }
        }
    }
    return count;
}

// Checks the diagonal through (row, col) in the given direction pair forms "MAS" or "SAM".
inline bool diagonalIsMas(const std::vector<std::string>& grid, int row, int col, Direction dir) {
    const int r1 = row - dir.dy;
    const int c1 = col - dir.dx;
    const int r2 = row + dir.dy;
    const int c2 = col + dir.dx;
    if (!inBounds(grid, r1, c1) || !inBounds(grid, r2, c2)) return false;
    const char a = grid[r1][c1];
    const char b = grid[r2][c2];
    return (a == 'M' && b == 'S') || (a == 'S' && b == 'M');
}

inline long long countXMas(const std::vector<std::string>& grid) {
    long long count = 0;
    for (int row = 0; row < static_cast<int>(grid.size()); ++row) {
        for (int col = 0; col < static_cast<int>(grid.front().size()); ++col) {
            if (grid[row][col] != 'A') continue;
            if (diagonalIsMas(grid, row, col, {1, 1}) && diagonalIsMas(grid, row, col, {1, -1})) {
                ++count;
            }
        }
    }
    return count;
}

}  // namespace detail

inline Answers solve(std::string_view input) {
    const auto grid = detail::parseGrid(input);
    return {detail::countXmas(grid), detail::countXMas(grid)};
}

}  // namespace day4
