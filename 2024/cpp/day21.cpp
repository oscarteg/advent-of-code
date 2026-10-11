#include "day21.hpp"

#include <exception>
#include <iostream>
#include <iterator>
#include <string>

int main() {
    const std::string input{std::istreambuf_iterator<char>{std::cin}, {}};
    try {
        const auto answers = day21::solve(input);
        std::cout << "Part 1: " << answers.part1 << "\nPart 2: " << answers.part2 << '\n';
    } catch (const std::exception& error) {
        std::cerr << error.what() << '\n';
        return 1;
    }
}
