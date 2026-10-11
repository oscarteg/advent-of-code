const std = @import("std");

const Stack = std.ArrayList(u8);

const Move = struct {
    count: usize,
    from: usize,
    to: usize,
};

const Drawing = struct {
    stacks: []Stack,
    moves: []Move,

    fn deinit(self: *Drawing, allocator: std.mem.Allocator) void {
        for (self.stacks) |*stack| stack.deinit(allocator);
        allocator.free(self.stacks);
        allocator.free(self.moves);
    }
};

const Answers = struct {
    part1: []const u8,
    part2: []const u8,
};

fn cloneStacks(allocator: std.mem.Allocator, stacks: []Stack) ![]Stack {
    const copy = try allocator.alloc(Stack, stacks.len);
    for (stacks, 0..) |stack, i| {
        copy[i] = try stack.clone(allocator);
    }
    return copy;
}

fn freeStacks(allocator: std.mem.Allocator, stacks: []Stack) void {
    for (stacks) |*stack| stack.deinit(allocator);
    allocator.free(stacks);
}

// Crate letters sit at fixed columns (index 1, 5, 9, ... every 4th column),
// so stacks are read by character position rather than by splitting on
// whitespace.
fn parseInput(allocator: std.mem.Allocator, input: []const u8) !Drawing {
    const separator = std.mem.indexOf(u8, input, "\n\n") orelse return error.MissingBlankLine;
    const stacks_input = input[0..separator];
    const moves_input = input[separator + 2 ..];

    var line_iter = std.mem.splitScalar(u8, stacks_input, '\n');
    var lines: std.ArrayList([]const u8) = .empty;
    defer lines.deinit(allocator);
    while (line_iter.next()) |line| {
        try lines.append(allocator, line);
    }
    if (lines.items.len == 0) return error.MissingStackLabels;

    const labels = lines.pop() orelse return error.MissingStackLabels;
    var label_iter = std.mem.tokenizeScalar(u8, labels, ' ');
    var num_stacks: usize = 0;
    while (label_iter.next()) |_| num_stacks += 1;
    if (num_stacks == 0) return error.MissingStackLabels;

    const stacks = try allocator.alloc(Stack, num_stacks);
    for (stacks) |*stack| stack.* = .empty;
    errdefer freeStacks(allocator, stacks);

    var i = lines.items.len;
    while (i > 0) {
        i -= 1;
        const line = lines.items[i];
        for (stacks, 0..) |*stack, index| {
            const column = 1 + index * 4;
            if (column < line.len and line[column] != ' ') {
                try stack.append(allocator, line[column]);
            }
        }
    }

    var moves: std.ArrayList(Move) = .empty;
    defer moves.deinit(allocator);
    var move_lines = std.mem.splitScalar(u8, moves_input, '\n');
    while (move_lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \r");
        if (trimmed.len == 0) continue;
        try moves.append(allocator, try parseMove(trimmed));
    }

    return Drawing{
        .stacks = stacks,
        .moves = try moves.toOwnedSlice(allocator),
    };
}

fn parseMove(line: []const u8) !Move {
    var parts: [6][]const u8 = undefined;
    var part_iter = std.mem.tokenizeScalar(u8, line, ' ');
    var count: usize = 0;
    while (part_iter.next()) |part| {
        if (count >= parts.len) return error.MalformedMove;
        parts[count] = part;
        count += 1;
    }
    if (count != 6 or !std.mem.eql(u8, parts[0], "move") or
        !std.mem.eql(u8, parts[2], "from") or !std.mem.eql(u8, parts[4], "to"))
    {
        return error.MalformedMove;
    }

    const move_count = std.fmt.parseInt(usize, parts[1], 10) catch return error.MalformedMove;
    const from = std.fmt.parseInt(usize, parts[3], 10) catch return error.MalformedMove;
    const to = std.fmt.parseInt(usize, parts[5], 10) catch return error.MalformedMove;
    if (from < 1 or to < 1) return error.MalformedMove;

    return Move{ .count = move_count, .from = from - 1, .to = to - 1 };
}

fn applyMoves(allocator: std.mem.Allocator, stacks: []Stack, moves: []const Move, preserve_order: bool) !void {
    for (moves) |mv| {
        if (mv.from >= stacks.len or mv.to >= stacks.len) return error.UnknownStack;
        const from_len = stacks[mv.from].items.len;
        if (mv.count > from_len) return error.NotEnoughCrates;

        const split_at = from_len - mv.count;
        const moved = try allocator.alloc(u8, mv.count);
        defer allocator.free(moved);
        @memcpy(moved, stacks[mv.from].items[split_at..]);
        stacks[mv.from].shrinkRetainingCapacity(split_at);

        if (!preserve_order) std.mem.reverse(u8, moved);
        try stacks[mv.to].appendSlice(allocator, moved);
    }
}

fn topCrates(allocator: std.mem.Allocator, stacks: []const Stack) ![]const u8 {
    const result = try allocator.alloc(u8, stacks.len);
    for (stacks, 0..) |stack, i| {
        result[i] = if (stack.items.len > 0) stack.items[stack.items.len - 1] else ' ';
    }
    return result;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var drawing = try parseInput(allocator, input);
    defer drawing.deinit(allocator);

    const part1_stacks = try cloneStacks(allocator, drawing.stacks);
    defer freeStacks(allocator, part1_stacks);
    try applyMoves(allocator, part1_stacks, drawing.moves, false);
    const part1 = try topCrates(allocator, part1_stacks);

    try applyMoves(allocator, drawing.stacks, drawing.moves, true);
    const part2 = try topCrates(allocator, drawing.stacks);

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
    try stdout.print("Part 1: {s}\nPart 2: {s}\n", .{ answers.part1, answers.part2 });
    try stdout.flush();
}

const example_input = @embedFile("examples/day_5.txt");

test "part 1 and part 2 example" {
    const allocator = std.testing.allocator;
    const answers = try solve(allocator, example_input);
    defer allocator.free(answers.part1);
    defer allocator.free(answers.part2);
    try std.testing.expectEqualStrings("CMZ", answers.part1);
    try std.testing.expectEqualStrings("MCD", answers.part2);
}

test "single crate single move" {
    const allocator = std.testing.allocator;
    const input = "[A]\n 1 \n\nmove 1 from 1 to 1";
    const answers = try solve(allocator, input);
    defer allocator.free(answers.part1);
    defer allocator.free(answers.part2);
    try std.testing.expectEqualStrings("A", answers.part1);
    try std.testing.expectEqualStrings("A", answers.part2);
}

test "malformed input missing blank line returns error" {
    const allocator = std.testing.allocator;
    const input = "[A]\n 1 \nmove 1 from 1 to 1";
    try std.testing.expectError(error.MissingBlankLine, parseInput(allocator, input));
}

test "malformed move line returns error" {
    const allocator = std.testing.allocator;
    const input = "[A]\n 1 \n\nmove one from 1 to 1";
    try std.testing.expectError(error.MalformedMove, parseInput(allocator, input));
}

test "moving more crates than available returns error" {
    const allocator = std.testing.allocator;
    var drawing = try parseInput(allocator, "[A]\n 1 \n\nmove 5 from 1 to 1");
    defer drawing.deinit(allocator);
    try std.testing.expectError(error.NotEnoughCrates, applyMoves(allocator, drawing.stacks, drawing.moves, false));
}
