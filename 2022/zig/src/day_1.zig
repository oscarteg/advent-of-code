const std = @import("std");
const testing = std.testing;

const Answers = struct {
    part1: u32,
    part2: u32,
};

fn elfTotals(allocator: std.mem.Allocator, input: []const u8) ![]u32 {
    var totals: std.ArrayList(u32) = .empty;
    errdefer totals.deinit(allocator);

    const trimmed = std.mem.trim(u8, input, "\n");
    if (trimmed.len == 0) return error.EmptyInput;

    var elves = std.mem.splitSequence(u8, trimmed, "\n\n");
    while (elves.next()) |elf| {
        var total: u32 = 0;
        var foods = std.mem.tokenizeScalar(u8, elf, '\n');
        while (foods.next()) |food| {
            total += try std.fmt.parseInt(u32, std.mem.trim(u8, food, " \t\r"), 10);
        }
        try totals.append(allocator, total);
    }

    if (totals.items.len == 0) return error.EmptyInput;
    return totals.toOwnedSlice(allocator);
}

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const totals = try elfTotals(allocator, input);
    defer allocator.free(totals);

    std.mem.sort(u32, totals, {}, std.sort.desc(u32));

    var part2: u32 = 0;
    for (totals[0..@min(3, totals.len)]) |total| part2 += total;

    return .{ .part1 = totals[0], .part2 = part2 };
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
    const input = @embedFile("examples/day_1.txt");
    const answers = try solve(testing.allocator, input);
    try testing.expectEqual(@as(u32, 24000), answers.part1);
    try testing.expectEqual(@as(u32, 45000), answers.part2);
}

test "single elf" {
    const answers = try solve(testing.allocator, "100\n200");
    try testing.expectEqual(@as(u32, 300), answers.part1);
    try testing.expectEqual(@as(u32, 300), answers.part2);
}

test "malformed calorie value returns error" {
    try testing.expectError(error.InvalidCharacter, solve(testing.allocator, "1000\nnot-a-number\n\n2000"));
}

test "empty input returns error" {
    try testing.expectError(error.EmptyInput, solve(testing.allocator, ""));
}
