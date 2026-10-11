const std = @import("std");

const Answers = struct {
    part1: i64,
    part2: i64,
};

const Step = struct {
    row: isize,
    col: isize,
};

/// Facings in clockwise order, so a 90 degree turn is +-1 modulo four.
const facings = [4]Step{
    .{ .row = -1, .col = 0 },
    .{ .row = 0, .col = 1 },
    .{ .row = 1, .col = 0 },
    .{ .row = 0, .col = -1 },
};
const facing_count = facings.len;
const east = 1;
const step_score = 1;
const turn_score = 1000;
const unreachable_score: i64 = std.math.maxInt(i64);

/// One state of the search: a (tile, facing) pair and what reaching it costs.
const Move = struct {
    score: i64,
    state: usize,
};

fn byScore(_: void, a: Move, b: Move) std.math.Order {
    return std.math.order(a.score, b.score);
}

const Frontier = std.PriorityQueue(Move, void, byScore);

/// The maze as a row-major grid of tiles; start and end are tile indices.
const Maze = struct {
    tiles: []const u8,
    width: usize,
    height: usize,
    start: usize,
    end: usize,

    fn isWall(self: Maze, tile: usize) bool {
        return self.tiles[tile] == '#';
    }

    /// The tile one step along the facing, or against it when walking
    /// backwards; null when that leaves the maze or hits a wall.
    fn stepFrom(self: Maze, tile: usize, facing: usize, backwards: bool) ?usize {
        const step = facings[facing];
        const row = @as(isize, @intCast(tile / self.width)) + if (backwards) -step.row else step.row;
        const col = @as(isize, @intCast(tile % self.width)) + if (backwards) -step.col else step.col;
        if (row < 0 or col < 0) return null;
        const next_row: usize = @intCast(row);
        const next_col: usize = @intCast(col);
        if (next_row >= self.height or next_col >= self.width) return null;
        const next = next_row * self.width + next_col;
        return if (self.isWall(next)) null else next;
    }

    /// Dijkstra over (tile, facing) states, where a turn costs 1000 and a step
    /// forward 1. Walking backwards explores the reverse graph, which scores
    /// every state by its distance to the nearest source instead of from it.
    fn scoresFrom(self: Maze, allocator: std.mem.Allocator, sources: []const usize, backwards: bool) ![]i64 {
        const best = try allocator.alloc(i64, self.tiles.len * facing_count);
        errdefer allocator.free(best);
        @memset(best, unreachable_score);

        var frontier: Frontier = .empty;
        defer frontier.deinit(allocator);
        for (sources) |source| {
            best[source] = 0;
            try frontier.push(allocator, .{ .score = 0, .state = source });
        }

        while (frontier.pop()) |current| {
            if (current.score > best[current.state]) continue;
            const tile = current.state / facing_count;
            const facing = current.state % facing_count;

            var moves: [3]Move = undefined;
            moves[0] = .{ .score = current.score + turn_score, .state = tile * facing_count + (facing + 1) % facing_count };
            moves[1] = .{ .score = current.score + turn_score, .state = tile * facing_count + (facing + 3) % facing_count };
            var move_count: usize = 2;
            if (self.stepFrom(tile, facing, backwards)) |next| {
                moves[move_count] = .{ .score = current.score + step_score, .state = next * facing_count + facing };
                move_count += 1;
            }

            for (moves[0..move_count]) |move| {
                if (move.score >= best[move.state]) continue;
                best[move.state] = move.score;
                try frontier.push(allocator, move);
            }
        }
        return best;
    }
};

fn parse(allocator: std.mem.Allocator, input: []const u8) !Maze {
    var tiles: std.ArrayList(u8) = .empty;
    errdefer tiles.deinit(allocator);

    var width: usize = 0;
    var height: usize = 0;
    var start: ?usize = null;
    var end: ?usize = null;
    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw| {
        const line = std.mem.trim(u8, raw, " \r");
        if (line.len == 0) continue;
        if (height == 0) {
            width = line.len;
        } else if (line.len != width) {
            return error.RaggedMaze;
        }
        for (line, 0..) |tile, col| {
            switch (tile) {
                '#', '.' => {},
                'S' => {
                    if (start != null) return error.DuplicateStart;
                    start = height * width + col;
                },
                'E' => {
                    if (end != null) return error.DuplicateEnd;
                    end = height * width + col;
                },
                else => return error.InvalidTile,
            }
        }
        try tiles.appendSlice(allocator, line);
        height += 1;
    }
    if (height == 0) return error.EmptyMaze;
    const start_tile = start orelse return error.MissingStart;
    const end_tile = end orelse return error.MissingEnd;

    return Maze{
        .tiles = try tiles.toOwnedSlice(allocator),
        .width = width,
        .height = height,
        .start = start_tile,
        .end = end_tile,
    };
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const maze = try parse(allocator, input);
    defer allocator.free(maze.tiles);

    const start_states = [1]usize{maze.start * facing_count + east};
    const from_start = try maze.scoresFrom(allocator, &start_states, false);
    defer allocator.free(from_start);

    var end_states: [facing_count]usize = undefined;
    for (&end_states, 0..) |*state, facing| state.* = maze.end * facing_count + facing;
    const to_end = try maze.scoresFrom(allocator, &end_states, true);
    defer allocator.free(to_end);

    var best = unreachable_score;
    for (end_states) |state| best = @min(best, from_start[state]);
    if (best == unreachable_score) return error.NoPath;

    // A state lies on a best path when reaching it and leaving it towards the
    // end together cost exactly the best score.
    var tiles_on_best_paths: i64 = 0;
    for (0..maze.tiles.len) |tile| {
        for (0..facing_count) |facing| {
            const state = tile * facing_count + facing;
            if (from_start[state] == unreachable_score) continue;
            if (to_end[state] == unreachable_score) continue;
            if (from_start[state] + to_end[state] != best) continue;
            tiles_on_best_paths += 1;
            break;
        }
    }
    return Answers{ .part1 = best, .part2 = tiles_on_best_paths };
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
    const answers = try solve(std.testing.allocator, @embedFile("examples/day16.txt"));
    try std.testing.expectEqual(@as(i64, 7036), answers.part1);
    try std.testing.expectEqual(@as(i64, 45), answers.part2);
}

test "second official example" {
    const answers = try solve(std.testing.allocator,
        \\#################
        \\#...#...#...#..E#
        \\#.#.#.#.#.#.#.#.#
        \\#.#.#.#...#...#.#
        \\#.#.#.#.###.#.#.#
        \\#...#.#.#.....#.#
        \\#.#.#.#.#.#####.#
        \\#.#...#.#.#.....#
        \\#.#.#####.#.###.#
        \\#.#.#.......#...#
        \\#.#.###.#####.###
        \\#.#.#...#.....#.#
        \\#.#.#.#####.###.#
        \\#.#.#.........#.#
        \\#.#.#.#########.#
        \\#S#.............#
        \\#################
    );
    try std.testing.expectEqual(@as(i64, 11048), answers.part1);
    try std.testing.expectEqual(@as(i64, 64), answers.part2);
}

test "corridor costs one turn plus two steps" {
    const answers = try solve(std.testing.allocator,
        \\###
        \\#E#
        \\#.#
        \\#S#
        \\###
    );
    try std.testing.expectEqual(@as(i64, 1002), answers.part1);
    try std.testing.expectEqual(@as(i64, 3), answers.part2);
}

test "end next to the start costs one step" {
    const answers = try solve(std.testing.allocator, "####\n#SE#\n####\n");
    try std.testing.expectEqual(@as(i64, 1), answers.part1);
    try std.testing.expectEqual(@as(i64, 2), answers.part2);
}

test "malformed input returns an error" {
    const testing = std.testing;
    try testing.expectError(error.EmptyMaze, solve(testing.allocator, ""));
    try testing.expectError(error.RaggedMaze, solve(testing.allocator, "####\n#SE#\n###\n"));
    try testing.expectError(error.InvalidTile, solve(testing.allocator, "####\n#SX#\n####\n"));
    try testing.expectError(error.MissingStart, solve(testing.allocator, "####\n#.E#\n####\n"));
    try testing.expectError(error.MissingEnd, solve(testing.allocator, "####\n#S.#\n####\n"));
    try testing.expectError(error.DuplicateStart, solve(testing.allocator, "####\n#SS#\n####\n"));
    try testing.expectError(error.NoPath, solve(testing.allocator, "#####\n#S#E#\n#####\n"));
}
