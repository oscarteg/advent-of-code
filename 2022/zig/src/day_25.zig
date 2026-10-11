const std = @import("std");

fn snafuDigitValue(digit: u8) !i64 {
    return switch (digit) {
        '2' => 2,
        '1' => 1,
        '0' => 0,
        '-' => -1,
        '=' => -2,
        else => error.InvalidDigit,
    };
}

fn parseSnafu(line: []const u8) !i64 {
    if (line.len == 0) return error.EmptyNumber;
    var value: i64 = 0;
    for (line) |digit| {
        value = value * 5 + try snafuDigitValue(digit);
    }
    return value;
}

fn toSnafu(allocator: std.mem.Allocator, value: i64) ![]u8 {
    if (value == 0) return allocator.dupe(u8, "0");

    var digits: std.ArrayList(u8) = .empty;
    defer digits.deinit(allocator);

    var remaining = value;
    while (remaining != 0) {
        const remainder = @mod(remaining, 5);
        const digit: u8, const carry: i64 = switch (remainder) {
            0 => .{ '0', 0 },
            1 => .{ '1', 0 },
            2 => .{ '2', 0 },
            3 => .{ '=', 1 },
            4 => .{ '-', 1 },
            else => unreachable,
        };
        try digits.append(allocator, digit);
        remaining = @divFloor(remaining, 5) + carry;
    }

    std.mem.reverse(u8, digits.items);
    return digits.toOwnedSlice(allocator);
}

const Answer = struct {
    part1: []u8,
};

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answer {
    var sum: i64 = 0;
    var count: usize = 0;
    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw_line| {
        const line = std.mem.trim(u8, raw_line, " \r\t");
        if (line.len == 0) continue;
        sum += try parseSnafu(line);
        count += 1;
    }
    if (count == 0) return error.EmptyInput;

    return .{ .part1 = try toSnafu(allocator, sum) };
}

pub fn main(init: std.process.Init) !void {
    const arena = init.arena.allocator();
    var stdin_buffer: [4096]u8 = undefined;
    var stdin_reader = std.Io.File.stdin().reader(init.io, &stdin_buffer);
    const input = try stdin_reader.interface.allocRemaining(arena, .unlimited);
    const answer = try solve(arena, input);
    var stdout_buffer: [256]u8 = undefined;
    var stdout_writer = std.Io.File.stdout().writer(init.io, &stdout_buffer);
    const stdout = &stdout_writer.interface;
    try stdout.print("Part 1: {s}\n", .{answer.part1});
    try stdout.flush();
}

test "converts single SNAFU digits" {
    try std.testing.expectEqual(@as(i64, 1), try parseSnafu("1"));
    try std.testing.expectEqual(@as(i64, 2), try parseSnafu("2"));
    try std.testing.expectEqual(@as(i64, 3), try parseSnafu("1="));
    try std.testing.expectEqual(@as(i64, 4), try parseSnafu("1-"));
    try std.testing.expectEqual(@as(i64, 5), try parseSnafu("10"));
}

test "round trips values through SNAFU" {
    const values = [_]i64{ 1, 2, 3, 4, 5, 10, 2022, 4890 };
    for (values) |value| {
        const snafu = try toSnafu(std.testing.allocator, value);
        defer std.testing.allocator.free(snafu);
        try std.testing.expectEqual(value, try parseSnafu(snafu));
    }
}

test "sums the official example" {
    const example = @embedFile("examples/day_25.txt");
    const answer = try solve(std.testing.allocator, example);
    defer std.testing.allocator.free(answer.part1);
    try std.testing.expectEqualStrings("2=-1=0", answer.part1);
}

test "handles a single number" {
    const answer = try solve(std.testing.allocator, "1=-0-2");
    defer std.testing.allocator.free(answer.part1);
    try std.testing.expectEqualStrings("1=-0-2", answer.part1);
}

test "rejects empty input" {
    try std.testing.expectError(error.EmptyInput, solve(std.testing.allocator, ""));
    try std.testing.expectError(error.EmptyInput, solve(std.testing.allocator, "   \n  \n"));
}

test "rejects invalid digits" {
    try std.testing.expectError(error.InvalidDigit, solve(std.testing.allocator, "1=-0-9"));
    try std.testing.expectError(error.InvalidDigit, solve(std.testing.allocator, "1a2"));
}
