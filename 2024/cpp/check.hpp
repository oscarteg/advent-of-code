#pragma once
#include <exception>
#include <iostream>
#include <source_location>

namespace check {

inline int failures = 0;

template <class Actual, class Expected>
void equal(const Actual& actual, const Expected& expected,
           std::source_location where = std::source_location::current()) {
    if (actual == expected) return;
    ++failures;
    std::cerr << where.file_name() << ':' << where.line() << ": expected " << expected
              << ", got " << actual << '\n';
}

template <class Fn>
void throws(Fn&& fn, std::source_location where = std::source_location::current()) {
    try {
        fn();
    } catch (const std::exception&) {
        return;
    }
    ++failures;
    std::cerr << where.file_name() << ':' << where.line() << ": expected an exception\n";
}

inline int exit_code() { return failures == 0 ? 0 : 1; }

}
