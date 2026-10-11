const std = @import("std");

/// Facings are the puzzle's own codes: 0 right, 1 down, 2 left, 3 up.
const Facing = u2;
const deltas = [4][2]i64{ .{ 1, 0 }, .{ 0, 1 }, .{ -1, 0 }, .{ 0, -1 } };

const Tile = enum { empty, open, wall };
const Turn = enum { left, right };
const Move = union(enum) { forward: u32, turn: Turn };
const Point = struct { x: i64, y: i64 };
const Crossing = struct { pos: Point, facing: Facing };
const Vec3 = [3]i64;
const Answers = struct { part1: i64, part2: i64 };

fn add(a: Vec3, b: Vec3) Vec3 {
    return .{ a[0] + b[0], a[1] + b[1], a[2] + b[2] };
}

fn subtract(a: Vec3, b: Vec3) Vec3 {
    return .{ a[0] - b[0], a[1] - b[1], a[2] - b[2] };
}

fn negate(a: Vec3) Vec3 {
    return .{ -a[0], -a[1], -a[2] };
}

fn scale(a: Vec3, factor: i64) Vec3 {
    return .{ a[0] * factor, a[1] * factor, a[2] * factor };
}

fn dot(a: Vec3, b: Vec3) i64 {
    return a[0] * b[0] + a[1] * b[1] + a[2] * b[2];
}

fn equal(a: Vec3, b: Vec3) bool {
    return a[0] == b[0] and a[1] == b[1] and a[2] == b[2];
}

/// Slot for one of the six axis-aligned unit normals.
fn axisIndex(normal: Vec3) usize {
    for (normal, 0..) |part, axis| {
        if (part != 0) return axis * 2 + @intFromBool(part < 0);
    }
    return 0;
}

const Notes = struct {
    tiles: []Tile,
    width: i64,
    height: i64,
    path: []Move,
    allocator: std.mem.Allocator,

    fn parse(allocator: std.mem.Allocator, input: []const u8) !Notes {
        var rows: std.ArrayList([]const u8) = .empty;
        defer rows.deinit(allocator);
        var path_text: ?[]const u8 = null;
        var lines = std.mem.splitScalar(u8, input, '\n');
        while (lines.next()) |raw| {
            const line = std.mem.trimEnd(u8, raw, "\r");
            if (line.len == 0) break;
            try rows.append(allocator, line);
        }
        while (lines.next()) |raw| {
            const line = std.mem.trimEnd(u8, raw, "\r");
            if (line.len != 0) {
                path_text = line;
                break;
            }
        }

        var width: usize = 0;
        for (rows.items) |row| width = @max(width, row.len);
        if (width == 0) return error.MissingMap;
        const text = path_text orelse return error.MissingPath;

        const height = rows.items.len;
        const tiles = try allocator.alloc(Tile, width * height);
        errdefer allocator.free(tiles);
        @memset(tiles, .empty);
        for (rows.items, 0..) |row, y| {
            for (row, 0..) |ch, x| {
                tiles[y * width + x] = switch (ch) {
                    ' ' => .empty,
                    '.' => .open,
                    '#' => .wall,
                    else => return error.InvalidMapCharacter,
                };
            }
        }
        return .{
            .tiles = tiles,
            .width = @intCast(width),
            .height = @intCast(height),
            .path = try parsePath(allocator, text),
            .allocator = allocator,
        };
    }

    fn deinit(self: *Notes) void {
        self.allocator.free(self.tiles);
        self.allocator.free(self.path);
    }

    fn tile(self: Notes, p: Point) Tile {
        if (p.x < 0 or p.y < 0 or p.x >= self.width or p.y >= self.height) return .empty;
        return self.tiles[@intCast(p.y * self.width + p.x)];
    }

    fn start(self: Notes) !Point {
        var x: i64 = 0;
        while (x < self.width) : (x += 1) {
            if (self.tile(.{ .x = x, .y = 0 }) == .open) return .{ .x = x, .y = 0 };
        }
        return error.NoOpenStartTile;
    }

    /// Part 1: leave the board and reappear on the opposite edge of the same
    /// row or column, keeping the facing.
    fn wrapFlat(self: Notes, from: Point, facing: Facing) Crossing {
        const delta = deltas[facing];
        var pos = from;
        while (true) {
            pos = .{
                .x = @mod(pos.x + delta[0], self.width),
                .y = @mod(pos.y + delta[1], self.height),
            };
            if (self.tile(pos) != .empty) return .{ .pos = pos, .facing = facing };
        }
    }
};

fn parsePath(allocator: std.mem.Allocator, text: []const u8) ![]Move {
    var path: std.ArrayList(Move) = .empty;
    errdefer path.deinit(allocator);
    var forward: ?u32 = null;
    for (text) |ch| switch (ch) {
        '0'...'9' => {
            const scaled = try std.math.mul(u32, forward orelse 0, 10);
            forward = try std.math.add(u32, scaled, ch - '0');
        },
        'L', 'R' => {
            if (forward) |count| {
                try path.append(allocator, .{ .forward = count });
                forward = null;
            }
            try path.append(allocator, .{ .turn = if (ch == 'L') .left else .right });
        },
        else => return error.InvalidPathCharacter,
    };
    if (forward) |count| try path.append(allocator, .{ .forward = count });
    if (path.items.len == 0) return error.EmptyPath;
    return path.toOwnedSlice(allocator);
}

/// One face of the folded cube: where it sits in the net plus the 3D directions
/// its grid axes point in. Every face keeps the invariant normal = down x right,
/// which is what makes the edge crossings exact.
const Face = struct { origin: Point, right: Vec3, down: Vec3, normal: Vec3 };

/// The net neighbour in the given facing rotates this basis a quarter turn
/// around the shared edge.
fn foldAcross(face: Face, facing: usize, origin: Point) Face {
    return switch (facing) {
        0 => .{ .origin = origin, .right = negate(face.normal), .down = face.down, .normal = face.right },
        1 => .{ .origin = origin, .right = face.right, .down = negate(face.normal), .normal = face.down },
        2 => .{ .origin = origin, .right = face.normal, .down = face.down, .normal = negate(face.right) },
        else => .{ .origin = origin, .right = face.right, .down = face.normal, .normal = negate(face.down) },
    };
}

fn indexOfCell(cells: [6]Point, cell: Point) ?usize {
    for (cells, 0..) |candidate, index| {
        if (candidate.x == cell.x and candidate.y == cell.y) return index;
    }
    return null;
}

const Cube = struct {
    size: i64,
    faces: [6]Face,
    by_normal: [6]usize,

    fn fold(notes: Notes) !Cube {
        var count: i64 = 0;
        for (notes.tiles) |t| {
            if (t != .empty) count += 1;
        }
        if (count == 0 or @rem(count, 6) != 0) return error.NotACubeNet;
        const area = @divExact(count, 6);
        var size: i64 = 1;
        while (size * size < area) size += 1;
        if (size * size != area) return error.NotACubeNet;

        var cells: [6]Point = undefined;
        var found: usize = 0;
        var row: i64 = 0;
        while (row * size < notes.height) : (row += 1) {
            var column: i64 = 0;
            while (column * size < notes.width) : (column += 1) {
                if (notes.tile(.{ .x = column * size, .y = row * size }) == .empty) continue;
                if (found == cells.len) return error.NotACubeNet;
                cells[found] = .{ .x = column, .y = row };
                found += 1;
            }
        }
        if (found != cells.len) return error.NotACubeNet;
        for (cells) |cell| {
            var j: i64 = 0;
            while (j < size) : (j += 1) {
                var i: i64 = 0;
                while (i < size) : (i += 1) {
                    const p = Point{ .x = cell.x * size + i, .y = cell.y * size + j };
                    if (notes.tile(p) == .empty) return error.NotACubeNet;
                }
            }
        }

        var faces: [6]Face = undefined;
        var assigned = [_]bool{false} ** 6;
        var queue: [6]usize = undefined;
        faces[0] = .{
            .origin = .{ .x = cells[0].x * size, .y = cells[0].y * size },
            .right = .{ 1, 0, 0 },
            .down = .{ 0, 1, 0 },
            .normal = .{ 0, 0, -1 },
        };
        assigned[0] = true;
        queue[0] = 0;
        var head: usize = 0;
        var tail: usize = 1;
        while (head < tail) : (head += 1) {
            const index = queue[head];
            const face = faces[index];
            for (deltas, 0..) |delta, facing| {
                const neighbor = Point{ .x = cells[index].x + delta[0], .y = cells[index].y + delta[1] };
                const next = indexOfCell(cells, neighbor) orelse continue;
                if (assigned[next]) continue;
                const origin = Point{ .x = cells[next].x * size, .y = cells[next].y * size };
                faces[next] = foldAcross(face, facing, origin);
                assigned[next] = true;
                queue[tail] = next;
                tail += 1;
            }
        }
        if (tail != cells.len) return error.NotACubeNet;

        var by_normal = [_]usize{std.math.maxInt(usize)} ** 6;
        for (faces, 0..) |face, index| {
            const slot = axisIndex(face.normal);
            if (by_normal[slot] != std.math.maxInt(usize)) return error.NotACubeNet;
            by_normal[slot] = index;
        }
        return .{ .size = size, .faces = faces, .by_normal = by_normal };
    }

    fn faceAt(self: Cube, p: Point) ?Face {
        for (self.faces) |face| {
            if (p.x >= face.origin.x and p.x < face.origin.x + self.size and
                p.y >= face.origin.y and p.y < face.origin.y + self.size) return face;
        }
        return null;
    }

    /// Part 2: step over a cube edge. Cell (i, j) of a face sits at the doubled
    /// lattice point size * normal + (2i - size + 1) * right + (2j - size + 1) *
    /// down. Crossing an edge swaps the components along the old normal and the
    /// travel direction, which reduces to adding travel - normal, and the walk
    /// continues along -normal on the face whose normal is the travel direction.
    fn wrap(self: Cube, from: Point, facing: Facing) !Crossing {
        const face = self.faceAt(from) orelse return error.OffTheCube;
        const travel = switch (facing) {
            0 => face.right,
            1 => face.down,
            2 => negate(face.right),
            3 => negate(face.down),
        };
        const i = from.x - face.origin.x;
        const j = from.y - face.origin.y;
        const here = add(scale(face.normal, self.size), add(
            scale(face.right, 2 * i - self.size + 1),
            scale(face.down, 2 * j - self.size + 1),
        ));
        const there = add(here, subtract(travel, face.normal));

        const dest = self.faces[self.by_normal[axisIndex(travel)]];
        const local = subtract(there, scale(dest.normal, self.size));
        const entry = negate(face.normal);
        const entry_facing: Facing = if (equal(entry, dest.right))
            0
        else if (equal(entry, dest.down))
            1
        else if (equal(entry, negate(dest.right)))
            2
        else
            3;
        return .{
            .pos = .{
                .x = dest.origin.x + @divFloor(dot(local, dest.right) + self.size - 1, 2),
                .y = dest.origin.y + @divFloor(dot(local, dest.down) + self.size - 1, 2),
            },
            .facing = entry_facing,
        };
    }
};

const Wrap = union(enum) {
    flat,
    cube: Cube,

    fn cross(self: Wrap, notes: Notes, from: Point, facing: Facing) !Crossing {
        switch (self) {
            .flat => return notes.wrapFlat(from, facing),
            .cube => |cube| return cube.wrap(from, facing),
        }
    }
};

fn walk(notes: Notes, wrapper: Wrap) !i64 {
    var pos = try notes.start();
    var facing: Facing = 0;
    for (notes.path) |move| switch (move) {
        .turn => |turn| facing = if (turn == .left) facing -% 1 else facing +% 1,
        .forward => |count| {
            var remaining = count;
            while (remaining > 0) : (remaining -= 1) {
                const delta = deltas[facing];
                const ahead = Point{ .x = pos.x + delta[0], .y = pos.y + delta[1] };
                const crossing = if (notes.tile(ahead) == .empty)
                    try wrapper.cross(notes, pos, facing)
                else
                    Crossing{ .pos = ahead, .facing = facing };
                if (notes.tile(crossing.pos) == .wall) break;
                pos = crossing.pos;
                facing = crossing.facing;
            }
        },
    };
    return 1000 * (pos.y + 1) + 4 * (pos.x + 1) + @as(i64, facing);
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    var notes = try Notes.parse(allocator, input);
    defer notes.deinit();
    const part1 = try walk(notes, .flat);
    const cube = try Cube.fold(notes);
    const part2 = try walk(notes, .{ .cube = cube });
    return .{ .part1 = part1, .part2 = part2 };
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

/// A 1-4-1 net with one tile per face: a layout nothing like the example's.
const unit_net = " .\n....\n  .\n\n2";
const unit_net_wall = " .\n....\n  #\n\n2";

test "part 1 example" {
    const answers = try solve(std.testing.allocator, @embedFile("examples/day_22.txt"));
    try std.testing.expectEqual(6032, answers.part1);
}

test "part 2 example" {
    const answers = try solve(std.testing.allocator, @embedFile("examples/day_22.txt"));
    try std.testing.expectEqual(5031, answers.part2);
}

test "folds a net with single tile faces" {
    // Flat wrapping never leaves the top row; on the cube the walk crosses onto
    // the right face and then onto the front face.
    const answers = try solve(std.testing.allocator, unit_net);
    try std.testing.expectEqual(1008, answers.part1);
    try std.testing.expectEqual(3013, answers.part2);
}

test "wall stops the walk after a crossing" {
    const answers = try solve(std.testing.allocator, unit_net_wall);
    try std.testing.expectEqual(1008, answers.part1);
    try std.testing.expectEqual(2013, answers.part2);
}

test "rejects malformed notes" {
    const allocator = std.testing.allocator;
    try std.testing.expectError(error.MissingPath, solve(allocator, "....\n...."));
    try std.testing.expectError(error.MissingMap, solve(allocator, ""));
    try std.testing.expectError(error.InvalidMapCharacter, solve(allocator, " .\n..x.\n  .\n\n2"));
    try std.testing.expectError(error.InvalidPathCharacter, solve(allocator, " .\n....\n  .\n\n2X"));
    try std.testing.expectError(error.NotACubeNet, solve(allocator, "....\n....\n\n1"));
    try std.testing.expectError(error.NotACubeNet, solve(allocator, "......\n......\n\n1"));
}
