const file = @embedFile("./day6-input.txt");
const std = @import("std");
const Allocator = std.heap.GeneralPurposeAllocator;
const List = std.ArrayList;
const Map = std.AutoHashMap;
const StrMap = std.StringHashMap;
const BitSet = std.DynamicBitSet;

const data = @embedFile("data/day06.txt");

// Useful stdlib functions
const tokenizeAny = std.mem.tokenizeAny;
const tokenizeSeq = std.mem.tokenizeSequence;
const tokenizeSca = std.mem.tokenizeScalar;
const splitAny = std.mem.splitAny;
const splitSeq = std.mem.splitSequence;
const splitSca = std.mem.splitScalar;
const indexOf = std.mem.indexOfScalar;
const indexOfAny = std.mem.indexOfAny;
const indexOfStr = std.mem.indexOfPosLinear;
const lastIndexOf = std.mem.lastIndexOfScalar;
const lastIndexOfAny = std.mem.lastIndexOfAny;
const lastIndexOfStr = std.mem.lastIndexOfLinear;
const trim = std.mem.trim;
const sliceMin = std.mem.min;
const sliceMax = std.mem.max;

const parseInt = std.fmt.parseInt;
const parseFloat = std.fmt.parseFloat;

const print = std.debug.print;
const assert = std.debug.assert;

const sort = std.sort.block;
const asc = std.sort.asc;
const desc = std.sort.desc;

const Direction = enum {
    Up,
    Down,
    Left,
    Right,

    pub fn turnRight(self: Direction) Direction {
        return switch (self) {
            .Up => .Right,
            .Right => .Down,
            .Down => .Left,
            .Left => .Up,
        };
    }

    pub fn moveDelta(self: Direction, delta: u32) struct { dcol: u32, drow: u32 } {
        return switch (self) {
            .Up => .{ .x = -1, .y = 0 },
            .Right => .{ .x = 1, .y = 0 },
            .Down => .{ .x = 0, .y = -delta },
            .Left => .{ .x = -delta, .y = 0 },
        };
    }
};

const Position = struct {
    col: u32,
    row: u32,

    fn toKey(self: Position) u64 {
        return @intCast(self.row) << 32 | @intCast(u64, @as(u32, @intCast(u32, self.col)));
    }
};

// Guard rules
// If there's an obstacle ahead, turn 90° right
// Otherwise, move forward one step
// The guard starts at a position marked with a direction (^, >, v, <)

pub fn main() !void {
    print("test\n", .{});
}
