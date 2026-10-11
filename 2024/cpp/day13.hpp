#pragma once

#include <regex>
#include <stdexcept>
#include <string_view>
#include <vector>

namespace day13 {

struct Machine {
    long long ax, ay, bx, by, px, py;
};

struct Answers {
    long long part1;
    long long part2;
};

inline std::vector<Machine> parse(std::string_view input) {
    std::vector<Machine> machines;
    const std::regex machine_regex(
        R"(Button A: X\+(\d+), Y\+(\d+)\nButton B: X\+(\d+), Y\+(\d+)\nPrize: X=(\d+), Y=(\d+))",
        std::regex::optimize);

    std::string input_str(input);
    std::sregex_iterator it(input_str.begin(), input_str.end(), machine_regex);
    std::sregex_iterator end;

    for (; it != end; ++it) {
        const auto& match = *it;
        if (match.size() != 7) {
            throw std::invalid_argument("Invalid machine format");
        }
        machines.push_back({std::stoll(match[1]), std::stoll(match[2]),
                            std::stoll(match[3]), std::stoll(match[4]),
                            std::stoll(match[5]), std::stoll(match[6])});
    }

    if (machines.empty()) {
        throw std::invalid_argument("No valid machines found");
    }

    return machines;
}

inline long long solve_machine_part1(const Machine& m) {
    for (long long a = 0; a <= 100; ++a) {
        for (long long b = 0; b <= 100; ++b) {
            if (a * m.ax + b * m.bx == m.px && a * m.ay + b * m.by == m.py) {
                return a * 3 + b * 1;
            }
        }
    }
    return 0;
}

inline long long solve_machine_part2(const Machine& m) {
    // Cramer's rule for a*ax + b*bx = px, a*ay + b*by = py
    const long long det = m.ax * m.by - m.ay * m.bx;
    if (det == 0) {
        return 0;
    }

    const long long num_a = m.px * m.by - m.py * m.bx;
    const long long num_b = m.ax * m.py - m.ay * m.px;

    if (num_a % det != 0 || num_b % det != 0) {
        return 0;
    }

    const long long a = num_a / det;
    const long long b = num_b / det;

    if (a >= 0 && b >= 0 && a * m.ax + b * m.bx == m.px &&
        a * m.ay + b * m.by == m.py) {
        return a * 3 + b * 1;
    }

    return 0;
}

inline Answers solve(std::string_view input) {
    const auto machines = parse(input);

    long long part1 = 0;
    long long part2 = 0;

    for (auto m : machines) {
        part1 += solve_machine_part1(m);

        m.px += 10000000000000LL;
        m.py += 10000000000000LL;
        part2 += solve_machine_part2(m);
    }

    return {part1, part2};
}

}
