const std = @import("std");

const Answers = struct {
    part1: i64,
    part2: i64,
};

const Direction = struct { dx: i32, dy: i32 };

const all_directions = [_]Direction{
    .{ .dx = 1, .dy = 0 }, .{ .dx = -1, .dy = 0 }, .{ .dx = 0, .dy = 1 },  .{ .dx = 0, .dy = -1 },
    .{ .dx = 1, .dy = 1 }, .{ .dx = 1, .dy = -1 }, .{ .dx = -1, .dy = 1 }, .{ .dx = -1, .dy = -1 },
};

const Grid = struct {
    rows: std.ArrayList([]const u8),

    fn deinit(self: *Grid, allocator: std.mem.Allocator) void {
        self.rows.deinit(allocator);
    }

    fn width(self: Grid) usize {
        return self.rows.items[0].len;
    }

    fn height(self: Grid) usize {
        return self.rows.items.len;
    }

    fn at(self: Grid, row: i32, col: i32) ?u8 {
        if (row < 0 or col < 0 or row >= self.height() or col >= self.width()) return null;
        return self.rows.items[@intCast(row)][@intCast(col)];
    }
};

fn parseGrid(allocator: std.mem.Allocator, input: []const u8) !Grid {
    var rows: std.ArrayList([]const u8) = .empty;
    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw_line| {
        const line = std.mem.trimEnd(u8, raw_line, "\r");
        if (line.len == 0) continue;
        try rows.append(allocator, line);
    }
    if (rows.items.len == 0) {
        rows.deinit(allocator);
        return error.EmptyGrid;
    }
    const width = rows.items[0].len;
    for (rows.items) |row| {
        if (row.len != width) {
            rows.deinit(allocator);
            return error.InconsistentWidth;
        }
    }
    return Grid{ .rows = rows };
}

fn matchesWord(grid: Grid, row: i32, col: i32, dir: Direction, word: []const u8) bool {
    for (word, 0..) |expected, i| {
        const r = row + dir.dy * @as(i32, @intCast(i));
        const c = col + dir.dx * @as(i32, @intCast(i));
        const actual = grid.at(r, c) orelse return false;
        if (actual != expected) return false;
    }
    return true;
}

fn countXmas(grid: Grid) i64 {
    var count: i64 = 0;
    for (0..grid.height()) |row_u| {
        for (0..grid.width()) |col_u| {
            const row: i32 = @intCast(row_u);
            const col: i32 = @intCast(col_u);
            if (grid.at(row, col) != 'X') continue;
            for (all_directions) |dir| {
                if (matchesWord(grid, row, col, dir, "XMAS")) count += 1;
            }
        }
    }
    return count;
}

fn diagonalIsMas(grid: Grid, row: i32, col: i32, dir: Direction) bool {
    const a = grid.at(row - dir.dy, col - dir.dx) orelse return false;
    const b = grid.at(row + dir.dy, col + dir.dx) orelse return false;
    return (a == 'M' and b == 'S') or (a == 'S' and b == 'M');
}

fn countXMas(grid: Grid) i64 {
    var count: i64 = 0;
    for (0..grid.height()) |row_u| {
        for (0..grid.width()) |col_u| {
            const row: i32 = @intCast(row_u);
            const col: i32 = @intCast(col_u);
            if (grid.at(row, col) != 'A') continue;
            if (diagonalIsMas(grid, row, col, .{ .dx = 1, .dy = 1 }) and
                diagonalIsMas(grid, row, col, .{ .dx = 1, .dy = -1 })) count += 1;
        }
    }
    return count;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var grid = try parseGrid(allocator, input);
    defer grid.deinit(allocator);
    return Answers{ .part1 = countXmas(grid), .part2 = countXMas(grid) };
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

const example_input = @embedFile("examples/day4.txt");

test "official example" {
    const answers = try solve(std.testing.allocator, example_input);
    try std.testing.expectEqual(@as(i64, 18), answers.part1);
    try std.testing.expectEqual(@as(i64, 9), answers.part2);
}

test "single diagonal x-mas" {
    const input = "M.S\n.A.\nM.S\n";
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(i64, 1), answers.part2);
}

test "ragged grid is an error" {
    const input = "XMAS\nXM\n";
    try std.testing.expectError(error.InconsistentWidth, solve(std.testing.allocator, input));
}

test "empty input is an error" {
    try std.testing.expectError(error.EmptyGrid, solve(std.testing.allocator, ""));
}
