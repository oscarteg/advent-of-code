const std = @import("std");

/// Playground: junction boxes hang in 3D space and get connected by strings
/// of lights, shortest straight-line distance first; every connection merges
/// the circuits of the two boxes it joins. Part 1 makes the closest
/// `connection_count` connections - a pair whose boxes already share a
/// circuit still counts, it just changes nothing - and multiplies the sizes
/// of the three largest circuits. Part 2 keeps connecting until a single
/// circuit remains and multiplies the X coordinates of the last pair that
/// merged two circuits.
const Answers = struct { part1: u64, part2: i64 };

/// The puzzle connects the 1000 closest pairs of the real input and the 10
/// closest of its 20-box example, so the count follows the input size.
const real_connections: usize = 1000;
const example_boxes: usize = 20;
const example_connections: usize = 10;

const Box = struct {
    x: i32,
    y: i32,
    z: i32,

    fn distanceSquared(self: Box, other: Box) u64 {
        const dx = @as(i64, self.x) - other.x;
        const dy = @as(i64, self.y) - other.y;
        const dz = @as(i64, self.z) - other.z;
        return @intCast(dx * dx + dy * dy + dz * dz);
    }
};

/// A candidate string of lights between two boxes. Squared distances order
/// the same way straight-line distances do, so no square root is needed.
const Connection = struct {
    first: u32,
    second: u32,
    distance_squared: u64,

    /// Equal distances fall back to the box order, keeping ties deterministic.
    fn closer(_: void, lhs: Connection, rhs: Connection) bool {
        if (lhs.distance_squared != rhs.distance_squared) {
            return lhs.distance_squared < rhs.distance_squared;
        }
        if (lhs.first != rhs.first) return lhs.first < rhs.first;
        return lhs.second < rhs.second;
    }
};

/// Which boxes share a circuit, as a union-find forest merged by size.
const Circuits = struct {
    parent: []u32,
    size: []u32,
    count: usize,

    fn init(allocator: std.mem.Allocator, boxes: usize) !Circuits {
        const parent = try allocator.alloc(u32, boxes);
        errdefer allocator.free(parent);
        const size = try allocator.alloc(u32, boxes);
        for (parent, size, 0..) |*box_parent, *box_size, index| {
            box_parent.* = @intCast(index);
            box_size.* = 1;
        }
        return .{ .parent = parent, .size = size, .count = boxes };
    }

    fn deinit(self: Circuits, allocator: std.mem.Allocator) void {
        allocator.free(self.parent);
        allocator.free(self.size);
    }

    fn find(self: Circuits, box: u32) u32 {
        var root = box;
        while (self.parent[root] != root) {
            self.parent[root] = self.parent[self.parent[root]]; // path halving
            root = self.parent[root];
        }
        return root;
    }

    /// Connects two boxes, reporting whether this merged two circuits.
    fn connect(self: *Circuits, first: u32, second: u32) bool {
        var root = self.find(first);
        var merged = self.find(second);
        if (root == merged) return false;
        if (self.size[root] < self.size[merged]) std.mem.swap(u32, &root, &merged);
        self.parent[merged] = root;
        self.size[root] += self.size[merged];
        self.count -= 1;
        return true;
    }

    fn threeLargestProduct(self: Circuits) !u64 {
        if (self.count < 3) return error.TooFewCircuits;
        var largest = [3]u64{ 0, 0, 0 };
        for (self.parent, self.size, 0..) |box_parent, box_size, index| {
            if (box_parent != index) continue;
            var candidate: u64 = box_size;
            for (&largest) |*slot| {
                if (candidate > slot.*) std.mem.swap(u64, slot, &candidate);
            }
        }
        return largest[0] * largest[1] * largest[2];
    }
};

fn parseBoxes(allocator: std.mem.Allocator, input: []const u8) ![]Box {
    var boxes: std.ArrayList(Box) = .empty;
    errdefer boxes.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, &std.ascii.whitespace);
        if (trimmed.len == 0) continue;

        var fields = std.mem.splitScalar(u8, trimmed, ',');
        var coordinates: [3]i32 = undefined;
        for (&coordinates) |*coordinate| {
            const field = fields.next() orelse return error.InvalidBox;
            coordinate.* = try std.fmt.parseInt(i32, std.mem.trim(u8, field, &std.ascii.whitespace), 10);
        }
        if (fields.next() != null) return error.InvalidBox;

        try boxes.append(allocator, .{ .x = coordinates[0], .y = coordinates[1], .z = coordinates[2] });
    }
    if (boxes.items.len < 2) return error.NotEnoughBoxes;

    return boxes.toOwnedSlice(allocator);
}

/// Every pair of boxes, closest first.
fn sortedConnections(allocator: std.mem.Allocator, boxes: []const Box) ![]Connection {
    const connections = try allocator.alloc(Connection, boxes.len * (boxes.len - 1) / 2);
    var next: usize = 0;
    for (boxes, 0..) |box, first| {
        for (boxes[first + 1 ..], first + 1..) |other, second| {
            connections[next] = .{
                .first = @intCast(first),
                .second = @intCast(second),
                .distance_squared = box.distanceSquared(other),
            };
            next += 1;
        }
    }
    std.mem.sortUnstable(Connection, connections, {}, Connection.closer);
    return connections;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    return solveWith(allocator, input, null);
}

/// `connection_override` replaces the number of connections part 1 makes;
/// null follows the count the puzzle prescribes for this input size.
fn solveWith(allocator: std.mem.Allocator, input: []const u8, connection_override: ?usize) !Answers {
    const boxes = try parseBoxes(allocator, input);
    defer allocator.free(boxes);

    const connection_count = connection_override orelse
        (if (boxes.len <= example_boxes) example_connections else real_connections);

    const connections = try sortedConnections(allocator, boxes);
    defer allocator.free(connections);

    var circuits = try Circuits.init(allocator, boxes.len);
    defer circuits.deinit(allocator);

    // The first connection always merges, so it is a safe starting value.
    var last_merge = connections[0];
    var part1: ?u64 = null;
    for (connections, 0..) |connection, index| {
        if (index == connection_count) part1 = try circuits.threeLargestProduct();
        if (circuits.connect(connection.first, connection.second)) last_merge = connection;
    }

    return .{
        .part1 = part1 orelse try circuits.threeLargestProduct(),
        .part2 = @as(i64, boxes[last_merge.first].x) * boxes[last_merge.second].x,
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

test "official example" {
    const input = @embedFile("examples/day8.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u64, 40), answers.part1);
    try std.testing.expectEqual(@as(i64, 25272), answers.part2);
}

test "two close pairs are bridged by the long connection between them" {
    const input = "0,0,0\n1,0,0\n10,0,0\n11,0,0\n";
    const answers = try solveWith(std.testing.allocator, input, 1);
    try std.testing.expectEqual(@as(u64, 2), answers.part1);
    try std.testing.expectEqual(@as(i64, 10), answers.part2);
}

test "a pair already sharing a circuit counts but changes nothing" {
    const input = "0,0,0\n1,0,0\n2,0,0\n100,0,0\n200,0,0\n";
    const answers = try solveWith(std.testing.allocator, input, 3);
    try std.testing.expectEqual(@as(u64, 3), answers.part1);
    try std.testing.expectEqual(@as(i64, 20_000), answers.part2);
}

test "surrounding whitespace is ignored" {
    const input = "  0,0,0 \n\n 1,0,0\n10,0,0\n 11,0,0 \n\n";
    const answers = try solveWith(std.testing.allocator, input, 1);
    try std.testing.expectEqual(@as(u64, 2), answers.part1);
    try std.testing.expectEqual(@as(i64, 10), answers.part2);
}

test "rejects an input left with fewer than three circuits" {
    const input = "0,0,0\n1,0,0\n";
    try std.testing.expectError(error.TooFewCircuits, solveWith(std.testing.allocator, input, 1));
}

test "rejects a box missing a coordinate" {
    const input = "1,2\n3,4,5\n";
    try std.testing.expectError(error.InvalidBox, solve(std.testing.allocator, input));
}

test "rejects a box with a fourth coordinate" {
    const input = "1,2,3,4\n5,6,7\n";
    try std.testing.expectError(error.InvalidBox, solve(std.testing.allocator, input));
}

test "rejects a non-numeric coordinate" {
    const input = "1,2,x\n3,4,5\n";
    try std.testing.expectError(error.InvalidCharacter, solve(std.testing.allocator, input));
}

test "rejects fewer than two boxes" {
    try std.testing.expectError(error.NotEnoughBoxes, solve(std.testing.allocator, "1,2,3\n"));
    try std.testing.expectError(error.NotEnoughBoxes, solve(std.testing.allocator, "  \n \n"));
}
