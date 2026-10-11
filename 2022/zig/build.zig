const std = @import("std");

const days = [_][]const u8{
    "day_1",  "day_2",  "day_3",  "day_4",  "day_5",
    "day_6",  "day_7",  "day_8",  "day_9",  "day_10",
    "day_11", "day_12", "day_13", "day_14", "day_15",
    "day_16", "day_17", "day_18", "day_19", "day_20",
    "day_21", "day_22", "day_23", "day_24", "day_25",
};

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const test_step = b.step("test", "Run unit tests for every day");

    for (days) |day| {
        const root_source_file = b.path(b.fmt("src/{s}.zig", .{day}));

        const exe = b.addExecutable(.{
            .name = day,
            .root_module = b.createModule(.{
                .root_source_file = root_source_file,
                .target = target,
                .optimize = optimize,
            }),
        });
        b.installArtifact(exe);

        const run_cmd = b.addRunArtifact(exe);
        run_cmd.step.dependOn(b.getInstallStep());
        if (b.args) |args| run_cmd.addArgs(args);
        const run_step = b.step(b.fmt("run-{s}", .{day}), b.fmt("Run {s}", .{day}));
        run_step.dependOn(&run_cmd.step);

        const exe_tests = b.addTest(.{
            .root_module = b.createModule(.{
                .root_source_file = root_source_file,
                .target = target,
                .optimize = optimize,
            }),
        });
        const run_tests = b.addRunArtifact(exe_tests);
        test_step.dependOn(&run_tests.step);
    }
}
