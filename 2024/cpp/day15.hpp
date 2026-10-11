#pragma once

#include <algorithm>
#include <cstddef>
#include <initializer_list>
#include <stdexcept>
#include <string>
#include <string_view>
#include <utility>
#include <vector>

namespace day15 {

struct Answers {
    long long part1;
    long long part2;
};

namespace detail {

struct Vector2 {
    int row = 0;
    int col = 0;

    friend bool operator==(Vector2, Vector2) = default;
};

inline Vector2 operator+(Vector2 left, Vector2 right) {
    return {left.row + right.row, left.col + right.col};
}

inline Vector2 operator-(Vector2 left, Vector2 right) {
    return {left.row - right.row, left.col - right.col};
}

inline std::size_t asIndex(int value) {
    return static_cast<std::size_t>(value);
}

// A rectangular warehouse map enclosed by walls, plus the robot's position.
// That wall border is what makes every push stop instead of leaving the map.
struct Warehouse {
    std::vector<std::string> rows;
    Vector2 robot;

    int height() const { return static_cast<int>(rows.size()); }
    int width() const { return static_cast<int>(rows.front().size()); }

    char tile(Vector2 at) const { return rows[asIndex(at.row)][asIndex(at.col)]; }
    char& tile(Vector2 at) { return rows[asIndex(at.row)][asIndex(at.col)]; }

    bool isEnclosed() const {
        for (int col = 0; col < width(); ++col) {
            if (tile({0, col}) != '#' || tile({height() - 1, col}) != '#') return false;
        }
        for (int row = 0; row < height(); ++row) {
            if (tile({row, 0}) != '#' || tile({row, width() - 1}) != '#') return false;
        }
        return true;
    }

    long long gpsSum(char boxTile) const {
        long long sum = 0;
        for (int row = 0; row < height(); ++row) {
            for (int col = 0; col < width(); ++col) {
                if (tile({row, col}) == boxTile) sum += 100LL * row + col;
            }
        }
        return sum;
    }
};

struct Puzzle {
    Warehouse warehouse;
    std::vector<Vector2> moves;
};

inline Vector2 toDirection(char move) {
    switch (move) {
        case '^': return {-1, 0};
        case 'v': return {1, 0};
        case '<': return {0, -1};
        case '>': return {0, 1};
        default: throw std::invalid_argument(std::string{"invalid robot move: "} + move);
    }
}

inline std::string_view trimTrailing(std::string_view line) {
    while (!line.empty() && (line.back() == '\r' || line.back() == ' ')) line.remove_suffix(1);
    return line;
}

// The map comes first; the blank line after it starts the move sequence, whose
// own line breaks carry no meaning.
inline Puzzle parse(std::string_view input) {
    Puzzle puzzle;
    bool readingMap = true;
    int robots = 0;

    while (!input.empty()) {
        const std::size_t newline = input.find('\n');
        const bool lastLine = newline == std::string_view::npos;
        const std::string_view line = trimTrailing(lastLine ? input : input.substr(0, newline));
        input = lastLine ? std::string_view{} : input.substr(newline + 1);

        if (line.empty()) {
            readingMap = readingMap && puzzle.warehouse.rows.empty();
            continue;
        }
        if (!readingMap) {
            for (const char move : line) puzzle.moves.push_back(toDirection(move));
            continue;
        }
        if (!puzzle.warehouse.rows.empty() && line.size() != puzzle.warehouse.rows.front().size()) {
            throw std::invalid_argument("warehouse rows have different widths");
        }
        for (std::size_t col = 0; col < line.size(); ++col) {
            const char current = line[col];
            if (current == '@') {
                ++robots;
                puzzle.warehouse.robot = {static_cast<int>(puzzle.warehouse.rows.size()),
                                          static_cast<int>(col)};
            } else if (current != '#' && current != '.' && current != 'O') {
                throw std::invalid_argument(std::string{"invalid warehouse tile: "} + current);
            }
        }
        puzzle.warehouse.rows.emplace_back(line);
    }

    if (puzzle.warehouse.rows.empty()) throw std::invalid_argument("input has no warehouse map");
    if (robots != 1) throw std::invalid_argument("warehouse needs exactly one robot");
    if (puzzle.moves.empty()) throw std::invalid_argument("input has no robot moves");
    if (!puzzle.warehouse.isEnclosed()) {
        throw std::invalid_argument("warehouse is not enclosed by walls");
    }
    return puzzle;
}

// Part 1: the robot shoves the whole run of boxes ahead of it one tile along
// the move, and only moves at all when the tile past that run is free. Hopping
// the near box to the far end leaves the same map as shifting every box.
inline void push(Warehouse& warehouse, Vector2 direction) {
    const Vector2 target = warehouse.robot + direction;
    Vector2 free = target;
    while (warehouse.tile(free) == 'O') free = free + direction;
    if (warehouse.tile(free) == '#') return;

    warehouse.tile(free) = 'O';
    warehouse.tile(target) = '@';
    warehouse.tile(warehouse.robot) = '.';
    warehouse.robot = target;
}

// Part 2 doubles every tile's width: walls and floor repeat, a box becomes the
// pair "[]", and the robot keeps a single tile of its own.
inline Warehouse widen(const Warehouse& narrow) {
    Warehouse wide;
    wide.rows.reserve(narrow.rows.size());
    for (const std::string& row : narrow.rows) {
        std::string wideRow;
        wideRow.reserve(row.size() * 2);
        for (const char current : row) {
            switch (current) {
                case 'O': wideRow += "[]"; break;
                case '@': wideRow += "@."; break;
                default: wideRow.append(2, current); break;
            }
        }
        wide.rows.push_back(std::move(wideRow));
    }
    wide.robot = {narrow.robot.row, narrow.robot.col * 2};
    return wide;
}

// Sideways the halves stay paired, so the whole run shifts one tile.
inline void pushWideRow(Warehouse& warehouse, Vector2 direction) {
    const Vector2 target = warehouse.robot + direction;
    Vector2 free = target;
    while (warehouse.tile(free) == '[' || warehouse.tile(free) == ']') free = free + direction;
    if (warehouse.tile(free) == '#') return;

    for (Vector2 at = free; at != target; at = at - direction) {
        warehouse.tile(at) = warehouse.tile(at - direction);
    }
    warehouse.tile(target) = '@';
    warehouse.tile(warehouse.robot) = '.';
    warehouse.robot = target;
}

// Collects the left halves of every box a vertical push would move, nearest
// first; returns false when a wall blocks any of them. A box can be reached
// through either half and the list stays small, so a linear duplicate check
// beats tracking a grid of visited tiles.
inline bool collectPushedBoxes(const Warehouse& warehouse, Vector2 direction,
                               std::vector<Vector2>& boxes) {
    const auto enqueue = [&boxes](Vector2 left) {
        if (std::find(boxes.begin(), boxes.end(), left) == boxes.end()) boxes.push_back(left);
    };

    const Vector2 target = warehouse.robot + direction;
    switch (warehouse.tile(target)) {
        case '#': return false;
        case '[': enqueue(target); break;
        case ']': enqueue(target + Vector2{0, -1}); break;
        default: return true;
    }

    for (std::size_t index = 0; index < boxes.size(); ++index) {
        const Vector2 left = boxes[index];
        for (const Vector2 half : {left, left + Vector2{0, 1}}) {
            const Vector2 next = half + direction;
            switch (warehouse.tile(next)) {
                case '#': return false;
                case '[': enqueue(next); break;
                case ']': enqueue(next + Vector2{0, -1}); break;
                default: break;
            }
        }
    }
    return true;
}

inline void pushWideColumn(Warehouse& warehouse, Vector2 direction) {
    std::vector<Vector2> boxes;
    if (!collectPushedBoxes(warehouse, direction, boxes)) return;

    // Farthest box first, so each one lands on tiles the previous box left.
    for (auto left = boxes.rbegin(); left != boxes.rend(); ++left) {
        const Vector2 right = *left + Vector2{0, 1};
        warehouse.tile(*left) = '.';
        warehouse.tile(right) = '.';
        warehouse.tile(*left + direction) = '[';
        warehouse.tile(right + direction) = ']';
    }
    warehouse.tile(warehouse.robot) = '.';
    warehouse.robot = warehouse.robot + direction;
    warehouse.tile(warehouse.robot) = '@';
}

inline void pushWide(Warehouse& warehouse, Vector2 direction) {
    if (direction.row == 0) {
        pushWideRow(warehouse, direction);
    } else {
        pushWideColumn(warehouse, direction);
    }
}

}  // namespace detail

inline Answers solve(std::string_view input) {
    const detail::Puzzle puzzle = detail::parse(input);

    detail::Warehouse narrow = puzzle.warehouse;
    detail::Warehouse wide = detail::widen(puzzle.warehouse);
    for (const detail::Vector2 direction : puzzle.moves) {
        detail::push(narrow, direction);
        detail::pushWide(wide, direction);
    }
    return {.part1 = narrow.gpsSum('O'), .part2 = wide.gpsSum('[')};
}

}  // namespace day15
