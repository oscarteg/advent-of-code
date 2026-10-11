const std = @import("std");

const Answers = struct {
    part1: i64,
    part2: []const u8,
};

const Op = enum(u2) { and_gate, or_gate, xor_gate };

const WireId = u32;
const no_wire: WireId = std.math.maxInt(WireId);
const no_gate: usize = std.math.maxInt(usize);
const max_swaps = 4;
const max_bits: u32 = 63;

const Gate = struct {
    left: WireId,
    right: WireId,
    op: Op,
};

const Circuit = struct {
    allocator: std.mem.Allocator,
    names: std.ArrayList([]const u8),
    ids: std.StringHashMap(WireId),
    initial: std.ArrayList(i8), // per wire: -1 unknown, else 0 or 1
    driver: std.ArrayList(usize), // per wire: gate index, or no_gate
    gates: std.ArrayList(Gate),
    outputs: std.ArrayList(WireId), // outputs[i] is the wire gate i drives
    bits: u32 = 0, // number of x (and y) input wires

    fn init(allocator: std.mem.Allocator) Circuit {
        return .{
            .allocator = allocator,
            .names = .empty,
            .ids = std.StringHashMap(WireId).init(allocator),
            .initial = .empty,
            .driver = .empty,
            .gates = .empty,
            .outputs = .empty,
        };
    }

    fn deinit(self: *Circuit) void {
        self.names.deinit(self.allocator);
        self.ids.deinit();
        self.initial.deinit(self.allocator);
        self.driver.deinit(self.allocator);
        self.gates.deinit(self.allocator);
        self.outputs.deinit(self.allocator);
    }

    fn intern(self: *Circuit, name: []const u8) !WireId {
        const entry = try self.ids.getOrPut(name);
        if (!entry.found_existing) {
            entry.value_ptr.* = @intCast(self.names.items.len);
            try self.names.append(self.allocator, name);
            try self.initial.append(self.allocator, -1);
            try self.driver.append(self.allocator, no_gate);
        }
        return entry.value_ptr.*;
    }

    fn wireId(self: *const Circuit, prefix: u8, bit: u32) WireId {
        const index: u8 = @intCast(bit);
        const name = [3]u8{ prefix, '0' + index / 10, '0' + index % 10 };
        return self.ids.get(&name) orelse no_wire;
    }
};

fn parseOp(text: []const u8) !Op {
    if (std.mem.eql(u8, text, "AND")) return .and_gate;
    if (std.mem.eql(u8, text, "OR")) return .or_gate;
    if (std.mem.eql(u8, text, "XOR")) return .xor_gate;
    return error.UnknownOperator;
}

fn parseGate(circuit: *Circuit, line: []const u8) !void {
    var words: [5][]const u8 = undefined;
    var count: usize = 0;
    var parts = std.mem.tokenizeScalar(u8, line, ' ');
    while (parts.next()) |word| {
        if (count == words.len) return error.MalformedGate;
        words[count] = word;
        count += 1;
    }
    if (count != words.len or !std.mem.eql(u8, words[3], "->")) return error.MalformedGate;

    const op = try parseOp(words[1]);
    const left = try circuit.intern(words[0]);
    const right = try circuit.intern(words[2]);
    const output = try circuit.intern(words[4]);
    if (circuit.driver.items[output] != no_gate) return error.WireDrivenTwice;
    circuit.driver.items[output] = circuit.gates.items.len;
    try circuit.gates.append(circuit.allocator, .{ .left = left, .right = right, .op = op });
    try circuit.outputs.append(circuit.allocator, output);
}

fn parse(allocator: std.mem.Allocator, input: []const u8) !Circuit {
    var circuit = Circuit.init(allocator);
    errdefer circuit.deinit();

    var lines = std.mem.splitScalar(u8, input, '\n');
    while (lines.next()) |raw| {
        const line = std.mem.trim(u8, raw, " \t\r");
        if (line.len == 0) continue;

        if (std.mem.indexOf(u8, line, "->") != null) {
            try parseGate(&circuit, line);
        } else if (std.mem.indexOfScalar(u8, line, ':')) |colon| {
            const name = std.mem.trim(u8, line[0..colon], " \t");
            const value = std.mem.trim(u8, line[colon + 1 ..], " \t");
            if (name.len == 0 or value.len != 1 or (value[0] != '0' and value[0] != '1')) {
                return error.MalformedInitialValue;
            }
            const wire = try circuit.intern(name);
            circuit.initial.items[wire] = if (value[0] == '1') 1 else 0;
        } else {
            return error.MalformedLine;
        }
    }

    var has_output_wire = false;
    for (circuit.outputs.items) |wire| {
        if (circuit.names.items[wire][0] == 'z') has_output_wire = true;
    }
    if (!has_output_wire) return error.NoOutputWires;

    while (circuit.bits < max_bits and circuit.wireId('x', circuit.bits) != no_wire) {
        circuit.bits += 1;
    }
    return circuit;
}

const Evaluator = struct {
    circuit: *const Circuit,
    values: []i8,
    visiting: []bool,

    fn valueOf(self: *Evaluator, wire: WireId) !bool {
        if (self.values[wire] >= 0) return self.values[wire] == 1;
        const driver = self.circuit.driver.items[wire];
        if (driver == no_gate) return error.MissingWireValue;
        if (self.visiting[wire]) return error.CircuitHasCycle;

        self.visiting[wire] = true;
        const gate = self.circuit.gates.items[driver];
        const left = try self.valueOf(gate.left);
        const right = try self.valueOf(gate.right);
        self.visiting[wire] = false;

        const value = switch (gate.op) {
            .and_gate => left and right,
            .or_gate => left or right,
            .xor_gate => left != right,
        };
        self.values[wire] = if (value) 1 else 0;
        return value;
    }
};

fn simulate(allocator: std.mem.Allocator, circuit: *const Circuit) !i64 {
    const values = try allocator.dupe(i8, circuit.initial.items);
    defer allocator.free(values);
    const visiting = try allocator.alloc(bool, circuit.names.items.len);
    defer allocator.free(visiting);
    @memset(visiting, false);

    var evaluator = Evaluator{ .circuit = circuit, .values = values, .visiting = visiting };
    var result: i64 = 0;
    for (circuit.outputs.items) |wire| {
        const name = circuit.names.items[wire];
        if (name[0] != 'z') continue;
        const bit = try std.fmt.parseInt(u6, name[1..], 10);
        if (try evaluator.valueOf(wire)) result |= @as(i64, 1) << bit;
    }
    return result;
}

fn gateKey(left: WireId, right: WireId, op: Op) u64 {
    const low: u64 = @min(left, right);
    const high: u64 = @max(left, right);
    return (low << 32) | (high << 2) | @intFromEnum(op);
}

const GateLookup = std.AutoHashMap(u64, usize);

fn buildLookup(allocator: std.mem.Allocator, circuit: *const Circuit) !GateLookup {
    var lookup = GateLookup.init(allocator);
    errdefer lookup.deinit();
    for (circuit.gates.items, 0..) |gate, index| {
        const entry = try lookup.getOrPut(gateKey(gate.left, gate.right, gate.op));
        if (!entry.found_existing) entry.value_ptr.* = index;
    }
    return lookup;
}

const Adder = struct {
    circuit: *const Circuit,
    lookup: *const GateLookup,
    outputs: []WireId,

    fn driven(self: Adder, left: WireId, right: WireId, op: Op) WireId {
        if (left == no_wire or right == no_wire) return no_wire;
        const index = self.lookup.get(gateKey(left, right, op)) orelse return no_wire;
        return self.outputs[index];
    }

    // A correct ripple-carry adder computes z0 = x0 XOR y0 and carry0 = x0 AND y0, then for
    // every higher bit z_i = (x_i XOR y_i) XOR carry_(i-1) and
    // carry_i = (x_i AND y_i) OR ((x_i XOR y_i) AND carry_(i-1)), with the last carry on the
    // top z wire. Returns the lowest bit that does not match, or bits + 1 when all of them do.
    fn firstBrokenBit(self: Adder) u32 {
        var carry: WireId = no_wire;
        var bit: u32 = 0;
        while (bit < self.circuit.bits) : (bit += 1) {
            const output = self.circuit.wireId('z', bit);
            const x = self.circuit.wireId('x', bit);
            const y = self.circuit.wireId('y', bit);
            const sum = self.driven(x, y, .xor_gate);
            const generated = self.driven(x, y, .and_gate);
            if (output == no_wire or sum == no_wire or generated == no_wire) return bit;

            if (bit == 0) {
                if (sum != output) return 0;
                carry = generated;
                continue;
            }
            if (self.driven(sum, carry, .xor_gate) != output) return bit;
            const propagated = self.driven(sum, carry, .and_gate);
            if (propagated == no_wire) return bit;
            const next_carry = self.driven(propagated, generated, .or_gate);
            if (next_carry == no_wire) return bit;
            carry = next_carry;
        }

        const top = self.circuit.wireId('z', self.circuit.bits);
        if (top == no_wire or carry != top) return self.circuit.bits;
        return self.circuit.bits + 1;
    }

    // Depth-first search over output swaps, keeping only swaps that push the first broken
    // bit higher; the puzzle input needs exactly four of them.
    fn repair(self: Adder, allocator: std.mem.Allocator, remaining: u32, swapped: *std.ArrayList(WireId)) !bool {
        const broken = self.firstBrokenBit();
        if (broken > self.circuit.bits) return true;
        if (remaining == 0) return false;

        for (0..self.outputs.len) |a| {
            for (a + 1..self.outputs.len) |b| {
                std.mem.swap(WireId, &self.outputs[a], &self.outputs[b]);
                if (self.firstBrokenBit() > broken) {
                    try swapped.append(allocator, self.outputs[a]);
                    try swapped.append(allocator, self.outputs[b]);
                    if (try self.repair(allocator, remaining - 1, swapped)) return true;
                    swapped.items.len -= 2;
                }
                std.mem.swap(WireId, &self.outputs[a], &self.outputs[b]);
            }
        }
        return false;
    }
};

// Empty when the circuit is not an adder that four swaps can repair, which includes both
// part 1 examples.
fn findSwappedWires(allocator: std.mem.Allocator, circuit: *const Circuit) ![]const u8 {
    var lookup = try buildLookup(allocator, circuit);
    defer lookup.deinit();
    const outputs = try allocator.dupe(WireId, circuit.outputs.items);
    defer allocator.free(outputs);

    var swapped: std.ArrayList(WireId) = .empty;
    defer swapped.deinit(allocator);

    const adder = Adder{ .circuit = circuit, .lookup = &lookup, .outputs = outputs };
    if (!try adder.repair(allocator, max_swaps, &swapped)) return allocator.alloc(u8, 0);

    var names: std.ArrayList([]const u8) = .empty;
    defer names.deinit(allocator);
    for (swapped.items) |wire| try names.append(allocator, circuit.names.items[wire]);
    std.mem.sort([]const u8, names.items, {}, struct {
        fn lessThan(_: void, a: []const u8, b: []const u8) bool {
            return std.mem.lessThan(u8, a, b);
        }
    }.lessThan);

    var joined: std.ArrayList(u8) = .empty;
    errdefer joined.deinit(allocator);
    for (names.items, 0..) |name, index| {
        if (index > 0) try joined.append(allocator, ',');
        try joined.appendSlice(allocator, name);
    }
    return joined.toOwnedSlice(allocator);
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var circuit = try parse(allocator, input);
    defer circuit.deinit();
    return .{
        .part1 = try simulate(allocator, &circuit),
        .part2 = try findSwappedWires(allocator, &circuit),
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
    try stdout.print("Part 1: {d}\nPart 2: {s}\n", .{ answers.part1, answers.part2 });
    try stdout.flush();
}

fn adderWireName(allocator: std.mem.Allocator, prefix: u8, bit: u32) ![]const u8 {
    return std.fmt.allocPrint(allocator, "{c}{d:0>2}", .{ prefix, bit });
}

fn swappedName(name: []const u8, swaps: []const [2][]const u8) []const u8 {
    for (swaps) |pair| {
        if (std.mem.eql(u8, name, pair[0])) return pair[1];
        if (std.mem.eql(u8, name, pair[1])) return pair[0];
    }
    return name;
}

// Builds a correct ripple-carry adder for x + y, then renames the gate outputs listed in
// `swaps` so the circuit is wrong in exactly the way part 2 describes.
fn adderInput(allocator: std.mem.Allocator, bits: u32, x: u64, y: u64, swaps: []const [2][]const u8) ![]u8 {
    var arena = std.heap.ArenaAllocator.init(allocator);
    defer arena.deinit();
    const scratch = arena.allocator();

    var text: std.ArrayList(u8) = .empty;
    errdefer text.deinit(allocator);

    for ([2]u8{ 'x', 'y' }, [2]u64{ x, y }) |prefix, value| {
        var bit: u32 = 0;
        while (bit < bits) : (bit += 1) {
            const shift: u6 = @intCast(bit);
            const line = try std.fmt.allocPrint(scratch, "{c}{d:0>2}: {d}\n", .{ prefix, bit, (value >> shift) & 1 });
            try text.appendSlice(allocator, line);
        }
    }
    try text.append(allocator, '\n');

    var gates: std.ArrayList([4][]const u8) = .empty;
    try gates.append(scratch, .{
        try adderWireName(scratch, 'x', 0),
        "XOR",
        try adderWireName(scratch, 'y', 0),
        try adderWireName(scratch, 'z', 0),
    });
    const first_carry = if (bits == 1)
        try adderWireName(scratch, 'z', 1)
    else
        try adderWireName(scratch, 'c', 0);
    try gates.append(scratch, .{
        try adderWireName(scratch, 'x', 0),
        "AND",
        try adderWireName(scratch, 'y', 0),
        first_carry,
    });

    var bit: u32 = 1;
    while (bit < bits) : (bit += 1) {
        const carry_in = try adderWireName(scratch, 'c', bit - 1);
        const sum = try adderWireName(scratch, 's', bit);
        const generated = try adderWireName(scratch, 'a', bit);
        const last = bit + 1 == bits;
        try gates.append(scratch, .{ try adderWireName(scratch, 'x', bit), "XOR", try adderWireName(scratch, 'y', bit), sum });
        try gates.append(scratch, .{ try adderWireName(scratch, 'x', bit), "AND", try adderWireName(scratch, 'y', bit), generated });
        try gates.append(scratch, .{ sum, "XOR", carry_in, try adderWireName(scratch, 'z', bit) });
        try gates.append(scratch, .{ sum, "AND", carry_in, try adderWireName(scratch, 'p', bit) });
        const carry_out = if (last)
            try adderWireName(scratch, 'z', bits)
        else
            try adderWireName(scratch, 'c', bit);
        try gates.append(scratch, .{ try adderWireName(scratch, 'p', bit), "OR", generated, carry_out });
    }

    for (gates.items) |gate| {
        const line = try std.fmt.allocPrint(scratch, "{s} {s} {s} -> {s}\n", .{ gate[0], gate[1], gate[2], swappedName(gate[3], swaps) });
        try text.appendSlice(allocator, line);
    }
    return text.toOwnedSlice(allocator);
}

test "official part one example" {
    const answers = try solve(std.testing.allocator, @embedFile("examples/day24.txt"));
    defer std.testing.allocator.free(answers.part2);
    try std.testing.expectEqual(@as(i64, 2024), answers.part1);
}

test "small part one example" {
    const answers = try solve(std.testing.allocator,
        \\x00: 1
        \\x01: 1
        \\x02: 1
        \\y00: 0
        \\y01: 1
        \\y02: 0
        \\
        \\x00 AND y00 -> z00
        \\x01 XOR y01 -> z01
        \\x02 OR y02 -> z02
        \\
    );
    defer std.testing.allocator.free(answers.part2);
    try std.testing.expectEqual(@as(i64, 4), answers.part1);
    try std.testing.expectEqualStrings("", answers.part2);
}

test "a correct adder adds its inputs and needs no swaps" {
    const input = try adderInput(std.testing.allocator, 8, 0b1011_0101, 0b0110_1110, &.{});
    defer std.testing.allocator.free(input);

    const answers = try solve(std.testing.allocator, input);
    defer std.testing.allocator.free(answers.part2);
    try std.testing.expectEqual(@as(i64, 0b1011_0101 + 0b0110_1110), answers.part1);
    try std.testing.expectEqualStrings("", answers.part2);
}

test "four crossed wires in an adder are reported sorted" {
    const swaps = [_][2][]const u8{
        .{ "s01", "a01" },
        .{ "z03", "a03" },
        .{ "z05", "c05" },
        .{ "s06", "a06" },
    };
    const input = try adderInput(std.testing.allocator, 8, 0b1011_0101, 0b0110_1110, &swaps);
    defer std.testing.allocator.free(input);

    const answers = try solve(std.testing.allocator, input);
    defer std.testing.allocator.free(answers.part2);
    try std.testing.expectEqualStrings("a01,a03,a06,c05,s01,s06,z03,z05", answers.part2);
}

test "edge case: a half adder is the smallest correct adder" {
    const input = try adderInput(std.testing.allocator, 1, 1, 1, &.{});
    defer std.testing.allocator.free(input);

    const answers = try solve(std.testing.allocator, input);
    defer std.testing.allocator.free(answers.part2);
    try std.testing.expectEqual(@as(i64, 2), answers.part1);
    try std.testing.expectEqualStrings("", answers.part2);
}

test "malformed input returns an error" {
    const allocator = std.testing.allocator;
    try std.testing.expectError(error.NoOutputWires, solve(allocator, ""));
    try std.testing.expectError(error.NoOutputWires, solve(allocator, "x00: 1\ny00: 1\n"));
    try std.testing.expectError(error.MalformedInitialValue, solve(allocator, "x00: 2\ny00: 1\n\nx00 AND y00 -> z00\n"));
    try std.testing.expectError(error.UnknownOperator, solve(allocator, "x00: 1\ny00: 1\n\nx00 NAND y00 -> z00\n"));
    try std.testing.expectError(error.MalformedLine, solve(allocator, "x00 1\n\nx00 AND y00 -> z00\n"));
    try std.testing.expectError(error.MalformedGate, solve(allocator, "x00: 1\ny00: 1\n\nx00 XOR -> z00\n"));
    try std.testing.expectError(error.MissingWireValue, solve(allocator, "x00: 1\n\nx00 AND y00 -> z00\n"));
    try std.testing.expectError(error.WireDrivenTwice, solve(allocator, "x00: 1\ny00: 1\n\nx00 AND y00 -> z00\nx00 OR y00 -> z00\n"));
    try std.testing.expectError(error.CircuitHasCycle, solve(allocator, "y00: 1\n\nz01 AND y00 -> z00\nz00 OR y00 -> z01\n"));
}
