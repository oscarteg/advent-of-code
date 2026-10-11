const std = @import("std");

pub const Answers = struct {
    part1: i64,
    part2: i64,
};

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var grid: std.ArrayList([]const u8) = .empty;
    defer grid.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        if (line.len == 0) continue;
        try grid.append(allocator, line);

        for (line) |c| {
            if (c < '0' or c > '9') {
                return error.InvalidCharacter;
            }
        }
    }

    if (grid.items.len == 0) {
        return error.EmptyInput;
    }

    const rows = grid.items.len;
    const cols = grid.items[0].len;

    for (grid.items) |row| {
        if (row.len != cols) {
            return error.IrregularGrid;
        }
    }

    var part1: i64 = 0;
    var part2: i64 = 0;

    for (0..rows) |r| {
        for (0..cols) |c| {
            if (grid.items[r][c] == '0') {
                var reachable_nines = std.AutoHashMap([2]usize, void).init(allocator);
                defer reachable_nines.deinit();

                var path_count: i64 = 0;

                var stack: std.ArrayList([2]usize) = .empty;
                defer stack.deinit(allocator);
                try stack.append(allocator, .{ r, c });

                while (stack.pop()) |pos| {
                    const cr = pos[0];
                    const cc = pos[1];

                    if (grid.items[cr][cc] == '9') {
                        try reachable_nines.put(.{ cr, cc }, {});
                        path_count += 1;
                        continue;
                    }

                    const current_elev = grid.items[cr][cc];
                    const next_elev = current_elev + 1;

                    const dr = [_]i32{ -1, 1, 0, 0 };
                    const dc = [_]i32{ 0, 0, -1, 1 };

                    for (&dr, &dc) |d_r, d_c| {
                        const nr: i32 = @as(i32, @intCast(cr)) + d_r;
                        const nc: i32 = @as(i32, @intCast(cc)) + d_c;

                        if (nr >= 0 and nr < @as(i32, @intCast(rows)) and
                            nc >= 0 and nc < @as(i32, @intCast(cols)))
                        {
                            const ur = @as(usize, @intCast(nr));
                            const uc = @as(usize, @intCast(nc));
                            if (grid.items[ur][uc] == next_elev) {
                                try stack.append(allocator, .{ ur, uc });
                            }
                        }
                    }
                }

                part1 += @as(i64, @intCast(reachable_nines.count()));
                part2 += path_count;
            }
        }
    }

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

test "example" {
    const example = @embedFile("examples/day10.txt");
    const answers = try solve(std.testing.allocator, example);
    try std.testing.expectEqual(36, answers.part1);
    try std.testing.expectEqual(81, answers.part2);
}

test "single trail" {
    const input = "0123456789\n";
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(1, answers.part1);
    try std.testing.expectEqual(1, answers.part2);
}

test "vertical trail" {
    const input = "0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n";
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(1, answers.part1);
    try std.testing.expectEqual(1, answers.part2);
}

test "empty input" {
    try std.testing.expectError(error.EmptyInput, solve(std.testing.allocator, ""));
}

test "invalid character" {
    try std.testing.expectError(error.InvalidCharacter, solve(std.testing.allocator, "0a1\n"));
}

test "irregular grid" {
    try std.testing.expectError(error.IrregularGrid, solve(std.testing.allocator, "012\n01\n"));
}
