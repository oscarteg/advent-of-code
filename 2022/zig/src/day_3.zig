const std = @import("std");

/// Returns the priority of an item: a-z are 1-26, A-Z are 27-52.
fn priority(c: u8) !u32 {
    if (c >= 'a' and c <= 'z') return c - 'a' + 1;
    if (c >= 'A' and c <= 'Z') return c - 'A' + 27;
    return error.InvalidItem;
}

/// Builds a 52-bit mask with one bit set per distinct item type present in `s`.
fn itemMask(s: []const u8) !u64 {
    var mask: u64 = 0;
    for (s) |c| {
        const p = try priority(c);
        mask |= @as(u64, 1) << @intCast(p - 1);
    }
    return mask;
}

/// Finds the single item priority common to every given string.
fn commonPriority(groups: []const []const u8) !u32 {
    if (groups.len == 0) return error.NoGroups;

    var common = try itemMask(groups[0]);
    for (groups[1..]) |group| {
        common &= try itemMask(group);
    }

    if (common == 0) return error.NoCommonItem;
    if (@popCount(common) != 1) return error.MultipleCommonItems;
    return @ctz(common) + 1;
}

fn part1(input: []const u8) !u32 {
    var sum: u32 = 0;
    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \r\t");
        if (trimmed.len == 0) continue;
        if (trimmed.len % 2 != 0) return error.OddLengthRucksack;
        const half = trimmed.len / 2;
        sum += try commonPriority(&.{ trimmed[0..half], trimmed[half..] });
    }
    return sum;
}

fn part2(input: []const u8) !u32 {
    var lines_list: std.ArrayList([]const u8) = .empty;
    defer lines_list.deinit(std.heap.page_allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \r\t");
        if (trimmed.len == 0) continue;
        try lines_list.append(std.heap.page_allocator, trimmed);
    }

    if (lines_list.items.len % 3 != 0 or lines_list.items.len == 0) {
        return error.NotMultipleOfThree;
    }

    var sum: u32 = 0;
    var i: usize = 0;
    while (i < lines_list.items.len) : (i += 3) {
        sum += try commonPriority(lines_list.items[i .. i + 3]);
    }
    return sum;
}

const Answers = struct {
    part1: u32,
    part2: u32,
};

fn solve(arena: std.mem.Allocator, input: []const u8) !Answers {
    _ = arena;
    return Answers{
        .part1 = try part1(input),
        .part2 = try part2(input),
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

test "part1 example" {
    const input = @embedFile("examples/day_3.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u32, 157), answers.part1);
}

test "part2 example" {
    const input = @embedFile("examples/day_3.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u32, 70), answers.part2);
}

test "part1 rejects odd length rucksack" {
    try std.testing.expectError(error.OddLengthRucksack, part1("abc\n"));
}

test "part2 rejects non multiple of three" {
    try std.testing.expectError(error.NotMultipleOfThree, part2("abc\ndef\n"));
}

test "priority rejects non letter" {
    try std.testing.expectError(error.InvalidItem, priority('1'));
}
