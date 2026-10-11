const std = @import("std");

fn trim(s: []const u8) []const u8 {
    var start: usize = 0;
    var end: usize = s.len;

    while (start < end and (s[start] == ' ' or s[start] == '\t')) {
        start += 1;
    }
    while (end > start and (s[end - 1] == ' ' or s[end - 1] == '\t')) {
        end -= 1;
    }

    return s[start..end];
}

fn count_ways(design: []const u8, patterns: []const []const u8, arena: std.mem.Allocator, memo: *std.StringHashMap(u64)) !u64 {
    if (design.len == 0) {
        return 1;
    }

    if (memo.get(design)) |cached| {
        return cached;
    }

    var ways: u64 = 0;
    for (patterns) |pattern| {
        if (design.len >= pattern.len and std.mem.eql(u8, design[0..pattern.len], pattern)) {
            ways += try count_ways(design[pattern.len..], patterns, arena, memo);
        }
    }

    try memo.put(try arena.dupe(u8, design), ways);
    return ways;
}

fn solve(arena: std.mem.Allocator, input: []const u8) !struct { part1: u64, part2: u64 } {
    var lines: std.ArrayList([]const u8) = .empty;
    defer lines.deinit(arena);

    var line_iter = std.mem.splitScalar(u8, input, '\n');
    while (line_iter.next()) |line| {
        const trimmed = std.mem.trim(u8, line, "\r\t ");
        if (trimmed.len > 0) {
            try lines.append(arena, trimmed);
        }
    }

    if (lines.items.len == 0) {
        return error.EmptyInput;
    }

    var patterns: std.ArrayList([]const u8) = .empty;
    defer patterns.deinit(arena);

    var pattern_iter = std.mem.splitScalar(u8, lines.items[0], ',');
    while (pattern_iter.next()) |p| {
        const trimmed_pattern = trim(p);
        if (trimmed_pattern.len > 0) {
            try patterns.append(arena, try arena.dupe(u8, trimmed_pattern));
        }
    }

    var part1: u64 = 0;
    var part2: u64 = 0;

    for (lines.items[1..]) |design| {
        var memo = std.StringHashMap(u64).init(arena);
        defer memo.deinit();

        const ways = try count_ways(design, patterns.items, arena, &memo);
        if (ways > 0) {
            part1 += 1;
            part2 += ways;
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

test "example" {
    var arena_instance = std.heap.ArenaAllocator.init(std.testing.allocator);
    defer arena_instance.deinit();
    const arena = arena_instance.allocator();

    const example_input = @embedFile("examples/day19.txt");
    const answers = try solve(arena, example_input);
    try std.testing.expectEqual(answers.part1, 6);
    try std.testing.expectEqual(answers.part2, 16);
}

test "single pattern" {
    var arena_instance = std.heap.ArenaAllocator.init(std.testing.allocator);
    defer arena_instance.deinit();
    const arena = arena_instance.allocator();

    const input = "a\na\naa\naaa";
    const answers = try solve(arena, input);
    try std.testing.expectEqual(answers.part1, 3);
    try std.testing.expectEqual(answers.part2, 3);
}

test "impossible design" {
    var arena_instance = std.heap.ArenaAllocator.init(std.testing.allocator);
    defer arena_instance.deinit();
    const arena = arena_instance.allocator();

    const input = "a, b\nc";
    const answers = try solve(arena, input);
    try std.testing.expectEqual(answers.part1, 0);
    try std.testing.expectEqual(answers.part2, 0);
}

test "empty input" {
    var arena_instance = std.heap.ArenaAllocator.init(std.testing.allocator);
    defer arena_instance.deinit();
    const arena = arena_instance.allocator();

    const result = solve(arena, "");
    try std.testing.expectError(error.EmptyInput, result);
}
