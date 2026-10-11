#pragma once

#include <algorithm>
#include <map>
#include <stdexcept>
#include <string_view>
#include <vector>

namespace day5 {

struct Answers {
    long long part1;
    long long part2;
};

using Page = int;
using Rule = std::pair<Page, Page>;
using Update = std::vector<Page>;

struct Graph {
    std::vector<Rule> rules;

    bool isValidOrder(const Update& update) const {
        std::map<Page, size_t> positions;
        for (size_t i = 0; i < update.size(); ++i) {
            positions[update[i]] = i;
        }

        for (const auto& [from, to] : rules) {
            auto from_it = positions.find(from);
            auto to_it = positions.find(to);

            if (from_it != positions.end() && to_it != positions.end()) {
                if (from_it->second > to_it->second) {
                    return false;
                }
            }
        }
        return true;
    }

    Update sortUpdate(Update update) const {
        std::map<Page, std::vector<Page>> adj;
        std::map<Page, int> in_degree;

        for (Page page : update) {
            adj[page] = {};
            in_degree[page] = 0;
        }

        for (const auto& [from, to] : rules) {
            if (adj.count(from) && adj.count(to)) {
                adj[from].push_back(to);
            }
        }

        for (const auto& [node, neighbors] : adj) {
            for (Page neighbor : neighbors) {
                in_degree[neighbor]++;
            }
        }

        std::vector<Page> sorted_pages;
        std::vector<Page> queue;

        for (const auto& [page, degree] : in_degree) {
            if (degree == 0) {
                queue.push_back(page);
            }
        }

        while (!queue.empty()) {
            Page current = queue.back();
            queue.pop_back();
            sorted_pages.push_back(current);

            for (Page neighbor : adj[current]) {
                in_degree[neighbor]--;
                if (in_degree[neighbor] == 0) {
                    queue.push_back(neighbor);
                }
            }
        }

        return sorted_pages;
    }
};

long long getMiddle(const Update& update) {
    if (update.empty()) return 0;
    return update[update.size() / 2];
}

Answers solve(std::string_view input) {
    Graph graph;
    std::vector<Update> updates;
    bool parsing_rules = true;
    size_t pos = 0;

    while (pos < input.size()) {
        size_t newline = input.find('\n', pos);
        std::string_view line = newline == std::string_view::npos ? input.substr(pos)
                                                                    : input.substr(pos, newline - pos);
        pos = newline == std::string_view::npos ? input.size() : newline + 1;

        if (line.empty() || line == "\r") {
            parsing_rules = false;
            continue;
        }

        if (!line.empty() && line.back() == '\r') {
            line.remove_suffix(1);
        }

        if (parsing_rules) {
            size_t pipe_pos = line.find('|');
            if (pipe_pos == std::string_view::npos) {
                throw std::invalid_argument("Invalid rule format");
            }

            auto parse_int = [](std::string_view s) {
                int result = 0;
                for (char c : s) {
                    if (c < '0' || c > '9') {
                        throw std::invalid_argument("Invalid page number");
                    }
                    result = result * 10 + (c - '0');
                }
                return result;
            };

            try {
                Page from = parse_int(line.substr(0, pipe_pos));
                Page to = parse_int(line.substr(pipe_pos + 1));
                graph.rules.push_back({from, to});
            } catch (const std::exception&) {
                throw std::invalid_argument("Invalid page number in rule");
            }
        } else {
            Update update;
            size_t start = 0;

            while (start < line.size()) {
                size_t comma = line.find(',', start);
                std::string_view num_str = comma == std::string_view::npos ? line.substr(start)
                                                                             : line.substr(start, comma - start);
                start = comma == std::string_view::npos ? line.size() : comma + 1;

                auto parse_int = [](std::string_view s) {
                    int result = 0;
                    for (char c : s) {
                        if (c < '0' || c > '9') {
                            throw std::invalid_argument("Invalid page number");
                        }
                        result = result * 10 + (c - '0');
                    }
                    return result;
                };

                try {
                    update.push_back(parse_int(num_str));
                } catch (const std::exception&) {
                    throw std::invalid_argument("Invalid page number in update");
                }
            }

            if (!update.empty()) {
                updates.push_back(update);
            }
        }
    }

    long long part1 = 0;
    long long part2 = 0;

    for (const auto& update : updates) {
        if (graph.isValidOrder(update)) {
            part1 += getMiddle(update);
        } else {
            Update sorted = graph.sortUpdate(update);
            part2 += getMiddle(sorted);
        }
    }

    return Answers{part1, part2};
}

}
