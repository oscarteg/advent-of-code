const std = @import("std");

const Answers = struct {
    part1: i64,
    part2: i64,
};

const Equation = struct {
    target: i64,
    operands: []i64,
};

fn parseEquation(allocator: std.mem.Allocator, line: []const u8) !Equation {
    const colon = std.mem.indexOfScalar(u8, line, ':') orelse return error.InvalidEquation;
    const target = try std.fmt.parseInt(i64, line[0..colon], 10);

    var operands: std.ArrayList(i64) = .empty;
    errdefer operands.deinit(allocator);

    var it = std.mem.tokenizeScalar(u8, line[colon + 1 ..], ' ');
    while (it.next()) |token| {
        const value = try std.fmt.parseInt(i64, token, 10);
        try operands.append(allocator, value);
    }
    if (operands.items.len == 0) return error.InvalidEquation;

    return Equation{ .target = target, .operands = try operands.toOwnedSlice(allocator) };
}

fn parse(allocator: std.mem.Allocator, input: []const u8) ![]Equation {
    var equations: std.ArrayList(Equation) = .empty;
    errdefer equations.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw| {
        const line = std.mem.trim(u8, raw, " \r");
        if (line.len == 0) continue;
        try equations.append(allocator, try parseEquation(allocator, line));
    }
    if (equations.items.len == 0) return error.InvalidEquation;

    return equations.toOwnedSlice(allocator);
}

fn concatenate(left: i64, right: i64) i64 {
    var multiplier: i64 = 1;
    while (multiplier <= right) multiplier *= 10;
    return left * multiplier + right;
}

fn canReach(target: i64, operands: []const i64, index: usize, accumulated: i64, allow_concat: bool) bool {
    if (accumulated > target) return false;
    if (index == operands.len) return accumulated == target;

    const next = operands[index];
    if (canReach(target, operands, index + 1, accumulated + next, allow_concat)) return true;
    if (canReach(target, operands, index + 1, accumulated * next, allow_concat)) return true;
    if (allow_concat and canReach(target, operands, index + 1, concatenate(accumulated, next), allow_concat)) {
        return true;
    }
    return false;
}

fn isSolvable(equation: Equation, allow_concat: bool) bool {
    return canReach(equation.target, equation.operands, 1, equation.operands[0], allow_concat);
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const equations = try parse(allocator, input);
    defer {
        for (equations) |equation| allocator.free(equation.operands);
        allocator.free(equations);
    }

    var part1: i64 = 0;
    var part2: i64 = 0;
    for (equations) |equation| {
        if (isSolvable(equation, false)) {
            part1 += equation.target;
            part2 += equation.target;
        } else if (isSolvable(equation, true)) {
            part2 += equation.target;
        }
    }
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
    const input = @embedFile("examples/day7.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(i64, 3749), answers.part1);
    try std.testing.expectEqual(@as(i64, 11387), answers.part2);
}

test "single operand equation is trivially solvable" {
    const answers = try solve(std.testing.allocator, "5: 5\n");
    try std.testing.expectEqual(@as(i64, 5), answers.part1);
    try std.testing.expectEqual(@as(i64, 5), answers.part2);
}

test "unsolvable equation contributes nothing" {
    const answers = try solve(std.testing.allocator, "10: 2 2\n");
    try std.testing.expectEqual(@as(i64, 0), answers.part1);
    try std.testing.expectEqual(@as(i64, 0), answers.part2);
}

test "malformed input returns an error" {
    try std.testing.expectError(error.InvalidEquation, solve(std.testing.allocator, ""));
    try std.testing.expectError(error.InvalidEquation, solve(std.testing.allocator, "190 10 19\n"));
    try std.testing.expectError(error.InvalidEquation, solve(std.testing.allocator, "190:\n"));
    try std.testing.expectError(error.InvalidCharacter, solve(std.testing.allocator, "abc: 1 2\n"));
}
