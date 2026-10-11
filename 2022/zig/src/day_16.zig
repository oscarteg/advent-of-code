const std = @import("std");
const testing = std.testing;

const start_valve = "AA";
const minutes_alone: u32 = 30;
const minutes_with_elephant: u32 = 26;
/// Bounds the bitmask table below; real scans stay far under this.
const max_openable_valves = 20;
/// Travel time standing in for "no tunnel route at all".
const unreachable_minutes: u32 = std.math.maxInt(u32) / 4;

const Answers = struct {
    part1: u32,
    part2: u32,
};

const Error = error{
    MalformedValve,
    UnknownValve,
    DuplicateValve,
    MissingStartValve,
    NoValves,
    TooManyOpenableValves,
};

/// One line of the scan: a valve, its flow rate and the valve list its tunnels reach.
const ScannedValve = struct {
    name: []const u8,
    flow: u32,
    tunnels: []const u8,

    fn tunnelNames(self: ScannedValve) std.mem.SplitIterator(u8, .sequence) {
        return std.mem.splitSequence(u8, self.tunnels, ", ");
    }
};

/// Splits the valve list off its prefix, which is singular for a lone tunnel.
fn stripTunnelPrefix(text: []const u8) ?[]const u8 {
    const plural = "tunnels lead to valves ";
    const singular = "tunnel leads to valve ";
    if (std.mem.startsWith(u8, text, plural)) return text[plural.len..];
    if (std.mem.startsWith(u8, text, singular)) return text[singular.len..];
    return null;
}

fn parseValve(line: []const u8) Error!ScannedValve {
    if (!std.mem.startsWith(u8, line, "Valve ")) return Error.MalformedValve;
    const after_valve = line["Valve ".len..];

    const flow_marker = " has flow rate=";
    const name_end = std.mem.indexOf(u8, after_valve, flow_marker) orelse return Error.MalformedValve;
    const name = after_valve[0..name_end];
    if (name.len == 0) return Error.MalformedValve;

    const after_name = after_valve[name_end + flow_marker.len ..];
    const flow_end = std.mem.indexOf(u8, after_name, "; ") orelse return Error.MalformedValve;
    const flow = std.fmt.parseInt(u32, after_name[0..flow_end], 10) catch return Error.MalformedValve;

    const tunnels = stripTunnelPrefix(after_name[flow_end + 2 ..]) orelse return Error.MalformedValve;
    if (tunnels.len == 0) return Error.MalformedValve;
    return .{ .name = name, .flow = flow, .tunnels = tunnels };
}

/// The valves worth opening and the walking times the search needs; every index
/// is a bit position in the sets of opened valves.
const Network = struct {
    allocator: std.mem.Allocator,
    flows: []u32,
    travel: []u32,
    from_start: []u32,

    fn deinit(self: Network) void {
        self.allocator.free(self.flows);
        self.allocator.free(self.travel);
        self.allocator.free(self.from_start);
    }

    fn travelMinutes(self: Network, from: usize, to: usize) u32 {
        return self.travel[from * self.flows.len + to];
    }
};

/// Shortest walking time between every pair of valves, by Floyd-Warshall.
fn shortestWalkingTimes(
    allocator: std.mem.Allocator,
    valves: []const ScannedValve,
    index: std.StringHashMap(usize),
) ![]u32 {
    const count = valves.len;
    const minutes = try allocator.alloc(u32, count * count);
    errdefer allocator.free(minutes);
    @memset(minutes, unreachable_minutes);

    for (valves, 0..) |valve, from| {
        var tunnels = valve.tunnelNames();
        while (tunnels.next()) |name| {
            const to = index.get(name) orelse return Error.UnknownValve;
            minutes[from * count + to] = 1;
        }
    }
    for (0..count) |valve| minutes[valve * count + valve] = 0;

    for (0..count) |via| {
        for (0..count) |from| {
            const hop = minutes[from * count + via];
            if (hop >= unreachable_minutes) continue;
            for (0..count) |to| {
                const detour = hop + minutes[via * count + to];
                if (detour < minutes[from * count + to]) minutes[from * count + to] = detour;
            }
        }
    }
    return minutes;
}

fn buildNetwork(allocator: std.mem.Allocator, input: []const u8) !Network {
    var valves: std.ArrayList(ScannedValve) = .empty;
    defer valves.deinit(allocator);
    var index = std.StringHashMap(usize).init(allocator);
    defer index.deinit();

    var lines = std.mem.tokenizeScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;
        const valve = try parseValve(trimmed);
        if ((try index.fetchPut(valve.name, valves.items.len)) != null) return Error.DuplicateValve;
        try valves.append(allocator, valve);
    }
    if (valves.items.len == 0) return Error.NoValves;
    const start = index.get(start_valve) orelse return Error.MissingStartValve;

    const count = valves.items.len;
    const minutes = try shortestWalkingTimes(allocator, valves.items, index);
    defer allocator.free(minutes);

    var openable: std.ArrayList(usize) = .empty;
    defer openable.deinit(allocator);
    for (valves.items, 0..) |valve, position| {
        if (valve.flow > 0) try openable.append(allocator, position);
    }
    if (openable.items.len > max_openable_valves) return Error.TooManyOpenableValves;

    const openable_count = openable.items.len;
    const flows = try allocator.alloc(u32, openable_count);
    errdefer allocator.free(flows);
    const travel = try allocator.alloc(u32, openable_count * openable_count);
    errdefer allocator.free(travel);
    const from_start = try allocator.alloc(u32, openable_count);
    errdefer allocator.free(from_start);

    for (openable.items, 0..) |position, from| {
        flows[from] = valves.items[position].flow;
        from_start[from] = minutes[start * count + position];
        for (openable.items, 0..) |other, to| {
            travel[from * openable_count + to] = minutes[position * count + other];
        }
    }
    return .{ .allocator = allocator, .flows = flows, .travel = travel, .from_start = from_start };
}

/// Depth-first walk over the valves still closed, recording for every set of
/// opened valves the most pressure that set can release.
const Search = struct {
    network: Network,
    best: []u32,

    fn explore(self: Search, valve: usize, minutes_left: u32, opened: usize, released: u32) void {
        self.best[opened] = @max(self.best[opened], released);
        for (self.network.flows, 0..) |flow, next| {
            const bit = @as(usize, 1) << @intCast(next);
            if (opened & bit != 0) continue;
            const cost = self.network.travelMinutes(valve, next) + 1;
            if (cost >= minutes_left) continue;
            const left = minutes_left - cost;
            self.explore(next, left, opened | bit, released + left * flow);
        }
    }
};

/// Most pressure released per set of opened valves within the minutes given.
fn bestPerOpenedSet(allocator: std.mem.Allocator, network: Network, minutes: u32) ![]u32 {
    const best = try allocator.alloc(u32, @as(usize, 1) << @intCast(network.flows.len));
    @memset(best, 0);

    const search = Search{ .network = network, .best = best };
    for (network.flows, 0..) |flow, valve| {
        const cost = network.from_start[valve] + 1;
        if (cost >= minutes) continue;
        const left = minutes - cost;
        search.explore(valve, left, @as(usize, 1) << @intCast(valve), left * flow);
    }
    return best;
}

/// Turns best[set] into the best over every subset of set, so two disjoint
/// routes can be combined by looking up a set and its complement.
fn maximizeOverSubsets(best: []u32, valve_count: usize) void {
    for (0..valve_count) |position| {
        const bit = @as(usize, 1) << @intCast(position);
        for (0..best.len) |set| {
            if (set & bit == 0) continue;
            best[set] = @max(best[set], best[set ^ bit]);
        }
    }
}

fn solve(allocator: std.mem.Allocator, input: []const u8) !Answers {
    const network = try buildNetwork(allocator, input);
    defer network.deinit();

    const alone = try bestPerOpenedSet(allocator, network, minutes_alone);
    defer allocator.free(alone);
    var part1: u32 = 0;
    for (alone) |released| part1 = @max(part1, released);

    const shared = try bestPerOpenedSet(allocator, network, minutes_with_elephant);
    defer allocator.free(shared);
    maximizeOverSubsets(shared, network.flows.len);
    const every_valve = shared.len - 1;
    var part2: u32 = 0;
    for (shared, 0..) |released, mine| part2 = @max(part2, released + shared[every_valve ^ mine]);

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

const example = @embedFile("examples/day_16.txt");

test "official example" {
    const answers = try solve(testing.allocator, example);
    try testing.expectEqual(@as(u32, 1651), answers.part1);
    try testing.expectEqual(@as(u32, 1707), answers.part2);
}

test "lone valve is opened as early as possible" {
    // A minute to walk to BB and a minute to open it leaves 28 of the 30
    // minutes at flow 10, and 24 of the 26 when sharing the work.
    const input =
        \\Valve AA has flow rate=0; tunnel leads to valve BB
        \\Valve BB has flow rate=10; tunnel leads to valve AA
    ;
    const answers = try solve(testing.allocator, input);
    try testing.expectEqual(@as(u32, 280), answers.part1);
    try testing.expectEqual(@as(u32, 240), answers.part2);
}

test "valves without flow release nothing" {
    const input =
        \\Valve AA has flow rate=0; tunnels lead to valves BB, CC
        \\Valve BB has flow rate=0; tunnel leads to valve AA
        \\Valve CC has flow rate=0; tunnel leads to valve AA
    ;
    const answers = try solve(testing.allocator, input);
    try testing.expectEqual(@as(u32, 0), answers.part1);
    try testing.expectEqual(@as(u32, 0), answers.part2);
}

test "valves without a tunnel route stay closed" {
    const input =
        \\Valve AA has flow rate=0; tunnel leads to valve BB
        \\Valve BB has flow rate=0; tunnel leads to valve AA
        \\Valve CC has flow rate=5; tunnel leads to valve DD
        \\Valve DD has flow rate=0; tunnel leads to valve CC
    ;
    const answers = try solve(testing.allocator, input);
    try testing.expectEqual(@as(u32, 0), answers.part1);
    try testing.expectEqual(@as(u32, 0), answers.part2);
}

test "malformed scans are rejected" {
    try testing.expectError(Error.NoValves, solve(testing.allocator, "\n\n"));
    try testing.expectError(Error.MalformedValve, solve(testing.allocator, "Valve AA has flow rate=0"));
    try testing.expectError(Error.MalformedValve, solve(testing.allocator, "Valve AA has flow rate=x; tunnel leads to valve BB"));
    try testing.expectError(Error.MalformedValve, solve(testing.allocator, "Valve AA has flow rate=0; leads to BB"));
    try testing.expectError(Error.MalformedValve, solve(testing.allocator, "AA has flow rate=0; tunnel leads to valve AA"));
}

test "inconsistent scans are rejected" {
    try testing.expectError(Error.UnknownValve, solve(testing.allocator, "Valve AA has flow rate=0; tunnel leads to valve ZZ"));
    try testing.expectError(Error.MissingStartValve, solve(testing.allocator, "Valve BB has flow rate=1; tunnel leads to valve BB"));
    const duplicate =
        \\Valve AA has flow rate=0; tunnel leads to valve AA
        \\Valve AA has flow rate=1; tunnel leads to valve AA
    ;
    try testing.expectError(Error.DuplicateValve, solve(testing.allocator, duplicate));
}
