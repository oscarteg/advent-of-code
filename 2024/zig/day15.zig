const std = @import("std");

const Answers = struct {
    part1: i64,
    part2: i64,
};

const Vector2 = struct {
    row: isize,
    col: isize,

    fn plus(self: Vector2, other: Vector2) Vector2 {
        return .{ .row = self.row + other.row, .col = self.col + other.col };
    }

    fn minus(self: Vector2, other: Vector2) Vector2 {
        return .{ .row = self.row - other.row, .col = self.col - other.col };
    }

    fn equals(self: Vector2, other: Vector2) bool {
        return self.row == other.row and self.col == other.col;
    }
};

const right_half = Vector2{ .row = 0, .col = 1 };
const left_half = Vector2{ .row = 0, .col = -1 };

/// A rectangular warehouse map enclosed by walls, stored row-major, plus the
/// robot's position. That wall border is what makes every push stop instead of
/// leaving the map.
const Warehouse = struct {
    tiles: []u8,
    width: usize,
    height: usize,
    robot: Vector2,

    fn deinit(self: Warehouse, allocator: std.mem.Allocator) void {
        allocator.free(self.tiles);
    }

    fn indexOf(self: Warehouse, at: Vector2) usize {
        return @as(usize, @intCast(at.row)) * self.width + @as(usize, @intCast(at.col));
    }

    fn tile(self: Warehouse, at: Vector2) u8 {
        return self.tiles[self.indexOf(at)];
    }

    fn setTile(self: *Warehouse, at: Vector2, value: u8) void {
        self.tiles[self.indexOf(at)] = value;
    }

    fn isEnclosed(self: Warehouse) bool {
        for (0..self.width) |col| {
            if (self.tiles[col] != '#') return false;
            if (self.tiles[(self.height - 1) * self.width + col] != '#') return false;
        }
        for (0..self.height) |row| {
            if (self.tiles[row * self.width] != '#') return false;
            if (self.tiles[row * self.width + self.width - 1] != '#') return false;
        }
        return true;
    }

    fn gpsSum(self: Warehouse, box: u8) i64 {
        var sum: i64 = 0;
        for (0..self.height) |row| {
            for (0..self.width) |col| {
                if (self.tiles[row * self.width + col] == box) {
                    sum += @intCast(100 * row + col);
                }
            }
        }
        return sum;
    }
};

const Puzzle = struct {
    warehouse: Warehouse,
    moves: []const Vector2,

    fn deinit(self: Puzzle, allocator: std.mem.Allocator) void {
        self.warehouse.deinit(allocator);
        allocator.free(self.moves);
    }
};

fn directionOf(move: u8) !Vector2 {
    return switch (move) {
        '^' => .{ .row = -1, .col = 0 },
        'v' => .{ .row = 1, .col = 0 },
        '<' => .{ .row = 0, .col = -1 },
        '>' => .{ .row = 0, .col = 1 },
        else => error.InvalidMove,
    };
}

/// The map comes first; the blank line after it starts the move sequence, whose
/// own line breaks carry no meaning.
fn parse(allocator: std.mem.Allocator, input: []const u8) !Puzzle {
    var tiles: std.ArrayList(u8) = .empty;
    errdefer tiles.deinit(allocator);
    var moves: std.ArrayList(Vector2) = .empty;
    errdefer moves.deinit(allocator);

    var width: usize = 0;
    var height: usize = 0;
    var robot: ?Vector2 = null;
    var reading_map = true;

    var lines = std.mem.splitScalar(u8, input, '\n');
    while (lines.next()) |raw| {
        const line = std.mem.trim(u8, raw, " \r");
        if (line.len == 0) {
            reading_map = reading_map and height == 0;
            continue;
        }
        if (!reading_map) {
            for (line) |move| try moves.append(allocator, try directionOf(move));
            continue;
        }
        if (height == 0) {
            width = line.len;
        } else if (line.len != width) {
            return error.RaggedWarehouse;
        }
        for (line, 0..) |current, col| switch (current) {
            '#', '.', 'O' => {},
            '@' => {
                if (robot != null) return error.TooManyRobots;
                robot = .{ .row = @intCast(height), .col = @intCast(col) };
            },
            else => return error.InvalidTile,
        };
        try tiles.appendSlice(allocator, line);
        height += 1;
    }

    if (height == 0) return error.EmptyWarehouse;
    if (moves.items.len == 0) return error.NoMoves;
    const robot_at = robot orelse return error.NoRobot;

    const warehouse = Warehouse{
        .tiles = try tiles.toOwnedSlice(allocator),
        .width = width,
        .height = height,
        .robot = robot_at,
    };
    errdefer warehouse.deinit(allocator);
    if (!warehouse.isEnclosed()) return error.OpenWarehouse;

    return .{ .warehouse = warehouse, .moves = try moves.toOwnedSlice(allocator) };
}

/// Part 1: the robot shoves the whole run of boxes ahead of it one tile along
/// the move, and only moves at all when the tile past that run is free. Hopping
/// the near box to the far end leaves the same map as shifting every box.
fn push(warehouse: *Warehouse, direction: Vector2) void {
    const target = warehouse.robot.plus(direction);
    var free = target;
    while (warehouse.tile(free) == 'O') free = free.plus(direction);
    if (warehouse.tile(free) == '#') return;

    warehouse.setTile(free, 'O');
    warehouse.setTile(target, '@');
    warehouse.setTile(warehouse.robot, '.');
    warehouse.robot = target;
}

/// Part 2 doubles every tile's width: walls and floor repeat, a box becomes the
/// pair "[]", and the robot keeps a single tile of its own.
fn widen(allocator: std.mem.Allocator, narrow: Warehouse) !Warehouse {
    const tiles = try allocator.alloc(u8, narrow.tiles.len * 2);
    for (narrow.tiles, 0..) |current, index| {
        tiles[index * 2] = if (current == 'O') '[' else current;
        tiles[index * 2 + 1] = switch (current) {
            'O' => ']',
            '@' => '.',
            else => current,
        };
    }
    return .{
        .tiles = tiles,
        .width = narrow.width * 2,
        .height = narrow.height,
        .robot = .{ .row = narrow.robot.row, .col = narrow.robot.col * 2 },
    };
}

/// Sideways the halves stay paired, so the whole run shifts one tile.
fn pushWideRow(warehouse: *Warehouse, direction: Vector2) void {
    const target = warehouse.robot.plus(direction);
    var free = target;
    while (warehouse.tile(free) == '[' or warehouse.tile(free) == ']') free = free.plus(direction);
    if (warehouse.tile(free) == '#') return;

    var at = free;
    while (!at.equals(target)) {
        const previous = at.minus(direction);
        warehouse.setTile(at, warehouse.tile(previous));
        at = previous;
    }
    warehouse.setTile(target, '@');
    warehouse.setTile(warehouse.robot, '.');
    warehouse.robot = target;
}

fn appendBox(allocator: std.mem.Allocator, boxes: *std.ArrayList(Vector2), left: Vector2) !void {
    for (boxes.items) |existing| {
        if (existing.equals(left)) return;
    }
    try boxes.append(allocator, left);
}

/// Collects the left halves of every box a vertical push would move, nearest
/// first; returns false when a wall blocks any of them. A box can be reached
/// through either half and the list stays small, so a linear duplicate check
/// beats tracking a grid of visited tiles.
fn collectPushedBoxes(
    allocator: std.mem.Allocator,
    warehouse: Warehouse,
    direction: Vector2,
    boxes: *std.ArrayList(Vector2),
) !bool {
    const target = warehouse.robot.plus(direction);
    switch (warehouse.tile(target)) {
        '#' => return false,
        '[' => try boxes.append(allocator, target),
        ']' => try boxes.append(allocator, target.plus(left_half)),
        else => return true,
    }

    var index: usize = 0;
    while (index < boxes.items.len) : (index += 1) {
        const left = boxes.items[index];
        for ([2]Vector2{ left, left.plus(right_half) }) |half| {
            const next = half.plus(direction);
            switch (warehouse.tile(next)) {
                '#' => return false,
                '[' => try appendBox(allocator, boxes, next),
                ']' => try appendBox(allocator, boxes, next.plus(left_half)),
                else => {},
            }
        }
    }
    return true;
}

fn pushWideColumn(allocator: std.mem.Allocator, warehouse: *Warehouse, direction: Vector2) !void {
    var boxes: std.ArrayList(Vector2) = .empty;
    defer boxes.deinit(allocator);
    if (!try collectPushedBoxes(allocator, warehouse.*, direction, &boxes)) return;

    // Farthest box first, so each one lands on tiles the previous box left.
    var index = boxes.items.len;
    while (index > 0) {
        index -= 1;
        const left = boxes.items[index];
        const right = left.plus(right_half);
        warehouse.setTile(left, '.');
        warehouse.setTile(right, '.');
        warehouse.setTile(left.plus(direction), '[');
        warehouse.setTile(right.plus(direction), ']');
    }
    warehouse.setTile(warehouse.robot, '.');
    warehouse.robot = warehouse.robot.plus(direction);
    warehouse.setTile(warehouse.robot, '@');
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var puzzle = try parse(allocator, input);
    defer puzzle.deinit(allocator);

    var wide = try widen(allocator, puzzle.warehouse);
    defer wide.deinit(allocator);

    for (puzzle.moves) |direction| {
        push(&puzzle.warehouse, direction);
        if (direction.row == 0) {
            pushWideRow(&wide, direction);
        } else {
            try pushWideColumn(allocator, &wide, direction);
        }
    }

    return .{
        .part1 = puzzle.warehouse.gpsSum('O'),
        .part2 = wide.gpsSum('['),
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
    try stdout.print("Part 1: {d}\nPart 2: {d}\n", .{ answers.part1, answers.part2 });
    try stdout.flush();
}

test "official examples" {
    // The puzzle states only the narrow sum for the smaller map.
    const small = try solve(std.testing.allocator,
        \\########
        \\#..O.O.#
        \\##@.O..#
        \\#...O..#
        \\#.#.O..#
        \\#...O..#
        \\#......#
        \\########
        \\
        \\<^^>>>vv<v>>v<<
    );
    try std.testing.expectEqual(@as(i64, 2028), small.part1);

    const large = try solve(std.testing.allocator, @embedFile("examples/day15.txt"));
    try std.testing.expectEqual(@as(i64, 10092), large.part1);
    try std.testing.expectEqual(@as(i64, 9021), large.part2);
}

test "widened boxes fan out on a vertical push" {
    // The puzzle states only the widened sum for this map.
    const answers = try solve(std.testing.allocator,
        \\#######
        \\#...#.#
        \\#.....#
        \\#..OO@#
        \\#..O..#
        \\#.....#
        \\#######
        \\
        \\<vv<<^^<<^^
    );
    try std.testing.expectEqual(@as(i64, 618), answers.part2);
}

test "a push into a wall moves nothing" {
    // The second push runs into the wall and leaves the narrow map untouched,
    // while the widened map has one floor tile more and takes both pushes.
    const answers = try solve(std.testing.allocator,
        \\#####
        \\#.O@#
        \\#####
        \\
        \\<<
    );
    try std.testing.expectEqual(@as(i64, 101), answers.part1);
    try std.testing.expectEqual(@as(i64, 102), answers.part2);
}

test "malformed input returns an error" {
    const allocator = std.testing.allocator;
    try std.testing.expectError(error.EmptyWarehouse, solve(allocator, ""));
    try std.testing.expectError(error.RaggedWarehouse, solve(allocator, "#####\n###\n\n<\n"));
    try std.testing.expectError(error.InvalidTile, solve(allocator, "#####\n#.X@#\n#####\n\n<\n"));
    try std.testing.expectError(error.NoRobot, solve(allocator, "#####\n#...#\n#####\n\n<\n"));
    try std.testing.expectError(error.TooManyRobots, solve(allocator, "#####\n#@@.#\n#####\n\n<\n"));
    try std.testing.expectError(error.OpenWarehouse, solve(allocator, "#####\n#.@..\n#####\n\n<\n"));
    try std.testing.expectError(error.NoMoves, solve(allocator, "#####\n#.@.#\n#####\n"));
    try std.testing.expectError(error.InvalidMove, solve(allocator, "#####\n#.@.#\n#####\n\n<x\n"));
}
