#pragma once

#include <array>
#include <stdexcept>
#include <string>
#include <string_view>
#include <vector>

namespace day25 {

struct Answers {
    long long part1;
    long long part2;
};

namespace detail {

constexpr int kRows = 7;
constexpr int kColumns = 5;
constexpr int kMaxHeight = kRows - 2;

using Heights = std::array<int, kColumns>;

inline std::vector<std::string_view> splitSchematics(std::string_view input) {
    std::vector<std::string_view> schematics;
    std::size_t start = 0;
    while (start < input.size()) {
        std::size_t end = input.find("\n\n", start);
        if (end == std::string_view::npos) end = input.size();
        schematics.push_back(input.substr(start, end - start));
        start = end + 2;
    }
    return schematics;
}

inline std::vector<std::string_view> splitLines(std::string_view schematic) {
    std::vector<std::string_view> lines;
    std::size_t start = 0;
    while (start <= schematic.size()) {
        std::size_t end = schematic.find('\n', start);
        if (end == std::string_view::npos) end = schematic.size();
        if (end > start) lines.push_back(schematic.substr(start, end - start));
        if (end == schematic.size()) break;
        start = end + 1;
    }
    return lines;
}

inline Heights columnHeights(const std::vector<std::string_view>& rows) {
    Heights heights{};
    for (const auto& row : rows) {
        if (row.size() != kColumns) throw std::invalid_argument("row must have 5 columns");
        for (int column = 0; column < kColumns; ++column) {
            if (row[column] == '#') ++heights[column];
        }
    }
    for (int& height : heights) --height;
    return heights;
}

struct Schematic {
    bool isLock;
    Heights heights;
};

inline Schematic parseSchematic(std::string_view text) {
    const auto rows = splitLines(text);
    if (rows.size() != kRows) throw std::invalid_argument("schematic must have 7 rows");
    const bool topFull = rows.front().find('.') == std::string_view::npos;
    const bool bottomFull = rows.back().find('.') == std::string_view::npos;
    if (topFull == bottomFull) throw std::invalid_argument("schematic is neither lock nor key");
    return Schematic{.isLock = topFull, .heights = columnHeights(rows)};
}

inline bool fits(const Heights& lock, const Heights& key) {
    for (int column = 0; column < kColumns; ++column) {
        if (lock[column] + key[column] > kMaxHeight) return false;
    }
    return true;
}

}  // namespace detail

inline Answers solve(std::string_view input) {
    std::vector<detail::Heights> locks;
    std::vector<detail::Heights> keys;
    for (const auto& block : detail::splitSchematics(input)) {
        if (block.empty()) continue;
        const auto schematic = detail::parseSchematic(block);
        (schematic.isLock ? locks : keys).push_back(schematic.heights);
    }
    if (locks.empty() || keys.empty()) throw std::invalid_argument("input needs locks and keys");

    long long fittingPairs = 0;
    for (const auto& lock : locks) {
        for (const auto& key : keys) {
            if (detail::fits(lock, key)) ++fittingPairs;
        }
    }
    return Answers{.part1 = fittingPairs, .part2 = 0};
}

}  // namespace day25
