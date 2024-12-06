#include <cstdio>
#include <cstdlib>
#include <iostream>
#include <sstream>
#include <string>
#include <vector>

bool is_safe(const std::vector<int> &levels) {
  if (levels.size() < 2)
    return true;
  bool should_increase = levels[1] > levels[0];
  for (size_t i = 1; i < levels.size(); i++) {
    int diff = levels[i] - levels[i - 1];
    if (abs(diff) < 1 || abs(diff) > 3)
      return false;
    if (should_increase && diff < 0)
      return false;
    if (!should_increase && diff > 0)
      return false;
  }
  return true;
}

std::vector<int> parse_line(const std::string &line) {
  std::vector<int> levels;
  std::stringstream iss(line);
  int level;
  while (iss >> level) {
    levels.push_back(level);
  }
  return levels;
}

bool is_safe_with_dampener(const std::vector<int> &levels) {
  if (is_safe(levels))
    return true;
  for (size_t i = 0; i < levels.size(); i++) {
    std::vector<int> modified = levels;
    modified.erase(modified.begin() + i);
    if (is_safe(modified))
      return true;
  }
  return false;
}

int main() {
  std::string line;
  int safe_count = 0;
  int safe_count_with_dampener = 0;
  std::vector<std::string> lines;

  while (std::getline(std::cin, line)) {
    if (line.empty())
      continue;
    lines.push_back(line);
  }

  for (const auto &line : lines) {
    std::vector<int> levels = parse_line(line);
    if (is_safe(levels))
      safe_count++;
    if (is_safe_with_dampener(levels)) {
      safe_count_with_dampener++;
    }
  }

  std::cout << "Part 1: " << safe_count << std::endl;
  std::cout << "Part 2: " << safe_count_with_dampener << std::endl;
  return 0;
}
