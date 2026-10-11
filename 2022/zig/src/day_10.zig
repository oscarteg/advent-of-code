const std = @import("std");

const Instruction = union(enum) {
    noop,
    addx: i32,
};

fn parseInput(allocator: std.mem.Allocator, input: []const u8) !std.ArrayList(Instruction) {
    var instructions = std.ArrayList(Instruction).empty;
    var lines = std.mem.tokenizeScalar(u8, input, '\n');

    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;

        if (std.mem.eql(u8, trimmed, "noop")) {
            try instructions.append(allocator, .noop);
        } else if (std.mem.startsWith(u8, trimmed, "addx ")) {
            const value_str = trimmed[5..];
            const value = try std.fmt.parseInt(i32, value_str, 10);
            try instructions.append(allocator, .{ .addx = value });
        } else {
            return error.InvalidInstruction;
        }
    }

    return instructions;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) ![2][]const u8 {
    var instructions = try parseInput(allocator, input);
    defer instructions.deinit(allocator);

    var x_values = std.ArrayList(i32).empty;
    defer x_values.deinit(allocator);
    try x_values.append(allocator, 1);

    var x: i32 = 1;
    for (instructions.items) |instruction| {
        switch (instruction) {
            .noop => {
                try x_values.append(allocator, x);
            },
            .addx => |value| {
                try x_values.append(allocator, x);
                try x_values.append(allocator, x);
                x += value;
            },
        }
    }

    var part1: i32 = 0;
    const cycles = [_]usize{ 20, 60, 100, 140, 180, 220 };
    for (cycles) |cycle| {
        if (cycle < x_values.items.len) {
            part1 += @as(i32, @intCast(cycle)) * x_values.items[cycle];
        }
    }

    var part2 = std.ArrayList(u8).empty;
    for (0..6) |row| {
        if (row > 0) {
            try part2.append(allocator, '\n');
        }
        for (0..40) |col| {
            const cycle = row * 40 + col + 1;
            const pixel_x = @as(i32, @intCast(col));
            if (cycle < x_values.items.len) {
                const sprite_x = x_values.items[cycle];
                if (sprite_x - 1 <= pixel_x and pixel_x <= sprite_x + 1) {
                    try part2.append(allocator, '#');
                } else {
                    try part2.append(allocator, '.');
                }
            } else {
                try part2.append(allocator, '.');
            }
        }
    }

    const part1_str = try std.fmt.allocPrint(allocator, "{d}", .{part1});

    return [2][]const u8{ part1_str, try part2.toOwnedSlice(allocator) };
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

    try stdout.print("Part 1: {s}\nPart 2:\n{s}\n", .{ answers[0], answers[1] });
    try stdout.flush();
}

test "example part 1" {
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day_10.txt");
    const answers = try solve(allocator, input);
    defer allocator.free(answers[0]);
    defer allocator.free(answers[1]);

    try std.testing.expectEqualStrings("13140", answers[0]);
}

test "example part 2" {
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day_10.txt");
    const answers = try solve(allocator, input);
    defer allocator.free(answers[0]);
    defer allocator.free(answers[1]);

    const expected = "##..##..##..##..##..##..##..##..##..##..\n" ++
        "###...###...###...###...###...###...###.\n" ++
        "####....####....####....####....####....\n" ++
        "#####.....#####.....#####.....#####.....\n" ++
        "######......######......######......####\n" ++
        "#######.......#######.......#######.....";

    try std.testing.expectEqualStrings(expected, answers[1]);
}

test "noop only" {
    const allocator = std.testing.allocator;
    const answers = try solve(allocator, "noop\nnoop\nnoop");
    defer allocator.free(answers[0]);
    defer allocator.free(answers[1]);
    try std.testing.expectEqualStrings("0", answers[0]);
}

test "malformed instruction" {
    const allocator = std.testing.allocator;
    const result = solve(allocator, "invalid");
    try std.testing.expectError(error.InvalidInstruction, result);
}
