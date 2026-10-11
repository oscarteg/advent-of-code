const std = @import("std");

pub fn main(init: std.process.Init) !void {
    const arena = init.arena.allocator();
    var stdin_buffer: [4096]u8 = undefined;
    var stdin_reader = std.Io.File.stdin().reader(init.io, &stdin_buffer);
    const input = try stdin_reader.interface.allocRemaining(arena, .unlimited);

    const trimmed = std.mem.trim(u8, input, "\n\r ");
    const answers = try solve(arena, trimmed);

    var stdout_buffer: [256]u8 = undefined;
    var stdout_writer = std.Io.File.stdout().writer(init.io, &stdout_buffer);
    const stdout = &stdout_writer.interface;
    try stdout.print("Part 1: {d}\nPart 2: {d}\n", .{ answers.part1, answers.part2 });
    try stdout.flush();
}

pub fn solve(allocator: std.mem.Allocator, input: []const u8) !struct { part1: usize, part2: usize } {
    if (input.len == 0) return error.EmptyInput;

    const part1 = try findMarker(input, 4);
    const part2 = try findMarker(input, 14);

    _ = allocator;
    return .{ .part1 = part1, .part2 = part2 };
}

fn findMarker(signal: []const u8, window_size: usize) !usize {
    if (signal.len < window_size) return error.SignalTooShort;

    var seen: [256]u8 = [_]u8{0} ** 256;
    var unique_count: usize = 0;

    for (0..signal.len) |i| {
        const ch = signal[i];

        if (seen[ch] == 0) {
            unique_count += 1;
        }
        seen[ch] += 1;

        if (i >= window_size) {
            const old_ch = signal[i - window_size];
            seen[old_ch] -= 1;
            if (seen[old_ch] == 0) {
                unique_count -= 1;
            }
        }

        if (i >= window_size - 1 and unique_count == window_size) {
            return i + 1;
        }
    }

    return error.NoMarkerFound;
}

test "part1 examples" {
    const allocator = std.testing.allocator;
    try std.testing.expectEqual(try findMarker("mjqjpqmgbljsphdztnvjfqwrcgsmlb", 4), 7);
    try std.testing.expectEqual(try findMarker("bvwbjplbgvbhsrlpgdmjqwftvncz", 4), 5);
    try std.testing.expectEqual(try findMarker("nppdvjthqldpwncqszvftbrmjlhg", 4), 6);
    try std.testing.expectEqual(try findMarker("nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 4), 10);
    try std.testing.expectEqual(try findMarker("zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 4), 11);
    _ = allocator;
}

test "part2 examples" {
    const allocator = std.testing.allocator;
    try std.testing.expectEqual(try findMarker("mjqjpqmgbljsphdztnvjfqwrcgsmlb", 14), 19);
    try std.testing.expectEqual(try findMarker("bvwbjplbgvbhsrlpgdmjqwftvncz", 14), 23);
    try std.testing.expectEqual(try findMarker("nppdvjthqldpwncqszvftbrmjlhg", 14), 23);
    try std.testing.expectEqual(try findMarker("nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 14), 29);
    try std.testing.expectEqual(try findMarker("zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 14), 26);
    _ = allocator;
}

test "example from file" {
    const allocator = std.testing.allocator;
    const input = @embedFile("examples/day_6.txt");
    const trimmed = std.mem.trim(u8, input, "\n\r ");
    const answers = try solve(allocator, trimmed);
    try std.testing.expectEqual(answers.part1, 7);
    try std.testing.expectEqual(answers.part2, 19);
}

test "edge cases" {
    try std.testing.expectEqual(try findMarker("abcd", 4), 4);
    try std.testing.expectEqual(try findMarker("abcde", 4), 4);
}

test "malformed input" {
    try std.testing.expectError(error.SignalTooShort, findMarker("ab", 4));
    try std.testing.expectError(error.EmptyInput, solve(std.testing.allocator, ""));
}

test "no marker found" {
    try std.testing.expectError(error.NoMarkerFound, findMarker("aaaaaaaaaa", 4));
}
