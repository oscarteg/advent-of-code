const std = @import("std");

const Answers = struct {
    part1: i64,
    part2: i64,
};

const modulus: i64 = 16777216;

fn nextSecret(secret: i64) i64 {
    var value = secret;
    value = @mod((value * 64) ^ value, modulus);
    value = @mod(@divFloor(value, 32) ^ value, modulus);
    value = @mod((value * 2048) ^ value, modulus);
    return value;
}

fn parseSeeds(allocator: std.mem.Allocator, input: []const u8) ![]i64 {
    var seeds: std.ArrayList(i64) = .empty;
    errdefer seeds.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw| {
        const line = std.mem.trim(u8, raw, " \r");
        if (line.len == 0) continue;
        try seeds.append(allocator, try std.fmt.parseInt(i64, line, 10));
    }
    if (seeds.items.len == 0) return error.InvalidInput;

    return seeds.toOwnedSlice(allocator);
}

fn encodeSequence(deltas: [4]i8) i32 {
    var key: i32 = 0;
    for (deltas) |delta| key = key * 19 + (@as(i32, delta) + 9);
    return key;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const seeds = try parseSeeds(allocator, input);
    defer allocator.free(seeds);

    var part1: i64 = 0;
    var total_bananas_by_sequence = std.AutoHashMap(i32, i64).init(allocator);
    defer total_bananas_by_sequence.deinit();

    for (seeds) |seed| {
        var secret = seed;
        var deltas = [4]i8{ 0, 0, 0, 0 };
        var previous_price: i8 = @intCast(@mod(secret, 10));
        var seen_by_sequence = std.AutoHashMap(i32, void).init(allocator);
        defer seen_by_sequence.deinit();

        var step: usize = 1;
        while (step <= 2000) : (step += 1) {
            secret = nextSecret(secret);
            const price: i8 = @intCast(@mod(secret, 10));
            deltas[0] = deltas[1];
            deltas[1] = deltas[2];
            deltas[2] = deltas[3];
            deltas[3] = price - previous_price;
            previous_price = price;

            if (step >= 4) {
                const key = encodeSequence(deltas);
                if (!seen_by_sequence.contains(key)) {
                    try seen_by_sequence.put(key, {});
                    const entry = try total_bananas_by_sequence.getOrPutValue(key, 0);
                    entry.value_ptr.* += price;
                }
            }
        }
        part1 += secret;
    }

    var part2: i64 = 0;
    var totals = total_bananas_by_sequence.valueIterator();
    while (totals.next()) |bananas| {
        if (bananas.* > part2) part2 = bananas.*;
    }

    return Answers{ .part1 = part1, .part2 = part2 };
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

test "official example part one" {
    const input = @embedFile("examples/day22.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(i64, 37327623), answers.part1);
}

test "part two best sequence" {
    const answers = try solve(std.testing.allocator, "1\n2\n3\n2024\n");
    try std.testing.expectEqual(@as(i64, 23), answers.part2);
}

test "single seed matches the official per-seed example" {
    const answers = try solve(std.testing.allocator, "1\n");
    try std.testing.expectEqual(@as(i64, 8685429), answers.part1);
}

test "malformed input returns an error" {
    try std.testing.expectError(error.InvalidInput, solve(std.testing.allocator, ""));
    try std.testing.expectError(error.InvalidCharacter, solve(std.testing.allocator, "not-a-number\n"));
}
