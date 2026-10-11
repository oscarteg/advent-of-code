#pragma once

#include <string_view>
#include <vector>
#include <map>
#include <set>
#include <stdexcept>
#include <numeric>

namespace day8 {

struct Answers {
    long long part1;
    long long part2;
};

Answers solve(std::string_view input) {
    std::vector<std::string> grid;
    std::string line;
    for (auto ch : input) {
        if (ch == '\n') {
            if (!line.empty()) {
                grid.push_back(line);
                line.clear();
            }
        } else {
            line += ch;
        }
    }
    if (!line.empty()) {
        grid.push_back(line);
    }

    if (grid.empty() || grid[0].empty()) {
        throw std::invalid_argument("Empty grid");
    }

    int rows = grid.size();
    int cols = grid[0].size();

    for (const auto& row : grid) {
        if (static_cast<int>(row.size()) != cols) {
            throw std::invalid_argument("Irregular grid dimensions");
        }
    }

    std::map<char, std::vector<std::pair<int, int>>> antennas_by_freq;
    for (int r = 0; r < rows; ++r) {
        for (int c = 0; c < cols; ++c) {
            char ch = grid[r][c];
            if (ch != '.') {
                antennas_by_freq[ch].push_back({r, c});
            }
        }
    }

    std::set<std::pair<int, int>> antinodes_part1;
    std::set<std::pair<int, int>> antinodes_part2;

    for (const auto& [freq, positions] : antennas_by_freq) {
        int n = positions.size();
        for (int i = 0; i < n; ++i) {
            for (int j = i + 1; j < n; ++j) {
                int r1 = positions[i].first, c1 = positions[i].second;
                int r2 = positions[j].first, c2 = positions[j].second;

                int dr = r2 - r1;
                int dc = c2 - c1;

                int antinode_r1 = r1 - dr;
                int antinode_c1 = c1 - dc;
                if (antinode_r1 >= 0 && antinode_r1 < rows && antinode_c1 >= 0 && antinode_c1 < cols) {
                    antinodes_part1.insert({antinode_r1, antinode_c1});
                }

                int antinode_r2 = r2 + dr;
                int antinode_c2 = c2 + dc;
                if (antinode_r2 >= 0 && antinode_r2 < rows && antinode_c2 >= 0 && antinode_c2 < cols) {
                    antinodes_part1.insert({antinode_r2, antinode_c2});
                }

                int gcd_val = std::gcd(std::abs(dr), std::abs(dc));
                if (gcd_val == 0) continue;
                
                int step_r = dr / gcd_val;
                int step_c = dc / gcd_val;

                for (int r = r1, c = c1; r >= 0 && r < rows && c >= 0 && c < cols; r += step_r, c += step_c) {
                    antinodes_part2.insert({r, c});
                }
                
                for (int r = r1 - step_r, c = c1 - step_c; r >= 0 && r < rows && c >= 0 && c < cols; r -= step_r, c -= step_c) {
                    antinodes_part2.insert({r, c});
                }
            }
        }
    }

    return {static_cast<long long>(antinodes_part1.size()), static_cast<long long>(antinodes_part2.size())};
}

}
