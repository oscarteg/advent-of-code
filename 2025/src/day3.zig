const std = @import("std");
const utils = @import("utils.zig");

const Answers = struct {
    part1: u64,
    part2: u64,
};

/// Pick k digits from a bank string greedily to form the largest k-digit number.
/// For each position, picks the largest digit such that enough digits remain.
fn pickKDigits(bank: []const u8, k: usize, allocator: std.mem.Allocator) !u64 {
    if (k > bank.len) return error.NotEnoughDigits;

    var result = std.ArrayList(u8).empty;
    defer result.deinit(allocator);

    var start_idx: usize = 0;

    for (0..k) |pos| {
        const remaining_needed = k - pos - 1;
        const end_idx = bank.len - remaining_needed - 1;

        var best_digit: ?u8 = null;
        var best_idx: usize = 0;

        for (start_idx..end_idx + 1) |i| {
            const digit = bank[i];
            if (best_digit == null or digit > best_digit.?) {
                best_digit = digit;
                best_idx = i;
            }
        }

        if (best_digit == null) return error.NoDigitFound;

        try result.append(allocator, best_digit.?);
        start_idx = best_idx + 1;
    }

    // Parse the result as a number
    var num: u64 = 0;
    for (result.items) |digit| {
        num = num * 10 + (digit - '0');
    }

    return num;
}

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var lines = utils.LineIterator.init(input);
    var part1_sum: u64 = 0;
    var part2_sum: u64 = 0;

    while (lines.next()) |line| {
        if (line.len == 0) continue;

        // Validate that all characters are digits
        for (line) |c| {
            if (c < '0' or c > '9') {
                return error.InvalidCharacter;
            }
        }

        const p1 = try pickKDigits(line, 2, allocator);
        const p2 = try pickKDigits(line, 12, allocator);

        part1_sum += p1;
        part2_sum += p2;
    }

    return .{
        .part1 = part1_sum,
        .part2 = part2_sum,
    };
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

test "example input" {
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day3.txt");
    const result = try solve(allocator, input);
    try std.testing.expectEqual(@as(u64, 357), result.part1);
    try std.testing.expectEqual(@as(u64, 3121910778619), result.part2);
}

test "single bank" {
    const allocator = std.testing.allocator;
    const input = "987654321111111\n";
    const result = try solve(allocator, input);
    try std.testing.expectEqual(@as(u64, 98), result.part1);
    try std.testing.expectEqual(@as(u64, 987654321111), result.part2);
}

test "invalid character" {
    const allocator = std.testing.allocator;
    const input = "12a34\n";
    const err = solve(allocator, input);
    try std.testing.expectError(error.InvalidCharacter, err);
}

test "not enough digits for part 2" {
    const allocator = std.testing.allocator;
    const input = "12345\n";
    const err = solve(allocator, input);
    try std.testing.expectError(error.NotEnoughDigits, err);
}
