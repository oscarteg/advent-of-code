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

const source = Point{ .x = 500, .y = 0 };

fn parsePoint(text: []const u8) !Point {
    const comma = std.mem.indexOfScalar(u8, text, ',') orelse return error.MissingComma;
    const x = try std.fmt.parseInt(i32, std.mem.trim(u8, text[0..comma], " \t"), 10);
    const y = try std.fmt.parseInt(i32, std.mem.trim(u8, text[comma + 1 ..], " \t"), 10);
    return .{ .x = x, .y = y };
}

fn addSegment(rocks: *std.AutoHashMap(Point, void), a: Point, b: Point) !void {
    if (a.x == b.x) {
        const lo = @min(a.y, b.y);
        const hi = @max(a.y, b.y);
        var y = lo;
        while (y <= hi) : (y += 1) try rocks.put(.{ .x = a.x, .y = y }, {});
    } else if (a.y == b.y) {
        const lo = @min(a.x, b.x);
        const hi = @max(a.x, b.x);
        var x = lo;
        while (x <= hi) : (x += 1) try rocks.put(.{ .x = x, .y = a.y }, {});
    } else {
        return error.NotAxisAligned;
    }
}

fn parseRocks(allocator: std.mem.Allocator, input: []const u8) !std.AutoHashMap(Point, void) {
    var rocks = std.AutoHashMap(Point, void).init(allocator);
    errdefer rocks.deinit();

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw_line| {
        const line = std.mem.trim(u8, raw_line, " \t\r");
        if (line.len == 0) continue;
        var points = std.mem.tokenizeSequence(u8, line, "->");
        var previous = try parsePoint(points.next() orelse return error.MalformedLine);
        var count: usize = 1;
        while (points.next()) |raw_point| {
            const point = try parsePoint(raw_point);
            try addSegment(&rocks, previous, point);
            previous = point;
            count += 1;
        }
        if (count < 2) return error.MalformedLine;
    }

    if (rocks.count() == 0) return error.EmptyInput;
    return rocks;
}

fn isBlocked(settled: *const std.AutoHashMap(Point, void), floor_y: ?i32, point: Point) bool {
    if (floor_y) |fy| if (point.y == fy) return true;
    return settled.contains(point);
}

/// Drops sand from the source until it either falls past `max_y` (no floor)
/// or comes to rest on the source itself (with a floor), returning the count
/// of settled units.
fn simulate(rocks: *const std.AutoHashMap(Point, void), max_y: i32, floor_y: ?i32) !u64 {
    var settled = try rocks.clone();
    defer settled.deinit();

    var count: u64 = 0;
    while (!settled.contains(source)) {
        var sand = source;
        while (true) {
            if (floor_y == null and sand.y > max_y) return count;
            const below = Point{ .x = sand.x, .y = sand.y + 1 };
            const below_left = Point{ .x = sand.x - 1, .y = sand.y + 1 };
            const below_right = Point{ .x = sand.x + 1, .y = sand.y + 1 };
            if (!isBlocked(&settled, floor_y, below)) {
                sand = below;
            } else if (!isBlocked(&settled, floor_y, below_left)) {
                sand = below_left;
            } else if (!isBlocked(&settled, floor_y, below_right)) {
                sand = below_right;
            } else {
                break;
            }
        }
        try settled.put(sand, {});
        count += 1;
    }
    return count;
}

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var rocks = try parseRocks(allocator, input);
    defer rocks.deinit();

    var max_y: i32 = std.math.minInt(i32);
    var it = rocks.keyIterator();
    while (it.next()) |point| {
        if (point.y > max_y) max_y = point.y;
    }

    return .{
        .part1 = try simulate(&rocks, max_y, null),
        .part2 = try simulate(&rocks, max_y, max_y + 2),
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
    const input = @embedFile("examples/day_14.txt");
    const answers = try solve(testing.allocator, input);
    try testing.expectEqual(@as(u64, 24), answers.part1);
    try testing.expectEqual(@as(u64, 93), answers.part2);
}

test "flat ledge catches a bounded pile" {
    const answers = try solve(testing.allocator, "498,3 -> 502,3\n");
    try testing.expectEqual(@as(u64, 4), answers.part1);
}

test "empty input returns error" {
    try testing.expectError(error.EmptyInput, solve(testing.allocator, ""));
}

test "missing comma returns error" {
    try testing.expectError(error.MissingComma, solve(testing.allocator, "498 4 -> 498 6\n"));
}

test "diagonal segment returns error" {
    try testing.expectError(error.NotAxisAligned, solve(testing.allocator, "498,4 -> 500,6\n"));
}
