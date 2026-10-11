const std = @import("std");

const ParseError = error{ UnexpectedEof, InvalidFormat, InvalidNumber, OutOfMemory, InvalidCharacter, Overflow };

const Packet = union(enum) {
    int: i32,
    list: std.ArrayList(Packet),
};

const Parser = struct {
    input: []const u8,
    pos: usize = 0,

    fn parse(self: *Parser, allocator: std.mem.Allocator) ParseError!Packet {
        self.skipWhitespace();
        if (self.pos >= self.input.len) return error.UnexpectedEof;

        if (self.input[self.pos] == '[') {
            return self.parseList(allocator);
        } else {
            return self.parseInt();
        }
    }

    fn parseList(self: *Parser, allocator: std.mem.Allocator) ParseError!Packet {
        self.pos += 1;
        self.skipWhitespace();

        var list: std.ArrayList(Packet) = .empty;

        if (self.pos < self.input.len and self.input[self.pos] == ']') {
            self.pos += 1;
            return Packet{ .list = list };
        }

        while (true) {
            const packet = try self.parse(allocator);
            try list.append(allocator, packet);
            self.skipWhitespace();

            if (self.pos >= self.input.len) return error.UnexpectedEof;

            if (self.input[self.pos] == ',') {
                self.pos += 1;
                self.skipWhitespace();
            } else if (self.input[self.pos] == ']') {
                self.pos += 1;
                return Packet{ .list = list };
            } else {
                return error.InvalidFormat;
            }
        }
    }

    fn parseInt(self: *Parser) ParseError!Packet {
        const start = self.pos;
        if (self.pos < self.input.len and self.input[self.pos] == '-') {
            self.pos += 1;
        }
        while (self.pos < self.input.len and std.ascii.isDigit(self.input[self.pos])) {
            self.pos += 1;
        }
        if (start == self.pos) return error.InvalidNumber;

        const numStr = self.input[start..self.pos];
        const num = try std.fmt.parseInt(i32, numStr, 10);
        return Packet{ .int = num };
    }

    fn skipWhitespace(self: *Parser) void {
        while (self.pos < self.input.len and std.ascii.isWhitespace(self.input[self.pos])) {
            self.pos += 1;
        }
    }
};

fn compare(a: Packet, b: Packet) std.math.Order {
    switch (a) {
        .int => |av| {
            switch (b) {
                .int => |bv| return std.math.order(av, bv),
                .list => |blist| {
                    var alist: std.ArrayList(Packet) = .empty;
                    defer alist.deinit(std.heap.page_allocator);
                    alist.append(std.heap.page_allocator, Packet{ .int = av }) catch unreachable;
                    return compareSlices(alist.items, blist.items);
                },
            }
        },
        .list => |alist| {
            switch (b) {
                .int => |bv| {
                    var blist: std.ArrayList(Packet) = .empty;
                    defer blist.deinit(std.heap.page_allocator);
                    blist.append(std.heap.page_allocator, Packet{ .int = bv }) catch unreachable;
                    return compareSlices(alist.items, blist.items);
                },
                .list => |blist| return compareSlices(alist.items, blist.items),
            }
        },
    }
}

fn compareSlices(a: []const Packet, b: []const Packet) std.math.Order {
    const minLen = @min(a.len, b.len);
    for (0..minLen) |i| {
        const ord = compare(a[i], b[i]);
        if (ord != .eq) return ord;
    }
    return std.math.order(a.len, b.len);
}

fn packetOrder(ctx: void, a: Packet, b: Packet) bool {
    _ = ctx;
    return compare(a, b) == .lt;
}

fn solvePart1(allocator: std.mem.Allocator, input: []const u8) !usize {
    var sum: usize = 0;
    var pairIdx: usize = 1;

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        if (std.mem.trim(u8, line, " \t").len == 0) continue;

        var parser1: Parser = .{ .input = line };
        const packet1 = try parser1.parse(allocator);

        if (lines.next()) |line2| {
            var parser2: Parser = .{ .input = line2 };
            const packet2 = try parser2.parse(allocator);

            if (compare(packet1, packet2) == .lt) {
                sum += pairIdx;
            }
            pairIdx += 1;
        }
    }

    return sum;
}

fn solvePart2(allocator: std.mem.Allocator, input: []const u8) !usize {
    var packets: std.ArrayList(Packet) = .empty;

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t");
        if (trimmed.len == 0) continue;

        var parser: Parser = .{ .input = trimmed };
        const packet = try parser.parse(allocator);
        try packets.append(allocator, packet);
    }

    var divider1: std.ArrayList(Packet) = .empty;
    try divider1.append(allocator, Packet{ .int = 2 });
    var div1: std.ArrayList(Packet) = .empty;
    try div1.append(allocator, Packet{ .list = divider1 });
    try packets.append(allocator, Packet{ .list = div1 });

    var divider2: std.ArrayList(Packet) = .empty;
    try divider2.append(allocator, Packet{ .int = 6 });
    var div2: std.ArrayList(Packet) = .empty;
    try div2.append(allocator, Packet{ .list = divider2 });
    try packets.append(allocator, Packet{ .list = div2 });

    std.mem.sortUnstable(Packet, packets.items, {}, packetOrder);

    var pos1: usize = 0;
    var pos2: usize = 0;
    for (packets.items, 1..) |p, idx| {
        if (isDiv1(p)) pos1 = idx;
        if (isDiv2(p)) pos2 = idx;
    }

    return pos1 * pos2;
}

fn isDiv1(p: Packet) bool {
    if (p != .list) return false;
    if (p.list.items.len != 1) return false;
    if (p.list.items[0] != .list) return false;
    if (p.list.items[0].list.items.len != 1) return false;
    if (p.list.items[0].list.items[0] != .int) return false;
    return p.list.items[0].list.items[0].int == 2;
}

fn isDiv2(p: Packet) bool {
    if (p != .list) return false;
    if (p.list.items.len != 1) return false;
    if (p.list.items[0] != .list) return false;
    if (p.list.items[0].list.items.len != 1) return false;
    if (p.list.items[0].list.items[0] != .int) return false;
    return p.list.items[0].list.items[0].int == 6;
}

pub fn main(init: std.process.Init) !void {
    const arena = init.arena.allocator();
    var stdin_buffer: [4096]u8 = undefined;
    var stdin_reader = std.Io.File.stdin().reader(init.io, &stdin_buffer);
    const input = try stdin_reader.interface.allocRemaining(arena, .unlimited);

    const part1 = try solvePart1(arena, input);
    const part2 = try solvePart2(arena, input);

    var stdout_buffer: [256]u8 = undefined;
    var stdout_writer = std.Io.File.stdout().writer(init.io, &stdout_buffer);
    const stdout = &stdout_writer.interface;
    try stdout.print("Part 1: {d}\nPart 2: {d}\n", .{ part1, part2 });
    try stdout.flush();
}

test "parse int" {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();
    var parser: Parser = .{ .input = "42" };
    const packet = try parser.parse(allocator);
    try std.testing.expectEqual(packet.int, 42);
}

test "parse list" {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();
    var parser: Parser = .{ .input = "[]" };
    const packet = try parser.parse(allocator);
    try std.testing.expectEqual(packet.list.items.len, 0);
}

test "parse complex" {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();
    var parser: Parser = .{ .input = "[1,[2,3],4]" };
    const packet = try parser.parse(allocator);
    try std.testing.expectEqual(packet.list.items.len, 3);
}

test "part1 example" {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();
    const input = @embedFile("examples/day_13.txt");
    const result = try solvePart1(allocator, input);
    try std.testing.expectEqual(result, 13);
}

test "part2 example" {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();
    const input = @embedFile("examples/day_13.txt");
    const result = try solvePart2(allocator, input);
    try std.testing.expectEqual(result, 140);
}

test "rejects an unterminated list" {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    try std.testing.expectError(error.UnexpectedEof, solvePart1(arena.allocator(), "[1,2\n[3]\n"));
}

test "rejects a non-numeric value" {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    try std.testing.expectError(error.InvalidNumber, solvePart2(arena.allocator(), "[a]\n[1]\n"));
}
