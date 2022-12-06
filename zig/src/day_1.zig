const std = @import("std");
const testing = @import("std").testing;

pub fn max_calorie_elf(calorie_list: []const u64)  @as(u64, u64) {
    // Split list into Elves' inventories
    var elf_inventories = []const []const u64{};
    var current_elf_inventory = []const u64{};
    for (calorie_list) |calorie| {
        if (calorie == 0) {
            elf_inventories = std.mem.append(elf_inventories, current_elf_inventory);
            current_elf_inventory = []const u64{};
        } else {
            current_elf_inventory = std.mem.append(current_elf_inventory, calorie);
        }
    }

    // Find the Elf with the most Calories
    var max_calories = 0;
    var max_elf = u64(0);
    var i = u64(0);
    for (elf_inventories) |elf_inventory| {
        var total_calories = std.math.sum(elf_inventory);
        if (total_calories > max_calories) {
            max_calories = total_calories;
            max_elf = i;
        }
        i += 1;
    }

    return @as(max_elf, max_calories);
}


test "max_calorie_elf" {
    var result = max_calorie_elf([1000, 2000, 3000, 0, 4000, 0, 5000, 6000, 0, 7000, 8000, 9000, 0, 10000]);
    testing.expectEqual(result.0, u64(3));
    testing.expectEqual(result.1, u64(24000));
}

// pub fn main() void {
//     // Read input list of food Calories
//     var calorie_list = []const u64{};
//     var line: []const u8 = undefined;
//     while (line = std.io.readLine()) |line| {
//         var calorie: u64 = 0;
//         try std.str.parseInt(line, 10, &calorie);
//         calorie_list = std.mem.append(calorie_list, calorie);
//     }
//
//     var (max_elf, max_calories) = max_calorie_elf(calorie_list);
//
//     // Print the number of Calories carried by the Elf with the most Calories
//     std.debug.warn("Elf {} carries {} calories.\n", max_elf, max_calories);
// }
