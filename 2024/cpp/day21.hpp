#pragma once

#include <algorithm>
#include <array>
#include <charconv>
#include <cstddef>
#include <cstdlib>
#include <map>
#include <optional>
#include <span>
#include <stdexcept>
#include <string>
#include <string_view>
#include <tuple>
#include <vector>

namespace day21 {

struct Answers {
    long long part1;
    long long part2;
};

// Directional-keypad robots between the door's numeric keypad and the human.
inline constexpr int part1Robots = 2;
inline constexpr int part2Robots = 25;

namespace detail {

// Keypad rows, with ' ' marking the gap no robot arm may cross.
inline constexpr std::array<std::string_view, 4> numericRows{"789", "456", "123", " 0A"};
inline constexpr std::array<std::string_view, 2> directionalRows{" ^A", "<v>"};

enum class Pad { numeric, directional };

inline std::span<const std::string_view> rowsOf(Pad pad) {
    return pad == Pad::numeric ? std::span<const std::string_view>{numericRows}
                               : std::span<const std::string_view>{directionalRows};
}

struct Position {
    int row;
    int column;
};

inline Position locate(Pad pad, char key) {
    const auto rows = rowsOf(pad);
    for (std::size_t row = 0; row < rows.size(); ++row) {
        const std::size_t column = rows[row].find(key);
        if (column != std::string_view::npos) {
            return Position{static_cast<int>(row), static_cast<int>(column)};
        }
    }
    throw std::invalid_argument(std::string{"key not on the keypad: "} + key);
}

inline bool hasKey(Pad pad, Position at) {
    if (at.row < 0 || at.column < 0) return false;
    const auto rows = rowsOf(pad);
    const auto row = static_cast<std::size_t>(at.row);
    const auto column = static_cast<std::size_t>(at.column);
    return row < rows.size() && column < rows[row].size() && rows[row][column] != ' ';
}

// Shortest routes for one move, each ending with the A that presses the key. Only two orders
// can be shortest - every horizontal step first, or every vertical step first - and an order
// is impossible when its turning point lands on the keypad's gap.
inline std::vector<std::string> shortestRoutes(Pad pad, char from, char to) {
    const Position start = locate(pad, from);
    const Position target = locate(pad, to);
    const std::string horizontal(static_cast<std::size_t>(std::abs(target.column - start.column)),
                                 target.column > start.column ? '>' : '<');
    const std::string vertical(static_cast<std::size_t>(std::abs(target.row - start.row)),
                               target.row > start.row ? 'v' : '^');

    std::vector<std::string> routes;
    if (hasKey(pad, Position{start.row, target.column})) {
        routes.push_back(horizontal + vertical + 'A');
    }
    if (!horizontal.empty() && !vertical.empty() &&
        hasKey(pad, Position{target.row, start.column})) {
        routes.push_back(vertical + horizontal + 'A');
    }
    return routes;
}

// Human button presses, memoised per keypad, key pair and robot count.
class Presses {
public:
    [[nodiscard]] long long forSequence(Pad pad, std::string_view sequence, int robots) {
        long long presses = 0;
        char from = 'A';  // every arm starts on its own A button
        for (const char to : sequence) {
            presses += forMove(pad, from, to, robots);
            from = to;
        }
        return presses;
    }

private:
    [[nodiscard]] long long forMove(Pad pad, char from, char to, int robots) {
        const CacheKey key{pad, from, to, robots};
        if (const auto cached = cache_.find(key); cached != cache_.end()) return cached->second;

        std::optional<long long> cheapest;
        for (const std::string& route : shortestRoutes(pad, from, to)) {
            // The human types the last keypad itself, one press per button of the route.
            const long long presses = robots == 0
                                          ? static_cast<long long>(route.size())
                                          : forSequence(Pad::directional, route, robots - 1);
            cheapest = cheapest ? std::min(*cheapest, presses) : presses;
        }
        if (!cheapest) throw std::invalid_argument("the keypad gap blocks every route");

        cache_.emplace(key, *cheapest);
        return *cheapest;
    }

    using CacheKey = std::tuple<Pad, char, char, int>;
    std::map<CacheKey, long long> cache_;
};

struct Code {
    std::string_view buttons;
    long long numericPart;
};

inline Code parseCode(std::string_view line) {
    const auto reject = [line] {
        throw std::invalid_argument("door code must be digits followed by A: " + std::string{line});
    };
    if (line.empty() || line.back() != 'A') reject();

    const std::string_view digits = line.substr(0, line.size() - 1);
    const auto isDigit = [](char character) { return character >= '0' && character <= '9'; };
    long long numericPart = 0;
    if (!std::ranges::all_of(digits, isDigit) ||
        std::from_chars(digits.data(), digits.data() + digits.size(), numericPart).ec !=
            std::errc{}) {
        reject();
    }
    return Code{.buttons = line, .numericPart = numericPart};
}

inline std::vector<Code> parse(std::string_view input) {
    std::vector<Code> codes;
    std::size_t lineStart = 0;
    while (lineStart < input.size()) {
        std::size_t lineEnd = input.find('\n', lineStart);
        if (lineEnd == std::string_view::npos) lineEnd = input.size();
        std::string_view line = input.substr(lineStart, lineEnd - lineStart);
        lineStart = lineEnd + 1;
        while (!line.empty() && (line.back() == '\r' || line.back() == ' ')) line.remove_suffix(1);
        if (line.empty()) continue;
        codes.push_back(parseCode(line));
    }
    if (codes.empty()) throw std::invalid_argument("input has no door codes");
    return codes;
}

}  // namespace detail

// Sums the complexities of every code: the presses the human needs for the shortest sequence
// that types it, times the code's numeric part.
inline Answers solve(std::string_view input) {
    const std::vector<detail::Code> codes = detail::parse(input);
    detail::Presses presses;

    Answers answers{.part1 = 0, .part2 = 0};
    for (const detail::Code& code : codes) {
        const auto complexity = [&](int robots) {
            return presses.forSequence(detail::Pad::numeric, code.buttons, robots) *
                   code.numericPart;
        };
        answers.part1 += complexity(part1Robots);
        answers.part2 += complexity(part2Robots);
    }
    return answers;
}

}  // namespace day21
