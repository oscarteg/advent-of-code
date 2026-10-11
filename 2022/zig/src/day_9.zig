const std = @import("std");
const testing = std.testing;

const Answers = struct {
    part1: u64,
    part2: u64,
};

const Point = struct {
    x: i32,
    y: i32,
};

const Motion = struct {
    dx: i32,
    dy: i32,
    steps: u32,
};

fn parseMotions(allocator: std.mem.Allocator, input: []const u8) ![]Motion {
    var motions: std.ArrayList(Motion) = .empty;
    errdefer motions.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw_line| {
        const line = std.mem.trim(u8, raw_line, " \t\r");
        if (line.len == 0) continue;
        var fields = std.mem.tokenizeScalar(u8, line, ' ');
        const direction = fields.next() orelse return error.MalformedLine;
        const steps_str = fields.next() orelse return error.MalformedLine;
        if (fields.next() != null) return error.MalformedLine;
        const steps = try std.fmt.parseInt(u32, steps_str, 10);
        const delta: struct { dx: i32, dy: i32 } = switch (direction[0]) {
            'R' => .{ .dx = 1, .dy = 0 },
            'L' => .{ .dx = -1, .dy = 0 },
            'U' => .{ .dx = 0, .dy = 1 },
            'D' => .{ .dx = 0, .dy = -1 },
            else => return error.UnknownDirection,
        };
        if (direction.len != 1) return error.UnknownDirection;
        try motions.append(allocator, .{ .dx = delta.dx, .dy = delta.dy, .steps = steps });
    }

    if (motions.items.len == 0) return error.EmptyInput;
    return motions.toOwnedSlice(allocator);
}

fn signumI32(n: i32) i32 {
    if (n > 0) return 1;
    if (n < 0) return -1;
    return 0;
}

fn follow(head: Point, tail: Point) Point {
    const dx = head.x - tail.x;
    const dy = head.y - tail.y;
    if (@abs(dx) <= 1 and @abs(dy) <= 1) return tail;
    return .{ .x = tail.x + signumI32(dx), .y = tail.y + signumI32(dy) };
}

fn simulate(allocator: std.mem.Allocator, motions: []const Motion, rope_len: usize) !usize {
    const knots = try allocator.alloc(Point, rope_len);
    defer allocator.free(knots);
    for (knots) |*knot| knot.* = .{ .x = 0, .y = 0 };

    var visited = std.AutoHashMap(Point, void).init(allocator);
    defer visited.deinit();
    try visited.put(knots[rope_len - 1], {});

    for (motions) |motion| {
        var step: u32 = 0;
        while (step < motion.steps) : (step += 1) {
            knots[0].x += motion.dx;
            knots[0].y += motion.dy;
            var k: usize = 1;
            while (k < rope_len) : (k += 1) {
                knots[k] = follow(knots[k - 1], knots[k]);
            }
            try visited.put(knots[rope_len - 1], {});
        }
    }

    return visited.count();
}

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const motions = try parseMotions(allocator, input);
    defer allocator.free(motions);
    return .{
        .part1 = try simulate(allocator, motions, 2),
        .part2 = try simulate(allocator, motions, 10),
    };
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
    const input = @embedFile("examples/day_9.txt");
    const answers = try solve(testing.allocator, input);
    try testing.expectEqual(@as(u64, 13), answers.part1);
    try testing.expectEqual(@as(u64, 1), answers.part2);
}

test "larger example part 2" {
    const input = "R 5\nU 8\nL 8\nD 3\nR 17\nD 10\nL 25\nU 20\n";
    const answers = try solve(testing.allocator, input);
    try testing.expectEqual(@as(u64, 36), answers.part2);
}

test "single step does not move tail" {
    const answers = try solve(testing.allocator, "R 1\n");
    try testing.expectEqual(@as(u64, 1), answers.part1);
    try testing.expectEqual(@as(u64, 1), answers.part2);
}

test "empty input returns error" {
    try testing.expectError(error.EmptyInput, solve(testing.allocator, ""));
}

test "malformed direction returns error" {
    try testing.expectError(error.UnknownDirection, solve(testing.allocator, "X 4\n"));
}

test "malformed step count returns error" {
    try testing.expectError(error.InvalidCharacter, solve(testing.allocator, "R four\n"));
}
