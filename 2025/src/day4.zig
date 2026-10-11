const std = @import("std");

/// Printing Department: a grid of paper rolls ('@') and empty space ('.').
/// Part 1 counts rolls with fewer than 4 of their 8 neighbours occupied
/// (accessible to a forklift). Part 2 repeatedly removes every accessible
/// roll, simultaneously, until none remain, and totals the removals.
const Answers = struct { part1: usize, part2: usize };

const Grid = [][]u8;

fn parseGrid(allocator: std.mem.Allocator, input: []const u8) !Grid {
    var rows: std.ArrayList([]u8) = .empty;
    errdefer {
        for (rows.items) |row| allocator.free(row);
        rows.deinit(allocator);
    }

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
            if (ch != '@' and ch != '.') return error.InvalidGrid;
        }
        const row = try allocator.dupe(u8, line);
        try rows.append(allocator, row);
    }
    if (rows.items.len == 0) return error.InvalidGrid;

    return try rows.toOwnedSlice(allocator);
}

fn freeGrid(allocator: std.mem.Allocator, grid: Grid) void {
    for (grid) |row| allocator.free(row);
    allocator.free(grid);
}

fn dupeGrid(allocator: std.mem.Allocator, grid: Grid) !Grid {
    var rows: std.ArrayList([]u8) = .empty;
    errdefer {
        for (rows.items) |row| allocator.free(row);
        rows.deinit(allocator);
    }
    for (grid) |row| {
        try rows.append(allocator, try allocator.dupe(u8, row));
    }
    return try rows.toOwnedSlice(allocator);
}

fn countNeighbours(grid: Grid, row: usize, col: usize) u8 {
    const height = grid.len;
    const width = grid[0].len;
    var count: u8 = 0;
    var dr: i32 = -1;
    while (dr <= 1) : (dr += 1) {
        var dc: i32 = -1;
        while (dc <= 1) : (dc += 1) {
            if (dr == 0 and dc == 0) continue;
            const r = @as(i32, @intCast(row)) + dr;
            const c = @as(i32, @intCast(col)) + dc;
            if (r < 0 or c < 0 or r >= @as(i32, @intCast(height)) or c >= @as(i32, @intCast(width))) continue;
            if (grid[@intCast(r)][@intCast(c)] == '@') count += 1;
        }
    }
    return count;
}

fn isAccessible(grid: Grid, row: usize, col: usize) bool {
    return grid[row][col] == '@' and countNeighbours(grid, row, col) < 4;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const grid = try parseGrid(allocator, input);
    defer freeGrid(allocator, grid);
    const height = grid.len;
    const width = grid[0].len;

    var part1: usize = 0;
    for (0..height) |row| {
        for (0..width) |col| {
            if (isAccessible(grid, row, col)) part1 += 1;
        }
    }

    var working = try dupeGrid(allocator, grid);
    defer freeGrid(allocator, working);

    var total: usize = 0;
    while (true) {
        var removable: std.ArrayList([2]usize) = .empty;
        defer removable.deinit(allocator);
        for (0..height) |row| {
            for (0..width) |col| {
                if (isAccessible(working, row, col)) try removable.append(allocator, .{ row, col });
            }
        }
        if (removable.items.len == 0) break;
        for (removable.items) |rc| working[rc[0]][rc[1]] = '.';
        total += removable.items.len;
    }

    return .{ .part1 = part1, .part2 = total };
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
    const input = @embedFile("examples/day4.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(usize, 13), answers.part1);
    try std.testing.expectEqual(@as(usize, 43), answers.part2);
}

test "fully isolated rolls are all accessible and removed in one pass" {
    const input = "@.@\n...\n@.@\n";
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(usize, 4), answers.part1);
    try std.testing.expectEqual(@as(usize, 4), answers.part2);
}

test "a fully packed block is never accessible" {
    const input = "@@@@\n@@@@\n@@@@\n@@@@\n";
    const answers = try solve(std.testing.allocator, input);
    // Only the 4 corners have fewer than 4 occupied neighbours, and removing
    // them leaves every remaining roll with at least 4 neighbours, so a
    // single pass removes exactly the same 4 rolls.
    try std.testing.expectEqual(@as(usize, 4), answers.part1);
    try std.testing.expectEqual(@as(usize, 4), answers.part2);
}

test "rejects rows of inconsistent width" {
    const input = "@@@\n@@\n";
    try std.testing.expectError(error.InvalidGrid, solve(std.testing.allocator, input));
}

test "rejects invalid characters" {
    const input = "@.x\n@.@\n";
    try std.testing.expectError(error.InvalidGrid, solve(std.testing.allocator, input));
}

test "rejects empty input" {
    try std.testing.expectError(error.InvalidGrid, solve(std.testing.allocator, "   \n  "));
}
