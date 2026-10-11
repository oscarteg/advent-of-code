const std = @import("std");
const testing = std.testing;

const Rule = struct {
    from: u32,
    to: u32,
};

const Update = struct {
    pages: std.ArrayList(u32),

    pub fn deinit(self: *Update, allocator: std.mem.Allocator) void {
        self.pages.deinit(allocator);
    }
};

const Graph = struct {
    rules: std.ArrayList(Rule),
    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator) Graph {
        return Graph{
            .rules = .empty,
            .allocator = allocator,
        };
    }

    pub fn deinit(self: *Graph) void {
        self.rules.deinit(self.allocator);
    }

    pub fn addRule(self: *Graph, from: u32, to: u32) !void {
        try self.rules.append(self.allocator, Rule{
            .from = from,
            .to = to,
        });
    }

    pub fn isValidOrder(self: *Graph, update: *Update) !bool {
        var positions = std.AutoHashMap(u32, usize).init(self.allocator);
        defer positions.deinit();

        for (update.pages.items, 0..) |page, i| {
            try positions.put(page, i);
        }

        for (self.rules.items) |rule| {
            const from_pos = positions.get(rule.from);
            const to_pos = positions.get(rule.to);

            if (from_pos == null or to_pos == null) continue;
            if (from_pos.? > to_pos.?) return false;
        }

        return true;
    }

    pub fn sortUpdate(self: *Graph, update: *Update) !void {
        var adj_list = std.AutoHashMap(u32, std.ArrayList(u32)).init(self.allocator);
        defer {
            var it = adj_list.valueIterator();
            while (it.next()) |list| {
                list.deinit(self.allocator);
            }
            adj_list.deinit();
        }

        for (update.pages.items) |page| {
            try adj_list.put(page, .empty);
        }

        for (self.rules.items) |rule| {
            if (adj_list.getPtr(rule.from)) |list| {
                if (adj_list.contains(rule.to)) {
                    try list.append(self.allocator, rule.to);
                }
            }
        }

        var in_degree = std.AutoHashMap(u32, u32).init(self.allocator);
        defer in_degree.deinit();

        for (update.pages.items) |page| {
            try in_degree.put(page, 0);
        }

        var it = adj_list.iterator();
        while (it.next()) |entry| {
            for (entry.value_ptr.items) |to| {
                if (in_degree.getPtr(to)) |degree| {
                    degree.* += 1;
                }
            }
        }

        var sorted_pages: std.ArrayList(u32) = .empty;
        defer sorted_pages.deinit(self.allocator);

        var queue: std.ArrayList(u32) = .empty;
        defer queue.deinit(self.allocator);

        var deg_it = in_degree.iterator();
        while (deg_it.next()) |entry| {
            if (entry.value_ptr.* == 0) {
                try queue.append(self.allocator, entry.key_ptr.*);
            }
        }

        while (queue.items.len > 0) {
            const current = queue.items[queue.items.len - 1];
            queue.shrinkRetainingCapacity(queue.items.len - 1);
            try sorted_pages.append(self.allocator, current);

            if (adj_list.getPtr(current)) |neighbors| {
                for (neighbors.items) |neighbor| {
                    if (in_degree.getPtr(neighbor)) |degree| {
                        degree.* -= 1;
                        if (degree.* == 0) {
                            try queue.append(self.allocator, neighbor);
                        }
                    }
                }
            }
        }

        update.pages.clearRetainingCapacity();
        try update.pages.appendSlice(self.allocator, sorted_pages.items);
    }
};

fn parseUpdate(allocator: std.mem.Allocator, line: []const u8) !Update {
    var update = Update{ .pages = .empty };
    var it = std.mem.splitScalar(u8, line, ',');

    while (it.next()) |num_str| {
        const trimmed = std.mem.trim(u8, num_str, " ");
        const num = try std.fmt.parseInt(u32, trimmed, 10);
        try update.pages.append(allocator, num);
    }

    return update;
}

fn parseRule(line: []const u8) !Rule {
    var it = std.mem.splitScalar(u8, line, '|');
    const from_str = it.next() orelse return error.InvalidFormat;
    const to_str = it.next() orelse return error.InvalidFormat;

    return Rule{
        .from = try std.fmt.parseInt(u32, std.mem.trim(u8, from_str, " "), 10),
        .to = try std.fmt.parseInt(u32, std.mem.trim(u8, to_str, " "), 10),
    };
}

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !struct { u32, u32 } {
    var graph = Graph.init(allocator);
    defer graph.deinit();

    var updates: std.ArrayList(Update) = .empty;
    defer {
        for (updates.items) |*update| {
            update.deinit(allocator);
        }
        updates.deinit(allocator);
    }

    var parsing_rules = true;
    var it = std.mem.splitScalar(u8, input, '\n');

    while (it.next()) |line| {
        if (line.len == 0 or std.mem.eql(u8, line, "\r")) {
            parsing_rules = false;
            continue;
        }

        if (parsing_rules) {
            const rule = try parseRule(line);
            try graph.addRule(rule.from, rule.to);
        } else {
            const update = try parseUpdate(allocator, line);
            try updates.append(allocator, update);
        }
    }

    var sum1: u32 = 0;
    for (updates.items) |*update| {
        if (try graph.isValidOrder(update)) {
            sum1 += getMiddleNumber(update);
        }
    }

    var sum2: u32 = 0;
    for (updates.items) |*update| {
        if (!(try graph.isValidOrder(update))) {
            try graph.sortUpdate(update);
            sum2 += getMiddleNumber(update);
        }
    }

    return .{ sum1, sum2 };
}

fn getMiddleNumber(update: *const Update) u32 {
    if (update.pages.items.len == 0) return 0;
    return update.pages.items[update.pages.items.len / 2];
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

    try stdout.print("Part 1: {d}\nPart 2: {d}\n", .{ answers[0], answers[1] });
    try stdout.flush();
}

test "example part 1" {
    const example_input =
        \\47|53
        \\97|13
        \\97|61
        \\97|47
        \\75|29
        \\61|13
        \\75|53
        \\29|13
        \\97|29
        \\53|29
        \\61|53
        \\97|53
        \\61|29
        \\47|13
        \\75|47
        \\97|75
        \\47|61
        \\75|61
        \\47|29
        \\75|13
        \\53|13
        \\
        \\75,47,61,53,29
        \\97,61,53,29,13
        \\75,29,13
        \\75,97,47,61,53
        \\61,13,29
        \\97,13,75,29,47
    ;

    const result = try solve(testing.allocator, example_input);
    try testing.expectEqual(@as(u32, 143), result[0]);
}

test "example part 2" {
    const example_input =
        \\47|53
        \\97|13
        \\97|61
        \\97|47
        \\75|29
        \\61|13
        \\75|53
        \\29|13
        \\97|29
        \\53|29
        \\61|53
        \\97|53
        \\61|29
        \\47|13
        \\75|47
        \\97|75
        \\47|61
        \\75|61
        \\47|29
        \\75|13
        \\53|13
        \\
        \\75,97,47,61,53
        \\61,13,29
        \\97,13,75,29,47
    ;

    const result = try solve(testing.allocator, example_input);
    try testing.expectEqual(@as(u32, 123), result[1]);
}

test "single page" {
    const input = "1|2\n\n1\n";
    const result = try solve(testing.allocator, input);
    try testing.expectEqual(@as(u32, 1), result[0]);
}

test "invalid input" {
    const input = "invalid\n\n1,2\n";
    try testing.expectError(error.InvalidFormat, solve(testing.allocator, input));
}
