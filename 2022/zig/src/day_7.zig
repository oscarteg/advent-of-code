const std = @import("std");

const disk_space: u64 = 70_000_000;
const required_unused: u64 = 30_000_000;
const small_directory_limit: u64 = 100_000;

const Answers = struct {
    part1: u64,
    part2: u64,
};

fn pathKey(allocator: std.mem.Allocator, stack: []const []const u8) ![]u8 {
    if (stack.len == 0) return allocator.dupe(u8, "/");
    const joined = try std.mem.join(allocator, "/", stack);
    defer allocator.free(joined);
    return std.fmt.allocPrint(allocator, "/{s}", .{joined});
}

fn directorySizes(allocator: std.mem.Allocator, input: []const u8) !std.StringHashMap(u64) {
    var sizes = std.StringHashMap(u64).init(allocator);
    errdefer {
        var it = sizes.keyIterator();
        while (it.next()) |key| allocator.free(key.*);
        sizes.deinit();
    }

    var stack: std.ArrayList([]const u8) = .empty;
    defer stack.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw_line| {
        const line = std.mem.trim(u8, raw_line, " \r\t");
        if (line.len == 0) continue;
        if (std.mem.eql(u8, line, "$ ls")) continue;

        if (std.mem.startsWith(u8, line, "$ cd ")) {
            const target = line[5..];
            if (std.mem.eql(u8, target, "/")) {
                stack.clearRetainingCapacity();
            } else if (std.mem.eql(u8, target, "..")) {
                if (stack.items.len == 0) return error.CdAboveRoot;
                _ = stack.pop();
            } else {
                try stack.append(allocator, target);
            }
            const key = try pathKey(allocator, stack.items);
            if (sizes.contains(key)) {
                allocator.free(key);
            } else {
                try sizes.put(key, 0);
            }
            continue;
        }

        if (std.mem.startsWith(u8, line, "dir ")) {
            if (line.len == 4) return error.MissingDirectoryName;
            continue;
        }

        const space = std.mem.indexOfScalar(u8, line, ' ') orelse return error.MalformedListingLine;
        const size_str = line[0..space];
        const name = line[space + 1 ..];
        if (name.len == 0) return error.MissingFileName;
        const size = try std.fmt.parseInt(u64, size_str, 10);

        var depth: usize = 0;
        while (depth <= stack.items.len) : (depth += 1) {
            const key = try pathKey(allocator, stack.items[0..depth]);
            if (sizes.getPtr(key)) |existing| {
                existing.* += size;
                allocator.free(key);
            } else {
                try sizes.put(key, size);
            }
        }
    }

    if (sizes.count() == 0) return error.NoDirectoriesFound;
    return sizes;
}

fn solveWith(
    allocator: std.mem.Allocator,
    input: []const u8,
    disk: u64,
    unused_needed: u64,
    small_limit: u64,
) !Answers {
    var sizes = try directorySizes(allocator, input);
    defer {
        var it = sizes.keyIterator();
        while (it.next()) |key| allocator.free(key.*);
        sizes.deinit();
    }

    var part1: u64 = 0;
    var part2: u64 = std.math.maxInt(u64);

    const used = sizes.get("/") orelse return error.MissingRootDirectory;
    if (used > disk) return error.UsedSpaceExceedsDisk;
    const unused = disk - used;
    const need_to_free = if (unused_needed > unused) unused_needed - unused else 0;

    var it = sizes.valueIterator();
    while (it.next()) |size| {
        if (size.* <= small_limit) part1 += size.*;
        if (size.* >= need_to_free and size.* < part2) part2 = size.*;
    }
    if (part2 == std.math.maxInt(u64)) return error.NoDirectoryLargeEnoughToDelete;

    return .{ .part1 = part1, .part2 = part2 };
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    return solveWith(allocator, input, disk_space, required_unused, small_directory_limit);
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
    const input = @embedFile("examples/day_7.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u64, 95437), answers.part1);
    try std.testing.expectEqual(@as(u64, 24933642), answers.part2);
}

test "empty directories count as zero" {
    const input = "$ cd /\n$ ls\ndir a\n100000 big.txt\n$ cd a\n$ ls\n";
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(u64, 100000), answers.part1);
}

test "cd above root is rejected" {
    const input = "$ cd /\n$ cd ..\n";
    try std.testing.expectError(error.CdAboveRoot, solve(std.testing.allocator, input));
}

test "malformed listing line is rejected" {
    const input = "$ cd /\n$ ls\nnotanumber\n";
    try std.testing.expectError(error.MalformedListingLine, solve(std.testing.allocator, input));
}
