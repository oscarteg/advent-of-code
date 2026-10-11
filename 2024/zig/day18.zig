const std = @import("std");

const Point = struct { x: i32, y: i32 };

const Answers = struct {
    part1: i64,
    part2: []const u8,
};

fn parse(allocator: std.mem.Allocator, input: []const u8) ![]Point {
    var points: std.ArrayList(Point) = .empty;
    errdefer points.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw| {
        const line = std.mem.trim(u8, raw, " \r");
        if (line.len == 0) continue;
        const comma = std.mem.indexOfScalar(u8, line, ',') orelse return error.InvalidCoordinate;
        const x = try std.fmt.parseInt(i32, line[0..comma], 10);
        const y = try std.fmt.parseInt(i32, line[comma + 1 ..], 10);
        try points.append(allocator, Point{ .x = x, .y = y });
    }
    if (points.items.len == 0) return error.InvalidCoordinate;

    return points.toOwnedSlice(allocator);
}

fn shortestPath(allocator: std.mem.Allocator, corrupted: []const Point, size: i32) !i64 {
    const cells: usize = @intCast(size * size);
    const blocked = try allocator.alloc(bool, cells);
    defer allocator.free(blocked);
    @memset(blocked, false);
    for (corrupted) |p| {
        if (p.x >= 0 and p.x < size and p.y >= 0 and p.y < size) {
            blocked[@intCast(p.y * size + p.x)] = true;
        }
    }

    const distance = try allocator.alloc(i64, cells);
    defer allocator.free(distance);
    @memset(distance, -1);

    var queue: std.ArrayList(Point) = .empty;
    defer queue.deinit(allocator);
    distance[0] = 0;
    try queue.append(allocator, Point{ .x = 0, .y = 0 });

    const deltas = [4]Point{ .{ .x = 1, .y = 0 }, .{ .x = -1, .y = 0 }, .{ .x = 0, .y = 1 }, .{ .x = 0, .y = -1 } };

    var head: usize = 0;
    while (head < queue.items.len) {
        const current = queue.items[head];
        head += 1;
        if (current.x == size - 1 and current.y == size - 1) {
            return distance[@intCast(current.y * size + current.x)];
        }
        for (deltas) |delta| {
            const nx = current.x + delta.x;
            const ny = current.y + delta.y;
            if (nx < 0 or nx >= size or ny < 0 or ny >= size) continue;
            const index: usize = @intCast(ny * size + nx);
            if (blocked[index] or distance[index] != -1) continue;
            distance[index] = distance[@intCast(current.y * size + current.x)] + 1;
            try queue.append(allocator, Point{ .x = nx, .y = ny });
        }
    }
    return -1;
}

// Binary search over the byte count for the first one that blocks the exit:
// the exit is reachable for every prefix up to the blocker and unreachable after.
fn firstBlockingByte(allocator: std.mem.Allocator, points: []const Point, size: i32) !Point {
    var low: usize = 0;
    var high: usize = points.len;
    while (low < high) {
        const mid = low + (high - low) / 2;
        const reachable = try shortestPath(allocator, points[0 .. mid + 1], size);
        if (reachable == -1) {
            high = mid;
        } else {
            low = mid + 1;
        }
    }
    if (low >= points.len) return error.ExitNeverBlocked;
    return points[low];
}

fn solve(allocator: std.mem.Allocator, input: []const u8, size: i32, fallen_bytes: usize) !Answers {
    const points = try parse(allocator, input);
    defer allocator.free(points);
    if (fallen_bytes > points.len) return error.NotEnoughBytes;

    const part1 = try shortestPath(allocator, points[0..fallen_bytes], size);
    if (part1 == -1) return error.ExitUnreachable;

    const blocker = try firstBlockingByte(allocator, points, size);
    const part2 = try std.fmt.allocPrint(allocator, "{d},{d}", .{ blocker.x, blocker.y });

    return Answers{ .part1 = part1, .part2 = part2 };
}

pub fn main(init: std.process.Init) !void {
    const arena = init.arena.allocator();
    var stdin_buffer: [4096]u8 = undefined;
    var stdin_reader = std.Io.File.stdin().reader(init.io, &stdin_buffer);
    const input = try stdin_reader.interface.allocRemaining(arena, .unlimited);
    const answers = try solve(arena, input, 7, 12);
    var stdout_buffer: [256]u8 = undefined;
    var stdout_writer = std.Io.File.stdout().writer(init.io, &stdout_buffer);
    const stdout = &stdout_writer.interface;
    try stdout.print("Part 1: {d}\nPart 2: {s}\n", .{ answers.part1, answers.part2 });
    try stdout.flush();
}

test "official example" {
    const input = @embedFile("examples/day18.txt");
    const answers = try solve(std.testing.allocator, input, 7, 12);
    defer std.testing.allocator.free(answers.part2);
    try std.testing.expectEqual(@as(i64, 22), answers.part1);
    try std.testing.expectEqualStrings("6,1", answers.part2);
}

test "two byte grid blocks after the second byte" {
    const answers = try solve(std.testing.allocator, "1,0\n0,1\n", 2, 0);
    defer std.testing.allocator.free(answers.part2);
    try std.testing.expectEqual(@as(i64, 2), answers.part1);
    try std.testing.expectEqualStrings("0,1", answers.part2);
}

test "malformed input returns an error" {
    try std.testing.expectError(error.InvalidCoordinate, solve(std.testing.allocator, "", 7, 0));
    try std.testing.expectError(error.InvalidCoordinate, solve(std.testing.allocator, "not a coordinate\n", 7, 1));
    try std.testing.expectError(error.NotEnoughBytes, solve(std.testing.allocator, "0,0\n", 7, 5));
}
