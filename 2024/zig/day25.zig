const std = @import("std");

const pin_space: usize = 5;
const rows: usize = 7;
const columns: usize = 5;

const Answers = struct {
    part1: i64,
    part2: i64,
};

const Schematic = struct {
    is_lock: bool,
    heights: [columns]u8,
};

fn parseSchematic(text: []const u8) !Schematic {
    var lines = std.mem.tokenizeScalar(u8, text, '\n');
    var grid: [rows][]const u8 = undefined;
    var count: usize = 0;
    while (lines.next()) |line| {
        if (count >= rows) return error.InvalidSchematic;
        if (line.len != columns) return error.InvalidSchematic;
        grid[count] = line;
        count += 1;
    }
    if (count != rows) return error.InvalidSchematic;

    const top_full = std.mem.indexOfScalar(u8, grid[0], '.') == null;
    const bottom_full = std.mem.indexOfScalar(u8, grid[rows - 1], '.') == null;
    if (top_full == bottom_full) return error.InvalidSchematic;

    var heights: [columns]u8 = .{0} ** columns;
    for (grid) |row| {
        for (row, 0..) |cell, column| {
            if (cell == '#') heights[column] += 1;
        }
    }
    for (&heights) |*height| height.* -= 1;

    return Schematic{ .is_lock = top_full, .heights = heights };
}

fn fits(lock: [columns]u8, key: [columns]u8) bool {
    for (0..columns) |column| {
        if (@as(usize, lock[column]) + @as(usize, key[column]) > pin_space) return false;
    }
    return true;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var locks: std.ArrayList([columns]u8) = .empty;
    defer locks.deinit(allocator);
    var keys: std.ArrayList([columns]u8) = .empty;
    defer keys.deinit(allocator);

    var blocks = std.mem.splitSequence(u8, std.mem.trim(u8, input, "\n"), "\n\n");
    while (blocks.next()) |block| {
        const schematic = try parseSchematic(block);
        if (schematic.is_lock) {
            try locks.append(allocator, schematic.heights);
        } else {
            try keys.append(allocator, schematic.heights);
        }
    }
    if (locks.items.len == 0 or keys.items.len == 0) return error.InvalidSchematic;

    var fitting_pairs: i64 = 0;
    for (locks.items) |lock| {
        for (keys.items) |key| {
            if (fits(lock, key)) fitting_pairs += 1;
        }
    }
    return Answers{ .part1 = fitting_pairs, .part2 = 0 };
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
    try stdout.print("Part 1: {d}\n", .{answers.part1});
    try stdout.flush();
}

test "official example" {
    const input = @embedFile("examples/day25.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(i64, 3), answers.part1);
}

test "no pairs fit" {
    const input =
        "#####\n#####\n#####\n#####\n#####\n#####\n.....\n" ++
        "\n" ++
        ".....\n#####\n#####\n#####\n#####\n#####\n#####\n";
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(i64, 0), answers.part1);
}

test "malformed input rejected" {
    try std.testing.expectError(error.InvalidSchematic, solve(std.testing.allocator, "not a schematic"));
    try std.testing.expectError(error.InvalidSchematic, solve(std.testing.allocator, "#####\n.####\n.....\n"));
}
