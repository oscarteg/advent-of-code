const std = @import("std");

pub fn find_start_of_packet(buf: []const u8) u8 {
    var chars = [0, 0, 0, 0]; // The last four characters received

    // Iterate over the datastream buffer
    for (buf.len, i) in buf {
        // Add the current character to the array of last four characters
        chars[3] = chars[2];
        chars[2] = chars[1];
        chars[1] = chars[0];
        chars[0] = buf[i];

        // Check if the last four characters are all different
        if (chars[0] != chars[1]) && (chars[0] != chars[2]) && (chars[0] != chars[3]) &&
           (chars[1] != chars[2]) && (chars[1] != chars[3]) && (chars[2] != chars[3]) {
            return i + 1; // Return the number of characters processed
        }
    }

    return 0; // No start-of-packet marker was found
}

test "find_start_of_packet" {
    expect(find_start_of_packet("mjqjpqmgbljsphdztnvjfqwrcgsmlb".toSlice()) == 7);
    expect(find_start_of_packet("bvwbjplbgvbhsrlpgdmjqwftvncz".toSlice()) == 5);
    expect(find_start_of_packet("nppdvjthqldpwncqszvftbrmjlhg".toSlice()) == 6);
    expect(find_start_of_packet("nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg".toSlice()) == 10);
    expect(find_start_of_packet("zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw".toSlice()) == 11);
}


