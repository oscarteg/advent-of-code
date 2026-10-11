const std = @import("std");

const Move = enum(u32) {
    rock = 1,
    paper = 2,
    scissors = 3,
};

const Outcome = enum(u32) {
    loss = 0,
    draw = 3,
    win = 6,
};

fn parseMoveFromRound(ch: u8) !Move {
    return switch (ch) {
        'A', 'X' => .rock,
        'B', 'Y' => .paper,
        'C', 'Z' => .scissors,
        else => error.InvalidMove,
    };
}

fn parseOutcomeFromCode(ch: u8) !Outcome {
    return switch (ch) {
        'X' => .loss,
        'Y' => .draw,
        'Z' => .win,
        else => error.InvalidOutcome,
    };
}

fn moveThatDefeats(opponent: Move) Move {
    return switch (opponent) {
        .rock => .paper,
        .paper => .scissors,
        .scissors => .rock,
    };
}

fn moveThatLosesTo(opponent: Move) Move {
    return switch (opponent) {
        .rock => .scissors,
        .paper => .rock,
        .scissors => .paper,
    };
}

fn scoreRound(our_move: Move, opponent: Move) u32 {
    const outcome: u32 = if (our_move == opponent)
        @intFromEnum(Outcome.draw)
    else if ((our_move == .rock and opponent == .scissors) or
        (our_move == .paper and opponent == .rock) or
        (our_move == .scissors and opponent == .paper))
        @intFromEnum(Outcome.win)
    else
        @intFromEnum(Outcome.loss);
    return @intFromEnum(our_move) + outcome;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) ![2]u32 {
    var part1: u32 = 0;
    var part2: u32 = 0;

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        if (line.len == 0) continue;
        const trimmed = std.mem.trim(u8, line, " \r\t");
        if (trimmed.len < 3) continue;

        const opponent_char = trimmed[0];
        const code_char = trimmed[2];

        const opponent = try parseMoveFromRound(opponent_char);

        const our_move_part1 = try parseMoveFromRound(code_char);
        part1 += scoreRound(our_move_part1, opponent);

        const outcome = try parseOutcomeFromCode(code_char);
        const our_move_part2 = switch (outcome) {
            .loss => moveThatLosesTo(opponent),
            .draw => opponent,
            .win => moveThatDefeats(opponent),
        };
        part2 += scoreRound(our_move_part2, opponent);
    }

    _ = allocator;
    return [2]u32{ part1, part2 };
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
    try stdout.print("Part 1: {d}\nPart 2: {d}\n", .{ answers[0], answers[1] });
    try stdout.flush();
}

test "part1 example" {
    const input = @embedFile("examples/day_2.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u32, 15), answers[0]);
}

test "part2 example" {
    const input = @embedFile("examples/day_2.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u32, 12), answers[1]);
}

test "part1 win" {
    const answers = try solve(std.testing.allocator, "A Y\n");
    try std.testing.expectEqual(@as(u32, 8), answers[0]);
}

test "part1 loss" {
    const answers = try solve(std.testing.allocator, "B X\n");
    try std.testing.expectEqual(@as(u32, 1), answers[0]);
}

test "part1 draw" {
    const answers = try solve(std.testing.allocator, "C Z\n");
    try std.testing.expectEqual(@as(u32, 6), answers[0]);
}

test "part2 loss" {
    const answers = try solve(std.testing.allocator, "B X\n");
    try std.testing.expectEqual(@as(u32, 1), answers[1]);
}

test "part2 draw" {
    const answers = try solve(std.testing.allocator, "C Y\n");
    try std.testing.expectEqual(@as(u32, 6), answers[1]);
}

test "part2 win" {
    const answers = try solve(std.testing.allocator, "A Z\n");
    try std.testing.expectEqual(@as(u32, 8), answers[1]);
}

test "invalid move" {
    const result = solve(std.testing.allocator, "D X\n");
    try std.testing.expectError(error.InvalidMove, result);
}
