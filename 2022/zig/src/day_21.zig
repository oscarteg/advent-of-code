const std = @import("std");
const testing = std.testing;

const root_name = "root";
const human_name = "humn";

const Answers = struct {
    part1: i64,
    part2: i64,
};

const ArithmeticError = error{ Overflow, DivisionByZero };
const EvalError = ArithmeticError || error{ UnknownMonkey, CyclicRiddle };

const Operator = enum {
    add,
    subtract,
    multiply,
    divide,

    fn parse(token: []const u8) error{UnknownOperator}!Operator {
        if (token.len == 1) {
            switch (token[0]) {
                '+' => return .add,
                '-' => return .subtract,
                '*' => return .multiply,
                '/' => return .divide,
                else => {},
            }
        }
        return error.UnknownOperator;
    }

    fn apply(self: Operator, left: i64, right: i64) ArithmeticError!i64 {
        return switch (self) {
            .add => std.math.add(i64, left, right),
            .subtract => std.math.sub(i64, left, right),
            .multiply => std.math.mul(i64, left, right),
            .divide => std.math.divTrunc(i64, left, right),
        };
    }
};

const Job = union(enum) {
    number: i64,
    operation: struct {
        left: []const u8,
        operator: Operator,
        right: []const u8,
    },
};

const Monkeys = std.StringHashMap(Job);
const Listeners = std.StringHashMap([]const u8);

fn parse(allocator: std.mem.Allocator, input: []const u8) !Monkeys {
    var monkeys = Monkeys.init(allocator);
    errdefer monkeys.deinit();

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw_line| {
        const line = std.mem.trim(u8, raw_line, " \t\r");
        if (line.len == 0) continue;

        const colon = std.mem.indexOfScalar(u8, line, ':') orelse return error.MissingColon;
        const name = std.mem.trim(u8, line[0..colon], " \t");
        if (name.len == 0) return error.MissingMonkeyName;

        var fields = std.mem.tokenizeAny(u8, line[colon + 1 ..], " \t");
        const first = fields.next() orelse return error.MissingJob;
        const job: Job = if (fields.next()) |operator_token| .{ .operation = .{
            .left = first,
            .operator = try Operator.parse(operator_token),
            .right = fields.next() orelse return error.MalformedJob,
        } } else .{ .number = try std.fmt.parseInt(i64, first, 10) };
        if (fields.next() != null) return error.MalformedJob;

        if (try monkeys.fetchPut(name, job) != null) return error.DuplicateMonkey;
    }

    if (!monkeys.contains(root_name)) return error.MissingRoot;
    return monkeys;
}

/// Yells `name`'s number. `budget` bounds the depth by one step per monkey, so
/// a riddle that refers back to itself fails instead of looping forever.
fn evaluate(monkeys: *const Monkeys, name: []const u8, budget: usize) EvalError!i64 {
    if (budget == 0) return error.CyclicRiddle;
    const job = monkeys.get(name) orelse return error.UnknownMonkey;
    return switch (job) {
        .number => |number| number,
        .operation => |operation| operation.operator.apply(
            try evaluate(monkeys, operation.left, budget - 1),
            try evaluate(monkeys, operation.right, budget - 1),
        ),
    };
}

/// Maps every monkey to the monkey waiting for its number. Part 2 inverts a
/// single path, so a monkey heard by two monkeys is rejected.
fn collectListeners(allocator: std.mem.Allocator, monkeys: *const Monkeys) !Listeners {
    var listeners = Listeners.init(allocator);
    errdefer listeners.deinit();

    var entries = monkeys.iterator();
    while (entries.next()) |entry| {
        switch (entry.value_ptr.*) {
            .number => {},
            .operation => |operation| {
                for ([_][]const u8{ operation.left, operation.right }) |operand| {
                    if (try listeners.fetchPut(operand, entry.key_ptr.*) != null) {
                        return error.MonkeyHeardTwice;
                    }
                }
            },
        }
    }
    return listeners;
}

fn pathToHuman(allocator: std.mem.Allocator, monkeys: *const Monkeys) !std.ArrayList([]const u8) {
    var listeners = try collectListeners(allocator, monkeys);
    defer listeners.deinit();
    if (!monkeys.contains(human_name)) return error.MissingHuman;

    var path: std.ArrayList([]const u8) = .empty;
    errdefer path.deinit(allocator);
    try path.append(allocator, human_name);
    var current: []const u8 = human_name;
    while (!std.mem.eql(u8, current, root_name)) {
        current = listeners.get(current) orelse return error.HumanNotHeardByRoot;
        try path.append(allocator, current);
        if (path.items.len > monkeys.count()) return error.CyclicRiddle;
    }
    std.mem.reverse([]const u8, path.items);
    return path;
}

/// Solves `target = left <operator> right` for the operand on `humn`'s side.
/// Divisions have to be exact, otherwise the riddle has no integer answer.
fn invert(operator: Operator, target: i64, known: i64, human_on_left: bool) !i64 {
    return switch (operator) {
        .add => std.math.sub(i64, target, known),
        .subtract => if (human_on_left)
            std.math.add(i64, target, known)
        else
            std.math.sub(i64, known, target),
        .multiply => std.math.divExact(i64, target, known),
        .divide => if (human_on_left)
            std.math.mul(i64, target, known)
        else
            std.math.divExact(i64, known, target),
    };
}

/// Walks from `root` down to `humn`, inverting every operation on the way: the
/// branch without `humn` evaluates to a known number, which pins down what the
/// branch holding `humn` has to yell.
fn humanValue(allocator: std.mem.Allocator, monkeys: *const Monkeys) !i64 {
    var path = try pathToHuman(allocator, monkeys);
    defer path.deinit(allocator);

    const budget = monkeys.count();
    var target: i64 = 0;
    for (path.items[0 .. path.items.len - 1], path.items[1..]) |parent, child| {
        const job = monkeys.get(parent) orelse return error.UnknownMonkey;
        const operation = switch (job) {
            .number => return error.MalformedJob,
            .operation => |operation| operation,
        };
        const human_on_left = std.mem.eql(u8, operation.left, child);
        const other = if (human_on_left) operation.right else operation.left;
        const known = try evaluate(monkeys, other, budget);

        // `root` compares its operands instead of combining them.
        target = if (std.mem.eql(u8, parent, root_name))
            known
        else
            try invert(operation.operator, target, known, human_on_left);
    }
    return target;
}

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var monkeys = try parse(allocator, input);
    defer monkeys.deinit();
    return .{
        .part1 = try evaluate(&monkeys, root_name, monkeys.count()),
        .part2 = try humanValue(allocator, &monkeys),
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
    const answers = try solve(testing.allocator, @embedFile("examples/day_21.txt"));
    try testing.expectEqual(@as(i64, 152), answers.part1);
    try testing.expectEqual(@as(i64, 301), answers.part2);
}

test "human is a direct operand of root" {
    const answers = try solve(testing.allocator, "root: humn + other\nother: 7\nhumn: 3\n");
    try testing.expectEqual(@as(i64, 10), answers.part1);
    try testing.expectEqual(@as(i64, 7), answers.part2);
}

test "inverts operations with the human on the right" {
    const input =
        \\root: target + chain
        \\target: 10
        \\chain: base - inner
        \\base: 30
        \\inner: fixed / humn
        \\fixed: 40
        \\humn: 1
        \\
    ;
    const answers = try solve(testing.allocator, input);
    try testing.expectEqual(@as(i64, 0), answers.part1);
    try testing.expectEqual(@as(i64, 2), answers.part2);
}

test "inverts multiplication and division with the human on the left" {
    const input = "root: a + b\na: 36\nb: c * three\nthree: 3\nc: humn / two\ntwo: 2\nhumn: 0\n";
    const answers = try solve(testing.allocator, input);
    try testing.expectEqual(@as(i64, 36), answers.part1);
    try testing.expectEqual(@as(i64, 24), answers.part2);
}

test "ignores blank lines and surrounding whitespace" {
    const answers = try solve(testing.allocator, "\n  root: humn + other  \n\nother: 7\n humn: 3\n\n");
    try testing.expectEqual(@as(i64, 10), answers.part1);
    try testing.expectEqual(@as(i64, 7), answers.part2);
}

test "malformed input returns an error" {
    try testing.expectError(error.MissingRoot, solve(testing.allocator, ""));
    try testing.expectError(error.MissingColon, solve(testing.allocator, "root 1\n"));
    try testing.expectError(error.MissingMonkeyName, solve(testing.allocator, ": 1\n"));
    try testing.expectError(error.MissingJob, solve(testing.allocator, "root:\n"));
    try testing.expectError(error.MalformedJob, solve(testing.allocator, "root: a + b + c\n"));
    try testing.expectError(error.UnknownOperator, solve(testing.allocator, "root: a ? b\n"));
    try testing.expectError(error.InvalidCharacter, solve(testing.allocator, "root: humn + other\nother: twelve\nhumn: 1\n"));
    try testing.expectError(error.DuplicateMonkey, solve(testing.allocator, "root: humn + other\nother: 7\nother: 8\nhumn: 1\n"));
    try testing.expectError(error.UnknownMonkey, solve(testing.allocator, "root: humn + missing\nhumn: 1\n"));
    try testing.expectError(error.MissingHuman, solve(testing.allocator, "root: a + b\na: 1\nb: 2\n"));
}

test "cycles return an error instead of looping" {
    try testing.expectError(error.CyclicRiddle, solve(testing.allocator, "root: a + b\na: 1\nb: b - a\nhumn: 1\n"));
}

test "riddle without an integer answer returns an error" {
    try testing.expectError(
        error.UnexpectedRemainder,
        solve(testing.allocator, "root: odd + scaled\nodd: 7\nscaled: humn * two\ntwo: 2\nhumn: 0\n"),
    );
}

test "division by zero returns an error" {
    try testing.expectError(
        error.DivisionByZero,
        solve(testing.allocator, "root: a + b\na: 1\nb: a / zero\nzero: 0\nhumn: 1\n"),
    );
}
