const std = @import("std");

const whitespace = " \t\r";

const Machine = struct {
    a: u64,
    b: u64,
    c: u64,
    program: []const u3,
};

const Answers = struct {
    part1: []const u8,
    part2: u64,
};

const Registers = struct {
    a: u64,
    b: u64,
    c: u64,

    fn combo(self: Registers, operand: u3) !u64 {
        return switch (operand) {
            0, 1, 2, 3 => operand,
            4 => self.a,
            5 => self.b,
            6 => self.c,
            7 => error.ReservedComboOperand,
        };
    }
};

fn registerValue(line: ?[]const u8, label: []const u8) !u64 {
    const text = std.mem.trim(u8, line orelse return error.MissingRegister, whitespace);
    if (!std.mem.startsWith(u8, text, label)) return error.MissingRegister;
    return std.fmt.parseInt(u64, std.mem.trim(u8, text[label.len..], whitespace), 10);
}

fn parse(allocator: std.mem.Allocator, input: []const u8) !Machine {
    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    const a = try registerValue(lines.next(), "Register A:");
    const b = try registerValue(lines.next(), "Register B:");
    const c = try registerValue(lines.next(), "Register C:");

    const label = "Program:";
    const line = std.mem.trim(u8, lines.next() orelse return error.MissingProgram, whitespace);
    if (!std.mem.startsWith(u8, line, label)) return error.MissingProgram;

    var program: std.ArrayList(u3) = .empty;
    errdefer program.deinit(allocator);
    var values = std.mem.splitScalar(u8, line[label.len..], ',');
    while (values.next()) |value| {
        try program.append(allocator, try std.fmt.parseInt(u3, std.mem.trim(u8, value, whitespace), 10));
    }

    return .{ .a = a, .b = b, .c = c, .program = try program.toOwnedSlice(allocator) };
}

fn run(allocator: std.mem.Allocator, machine: Machine, initial_a: u64) ![]u3 {
    var registers = Registers{ .a = initial_a, .b = machine.b, .c = machine.c };
    var output: std.ArrayList(u3) = .empty;
    errdefer output.deinit(allocator);

    var ip: usize = 0;
    while (ip + 1 < machine.program.len) {
        const opcode = machine.program[ip];
        const operand = machine.program[ip + 1];
        var next = ip + 2;
        switch (opcode) {
            0 => registers.a = std.math.shr(u64, registers.a, try registers.combo(operand)),
            1 => registers.b ^= operand,
            2 => registers.b = try registers.combo(operand) % 8,
            3 => if (registers.a != 0) {
                next = operand;
            },
            4 => registers.b ^= registers.c,
            5 => try output.append(allocator, @intCast(try registers.combo(operand) % 8)),
            6 => registers.b = std.math.shr(u64, registers.a, try registers.combo(operand)),
            7 => registers.c = std.math.shr(u64, registers.a, try registers.combo(operand)),
        }
        ip = next;
    }

    return output.toOwnedSlice(allocator);
}

fn joinValues(allocator: std.mem.Allocator, values: []const u3) ![]u8 {
    var text: std.ArrayList(u8) = .empty;
    errdefer text.deinit(allocator);
    for (values, 0..) |value, index| {
        if (index > 0) try text.append(allocator, ',');
        try text.append(allocator, '0' + @as(u8, value));
    }
    return text.toOwnedSlice(allocator);
}

// A self-replicating program emits one value per loop and then drops the low
// three bits of A, so A is built three bits at a time from the last value
// backwards, keeping every candidate whose run reproduces the program's tail.
fn findSelfReplicatingA(allocator: std.mem.Allocator, machine: Machine) !u64 {
    var candidates: std.ArrayList(u64) = .empty;
    defer candidates.deinit(allocator);
    try candidates.append(allocator, 0);

    var index = machine.program.len;
    while (index > 0) {
        index -= 1;
        const tail = machine.program[index..];

        var next: std.ArrayList(u64) = .empty;
        errdefer next.deinit(allocator);
        for (candidates.items) |prefix| {
            for (0..8) |digit| {
                const candidate = prefix * 8 + @as(u64, @intCast(digit));
                const output = try run(allocator, machine, candidate);
                defer allocator.free(output);
                if (std.mem.eql(u3, output, tail)) try next.append(allocator, candidate);
            }
        }

        candidates.deinit(allocator);
        candidates = next;
    }

    return if (candidates.items.len == 0) 0 else std.mem.min(u64, candidates.items);
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const machine = try parse(allocator, input);
    defer allocator.free(machine.program);

    const output = try run(allocator, machine, machine.a);
    defer allocator.free(output);

    const part2 = try findSelfReplicatingA(allocator, machine);
    return .{ .part1 = try joinValues(allocator, output), .part2 = part2 };
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

    try stdout.print("Part 1: {s}\nPart 2: {d}\n", .{ answers.part1, answers.part2 });
    try stdout.flush();
}

test "official example" {
    const answers = try solve(std.testing.allocator, @embedFile("examples/day17.txt"));
    defer std.testing.allocator.free(answers.part1);
    try std.testing.expectEqualStrings("4,6,3,5,6,3,5,2,1,0", answers.part1);
    try std.testing.expectEqual(@as(u64, 0), answers.part2);
}

test "program that prints itself" {
    const answers = try solve(std.testing.allocator, "Register A: 2024\nRegister B: 0\nRegister C: 0\n\nProgram: 0,3,5,4,3,0\n");
    defer std.testing.allocator.free(answers.part1);
    try std.testing.expectEqualStrings("5,7,3,0", answers.part1);
    try std.testing.expectEqual(@as(u64, 117440), answers.part2);
}

test "instruction examples" {
    const literals = try solve(std.testing.allocator, "Register A: 10\nRegister B: 0\nRegister C: 0\n\nProgram: 5,0,5,1,5,4\n");
    defer std.testing.allocator.free(literals.part1);
    try std.testing.expectEqualStrings("0,1,2", literals.part1);

    const shifts = try solve(std.testing.allocator, "Register A: 2024\nRegister B: 0\nRegister C: 0\n\nProgram: 0,1,5,4,3,0\n");
    defer std.testing.allocator.free(shifts.part1);
    try std.testing.expectEqualStrings("4,2,5,6,7,7,7,7,3,1,0", shifts.part1);
}

test "program without an out instruction prints nothing" {
    const answers = try solve(std.testing.allocator, "Register A: 0\nRegister B: 29\nRegister C: 0\n\nProgram: 1,7\n");
    defer std.testing.allocator.free(answers.part1);
    try std.testing.expectEqualStrings("", answers.part1);
    try std.testing.expectEqual(@as(u64, 0), answers.part2);
}

test "malformed input" {
    const allocator = std.testing.allocator;
    try std.testing.expectError(error.MissingRegister, solve(allocator, ""));
    try std.testing.expectError(error.MissingProgram, solve(allocator, "Register A: 1\nRegister B: 0\nRegister C: 0\n"));
    try std.testing.expectError(error.InvalidCharacter, solve(allocator, "Register A: x\nRegister B: 0\nRegister C: 0\n\nProgram: 0,1\n"));
    try std.testing.expectError(error.Overflow, solve(allocator, "Register A: 1\nRegister B: 0\nRegister C: 0\n\nProgram: 0,8\n"));
}
