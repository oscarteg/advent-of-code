const std = @import("std");

/// Presents are given as 3x3 patterns, so every orientation fits a 3x3 box.
const present_size = 3;

const Pattern = [present_size][present_size]bool;

const Answers = struct { part1: usize };

/// One cell of an oriented present, relative to the cell the present is
/// anchored at: its own first cell in row-major order, so rows never go up.
const Offset = struct { row: usize, column: isize };

/// A present in one of its distinct orientations.
const Variant = struct {
    cells: []const Offset,
    height: usize,
};

const Present = struct {
    area: usize,
    variants: []const Variant,
};

/// A region under one tree: its size plus how many of each present must fit.
const Region = struct {
    width: usize,
    height: usize,
    counts: []const usize,
};

const Farm = struct {
    presents: []const Present,
    regions: []const Region,
};

pub fn solve(arena: std.mem.Allocator, input: []const u8) !Answers {
    const farm = try parse(arena, input);
    var fitting: usize = 0;
    for (farm.regions) |region| {
        if (try regionFits(arena, farm.presents, region)) fitting += 1;
    }
    return .{ .part1 = fitting };
}

fn parse(arena: std.mem.Allocator, input: []const u8) !Farm {
    var presents: std.ArrayList(Present) = .empty;
    var regions: std.ArrayList(Region) = .empty;

    var lines = std.mem.splitScalar(u8, input, '\n');
    while (lines.next()) |raw_line| {
        const line = std.mem.trim(u8, raw_line, &std.ascii.whitespace);
        if (line.len == 0) continue;
        if (line[line.len - 1] == ':') {
            if (regions.items.len != 0) return error.PresentAfterRegion;
            const index = std.fmt.parseInt(usize, line[0 .. line.len - 1], 10) catch return error.UnexpectedLine;
            if (index != presents.items.len) return error.PresentOutOfOrder;
            try presents.append(arena, try parsePresent(arena, &lines));
            continue;
        }
        if (presents.items.len == 0) return error.MissingPresents;
        try regions.append(arena, try parseRegion(arena, line, presents.items.len));
    }
    if (presents.items.len == 0) return error.MissingPresents;
    return .{ .presents = presents.items, .regions = regions.items };
}

fn parsePresent(arena: std.mem.Allocator, lines: *std.mem.SplitIterator(u8, .scalar)) !Present {
    var pattern = std.mem.zeroes(Pattern);
    var area: usize = 0;
    for (&pattern) |*pattern_row| {
        const row = std.mem.trim(u8, lines.next() orelse return error.TruncatedPresent, &std.ascii.whitespace);
        if (row.len != present_size) return error.InvalidPresentRow;
        for (row, 0..) |cell, column| {
            switch (cell) {
                '#' => {
                    pattern_row[column] = true;
                    area += 1;
                },
                '.' => {},
                else => return error.InvalidPresentCell,
            }
        }
    }
    if (area == 0) return error.EmptyPresent;
    return .{ .area = area, .variants = try variantsOf(arena, pattern) };
}

fn parseRegion(arena: std.mem.Allocator, line: []const u8, present_count: usize) !Region {
    const colon = std.mem.indexOfScalar(u8, line, ':') orelse return error.InvalidRegion;
    const size = std.mem.trim(u8, line[0..colon], &std.ascii.whitespace);
    const separator = std.mem.indexOfScalar(u8, size, 'x') orelse return error.InvalidRegion;
    const width = std.fmt.parseInt(usize, size[0..separator], 10) catch return error.InvalidRegion;
    const height = std.fmt.parseInt(usize, size[separator + 1 ..], 10) catch return error.InvalidRegion;
    if (width == 0 or height == 0) return error.InvalidRegion;

    var counts: std.ArrayList(usize) = .empty;
    var fields = std.mem.tokenizeAny(u8, line[colon + 1 ..], " \t");
    while (fields.next()) |field| {
        const count = std.fmt.parseInt(usize, field, 10) catch return error.InvalidRegion;
        try counts.append(arena, count);
    }
    if (counts.items.len != present_count) return error.CountMismatch;
    return .{ .width = width, .height = height, .counts = counts.items };
}

/// Every rotation and flip of a pattern, without duplicates.
fn variantsOf(arena: std.mem.Allocator, pattern: Pattern) ![]const Variant {
    var variants: std.ArrayList(Variant) = .empty;
    var oriented = pattern;
    for (0..2) |flip| {
        if (flip == 1) oriented = mirrored(oriented);
        for (0..4) |_| {
            oriented = rotated(oriented);
            const candidate = try variantFrom(arena, oriented);
            for (variants.items) |existing| {
                if (sameShape(existing.cells, candidate.cells)) break;
            } else {
                try variants.append(arena, candidate);
            }
        }
    }
    return variants.items;
}

fn rotated(pattern: Pattern) Pattern {
    var result = std.mem.zeroes(Pattern);
    for (0..present_size) |row| {
        for (0..present_size) |column| {
            result[row][column] = pattern[present_size - 1 - column][row];
        }
    }
    return result;
}

fn mirrored(pattern: Pattern) Pattern {
    var result = std.mem.zeroes(Pattern);
    for (0..present_size) |row| {
        for (0..present_size) |column| {
            result[row][column] = pattern[row][present_size - 1 - column];
        }
    }
    return result;
}

fn variantFrom(arena: std.mem.Allocator, pattern: Pattern) !Variant {
    var cells: std.ArrayList(Offset) = .empty;
    var anchor: ?[2]usize = null;
    var height: usize = 0;
    for (pattern, 0..) |pattern_row, row| {
        for (pattern_row, 0..) |filled, column| {
            if (!filled) continue;
            const origin = anchor orelse [2]usize{ row, column };
            anchor = origin;
            try cells.append(arena, .{
                .row = row - origin[0],
                .column = @as(isize, @intCast(column)) - @as(isize, @intCast(origin[1])),
            });
            height = row - origin[0] + 1;
        }
    }
    if (anchor == null) return error.EmptyPresent;
    return .{ .cells = cells.items, .height = height };
}

fn sameShape(left: []const Offset, right: []const Offset) bool {
    if (left.len != right.len) return false;
    for (left, right) |a, b| {
        if (a.row != b.row or a.column != b.column) return false;
    }
    return true;
}

/// Decides whether one region can hold every present listed for it.
fn regionFits(arena: std.mem.Allocator, presents: []const Present, region: Region) !bool {
    var wanted: usize = 0;
    var area: usize = 0;
    var min_height: usize = present_size;
    for (region.counts, presents) |count, present| {
        if (count == 0) continue;
        wanted += count;
        area += count * present.area;
        for (present.variants) |variant| min_height = @min(min_height, variant.height);
    }
    if (wanted == 0) return true;

    const cells = region.width * region.height;
    if (area > cells) return false;
    // Every present fits a 3x3 box, so this many need no puzzling at all.
    if (wanted <= (region.width / present_size) * (region.height / present_size)) return true;

    const occupied = try arena.alloc(bool, cells);
    defer arena.free(occupied);
    @memset(occupied, false);
    const counts = try arena.dupe(usize, region.counts);
    defer arena.free(counts);

    var packer = Packer{
        .width = region.width,
        .height = region.height,
        .min_height = min_height,
        .presents = presents,
        .counts = counts,
        .remaining = wanted,
        .occupied = occupied,
    };
    return packer.fits(0, cells - area);
}

const Packer = struct {
    width: usize,
    height: usize,
    min_height: usize,
    presents: []const Present,
    counts: []usize,
    remaining: usize,
    occupied: []bool,

    /// Walks the grid in row-major order: the first undecided cell is either
    /// left empty, spending one cell of the waste budget, or anchors a present.
    /// A present covering that cell must start there, so each packing is tried
    /// exactly once and the budget keeps the search tight.
    fn fits(self: *Packer, start: usize, waste: usize) bool {
        if (self.remaining == 0) return true;
        var budget = waste;
        var index = start;
        while (index < self.occupied.len) : (index += 1) {
            if (self.occupied[index]) continue;
            const row = index / self.width;
            if (row + self.min_height > self.height) return false;
            const column = index % self.width;
            for (self.counts, self.presents) |*count, present| {
                if (count.* == 0) continue;
                for (present.variants) |variant| {
                    if (!self.canPlace(variant, row, column)) continue;
                    self.place(variant, row, column, true);
                    count.* -= 1;
                    self.remaining -= 1;
                    if (self.fits(index + 1, budget)) return true;
                    count.* += 1;
                    self.remaining += 1;
                    self.place(variant, row, column, false);
                }
            }
            if (budget == 0) return false;
            budget -= 1;
        }
        return false;
    }

    fn canPlace(self: *const Packer, variant: Variant, row: usize, column: usize) bool {
        for (variant.cells) |cell| {
            const index = self.cellIndex(cell, row, column) orelse return false;
            if (self.occupied[index]) return false;
        }
        return true;
    }

    fn place(self: *Packer, variant: Variant, row: usize, column: usize, occupied: bool) void {
        for (variant.cells) |cell| {
            self.occupied[self.cellIndex(cell, row, column).?] = occupied;
        }
    }

    fn cellIndex(self: *const Packer, cell: Offset, row: usize, column: usize) ?usize {
        const target_row = row + cell.row;
        const target_column = @as(isize, @intCast(column)) + cell.column;
        if (target_row >= self.height) return null;
        if (target_column < 0 or target_column >= @as(isize, @intCast(self.width))) return null;
        return target_row * self.width + @as(usize, @intCast(target_column));
    }
};

pub fn main(init: std.process.Init) !void {
    const arena = init.arena.allocator();
    var stdin_buffer: [4096]u8 = undefined;
    var stdin_reader = std.Io.File.stdin().reader(init.io, &stdin_buffer);
    const input = try stdin_reader.interface.allocRemaining(arena, .unlimited);
    const answers = try solve(arena, input);
    var stdout_buffer: [256]u8 = undefined;
    var stdout_writer = std.Io.File.stdout().writer(init.io, &stdout_buffer);
    const stdout = &stdout_writer.interface;
    try stdout.print("Part 1: {d}\n", .{answers.part1});
    try stdout.flush();
}

const c_present = "0:\n###\n#..\n###\n\n";
const bar_present = "0:\n###\n...\n...\n\n";

test "counts the regions of the official example that fit" {
    var arena_state = std.heap.ArenaAllocator.init(std.testing.allocator);
    defer arena_state.deinit();

    const answers = try solve(arena_state.allocator(), @embedFile("examples/day12.txt"));
    try std.testing.expectEqual(@as(usize, 2), answers.part1);
}

test "interlocks presents that no 3x3 block would hold" {
    var arena_state = std.heap.ArenaAllocator.init(std.testing.allocator);
    defer arena_state.deinit();

    // The puzzle's own 4x4 example packs two of these; 3x3 lacks the area.
    const answers = try solve(arena_state.allocator(), c_present ++ "4x4: 2\n3x3: 2\n3x3: 1\n");
    try std.testing.expectEqual(@as(usize, 2), answers.part1);
}

test "rotates presents and ignores empty regions" {
    var arena_state = std.heap.ArenaAllocator.init(std.testing.allocator);
    defer arena_state.deinit();

    // A bar fits a one-wide region only on its side, and never in 2x2.
    const answers = try solve(arena_state.allocator(), bar_present ++ "1x3: 1\n3x1: 1\n2x2: 1\n1x1: 0\n");
    try std.testing.expectEqual(@as(usize, 3), answers.part1);
}

test "rejects malformed input" {
    var arena_state = std.heap.ArenaAllocator.init(std.testing.allocator);
    defer arena_state.deinit();
    const arena = arena_state.allocator();

    try std.testing.expectError(error.MissingPresents, solve(arena, ""));
    try std.testing.expectError(error.MissingPresents, solve(arena, "4x4: 1\n"));
    try std.testing.expectError(error.InvalidPresentRow, solve(arena, "0:\n##\n##.\n##.\n\n4x4: 1\n"));
    try std.testing.expectError(error.InvalidPresentCell, solve(arena, "0:\n###\n#x.\n###\n\n4x4: 1\n"));
    try std.testing.expectError(error.EmptyPresent, solve(arena, "0:\n...\n...\n...\n\n4x4: 1\n"));
    try std.testing.expectError(error.TruncatedPresent, solve(arena, "0:\n###\n###"));
    try std.testing.expectError(error.PresentOutOfOrder, solve(arena, "1:\n###\n###\n###\n\n4x4: 1\n"));
    try std.testing.expectError(error.PresentAfterRegion, solve(arena, c_present ++ "4x4: 1\n1:\n###\n###\n###\n"));
    try std.testing.expectError(error.CountMismatch, solve(arena, c_present ++ "4x4: 1 2\n"));
    try std.testing.expectError(error.InvalidRegion, solve(arena, c_present ++ "4by4: 1\n"));
    try std.testing.expectError(error.InvalidRegion, solve(arena, c_present ++ "0x4: 1\n"));
    try std.testing.expectError(error.InvalidRegion, solve(arena, c_present ++ "4x4: x\n"));
}
