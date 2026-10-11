const std = @import("std");

const Answers = struct {
    part1: u32,
    part2: u32,
};

fn parseReport(allocator: std.mem.Allocator, line: []const u8) !std.ArrayList(i32) {
    var levels: std.ArrayList(i32) = .empty;
    errdefer levels.deinit(allocator);

    var tokens = std.mem.tokenizeScalar(u8, line, ' ');
    while (tokens.next()) |token| {
        const value = try std.fmt.parseInt(i32, token, 10);
        try levels.append(allocator, value);
    }
    if (levels.items.len == 0) return error.EmptyReport;
    return levels;
}

fn isSafe(levels: []const i32) bool {
    if (levels.len < 2) return true;
    const increasing = levels[1] > levels[0];
    var i: usize = 1;
    while (i < levels.len) : (i += 1) {
        const diff = levels[i] - levels[i - 1];
        const magnitude = @abs(diff);
        if (magnitude < 1 or magnitude > 3) return false;
        if (increasing and diff < 0) return false;
        if (!increasing and diff > 0) return false;
    }
    return true;
}

fn isSafeWithDampener(allocator: std.mem.Allocator, levels: []const i32) !bool {
    if (isSafe(levels)) return true;
    var without: std.ArrayList(i32) = .empty;
    defer without.deinit(allocator);

    for (0..levels.len) |skip| {
        without.clearRetainingCapacity();
        for (levels, 0..) |level, i| {
            if (i != skip) try without.append(allocator, level);
        }
        if (isSafe(without.items)) return true;
    }
    return false;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var safe_count: u32 = 0;
    var safe_with_dampener_count: u32 = 0;

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw_line| {
        const line = std.mem.trimEnd(u8, raw_line, "\r");
        var levels = try parseReport(allocator, line);
        defer levels.deinit(allocator);

        if (isSafe(levels.items)) safe_count += 1;
        if (try isSafeWithDampener(allocator, levels.items)) safe_with_dampener_count += 1;
    }

    return Answers{ .part1 = safe_count, .part2 = safe_with_dampener_count };
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

test "official example" {
    const input = @embedFile("examples/day2.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u32, 2), answers.part1);
    try std.testing.expectEqual(@as(u32, 4), answers.part2);
}

test "single level report is safe for both parts" {
    const answers = try solve(std.testing.allocator, "5\n");
    try std.testing.expectEqual(@as(u32, 1), answers.part1);
    try std.testing.expectEqual(@as(u32, 1), answers.part2);
}

test "malformed input returns an error" {
    try std.testing.expectError(error.InvalidCharacter, solve(std.testing.allocator, "1 two 3\n"));
    try std.testing.expectError(error.EmptyReport, solve(std.testing.allocator, "   \n"));
}
