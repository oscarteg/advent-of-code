#pragma once

#include <algorithm>
#include <set>
#include <stdexcept>
#include <string>
#include <string_view>
#include <unordered_map>
#include <unordered_set>
#include <vector>

namespace day23 {

struct Answers {
    long long part1;
    std::string part2;
};

bool is_clique(const std::vector<std::string>& nodes,
               const std::unordered_map<std::string, std::unordered_set<std::string>>& graph) {
    for (size_t i = 0; i < nodes.size(); ++i) {
        for (size_t j = i + 1; j < nodes.size(); ++j) {
            if (!graph.at(nodes[i]).count(nodes[j])) {
                return false;
            }
        }
    }
    return true;
}

Answers solve(std::string_view input) {
    std::unordered_map<std::string, std::unordered_set<std::string>> graph;
    
    size_t pos = 0;
    while (pos < input.size()) {
        size_t newline = input.find('\n', pos);
        if (newline == std::string_view::npos) newline = input.size();
        
        std::string_view line = input.substr(pos, newline - pos);
        if (line.empty()) {
            pos = newline + 1;
            continue;
        }
        
        size_t dash = line.find('-');
        if (dash == std::string_view::npos) {
            throw std::invalid_argument("Invalid edge format");
        }
        
        std::string a(line.substr(0, dash));
        std::string b(line.substr(dash + 1));
        
        graph[a].insert(b);
        graph[b].insert(a);
        
        pos = newline + 1;
    }
    
    std::set<std::set<std::string>> triangles;
    for (const auto& [node, neighbors] : graph) {
        std::vector<std::string> neighbor_vec(neighbors.begin(), neighbors.end());
        for (size_t i = 0; i < neighbor_vec.size(); ++i) {
            for (size_t j = i + 1; j < neighbor_vec.size(); ++j) {
                const auto& n1 = neighbor_vec[i];
                const auto& n2 = neighbor_vec[j];
                if (graph.at(n1).count(n2)) {
                    std::set<std::string> triangle{node, n1, n2};
                    triangles.insert(triangle);
                }
            }
        }
    }
    
    long long part1 = 0;
    for (const auto& triangle : triangles) {
        if (std::any_of(triangle.begin(), triangle.end(),
                       [](const std::string& s) { return !s.empty() && s[0] == 't'; })) {
            ++part1;
        }
    }
    
    std::vector<std::string> nodes;
    for (const auto& [node, _] : graph) {
        nodes.push_back(node);
    }
    
    std::vector<std::string> max_clique;
    size_t n = nodes.size();
    
    for (size_t i = 0; i < n; ++i) {
        for (size_t j = i + 1; j < n; ++j) {
            for (size_t k = j + 1; k < n; ++k) {
                for (size_t l = k + 1; l < n; ++l) {
                    std::vector<std::string> cand{nodes[i], nodes[j], nodes[k], nodes[l]};
                    if (is_clique(cand, graph) && cand.size() > max_clique.size()) {
                        max_clique = cand;
                    }
                }
            }
        }
    }
    
    if (max_clique.empty()) {
        for (size_t i = 0; i < n; ++i) {
            for (size_t j = i + 1; j < n; ++j) {
                for (size_t k = j + 1; k < n; ++k) {
                    std::vector<std::string> cand{nodes[i], nodes[j], nodes[k]};
                    if (is_clique(cand, graph) && cand.size() > max_clique.size()) {
                        max_clique = cand;
                    }
                }
            }
        }
    }
    
    std::sort(max_clique.begin(), max_clique.end());
    std::string part2;
    for (size_t i = 0; i < max_clique.size(); ++i) {
        if (i > 0) part2 += ',';
        part2 += max_clique[i];
    }
    
    return {part1, part2};
}

}
