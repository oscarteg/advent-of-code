const std = @import("std");

const Pos = struct { x: i32, y: i32 };

const State = struct {
    elves: std.AutoHashMap(Pos, void),
    direction_idx: usize,
    arena: std.mem.Allocator,

    fn init(arena: std.mem.Allocator) State {
        return State{
            .elves = std.AutoHashMap(Pos, void).init(arena),
            .direction_idx = 0,
            .arena = arena,
        };
    }

    fn deinit(self: *State) void {
        self.elves.deinit();
    }

    fn hasNeighbors(self: *const State, pos: Pos) bool {
        var dx: i32 = -1;
        while (dx <= 1) : (dx += 1) {
            var dy: i32 = -1;
            while (dy <= 1) : (dy += 1) {
                if (dx == 0 and dy == 0) continue;
                if (self.elves.contains(Pos{ .x = pos.x + dx, .y = pos.y + dy })) {
                    return true;
                }
            }
        }
        return false;
    }

    fn canMoveNorth(self: *const State, pos: Pos) bool {
        return !self.elves.contains(Pos{ .x = pos.x - 1, .y = pos.y - 1 }) and
            !self.elves.contains(Pos{ .x = pos.x, .y = pos.y - 1 }) and
            !self.elves.contains(Pos{ .x = pos.x + 1, .y = pos.y - 1 });
    }

    fn canMoveSouth(self: *const State, pos: Pos) bool {
        return !self.elves.contains(Pos{ .x = pos.x - 1, .y = pos.y + 1 }) and
            !self.elves.contains(Pos{ .x = pos.x, .y = pos.y + 1 }) and
            !self.elves.contains(Pos{ .x = pos.x + 1, .y = pos.y + 1 });
    }

    fn canMoveWest(self: *const State, pos: Pos) bool {
        return !self.elves.contains(Pos{ .x = pos.x - 1, .y = pos.y - 1 }) and
            !self.elves.contains(Pos{ .x = pos.x - 1, .y = pos.y }) and
            !self.elves.contains(Pos{ .x = pos.x - 1, .y = pos.y + 1 });
    }

    fn canMoveEast(self: *const State, pos: Pos) bool {
        return !self.elves.contains(Pos{ .x = pos.x + 1, .y = pos.y - 1 }) and
            !self.elves.contains(Pos{ .x = pos.x + 1, .y = pos.y }) and
            !self.elves.contains(Pos{ .x = pos.x + 1, .y = pos.y + 1 });
    }

    fn proposeMove(self: *const State, pos: Pos) ?Pos {
        if (!self.hasNeighbors(pos)) return null;

        var i: usize = 0;
        while (i < 4) : (i += 1) {
            const dir = (self.direction_idx + i) % 4;
            switch (dir) {
                0 => if (self.canMoveNorth(pos)) return Pos{ .x = pos.x, .y = pos.y - 1 },
                1 => if (self.canMoveSouth(pos)) return Pos{ .x = pos.x, .y = pos.y + 1 },
                2 => if (self.canMoveWest(pos)) return Pos{ .x = pos.x - 1, .y = pos.y },
                3 => if (self.canMoveEast(pos)) return Pos{ .x = pos.x + 1, .y = pos.y },
                else => {},
            }
        }
        return null;
    }

    fn round(self: *State) bool {
        var proposals = std.AutoHashMap(Pos, std.ArrayList(Pos)).init(self.arena);
        defer {
            var it = proposals.valueIterator();
            while (it.next()) |list| {
                list.deinit(self.arena);
            }
            proposals.deinit();
        }

        var elves = self.elves.keyIterator();
        while (elves.next()) |elf| {
            if (self.proposeMove(elf.*)) |new_pos| {
                if (proposals.getPtr(new_pos)) |list| {
                    list.append(self.arena, elf.*) catch {};
                } else {
                    var list: std.ArrayList(Pos) = .empty;
                    list.append(self.arena, elf.*) catch {};
                    proposals.put(new_pos, list) catch {};
                }
            }
        }

        var moved = false;
        var prop_it = proposals.iterator();
        while (prop_it.next()) |entry| {
            if (entry.value_ptr.items.len == 1) {
                const old_pos = entry.value_ptr.items[0];
                const new_pos = entry.key_ptr.*;
                _ = self.elves.remove(old_pos);
                self.elves.put(new_pos, {}) catch {};
                moved = true;
            }
        }

        self.direction_idx = (self.direction_idx + 1) % 4;
        return moved;
    }

    fn countEmptyTiles(self: *const State) i32 {
        var min_x: i32 = std.math.maxInt(i32);
        var max_x: i32 = std.math.minInt(i32);
        var min_y: i32 = std.math.maxInt(i32);
        var max_y: i32 = std.math.minInt(i32);

        var elves = self.elves.keyIterator();
        while (elves.next()) |elf| {
            min_x = @min(min_x, elf.x);
            max_x = @max(max_x, elf.x);
            min_y = @min(min_y, elf.y);
            max_y = @max(max_y, elf.y);
        }

        const width = max_x - min_x + 1;
        const height = max_y - min_y + 1;
        const total = width * height;
        return total - @as(i32, @intCast(self.elves.count()));
    }
};

fn parseInput(arena: std.mem.Allocator, input: []const u8) !State {
    var state = State.init(arena);
    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    var y: i32 = 0;

    while (lines.next()) |line| {
        if (line.len == 0) continue;
        var x: i32 = 0;
        for (line) |ch| {
            if (ch == '#') {
                try state.elves.put(Pos{ .x = x, .y = y }, {});
            } else if (ch != '.') {
                state.deinit();
                return error.InvalidCharacter;
            }
            x += 1;
        }
        y += 1;
    }

    if (state.elves.count() == 0) {
        state.deinit();
        return error.NoElves;
    }
    return state;
}

fn solve(arena: std.mem.Allocator, input: []const u8) ![2]i32 {
    var state = try parseInput(arena, input);
    defer state.deinit();

    var part2: i32 = 0;
    var round_num: i32 = 1;
    while (round_num <= 10) : (round_num += 1) {
        if (!state.round()) {
            part2 = round_num;
            break;
        }
    }

    const part1 = state.countEmptyTiles();

    if (part2 == 0) {
        round_num = 10;
        while (true) {
            round_num += 1;
            if (!state.round()) {
                part2 = round_num;
                break;
            }
        }
    }

    return [2]i32{ part1, part2 };
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
    try stdout.print("Part 1: {d}\nPart 2: {d}\n", .{ answers[0], answers[1] });
    try stdout.flush();
}

test "part1_example" {
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day_23.txt");
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(110, answers[0]);
}

test "part2_example" {
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day_23.txt");
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(20, answers[1]);
}

test "single_elf" {
    const allocator = std.testing.allocator;
    const input = "#";
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(0, answers[0]);
    try std.testing.expectEqual(1, answers[1]);
}

test "malformed_input" {
    const allocator = std.testing.allocator;
    const input = "#\n$\n#";
    const result = solve(allocator, input);
    try std.testing.expectError(error.InvalidCharacter, result);
}
