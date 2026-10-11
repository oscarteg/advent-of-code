const std = @import("std");

const Answers = struct {
    part1: i64,
    part2: i64,
};

fn parseLists(allocator: std.mem.Allocator, input: []const u8) ![2]std.ArrayList(i32) {
    var left: std.ArrayList(i32) = .empty;
    var right: std.ArrayList(i32) = .empty;
    errdefer left.deinit(allocator);
    errdefer right.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw_line| {
        const line = std.mem.trimEnd(u8, raw_line, "\r");
        if (line.len == 0) continue;

        var tokens = std.mem.tokenizeScalar(u8, line, ' ');
        const left_str = tokens.next() orelse return error.MissingLeftNumber;
        const right_str = tokens.next() orelse return error.MissingRightNumber;

        const l = try std.fmt.parseInt(i32, left_str, 10);
        const r = try std.fmt.parseInt(i32, right_str, 10);

        try left.append(allocator, l);
        try right.append(allocator, r);
    }

    if (left.items.len == 0) return error.EmptyInput;

    return [2]std.ArrayList(i32){ left, right };
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var lists = try parseLists(allocator, input);
    defer lists[0].deinit(allocator);
    defer lists[1].deinit(allocator);

    std.mem.sort(i32, lists[0].items, {}, std.sort.asc(i32));
    std.mem.sort(i32, lists[1].items, {}, std.sort.asc(i32));

    var distance: i64 = 0;
    for (lists[0].items, lists[1].items) |l, r| {
        const diff = @as(i64, l) - @as(i64, r);
        distance += if (diff < 0) -diff else diff;
    }

    var right_counts: std.AutoHashMap(i32, i32) = .init(allocator);
    defer right_counts.deinit();
    for (lists[1].items) |r| {
        const entry = try right_counts.getOrPutValue(r, 0);
        entry.value_ptr.* += 1;
    }

    var score: i64 = 0;
    for (lists[0].items) |l| {
        const count = right_counts.get(l) orelse 0;
        score += @as(i64, l) * count;
    }

    return Answers{ .part1 = distance, .part2 = score };
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
    const input = @embedFile("examples/day1.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(i64, 11), answers.part1);
    try std.testing.expectEqual(@as(i64, 31), answers.part2);
}

test "single pair" {
    const answers = try solve(std.testing.allocator, "5 10\n");
    try std.testing.expectEqual(@as(i64, 5), answers.part1);
    try std.testing.expectEqual(@as(i64, 0), answers.part2);
}

test "malformed input" {
    try std.testing.expectError(error.MissingRightNumber, solve(std.testing.allocator, "5\n"));
    try std.testing.expectError(error.EmptyInput, solve(std.testing.allocator, "\n"));
}
