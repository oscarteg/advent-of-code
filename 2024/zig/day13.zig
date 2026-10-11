const std = @import("std");

const Machine = struct {
    ax: i64,
    ay: i64,
    bx: i64,
    by: i64,
    px: i64,
    py: i64,
};

const Answers = struct {
    part1: i64,
    part2: i64,
};

fn parse(allocator: std.mem.Allocator, input: []const u8) !std.ArrayList(Machine) {
    var machines = std.ArrayList(Machine).empty;
    errdefer machines.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;

        if (std.mem.startsWith(u8, trimmed, "Button A:")) {
            var ax: i64 = 0;
            var ay: i64 = 0;
            var bx: i64 = 0;
            var by: i64 = 0;
            var px: i64 = 0;
            var py: i64 = 0;

            if (std.mem.indexOf(u8, trimmed, "X+")) |pos| {
                var rest = trimmed[pos + 2 ..];
                if (std.mem.indexOf(u8, rest, ",")) |comma| {
                    ax = std.fmt.parseInt(i64, rest[0..comma], 10) catch return error.InvalidFormat;
                    rest = rest[comma + 1 ..];
                    if (std.mem.indexOf(u8, rest, "Y+")) |ypos| {
                        ay = std.fmt.parseInt(i64, rest[ypos + 2 ..], 10) catch return error.InvalidFormat;
                    }
                }
            }

            if (lines.next()) |button_b| {
                const trimmed_b = std.mem.trim(u8, button_b, " \t\r");
                if (std.mem.startsWith(u8, trimmed_b, "Button B:")) {
                    if (std.mem.indexOf(u8, trimmed_b, "X+")) |pos| {
                        var rest = trimmed_b[pos + 2 ..];
                        if (std.mem.indexOf(u8, rest, ",")) |comma| {
                            bx = std.fmt.parseInt(i64, rest[0..comma], 10) catch return error.InvalidFormat;
                            rest = rest[comma + 1 ..];
                            if (std.mem.indexOf(u8, rest, "Y+")) |ypos| {
                                by = std.fmt.parseInt(i64, rest[ypos + 2 ..], 10) catch return error.InvalidFormat;
                            }
                        }
                    }
                }
            }

            if (lines.next()) |prize_line| {
                const trimmed_p = std.mem.trim(u8, prize_line, " \t\r");
                if (std.mem.startsWith(u8, trimmed_p, "Prize:")) {
                    if (std.mem.indexOf(u8, trimmed_p, "X=")) |pos| {
                        var rest = trimmed_p[pos + 2 ..];
                        if (std.mem.indexOf(u8, rest, ",")) |comma| {
                            px = std.fmt.parseInt(i64, rest[0..comma], 10) catch return error.InvalidFormat;
                            rest = rest[comma + 1 ..];
                            if (std.mem.indexOf(u8, rest, "Y=")) |ypos| {
                                py = std.fmt.parseInt(i64, rest[ypos + 2 ..], 10) catch return error.InvalidFormat;
                            }
                        }
                    }
                }
            }

            try machines.append(allocator, Machine{
                .ax = ax,
                .ay = ay,
                .bx = bx,
                .by = by,
                .px = px,
                .py = py,
            });
        }
    }

    if (machines.items.len == 0) {
        return error.NoMachines;
    }

    return machines;
}

fn solve_machine_part1(m: Machine) i64 {
    var a: i64 = 0;
    while (a <= 100) : (a += 1) {
        var b: i64 = 0;
        while (b <= 100) : (b += 1) {
            if (a * m.ax + b * m.bx == m.px and a * m.ay + b * m.by == m.py) {
                return a * 3 + b * 1;
            }
        }
    }
    return 0;
}

fn solve_machine_part2(m: Machine) i64 {
    const det = m.ax * m.by - m.ay * m.bx;
    if (det == 0) {
        return 0;
    }

    const num_a = m.px * m.by - m.py * m.bx;
    const num_b = m.ax * m.py - m.ay * m.px;

    if (@rem(num_a, det) != 0 or @rem(num_b, det) != 0) {
        return 0;
    }

    const a = @divExact(num_a, det);
    const b = @divExact(num_b, det);

    if (a >= 0 and b >= 0 and a * m.ax + b * m.bx == m.px and a * m.ay + b * m.by == m.py) {
        return a * 3 + b * 1;
    }

    return 0;
}

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var machines = try parse(allocator, input);
    defer machines.deinit(allocator);

    var part1: i64 = 0;
    var part2: i64 = 0;

    for (machines.items) |m| {
        part1 += solve_machine_part1(m);

        var m2 = m;
        m2.px += 10000000000000;
        m2.py += 10000000000000;
        part2 += solve_machine_part2(m2);
    }

    return Answers{ .part1 = part1, .part2 = part2 };
}

pub fn main(init: std.process.Init) !void {
    const arena = init.arena.allocator();
    var stdin_buffer: [8192]u8 = undefined;
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
    const example = @embedFile("examples/day13.txt");
    const allocator = std.testing.allocator;
    const result = try solve(allocator, example);
    try std.testing.expectEqual(@as(i64, 480), result.part1);
}

test "part2 example" {
    const example = @embedFile("examples/day13.txt");
    const allocator = std.testing.allocator;
    const result = try solve(allocator, example);
    try std.testing.expectEqual(@as(i64, 875318608908), result.part2);
}

test "unsolvable machine" {
    const input = "Button A: X+1, Y+1\nButton B: X+2, Y+2\nPrize: X=3, Y=4\n";
    const allocator = std.testing.allocator;
    const result = try solve(allocator, input);
    try std.testing.expectEqual(@as(i64, 0), result.part1);
}

test "solvable machine" {
    const input = "Button A: X+1, Y+0\nButton B: X+0, Y+1\nPrize: X=3, Y=2\n";
    const allocator = std.testing.allocator;
    const result = try solve(allocator, input);
    try std.testing.expectEqual(@as(i64, 11), result.part1);
}

test "rejects a non-numeric button offset" {
    const input = "Button A: X+a, Y+1\nButton B: X+2, Y+2\nPrize: X=3, Y=4\n";
    try std.testing.expectError(error.InvalidFormat, solve(std.testing.allocator, input));
}

test "rejects input without machines" {
    try std.testing.expectError(error.NoMachines, solve(std.testing.allocator, "\n\n"));
}
