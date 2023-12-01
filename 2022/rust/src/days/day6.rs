/* fn run(s: &[char], window_size: usize) -> usize {
    let mut set = 0u32;
    for i in 0..s.len() {
        // Turn on bits as they enter the window
        set ^= 1 << (s[i] as u32 - 'a' as u32);

        // Turn off bits as they leave the window
        if i >= window_size {
            set ^= 1 << (s[i - window_size] as u32 - 'a' as u32);
        }

        // Check the current window and see if we're done
        if set.count_ones() as usize == window_size {
            return i + 1;
        }
    }
    panic!("No unique window found");
} */
