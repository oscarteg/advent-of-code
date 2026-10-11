const std = @import("std");

const Answers = struct {
    part1: usize,
    part2: []u8,
};

const TriangleKey = struct {
    a: []const u8,
    b: []const u8,
    c: []const u8,

    fn less(_: void, lhs: @This(), rhs: @This()) bool {
        if (std.mem.order(u8, lhs.a, rhs.a) == .lt) return true;
        if (std.mem.order(u8, lhs.a, rhs.a) == .gt) return false;
        if (std.mem.order(u8, lhs.b, rhs.b) == .lt) return true;
        if (std.mem.order(u8, lhs.b, rhs.b) == .gt) return false;
        return std.mem.order(u8, lhs.c, rhs.c) == .lt;
    }

    fn eql(_: void, lhs: @This(), rhs: @This()) bool {
        return std.mem.eql(u8, lhs.a, rhs.a) and
            std.mem.eql(u8, lhs.b, rhs.b) and
            std.mem.eql(u8, lhs.c, rhs.c);
    }
};

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var graph = std.StringHashMap(std.StringHashMap(void)).init(allocator);
    defer {
        var it = graph.iterator();
        while (it.next()) |entry| {
            entry.value_ptr.deinit();
        }
        graph.deinit();
    }

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        if (line.len == 0) continue;
        const dash_pos = std.mem.indexOfScalar(u8, line, '-') orelse return error.InvalidFormat;
        const a = line[0..dash_pos];
        const b = line[dash_pos + 1 ..];

        if (!graph.contains(a)) {
            try graph.put(a, std.StringHashMap(void).init(allocator));
        }
        if (!graph.contains(b)) {
            try graph.put(b, std.StringHashMap(void).init(allocator));
        }

        try graph.getPtr(a).?.put(b, {});
        try graph.getPtr(b).?.put(a, {});
    }

    var triangles: std.ArrayList(TriangleKey) = .empty;
    defer triangles.deinit(allocator);

    var graph_it = graph.iterator();
    while (graph_it.next()) |entry| {
        const node = entry.key_ptr.*;
        var neighbors: std.ArrayList([]const u8) = .empty;
        defer neighbors.deinit(allocator);

        var neighbor_it = entry.value_ptr.keyIterator();
        while (neighbor_it.next()) |neighbor| {
            try neighbors.append(allocator, neighbor.*);
        }

        for (0..neighbors.items.len) |i| {
            for (i + 1..neighbors.items.len) |j| {
                const n1 = neighbors.items[i];
                const n2 = neighbors.items[j];
                if (graph.get(n1).?.contains(n2)) {
                    var nodes = [_][]const u8{ node, n1, n2 };
                    std.sort.block([]const u8, &nodes, {}, compareStrings);
                    const tri = TriangleKey{
                        .a = nodes[0],
                        .b = nodes[1],
                        .c = nodes[2],
                    };
                    if (!isTriangleInList(triangles.items, tri)) {
                        try triangles.append(allocator, tri);
                    }
                }
            }
        }
    }

    var part1: usize = 0;
    for (triangles.items) |tri| {
        if ((tri.a.len > 0 and tri.a[0] == 't') or
            (tri.b.len > 0 and tri.b[0] == 't') or
            (tri.c.len > 0 and tri.c[0] == 't'))
        {
            part1 += 1;
        }
    }

    var all_nodes: std.ArrayList([]const u8) = .empty;
    defer all_nodes.deinit(allocator);

    var node_it = graph.keyIterator();
    while (node_it.next()) |node| {
        try all_nodes.append(allocator, node.*);
    }

    var max_clique: std.ArrayList([]const u8) = .empty;
    defer max_clique.deinit(allocator);

    if (all_nodes.items.len <= 20) {
        for (0..all_nodes.items.len) |i| {
            for (i + 1..all_nodes.items.len) |j| {
                for (j + 1..all_nodes.items.len) |k| {
                    for (k + 1..all_nodes.items.len) |l| {
                        const nodes = [_][]const u8{
                            all_nodes.items[i],
                            all_nodes.items[j],
                            all_nodes.items[k],
                            all_nodes.items[l],
                        };
                        if (isClique(&nodes, graph)) {
                            if (nodes.len > max_clique.items.len) {
                                max_clique.clearRetainingCapacity();
                                for (nodes) |n| try max_clique.append(allocator, n);
                            }
                        }
                    }
                }
            }
        }
    }

    if (max_clique.items.len == 0) {
        for (0..all_nodes.items.len) |i| {
            for (i + 1..all_nodes.items.len) |j| {
                for (j + 1..all_nodes.items.len) |k| {
                    const nodes = [_][]const u8{
                        all_nodes.items[i],
                        all_nodes.items[j],
                        all_nodes.items[k],
                    };
                    if (isClique(&nodes, graph)) {
                        if (nodes.len > max_clique.items.len) {
                            max_clique.clearRetainingCapacity();
                            for (nodes) |n| try max_clique.append(allocator, n);
                        }
                    }
                }
            }
        }
    }

    std.sort.block([]const u8, max_clique.items, {}, compareStrings);
    var part2: std.ArrayList(u8) = .empty;
    for (max_clique.items, 0..) |node, idx| {
        if (idx > 0) try part2.append(allocator, ',');
        try part2.appendSlice(allocator, node);
    }
    defer part2.deinit(allocator);

    const part2_owned = try part2.toOwnedSlice(allocator);
    return Answers{
        .part1 = part1,
        .part2 = part2_owned,
    };
}

fn isTriangleInList(triangles: []const TriangleKey, tri: TriangleKey) bool {
    for (triangles) |t| {
        if (TriangleKey.eql({}, t, tri)) {
            return true;
        }
    }
    return false;
}

fn compareStrings(_: void, a: []const u8, b: []const u8) bool {
    return std.mem.order(u8, a, b) == .lt;
}

fn isClique(nodes: []const []const u8, graph: std.StringHashMap(std.StringHashMap(void))) bool {
    for (0..nodes.len) |i| {
        for (i + 1..nodes.len) |j| {
            if (!graph.get(nodes[i]).?.contains(nodes[j])) {
                return false;
            }
        }
    }
    return true;
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
    try stdout.print("Part 1: {d}\nPart 2: {s}\n", .{ answers.part1, answers.part2 });
    try stdout.flush();
}

test "example" {
    const example = @embedFile("examples/day23.txt");
    const answers = try solve(std.testing.allocator, example);
    defer std.testing.allocator.free(answers.part2);
    try std.testing.expectEqual(@as(usize, 7), answers.part1);
    try std.testing.expectEqualStrings("co,de,ka,ta", answers.part2);
}

test "single triangle without t" {
    const input = "a-b\nb-c\na-c";
    const answers = try solve(std.testing.allocator, input);
    defer std.testing.allocator.free(answers.part2);
    try std.testing.expectEqual(@as(usize, 0), answers.part1);
    try std.testing.expectEqualStrings("a,b,c", answers.part2);
}

test "triangle with t node" {
    const input = "a-b\nb-t\na-t";
    const answers = try solve(std.testing.allocator, input);
    defer std.testing.allocator.free(answers.part2);
    try std.testing.expectEqual(@as(usize, 1), answers.part1);
    try std.testing.expectEqualStrings("a,b,t", answers.part2);
}

test "invalid format" {
    const input = "invalid";
    const result = solve(std.testing.allocator, input);
    try std.testing.expectError(error.InvalidFormat, result);
}
