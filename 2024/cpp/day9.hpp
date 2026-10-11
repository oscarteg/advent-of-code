#pragma once

#include <stdexcept>
#include <string>
#include <string_view>
#include <vector>

namespace day9 {

struct Answers {
    long long part1;
    long long part2;
};

namespace detail {

constexpr long long kFree = -1;

struct Segment {
    long long id;  // kFree for free space
    long long start;
    long long length;
};

inline std::vector<Segment> parse(std::string_view input) {
    std::string digits;
    for (char c : input) {
        if (c == '\n' || c == '\r' || c == ' ') continue;
        if (c < '0' || c > '9') throw std::invalid_argument("invalid character in disk map");
        digits.push_back(c);
    }
    if (digits.empty()) throw std::invalid_argument("empty disk map");

    std::vector<Segment> segments;
    long long position = 0;
    long long fileId = 0;
    for (std::size_t i = 0; i < digits.size(); ++i) {
        const long long length = digits[i] - '0';
        const bool isFile = i % 2 == 0;
        if (length > 0) {
            segments.push_back({isFile ? fileId : kFree, position, length});
            position += length;
        }
        if (isFile) ++fileId;
    }
    return segments;
}

inline long long totalLength(const std::vector<Segment>& segments) {
    long long total = 0;
    for (const auto& s : segments) total += s.length;
    return total;
}

// Part 1 fragments files block by block, so it is simulated on an explicit block array.
inline long long solvePart1(const std::vector<Segment>& segments) {
    std::vector<long long> blocks(static_cast<std::size_t>(totalLength(segments)), kFree);
    for (const auto& s : segments) {
        for (long long i = 0; i < s.length; ++i) blocks[static_cast<std::size_t>(s.start + i)] = s.id;
    }

    std::size_t left = 0;
    std::size_t right = blocks.size();
    while (true) {
        while (left < blocks.size() && blocks[left] != kFree) ++left;
        while (right > 0 && blocks[right - 1] == kFree) --right;
        if (left >= right) break;
        blocks[left] = blocks[right - 1];
        blocks[right - 1] = kFree;
        ++left;
        --right;
    }

    long long checksum = 0;
    for (std::size_t i = 0; i < blocks.size(); ++i) {
        if (blocks[i] != kFree) checksum += static_cast<long long>(i) * blocks[i];
    }
    return checksum;
}

// Part 2 moves whole files, highest id first, into the leftmost free span that fits.
// Free spans stay sorted by start as they only shrink from the left, so a single
// left-to-right scan per file finds the leftmost fit.
inline long long solvePart2(const std::vector<Segment>& segments) {
    std::vector<Segment> files;
    std::vector<Segment> freeSpans;
    for (const auto& s : segments) {
        if (s.id == kFree) freeSpans.push_back(s); else files.push_back(s);
    }

    long long checksum = 0;
    for (auto fileIt = files.rbegin(); fileIt != files.rend(); ++fileIt) {
        long long position = fileIt->start;
        bool moved = false;
        for (auto& free : freeSpans) {
            if (free.start >= fileIt->start) break;
            if (free.length >= fileIt->length) {
                position = free.start;
                free.start += fileIt->length;
                free.length -= fileIt->length;
                moved = true;
                break;
            }
        }
        (void)moved;
        for (long long i = 0; i < fileIt->length; ++i) checksum += (position + i) * fileIt->id;
    }
    return checksum;
}

}  // namespace detail

inline Answers solve(std::string_view input) {
    const auto segments = detail::parse(input);
    return {detail::solvePart1(segments), detail::solvePart2(segments)};
}

}  // namespace day9
