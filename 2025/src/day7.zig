const std = @import("std");

/// Laboratories: a beam starts at 'S' and travels straight down. Hitting a
/// splitter ('^') stops that beam and emits two new beams, immediately left
/// and right of the splitter, which continue downward from the next row.
/// Part 1 counts how many times a splitter fires (once per row/column it
/// is hit, regardless of how many merged beams hit it together). Part 2
/// treats every split as a timeline fork and counts how many timelines
/// (beams, with multiplicity) remain once they reach the bottom row.
const Answers = struct { part1: u64, part2: u64 };

const Grid = struct {
    rows: [][]const u8,
    width: usize,
    start_row: usize,
    start_col: usize,
};

fn parseGrid(allocator: std.mem.Allocator, input: []const u8) !Grid {
    var rows: std.ArrayList([]const u8) = .empty;
    defer rows.deinit(allocator);

    var width: ?usize = null;
    var lines = std.mem.tokenizeScalar(u8, std.mem.trim(u8, input, &std.ascii.whitespace), '\n');
    while (lines.next()) |line_raw| {
        const line = std.mem.trim(u8, line_raw, &std.ascii.whitespace);
        if (line.len == 0) continue;
        if (width) |w| {
            if (line.len != w) return error.InvalidGrid;
        } else {
            width = line.len;
        }
        for (line) |ch| {
            if (ch != '.' and ch != '^' and ch != 'S') return error.InvalidGrid;
        }
        try rows.append(allocator, line);
    }
    if (rows.items.len == 0) return error.InvalidGrid;

    var start_row: ?usize = null;
    var start_col: ?usize = null;
    for (rows.items, 0..) |row, r| {
        for (row, 0..) |ch, c| {
            if (ch != 'S') continue;
            if (start_row != null) return error.InvalidGrid; // more than one start
            start_row = r;
            start_col = c;
        }
    }
    if (start_row == null) return error.InvalidGrid; // no start

    return .{
        .rows = try rows.toOwnedSlice(allocator),
        .width = width.?,
        .start_row = start_row.?,
        .start_col = start_col.?,
    };
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const grid = try parseGrid(allocator, input);
    defer allocator.free(grid.rows);

    var counts: std.AutoHashMap(usize, u64) = .init(allocator);
    defer counts.deinit();
    try counts.put(grid.start_col, 1);

    var part1: u64 = 0;
    var row: usize = grid.start_row + 1;
    while (row < grid.rows.len) : (row += 1) {
        var next: std.AutoHashMap(usize, u64) = .init(allocator);

        var it = counts.iterator();
        while (it.next()) |entry| {
            const col = entry.key_ptr.*;
            const count = entry.value_ptr.*;
            if (grid.rows[row][col] == '^') {
                part1 += 1;
                if (col > 0) try addCount(&next, col - 1, count);
                if (col + 1 < grid.width) try addCount(&next, col + 1, count);
            } else {
                try addCount(&next, col, count);
            }
        }

        counts.deinit();
        counts = next;
    }

    var part2: u64 = 0;
    var it = counts.iterator();
    while (it.next()) |entry| part2 += entry.value_ptr.*;

    return .{ .part1 = part1, .part2 = part2 };
}

fn addCount(map: *std.AutoHashMap(usize, u64), col: usize, amount: u64) !void {
    const entry = try map.getOrPutValue(col, 0);
    entry.value_ptr.* += amount;
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
    const input = @embedFile("examples/day7.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u64, 21), answers.part1);
    try std.testing.expectEqual(@as(u64, 40), answers.part2);
}

test "no splitters leaves a single straight beam" {
    const input = ".S.\n...\n...\n";
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u64, 0), answers.part1);
    try std.testing.expectEqual(@as(u64, 1), answers.part2);
}

test "a splitter at the grid edge loses both offshoots" {
    const input = "S\n^\n";
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u64, 1), answers.part1);
    try std.testing.expectEqual(@as(u64, 0), answers.part2);
}

test "rejects a grid with no start" {
    const input = "...\n^..\n";
    try std.testing.expectError(error.InvalidGrid, solve(std.testing.allocator, input));
}

test "rejects a grid with more than one start" {
    const input = "S.S\n...\n";
    try std.testing.expectError(error.InvalidGrid, solve(std.testing.allocator, input));
}

test "rejects rows of inconsistent width" {
    const input = "S..\n^.\n";
    try std.testing.expectError(error.InvalidGrid, solve(std.testing.allocator, input));
}

test "rejects invalid characters" {
    const input = "S.x\n^..\n";
    try std.testing.expectError(error.InvalidGrid, solve(std.testing.allocator, input));
}

test "rejects empty input" {
    try std.testing.expectError(error.InvalidGrid, solve(std.testing.allocator, "   \n  "));
}
