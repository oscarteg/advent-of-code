const std = @import("std");

const decryption_key: i64 = 811589153;
const mix_rounds_part2: usize = 10;

const Answers = struct {
    part1: i64,
    part2: i64,
};

fn parse(allocator: std.mem.Allocator, input: []const u8) ![]i64 {
    var numbers: std.ArrayList(i64) = .empty;
    errdefer numbers.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;
        const value = try std.fmt.parseInt(i64, trimmed, 10);
        try numbers.append(allocator, value);
    }

    if (numbers.items.len == 0) return error.EmptyInput;
    return numbers.toOwnedSlice(allocator);
}

fn mix(allocator: std.mem.Allocator, numbers: []const i64, rounds: usize) ![]i64 {
    const len = numbers.len;
    var order = try allocator.alloc(usize, len);
    defer allocator.free(order);
    for (order, 0..) |*slot, i| slot.* = i;

    if (len > 1) {
        var round: usize = 0;
        while (round < rounds) : (round += 1) {
            for (numbers, 0..) |value, original_index| {
                const current_pos = std.mem.indexOfScalar(usize, order, original_index).?;
                var i = current_pos;
                while (i + 1 < len) : (i += 1) order[i] = order[i + 1];

                const modulus: i64 = @intCast(len - 1);
                const new_pos: usize = @intCast(@mod(@as(i64, @intCast(current_pos)) + value, modulus));

                i = len - 1;
                while (i > new_pos) : (i -= 1) order[i] = order[i - 1];
                order[new_pos] = original_index;
            }
        }
    }

    const mixed = try allocator.alloc(i64, len);
    for (order, 0..) |original_index, i| mixed[i] = numbers[original_index];
    return mixed;
}

fn groveCoordinatesSum(mixed: []const i64) !i64 {
    const zero_pos = std.mem.indexOfScalar(i64, mixed, 0) orelse return error.NoZeroValue;
    const len = mixed.len;
    var sum: i64 = 0;
    inline for (.{ 1000, 2000, 3000 }) |offset| {
        sum += mixed[(zero_pos + offset) % len];
    }
    return sum;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const numbers = try parse(allocator, input);
    defer allocator.free(numbers);

    const mixed_once = try mix(allocator, numbers, 1);
    defer allocator.free(mixed_once);
    const part1 = try groveCoordinatesSum(mixed_once);

    const decrypted = try allocator.alloc(i64, numbers.len);
    defer allocator.free(decrypted);
    for (numbers, 0..) |value, i| decrypted[i] = value * decryption_key;

    const mixed_many = try mix(allocator, decrypted, mix_rounds_part2);
    defer allocator.free(mixed_many);
    const part2 = try groveCoordinatesSum(mixed_many);

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

const testing = std.testing;

test "part1 example" {
    const input = @embedFile("examples/day_20.txt");
    const numbers = try parse(testing.allocator, input);
    defer testing.allocator.free(numbers);
    const mixed = try mix(testing.allocator, numbers, 1);
    defer testing.allocator.free(mixed);
    try testing.expectEqual(@as(i64, 3), try groveCoordinatesSum(mixed));
}

test "part2 example" {
    const input = @embedFile("examples/day_20.txt");
    const numbers = try parse(testing.allocator, input);
    defer testing.allocator.free(numbers);

    const decrypted = try testing.allocator.alloc(i64, numbers.len);
    defer testing.allocator.free(decrypted);
    for (numbers, 0..) |value, i| decrypted[i] = value * decryption_key;

    const mixed = try mix(testing.allocator, decrypted, mix_rounds_part2);
    defer testing.allocator.free(mixed);
    try testing.expectEqual(@as(i64, 1623178306), try groveCoordinatesSum(mixed));
}

test "solve example" {
    const input = @embedFile("examples/day_20.txt");
    const answers = try solve(testing.allocator, input);
    try testing.expectEqual(@as(i64, 3), answers.part1);
    try testing.expectEqual(@as(i64, 1623178306), answers.part2);
}

test "single element edge case" {
    const numbers = try parse(testing.allocator, "0\n");
    defer testing.allocator.free(numbers);
    const mixed = try mix(testing.allocator, numbers, 1);
    defer testing.allocator.free(mixed);
    try testing.expectEqualSlices(i64, &[_]i64{0}, mixed);
}

test "duplicate values are tracked by original index" {
    // Values repeat; mixing must still move each occurrence independently.
    const numbers = try parse(testing.allocator, "1\n1\n0\n-1\n");
    defer testing.allocator.free(numbers);
    const mixed = try mix(testing.allocator, numbers, 1);
    defer testing.allocator.free(mixed);
    try testing.expectEqual(@as(usize, 4), mixed.len);
    var sum: i64 = 0;
    for (mixed) |v| sum += v;
    try testing.expectEqual(@as(i64, 1), sum);
}

test "malformed input errors" {
    try testing.expectError(error.InvalidCharacter, parse(testing.allocator, "not-a-number\n"));
}

test "empty input errors" {
    try testing.expectError(error.EmptyInput, parse(testing.allocator, ""));
}
