const std = @import("std");

/// Reads all input from stdin into a buffer
pub fn readStdin(allocator: std.mem.Allocator) ![]const u8 {
    const input = try std.fs.File.stdin().readToEndAlloc(allocator, 10 * 1024 * 1024); // 10MB max
    return std.mem.trim(u8, input, &std.ascii.whitespace);
}

/// Iterator for lines in a string
pub const LineIterator = struct {
    buffer: []const u8,
    index: usize = 0,

    pub fn init(buffer: []const u8) LineIterator {
        return .{ .buffer = buffer };
    }

    pub fn next(self: *LineIterator) ?[]const u8 {
        if (self.index >= self.buffer.len) return null;

        const start = self.index;
        var end = start;

        while (end < self.buffer.len and self.buffer[end] != '\n') {
            end += 1;
        }

        // Trim \r if present (Windows line endings)
        var line_end = end;
        if (line_end > start and self.buffer[line_end - 1] == '\r') {
            line_end -= 1;
        }

        self.index = if (end < self.buffer.len) end + 1 else end;
        return self.buffer[start..line_end];
    }
};

/// Parse an integer from a string
pub fn parseInt(comptime T: type, s: []const u8) !T {
    return try std.fmt.parseInt(T, s, 10);
}

/// Parse an integer from a string, returning null on failure
pub fn parseIntOpt(comptime T: type, s: []const u8) ?T {
    return std.fmt.parseInt(T, s, 10) catch null;
}

/// Split a string by a delimiter
pub fn split(buffer: []const u8, delimiter: u8) std.mem.SplitIterator(u8, .scalar) {
    return std.mem.splitScalar(u8, buffer, delimiter);
}

/// Trim whitespace from both ends of a string
pub fn trim(s: []const u8) []const u8 {
    return std.mem.trim(u8, s, &std.ascii.whitespace);
}

pub fn divmod(a: i32, b: i32) struct {
    div: i32,
    rem: i32,
} {
    return .{ .div = @divTrunc(a, b), .rem = modPositive(a, b) };
}

/// Modulo operation that always returns positive result
pub fn modPositive(a: i32, n: i32) i32 {
    const r = @mod(a, n);
    return if (r < 0) r + n else r;
}

/// Common 2D grid directions (up, right, down, left)
pub const Direction = enum {
    up,
    right,
    down,
    left,

    pub fn delta(self: Direction) struct { dx: i32, dy: i32 } {
        return switch (self) {
            .up => .{ .dx = 0, .dy = -1 },
            .right => .{ .dx = 1, .dy = 0 },
            .down => .{ .dx = 0, .dy = 1 },
            .left => .{ .dx = -1, .dy = 0 },
        };
    }
};

/// 2D Point/Coordinate
pub const Point = struct {
    x: i32,
    y: i32,

    pub fn add(self: Point, other: Point) Point {
        return .{ .x = self.x + other.x, .y = self.y + other.y };
    }

    pub fn equals(self: Point, other: Point) bool {
        return self.x == other.x and self.y == other.y;
    }
};

/// GCD (Greatest Common Divisor)
pub fn gcd(a: i64, b: i64) i64 {
    var x = @abs(a);
    var y = @abs(b);
    while (y != 0) {
        const temp = y;
        y = @mod(x, y);
        x = temp;
    }
    return x;
}

/// LCM (Least Common Multiple)
pub fn lcm(a: i64, b: i64) i64 {
    if (a == 0 or b == 0) return 0;
    return @divExact(@abs(a * b), gcd(a, b));
}

/// Min of two values
pub fn min(comptime T: type, a: T, b: T) T {
    return if (a < b) a else b;
}

/// Max of two values
pub fn max(comptime T: type, a: T, b: T) T {
    return if (a > b) a else b;
}

/// Count occurrences of a value in a slice
pub fn count(comptime T: type, slice: []const T, value: T) usize {
    var cnt: usize = 0;
    for (slice) |item| {
        if (item == value) cnt += 1;
    }
    return cnt;
}

/// Sum all values in a slice
pub fn sum(comptime T: type, slice: []const T) T {
    var total: T = 0;
    for (slice) |value| {
        total += value;
    }
    return total;
}
