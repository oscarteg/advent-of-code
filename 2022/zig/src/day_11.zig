const std = @import("std");

const rounds_part1: usize = 20;
const rounds_part2: usize = 10000;
const relief_factor: u64 = 3;

const Operator = enum { add, multiply };

const Operation = struct {
    op: Operator,
    use_old: bool,
    operand: u64,

    fn apply(self: Operation, old: u64) u64 {
        const rhs = if (self.use_old) old else self.operand;
        return switch (self.op) {
            .add => old + rhs,
            .multiply => old * rhs,
        };
    }
};

const Monkey = struct {
    items: std.ArrayList(u64),
    op: Operation,
    divisor: u64,
    if_true: usize,
    if_false: usize,
    inspections: u64 = 0,
};

const Answers = struct {
    part1: u64,
    part2: u64,
};

fn parseSuffixNumbers(allocator: std.mem.Allocator, line: []const u8, prefix: []const u8) ![]u64 {
    const trimmed = std.mem.trim(u8, line, " \t\r");
    if (!std.mem.startsWith(u8, trimmed, prefix)) return error.InvalidMonkey;
    const rest = trimmed[prefix.len..];
    var numbers: std.ArrayList(u64) = .empty;
    errdefer numbers.deinit(allocator);
    var fields = std.mem.splitScalar(u8, rest, ',');
    while (fields.next()) |field| {
        const value = try std.fmt.parseInt(u64, std.mem.trim(u8, field, " \t\r"), 10);
        try numbers.append(allocator, value);
    }
    return numbers.toOwnedSlice(allocator);
}

fn parseSuffixNumber(line: []const u8, prefix: []const u8) !u64 {
    const trimmed = std.mem.trim(u8, line, " \t\r");
    if (!std.mem.startsWith(u8, trimmed, prefix)) return error.InvalidMonkey;
    return std.fmt.parseInt(u64, std.mem.trim(u8, trimmed[prefix.len..], " \t\r"), 10);
}

fn parseOperation(line: []const u8) !Operation {
    const prefix = "Operation: new = old ";
    const trimmed = std.mem.trim(u8, line, " \t\r");
    if (!std.mem.startsWith(u8, trimmed, prefix)) return error.InvalidMonkey;
    var fields = std.mem.tokenizeScalar(u8, trimmed[prefix.len..], ' ');
    const op_token = fields.next() orelse return error.InvalidMonkey;
    const op: Operator = if (std.mem.eql(u8, op_token, "+"))
        .add
    else if (std.mem.eql(u8, op_token, "*"))
        .multiply
    else
        return error.InvalidMonkey;
    const operand_token = fields.next() orelse return error.InvalidMonkey;
    if (std.mem.eql(u8, operand_token, "old")) {
        return .{ .op = op, .use_old = true, .operand = 0 };
    }
    const operand = try std.fmt.parseInt(u64, operand_token, 10);
    return .{ .op = op, .use_old = false, .operand = operand };
}

fn parseMonkey(allocator: std.mem.Allocator, block: []const u8) !Monkey {
    var lines: std.ArrayList([]const u8) = .empty;
    defer lines.deinit(allocator);
    var raw_lines = std.mem.splitScalar(u8, std.mem.trim(u8, block, " \t\r\n"), '\n');
    while (raw_lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len != 0) try lines.append(allocator, trimmed);
    }
    if (lines.items.len != 6) return error.InvalidMonkey;

    const items = try parseSuffixNumbers(allocator, lines.items[1], "Starting items:");
    var item_list: std.ArrayList(u64) = .empty;
    errdefer item_list.deinit(allocator);
    try item_list.appendSlice(allocator, items);
    allocator.free(items);

    const op = try parseOperation(lines.items[2]);
    const divisor = try parseSuffixNumber(lines.items[3], "Test: divisible by");
    const if_true = try parseSuffixNumber(lines.items[4], "If true: throw to monkey");
    const if_false = try parseSuffixNumber(lines.items[5], "If false: throw to monkey");

    return .{
        .items = item_list,
        .op = op,
        .divisor = divisor,
        .if_true = @intCast(if_true),
        .if_false = @intCast(if_false),
    };
}

fn parse(allocator: std.mem.Allocator, input: []const u8) !std.ArrayList(Monkey) {
    var monkeys: std.ArrayList(Monkey) = .empty;
    errdefer {
        for (monkeys.items) |*m| m.items.deinit(allocator);
        monkeys.deinit(allocator);
    }
    var blocks = std.mem.splitSequence(u8, std.mem.trim(u8, input, " \t\r\n"), "\n\n");
    while (blocks.next()) |block| {
        const trimmed = std.mem.trim(u8, block, " \t\r\n");
        if (trimmed.len == 0) continue;
        const monkey = try parseMonkey(allocator, trimmed);
        try monkeys.append(allocator, monkey);
    }
    if (monkeys.items.len == 0) return error.InvalidMonkey;
    for (monkeys.items) |m| {
        if (m.divisor == 0) return error.InvalidMonkey;
        if (m.if_true >= monkeys.items.len or m.if_false >= monkeys.items.len) return error.InvalidMonkey;
    }
    return monkeys;
}

fn deinitMonkeys(allocator: std.mem.Allocator, monkeys: *std.ArrayList(Monkey)) void {
    for (monkeys.items) |*m| m.items.deinit(allocator);
    monkeys.deinit(allocator);
}

fn run(allocator: std.mem.Allocator, monkeys: []Monkey, rounds: usize, relief: *const fn (u64) u64) !u64 {
    var round: usize = 0;
    while (round < rounds) : (round += 1) {
        for (monkeys) |*m| {
            while (m.items.items.len > 0) {
                const item = m.items.orderedRemove(0);
                m.inspections += 1;
                const worry = relief(m.op.apply(item));
                const target = if (worry % m.divisor == 0) m.if_true else m.if_false;
                try monkeys[target].items.append(allocator, worry);
            }
        }
    }
    var counts = try allocator.alloc(u64, monkeys.len);
    defer allocator.free(counts);
    for (monkeys, 0..) |m, i| counts[i] = m.inspections;
    std.mem.sort(u64, counts, {}, std.sort.desc(u64));
    return counts[0] * counts[1];
}

fn reliefPart1(worry: u64) u64 {
    return worry / relief_factor;
}

var part2_modulus: u64 = 1;
fn reliefPart2(worry: u64) u64 {
    return worry % part2_modulus;
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var part1_monkeys = try parse(allocator, input);
    defer deinitMonkeys(allocator, &part1_monkeys);
    const part1 = try run(allocator, part1_monkeys.items, rounds_part1, reliefPart1);

    var part2_monkeys = try parse(allocator, input);
    defer deinitMonkeys(allocator, &part2_monkeys);
    var modulus: u64 = 1;
    for (part2_monkeys.items) |m| modulus *= m.divisor;
    part2_modulus = modulus;
    const part2 = try run(allocator, part2_monkeys.items, rounds_part2, reliefPart2);

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
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day_11.txt");
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(@as(u64, 10605), answers.part1);
    try std.testing.expectEqual(@as(u64, 2713310158), answers.part2);
}

test "two monkeys swap items each round" {
    const allocator = std.testing.allocator;
    const input = "Monkey 0:\n  Starting items: 1, 2\n  Operation: new = old + 1\n  Test: divisible by 2\n    If true: throw to monkey 1\n    If false: throw to monkey 1\n\nMonkey 1:\n  Starting items: 3\n  Operation: new = old + 1\n  Test: divisible by 2\n    If true: throw to monkey 0\n    If false: throw to monkey 0\n";
    var monkeys = try parse(allocator, input);
    defer deinitMonkeys(allocator, &monkeys);
    const answer = try run(allocator, monkeys.items, 1, reliefPart1);
    // Monkey 0 inspects 2 items, monkey 1 inspects the 3 items it now holds: 2 * 3 = 6.
    try std.testing.expectEqual(@as(u64, 6), answer);
}

test "malformed input missing fields is error" {
    const allocator = std.testing.allocator;
    const input = "Monkey 0:\n  Starting items: 1\n  Operation: new = old + 1\n";
    try std.testing.expectError(error.InvalidMonkey, parse(allocator, input));
}

test "malformed input throw target out of range is error" {
    const allocator = std.testing.allocator;
    const input = "Monkey 0:\n  Starting items: 1\n  Operation: new = old + 1\n  Test: divisible by 2\n    If true: throw to monkey 1\n    If false: throw to monkey 0\n";
    try std.testing.expectError(error.InvalidMonkey, parse(allocator, input));
}
