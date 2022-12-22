const std = @import("std");
const testing = @import("std").testing;

/// Reads an entire file into memory, caller owns the returned slice.
pub fn slurp(allocator: std.mem.Allocator, file_path: []const u8) ![]u8 {
    var path_buffer: [std.fs.MAX_PATH_BYTES]u8 = undefined;
    const path = try std.fs.realpath(file_path, &path_buffer);

    const file = try std.fs.openFileAbsolute(path, .{});
    defer file.close();

    var buf = try file.readToEndAlloc(
        allocator,
        (try file.stat()).size,
    );

    return allocator.resize(buf, buf.len - 1).?;
}

pub const Result = union(enum) {
    int: i64,
    string: []const u8,

    pub fn cmp(self: Result, result: Result) bool {
        switch (self) {
            .int => |i| return i == result.int,
            .string => |s| return std.mem.eql(u8, s, result.string),
        }
    }

    pub fn deinit(self: Result, allocator: std.mem.Allocator) void {
        switch (self) {
            .string => |s| return allocator.free(s),
            else => {},
        }
    }
};

pub fn puzzle_1(input: []const u8) !Result {
    var iter = std.mem.split(u8, input, "\n");
    var count: i32 = 0;
    var max: i32 = 0;

    while (iter.next()) |line| {
        if (line.len == 0) {
            if (count > max) max = count;
            count = 0;
            continue;
        }
        const number = try std.fmt.parseInt(i32, line, 0);
        count += number;
    }

    return .{ .int = max };
}

test "part_1_test" {
    var input =
        \\1000
        \\2000
        \\3000
        \\
        \\4000
        \\
        \\5000
        \\6000
        \\
        \\7000
        \\8000
        \\9000
        \\
        \\10000
    ;
    defer std.testing.allocator.free(input);

    const result = try puzzle_1(input);
    defer result.deinit(std.testing.allocator);

    try testing.expectEqual(result, .{ .int = 538 });
}
