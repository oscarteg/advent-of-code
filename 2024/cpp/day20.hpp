#pragma once

#include <cstddef>
#include <cstdlib>
#include <queue>
#include <stdexcept>
#include <string>
#include <string_view>
#include <vector>

namespace day20 {

struct Answers {
    long long part1;
    long long part2;
};

namespace detail {

struct Point {
    int x;
    int y;
};

struct Track {
    int width = 0;
    int height = 0;
    std::vector<bool> wall;  // row-major, width * height cells
    Point start{};
    Point end{};

    [[nodiscard]] std::size_t index(Point p) const {
        return static_cast<std::size_t>(p.y) * static_cast<std::size_t>(width) +
               static_cast<std::size_t>(p.x);
    }

    [[nodiscard]] bool contains(Point p) const {
        return p.x >= 0 && p.x < width && p.y >= 0 && p.y < height;
    }
};

inline Track parse(std::string_view input) {
    Track track;
    bool foundStart = false;
    bool foundEnd = false;

    std::size_t lineStart = 0;
    while (lineStart < input.size()) {
        std::size_t lineEnd = input.find('\n', lineStart);
        if (lineEnd == std::string_view::npos) lineEnd = input.size();
        std::string_view line = input.substr(lineStart, lineEnd - lineStart);
        lineStart = lineEnd + 1;
        while (!line.empty() && (line.back() == '\r' || line.back() == ' ')) line.remove_suffix(1);
        if (line.empty()) continue;

        const int lineWidth = static_cast<int>(line.size());
        if (track.width == 0) {
            track.width = lineWidth;
        } else if (lineWidth != track.width) {
            throw std::invalid_argument("map rows have different widths");
        }

        for (int column = 0; column < lineWidth; ++column) {
            const char cell = line[static_cast<std::size_t>(column)];
            const Point here{column, track.height};
            switch (cell) {
                case '#':
                case '.':
                    break;
                case 'S':
                    if (foundStart) throw std::invalid_argument("map has more than one start");
                    track.start = here;
                    foundStart = true;
                    break;
                case 'E':
                    if (foundEnd) throw std::invalid_argument("map has more than one end");
                    track.end = here;
                    foundEnd = true;
                    break;
                default:
                    throw std::invalid_argument("invalid map character: " + std::string(1, cell));
            }
            track.wall.push_back(cell == '#');
        }
        ++track.height;
    }

    if (track.height == 0) throw std::invalid_argument("input has no map rows");
    if (!foundStart) throw std::invalid_argument("map has no start position");
    if (!foundEnd) throw std::invalid_argument("map has no end position");
    return track;
}

// Steps from origin to every track cell reachable from it, -1 where unreachable.
inline std::vector<long long> distancesFrom(const Track& track, Point origin) {
    std::vector<long long> distance(track.wall.size(), -1);
    std::queue<Point> frontier;
    distance[track.index(origin)] = 0;
    frontier.push(origin);

    constexpr Point deltas[] = {{1, 0}, {-1, 0}, {0, 1}, {0, -1}};
    while (!frontier.empty()) {
        const Point current = frontier.front();
        frontier.pop();
        const long long steps = distance[track.index(current)];
        for (const Point delta : deltas) {
            const Point next{current.x + delta.x, current.y + delta.y};
            if (!track.contains(next)) continue;
            const std::size_t at = track.index(next);
            if (track.wall[at] || distance[at] != -1) continue;
            distance[at] = steps + 1;
            frontier.push(next);
        }
    }
    return distance;
}

struct Race {
    Track track;
    std::vector<long long> fromStart;
    std::vector<long long> toEnd;
    long long uncheatedTime;
};

inline Race prepare(std::string_view input) {
    Track track = parse(input);
    std::vector<long long> fromStart = distancesFrom(track, track.start);
    const long long uncheatedTime = fromStart[track.index(track.end)];
    if (uncheatedTime < 0) throw std::invalid_argument("the end is unreachable from the start");
    std::vector<long long> toEnd = distancesFrom(track, track.end);
    return Race{.track = std::move(track),
                .fromStart = std::move(fromStart),
                .toEnd = std::move(toEnd),
                .uncheatedTime = uncheatedTime};
}

// A cheat ignores walls for at most maxCheatTime picoseconds, so landing on a track cell at
// Manhattan distance d costs exactly d picoseconds: the cheated race then takes
// fromStart[begin] + d + toEnd[landing]. Cheats are identified by those two cells.
inline long long countCheats(const Race& race, int maxCheatTime, long long minSaving) {
    long long cheats = 0;
    for (int y = 0; y < race.track.height; ++y) {
        for (int x = 0; x < race.track.width; ++x) {
            const long long toCheat = race.fromStart[race.track.index({x, y})];
            if (toCheat < 0) continue;
            for (int dy = -maxCheatTime; dy <= maxCheatTime; ++dy) {
                const int room = maxCheatTime - std::abs(dy);
                for (int dx = -room; dx <= room; ++dx) {
                    const Point landing{x + dx, y + dy};
                    if (!race.track.contains(landing)) continue;
                    const long long fromCheat = race.toEnd[race.track.index(landing)];
                    if (fromCheat < 0) continue;
                    const long long jump = std::abs(dx) + std::abs(dy);
                    if (race.uncheatedTime - (toCheat + jump + fromCheat) >= minSaving) ++cheats;
                }
            }
        }
    }
    return cheats;
}

}  // namespace detail

// Counts the cheats saving at least minSaving picoseconds: part 1 may ignore walls for
// 2 picoseconds, part 2 for 20.
inline Answers solve(std::string_view input, long long minSaving = 100) {
    const detail::Race race = detail::prepare(input);
    return Answers{.part1 = detail::countCheats(race, 2, minSaving),
                   .part2 = detail::countCheats(race, 20, minSaving)};
}

}  // namespace day20
