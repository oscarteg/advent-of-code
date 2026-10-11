const std = @import("std");
const utils = @import("utils.zig");

const RotationDir = enum {
    left,
    right,

    fn fromChar(c: u8) !RotationDir {
        return switch (c) {
            'L' => .left,
            'R' => .right,
            else => error.InvalidDirection,
        };
    }
};

const Rotation = struct {
    direction: RotationDir,
    distance: i32,

    fn parse(line: []const u8) !Rotation {
        if (line.len < 2) return error.InvalidFormat;

        const direction = try RotationDir.fromChar(line[0]);
        const distance = try utils.parseInt(i32, line[1..]);

        return .{
            .direction = direction,
            .distance = distance,
        };
    }

    fn applyAndCountCrossings(self: Rotation, start: i32) struct { position: i32, crossings: u32 } {
        return switch (self.direction) {
            .right => {
                const result = utils.divmod(start + self.distance, 100);
                return .{
                    .position = result.rem,
                    .crossings = @intCast(result.div),
                };
            },
            .left => {
                const result = utils.divmod(start - self.distance, 100);
                const crossings: u32 = if (start == 0)
                    @intCast(utils.divmod(self.distance, 100).div)
                else if (self.distance < start)
                    0
                else blk: {
                    const cross_result = utils.divmod(self.distance - start, 100);
                    break :blk @intCast(1 + cross_result.div);
                };
                return .{
                    .position = result.rem,
                    .crossings = crossings,
                };
            },
        };
    }
};

pub const Answer = struct {
    part1: u32,
    part2: u32,
};

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answer {
    _ = allocator;

    var position: i32 = 50;
    var part1: u32 = 0;
    var part2: u32 = 0;

    var lines = utils.LineIterator.init(input);
    while (lines.next()) |line| {
        if (line.len == 0) continue;
        const rotation = try Rotation.parse(line);

        // Get both position and crossings from divmod in one call!
        const result = rotation.applyAndCountCrossings(position);
        position = result.position;
        part2 += result.crossings;

        // Part 1: Count times we end at 0
        if (position == 0) {
            part1 += 1;
        }
    }

    return Answer{
        .part1 = part1,
        .part2 = part2,
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

test "example input" {
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day1.txt");
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(@as(u32, 3), answers.part1);
    try std.testing.expectEqual(@as(u32, 6), answers.part2);
}

test "single rotation ending at zero" {
    const allocator = std.testing.allocator;
    const input = "R50\n";
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(@as(u32, 1), answers.part1);
    try std.testing.expectEqual(@as(u32, 1), answers.part2);
}

test "rotation with multiple crossings" {
    const allocator = std.testing.allocator;
    const input = "R250\n";
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(@as(u32, 1), answers.part1);
    try std.testing.expectEqual(@as(u32, 3), answers.part2);
}

test "rotation from zero" {
    const allocator = std.testing.allocator;
    const input = "R50\nL50\n";
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(@as(u32, 1), answers.part1);
    try std.testing.expectEqual(@as(u32, 1), answers.part2);
}

test "malformed input returns error" {
    const allocator = std.testing.allocator;
    const input = "X50\n";
    const result = solve(allocator, input);
    try std.testing.expectError(error.InvalidDirection, result);
}

test "missing distance returns error" {
    const allocator = std.testing.allocator;
    const input = "L\n";
    const result = solve(allocator, input);
    try std.testing.expectError(error.InvalidFormat, result);
}
