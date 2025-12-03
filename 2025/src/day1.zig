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

    /// Apply rotation and count zero crossings - divmod gives us both!
    fn applyAndCountCrossings(self: Rotation, start: i32) struct { position: i32, crossings: u32 } {
        return switch (self.direction) {
            .right => {
                // divmod gives us both: .div = crossings, .rem = final position
                const result = utils.divmod(start + self.distance, 100);
                return .{
                    .position = result.rem,
                    .crossings = @intCast(result.div),
                };
            },
            .left => {
                // Use divmod for final position
                // Count crossings
                const result = utils.divmod(start - self.distance, 100);
                const crossings: u32 = if (self.distance <= start)
                    0
                else blk: {
                    const cross_result = utils.divmod(self.distance - start - 1, 100);
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

pub fn main() !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    const input = try utils.readStdin(allocator);

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

    // Output results
    var stdout_buffer: [512]u8 = undefined;
    var stdout_file_writer = std.fs.File.stdout().writer(&stdout_buffer);
    const stdout = &stdout_file_writer.interface;
    try stdout.print("Part 1: {d}\n", .{part1});
    try stdout.print("Part 2: {d}\n", .{part2});
    try stdout.flush();
}
