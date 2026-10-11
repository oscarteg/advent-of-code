const std = @import("std");

const Grid = struct {
    trees: [][]u32,
    rows: usize,
    cols: usize,
    arena: std.mem.Allocator,

    fn parse(allocator: std.mem.Allocator, input: []const u8) !Grid {
        var lines = std.mem.splitScalar(u8, input, '\n');
        var trees = std.ArrayList([]u32).empty;
        defer trees.deinit(allocator);

        while (lines.next()) |line| {
            const trimmed = std.mem.trim(u8, line, " \t\r");
            if (trimmed.len == 0) continue;

            var row = std.ArrayList(u32).empty;
            defer row.deinit(allocator);

            for (trimmed) |ch| {
                if (ch >= '0' and ch <= '9') {
                    try row.append(allocator, ch - '0');
                } else {
                    return error.InvalidDigit;
                }
            }

            try trees.append(allocator, try row.toOwnedSlice(allocator));
        }

        if (trees.items.len == 0) return error.EmptyGrid;

        const rows = trees.items.len;
        const cols = trees.items[0].len;

        for (trees.items) |row| {
            if (row.len != cols) return error.InconsistentRowLength;
        }

        return Grid{
            .trees = try trees.toOwnedSlice(allocator),
            .rows = rows,
            .cols = cols,
            .arena = allocator,
        };
    }

    fn isVisible(self: *const Grid, r: usize, c: usize) bool {
        const height = self.trees[r][c];

        // Check left
        {
            var col: usize = 0;
            while (col < c) : (col += 1) {
                if (self.trees[r][col] >= height) break;
            } else return true;
        }

        // Check right
        {
            var col = c + 1;
            while (col < self.cols) : (col += 1) {
                if (self.trees[r][col] >= height) break;
            } else return true;
        }

        // Check up
        {
            var row: usize = 0;
            while (row < r) : (row += 1) {
                if (self.trees[row][c] >= height) break;
            } else return true;
        }

        // Check down
        {
            var row = r + 1;
            while (row < self.rows) : (row += 1) {
                if (self.trees[row][c] >= height) break;
            } else return true;
        }

        return false;
    }

    fn scenicScore(self: *const Grid, r: usize, c: usize) usize {
        const height = self.trees[r][c];

        var left: usize = 0;
        if (c > 0) {
            var col = c - 1;
            while (true) {
                left += 1;
                if (self.trees[r][col] >= height) break;
                if (col == 0) break;
                col -= 1;
            }
        }

        var right: usize = 0;
        var col = c + 1;
        while (col < self.cols) : (col += 1) {
            right += 1;
            if (self.trees[r][col] >= height) break;
        }

        var up: usize = 0;
        if (r > 0) {
            var row = r - 1;
            while (true) {
                up += 1;
                if (self.trees[row][c] >= height) break;
                if (row == 0) break;
                row -= 1;
            }
        }

        var down: usize = 0;
        var row = r + 1;
        while (row < self.rows) : (row += 1) {
            down += 1;
            if (self.trees[row][c] >= height) break;
        }

        return left * right * up * down;
    }
};

fn solve(allocator: std.mem.Allocator, input: []const u8) ![2]usize {
    var grid = try Grid.parse(allocator, input);
    defer {
        for (grid.trees) |row| {
            allocator.free(row);
        }
        allocator.free(grid.trees);
    }

    var visible_count: usize = 0;
    var max_score: usize = 0;

    var r: usize = 0;
    while (r < grid.rows) : (r += 1) {
        var c: usize = 0;
        while (c < grid.cols) : (c += 1) {
            if (grid.isVisible(r, c)) {
                visible_count += 1;
            }
            const score = grid.scenicScore(r, c);
            if (score > max_score) {
                max_score = score;
            }
        }
    }

    return [2]usize{ visible_count, max_score };
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
    try stdout.print("Part 1: {d}\nPart 2: {d}\n", .{ answers[0], answers[1] });
    try stdout.flush();
}

test "part1_example" {
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day_8.txt");
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(answers[0], 21);
}

test "part2_example" {
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day_8.txt");
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(answers[1], 8);
}

test "single_tree" {
    const allocator = std.testing.allocator;
    const input = "5\n";
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(answers[0], 1);
    try std.testing.expectEqual(answers[1], 0);
}

test "malformed_input" {
    const allocator = std.testing.allocator;
    const input = "12a34\n";
    const result = solve(allocator, input);
    try std.testing.expectError(error.InvalidDigit, result);
}
