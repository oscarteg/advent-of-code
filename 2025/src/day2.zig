const std = @import("std");
const utils = @import("utils.zig");

const Answers = struct { part1: u64, part2: u64 };

const IdRange = struct { lo: u64, hi: u64 };

fn parseRange(token: []const u8) !IdRange {
    const trimmed = utils.trim(token);
    const dash = std.mem.indexOfScalar(u8, trimmed, '-') orelse return error.InvalidFormat;
    const lo = try utils.parseInt(u64, trimmed[0..dash]);
    const hi = try utils.parseInt(u64, trimmed[dash + 1 ..]);
    if (lo > hi) return error.InvalidFormat;
    return .{ .lo = lo, .hi = hi };
}

/// True if `digits` is exactly two copies of some shorter sequence.
fn isRepeatedExactlyTwice(digits: []const u8) bool {
    if (digits.len % 2 != 0) return false;
    const half = digits.len / 2;
    return std.mem.eql(u8, digits[0..half], digits[half..]);
}

/// True if `digits` has a repeating unit of length `period` covering the whole string.
fn hasPeriod(digits: []const u8, period: usize) bool {
    var i: usize = period;
    while (i < digits.len) : (i += 1) {
        if (digits[i] != digits[i % period]) return false;
    }
    return true;
}

/// True if `digits` is some shorter sequence repeated two or more times.
fn isRepeatedAtLeastTwice(digits: []const u8) bool {
    var period: usize = 1;
    while (period < digits.len) : (period += 1) {
        if (digits.len % period == 0 and hasPeriod(digits, period)) return true;
    }
    return false;
}

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    _ = allocator;
    var part1: u64 = 0;
    var part2: u64 = 0;

    var tokens = std.mem.tokenizeScalar(u8, utils.trim(input), ',');
    while (tokens.next()) |token| {
        const range = try parseRange(token);
        var n = range.lo;
        while (n <= range.hi) : (n += 1) {
            var buf: [20]u8 = undefined;
            const digits = std.fmt.bufPrint(&buf, "{d}", .{n}) catch unreachable;
            if (isRepeatedExactlyTwice(digits)) part1 += n;
            if (isRepeatedAtLeastTwice(digits)) part2 += n;
            if (n == range.hi) break;
        }
    }

    return .{ .part1 = part1, .part2 = part2 };
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
    try stdout.print("Part 1: {d}\nPart 2: {d}\n", .{ answers.part1, answers.part2 });
    try stdout.flush();
}

test "part 1 and part 2 on the official example" {
    const input = @embedFile("examples/day2.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u64, 1227775554), answers.part1);
    try std.testing.expectEqual(@as(u64, 4174379265), answers.part2);
}

test "part 2 counts sequences repeated more than twice" {
    const answers = try solve(std.testing.allocator, "123123123-123123123,1111111-1111111");
    try std.testing.expectEqual(@as(u64, 0), answers.part1);
    try std.testing.expectEqual(@as(u64, 123123123 + 1111111), answers.part2);
}

test "single-digit ids never count as repeated" {
    const answers = try solve(std.testing.allocator, "5-5");
    try std.testing.expectEqual(@as(u64, 0), answers.part1);
    try std.testing.expectEqual(@as(u64, 0), answers.part2);
}

test "malformed input returns an error" {
    try std.testing.expectError(error.InvalidFormat, solve(std.testing.allocator, "11-22,notarange"));
    try std.testing.expectError(error.InvalidFormat, solve(std.testing.allocator, "22-11"));
}
