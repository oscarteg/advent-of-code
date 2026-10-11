#pragma once

#include <regex>
#include <stdexcept>
#include <string>
#include <string_view>

namespace day3 {

struct Answers {
    long long part1;
    long long part2;
};

inline Answers solve(std::string_view input) {
    if (input.find_first_not_of(" \t\r\n") == std::string_view::npos) {
        throw std::invalid_argument("empty input");
    }
    static const std::regex instruction(R"(mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\))");
    Answers answers{0, 0};
    bool enabled = true;
    const std::cregex_iterator end;
    for (std::cregex_iterator it(input.data(), input.data() + input.size(), instruction); it != end; ++it) {
        const std::string token = it->str();
        if (token == "do()") {
            enabled = true;
        } else if (token == "don't()") {
            enabled = false;
        } else {
            const long long product = std::stoll((*it)[1].str()) * std::stoll((*it)[2].str());
            answers.part1 += product;
            if (enabled) answers.part2 += product;
        }
    }
    return answers;
}

}
