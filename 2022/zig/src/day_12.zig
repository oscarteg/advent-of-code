const std = @import("std");
const testing = std.testing;

const Answers = struct {
    part1: u32,
    part2: u32,
};

const Error = error{
    EmptyGrid,
    RaggedRow,
    InvalidElevation,
    MissingStart,
    DuplicateStart,
    MissingEnd,
    DuplicateEnd,
    NoPath,
};

const Grid = struct {
    heights: []u8,
    width: usize,
    height: usize,
    start: usize,
    end: usize,

    fn parse(allocator: std.mem.Allocator, input: []const u8) !Grid {
        var rows: std.ArrayList([]const u8) = .empty;
        defer rows.deinit(allocator);
        var lines = std.mem.tokenizeScalar(u8, input, '\n');
        while (lines.next()) |line| try rows.append(allocator, line);
        if (rows.items.len == 0) return Error.EmptyGrid;

        const width = rows.items[0].len;
        if (width == 0) return Error.EmptyGrid;

        var heights = try allocator.alloc(u8, width * rows.items.len);
        errdefer allocator.free(heights);
        var start: ?usize = null;
        var end: ?usize = null;

        for (rows.items, 0..) |row, row_index| {
            if (row.len != width) return Error.RaggedRow;
            for (row, 0..) |ch, col_index| {
                const index = row_index * width + col_index;
                heights[index] = switch (ch) {
                    'S' => blk: {
                        if (start != null) return Error.DuplicateStart;
                        start = index;
                        break :blk 'a' - 'a';
                    },
                    'E' => blk: {
                        if (end != null) return Error.DuplicateEnd;
                        end = index;
                        break :blk 'z' - 'a';
                    },
                    'a'...'z' => ch - 'a',
                    else => return Error.InvalidElevation,
                };
            }
        }

        return .{
            .heights = heights,
            .width = width,
            .height = rows.items.len,
            .start = start orelse return Error.MissingStart,
            .end = end orelse return Error.MissingEnd,
        };
    }

    fn deinit(self: *Grid, allocator: std.mem.Allocator) void {
        allocator.free(self.heights);
    }

    fn neighbors(self: Grid, index: usize, buffer: *[4]usize) []usize {
        const row = index / self.width;
        const col = index % self.width;
        var count: usize = 0;
        if (row > 0) {
            buffer[count] = index - self.width;
            count += 1;
        }
        if (row + 1 < self.height) {
            buffer[count] = index + self.width;
            count += 1;
        }
        if (col > 0) {
            buffer[count] = index - 1;
            count += 1;
        }
        if (col + 1 < self.width) {
            buffer[count] = index + 1;
            count += 1;
        }
        return buffer[0..count];
    }

    // Breadth-first search from the end over the reversed graph: an edge
    // `current -> neighbor` here corresponds to the forward edge
    // `neighbor -> current`, valid when it climbs by at most one elevation.
    fn shortestPathFromEnd(self: Grid, allocator: std.mem.Allocator, starts: []const usize) !u32 {
        const unvisited = std.math.maxInt(u32);
        var distances = try allocator.alloc(u32, self.heights.len);
        defer allocator.free(distances);
        @memset(distances, unvisited);

        var queue: std.ArrayList(usize) = .empty;
        defer queue.deinit(allocator);
        distances[self.end] = 0;
        try queue.append(allocator, self.end);

        var head: usize = 0;
        while (head < queue.items.len) {
            const current = queue.items[head];
            head += 1;
            var buffer: [4]usize = undefined;
            for (self.neighbors(current, &buffer)) |neighbor| {
                if (distances[neighbor] != unvisited) continue;
                const drop = @as(i16, self.heights[current]) - @as(i16, self.heights[neighbor]);
                if (drop > 1) continue;
                distances[neighbor] = distances[current] + 1;
                try queue.append(allocator, neighbor);
            }
        }

        var best: ?u32 = null;
        for (starts) |start| {
            if (distances[start] == unvisited) continue;
            if (best == null or distances[start] < best.?) best = distances[start];
        }
        return best orelse Error.NoPath;
    }
};

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var grid = try Grid.parse(allocator, input);
    defer grid.deinit(allocator);

    const part1 = try grid.shortestPathFromEnd(allocator, &[_]usize{grid.start});

    var lowest: std.ArrayList(usize) = .empty;
    defer lowest.deinit(allocator);
    for (grid.heights, 0..) |elevation, index| {
        if (elevation == 0) try lowest.append(allocator, index);
    }
    const part2 = try grid.shortestPathFromEnd(allocator, lowest.items);

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

const example = @embedFile("examples/day_12.txt");

test "official example" {
    const answers = try solve(testing.allocator, example);
    try testing.expectEqual(@as(u32, 31), answers.part1);
    try testing.expectEqual(@as(u32, 29), answers.part2);
}

test "single row staircase" {
    const answers = try solve(testing.allocator, "SbcdefghijklmnopqrstuvwxyE");
    try testing.expectEqual(@as(u32, 25), answers.part1);
    try testing.expectEqual(@as(u32, 25), answers.part2);
}

test "ragged rows are rejected" {
    try testing.expectError(Error.RaggedRow, solve(testing.allocator, "Sab\nabcE"));
}

test "missing start is rejected" {
    try testing.expectError(Error.MissingStart, solve(testing.allocator, "abz\nabE"));
}

test "invalid elevation character is rejected" {
    try testing.expectError(Error.InvalidElevation, solve(testing.allocator, "S1\naE"));
}
