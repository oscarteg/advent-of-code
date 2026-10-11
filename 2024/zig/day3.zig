const std = @import("std");

pub const Answers = struct {
    part1: i64,
    part2: i64,
};

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    if (input.len == 0) {
        return error.EmptyInput;
    }

    var part1_sum: i64 = 0;
    var part2_sum: i64 = 0;
    var enabled = true;
    var pos: usize = 0;

    while (pos < input.len) {
        const next_mul = std.mem.indexOf(u8, input[pos..], "mul(");
        const next_do = std.mem.indexOf(u8, input[pos..], "do()");
        const next_dont = std.mem.indexOf(u8, input[pos..], "don't()");

        var next_pos: ?usize = null;
        var action: u8 = 0;

        if (next_mul) |m| {
            next_pos = m;
            action = 0;
        }
        if (next_do) |d| {
            if (next_pos == null or d < next_pos.?) {
                next_pos = d;
                action = 1;
            }
        }
        if (next_dont) |d| {
            if (next_pos == null or d < next_pos.?) {
                next_pos = d;
                action = 2;
            }
        }

        if (next_pos == null) break;

        const actual_pos = pos + next_pos.?;

        if (action == 0) {
            const substr = input[actual_pos..];
            if (tryParseMul(substr)) |pair| {
                part1_sum += pair.x * pair.y;
                if (enabled) {
                    part2_sum += pair.x * pair.y;
                }
                pos = actual_pos + pair.len;
            } else {
                pos = actual_pos + 4;
            }
        } else if (action == 1) {
            enabled = true;
            pos = actual_pos + 4;
        } else if (action == 2) {
            enabled = false;
            pos = actual_pos + 7;
        }
    }

    _ = allocator;
    return Answers{ .part1 = part1_sum, .part2 = part2_sum };
}

const MulResult = struct {
    x: i64,
    y: i64,
    len: usize,
};

fn tryParseMul(substr: []const u8) ?MulResult {
    if (substr.len < 8) return null;
    if (!std.mem.startsWith(u8, substr, "mul(")) return null;

    var i: usize = 4;
    var x: i64 = 0;
    var digit_count: usize = 0;

    while (i < substr.len and digit_count < 3) {
        const c = substr[i];
        if (c >= '0' and c <= '9') {
            x = x * 10 + (c - '0');
            digit_count += 1;
            i += 1;
        } else {
            break;
        }
    }

    if (digit_count == 0 or digit_count > 3 or i >= substr.len) return null;
    if (substr[i] != ',') return null;

    i += 1;
    var y: i64 = 0;
    digit_count = 0;

    while (i < substr.len and digit_count < 3) {
        const c = substr[i];
        if (c >= '0' and c <= '9') {
            y = y * 10 + (c - '0');
            digit_count += 1;
            i += 1;
        } else {
            break;
        }
    }

    if (digit_count == 0 or digit_count > 3 or i >= substr.len) return null;
    if (substr[i] != ')') return null;

    return MulResult{ .x = x, .y = y, .len = i + 1 };
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

test "simple mul" {
    const input = "mul(2,3)mul(4,5)";
    const ans = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(ans.part1, 26);
    try std.testing.expectEqual(ans.part2, 26);
}

test "mul with do don't" {
    const input = "mul(2,3)don't()mul(4,5)do()mul(1,1)";
    const ans = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(ans.part1, 27);
    try std.testing.expectEqual(ans.part2, 7);
}

test "invalid mul patterns" {
    const input = "mul(1000,1000)mul(1,)mul(,1)mul(1 2)mul(1,2";
    const ans = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(ans.part1, 0);
    try std.testing.expectEqual(ans.part2, 0);
}

test "empty input throws" {
    const ans = solve(std.testing.allocator, "");
    try std.testing.expectError(error.EmptyInput, ans);
}

test "official part 1 example" {
    const ans = try solve(std.testing.allocator, @embedFile("examples/day3.txt"));
    try std.testing.expectEqual(ans.part1, 161);
}

test "official part 2 example" {
    const ans = try solve(std.testing.allocator, @embedFile("examples/day3-part2.txt"));
    try std.testing.expectEqual(ans.part1, 161);
    try std.testing.expectEqual(ans.part2, 48);
}
