const std = @import("std");

const real_threshold = 100;

const Point = struct { x: i32, y: i32 };

const Answers = struct {
    part1: i64,
    part2: i64,
};

const Track = struct {
    width: i32,
    height: i32,
    wall: []const bool,
    start: Point,
    end: Point,

    fn index(self: Track, p: Point) usize {
        return @intCast(p.y * self.width + p.x);
    }

    fn contains(self: Track, p: Point) bool {
        return p.x >= 0 and p.x < self.width and p.y >= 0 and p.y < self.height;
    }
};

const Race = struct {
    track: Track,
    from_start: []i64,
    to_end: []i64,
    uncheated_time: i64,

    fn deinit(self: Race, allocator: std.mem.Allocator) void {
        allocator.free(self.track.wall);
        allocator.free(self.from_start);
        allocator.free(self.to_end);
    }
};

fn parse(allocator: std.mem.Allocator, input: []const u8) !Track {
    var wall: std.ArrayList(bool) = .empty;
    errdefer wall.deinit(allocator);

    var width: i32 = 0;
    var height: i32 = 0;
    var start: ?Point = null;
    var end: ?Point = null;

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw| {
        const line = std.mem.trim(u8, raw, " \r");
        if (line.len == 0) continue;

        const line_width: i32 = @intCast(line.len);
        if (width == 0) {
            width = line_width;
        } else if (line_width != width) {
            return error.RaggedMap;
        }

        for (line, 0..) |cell, column| {
            const here = Point{ .x = @intCast(column), .y = height };
            switch (cell) {
                '#', '.' => {},
                'S' => {
                    if (start != null) return error.DuplicateStart;
                    start = here;
                },
                'E' => {
                    if (end != null) return error.DuplicateEnd;
                    end = here;
                },
                else => return error.InvalidMapCharacter,
            }
            try wall.append(allocator, cell == '#');
        }
        height += 1;
    }

    if (height == 0) return error.EmptyMap;
    const start_point = start orelse return error.MissingStart;
    const end_point = end orelse return error.MissingEnd;

    return Track{
        .width = width,
        .height = height,
        .wall = try wall.toOwnedSlice(allocator),
        .start = start_point,
        .end = end_point,
    };
}

// Steps from origin to every track cell reachable from it, -1 where unreachable.
fn distancesFrom(allocator: std.mem.Allocator, track: Track, origin: Point) ![]i64 {
    const distance = try allocator.alloc(i64, track.wall.len);
    errdefer allocator.free(distance);
    @memset(distance, -1);

    var frontier: std.ArrayList(Point) = .empty;
    defer frontier.deinit(allocator);
    distance[track.index(origin)] = 0;
    try frontier.append(allocator, origin);

    const deltas = [4]Point{ .{ .x = 1, .y = 0 }, .{ .x = -1, .y = 0 }, .{ .x = 0, .y = 1 }, .{ .x = 0, .y = -1 } };

    var head: usize = 0;
    while (head < frontier.items.len) : (head += 1) {
        const current = frontier.items[head];
        const steps = distance[track.index(current)];
        for (deltas) |delta| {
            const next = Point{ .x = current.x + delta.x, .y = current.y + delta.y };
            if (!track.contains(next)) continue;
            const at = track.index(next);
            if (track.wall[at] or distance[at] != -1) continue;
            distance[at] = steps + 1;
            try frontier.append(allocator, next);
        }
    }
    return distance;
}

fn prepare(allocator: std.mem.Allocator, input: []const u8) !Race {
    const track = try parse(allocator, input);
    errdefer allocator.free(track.wall);

    const from_start = try distancesFrom(allocator, track, track.start);
    errdefer allocator.free(from_start);

    const uncheated_time = from_start[track.index(track.end)];
    if (uncheated_time < 0) return error.EndUnreachable;

    return Race{
        .track = track,
        .from_start = from_start,
        .to_end = try distancesFrom(allocator, track, track.end),
        .uncheated_time = uncheated_time,
    };
}

// A cheat ignores walls for at most max_cheat_time picoseconds, so landing on a track cell
// at Manhattan distance d costs exactly d picoseconds: the cheated race then takes
// from_start[begin] + d + to_end[landing]. Cheats are identified by those two cells.
fn countCheats(race: Race, max_cheat_time: i32, min_saving: i64) i64 {
    var cheats: i64 = 0;
    var y: i32 = 0;
    while (y < race.track.height) : (y += 1) {
        var x: i32 = 0;
        while (x < race.track.width) : (x += 1) {
            const to_cheat = race.from_start[race.track.index(.{ .x = x, .y = y })];
            if (to_cheat < 0) continue;

            var dy: i32 = -max_cheat_time;
            while (dy <= max_cheat_time) : (dy += 1) {
                const room = max_cheat_time - @as(i32, @intCast(@abs(dy)));
                var dx: i32 = -room;
                while (dx <= room) : (dx += 1) {
                    const landing = Point{ .x = x + dx, .y = y + dy };
                    if (!race.track.contains(landing)) continue;
                    const from_cheat = race.to_end[race.track.index(landing)];
                    if (from_cheat < 0) continue;
                    const jump: i64 = @intCast(@abs(dx) + @abs(dy));
                    if (race.uncheated_time - (to_cheat + jump + from_cheat) >= min_saving) cheats += 1;
                }
            }
        }
    }
    return cheats;
}

fn solveWith(allocator: std.mem.Allocator, input: []const u8, min_saving: i64) !Answers {
    const race = try prepare(allocator, input);
    defer race.deinit(allocator);

    return Answers{
        .part1 = countCheats(race, 2, min_saving),
        .part2 = countCheats(race, 20, min_saving),
    };
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    return solveWith(allocator, input, real_threshold);
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

const example = @embedFile("examples/day20.txt");

test "official example with two picosecond cheats" {
    const cases = [_]struct { min_saving: i64, cheats: i64 }{
        .{ .min_saving = 64, .cheats = 1 },
        .{ .min_saving = 40, .cheats = 2 },
        .{ .min_saving = 20, .cheats = 5 },
        .{ .min_saving = 2, .cheats = 44 },
    };
    for (cases) |case| {
        const answers = try solveWith(std.testing.allocator, example, case.min_saving);
        try std.testing.expectEqual(case.cheats, answers.part1);
    }
}

test "official example with twenty picosecond cheats" {
    const cases = [_]struct { min_saving: i64, cheats: i64 }{
        .{ .min_saving = 76, .cheats = 3 },
        .{ .min_saving = 74, .cheats = 7 },
        .{ .min_saving = 72, .cheats = 29 },
        .{ .min_saving = 50, .cheats = 285 },
    };
    for (cases) |case| {
        const answers = try solveWith(std.testing.allocator, example, case.min_saving);
        try std.testing.expectEqual(case.cheats, answers.part2);
    }
}

test "the best example cheat saves 76, so the real threshold counts nothing" {
    const answers = try solve(std.testing.allocator, example);
    try std.testing.expectEqual(@as(i64, 0), answers.part1);
    try std.testing.expectEqual(@as(i64, 0), answers.part2);
}

test "u shaped track" {
    // Cheating through the middle wall column saves 4 picoseconds from the start
    // and 2 from the cell below it.
    const corridor = "#####\n#S#E#\n#.#.#\n#...#\n#####\n";
    const cheats = try solveWith(std.testing.allocator, corridor, 2);
    try std.testing.expectEqual(@as(i64, 2), cheats.part1);
    try std.testing.expectEqual(@as(i64, 4), cheats.part2);

    const best = try solveWith(std.testing.allocator, corridor, 4);
    try std.testing.expectEqual(@as(i64, 1), best.part1);
    try std.testing.expectEqual(@as(i64, 1), best.part2);
}

test "malformed input returns an error" {
    const allocator = std.testing.allocator;
    try std.testing.expectError(error.EmptyMap, solve(allocator, ""));
    try std.testing.expectError(error.MissingStart, solve(allocator, "###\n#.#\n###\n"));
    try std.testing.expectError(error.MissingEnd, solve(allocator, "###\n#S#\n###\n"));
    try std.testing.expectError(error.DuplicateStart, solve(allocator, "#####\n#SSE#\n#####\n"));
    try std.testing.expectError(error.RaggedMap, solve(allocator, "#####\n#S.E#\n####\n"));
    try std.testing.expectError(error.InvalidMapCharacter, solve(allocator, "#####\n#S?E#\n#####\n"));
    try std.testing.expectError(error.EndUnreachable, solve(allocator, "#####\n#S#E#\n#####\n"));
}
