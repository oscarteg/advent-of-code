#pragma once

#include <stdexcept>
#include <string>
#include <string_view>
#include <utility>
#include <vector>

namespace day6 {

struct Answers {
    long long part1;
    long long part2;
};

namespace detail {

enum class Direction { Up, Right, Down, Left };

inline Direction turnRight(Direction direction) {
    switch (direction) {
        case Direction::Up: return Direction::Right;
        case Direction::Right: return Direction::Down;
        case Direction::Down: return Direction::Left;
        case Direction::Left: return Direction::Up;
    }
    throw std::invalid_argument("unreachable direction");
}

inline std::pair<int, int> delta(Direction direction) {
    switch (direction) {
        case Direction::Up: return {-1, 0};
        case Direction::Right: return {0, 1};
        case Direction::Down: return {1, 0};
        case Direction::Left: return {0, -1};
    }
    throw std::invalid_argument("unreachable direction");
}

inline Direction directionFromChar(char c) {
    switch (c) {
        case '^': return Direction::Up;
        case '>': return Direction::Right;
        case 'v': return Direction::Down;
        case '<': return Direction::Left;
    }
    throw std::invalid_argument("invalid guard marker");
}

struct Grid {
    std::vector<std::string> rows;
    int rowCount;
    int colCount;
    int startRow;
    int startCol;
    Direction startDirection;

    bool inBounds(int r, int c) const { return r >= 0 && r < rowCount && c >= 0 && c < colCount; }
    bool isWall(int r, int c) const { return rows[r][c] == '#'; }
};

inline Grid parseGrid(std::string_view input) {
    std::vector<std::string> rows;
    size_t pos = 0;
    while (pos < input.size()) {
        size_t next = input.find('\n', pos);
        if (next == std::string_view::npos) next = input.size();
        if (next > pos) {
            rows.push_back(std::string(input.substr(pos, next - pos)));
        }
        pos = next + 1;
    }

    if (rows.empty()) {
        throw std::invalid_argument("empty input");
    }

    const size_t cols = rows[0].size();
    for (const auto& row : rows) {
        if (row.size() != cols) {
            throw std::invalid_argument("irregular grid");
        }
    }

    bool found = false;
    int startRow = 0, startCol = 0;
    Direction startDirection = Direction::Up;

    for (int r = 0; r < static_cast<int>(rows.size()); ++r) {
        for (int c = 0; c < static_cast<int>(cols); ++c) {
            const char ch = rows[r][c];
            if (ch == '.' || ch == '#') continue;
            const Direction direction = directionFromChar(ch);
            if (found) {
                throw std::invalid_argument("multiple guards");
            }
            found = true;
            startRow = r;
            startCol = c;
            startDirection = direction;
            rows[r][c] = '.';
        }
    }

    if (!found) {
        throw std::invalid_argument("missing guard");
    }

    const int rowCount = static_cast<int>(rows.size());
    const int colCount = static_cast<int>(cols);
    return Grid{std::move(rows), rowCount, colCount, startRow, startCol, startDirection};
}

// Walks the guard from the start until it leaves the grid or revisits a
// (position, direction) state. Calls visit(r, c) for every cell entered.
template <class Visit>
inline bool walk(const Grid& grid, int extraWallRow, int extraWallCol, Visit&& visit) {
    std::vector<std::vector<uint8_t>> seen(
        grid.rowCount, std::vector<uint8_t>(grid.colCount, 0));

    int r = grid.startRow, c = grid.startCol;
    Direction direction = grid.startDirection;

    while (true) {
        const uint8_t mask = 1u << static_cast<int>(direction);
        if (seen[r][c] & mask) return true;  // loop detected
        seen[r][c] |= mask;
        visit(r, c);

        const auto [dr, dc] = delta(direction);
        const int nr = r + dr, nc = c + dc;
        if (!grid.inBounds(nr, nc)) return false;  // left the grid

        const bool blocked = grid.isWall(nr, nc) || (nr == extraWallRow && nc == extraWallCol);
        if (blocked) {
            direction = turnRight(direction);
        } else {
            r = nr;
            c = nc;
        }
    }
}

}

inline Answers solve(std::string_view input) {
    const detail::Grid grid = detail::parseGrid(input);

    std::vector<std::vector<uint8_t>> visited(grid.rowCount,
                                               std::vector<uint8_t>(grid.colCount, 0));
    long long part1 = 0;
    detail::walk(grid, -1, -1, [&](int r, int c) {
        if (!visited[r][c]) {
            visited[r][c] = 1;
            ++part1;
        }
    });

    long long part2 = 0;
    for (int r = 0; r < grid.rowCount; ++r) {
        for (int c = 0; c < grid.colCount; ++c) {
            if (!visited[r][c]) continue;
            if (r == grid.startRow && c == grid.startCol) continue;
            if (detail::walk(grid, r, c, [](int, int) {})) {
                ++part2;
            }
        }
    }

    return {part1, part2};
}

}
