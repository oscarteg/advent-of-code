#pragma once

#include <functional>
#include <string_view>
#include <vector>
#include <unordered_map>

namespace day19 {

struct Answers {
    long long part1;
    long long part2;
};

Answers solve(std::string_view input) {
    auto lines = std::vector<std::string_view>{};
    size_t pos = 0;
    
    while (pos < input.size()) {
        size_t next = input.find('\n', pos);
        if (next == std::string_view::npos) {
            if (pos < input.size()) {
                lines.emplace_back(input.substr(pos));
            }
            break;
        }
        lines.emplace_back(input.substr(pos, next - pos));
        pos = next + 1;
    }
    
    if (lines.empty()) {
        throw std::invalid_argument("empty input");
    }
    
    auto patterns = std::vector<std::string>{};
    std::string_view patterns_line = lines[0];
    pos = 0;
    while (pos < patterns_line.size()) {
        size_t comma = patterns_line.find(',', pos);
        std::string_view pattern = patterns_line.substr(pos, comma == std::string_view::npos ? std::string_view::npos : comma - pos);
        
        while (!pattern.empty() && (pattern.front() == ' ' || pattern.front() == '\t')) {
            pattern = pattern.substr(1);
        }
        while (!pattern.empty() && (pattern.back() == ' ' || pattern.back() == '\t')) {
            pattern = pattern.substr(0, pattern.size() - 1);
        }
        
        if (!pattern.empty()) {
            patterns.emplace_back(pattern);
        }
        
        if (comma == std::string_view::npos) break;
        pos = comma + 1;
    }
    
    auto designs = std::vector<std::string>{};
    for (size_t i = 1; i < lines.size(); ++i) {
        if (!lines[i].empty()) {
            designs.emplace_back(lines[i]);
        }
    }
    
    long long part1_count = 0;
    long long part2_sum = 0;
    
    for (const auto& design : designs) {
        auto memo = std::unordered_map<std::string, long long>{};
        
        std::function<long long(std::string_view)> count_ways = [&](std::string_view remaining) -> long long {
            if (remaining.empty()) {
                return 1;
            }
            
            std::string key(remaining);
            if (memo.count(key)) {
                return memo[key];
            }
            
            long long ways = 0;
            for (const auto& pattern : patterns) {
                if (remaining.size() >= pattern.size() && 
                    remaining.substr(0, pattern.size()) == pattern) {
                    ways += count_ways(remaining.substr(pattern.size()));
                }
            }
            
            return memo[key] = ways;
        };
        
        long long ways = count_ways(design);
        if (ways > 0) {
            part1_count++;
            part2_sum += ways;
        }
    }
    
    return {part1_count, part2_sum};
}

}
