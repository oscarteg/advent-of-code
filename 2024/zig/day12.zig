const std = @import("std");

const Answers = struct {
    part1: i64,
    part2: i64,
};

const Plot = struct {
    row: isize,
    col: isize,
};

/// What one region costs to fence: part 1 prices area * perimeter, part 2
/// area * sides.
const Region = struct {
    area: i64 = 0,
    perimeter: i64 = 0,
    sides: i64 = 0,
};

const neighbour_steps = [4]Plot{
    .{ .row = -1, .col = 0 },
    .{ .row = 1, .col = 0 },
    .{ .row = 0, .col = -1 },
    .{ .row = 0, .col = 1 },
};
const diagonal_steps = [2]isize{ -1, 1 };

/// A fixed-width grid of plant letters, stored row-major.
const Garden = struct {
    plants: []const u8,
    width: usize,
    height: usize,

    fn indexOf(self: Garden, at: Plot) usize {
        return @as(usize, @intCast(at.row)) * self.width + @as(usize, @intCast(at.col));
    }

    fn hasPlant(self: Garden, at: Plot, plant: u8) bool {
        if (at.row < 0 or at.col < 0) return false;
        const row: usize = @intCast(at.row);
        const col: usize = @intCast(at.col);
        if (row >= self.height or col >= self.width) return false;
        return self.plants[row * self.width + col] == plant;
    }

    /// Counts the corners the plot contributes to its region: a corner is
    /// convex when both plots adjacent to a diagonal lie outside the region,
    /// and concave when both lie inside while the diagonal plot does not.
    fn cornersAt(self: Garden, at: Plot, plant: u8) i64 {
        var corners: i64 = 0;
        for (diagonal_steps) |row_step| {
            for (diagonal_steps) |col_step| {
                const vertical = self.hasPlant(.{ .row = at.row + row_step, .col = at.col }, plant);
                const horizontal = self.hasPlant(.{ .row = at.row, .col = at.col + col_step }, plant);
                const diagonal = self.hasPlant(.{ .row = at.row + row_step, .col = at.col + col_step }, plant);
                if (!vertical and !horizontal) corners += 1;
                if (vertical and horizontal and !diagonal) corners += 1;
            }
        }
        return corners;
    }

    /// Flood-fills the region containing start, marking every plot it reaches
    /// visited. A region has as many sides as corners.
    fn measureRegion(self: Garden, allocator: std.mem.Allocator, start: Plot, visited: []bool) !Region {
        const plant = self.plants[self.indexOf(start)];
        visited[self.indexOf(start)] = true;

        var queue: std.ArrayList(Plot) = .empty;
        defer queue.deinit(allocator);
        try queue.append(allocator, start);

        var measured = Region{};
        var head: usize = 0;
        while (head < queue.items.len) {
            const current = queue.items[head];
            head += 1;
            measured.area += 1;
            measured.sides += self.cornersAt(current, plant);
            for (neighbour_steps) |step| {
                const next = Plot{ .row = current.row + step.row, .col = current.col + step.col };
                if (!self.hasPlant(next, plant)) {
                    measured.perimeter += 1;
                    continue;
                }
                if (visited[self.indexOf(next)]) continue;
                visited[self.indexOf(next)] = true;
                try queue.append(allocator, next);
            }
        }
        return measured;
    }
};

fn parse(allocator: std.mem.Allocator, input: []const u8) !Garden {
    var plants: std.ArrayList(u8) = .empty;
    errdefer plants.deinit(allocator);

    var width: usize = 0;
    var height: usize = 0;
    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw| {
        const line = std.mem.trim(u8, raw, " \r");
        if (line.len == 0) continue;
        if (height == 0) {
            width = line.len;
        } else if (line.len != width) {
            return error.RaggedGarden;
        }
        for (line) |plant| {
            if (!std.ascii.isAlphabetic(plant)) return error.InvalidPlant;
        }
        try plants.appendSlice(allocator, line);
        height += 1;
    }
    if (height == 0) return error.EmptyGarden;

    return Garden{
        .plants = try plants.toOwnedSlice(allocator),
        .width = width,
        .height = height,
    };
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const garden = try parse(allocator, input);
    defer allocator.free(garden.plants);

    const visited = try allocator.alloc(bool, garden.plants.len);
    defer allocator.free(visited);
    @memset(visited, false);

    var answers = Answers{ .part1 = 0, .part2 = 0 };
    for (0..garden.height) |row| {
        for (0..garden.width) |col| {
            const at = Plot{ .row = @intCast(row), .col = @intCast(col) };
            if (visited[garden.indexOf(at)]) continue;
            const measured = try garden.measureRegion(allocator, at, visited);
            answers.part1 += measured.area * measured.perimeter;
            answers.part2 += measured.area * measured.sides;
        }
    }
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
    const answers = try solve(std.testing.allocator, @embedFile("examples/day12.txt"));
    try std.testing.expectEqual(@as(i64, 1930), answers.part1);
    try std.testing.expectEqual(@as(i64, 1206), answers.part2);
}

test "smaller examples" {
    const five_regions = try solve(std.testing.allocator,
        \\AAAA
        \\BBCD
        \\BBCC
        \\EEEC
    );
    try std.testing.expectEqual(@as(i64, 140), five_regions.part1);
    try std.testing.expectEqual(@as(i64, 80), five_regions.part2);

    const nested_regions = try solve(std.testing.allocator,
        \\OOOOO
        \\OXOXO
        \\OOOOO
        \\OXOXO
        \\OOOOO
    );
    try std.testing.expectEqual(@as(i64, 772), nested_regions.part1);
    try std.testing.expectEqual(@as(i64, 436), nested_regions.part2);
}

test "side counting examples" {
    const e_shape = try solve(std.testing.allocator,
        \\EEEEE
        \\EXXXX
        \\EEEEE
        \\EXXXX
        \\EEEEE
    );
    try std.testing.expectEqual(@as(i64, 236), e_shape.part2);

    const diagonally_touching = try solve(std.testing.allocator,
        \\AAAAAA
        \\AAABBA
        \\AAABBA
        \\ABBAAA
        \\ABBAAA
        \\AAAAAA
    );
    try std.testing.expectEqual(@as(i64, 368), diagonally_touching.part2);
}

test "single plot region" {
    const answers = try solve(std.testing.allocator, "A\n");
    try std.testing.expectEqual(@as(i64, 4), answers.part1);
    try std.testing.expectEqual(@as(i64, 4), answers.part2);
}

test "malformed input returns an error" {
    try std.testing.expectError(error.EmptyGarden, solve(std.testing.allocator, ""));
    try std.testing.expectError(error.RaggedGarden, solve(std.testing.allocator, "AAAA\nBB\n"));
    try std.testing.expectError(error.InvalidPlant, solve(std.testing.allocator, "AA\nA1\n"));
}
