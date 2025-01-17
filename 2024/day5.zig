const std = @import("std");
const testing = std.testing;

const Rule = struct {
    from: u32,
    to: u32,
};

const Update = struct {
    pages: std.ArrayList(u32),

    pub fn deinit(self: *Update) void {
        self.pages.deinit();
    }
};

const Graph = struct {
    rules: std.ArrayList(Rule),
    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator) Graph {
        return Graph{
            .rules = std.ArrayList(Rule).init(allocator),
            .allocator = allocator,
        };
    }

    pub fn deinit(self: *Graph) void {
        self.rules.deinit();
    }

    pub fn addRule(self: *Graph, from: u32, to: u32) !void {
        try self.rules.append(Rule{
            .from = from,
            .to = to,
        });
    }

    pub fn getRule(self: *Graph, index: usize) ?Rule {
        return self.rules.get(index);
    }

    pub fn isValidOrder(self: *Graph, update: *Update) bool {
        var positions = std.AutoHashMap(u32, usize).init(self.allocator);
        defer positions.deinit();

        for (update.pages.items, 0..) |page, i| {
            positions.put(page, i) catch return false;
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
                list.deinit();
            }
            adj_list.deinit();
        }

        for (update.pages.items) |page| {
            try adj_list.put(page, std.ArrayList(u32).init(self.allocator));
        }

        for (self.rules.items) |rule| {
            if (adj_list.getPtr(rule.from)) |list| {
                if (adj_list.contains(rule.to)) {
                    try list.append(rule.to);
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

        var sorted_pages = std.ArrayList(u32).init(self.allocator);
        defer sorted_pages.deinit();

        var queue = std.ArrayList(u32).init(self.allocator);
        defer queue.deinit();

        // Add all vertices with in-degree 0 to queue
        var deg_it = in_degree.iterator();
        while (deg_it.next()) |entry| {
            if (entry.value_ptr.* == 0) {
                try queue.append(entry.key_ptr.*);
            }
        }

        while (queue.items.len > 0) {
            const current = queue.pop();
            try sorted_pages.append(current);

            if (adj_list.getPtr(current)) |neighbors| {
                for (neighbors.items) |neighbor| {
                    if (in_degree.getPtr(neighbor)) |degree| {
                        degree.* -= 1;
                        if (degree.* == 0) {
                            try queue.append(neighbor);
                        }
                    }
                }
            }
        }

        update.pages.clearRetainingCapacity();
        try update.pages.appendSlice(sorted_pages.items);
    }
};

fn parseUpdate(allocator: std.mem.Allocator, line: []const u8) !Update {
    var update = Update{ .pages = std.ArrayList(u32).init(allocator) };
    var it = std.mem.split(u8, line, ",");

    while (it.next()) |num_str| {
        const num = try std.fmt.parseInt(u32, std.mem.trim(u8, num_str, " "), 10);
        try update.pages.append(num);
    }

    return update;
}

fn parseRule(line: []const u8) !Rule {
    var it = std.mem.split(u8, line, "|");
    const from_str = it.next() orelse return error.InvalidFormat;
    const to_str = it.next() orelse return error.InvalidFormat;

    return Rule{
        .from = try std.fmt.parseInt(u32, std.mem.trim(u8, from_str, " "), 10),
        .to = try std.fmt.parseInt(u32, std.mem.trim(u8, to_str, " "), 10),
    };
}

// Process input for Part 1
pub fn processInput(allocator: std.mem.Allocator, reader: anytype) !u32 {
    var graph = Graph.init(allocator);
    defer graph.deinit();

    var updates = std.ArrayList(Update).init(allocator);
    defer {
        for (updates.items) |*update| {
            update.deinit();
        }
        updates.deinit();
    }

    var buf: [1024]u8 = undefined;
    var parsing_rules = true;

    while (try reader.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        if (line.len == 0) {
            parsing_rules = false;
            continue;
        }

        if (parsing_rules) {
            const rule = try parseRule(line);
            try graph.addRule(rule.from, rule.to);
        } else {
            const update = try parseUpdate(allocator, line);
            try updates.append(update);
        }
    }

    var sum: u32 = 0;
    for (updates.items) |*update| {
        if (graph.isValidOrder(update)) {
            sum += getMiddleNumber(update);
        }
    }

    return sum;
}

// Process input for Part 2
pub fn processInputPart2(allocator: std.mem.Allocator, reader: anytype) !u32 {
    var graph = Graph.init(allocator);
    defer graph.deinit();

    var updates = std.ArrayList(Update).init(allocator);
    defer {
        for (updates.items) |*update| {
            update.deinit();
        }
        updates.deinit();
    }

    var buf: [1024]u8 = undefined;
    var parsing_rules = true;

    while (try reader.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        if (line.len == 0) {
            parsing_rules = false;
            continue;
        }

        if (parsing_rules) {
            const rule = try parseRule(line);
            try graph.addRule(rule.from, rule.to);
        } else {
            const update = try parseUpdate(allocator, line);
            try updates.append(update);
        }
    }

    var sum: u32 = 0;
    for (updates.items) |*update| {
        if (!graph.isValidOrder(update)) {
            try graph.sortUpdate(update);
            sum += getMiddleNumber(update);
        }
    }

    return sum;
}

fn getMiddleNumber(update: *const Update) u32 {
    if (update.pages.items.len == 0) return 0;
    return update.pages.items[update.pages.items.len / 2];
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const file = try std.fs.cwd().openFile("day5_input.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    const part1_result = try processInput(allocator, buf_reader.reader());

    try file.seekTo(0);
    buf_reader = std.io.bufferedReader(file.reader());
    const part2_result = try processInputPart2(allocator, buf_reader.reader());

    const stdout = std.io.getStdOut().writer();
    try stdout.print("Part 1: Sum of middle numbers from valid updates: {}\n", .{part1_result});
    try stdout.print("Part 2: Sum of middle numbers from corrected invalid updates: {}\n", .{part2_result});
}

test "example from problem description part 1" {
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

    var fbs = std.io.fixedBufferStream(example_input);
    const result = try processInput(testing.allocator, fbs.reader());
    try testing.expectEqual(@as(u32, 143), result);
}

test "example from problem description part 2" {
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

    var fbs = std.io.fixedBufferStream(example_input);
    const result = try processInputPart2(testing.allocator, fbs.reader());
    try testing.expectEqual(@as(u32, 123), result);
}
