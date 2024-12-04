#include <algorithm>
#include <cstdio>
#include <iostream>
#include <sstream>
#include <string>
#include <utility>
#include <vector>

std::pair<std::vector<int>, std::vector<int>> parse_input() {
  // Left and right column
  std::vector<int> left, right;
  std::string line;

  // While next line
  while (std::getline(std::cin, line)) {
    if (line.empty()) {
      continue;
    }

    std::istringstream iss(line);

    int l, r;
    iss >> l >> r;
    left.push_back(l);
    right.push_back(r);
  }

  return {left, right};
}

int main(int argc, char *argv[]) {
  //
  //   std::istringstream input(R"(3   4
  // 4   3
  // 2   5
  // 1   3
  // 3   9
  // 3   3)");

  auto [left, right] = parse_input();

  std::sort(left.begin(), left.end());
  std::sort(right.begin(), right.end());

  int distance = 0;

  for (size_t i = 0; i < left.size(); i++) {
    distance += std::abs(left[i] - right[i]);
  }

  std::cout << distance << std::endl;

  return 0;
}
