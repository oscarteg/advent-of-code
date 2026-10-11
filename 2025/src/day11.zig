//! Advent of Code 2025 day 11: Reactor.
//!
//! Each line names a device and the devices its outputs feed into; together
//! they form a directed acyclic graph. Part 1 counts the paths from `you` to
//! `out`. Part 2 counts the paths from `svr` to `out` that visit both `dac`
//! and `fft`, in either order.

const std = @import("std");

const Answers = struct { part1: u64, part2: u64 };

const Graph = std.StringHashMap([]const []const u8);

fn parseGraph(arena: std.mem.Allocator, input: []const u8) !Graph {
    var graph = Graph.init(arena);
    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, &std.ascii.whitespace);
        if (trimmed.len == 0) continue;

        const colon = std.mem.indexOfScalar(u8, trimmed, ':') orelse return error.MissingColon;
        const device = std.mem.trim(u8, trimmed[0..colon], &std.ascii.whitespace);
        if (device.len == 0) return error.MissingDevice;

        var outputs: std.ArrayList([]const u8) = .empty;
        var names = std.mem.tokenizeAny(u8, trimmed[colon + 1 ..], &std.ascii.whitespace);
        while (names.next()) |name| try outputs.append(arena, name);
        try graph.put(device, outputs.items);
    }
    if (graph.count() == 0) return error.EmptyInput;
    return graph;
}

fn countPaths(graph: Graph, memo: *std.StringHashMap(u64), from: []const u8, to: []const u8) !u64 {
    if (std.mem.eql(u8, from, to)) return 1;
    if (memo.get(from)) |known| return known;
    const outputs = graph.get(from) orelse return 0;
    var total: u64 = 0;
    for (outputs) |next| total += try countPaths(graph, memo, next, to);
    try memo.put(from, total);
    return total;
}

fn pathsVia(arena: std.mem.Allocator, graph: Graph, waypoints: []const []const u8) !u64 {
    var total: u64 = 1;
    for (waypoints[0 .. waypoints.len - 1], waypoints[1..]) |from, to| {
        var memo = std.StringHashMap(u64).init(arena);
        total *= try countPaths(graph, &memo, from, to);
    }
    return total;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var arena_state = std.heap.ArenaAllocator.init(allocator);
    defer arena_state.deinit();
    const arena = arena_state.allocator();

    const graph = try parseGraph(arena, input);
    // The graph is acyclic, so at most one of the two waypoint orders has any paths.
    const dac_first = try pathsVia(arena, graph, &.{ "svr", "dac", "fft", "out" });
    const fft_first = try pathsVia(arena, graph, &.{ "svr", "fft", "dac", "out" });
    return .{
        .part1 = try pathsVia(arena, graph, &.{ "you", "out" }),
        .part2 = dac_first + fft_first,
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

const part2_example =
    \\svr: aaa bbb
    \\aaa: fft
    \\fft: ccc
    \\bbb: tty
    \\tty: ccc
    \\ccc: ddd eee
    \\ddd: hub
    \\hub: fff
    \\eee: dac
    \\dac: fff
    \\fff: ggg hhh
    \\ggg: out
    \\hhh: out
;

test "official part 1 example" {
    const answers = try solve(std.testing.allocator, @embedFile("examples/day11.txt"));
    try std.testing.expectEqual(@as(u64, 5), answers.part1);
    try std.testing.expectEqual(@as(u64, 0), answers.part2);
}

test "official part 2 example visits fft before dac" {
    const answers = try solve(std.testing.allocator, part2_example);
    try std.testing.expectEqual(@as(u64, 2), answers.part2);
}

test "paths visiting dac before fft count too" {
    const answers = try solve(std.testing.allocator, "svr: a b\na: dac\nb: dac\ndac: fft\nfft: out\n");
    try std.testing.expectEqual(@as(u64, 2), answers.part2);
}

test "a path that skips a waypoint does not count" {
    const answers = try solve(std.testing.allocator, "svr: dac\ndac: out\nfft: out\n");
    try std.testing.expectEqual(@as(u64, 0), answers.part2);
}

test "a single direct path" {
    const answers = try solve(std.testing.allocator, "you: out");
    try std.testing.expectEqual(@as(u64, 1), answers.part1);
}

test "rejects a line without a colon" {
    try std.testing.expectError(error.MissingColon, solve(std.testing.allocator, "you out\n"));
}

test "rejects a line without a device name" {
    try std.testing.expectError(error.MissingDevice, solve(std.testing.allocator, ": out\n"));
}

test "rejects empty input" {
    try std.testing.expectError(error.EmptyInput, solve(std.testing.allocator, " \n\n"));
}
