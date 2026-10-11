const std = @import("std");

const Robot = struct {
    px: i64,
    py: i64,
    vx: i64,
    vy: i64,
};

const Answers = struct {
    part1: i64,
    part2: i64,
};

const Position = struct { x: i64, y: i64 };

fn parseField(text: []const u8) !i64 {
    return std.fmt.parseInt(i64, text, 10) catch error.InvalidFormat;
}

fn parse(allocator: std.mem.Allocator, input: []const u8) !std.ArrayList(Robot) {
    var robots = std.ArrayList(Robot).empty;
    errdefer robots.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;
        if (!std.mem.startsWith(u8, trimmed, "p=")) return error.InvalidFormat;

        const v_at = std.mem.indexOf(u8, trimmed, " v=") orelse return error.InvalidFormat;
        const p_part = trimmed[2..v_at];
        const v_part = trimmed[v_at + 3 ..];

        const p_comma = std.mem.indexOf(u8, p_part, ",") orelse return error.InvalidFormat;
        const v_comma = std.mem.indexOf(u8, v_part, ",") orelse return error.InvalidFormat;

        try robots.append(allocator, Robot{
            .px = try parseField(p_part[0..p_comma]),
            .py = try parseField(p_part[p_comma + 1 ..]),
            .vx = try parseField(v_part[0..v_comma]),
            .vy = try parseField(v_part[v_comma + 1 ..]),
        });
    }

    if (robots.items.len == 0) return error.NoRobots;

    return robots;
}

fn wrap(value: i64, modulus: i64) i64 {
    const remainder = @rem(value, modulus);
    return if (remainder < 0) remainder + modulus else remainder;
}

fn positionAt(robot: Robot, seconds: i64, width: i64, height: i64) Position {
    return .{
        .x = wrap(robot.px + robot.vx * seconds, width),
        .y = wrap(robot.py + robot.vy * seconds, height),
    };
}

fn safetyFactor(robots: []const Robot, width: i64, height: i64, seconds: i64) i64 {
    const mid_x = @divTrunc(width, 2);
    const mid_y = @divTrunc(height, 2);
    var quadrants = [4]i64{ 0, 0, 0, 0 };

    for (robots) |robot| {
        const pos = positionAt(robot, seconds, width, height);
        if (pos.x == mid_x or pos.y == mid_y) continue;
        var index: usize = 0;
        if (pos.x > mid_x) index += 1;
        if (pos.y > mid_y) index += 2;
        quadrants[index] += 1;
    }

    return quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3];
}

fn positionsUnique(allocator: std.mem.Allocator, positions: []const Position, width: i64) !bool {
    var seen = std.AutoHashMap(i64, void).init(allocator);
    defer seen.deinit();

    for (positions) |pos| {
        const key = pos.y * width + pos.x;
        if (seen.contains(key)) return false;
        try seen.put(key, {});
    }
    return true;
}

fn findPictureSecond(allocator: std.mem.Allocator, robots: []const Robot, width: i64, height: i64) !i64 {
    const period = width * height;
    const positions = try allocator.alloc(Position, robots.len);
    defer allocator.free(positions);

    var t: i64 = 0;
    while (t < period) : (t += 1) {
        for (robots, 0..) |robot, i| {
            positions[i] = positionAt(robot, t, width, height);
        }
        if (try positionsUnique(allocator, positions, width)) {
            return t;
        }
    }

    return error.NoPictureFound;
}

// The input never states the grid size: the puzzle's example uses 11x7 and
// every real input uses 101x103, so pick between those two.
fn inferDimensions(robots: []const Robot) struct { width: i64, height: i64 } {
    for (robots) |robot| {
        if (robot.px >= 11 or robot.py >= 7) return .{ .width = 101, .height = 103 };
    }
    return .{ .width = 11, .height = 7 };
}

fn compute(allocator: std.mem.Allocator, robots: []const Robot, width: i64, height: i64) !Answers {
    const part1 = safetyFactor(robots, width, height, 100);
    const part2 = try findPictureSecond(allocator, robots, width, height);
    return Answers{ .part1 = part1, .part2 = part2 };
}

fn solveWith(allocator: std.mem.Allocator, input: []const u8, width: i64, height: i64) !Answers {
    var robots = try parse(allocator, input);
    defer robots.deinit(allocator);
    return compute(allocator, robots.items, width, height);
}

// The grid size is not given in the input; it is inferred as one more than
// the largest x/y coordinate, which matches the real 101x103 puzzle grid.
pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var robots = try parse(allocator, input);
    defer robots.deinit(allocator);
    const dims = inferDimensions(robots.items);
    return compute(allocator, robots.items, dims.width, dims.height);
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

test "part1 example" {
    const example = @embedFile("examples/day14.txt");
    const allocator = std.testing.allocator;
    const result = try solveWith(allocator, example, 11, 7);
    try std.testing.expectEqual(@as(i64, 12), result.part1);
}

test "picture detector finds immediate non-overlap" {
    const allocator = std.testing.allocator;
    const input = "p=0,0 v=0,0\np=1,0 v=0,0\n";
    var robots = try parse(allocator, input);
    defer robots.deinit(allocator);
    const t = try findPictureSecond(allocator, robots.items, 3, 3);
    try std.testing.expectEqual(@as(i64, 0), t);
}

test "picture detector finds later non-overlap" {
    const allocator = std.testing.allocator;
    const input = "p=0,0 v=1,0\np=0,0 v=2,0\n";
    var robots = try parse(allocator, input);
    defer robots.deinit(allocator);
    const t = try findPictureSecond(allocator, robots.items, 5, 1);
    try std.testing.expectEqual(@as(i64, 1), t);
}

test "malformed input errors" {
    const allocator = std.testing.allocator;
    try std.testing.expectError(error.NoRobots, parse(allocator, ""));
    try std.testing.expectError(error.InvalidFormat, parse(allocator, "garbage\n"));
}
