const std = @import("std");

// Directional-keypad robots between the door's numeric keypad and the human.
const part1_robots = 2;
const part2_robots = 25;

// The longest move crosses the numeric keypad in five steps, plus the A that presses the key.
const max_route_len = 6;

const Answers = struct {
    part1: u64,
    part2: u64,
};

const Position = struct { row: usize, column: usize };

const Order = enum { horizontal_first, vertical_first };

const Step = struct { button: u8, count: usize };

// Keypad rows, with ' ' marking the gap no robot arm may cross.
const numeric_rows = [_][]const u8{ "789", "456", "123", " 0A" };
const directional_rows = [_][]const u8{ " ^A", "<v>" };

const Pad = enum {
    numeric,
    directional,

    fn rows(self: Pad) []const []const u8 {
        return switch (self) {
            .numeric => &numeric_rows,
            .directional => &directional_rows,
        };
    }

    fn locate(self: Pad, key: u8) !Position {
        for (self.rows(), 0..) |row, index| {
            if (std.mem.indexOfScalar(u8, row, key)) |column| {
                return .{ .row = index, .column = column };
            }
        }
        return error.KeyNotOnKeypad;
    }

    fn hasKey(self: Pad, at: Position) bool {
        const pad_rows = self.rows();
        if (at.row >= pad_rows.len) return false;
        const row = pad_rows[at.row];
        return at.column < row.len and row[at.column] != ' ';
    }
};

fn stepsBetween(from: usize, to: usize) usize {
    return if (from > to) from - to else to - from;
}

// The shortest route for one move, ending with the A that presses the key. Only two orders can
// be shortest - every horizontal step first, or every vertical step first - and an order is
// impossible when its turning point lands on the keypad's gap.
fn routeFor(pad: Pad, start: Position, target: Position, order: Order, buffer: []u8) ?[]const u8 {
    const corner: Position = switch (order) {
        .horizontal_first => .{ .row = start.row, .column = target.column },
        .vertical_first => .{ .row = target.row, .column = start.column },
    };
    if (!pad.hasKey(corner)) return null;

    const horizontal = Step{
        .button = if (target.column > start.column) '>' else '<',
        .count = stepsBetween(start.column, target.column),
    };
    const vertical = Step{
        .button = if (target.row > start.row) 'v' else '^',
        .count = stepsBetween(start.row, target.row),
    };
    const steps = switch (order) {
        .horizontal_first => [2]Step{ horizontal, vertical },
        .vertical_first => [2]Step{ vertical, horizontal },
    };

    var len: usize = 0;
    for (steps) |step| {
        @memset(buffer[len..][0..step.count], step.button);
        len += step.count;
    }
    buffer[len] = 'A';
    return buffer[0 .. len + 1];
}

const MoveKey = struct {
    pad: Pad,
    from: u8,
    to: u8,
    robots: u8,
};

// Spelled out because typeOn and moveAndPress call each other.
const PressError = error{ KeyNotOnKeypad, GapBlocksEveryRoute, OutOfMemory };

// Human button presses, memoised per keypad, key pair and robot count.
const Solver = struct {
    cache: std.AutoHashMap(MoveKey, u64),

    fn init(allocator: std.mem.Allocator) Solver {
        return .{ .cache = std.AutoHashMap(MoveKey, u64).init(allocator) };
    }

    fn deinit(self: *Solver) void {
        self.cache.deinit();
    }

    fn typeOn(self: *Solver, pad: Pad, sequence: []const u8, robots: u8) PressError!u64 {
        var presses: u64 = 0;
        var from: u8 = 'A'; // every arm starts on its own A button
        for (sequence) |to| {
            presses += try self.moveAndPress(pad, from, to, robots);
            from = to;
        }
        return presses;
    }

    fn moveAndPress(self: *Solver, pad: Pad, from: u8, to: u8, robots: u8) PressError!u64 {
        const key = MoveKey{ .pad = pad, .from = from, .to = to, .robots = robots };
        if (self.cache.get(key)) |cached| return cached;

        const start = try pad.locate(from);
        const target = try pad.locate(to);

        var cheapest: ?u64 = null;
        var buffer: [max_route_len]u8 = undefined;
        for ([_]Order{ .horizontal_first, .vertical_first }) |order| {
            const route = routeFor(pad, start, target, order, &buffer) orelse continue;
            // The human types the last keypad itself, one press per button of the route.
            const presses: u64 = if (robots == 0)
                route.len
            else
                try self.typeOn(.directional, route, robots - 1);
            if (cheapest == null or presses < cheapest.?) cheapest = presses;
        }
        const presses = cheapest orelse return error.GapBlocksEveryRoute;

        try self.cache.put(key, presses);
        return presses;
    }
};

const Code = struct {
    buttons: []const u8,
    numeric_part: u64,

    fn parse(line: []const u8) !Code {
        if (line.len == 0 or line[line.len - 1] != 'A') return error.InvalidCode;
        const digits = line[0 .. line.len - 1];
        if (digits.len == 0) return error.InvalidCode;
        for (digits) |digit| {
            if (!std.ascii.isDigit(digit)) return error.InvalidCode;
        }
        return .{ .buttons = line, .numeric_part = try std.fmt.parseInt(u64, digits, 10) };
    }
};

// Sums the complexities of every code: the presses the human needs for the shortest sequence
// that types it, times the code's numeric part.
fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var solver = Solver.init(allocator);
    defer solver.deinit();

    var answers = Answers{ .part1 = 0, .part2 = 0 };
    var codes: usize = 0;
    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw| {
        const line = std.mem.trim(u8, raw, " \r");
        if (line.len == 0) continue;

        const code = try Code.parse(line);
        answers.part1 += code.numeric_part * try solver.typeOn(.numeric, code.buttons, part1_robots);
        answers.part2 += code.numeric_part * try solver.typeOn(.numeric, code.buttons, part2_robots);
        codes += 1;
    }
    if (codes == 0) return error.NoCodes;
    return answers;
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
    const answers = try solve(std.testing.allocator, @embedFile("examples/day21.txt"));
    try std.testing.expectEqual(@as(u64, 126384), answers.part1);
    try std.testing.expectEqual(@as(u64, 154115708116294), answers.part2);
}

test "published sequence lengths" {
    const cases = [_]struct { input: []const u8, part1: u64 }{
        .{ .input = "029A\n", .part1 = 68 * 29 },
        .{ .input = "980A\n", .part1 = 60 * 980 },
        .{ .input = "179A\n", .part1 = 68 * 179 },
        .{ .input = "456A\n", .part1 = 64 * 456 },
        .{ .input = "379A\n", .part1 = 64 * 379 },
    };
    for (cases) |expected| {
        const answers = try solve(std.testing.allocator, expected.input);
        try std.testing.expectEqual(expected.part1, answers.part1);
    }
}

test "codes are independent" {
    const both = try solve(std.testing.allocator, "029A\n980A\n");
    const first = try solve(std.testing.allocator, "029A\n");
    const second = try solve(std.testing.allocator, "980A\n");
    try std.testing.expectEqual(first.part1 + second.part1, both.part1);
    try std.testing.expectEqual(first.part2 + second.part2, both.part2);
}

test "leading zeroes leave a code without complexity" {
    const answers = try solve(std.testing.allocator, "000A\n");
    try std.testing.expectEqual(@as(u64, 0), answers.part1);
    try std.testing.expectEqual(@as(u64, 0), answers.part2);
}

test "carriage returns and blank lines are ignored" {
    const answers = try solve(std.testing.allocator, "029A\r\n\r\n980A\r\n");
    const expected = try solve(std.testing.allocator, "029A\n980A\n");
    try std.testing.expectEqual(expected.part1, answers.part1);
    try std.testing.expectEqual(expected.part2, answers.part2);
}

test "malformed input returns an error" {
    const allocator = std.testing.allocator;
    try std.testing.expectError(error.NoCodes, solve(allocator, ""));
    try std.testing.expectError(error.NoCodes, solve(allocator, "   \n"));
    try std.testing.expectError(error.InvalidCode, solve(allocator, "029\n"));
    try std.testing.expectError(error.InvalidCode, solve(allocator, "02XA\n"));
    try std.testing.expectError(error.InvalidCode, solve(allocator, "A\n"));
    try std.testing.expectError(error.InvalidCode, solve(allocator, "-29A\n"));
    try std.testing.expectError(error.InvalidCode, solve(allocator, "029a\n"));
}
