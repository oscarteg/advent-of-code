#include <cstddef>
#include <cstdlib>
#include <fstream>
#include <iostream>
#include <string>
#include <vector>

struct Direction {
  int x, y;
};

const Direction directions[] = {
    {1, 0},
    {
        -1,
        0,
    },
    {0, 1},
    {0, -1},
    {1, 1},
    {1, -1},
    {-1, 1},
    {-1, -1},
};

const Direction directionsxMas[] = {
    {1, 1},   // DR
    {-1, -1}, // UL
    {-1, 1},  // UR
    {1, -1},  // UR
};

bool is_valid(int x, int y, int rows, int columns) {
  return x >= 0 && x < columns && y >= 0 && y < rows;
}

bool check_direction(const std::vector<std::string> &grid, int x, int y,
                     const Direction &dir) {
  // We're starting from the intersection point ('A'), so check both endpoints
  // First endpoint
  int x1 = x - dir.x;
  int y1 = y - dir.y;
  // Second endpoint
  int x2 = x + dir.x;
  int y2 = y + dir.y;

  if (!is_valid(x1, y1, grid.size(), grid[0].size()) ||
      !is_valid(x2, y2, grid.size(), grid[0].size())) {
    return false;
  }

  // Check for "MAS" pattern
  if (grid[y1][x1] == 'M' && grid[y2][x2] == 'S') {
    return true;
  }
  // Check for "SAM" pattern
  if (grid[y1][x1] == 'S' && grid[y2][x2] == 'M') {
    return true;
  }

  return false;
}

bool check_x_mas(const std::vector<std::string> &grid, int x, int y) {
  // Define the diagonal direction pairs that form an X
  const std::pair<Direction, Direction> pairs[] = {
      {{1, 1}, {1, -1}},  // First diagonal pair
      {{-1, 1}, {-1, -1}} // Second diagonal pair
  };

  for (const auto &[dir1, dir2] : pairs) {
    if (check_direction(grid, x, y, dir1) &&
        check_direction(grid, x, y, dir2)) {
      return true;
    }
  }
  return false;
}

int count_x_mas(const std::vector<std::string> &grid) {
  int count = 0;
  for (int y = 0; y < grid.size(); y++) {
    for (int x = 0; x < grid[0].size(); x++) {
      // Only check positions with 'A' as they must be the intersection
      if (grid[y][x] == 'A' && check_x_mas(grid, x, y)) {
        count++;
      }
    }
  }
  return count;
}

// Check for XMAS starting at position (x,y) in direction d
bool check_xmas(const std::vector<std::string> &grid, int x, int y, int d) {
  const std::string target = "XMAS";

  for (size_t i = 0; i < target.length(); i++) {
    int new_x = x + i * directions[d].x;
    int new_y = y + i * directions[d].y;

    if (!is_valid(new_x, new_y, grid.size(), grid[0].size()) ||
        grid[new_y][new_x] != target[i]) {
      return false;
    }
  }
  return true;
}

int count_xmas(const std::vector<std::string> &grid) {
  int count = 0;
  int rows = grid.size();
  int cols = grid[0].size();

  // Check each starting position
  for (int y = 0; y < rows; y++) {
    for (int x = 0; x < cols; x++) {
      // If we find an 'X', check all 8 directions
      if (grid[y][x] == 'X') {
        for (int d = 0; d < 8; d++) {
          if (check_xmas(grid, x, y, d)) {
            count++;
          }
        }
      }
    }
  }
  return count;
}

int main() {
  std::vector<std::string> exampleGrid = {
      "MMMSXXMASM", "MSAMXMSMSA", "AMXSXMAAMM", "MSAMASMSMX", "XMASAMXAMM",
      "XXAMMXXAMA", "SMSMSASXSS", "SAXAMASAAA", "MAMMMXMMMM", "MXMXAXMASX"};

  std::ifstream input("day4_input.txt");

  if (!input.is_open()) {
    std::cerr << "Could not open input file" << std::endl;
    return 1;
  }

  std::vector<std::string> grid;
  std::string line;

  while (std::getline(input, line)) {
    if (line.empty()) {
      break;
    }
    grid.push_back(line);
  }

  int total_part_1 = count_xmas(grid);

  int total_part_2 = count_x_mas(grid);

  std::cout << "Part 1: " << total_part_1 << std::endl;
  std::cout << "Part 2: " << total_part_2 << std::endl;
  return 0;
}
