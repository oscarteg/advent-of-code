const std = @import("std");

const Answers = struct {
    part1: i64,
    part2: i64,
};

fn parseStones(allocator: std.mem.Allocator, input: []const u8) ![]i64 {
    var stones: std.ArrayList(i64) = .empty;
    errdefer stones.deinit(allocator);

    var it = std.mem.tokenizeAny(u8, input, " \n\r\t");
    while (it.next()) |token| {
        const value = try std.fmt.parseInt(i64, token, 10);
        if (value < 0) return error.InvalidStone;
        try stones.append(allocator, value);
    }
    if (stones.items.len == 0) return error.InvalidStone;

    return stones.toOwnedSlice(allocator);
}

fn digitCount(value: i64) u32 {
    var remaining = value;
    var digits: u32 = 1;
    while (remaining >= 10) {
        remaining = @divTrunc(remaining, 10);
        digits += 1;
    }
    return digits;
}

fn addCount(map: *std.AutoHashMap(i64, i64), value: i64, count: i64) !void {
    const entry = try map.getOrPutValue(value, 0);
    entry.value_ptr.* += count;
}

fn blink(allocator: std.mem.Allocator, stones: *std.AutoHashMap(i64, i64)) !void {
    var next = std.AutoHashMap(i64, i64).init(allocator);
    errdefer next.deinit();

    var it = stones.iterator();
    while (it.next()) |entry| {
        const value = entry.key_ptr.*;
        const count = entry.value_ptr.*;
        if (value == 0) {
            try addCount(&next, 1, count);
            continue;
        }
        const digits = digitCount(value);
        if (digits % 2 == 0) {
            var divisor: i64 = 1;
            var i: u32 = 0;
            while (i < digits / 2) : (i += 1) divisor *= 10;
            try addCount(&next, @divTrunc(value, divisor), count);
            try addCount(&next, @mod(value, divisor), count);
        } else {
            try addCount(&next, value * 2024, count);
        }
    }

    stones.deinit();
    stones.* = next;
}

fn totalStones(stones: *const std.AutoHashMap(i64, i64)) i64 {
    var total: i64 = 0;
    var it = stones.valueIterator();
    while (it.next()) |count| total += count.*;
    return total;
}

fn solveWith(allocator: std.mem.Allocator, input: []const u8, blinks1: u32, blinks2: u32) !Answers {
    const parsed = try parseStones(allocator, input);
    defer allocator.free(parsed);

    var stones = std.AutoHashMap(i64, i64).init(allocator);
    defer stones.deinit();
    for (parsed) |value| try addCount(&stones, value, 1);

    var part1: i64 = 0;
    var part2: i64 = 0;
    const total_blinks = @max(blinks1, blinks2);
    var blinkNumber: u32 = 1;
    while (blinkNumber <= total_blinks) : (blinkNumber += 1) {
        try blink(allocator, &stones);
        if (blinkNumber == blinks1) part1 = totalStones(&stones);
        if (blinkNumber == blinks2) part2 = totalStones(&stones);
    }
    return Answers{ .part1 = part1, .part2 = part2 };
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    return solveWith(allocator, input, 25, 75);
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

test "official example after six blinks" {
    const input = @embedFile("examples/day11.txt");
    const answers = try solveWith(std.testing.allocator, input, 6, 6);
    try std.testing.expectEqual(@as(i64, 22), answers.part1);
}

test "official example" {
    const input = @embedFile("examples/day11.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(i64, 55312), answers.part1);
    try std.testing.expectEqual(@as(i64, 65601038650482), answers.part2);
}

test "single zero stone becomes one after one blink" {
    const answers = try solveWith(std.testing.allocator, "0\n", 1, 1);
    try std.testing.expectEqual(@as(i64, 1), answers.part1);
}

test "malformed input returns an error" {
    try std.testing.expectError(error.InvalidStone, solve(std.testing.allocator, ""));
    try std.testing.expectError(error.InvalidCharacter, solve(std.testing.allocator, "abc\n"));
    try std.testing.expectError(error.InvalidStone, solve(std.testing.allocator, "-5\n"));
}
