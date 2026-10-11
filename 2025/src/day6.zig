const std = @import("std");

pub const Answers = struct { part1: i64, part2: i64 };

const Problem = struct { start: usize, end: usize };

/// Copies each line into a fixed-width buffer padded with spaces, so every
/// line can be indexed by column without bounds checks.
fn padLines(allocator: std.mem.Allocator, lines: []const []const u8, width: usize) ![]const []const u8 {
    const out = try allocator.alloc([]const u8, lines.len);
    for (lines, 0..) |line, i| {
        const buf = try allocator.alloc(u8, width);
        @memset(buf, ' ');
        @memcpy(buf[0..line.len], line);
        out[i] = buf;
    }
    return out;
}

fn isSeparatorColumn(lines: []const []const u8, col: usize) bool {
    for (lines) |line| {
        if (line[col] != ' ') return false;
    }
    return true;
}

/// Problems are contiguous runs of columns that are not entirely spaces.
fn findProblems(allocator: std.mem.Allocator, lines: []const []const u8, width: usize) ![]const Problem {
    var list: std.ArrayList(Problem) = .empty;
    var col: usize = 0;
    while (col < width) {
        if (isSeparatorColumn(lines, col)) {
            col += 1;
            continue;
        }
        const start = col;
        while (col < width and !isSeparatorColumn(lines, col)) col += 1;
        try list.append(allocator, .{ .start = start, .end = col });
    }
    return list.toOwnedSlice(allocator);
}

fn findOperator(op_row: []const u8, start: usize, end: usize) !u8 {
    for (op_row[start..end]) |ch| {
        if (ch == '+' or ch == '*') return ch;
    }
    return error.MissingOperator;
}

fn apply(op: u8, acc: i64, n: i64) i64 {
    return if (op == '+') acc + n else acc * n;
}

/// Part 1: each number row within a problem's columns is one operand.
fn solveRows(num_rows: []const []const u8, op: u8, start: usize, end: usize) !i64 {
    var acc: ?i64 = null;
    for (num_rows) |row| {
        const seg = std.mem.trim(u8, row[start..end], " ");
        if (seg.len == 0) continue;
        const n = try std.fmt.parseInt(i64, seg, 10);
        acc = if (acc) |a| apply(op, a, n) else n;
    }
    return acc orelse error.EmptyProblem;
}

/// Part 2: each column (right to left) within a problem's columns is one
/// operand, read top to bottom.
fn solveColumns(num_rows: []const []const u8, op: u8, start: usize, end: usize) !i64 {
    var acc: ?i64 = null;
    var col = end;
    while (col > start) {
        col -= 1;
        var buf: [32]u8 = undefined;
        var len: usize = 0;
        for (num_rows) |row| {
            const ch = row[col];
            if (ch == ' ') continue;
            if (!std.ascii.isDigit(ch)) return error.InvalidDigit;
            buf[len] = ch;
            len += 1;
        }
        if (len == 0) continue;
        const n = try std.fmt.parseInt(i64, buf[0..len], 10);
        acc = if (acc) |a| apply(op, a, n) else n;
    }
    return acc orelse error.EmptyProblem;
}

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var arena_state = std.heap.ArenaAllocator.init(allocator);
    defer arena_state.deinit();
    const arena = arena_state.allocator();

    var raw_lines: std.ArrayList([]const u8) = .empty;
    var it = std.mem.splitScalar(u8, input, '\n');
    while (it.next()) |line| {
        try raw_lines.append(arena, std.mem.trimEnd(u8, line, "\r"));
    }

    var count = raw_lines.items.len;
    while (count > 0 and raw_lines.items[count - 1].len == 0) count -= 1;
    if (count < 2) return error.InvalidInput;
    const trimmed_lines = raw_lines.items[0..count];

    var width: usize = 0;
    for (trimmed_lines) |line| width = @max(width, line.len);
    if (width == 0) return error.InvalidInput;

    const lines = try padLines(arena, trimmed_lines, width);
    const op_row = lines[lines.len - 1];
    const num_rows = lines[0 .. lines.len - 1];

    const problems = try findProblems(arena, lines, width);
    if (problems.len == 0) return error.InvalidInput;

    var part1: i64 = 0;
    var part2: i64 = 0;
    for (problems) |p| {
        const op = try findOperator(op_row, p.start, p.end);
        part1 += try solveRows(num_rows, op, p.start, p.end);
        part2 += try solveColumns(num_rows, op, p.start, p.end);
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

test "official example" {
    const input = @embedFile("examples/day6.txt");
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(i64, 4277556), answers.part1);
    try std.testing.expectEqual(@as(i64, 3263827), answers.part2);
}

test "single problem edge case" {
    const input = "5\n+\n";
    const answers = try solve(std.testing.allocator, input);
    try std.testing.expectEqual(@as(i64, 5), answers.part1);
    try std.testing.expectEqual(@as(i64, 5), answers.part2);
}

test "malformed input returns errors" {
    try std.testing.expectError(error.InvalidInput, solve(std.testing.allocator, ""));
    try std.testing.expectError(error.MissingOperator, solve(std.testing.allocator, "12\n34\n"));
    try std.testing.expectError(error.InvalidCharacter, solve(std.testing.allocator, "1x\n+ \n"));
}
