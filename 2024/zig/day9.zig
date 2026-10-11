const std = @import("std");

const Answers = struct {
    part1: i64,
    part2: i64,
};

const free_id: i64 = -1;

const Segment = struct {
    id: i64,
    start: i64,
    length: i64,
};

fn parse(allocator: std.mem.Allocator, input: []const u8) ![]Segment {
    var digits: std.ArrayList(u8) = .empty;
    defer digits.deinit(allocator);
    for (input) |c| {
        if (c == '\n' or c == '\r' or c == ' ') continue;
        if (c < '0' or c > '9') return error.InvalidCharacter;
        try digits.append(allocator, c);
    }
    if (digits.items.len == 0) return error.EmptyDiskMap;

    var segments: std.ArrayList(Segment) = .empty;
    errdefer segments.deinit(allocator);

    var position: i64 = 0;
    var file_id: i64 = 0;
    for (digits.items, 0..) |digit, i| {
        const length: i64 = @intCast(digit - '0');
        const is_file = i % 2 == 0;
        if (length > 0) {
            try segments.append(allocator, .{ .id = if (is_file) file_id else free_id, .start = position, .length = length });
            position += length;
        }
        if (is_file) file_id += 1;
    }
    return segments.toOwnedSlice(allocator);
}

fn totalLength(segments: []const Segment) i64 {
    var total: i64 = 0;
    for (segments) |s| total += s.length;
    return total;
}

fn solvePart1(allocator: std.mem.Allocator, segments: []const Segment) !i64 {
    const blocks = try allocator.alloc(i64, @intCast(totalLength(segments)));
    defer allocator.free(blocks);
    for (blocks) |*b| b.* = free_id;
    for (segments) |s| {
        var i: i64 = 0;
        while (i < s.length) : (i += 1) blocks[@intCast(s.start + i)] = s.id;
    }

    var left: usize = 0;
    var right: usize = blocks.len;
    while (true) {
        while (left < blocks.len and blocks[left] != free_id) left += 1;
        while (right > 0 and blocks[right - 1] == free_id) right -= 1;
        if (left >= right) break;
        blocks[left] = blocks[right - 1];
        blocks[right - 1] = free_id;
        left += 1;
        right -= 1;
    }

    var sum: i64 = 0;
    for (blocks, 0..) |block, idx| {
        if (block != free_id) sum += @as(i64, @intCast(idx)) * block;
    }
    return sum;
}

fn solvePart2(allocator: std.mem.Allocator, segments: []const Segment) !i64 {
    var files: std.ArrayList(Segment) = .empty;
    defer files.deinit(allocator);
    var free_spans: std.ArrayList(Segment) = .empty;
    defer free_spans.deinit(allocator);

    for (segments) |s| {
        if (s.id == free_id) {
            try free_spans.append(allocator, s);
        } else {
            try files.append(allocator, s);
        }
    }

    var sum: i64 = 0;
    var i = files.items.len;
    while (i > 0) {
        i -= 1;
        const file = files.items[i];
        var position = file.start;
        for (free_spans.items) |*free| {
            if (free.start >= file.start) break;
            if (free.length >= file.length) {
                position = free.start;
                free.start += file.length;
                free.length -= file.length;
                break;
            }
        }
        var k: i64 = 0;
        while (k < file.length) : (k += 1) sum += (position + k) * file.id;
    }
    return sum;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const segments = try parse(allocator, input);
    defer allocator.free(segments);
    const part1 = try solvePart1(allocator, segments);
    const part2 = try solvePart2(allocator, segments);
    return Answers{ .part1 = part1, .part2 = part2 };
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
    const input = @embedFile("examples/day9.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(i64, 1928), answers.part1);
    try std.testing.expectEqual(@as(i64, 2858), answers.part2);
}

test "single file with no free space" {
    const answers = try solve(std.testing.allocator, "5");
    try std.testing.expectEqual(@as(i64, 0), answers.part1);
    try std.testing.expectEqual(@as(i64, 0), answers.part2);
}

test "malformed input returns an error" {
    try std.testing.expectError(error.EmptyDiskMap, solve(std.testing.allocator, ""));
    try std.testing.expectError(error.InvalidCharacter, solve(std.testing.allocator, "12a45"));
}
