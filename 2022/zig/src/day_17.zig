const std = @import("std");

const chamber_width: i32 = 7;
const part1_rocks: u64 = 2022;
const part2_rocks: u64 = 1_000_000_000_000;
const profile_depth: usize = 50;

const Rock = struct { rows: []const u8, width: i32 };

const rocks = [5]Rock{
    .{ .rows = &[_]u8{0b1111}, .width = 4 },
    .{ .rows = &[_]u8{ 0b010, 0b111, 0b010 }, .width = 3 },
    .{ .rows = &[_]u8{ 0b111, 0b100, 0b100 }, .width = 3 },
    .{ .rows = &[_]u8{ 0b1, 0b1, 0b1, 0b1 }, .width = 1 },
    .{ .rows = &[_]u8{ 0b11, 0b11 }, .width = 2 },
};

const ProfileKey = struct {
    rock_idx: usize,
    jet_idx: usize,
    profile: [profile_depth]u8,
};

const SeenEntry = struct { rock_count: u64, height: u64 };

fn parseJets(allocator: std.mem.Allocator, input: []const u8) ![]i32 {
    const trimmed = std.mem.trim(u8, input, " \t\r\n");
    if (trimmed.len == 0) return error.EmptyJetPattern;
    var jets = try allocator.alloc(i32, trimmed.len);
    errdefer allocator.free(jets);
    for (trimmed, 0..) |c, i| {
        jets[i] = switch (c) {
            '<' => -1,
            '>' => 1,
            else => return error.InvalidJetCharacter,
        };
    }
    return jets;
}

fn collides(shape: Rock, x: i32, y: i32, rows: []const u8) bool {
    if (x < 0 or x + shape.width > chamber_width or y < 0) return true;
    for (shape.rows, 0..) |row_mask, i| {
        const row_y = y + @as(i32, @intCast(i));
        if (row_y < @as(i32, @intCast(rows.len))) {
            const shifted: u16 = @as(u16, row_mask) << @as(u4, @intCast(x));
            if (@as(u16, rows[@intCast(row_y)]) & shifted != 0) return true;
        }
    }
    return false;
}

fn settle(allocator: std.mem.Allocator, shape: Rock, x: i32, y: i32, rows: *std.ArrayList(u8)) !void {
    for (shape.rows, 0..) |row_mask, i| {
        const row_y: usize = @intCast(y + @as(i32, @intCast(i)));
        while (rows.items.len <= row_y) try rows.append(allocator, 0);
        rows.items[row_y] |= row_mask << @as(u3, @intCast(x));
    }
}

fn surfaceProfile(rows: []const u8) [profile_depth]u8 {
    var profile: [profile_depth]u8 = [_]u8{0xFF} ** profile_depth;
    const depth = @min(profile_depth, rows.len);
    const start = rows.len - depth;
    @memcpy(profile[profile_depth - depth ..], rows[start..]);
    return profile;
}

fn towerHeight(allocator: std.mem.Allocator, jets: []const i32, target: u64) !u64 {
    var rows: std.ArrayList(u8) = .empty;
    defer rows.deinit(allocator);
    var heights: std.ArrayList(u64) = .empty;
    defer heights.deinit(allocator);
    try heights.append(allocator, 0);

    var seen = std.AutoHashMap(ProfileKey, SeenEntry).init(allocator);
    defer seen.deinit();

    var jet_index: usize = 0;
    var n: u64 = 0;
    while (n < target) {
        const rock_idx: usize = @intCast(n % 5);
        const shape = rocks[rock_idx];
        var x: i32 = 2;
        var y: i32 = @as(i32, @intCast(rows.items.len)) + 3;
        while (true) {
            const dir = jets[jet_index];
            jet_index = (jet_index + 1) % jets.len;
            const nx = x + dir;
            if (!collides(shape, nx, y, rows.items)) x = nx;
            if (!collides(shape, x, y - 1, rows.items)) {
                y -= 1;
            } else {
                try settle(allocator, shape, x, y, &rows);
                break;
            }
        }
        n += 1;
        try heights.append(allocator, rows.items.len);

        const key = ProfileKey{ .rock_idx = rock_idx, .jet_idx = jet_index, .profile = surfaceProfile(rows.items) };
        if (seen.get(key)) |prev| {
            const cycle_len = n - prev.rock_count;
            const cycle_gain = rows.items.len - prev.height;
            const remaining = target - n;
            const cycles = remaining / cycle_len;
            const remainder = remaining % cycle_len;
            const remainder_gain = heights.items[@intCast(prev.rock_count + remainder)] - prev.height;
            return @as(u64, rows.items.len) + cycles * cycle_gain + remainder_gain;
        }
        try seen.put(key, SeenEntry{ .rock_count = n, .height = rows.items.len });
    }
    return rows.items.len;
}

const Answers = struct { part1: u64, part2: u64 };

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const jets = try parseJets(allocator, input);
    defer allocator.free(jets);
    const part1 = try towerHeight(allocator, jets, part1_rocks);
    const part2 = try towerHeight(allocator, jets, part2_rocks);
    return Answers{ .part1 = part1, .part2 = part2 };
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

test "official_example" {
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day_17.txt");
    const answers = try solve(allocator, input);
    try std.testing.expectEqual(3068, answers.part1);
    try std.testing.expectEqual(1514285714288, answers.part2);
}

test "single_jet_pushes_rock_to_the_wall" {
    const allocator = std.testing.allocator;
    const jets = try parseJets(allocator, "<");
    defer allocator.free(jets);
    const height = try towerHeight(allocator, jets, 1);
    try std.testing.expectEqual(1, height);
}

test "malformed_jet_pattern_is_rejected" {
    const allocator = std.testing.allocator;
    try std.testing.expectError(error.InvalidJetCharacter, solve(allocator, "abc"));
    try std.testing.expectError(error.EmptyJetPattern, solve(allocator, ""));
}
