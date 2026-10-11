const std = @import("std");
const testing = std.testing;

const ore = 0;
const clay = 1;
const obsidian = 2;
const geode = 3;
const resource_count = 3;

const part1_minutes = 24;
const part2_minutes = 32;
const part2_blueprints = 3;

const build_order = [_]usize{ geode, obsidian, clay, ore };

const Answers = struct {
    part1: u64,
    part2: u64,
};

fn skipSpaces(text: []const u8) []const u8 {
    var index: usize = 0;
    while (index < text.len and (text[index] == ' ' or text[index] == '\t' or text[index] == '\r')) index += 1;
    return text[index..];
}

/// Reads a blueprint line literal by literal so malformed input is rejected
/// instead of silently misparsed.
const Scanner = struct {
    rest: []const u8,

    fn literal(self: *Scanner, expected: []const u8) !void {
        const rest = skipSpaces(self.rest);
        if (!std.mem.startsWith(u8, rest, expected)) return error.MalformedBlueprint;
        self.rest = rest[expected.len..];
    }

    fn number(self: *Scanner) !u32 {
        const rest = skipSpaces(self.rest);
        var digits: usize = 0;
        while (digits < rest.len and std.ascii.isDigit(rest[digits])) digits += 1;
        if (digits == 0) return error.MalformedBlueprint;
        self.rest = rest[digits..];
        return std.fmt.parseInt(u32, rest[0..digits], 10);
    }

    fn finish(self: Scanner) !void {
        if (std.mem.trim(u8, self.rest, " \t\r").len != 0) return error.MalformedBlueprint;
    }
};

/// A moment in a blueprint's schedule. Geode robots are not tracked: the
/// geodes they will have cracked once the clock runs out are added to
/// `geodes` the minute they are finished.
const State = struct {
    minutes_left: u32,
    resources: [resource_count]u32,
    robots: [resource_count]u32,
    geodes: u32,

    /// Waits until `robot` is affordable and builds it, or returns null when
    /// it can never pay off before the clock runs out.
    fn build(self: State, robot: usize, cost: [resource_count]u32) ?State {
        var wait: u32 = 0;
        for (0..resource_count) |resource| {
            const missing = cost[resource] -| self.resources[resource];
            if (missing == 0) continue;
            const rate = self.robots[resource];
            if (rate == 0) return null;
            wait = @max(wait, (missing + rate - 1) / rate);
        }

        const elapsed = wait + 1;
        if (elapsed >= self.minutes_left) return null;

        var next = self;
        next.minutes_left -= elapsed;
        for (0..resource_count) |resource| {
            next.resources[resource] = self.resources[resource] + self.robots[resource] * elapsed - cost[resource];
        }
        if (robot == geode) {
            next.geodes += next.minutes_left;
        } else {
            next.robots[robot] += 1;
        }
        return next;
    }
};

/// Robot costs in ore, clay and obsidian, indexed by robot type.
const Blueprint = struct {
    id: u32,
    costs: [build_order.len][resource_count]u32,
    /// Highest useful number of ore, clay and obsidian robots: no recipe can
    /// spend more of a resource per minute than its largest cost.
    caps: [resource_count]u32,

    fn parse(line: []const u8) !Blueprint {
        var scanner = Scanner{ .rest = line };
        try scanner.literal("Blueprint");
        const id = try scanner.number();
        try scanner.literal(":");
        try scanner.literal("Each ore robot costs");
        const ore_robot = try scanner.number();
        try scanner.literal("ore.");
        try scanner.literal("Each clay robot costs");
        const clay_robot = try scanner.number();
        try scanner.literal("ore.");
        try scanner.literal("Each obsidian robot costs");
        const obsidian_robot_ore = try scanner.number();
        try scanner.literal("ore and");
        const obsidian_robot_clay = try scanner.number();
        try scanner.literal("clay.");
        try scanner.literal("Each geode robot costs");
        const geode_robot_ore = try scanner.number();
        try scanner.literal("ore and");
        const geode_robot_obsidian = try scanner.number();
        try scanner.literal("obsidian.");
        try scanner.finish();

        return .{
            .id = id,
            .costs = .{
                .{ ore_robot, 0, 0 },
                .{ clay_robot, 0, 0 },
                .{ obsidian_robot_ore, obsidian_robot_clay, 0 },
                .{ geode_robot_ore, 0, geode_robot_obsidian },
            },
            .caps = .{
                @max(ore_robot, clay_robot, obsidian_robot_ore, geode_robot_ore),
                obsidian_robot_clay,
                geode_robot_obsidian,
            },
        };
    }

    fn maxGeodes(self: Blueprint, minutes: u32) u32 {
        var best: u32 = 0;
        self.search(.{
            .minutes_left = minutes,
            .resources = .{ 0, 0, 0 },
            .robots = .{ 1, 0, 0 },
            .geodes = 0,
        }, &best);
        return best;
    }

    /// Depth-first search over the next robot to build, geode robots first,
    /// dropping every branch whose optimistic bound cannot beat the best
    /// schedule found so far.
    fn search(self: Blueprint, state: State, best: *u32) void {
        if (state.geodes > best.*) best.* = state.geodes;
        if (self.optimisticGeodes(state) <= best.*) return;
        for (build_order) |robot| {
            if (robot != geode and state.robots[robot] >= self.caps[robot]) continue;
            if (state.build(robot, self.costs[robot])) |next| self.search(next, best);
        }
    }

    /// Upper bound on the geodes reachable from `state`, from a relaxed game
    /// that ignores ore costs and lets one robot of every type be built per
    /// minute. Building everything affordable as early as possible is optimal
    /// there (clay is only spent on obsidian robots, obsidian only on geode
    /// robots), so no real schedule can beat the result.
    fn optimisticGeodes(self: Blueprint, state: State) u32 {
        const obsidian_price = self.costs[geode][obsidian];
        const clay_price = self.costs[obsidian][clay];
        var clay_amount = state.resources[clay];
        var obsidian_amount = state.resources[obsidian];
        var clay_robots = state.robots[clay];
        var obsidian_robots = state.robots[obsidian];
        var geodes = state.geodes;

        var remaining = state.minutes_left;
        while (remaining > 0) : (remaining -= 1) {
            var built_obsidian_robot: u32 = 0;
            if (obsidian_amount >= obsidian_price) {
                obsidian_amount -= obsidian_price;
                geodes += remaining - 1;
            }
            if (clay_amount >= clay_price) {
                clay_amount -= clay_price;
                built_obsidian_robot = 1;
            }
            clay_amount += clay_robots;
            obsidian_amount += obsidian_robots;
            clay_robots += 1;
            obsidian_robots += built_obsidian_robot;
        }
        return geodes;
    }
};

fn parseBlueprints(allocator: std.mem.Allocator, input: []const u8) ![]Blueprint {
    var blueprints: std.ArrayList(Blueprint) = .empty;
    errdefer blueprints.deinit(allocator);

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |raw_line| {
        const line = std.mem.trim(u8, raw_line, " \t\r");
        if (line.len == 0) continue;
        try blueprints.append(allocator, try Blueprint.parse(line));
    }

    if (blueprints.items.len == 0) return error.NoBlueprints;
    return blueprints.toOwnedSlice(allocator);
}

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const blueprints = try parseBlueprints(allocator, input);
    defer allocator.free(blueprints);

    var part1: u64 = 0;
    var part2: u64 = 1;
    for (blueprints, 0..) |blueprint, index| {
        part1 += @as(u64, blueprint.id) * @as(u64, blueprint.maxGeodes(part1_minutes));
        if (index < part2_blueprints) part2 *= @as(u64, blueprint.maxGeodes(part2_minutes));
    }
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

const second_blueprint = "Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian.";

test "official example" {
    const answers = try solve(testing.allocator, @embedFile("examples/day_19.txt"));
    try testing.expectEqual(@as(u64, 33), answers.part1);
    try testing.expectEqual(@as(u64, 3472), answers.part2);
}

test "quality level uses the blueprint id" {
    // The example's second blueprint cracks 12 geodes in 24 minutes and 62 in
    // 32 minutes, so on its own it scores 2 * 12 and 62.
    const answers = try solve(testing.allocator, second_blueprint);
    try testing.expectEqual(@as(u64, 24), answers.part1);
    try testing.expectEqual(@as(u64, 62), answers.part2);
}

test "unaffordable blueprint cracks nothing" {
    const input = "Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 50 clay. Each geode robot costs 90 ore and 90 obsidian.";
    const answers = try solve(testing.allocator, input);
    try testing.expectEqual(@as(u64, 0), answers.part1);
    try testing.expectEqual(@as(u64, 0), answers.part2);
}

test "blank input returns error" {
    try testing.expectError(error.NoBlueprints, solve(testing.allocator, "\n  \n"));
}

test "truncated blueprint returns error" {
    try testing.expectError(error.MalformedBlueprint, solve(testing.allocator, "Blueprint 1: Each ore robot costs 4 ore."));
}

test "non numeric id returns error" {
    try testing.expectError(error.MalformedBlueprint, solve(testing.allocator, "Blueprint two: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian."));
}

test "trailing text returns error" {
    try testing.expectError(error.MalformedBlueprint, solve(testing.allocator, second_blueprint ++ " and a pony."));
}
