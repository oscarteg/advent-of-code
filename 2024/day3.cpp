

#include <cstddef>
#include <cstdlib>
#include <fstream>
#include <iostream>
#include <regex>
#include <sstream>
#include <string>

const std::regex MUL_PATTERN("mul\\((\\d{1,3}),(\\d{1,3})\\)");
const std::regex DO_PATTERN("do\\(\\)");
const std::regex DONT_PATTERN("don't\\(\\)");

int part_1(const std::string &line) {
  std::sregex_iterator iter(line.begin(), line.end(), MUL_PATTERN);
  std::sregex_iterator end;
  int sum = 0;
  while (iter != end) {
    std::smatch match = *iter;
    int x = std::stoi(match[1]);
    int y = std::stoi(match[2]);
    int result = x * y;
    sum += result;
    iter++;
  }
  return sum;
}

int part_2(const std::string &line, bool &enabled) {

  int total = 0;

  size_t pos = 0;
  while (pos < line.length()) {
    std::smatch match;
    std::string remainder = line.substr(pos);

    if (std::regex_search(remainder, match, DO_PATTERN)) {
      if (match.position() == 0) {
        enabled = true;
        pos += match.length();
        continue;
      }
    }
    if (std::regex_search(remainder, match, DONT_PATTERN)) {
      if (match.position() == 0) {
        enabled = false;
        pos += match.length();
        continue;
      }
    }
    if (std::regex_search(remainder, match, MUL_PATTERN)) {
      if (match.position() == 0) {
        if (enabled) {
          int x = std::stoi(match[1]);
          int y = std::stoi(match[2]);
          total += x * y;
        }
        pos += match.length();
        continue;
      }
    }
    pos++;
  }
  return total;
}

int main() {

  std::ifstream input("day3_input.txt");

  if (!input.is_open()) {
    std::cerr << "Could not open input file" << std::endl;
    return 1;
  }

  std::string line;

  int total_part_1 = 0;
  int total_part_2 = 0;

  bool enabled = true;

  while (std::getline(input, line)) {

    if (line.empty()) {
      break;
    }

    int sum_part_1 = part_1(line);
    int sum_part_2 = part_2(line, enabled);

    std::cout << sum_part_1 << std::endl;
    std::cout << sum_part_2 << std::endl;

    total_part_1 += sum_part_1;
    total_part_2 += sum_part_2;
  }

  std::cout << "Part 1: " << total_part_1 << std::endl;
  std::cout << "Part 2: " << total_part_2 << std::endl;
  return 0;
}
