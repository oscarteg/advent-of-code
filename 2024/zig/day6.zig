const std = @import("std");

pub const Answers = struct {
    part1: i64,
    part2: i64,
};

const Direction = enum(u2) {
    up,
    right,
    down,
    left,

    fn turnRight(self: Direction) Direction {
        const next: u2 = @intFromEnum(self) +% 1;
        return @enumFromInt(next);
    }

    fn delta(self: Direction) struct { dr: i32, dc: i32 } {
        return switch (self) {
            .up => .{ .dr = -1, .dc = 0 },
            .right => .{ .dr = 0, .dc = 1 },
            .down => .{ .dr = 1, .dc = 0 },
            .left => .{ .dr = 0, .dc = -1 },
        };
    }

    fn fromChar(c: u8) !Direction {
        return switch (c) {
            '^' => .up,
            '>' => .right,
            'v' => .down,
            '<' => .left,
            else => error.InvalidGuardMarker,
        };
    }
};

const Grid = struct {
    rows: [][]u8,
    row_count: usize,
    col_count: usize,
    start_row: usize,
    start_col: usize,
    start_direction: Direction,

    fn inBounds(self: Grid, r: i32, c: i32) bool {
        return r >= 0 and c >= 0 and r < @as(i32, @intCast(self.row_count)) and
            c < @as(i32, @intCast(self.col_count));
    }

    fn isWall(self: Grid, r: usize, c: usize) bool {
        return self.rows[r][c] == '#';
    }
};

fn parseGrid(allocator: std.mem.Allocator, input: []const u8) !Grid {
    var rows: std.ArrayList([]u8) = .empty;
    defer rows.deinit(allocator);
    errdefer for (rows.items) |row| allocator.free(row);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        try rows.append(allocator, try allocator.dupe(u8, line));
    }

    if (rows.items.len == 0) return error.EmptyInput;

    const cols = rows.items[0].len;
    for (rows.items) |row| {
        if (row.len != cols) return error.IrregularGrid;
    }

    var found = false;
    var start_row: usize = 0;
    var start_col: usize = 0;
    var start_direction: Direction = .up;

    for (rows.items, 0..) |row, r| {
        for (row, 0..) |ch, c| {
            if (ch == '.' or ch == '#') continue;
            const direction = try Direction.fromChar(ch);
            if (found) return error.MultipleGuards;
            found = true;
            start_row = r;
            start_col = c;
            start_direction = direction;
            row[c] = '.';
        }
    }

    if (!found) return error.MissingGuard;

    const row_count = rows.items.len;
    return .{
        .rows = try rows.toOwnedSlice(allocator),
        .row_count = row_count,
        .col_count = cols,
        .start_row = start_row,
        .start_col = start_col,
        .start_direction = start_direction,
    };
}

fn freeGrid(allocator: std.mem.Allocator, grid: Grid) void {
    for (grid.rows) |row| allocator.free(row);
    allocator.free(grid.rows);
}

fn walk(
    allocator: std.mem.Allocator,
    grid: Grid,
    extra_wall_row: ?usize,
    extra_wall_col: ?usize,
    visited: ?*std.AutoHashMap([2]usize, void),
) !bool {
    const seen = try allocator.alloc(u8, grid.row_count * grid.col_count);
    defer allocator.free(seen);
    @memset(seen, 0);

    var r = grid.start_row;
    var c = grid.start_col;
    var direction = grid.start_direction;

    while (true) {
        const index = r * grid.col_count + c;
        const mask = @as(u8, 1) << @intFromEnum(direction);
        if (seen[index] & mask != 0) return true;
        seen[index] |= mask;
        if (visited) |v| try v.put(.{ r, c }, {});

        const d = direction.delta();
        const nr = @as(i32, @intCast(r)) + d.dr;
        const nc = @as(i32, @intCast(c)) + d.dc;
        if (!grid.inBounds(nr, nc)) return false;

        const ur: usize = @intCast(nr);
        const uc: usize = @intCast(nc);
        const blocked = grid.isWall(ur, uc) or
            (extra_wall_row != null and extra_wall_row.? == ur and extra_wall_col.? == uc);
        if (blocked) {
            direction = direction.turnRight();
        } else {
            r = ur;
            c = uc;
        }
    }
}

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const grid = try parseGrid(allocator, input);
    defer freeGrid(allocator, grid);

    var visited = std.AutoHashMap([2]usize, void).init(allocator);
    defer visited.deinit();
    _ = try walk(allocator, grid, null, null, &visited);

    var part2: i64 = 0;
    var it = visited.keyIterator();
    while (it.next()) |pos| {
        if (pos[0] == grid.start_row and pos[1] == grid.start_col) continue;
        if (try walk(allocator, grid, pos[0], pos[1], null)) part2 += 1;
    }

    return .{ .part1 = @intCast(visited.count()), .part2 = part2 };
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
    const example = @embedFile("examples/day6.txt");
    const answers = try solve(std.testing.allocator, example);
    try std.testing.expectEqual(41, answers.part1);
    try std.testing.expectEqual(6, answers.part2);
}

test "straight line" {
    const answers = try solve(std.testing.allocator, "...\n.^.\n...\n");
    try std.testing.expectEqual(2, answers.part1);
    try std.testing.expectEqual(0, answers.part2);
}

test "empty input" {
    try std.testing.expectError(error.EmptyInput, solve(std.testing.allocator, ""));
}

test "irregular grid" {
    try std.testing.expectError(error.IrregularGrid, solve(std.testing.allocator, "..\n...\n"));
}

test "invalid character" {
    try std.testing.expectError(error.InvalidGuardMarker, solve(std.testing.allocator, "..\n.?\n"));
}

test "missing guard" {
    try std.testing.expectError(error.MissingGuard, solve(std.testing.allocator, "...\n...\n"));
}

test "multiple guards" {
    try std.testing.expectError(error.MultipleGuards, solve(std.testing.allocator, "^.\n.^\n"));
}
