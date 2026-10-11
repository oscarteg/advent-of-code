//! Advent of Code 2025 day 9: Movie Theater.
//!
//! The input lists the theater's red tiles as `x,y`, in the order they are
//! walked around the room. Part 1 uses any two red tiles as the opposite
//! corners of a rectangle and looks for the largest area, counting tiles
//! inclusively. Part 2 joins consecutive red tiles (wrapping around) with
//! straight green lines and fills the enclosed area green: the rectangle
//! still needs two red tiles as opposite corners, but may only cover red
//! and green tiles.

const std = @import("std");

const Answers = struct { part1: u64, part2: u64 };

const Tile = struct { x: i32, y: i32 };

/// Tiles of the rectangle spanned by two tiles, both corners included.
fn areaBetween(a: Tile, b: Tile) u64 {
    const width = @as(u64, @abs(a.x - b.x)) + 1;
    const height = @as(u64, @abs(a.y - b.y)) + 1;
    return width * height;
}

/// A closed loop of red tiles: consecutive tiles, wrapping around, are
/// joined by a horizontal or vertical line.
const Loop = struct {
    tiles: []const Tile,

    fn parse(allocator: std.mem.Allocator, input: []const u8) !Loop {
        var tiles: std.ArrayList(Tile) = .empty;
        errdefer tiles.deinit(allocator);

        var lines = std.mem.tokenizeScalar(u8, input, '\n');
        while (lines.next()) |line| {
            const trimmed = std.mem.trim(u8, line, &std.ascii.whitespace);
            if (trimmed.len == 0) continue;

            var fields = std.mem.splitScalar(u8, trimmed, ',');
            const x_field = fields.next() orelse return error.InvalidTile;
            const y_field = fields.next() orelse return error.InvalidTile;
            if (fields.next() != null) return error.InvalidTile;

            try tiles.append(allocator, .{
                .x = try std.fmt.parseInt(i32, std.mem.trim(u8, x_field, &std.ascii.whitespace), 10),
                .y = try std.fmt.parseInt(i32, std.mem.trim(u8, y_field, &std.ascii.whitespace), 10),
            });
        }
        if (tiles.items.len < 4) return error.NotALoop;

        const loop = Loop{ .tiles = try tiles.toOwnedSlice(allocator) };
        errdefer loop.deinit(allocator);

        var previous = loop.tiles[loop.tiles.len - 1];
        for (loop.tiles) |tile| {
            const straight_line = (tile.x == previous.x) != (tile.y == previous.y);
            if (!straight_line) return error.NotALoop;
            previous = tile;
        }
        return loop;
    }

    fn deinit(self: Loop, allocator: std.mem.Allocator) void {
        allocator.free(self.tiles);
    }
};

/// Sorted grid lines of one axis: every distinct coordinate, plus one
/// representative for each gap between neighbouring coordinates. All red and
/// green lines lie on coordinates of the loop, so the tiles a representative
/// stands in for share its colour. Sorts `values` in place.
fn stripLines(allocator: std.mem.Allocator, values: []i32) ![]i32 {
    std.mem.sortUnstable(i32, values, {}, std.sort.asc(i32));

    var lines: std.ArrayList(i32) = .empty;
    errdefer lines.deinit(allocator);
    for (values, 0..) |value, i| {
        if (i > 0) {
            const previous = values[i - 1];
            if (value == previous) continue;
            if (value > previous + 1) try lines.append(allocator, previous + 1);
        }
        try lines.append(allocator, value);
    }
    return lines.toOwnedSlice(allocator);
}

/// Index of the first line that is not below `value`.
fn lowerBound(lines: []const i32, value: i32) usize {
    var low: usize = 0;
    var high: usize = lines.len;
    while (low < high) {
        const middle = low + (high - low) / 2;
        if (lines[middle] < value) low = middle + 1 else high = middle;
    }
    return low;
}

/// Where the loop covers the floor, on a grid compressed to the loop's own
/// lines: `tile_column` and `tile_row` hold the strip indices of every tile,
/// and `uncovered` is a running count of uncovered strips, which turns the
/// rectangle between two red tiles into four lookups.
const Coverage = struct {
    tile_column: []usize,
    tile_row: []usize,
    stride: usize,
    uncovered: []u32,

    fn init(allocator: std.mem.Allocator, loop: Loop) !Coverage {
        const scratch = try allocator.alloc(i32, loop.tiles.len);
        defer allocator.free(scratch);

        for (loop.tiles, 0..) |tile, i| scratch[i] = tile.x;
        const columns = try stripLines(allocator, scratch);
        defer allocator.free(columns);

        for (loop.tiles, 0..) |tile, i| scratch[i] = tile.y;
        const rows = try stripLines(allocator, scratch);
        defer allocator.free(rows);

        const tile_column = try allocator.alloc(usize, loop.tiles.len);
        errdefer allocator.free(tile_column);
        const tile_row = try allocator.alloc(usize, loop.tiles.len);
        errdefer allocator.free(tile_row);
        for (loop.tiles, 0..) |tile, i| {
            tile_column[i] = lowerBound(columns, tile.x);
            tile_row[i] = lowerBound(rows, tile.y);
        }

        const stride = columns.len + 1;
        const uncovered = try allocator.alloc(u32, (rows.len + 1) * stride);
        errdefer allocator.free(uncovered);
        @memset(uncovered[0..stride], 0);

        const crossings = try allocator.alloc(u32, columns.len);
        defer allocator.free(crossings);
        const on_line = try allocator.alloc(i32, stride);
        defer allocator.free(on_line);

        for (rows, 0..) |y, row| {
            @memset(crossings, 0);
            @memset(on_line, 0);
            var total_crossings: u32 = 0;

            var previous = loop.tiles.len - 1;
            for (loop.tiles, 0..) |tile, current| {
                const other = loop.tiles[previous];
                if (tile.x == other.x) {
                    const top = @min(other.y, tile.y);
                    const bottom = @max(other.y, tile.y);
                    if (y >= top and y <= bottom) {
                        on_line[tile_column[current]] += 1;
                        on_line[tile_column[current] + 1] -= 1;
                        // Half-open in y, so a corner counts as one crossing.
                        if (y < bottom) {
                            crossings[tile_column[current]] += 1;
                            total_crossings += 1;
                        }
                    }
                } else if (tile.y == y) {
                    on_line[@min(tile_column[previous], tile_column[current])] += 1;
                    on_line[@max(tile_column[previous], tile_column[current]) + 1] -= 1;
                }
                previous = current;
            }

            // A strip is covered when a line runs through it, or when an odd
            // number of vertical lines separates it from the right edge.
            uncovered[(row + 1) * stride] = 0;
            var crossings_left: u32 = 0;
            var lines_here: i32 = 0;
            for (0..columns.len) |column| {
                crossings_left += crossings[column];
                lines_here += on_line[column];
                const covered = lines_here > 0 or (total_crossings - crossings_left) % 2 == 1;
                const strip: u32 = if (covered) 0 else 1;
                uncovered[(row + 1) * stride + column + 1] = uncovered[(row + 1) * stride + column] +
                    (uncovered[row * stride + column + 1] - uncovered[row * stride + column]) + strip;
            }
        }
        return .{
            .tile_column = tile_column,
            .tile_row = tile_row,
            .stride = stride,
            .uncovered = uncovered,
        };
    }

    fn deinit(self: Coverage, allocator: std.mem.Allocator) void {
        allocator.free(self.tile_column);
        allocator.free(self.tile_row);
        allocator.free(self.uncovered);
    }

    /// Whether every tile between the two red tiles is red or green.
    fn isFullyCovered(self: Coverage, a: usize, b: usize) bool {
        const left = @min(self.tile_column[a], self.tile_column[b]);
        const right = @max(self.tile_column[a], self.tile_column[b]) + 1;
        const top = @min(self.tile_row[a], self.tile_row[b]) * self.stride;
        const bottom = (@max(self.tile_row[a], self.tile_row[b]) + 1) * self.stride;
        const strips = (self.uncovered[bottom + right] - self.uncovered[bottom + left]) -
            (self.uncovered[top + right] - self.uncovered[top + left]);
        return strips == 0;
    }
};

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const loop = try Loop.parse(allocator, input);
    defer loop.deinit(allocator);

    const coverage = try Coverage.init(allocator, loop);
    defer coverage.deinit(allocator);

    var answers = Answers{ .part1 = 0, .part2 = 0 };
    for (loop.tiles, 0..) |a, i| {
        for (loop.tiles[i + 1 ..], i + 1..) |b, j| {
            const area = areaBetween(a, b);
            answers.part1 = @max(answers.part1, area);
            if (area > answers.part2 and coverage.isFullyCovered(i, j)) answers.part2 = area;
        }
    }
    return answers;
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
    const answers = try solve(std.testing.allocator, @embedFile("examples/day9.txt"));
    try std.testing.expectEqual(@as(u64, 50), answers.part1);
    try std.testing.expectEqual(@as(u64, 24), answers.part2);
}

test "a square loop covers every tile it spans" {
    const answers = try solve(std.testing.allocator, "0,0\n4,0\n4,4\n0,4\n");
    try std.testing.expectEqual(@as(u64, 25), answers.part1);
    try std.testing.expectEqual(@as(u64, 25), answers.part2);
}

test "an L-shaped loop rules out rectangles reaching into the notch" {
    const answers = try solve(std.testing.allocator, "0,0\n4,0\n4,2\n2,2\n2,4\n0,4\n");
    try std.testing.expectEqual(@as(u64, 25), answers.part1);
    try std.testing.expectEqual(@as(u64, 15), answers.part2);
}

test "a comb loop keeps rectangles out of its gaps" {
    const input = "0,0\n9,0\n9,6\n7,6\n7,2\n5,2\n5,6\n3,6\n3,2\n1,2\n1,6\n0,6\n";
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u64, 70), answers.part1);
    try std.testing.expectEqual(@as(u64, 27), answers.part2);
}

test "surrounding whitespace is ignored" {
    const answers = try solve(std.testing.allocator, "  0,0 \n\n 4,0\n4,4\n 0,4 \n\n");
    try std.testing.expectEqual(@as(u64, 25), answers.part1);
    try std.testing.expectEqual(@as(u64, 25), answers.part2);
}

test "rejects a tile missing a coordinate" {
    try std.testing.expectError(error.InvalidTile, solve(std.testing.allocator, "0,0\n4,0\n4\n0,4\n"));
}

test "rejects a tile with a third coordinate" {
    try std.testing.expectError(error.InvalidTile, solve(std.testing.allocator, "0,0\n4,0,1\n4,4\n0,4\n"));
}

test "rejects a non-numeric coordinate" {
    try std.testing.expectError(error.InvalidCharacter, solve(std.testing.allocator, "0,0\n4,x\n4,4\n0,4\n"));
}

test "rejects tiles that do not form a loop" {
    try std.testing.expectError(error.NotALoop, solve(std.testing.allocator, "1,2\n3,4\n"));
    try std.testing.expectError(error.NotALoop, solve(std.testing.allocator, "  \n \n"));
    try std.testing.expectError(error.NotALoop, solve(std.testing.allocator, "0,0\n4,1\n4,4\n0,4\n"));
}
