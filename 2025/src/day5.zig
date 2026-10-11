const std = @import("std");

const Range = struct {
    start: u32,
    end: u32,
};

const Answers = struct {
    part1: u32,
    part2: u32,
};

fn countLines(text: []const u8) usize {
    var count: usize = 0;
    for (text) |c| {
        if (c == '\n') count += 1;
    }
    if (text.len > 0 and text[text.len - 1] != '\n') count += 1;
    return if (count > 0) count else 0;
}

fn parseInput(arena: std.mem.Allocator, input: []const u8) !struct {
    ranges: []Range,
    ids: []u32,
} {
    var parts = std.mem.splitSequence(u8, input, "\n\n");
    const ranges_text = parts.next() orelse return error.InvalidFormat;
    const ids_text = parts.next() orelse return error.InvalidFormat;

    const range_count = countLines(ranges_text);
    const id_count = countLines(ids_text);

    var ranges = try arena.alloc(Range, range_count);
    var ids = try arena.alloc(u32, id_count);

    var range_idx: usize = 0;
    var range_lines = std.mem.splitSequence(u8, ranges_text, "\n");
    while (range_lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;

        var dash_parts = std.mem.splitSequence(u8, trimmed, "-");
        const start_str = dash_parts.next() orelse return error.InvalidFormat;
        const end_str = dash_parts.next() orelse return error.InvalidFormat;

        const start = std.fmt.parseInt(u32, start_str, 10) catch return error.InvalidFormat;
        const end = std.fmt.parseInt(u32, end_str, 10) catch return error.InvalidFormat;

        ranges[range_idx] = Range{ .start = start, .end = end };
        range_idx += 1;
    }

    var id_idx: usize = 0;
    var id_lines = std.mem.splitSequence(u8, ids_text, "\n");
    while (id_lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;

        const id = std.fmt.parseInt(u32, trimmed, 10) catch return error.InvalidFormat;
        ids[id_idx] = id;
        id_idx += 1;
    }

    return .{
        .ranges = ranges[0..range_idx],
        .ids = ids[0..id_idx],
    };
}

fn isInRange(id: u32, ranges: []const Range) bool {
    for (ranges) |r| {
        if (id >= r.start and id <= r.end) {
            return true;
        }
    }
    return false;
}

fn solvePart1(ranges: []const Range, ids: []const u32) u32 {
    var count: u32 = 0;
    for (ids) |id| {
        if (isInRange(id, ranges)) {
            count += 1;
        }
    }
    return count;
}

fn solvePart2(arena: std.mem.Allocator, ranges: []const Range) !u32 {
    if (ranges.len == 0) return 0;

    var sorted_ranges = try arena.dupe(Range, ranges);
    std.mem.sort(Range, sorted_ranges, {}, struct {
        fn compare(_: void, a: Range, b: Range) bool {
            return a.start < b.start;
        }
    }.compare);

    var merged_ranges = try arena.alloc(Range, ranges.len);
    var merged_count: usize = 0;

    merged_ranges[0] = sorted_ranges[0];
    merged_count = 1;

    for (sorted_ranges[1..]) |r| {
        var last = &merged_ranges[merged_count - 1];

        if (r.start <= last.end + 1) {
            if (r.end > last.end) {
                last.end = r.end;
            }
        } else {
            merged_ranges[merged_count] = r;
            merged_count += 1;
        }
    }

    var total: u32 = 0;
    for (merged_ranges[0..merged_count]) |r| {
        total += r.end - r.start + 1;
    }

    return total;
}

pub fn solve(arena: std.mem.Allocator, input: []const u8) !Answers {
    const parsed = try parseInput(arena, input);
    const part1 = solvePart1(parsed.ranges, parsed.ids);
    const part2 = try solvePart2(arena, parsed.ranges);

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

test "example input" {
    const example = @embedFile("examples/day5.txt");
    var buffer: [4096]u8 = undefined;
    var fba = std.heap.FixedBufferAllocator.init(&buffer);
    const answers = try solve(fba.allocator(), example);
    try std.testing.expectEqual(answers.part1, 3);
    try std.testing.expectEqual(answers.part2, 14);
}

test "single range" {
    const input = "5-10\n\n7\n3\n12";
    var buffer: [4096]u8 = undefined;
    var fba = std.heap.FixedBufferAllocator.init(&buffer);
    const answers = try solve(fba.allocator(), input);
    try std.testing.expectEqual(answers.part1, 1);
    try std.testing.expectEqual(answers.part2, 6);
}

test "no ranges" {
    const input = "\n\n1\n2\n3";
    var buffer: [4096]u8 = undefined;
    var fba = std.heap.FixedBufferAllocator.init(&buffer);
    const answers = try solve(fba.allocator(), input);
    try std.testing.expectEqual(answers.part1, 0);
    try std.testing.expectEqual(answers.part2, 0);
}

test "malformed range" {
    const input = "5-\n\n1";
    var buffer: [4096]u8 = undefined;
    var fba = std.heap.FixedBufferAllocator.init(&buffer);
    try std.testing.expectError(error.InvalidFormat, solve(fba.allocator(), input));
}
