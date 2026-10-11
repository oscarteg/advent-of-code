#pragma once

#include <string_view>
#include <vector>
#include <set>
#include <stack>

namespace day10 {

struct Answers {
    long long part1;
    long long part2;
};

Answers solve(std::string_view input) {
    std::vector<std::string> grid;
    size_t pos = 0;
    while (pos < input.size()) {
        size_t next = input.find('\n', pos);
        if (next == std::string_view::npos) next = input.size();
        if (next > pos) {
            grid.push_back(std::string(input.substr(pos, next - pos)));
        }
        pos = next + 1;
    }

    if (grid.empty()) {
        throw std::invalid_argument("empty input");
    }

    int rows = grid.size();
    int cols = grid[0].size();

    for (const auto& row : grid) {
        if (row.size() != cols) {
            throw std::invalid_argument("irregular grid");
        }
        for (char c : row) {
            if (c < '0' || c > '9') {
                throw std::invalid_argument("invalid character in grid");
            }
        }
    }

    long long part1 = 0, part2 = 0;

    for (int r = 0; r < rows; ++r) {
        for (int c = 0; c < cols; ++c) {
            if (grid[r][c] == '0') {
                std::set<std::pair<int, int>> reachable_nines;
                long long path_count = 0;

                std::stack<std::pair<int, int>> st;
                st.push({r, c});

                while (!st.empty()) {
                    auto [cr, cc] = st.top();
                    st.pop();

                    if (grid[cr][cc] == '9') {
                        reachable_nines.insert({cr, cc});
                        path_count++;
                        continue;
                    }

                    char current = grid[cr][cc];
                    char next_elev = current + 1;

                    int dr[] = {-1, 1, 0, 0};
                    int dc[] = {0, 0, -1, 1};

                    for (int i = 0; i < 4; ++i) {
                        int nr = cr + dr[i];
                        int nc = cc + dc[i];
                        if (nr >= 0 && nr < rows && nc >= 0 && nc < cols &&
                            grid[nr][nc] == next_elev) {
                            st.push({nr, nc});
                        }
                    }
                }

                part1 += reachable_nines.size();
                part2 += path_count;
            }
        }
    }

    return {part1, part2};
}

}
