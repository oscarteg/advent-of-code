const std = @import("std");
const os = std.os;
const warn = std.debug.warn;
const fs = std.fs;
const testing = std.testing;
const expect = testing.expect;

pub fn HashSet(comptime T: type) type {
    const mapType = if (T == []const u8) std.StringHashMap else std.HashMap(T, void);

    return struct {
        const Self = @This();

        map: mapType,

        pub fn init(allocator: *std.mem.Allocator) Self {
            return Self{
                .map = mapType.init(allocator),
            };
        }

        pub fn deinit(self: *Self) void {
            return self.map.deinit();
        }

        pub fn count(self: *const Self) usize {
            return self.map.count();
        }

        pub fn insertCheck(self: *Self, value: T) bool {
            return self.map.putNoClobber(value, undefined);
        }

        pub fn insert(self: *Self, value: T) void {
            return self.map.putAssumeCapacity(value, undefined);
        }
    };
}
