const std = @import("std");

const Direction = enum { up, down, left, right };

const Blizzard = struct {
    row: usize,
    col: usize,
    dir: Direction,
};

const Grid = struct {
    width: usize,
    height: usize,
    start_col: usize,
    end_col: usize,
    blizzards: []Blizzard,

    fn parse(allocator: std.mem.Allocator, input: []const u8) !Grid {
        var lines: std.ArrayList([]const u8) = .empty;
        defer lines.deinit(allocator);
        var line_it = std.mem.tokenizeScalar(u8, input, '\n');
        while (line_it.next()) |line| {
            try lines.append(allocator, line);
        }
        if (lines.items.len < 3) return error.InvalidFormat;

        const height = lines.items.len - 2;
        if (lines.items[0].len < 2) return error.InvalidFormat;
        const width = lines.items[0].len - 2;
        if (width == 0) return error.InvalidFormat;

        const top = lines.items[0];
        const bottom = lines.items[lines.items.len - 1];
        if (top.len != width + 2 or bottom.len != width + 2) return error.InvalidFormat;

        const start_idx = std.mem.indexOfScalar(u8, top, '.') orelse return error.InvalidFormat;
        if (start_idx == 0 or start_idx - 1 >= width) return error.InvalidFormat;
        const start_col = start_idx - 1;

        const end_idx = std.mem.indexOfScalar(u8, bottom, '.') orelse return error.InvalidFormat;
        if (end_idx == 0 or end_idx - 1 >= width) return error.InvalidFormat;
        const end_col = end_idx - 1;

        var blizzards: std.ArrayList(Blizzard) = .empty;
        errdefer blizzards.deinit(allocator);
        for (lines.items[1 .. 1 + height], 0..) |line, row| {
            if (line.len != width + 2) return error.InvalidFormat;
            if (line[0] != '#' or line[width + 1] != '#') return error.InvalidFormat;
            for (line[1 .. width + 1], 0..) |ch, col| {
                const dir: Direction = switch (ch) {
                    '^' => .up,
                    'v' => .down,
                    '<' => .left,
                    '>' => .right,
                    '.' => continue,
                    else => return error.InvalidFormat,
                };
                try blizzards.append(allocator, .{ .row = row, .col = col, .dir = dir });
            }
        }

        return .{
            .width = width,
            .height = height,
            .start_col = start_col,
            .end_col = end_col,
            .blizzards = try blizzards.toOwnedSlice(allocator),
        };
    }

    fn deinit(self: *Grid, allocator: std.mem.Allocator) void {
        allocator.free(self.blizzards);
    }

    fn occupiedAt(self: Grid, allocator: std.mem.Allocator, time: usize) !std.AutoHashMap(Position, void) {
        var set = std.AutoHashMap(Position, void).init(allocator);
        errdefer set.deinit();
        for (self.blizzards) |b| {
            const dr: i64, const dc: i64 = switch (b.dir) {
                .up => .{ -1, 0 },
                .down => .{ 1, 0 },
                .left => .{ 0, -1 },
                .right => .{ 0, 1 },
            };
            const height_i: i64 = @intCast(self.height);
            const width_i: i64 = @intCast(self.width);
            const row = @mod(@as(i64, @intCast(b.row)) + dr * @as(i64, @intCast(time)), height_i);
            const col = @mod(@as(i64, @intCast(b.col)) + dc * @as(i64, @intCast(time)), width_i);
            try set.put(.{ .row = @intCast(row), .col = @intCast(col) }, {});
        }
        return set;
    }

    fn isValid(self: Grid, row: i64, col: i64) bool {
        if (row == -1) return col == @as(i64, @intCast(self.start_col));
        if (row == @as(i64, @intCast(self.height))) return col == @as(i64, @intCast(self.end_col));
        return row >= 0 and col >= 0 and col < @as(i64, @intCast(self.width));
    }
};

const Position = struct { row: i64, col: i64 };
const State = struct { row: i64, col: i64, time_mod: usize };

fn gcd(a: usize, b: usize) usize {
    if (b == 0) return a;
    return gcd(b, a % b);
}

fn lcm(a: usize, b: usize) usize {
    return a / gcd(a, b) * b;
}

fn shortestTime(
    allocator: std.mem.Allocator,
    grid: Grid,
    occupied: []std.AutoHashMap(Position, void),
    period: usize,
    start_time: usize,
    from: Position,
    to: Position,
) !usize {
    var queue: std.ArrayList(struct { pos: Position, time: usize }) = .empty;
    defer queue.deinit(allocator);
    var visited = std.AutoHashMap(State, void).init(allocator);
    defer visited.deinit();

    try queue.append(allocator, .{ .pos = from, .time = start_time });
    try visited.put(.{ .row = from.row, .col = from.col, .time_mod = start_time % period }, {});

    var head: usize = 0;
    while (head < queue.items.len) {
        const current = queue.items[head];
        head += 1;
        if (current.pos.row == to.row and current.pos.col == to.col) return current.time;

        const next_time = current.time + 1;
        const occupied_next = &occupied[next_time % period];
        const deltas = [_][2]i64{ .{ 0, 0 }, .{ -1, 0 }, .{ 1, 0 }, .{ 0, -1 }, .{ 0, 1 } };
        for (deltas) |delta| {
            const next = Position{ .row = current.pos.row + delta[0], .col = current.pos.col + delta[1] };
            if (!grid.isValid(next.row, next.col)) continue;
            if (next.row >= 0 and next.row < @as(i64, @intCast(grid.height))) {
                if (occupied_next.contains(.{ .row = next.row, .col = next.col })) continue;
            }
            const state = State{ .row = next.row, .col = next.col, .time_mod = next_time % period };
            if (!visited.contains(state)) {
                try visited.put(state, {});
                try queue.append(allocator, .{ .pos = next, .time = next_time });
            }
        }
    }

    return error.NoPathFound;
}

const Answers = struct { part1: usize, part2: usize };

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var grid = try Grid.parse(allocator, input);
    defer grid.deinit(allocator);

    const period = lcm(grid.width, grid.height);
    var occupied = try allocator.alloc(std.AutoHashMap(Position, void), period);
    defer {
        for (occupied) |*set| set.deinit();
        allocator.free(occupied);
    }
    for (0..period) |t| {
        occupied[t] = try grid.occupiedAt(allocator, t);
    }

    const start = Position{ .row = -1, .col = @intCast(grid.start_col) };
    const end = Position{ .row = @intCast(grid.height), .col = @intCast(grid.end_col) };

    const there = try shortestTime(allocator, grid, occupied, period, 0, start, end);
    const back = try shortestTime(allocator, grid, occupied, period, there, end, start);
    const there_again = try shortestTime(allocator, grid, occupied, period, back, start, end);

    return .{ .part1 = there, .part2 = there_again };
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
    const input = @embedFile("examples/day_24.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(usize, 18), answers.part1);
    try std.testing.expectEqual(@as(usize, 54), answers.part2);
}

test "no blizzards is a direct walk" {
    const input = "#.#\n#.#\n#.#\n#.#\n#.#\n";
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(usize, 4), answers.part1);
    try std.testing.expect(answers.part2 > answers.part1);
}

test "missing openings is an error" {
    const input = "###\n#.#\n#.#\n#.#\n###\n";
    try std.testing.expectError(error.InvalidFormat, solve(std.testing.allocator, input));
}

test "unknown character is an error" {
    const input = "#.#\n#x#\n#.#\n#.#\n#.#\n";
    try std.testing.expectError(error.InvalidFormat, solve(std.testing.allocator, input));
}
