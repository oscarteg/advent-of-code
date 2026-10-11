const std = @import("std");

fn gcd(a: usize, b: usize) usize {
    if (b == 0) return a;
    return gcd(b, a % b);
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !struct { part1: i64, part2: i64 } {
    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    var grid: std.ArrayList([]const u8) = .empty;
    defer grid.deinit(allocator);

    while (lines.next()) |line| {
        if (line.len > 0) {
            try grid.append(allocator, line);
        }
    }

    if (grid.items.len == 0 or grid.items[0].len == 0) {
        return error.EmptyGrid;
    }

    const rows = grid.items.len;
    const cols = grid.items[0].len;

    for (grid.items) |row| {
        if (row.len != cols) {
            return error.IrregularGrid;
        }
    }

    var antennas_by_freq = std.AutoHashMap(u8, std.ArrayList(struct { r: i64, c: i64 })).init(allocator);
    defer {
        var it = antennas_by_freq.valueIterator();
        while (it.next()) |list| {
            list.deinit(allocator);
        }
        antennas_by_freq.deinit();
    }

    for (0..rows) |r| {
        for (0..cols) |c| {
            const ch = grid.items[r][c];
            if (ch != '.') {
                if (!antennas_by_freq.contains(ch)) {
                    try antennas_by_freq.put(ch, .empty);
                }
                try antennas_by_freq.getPtr(ch).?.append(allocator, .{ .r = @intCast(r), .c = @intCast(c) });
            }
        }
    }

    var antinodes_part1 = std.AutoHashMap(struct { r: i64, c: i64 }, void).init(allocator);
    defer antinodes_part1.deinit();

    var antinodes_part2 = std.AutoHashMap(struct { r: i64, c: i64 }, void).init(allocator);
    defer antinodes_part2.deinit();

    var freq_it = antennas_by_freq.valueIterator();
    while (freq_it.next()) |positions| {
        const n = positions.items.len;
        for (0..n) |i| {
            for (i + 1..n) |j| {
                const r1 = positions.items[i].r;
                const c1 = positions.items[i].c;
                const r2 = positions.items[j].r;
                const c2 = positions.items[j].c;

                const dr = r2 - r1;
                const dc = c2 - c1;

                const antinode_r1 = r1 - dr;
                const antinode_c1 = c1 - dc;
                if (antinode_r1 >= 0 and antinode_r1 < rows and antinode_c1 >= 0 and antinode_c1 < cols) {
                    try antinodes_part1.put(.{ .r = antinode_r1, .c = antinode_c1 }, {});
                }

                const antinode_r2 = r2 + dr;
                const antinode_c2 = c2 + dc;
                if (antinode_r2 >= 0 and antinode_r2 < rows and antinode_c2 >= 0 and antinode_c2 < cols) {
                    try antinodes_part1.put(.{ .r = antinode_r2, .c = antinode_c2 }, {});
                }

                if (dr == 0 and dc == 0) continue;

                const g = gcd(@abs(dr), @abs(dc));
                const g_i64: i64 = @intCast(g);
                const step_r = @divTrunc(dr, g_i64);
                const step_c = @divTrunc(dc, g_i64);

                var r = r1;
                var c = c1;
                while (r >= 0 and r < rows and c >= 0 and c < cols) {
                    try antinodes_part2.put(.{ .r = r, .c = c }, {});
                    r += step_r;
                    c += step_c;
                }

                r = r1 - step_r;
                c = c1 - step_c;
                while (r >= 0 and r < rows and c >= 0 and c < cols) {
                    try antinodes_part2.put(.{ .r = r, .c = c }, {});
                    r -= step_r;
                    c -= step_c;
                }
            }
        }
    }

    return .{
        .part1 = @intCast(antinodes_part1.count()),
        .part2 = @intCast(antinodes_part2.count()),
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

test "example" {
    const example = @embedFile("examples/day8.txt");
    const answers = try solve(std.testing.allocator, example);
    try std.testing.expectEqual(answers.part1, 14);
    try std.testing.expectEqual(answers.part2, 34);
}

test "empty grid throws" {
    const result = solve(std.testing.allocator, "");
    try std.testing.expectError(error.EmptyGrid, result);
}

test "irregular grid throws" {
    const result = solve(std.testing.allocator, "...\n..\n...");
    try std.testing.expectError(error.IrregularGrid, result);
}
