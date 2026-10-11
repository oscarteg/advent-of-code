const std = @import("std");

const Range = struct {
    start: u32,
    end: u32,
};

fn parseRange(s: []const u8) !Range {
    var parts = std.mem.tokenizeScalar(u8, s, '-');
    const start_str = parts.next() orelse return error.InvalidRange;
    const end_str = parts.next() orelse return error.InvalidRange;

    const start = try std.fmt.parseInt(u32, start_str, 10);
    const end = try std.fmt.parseInt(u32, end_str, 10);

    return Range{ .start = start, .end = end };
}

fn contains(a: Range, b: Range) bool {
    return a.start <= b.start and a.end >= b.end;
}

fn overlaps(a: Range, b: Range) bool {
    return @max(a.start, b.start) <= @min(a.end, b.end);
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !struct { u32, u32 } {
    var part1: u32 = 0;
    var part2: u32 = 0;

    var lines = std.mem.tokenizeScalar(u8, std.mem.trim(u8, input, "\n"), '\n');
    while (lines.next()) |line| {
        var parts = std.mem.tokenizeScalar(u8, line, ',');
        const a_str = parts.next() orelse return error.InvalidLine;
        const b_str = parts.next() orelse return error.InvalidLine;

        const a = try parseRange(a_str);
        const b = try parseRange(b_str);

        if (contains(a, b) or contains(b, a)) {
            part1 += 1;
        }

        if (overlaps(a, b)) {
            part2 += 1;
        }
    }

    _ = allocator;
    return .{ part1, part2 };
}

pub fn main(init: std.process.Init) !void {
    const arena = init.arena.allocator();
    var stdin_buffer: [4096]u8 = undefined;
    var stdin_reader = std.Io.File.stdin().reader(init.io, &stdin_buffer);
    const input = try stdin_reader.interface.allocRemaining(arena, .unlimited);

    const answers = try solve(arena, input);

    var stdout_buffer: [256]u8 = undefined;
    var stdout_writer = std.Io.File.stdout().writer(init.io, &stdout_buffer);
    const stdout = &stdout_writer.interface;
    try stdout.print("Part 1: {d}\nPart 2: {d}\n", .{ answers[0], answers[1] });
    try stdout.flush();
}

test "parseRange" {
    const range = try parseRange("2-4");
    try std.testing.expectEqual(Range{ .start = 2, .end = 4 }, range);
}

test "contains" {
    const a = Range{ .start = 2, .end = 8 };
    const b = Range{ .start = 3, .end = 7 };
    try std.testing.expect(contains(a, b));
    try std.testing.expect(!contains(b, a));
}

test "overlaps" {
    const a = Range{ .start = 5, .end = 7 };
    const b = Range{ .start = 7, .end = 9 };
    try std.testing.expect(overlaps(a, b));
}

test "example" {
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day_4.txt");
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(answers[0], 2);
    try std.testing.expectEqual(answers[1], 4);
}

test "malformed input" {
    const allocator = std.testing.allocator;
    _ = allocator;
    const invalid = "2-4";
    try std.testing.expectError(error.InvalidLine, solve(std.testing.allocator, invalid));
}
