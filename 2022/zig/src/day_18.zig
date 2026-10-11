const std = @import("std");

const Cube = struct {
    x: i32,
    y: i32,
    z: i32,

    fn eql(a: Cube, b: Cube) bool {
        return a.x == b.x and a.y == b.y and a.z == b.z;
    }

    fn hash(self: Cube) u64 {
        var hasher = std.hash.Wyhash.init(0);
        hasher.update(std.mem.asBytes(&self.x));
        hasher.update(std.mem.asBytes(&self.y));
        hasher.update(std.mem.asBytes(&self.z));
        return hasher.final();
    }
};

fn parseInput(allocator: std.mem.Allocator, input: []const u8) !std.AutoHashMap(Cube, void) {
    var cubes = std.AutoHashMap(Cube, void).init(allocator);
    var lines = std.mem.splitScalar(u8, input, '\n');
    var count: usize = 0;

    while (lines.next()) |line| {
        if (line.len == 0) continue;
        var parts = std.mem.splitScalar(u8, line, ',');
        const x_str = parts.next() orelse return error.InvalidFormat;
        const y_str = parts.next() orelse return error.InvalidFormat;
        const z_str = parts.next() orelse return error.InvalidFormat;

        const x = try std.fmt.parseInt(i32, x_str, 10);
        const y = try std.fmt.parseInt(i32, y_str, 10);
        const z = try std.fmt.parseInt(i32, z_str, 10);

        try cubes.put(Cube{ .x = x, .y = y, .z = z }, {});
        count += 1;
    }

    if (count == 0) return error.NoInput;
    return cubes;
}

fn surfaceArea(cubes: std.AutoHashMap(Cube, void)) usize {
    var area: usize = 0;
    var it = cubes.keyIterator();
    while (it.next()) |cube| {
        const neighbors = [_]Cube{
            .{ .x = cube.x + 1, .y = cube.y, .z = cube.z },
            .{ .x = cube.x - 1, .y = cube.y, .z = cube.z },
            .{ .x = cube.x, .y = cube.y + 1, .z = cube.z },
            .{ .x = cube.x, .y = cube.y - 1, .z = cube.z },
            .{ .x = cube.x, .y = cube.y, .z = cube.z + 1 },
            .{ .x = cube.x, .y = cube.y, .z = cube.z - 1 },
        };
        for (neighbors) |n| {
            if (!cubes.contains(n)) {
                area += 1;
            }
        }
    }
    return area;
}

fn externalSurfaceArea(allocator: std.mem.Allocator, cubes: std.AutoHashMap(Cube, void)) !usize {
    var min_x: i32 = std.math.maxInt(i32);
    var max_x: i32 = std.math.minInt(i32);
    var min_y: i32 = std.math.maxInt(i32);
    var max_y: i32 = std.math.minInt(i32);
    var min_z: i32 = std.math.maxInt(i32);
    var max_z: i32 = std.math.minInt(i32);

    var it = cubes.keyIterator();
    while (it.next()) |cube| {
        min_x = @min(min_x, cube.x);
        max_x = @max(max_x, cube.x);
        min_y = @min(min_y, cube.y);
        max_y = @max(max_y, cube.y);
        min_z = @min(min_z, cube.z);
        max_z = @max(max_z, cube.z);
    }

    min_x -= 1;
    max_x += 1;
    min_y -= 1;
    max_y += 1;
    min_z -= 1;
    max_z += 1;

    var external = std.AutoHashMap(Cube, void).init(allocator);
    defer external.deinit();

    var queue: std.ArrayList(Cube) = .empty;
    defer queue.deinit(allocator);

    const start = Cube{ .x = min_x, .y = min_y, .z = min_z };
    try external.put(start, {});
    try queue.append(allocator, start);

    while (queue.pop()) |cube| {
        const neighbors = [_]Cube{
            .{ .x = cube.x + 1, .y = cube.y, .z = cube.z },
            .{ .x = cube.x - 1, .y = cube.y, .z = cube.z },
            .{ .x = cube.x, .y = cube.y + 1, .z = cube.z },
            .{ .x = cube.x, .y = cube.y - 1, .z = cube.z },
            .{ .x = cube.x, .y = cube.y, .z = cube.z + 1 },
            .{ .x = cube.x, .y = cube.y, .z = cube.z - 1 },
        };

        for (neighbors) |n| {
            if (n.x >= min_x and n.x <= max_x and
                n.y >= min_y and n.y <= max_y and
                n.z >= min_z and n.z <= max_z and
                !cubes.contains(n) and
                !external.contains(n))
            {
                try external.put(n, {});
                try queue.append(allocator, n);
            }
        }
    }

    var area: usize = 0;
    it = cubes.keyIterator();
    while (it.next()) |cube| {
        const neighbors = [_]Cube{
            .{ .x = cube.x + 1, .y = cube.y, .z = cube.z },
            .{ .x = cube.x - 1, .y = cube.y, .z = cube.z },
            .{ .x = cube.x, .y = cube.y + 1, .z = cube.z },
            .{ .x = cube.x, .y = cube.y - 1, .z = cube.z },
            .{ .x = cube.x, .y = cube.y, .z = cube.z + 1 },
            .{ .x = cube.x, .y = cube.y, .z = cube.z - 1 },
        };
        for (neighbors) |n| {
            if (external.contains(n)) {
                area += 1;
            }
        }
    }
    return area;
}

const Answers = struct {
    part1: usize,
    part2: usize,
};

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var cubes = try parseInput(allocator, input);
    defer cubes.deinit();

    const part1 = surfaceArea(cubes);
    const part2 = try externalSurfaceArea(allocator, cubes);

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

test "tiny example" {
    const allocator = std.testing.allocator;
    const input = "1,1,1\n2,1,1";
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(@as(usize, 10), answers.part1);
}

test "example" {
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day_18.txt");
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(@as(usize, 64), answers.part1);
    try std.testing.expectEqual(@as(usize, 58), answers.part2);
}

test "malformed input" {
    const allocator = std.testing.allocator;
    try std.testing.expectError(error.InvalidFormat, solve(allocator, "1,1"));
    try std.testing.expectError(error.NoInput, solve(allocator, ""));
}
