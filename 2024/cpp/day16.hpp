#pragma once

#include <algorithm>
#include <array>
#include <cstddef>
#include <functional>
#include <limits>
#include <queue>
#include <sstream>
#include <stdexcept>
#include <string>
#include <string_view>
#include <utility>
#include <vector>

namespace day16 {

struct Answers {
    long long part1;
    long long part2;
};

namespace detail {

struct Step {
    int row;
    int col;
};

// Facings in clockwise order, so a 90 degree turn is +-1 modulo four.
inline constexpr std::array<Step, 4> kFacings = {{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}};
inline constexpr int kFacingCount = static_cast<int>(kFacings.size());
inline constexpr int kEast = 1;
inline constexpr long long kStepScore = 1;
inline constexpr long long kTurnScore = 1000;
inline constexpr long long kUnreachable = std::numeric_limits<long long>::max();

// The maze as a row-major grid of tiles; start and end are tile indices.
struct Maze {
    std::string tiles;
    int width = 0;
    int height = 0;
    int start = 0;
    int end = 0;

    int tileCount() const { return static_cast<int>(tiles.size()); }

    bool isWall(int tile) const { return tiles[tile] == '#'; }

    // The tile one step along the facing, or against it when walking
    // backwards; -1 when that leaves the maze or hits a wall.
    int stepFrom(int tile, int facing, bool backwards) const {
        const Step step = kFacings[facing];
        const int row = tile / width + (backwards ? -step.row : step.row);
        const int col = tile % width + (backwards ? -step.col : step.col);
        if (row < 0 || col < 0 || row >= height || col >= width) return -1;
        const int next = row * width + col;
        return isWall(next) ? -1 : next;
    }
};

inline Maze parseMaze(std::string_view input) {
    Maze maze;
    int start = -1;
    int end = -1;
    std::istringstream stream{std::string{input}};
    std::string line;
    while (std::getline(stream, line)) {
        if (!line.empty() && line.back() == '\r') line.pop_back();
        if (line.empty()) continue;
        if (maze.height == 0) {
            maze.width = static_cast<int>(line.size());
        } else if (static_cast<int>(line.size()) != maze.width) {
            throw std::invalid_argument("day16: maze rows have inconsistent width");
        }
        for (int col = 0; col < maze.width; ++col) {
            const int tile = maze.height * maze.width + col;
            switch (line[col]) {
                case '#':
                case '.':
                    break;
                case 'S':
                    if (start >= 0) throw std::invalid_argument("day16: maze has two start tiles");
                    start = tile;
                    break;
                case 'E':
                    if (end >= 0) throw std::invalid_argument("day16: maze has two end tiles");
                    end = tile;
                    break;
                default:
                    throw std::invalid_argument("day16: maze tile is not #, ., S or E");
            }
        }
        maze.tiles += line;
        ++maze.height;
    }
    if (maze.height == 0) throw std::invalid_argument("day16: input has no maze rows");
    if (start < 0) throw std::invalid_argument("day16: maze has no start tile");
    if (end < 0) throw std::invalid_argument("day16: maze has no end tile");
    maze.start = start;
    maze.end = end;
    return maze;
}

// Dijkstra over (tile, facing) states, where a turn costs 1000 and a step
// forward 1. Walking backwards explores the reverse graph, which scores every
// state by its distance to the nearest source instead of from it.
inline std::vector<long long> scoresFrom(const Maze& maze, const std::vector<int>& sources,
                                         bool backwards) {
    std::vector<long long> best(
        static_cast<std::size_t>(maze.tileCount()) * static_cast<std::size_t>(kFacingCount),
        kUnreachable);
    using Move = std::pair<long long, int>;  // score first, so the cheapest state pops
    std::priority_queue<Move, std::vector<Move>, std::greater<Move>> frontier;
    for (const int source : sources) {
        best[source] = 0;
        frontier.emplace(0, source);
    }

    while (!frontier.empty()) {
        const Move current = frontier.top();
        frontier.pop();
        const long long score = current.first;
        if (score > best[current.second]) continue;

        const int tile = current.second / kFacingCount;
        const int facing = current.second % kFacingCount;
        std::array<Move, 3> moves{
            Move{score + kTurnScore, tile * kFacingCount + (facing + 1) % kFacingCount},
            Move{score + kTurnScore, tile * kFacingCount + (facing + 3) % kFacingCount},
            Move{},
        };
        std::size_t moveCount = 2;
        const int next = maze.stepFrom(tile, facing, backwards);
        if (next >= 0) moves[moveCount++] = Move{score + kStepScore, next * kFacingCount + facing};

        for (std::size_t move = 0; move < moveCount; ++move) {
            if (moves[move].first >= best[moves[move].second]) continue;
            best[moves[move].second] = moves[move].first;
            frontier.push(moves[move]);
        }
    }
    return best;
}

}  // namespace detail

inline Answers solve(std::string_view input) {
    const auto maze = detail::parseMaze(input);
    const std::vector<int> startStates{maze.start * detail::kFacingCount + detail::kEast};
    const auto fromStart = detail::scoresFrom(maze, startStates, false);

    std::vector<int> endStates;
    endStates.reserve(detail::kFacingCount);
    for (int facing = 0; facing < detail::kFacingCount; ++facing) {
        endStates.push_back(maze.end * detail::kFacingCount + facing);
    }
    const auto toEnd = detail::scoresFrom(maze, endStates, true);

    long long best = detail::kUnreachable;
    for (const int state : endStates) best = std::min(best, fromStart[state]);
    if (best == detail::kUnreachable) {
        throw std::invalid_argument("day16: no path from the start to the end tile");
    }

    // A state lies on a best path when reaching it and leaving it towards the
    // end together cost exactly the best score.
    long long tilesOnBestPaths = 0;
    for (int tile = 0; tile < maze.tileCount(); ++tile) {
        for (int facing = 0; facing < detail::kFacingCount; ++facing) {
            const int state = tile * detail::kFacingCount + facing;
            if (fromStart[state] == detail::kUnreachable) continue;
            if (toEnd[state] == detail::kUnreachable) continue;
            if (fromStart[state] + toEnd[state] != best) continue;
            ++tilesOnBestPaths;
            break;
        }
    }
    return Answers{best, tilesOnBestPaths};
}

}  // namespace day16
