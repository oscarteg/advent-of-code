const std = @import("std");
const testing = std.testing;

const real_row: i64 = 2_000_000;
const real_search_limit: i64 = 4_000_000;
const tuning_multiplier: i64 = 4_000_000;

const Answers = struct {
    part1: i64,
    part2: i64,
};

const Error = error{
    MalformedReading,
    NoReadings,
    NoUncoveredPosition,
};

/// An inclusive range of x positions.
const Span = struct {
    start: i64,
    end: i64,

    fn lessThan(_: void, a: Span, b: Span) bool {
        return a.start < b.start;
    }
};

/// A sensor, its closest beacon and the Manhattan distance between the two.
const Reading = struct {
    sensor_x: i64,
    sensor_y: i64,
    beacon_x: i64,
    beacon_y: i64,
    radius: i64,

    /// The inclusive x-range this sensor rules out on the given row, if it reaches it.
    fn rowSpan(self: Reading, row: i64) ?Span {
        const spread = self.radius - absDiff(self.sensor_y, row);
        if (spread < 0) return null;
        return .{ .start = self.sensor_x - spread, .end = self.sensor_x + spread };
    }
};

fn absDiff(a: i64, b: i64) i64 {
    return if (a < b) b - a else a - b;
}

/// Consumes the prefix and the integer up to the terminator, advancing rest past it.
fn takeValue(rest: *[]const u8, prefix: []const u8, terminator: []const u8) Error!i64 {
    if (!std.mem.startsWith(u8, rest.*, prefix)) return Error.MalformedReading;
    const tail = rest.*[prefix.len..];
    const end = std.mem.indexOf(u8, tail, terminator) orelse return Error.MalformedReading;
    const value = std.fmt.parseInt(i64, tail[0..end], 10) catch return Error.MalformedReading;
    rest.* = tail[end + terminator.len ..];
    return value;
}

fn parseReading(line: []const u8) Error!Reading {
    var rest = line;
    const sensor_x = try takeValue(&rest, "Sensor at x=", ", ");
    const sensor_y = try takeValue(&rest, "y=", ": ");
    const beacon_x = try takeValue(&rest, "closest beacon is at x=", ", ");
    if (!std.mem.startsWith(u8, rest, "y=")) return Error.MalformedReading;
    const beacon_y = std.fmt.parseInt(i64, rest[2..], 10) catch return Error.MalformedReading;
    return .{
        .sensor_x = sensor_x,
        .sensor_y = sensor_y,
        .beacon_x = beacon_x,
        .beacon_y = beacon_y,
        .radius = absDiff(sensor_x, beacon_x) + absDiff(sensor_y, beacon_y),
    };
}

fn parse(allocator: std.mem.Allocator, input: []const u8) ![]Reading {
    var readings: std.ArrayList(Reading) = .empty;
    errdefer readings.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;
        try readings.append(allocator, try parseReading(trimmed));
    }
    if (readings.items.len == 0) return Error.NoReadings;
    return readings.toOwnedSlice(allocator);
}

/// Fills scratch with the disjoint x-ranges covering the row, sorted by start and
/// with touching ranges joined, so that any boundary left between them is a real
/// gap. Scratch must have room for one span per reading.
fn mergedRowSpans(scratch: []Span, readings: []const Reading, row: i64) []Span {
    var count: usize = 0;
    for (readings) |reading| {
        if (reading.rowSpan(row)) |span| {
            scratch[count] = span;
            count += 1;
        }
    }

    const spans = scratch[0..count];
    if (spans.len == 0) return spans;
    std.mem.sort(Span, spans, {}, Span.lessThan);

    var write: usize = 0;
    for (spans[1..]) |span| {
        if (span.start <= spans[write].end + 1) {
            spans[write].end = @max(spans[write].end, span.end);
        } else {
            write += 1;
            spans[write] = span;
        }
    }
    return spans[0 .. write + 1];
}

/// Positions on the row that cannot hold a beacon: every covered position minus the
/// distinct known beacons there, which each sit inside their own sensor's range.
fn excludedPositions(scratch: []Span, readings: []const Reading, row: i64) i64 {
    var covered: i64 = 0;
    for (mergedRowSpans(scratch, readings, row)) |span| covered += span.end - span.start + 1;

    var beacons: i64 = 0;
    for (readings, 0..) |reading, index| {
        if (reading.beacon_y != row) continue;
        var seen = false;
        for (readings[0..index]) |earlier| {
            seen = seen or (earlier.beacon_y == row and earlier.beacon_x == reading.beacon_x);
        }
        if (!seen) beacons += 1;
    }
    return covered - beacons;
}

/// The smallest x between 0 and the limit on the given row that no sensor covers.
fn uncoveredX(scratch: []Span, readings: []const Reading, row: i64, limit: i64) ?i64 {
    var x: i64 = 0;
    for (mergedRowSpans(scratch, readings, row)) |span| {
        if (span.start > x) return x;
        x = @max(x, span.end + 1);
        if (x > limit) return null;
    }
    return if (x <= limit) x else null;
}

/// Tuning frequency of the one position in the search square no sensor covers.
fn tuningFrequency(scratch: []Span, readings: []const Reading, limit: i64) Error!i64 {
    var row: i64 = 0;
    while (row <= limit) : (row += 1) {
        if (uncoveredX(scratch, readings, row, limit)) |x| return x * tuning_multiplier + row;
    }
    return Error.NoUncoveredPosition;
}

fn solveWith(allocator: std.mem.Allocator, input: []const u8, row: i64, search_limit: i64) !Answers {
    const readings = try parse(allocator, input);
    defer allocator.free(readings);
    const scratch = try allocator.alloc(Span, readings.len);
    defer allocator.free(scratch);

    return .{
        .part1 = excludedPositions(scratch, readings, row),
        .part2 = try tuningFrequency(scratch, readings, search_limit),
    };
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    return solveWith(allocator, input, real_row, real_search_limit);
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

const example = @embedFile("examples/day_15.txt");
const example_row: i64 = 10;
const example_search_limit: i64 = 20;

test "official example" {
    const answers = try solveWith(testing.allocator, example, example_row, example_search_limit);
    try testing.expectEqual(@as(i64, 26), answers.part1);
    try testing.expectEqual(@as(i64, 56_000_011), answers.part2);
}

test "solve uses the real puzzle row" {
    const answers = try solve(testing.allocator, example);
    try testing.expectEqual(@as(i64, 0), answers.part1);
}

test "single sensor rules out its diamond row" {
    // Row 0 covers x in -2...2; the known beacon at x=2 stays a valid position,
    // and row 1 leaves x=2 uncovered.
    const answers = try solveWith(testing.allocator, "Sensor at x=0, y=0: closest beacon is at x=2, y=0", 0, 2);
    try testing.expectEqual(@as(i64, 4), answers.part1);
    try testing.expectEqual(@as(i64, 8_000_001), answers.part2);
}

test "row outside every diamond rules out nothing" {
    const answers = try solveWith(testing.allocator, example, 1_000, example_search_limit);
    try testing.expectEqual(@as(i64, 0), answers.part1);
}

test "fully covered search area has no distress beacon" {
    const input = "Sensor at x=0, y=0: closest beacon is at x=10, y=0";
    try testing.expectError(Error.NoUncoveredPosition, solveWith(testing.allocator, input, 0, 2));
}

test "empty input is rejected" {
    try testing.expectError(Error.NoReadings, solve(testing.allocator, "\n\n"));
}

test "malformed readings are rejected" {
    try testing.expectError(Error.MalformedReading, parseReading("Sensor at x=2, y=18"));
    try testing.expectError(Error.MalformedReading, parseReading("Sensor at x=a, y=18: closest beacon is at x=-2, y=15"));
    try testing.expectError(Error.MalformedReading, parseReading("at x=2, y=18: closest beacon is at x=-2, y=15"));
}
