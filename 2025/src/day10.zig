const std = @import("std");

/// Factory: every line describes one machine with an indicator light diagram
/// in `[]`, button wiring schematics in `()` and joltage requirements in `{}`.
/// Part 1 toggles the lights into the diagram with as few presses as possible;
/// part 2 ignores the lights and counts the presses that make every joltage
/// counter reach its target exactly.
const Answers = struct { part1: u64, part2: u64 };

/// The indicator diagram is held in a `u64` bitmask.
const max_lights = 64;
/// Button subsets are enumerated as `u64` bitmasks, so `2^buttons` has to fit.
const max_buttons = 63;

const Machine = struct {
    lights: usize,
    /// Bit `j` is set when light `j` has to end up on.
    target: u64,
    /// Per button: the lights it toggles and the counters it increments.
    buttons: []const []const usize,
    joltage: []const i64,
};

/// Strips the brackets off a field like `[.##.]`, `(0,2)` or `{3,5}`.
fn unwrap(field: []const u8, open: u8, close: u8) ![]const u8 {
    if (field.len < 3 or field[0] != open or field[field.len - 1] != close) return error.InvalidMachine;
    return field[1 .. field.len - 1];
}

fn parseIndices(allocator: std.mem.Allocator, field: []const u8, lights: usize) ![]const usize {
    var indices: std.ArrayList(usize) = .empty;
    errdefer indices.deinit(allocator);

    var parts = std.mem.tokenizeScalar(u8, field, ',');
    while (parts.next()) |part| {
        const index = try std.fmt.parseInt(usize, part, 10);
        if (index >= lights) return error.LightOutOfRange;
        try indices.append(allocator, index);
    }
    if (indices.items.len == 0) return error.InvalidMachine;

    return indices.toOwnedSlice(allocator);
}

fn parseMachine(allocator: std.mem.Allocator, line: []const u8) !Machine {
    var fields: std.ArrayList([]const u8) = .empty;
    defer fields.deinit(allocator);
    var tokens = std.mem.tokenizeAny(u8, line, " \t");
    while (tokens.next()) |token| try fields.append(allocator, token);
    if (fields.items.len < 3) return error.InvalidMachine;

    const diagram = try unwrap(fields.items[0], '[', ']');
    if (diagram.len > max_lights) return error.TooManyLights;
    var target: u64 = 0;
    for (diagram, 0..) |cell, light| switch (cell) {
        '.' => {},
        '#' => target |= @as(u64, 1) << @intCast(light),
        else => return error.InvalidMachine,
    };

    const button_fields = fields.items[1 .. fields.items.len - 1];
    if (button_fields.len > max_buttons) return error.TooManyButtons;
    const buttons = try allocator.alloc([]const usize, button_fields.len);
    for (button_fields, buttons) |field, *button| {
        button.* = try parseIndices(allocator, try unwrap(field, '(', ')'), diagram.len);
    }

    var joltage: std.ArrayList(i64) = .empty;
    errdefer joltage.deinit(allocator);
    var requirements = std.mem.tokenizeScalar(u8, try unwrap(fields.items[fields.items.len - 1], '{', '}'), ',');
    while (requirements.next()) |requirement| {
        const value = try std.fmt.parseInt(i64, requirement, 10);
        if (value < 0) return error.NegativeJoltage;
        try joltage.append(allocator, value);
    }
    if (joltage.items.len != diagram.len) return error.JoltageCountMismatch;

    return .{
        .lights = diagram.len,
        .target = target,
        .buttons = buttons,
        .joltage = try joltage.toOwnedSlice(allocator),
    };
}

/// Fewest presses that toggle the all-off lights into the diagram. Pressing a
/// button twice cancels out, so every button is pressed once or not at all and
/// the lights form a system over GF(2). Elimination yields one solution plus a
/// basis of the null space - button sets that toggle nothing - and the answer
/// is the lightest solution in that coset.
fn fewestToggles(allocator: std.mem.Allocator, machine: Machine) !u64 {
    const Row = struct { lights: u64, buttons: u64 };
    const rows = try allocator.alloc(Row, machine.buttons.len);
    defer allocator.free(rows);
    for (machine.buttons, rows, 0..) |button, *row, index| {
        var lights: u64 = 0;
        for (button) |light| lights |= @as(u64, 1) << @intCast(light);
        row.* = .{ .lights = lights, .buttons = @as(u64, 1) << @intCast(index) };
    }

    const pivot_lights = try allocator.alloc(usize, rows.len);
    defer allocator.free(pivot_lights);
    var rank: usize = 0;
    for (0..machine.lights) |light| {
        if (rank == rows.len) break;
        const bit = @as(u64, 1) << @intCast(light);
        const pivot = for (rows[rank..], rank..) |row, index| {
            if (row.lights & bit != 0) break index;
        } else continue;

        std.mem.swap(Row, &rows[rank], &rows[pivot]);
        for (rows, 0..) |*row, index| {
            if (index == rank or row.lights & bit == 0) continue;
            row.lights ^= rows[rank].lights;
            row.buttons ^= rows[rank].buttons;
        }
        pivot_lights[rank] = light;
        rank += 1;
    }

    var residual = machine.target;
    var pressed: u64 = 0;
    for (pivot_lights[0..rank], rows[0..rank]) |light, row| {
        if (residual & (@as(u64, 1) << @intCast(light)) == 0) continue;
        residual ^= row.lights;
        pressed ^= row.buttons;
    }
    if (residual != 0) return error.LightsUnreachable;

    const basis = rows[rank..];
    var best: u64 = @popCount(pressed);
    var combination: u64 = 1;
    while (combination < @as(u64, 1) << @intCast(basis.len)) : (combination += 1) {
        var buttons = pressed;
        for (basis, 0..) |row, index| {
            if (combination & (@as(u64, 1) << @intCast(index)) != 0) buttons ^= row.buttons;
        }
        best = @min(best, @popCount(buttons));
    }
    return best;
}

/// A `lights` by `buttons + 1` matrix: one row per joltage counter, the last
/// column holding that counter's target.
const System = struct {
    cells: []i64,
    columns: usize,

    fn row(self: System, index: usize) []i64 {
        return self.cells[index * self.columns ..][0..self.columns];
    }
};

/// Depth-first search over the free presses of an eliminated system; the
/// pivot presses follow from them and the total is pruned against the best
/// solution found so far.
const PressSearch = struct {
    system: System,
    pivot_rows: []const ?usize,
    pivot_columns: []const usize,
    free_columns: []const usize,
    bounds: []const i64,
    presses: []i64,
    best: i64,

    fn run(self: *PressSearch, index: usize, assigned: i64) void {
        if (assigned >= self.best) return;
        if (index == self.free_columns.len) {
            var total = assigned;
            for (self.pivot_columns) |column| {
                const row = self.system.row(self.pivot_rows[column].?);
                var remainder = row[self.system.columns - 1];
                for (self.free_columns) |free| remainder -= row[free] * self.presses[free];
                if (@rem(remainder, row[column]) != 0) return;
                const presses = @divExact(remainder, row[column]);
                if (presses < 0) return;
                total += presses;
            }
            self.best = @min(self.best, total);
            return;
        }

        const column = self.free_columns[index];
        var presses: i64 = 0;
        while (presses <= self.bounds[index]) : (presses += 1) {
            self.presses[column] = presses;
            self.run(index + 1, assigned + presses);
        }
        self.presses[column] = 0;
    }
};

fn gcd(a: i64, b: i64) i64 {
    var x = a;
    var y = b;
    while (y != 0) {
        const remainder = @rem(x, y);
        x = y;
        y = remainder;
    }
    return x;
}

fn rowGcd(row: []const i64) i64 {
    var divisor: i64 = 0;
    for (row) |cell| divisor = gcd(divisor, @intCast(@abs(cell)));
    return divisor;
}

/// Fewest presses that make every counter reach its joltage target exactly.
/// A press adds one to each counter its button feeds, so the targets are a
/// linear system `A x = joltage` to be solved over the nonnegative integers.
/// Integer-preserving Gaussian elimination turns the pivot presses into a
/// function of the free ones; each free press is bounded by the smallest
/// target its button feeds, so a depth-first search over that box finds the
/// cheapest solution.
fn fewestPresses(allocator: std.mem.Allocator, machine: Machine) !u64 {
    const buttons = machine.buttons.len;
    const system = System{
        .cells = try allocator.alloc(i64, machine.lights * (buttons + 1)),
        .columns = buttons + 1,
    };
    defer allocator.free(system.cells);
    @memset(system.cells, 0);
    for (machine.buttons, 0..) |button, column| {
        for (button) |light| system.row(light)[column] = 1;
    }
    for (machine.joltage, 0..) |requirement, light| system.row(light)[buttons] = requirement;

    const pivot_rows = try allocator.alloc(?usize, buttons);
    defer allocator.free(pivot_rows);
    @memset(pivot_rows, null);
    var rank: usize = 0;
    for (0..buttons) |column| {
        if (rank == machine.lights) break;
        const pivot = for (rank..machine.lights) |row| {
            if (system.row(row)[column] != 0) break row;
        } else continue;

        for (system.row(rank), system.row(pivot)) |*a, *b| std.mem.swap(i64, a, b);
        const pivot_value = system.row(rank)[column];
        for (0..machine.lights) |row| {
            if (row == rank) continue;
            const factor = system.row(row)[column];
            if (factor == 0) continue;
            for (system.row(row), system.row(rank)) |*cell, pivot_cell| {
                cell.* = cell.* * pivot_value - pivot_cell * factor;
            }
            // Keep the entries small; the row stays an equivalent equation.
            const divisor = rowGcd(system.row(row));
            if (divisor > 1) for (system.row(row)) |*cell| {
                cell.* = @divExact(cell.*, divisor);
            };
        }
        pivot_rows[column] = rank;
        rank += 1;
    }

    for (0..machine.lights) |row| {
        const cells = system.row(row);
        if (std.mem.allEqual(i64, cells[0..buttons], 0) and cells[buttons] != 0) {
            return error.JoltageUnreachable;
        }
    }

    var pivot_columns: std.ArrayList(usize) = .empty;
    defer pivot_columns.deinit(allocator);
    var free_columns: std.ArrayList(usize) = .empty;
    defer free_columns.deinit(allocator);
    var bounds: std.ArrayList(i64) = .empty;
    defer bounds.deinit(allocator);
    for (pivot_rows, 0..) |pivot_row, column| {
        if (pivot_row != null) {
            try pivot_columns.append(allocator, column);
            continue;
        }
        // A button cannot be pressed more often than the smallest target it feeds.
        var bound: i64 = std.math.maxInt(i64);
        for (machine.buttons[column]) |light| bound = @min(bound, machine.joltage[light]);
        try free_columns.append(allocator, column);
        try bounds.append(allocator, bound);
    }

    const presses = try allocator.alloc(i64, buttons);
    defer allocator.free(presses);
    @memset(presses, 0);
    var search = PressSearch{
        .system = system,
        .pivot_rows = pivot_rows,
        .pivot_columns = pivot_columns.items,
        .free_columns = free_columns.items,
        .bounds = bounds.items,
        .presses = presses,
        .best = std.math.maxInt(i64),
    };
    search.run(0, 0);
    if (search.best == std.math.maxInt(i64)) return error.JoltageUnreachable;

    return @intCast(search.best);
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var arena_state = std.heap.ArenaAllocator.init(allocator);
    defer arena_state.deinit();
    const arena = arena_state.allocator();

    var answers = Answers{ .part1 = 0, .part2 = 0 };
    var machines: usize = 0;
    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, &std.ascii.whitespace);
        if (trimmed.len == 0) continue;

        const machine = try parseMachine(arena, trimmed);
        answers.part1 += try fewestToggles(arena, machine);
        answers.part2 += try fewestPresses(arena, machine);
        machines += 1;
    }
    if (machines == 0) return error.NoMachines;

    return answers;
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
    const answers = try solve(std.testing.allocator, @embedFile("examples/day10.txt"));
    try std.testing.expectEqual(@as(u64, 7), answers.part1);
    try std.testing.expectEqual(@as(u64, 33), answers.part2);
}

test "per-machine presses of the official example" {
    const cases = [_]struct { line: []const u8, part1: u64, part2: u64 }{
        .{ .line = "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}", .part1 = 2, .part2 = 10 },
        .{ .line = "[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}", .part1 = 3, .part2 = 12 },
        .{ .line = "[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}", .part1 = 2, .part2 = 11 },
    };
    for (cases) |case| {
        const answers = try solve(std.testing.allocator, case.line);
        try std.testing.expectEqual(case.part1, answers.part1);
        try std.testing.expectEqual(case.part2, answers.part2);
    }
}

test "a configured machine with zero joltage needs no presses" {
    const answers = try solve(std.testing.allocator, "[..] (0) (1) {0,0}\n");
    try std.testing.expectEqual(@as(u64, 0), answers.part1);
    try std.testing.expectEqual(@as(u64, 0), answers.part2);
}

test "the only button is pressed once for lights and three times for joltage" {
    const answers = try solve(std.testing.allocator, "[##] (0,1) {3,3}\n");
    try std.testing.expectEqual(@as(u64, 1), answers.part1);
    try std.testing.expectEqual(@as(u64, 3), answers.part2);
}

test "surrounding whitespace and blank lines are ignored" {
    const answers = try solve(std.testing.allocator, "\n  [##] (0,1) {3,3}  \n\n");
    try std.testing.expectEqual(@as(u64, 1), answers.part1);
    try std.testing.expectEqual(@as(u64, 3), answers.part2);
}

test "rejects a machine without a diagram" {
    try std.testing.expectError(error.InvalidMachine, solve(std.testing.allocator, "(0) (1) {1,1}\n"));
}

test "rejects a machine without joltage requirements" {
    try std.testing.expectError(error.InvalidMachine, solve(std.testing.allocator, "[.#] (0) (1)\n"));
}

test "rejects a machine without buttons" {
    try std.testing.expectError(error.InvalidMachine, solve(std.testing.allocator, "[.#] {1,1}\n"));
}

test "rejects an unparseable light index" {
    try std.testing.expectError(error.InvalidCharacter, solve(std.testing.allocator, "[.#] (a) {1,1}\n"));
}

test "rejects a diagram with an unknown light state" {
    try std.testing.expectError(error.InvalidMachine, solve(std.testing.allocator, "[.x] (0) {1,1}\n"));
}

test "rejects a button wired past the last light" {
    try std.testing.expectError(error.LightOutOfRange, solve(std.testing.allocator, "[.#] (5) {1,1}\n"));
}

test "rejects a negative joltage requirement" {
    try std.testing.expectError(error.NegativeJoltage, solve(std.testing.allocator, "[.#] (0) {1,-1}\n"));
}

test "rejects joltage requirements that do not match the lights" {
    try std.testing.expectError(error.JoltageCountMismatch, solve(std.testing.allocator, "[.#] (0) {1}\n"));
}

test "rejects input without machines" {
    try std.testing.expectError(error.NoMachines, solve(std.testing.allocator, "\n  \n"));
}

test "rejects lights the buttons cannot reach" {
    try std.testing.expectError(error.LightsUnreachable, solve(std.testing.allocator, "[#.] (0,1) {1,1}\n"));
}

test "rejects joltage the buttons cannot reach" {
    try std.testing.expectError(error.JoltageUnreachable, solve(std.testing.allocator, "[#.] (0) {1,2}\n"));
}
