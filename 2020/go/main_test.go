package main

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunner(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "aoc2020")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	build := exec.CommandContext(ctx, "go", "build", "-race", "-o", binary, ".")
	output, err := build.CombinedOutput()
	require.NoError(t, err, string(output))

	expenses := "1721\n979\n366\n299\n675\n1456\n"
	passwords := "1-3 a: abcde\n1-3 b: cdefg\n2-9 c: ccccccccc\n"
	trees := "..##.......\n#...#...#..\n.#....#..#.\n..#.#...#.#\n.#...##..#.\n" +
		"..#.##.....\n.#.#.#....#\n.#........#\n#.##...#...\n#...##....#\n.#..#...#.#\n"
	passports := "ecl:gry pid:860033327 eyr:2020 hcl:#fffffd\nbyr:1937 iyr:2017 cid:147 hgt:183cm\n\n" +
		"iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884\nhcl:#cfa07d byr:1929\n\n" +
		"hcl:#ae17e1 iyr:2013\neyr:2024\necl:brn pid:760753108 byr:1931\nhgt:179cm\n\n" +
		"hcl:#cfa07d eyr:2025 pid:166559648\niyr:2011 ecl:brn hgt:59in\n"
	seatsWithGap := "FBFBBFFRLR\nFBFBBFFRRL\nFBFBBFBLLL\nFBFBBFBLLR\n"
	groups := "abc\n\na\nb\nc\n\nab\nac\n\na\na\na\na\n\nb\n"
	bags := "light red bags contain 1 bright white bag, 2 muted yellow bags.\n" +
		"dark orange bags contain 3 bright white bags, 4 muted yellow bags.\n" +
		"bright white bags contain 1 shiny gold bag.\n" +
		"muted yellow bags contain 2 shiny gold bags, 9 faded blue bags.\n" +
		"shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.\n" +
		"dark olive bags contain 3 faded blue bags, 4 dotted black bags.\n" +
		"vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.\n" +
		"faded blue bags contain no other bags.\ndotted black bags contain no other bags.\n"
	bootCode := "nop +0\nacc +1\njmp +4\nacc +3\njmp -3\nacc -99\nacc +1\njmp -4\nacc +6\n"
	jumps := "35\n20\n15\n25\n47\n40\n62\n55\n65\n95\n102\n117\n150\n182\n127\n219\n299\n277\n309\n576\n"
	adapters := "28\n33\n18\n42\n31\n14\n46\n20\n48\n47\n24\n23\n49\n45\n19\n38\n39\n11\n1\n32\n25\n35\n8\n17\n7\n9\n4\n2\n34\n10\n3\n"
	seatLayout := "L.LL.LL.LL\nLLLLLLL.LL\nL.L.L..L..\nLLLL.LL.LL\nL.LL.LL.LL\n" +
		"L.LLLLL.LL\n..L.L.....\nLLLLLLLLLL\nL.LLLLLL.L\nL.LLLLL.LL\n"
	navigation := "F10\nN3\nF7\nR90\nF11\n"
	buses := "939\n7,13,x,x,59,x,31,19\n"
	maskB := "mask = 000000000000000000000000000000X1001X\nmem[42] = 100\n" +
		"mask = 00000000000000000000000000000000X0XX\nmem[26] = 1\n"
	memoryGame := "0,3,6\n"
	ticketNotesB := "departure gate: 1-3\ndeparture track: 1-6\nseat: 1-9\n\nyour ticket:\n2,5,7\n\n" +
		"nearby tickets:\n1,4,9\n3,6,8\n"
	seaMonsters := ".#.\n..#\n###\n"
	bitmask := "1 + 2 * 3 + 4 * 5 + 6\n1 + (2 * 3) + (4 * (5 + 6))\n2 * 3 + (4 * 5)\n" +
		"5 + (8 * 3 + 9 + 3 * 4 * 3)\n5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))\n" +
		"((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2\n"
	messagesLoop := "42: 9 14 | 10 1\n9: 14 27 | 1 26\n10: 23 14 | 28 1\n1: \"a\"\n11: 42 31\n" +
		"5: 1 14 | 15 1\n19: 14 1 | 14 14\n12: 24 14 | 19 1\n16: 15 1 | 14 14\n31: 14 17 | 1 13\n" +
		"6: 14 14 | 1 14\n2: 1 24 | 14 4\n0: 8 11\n13: 14 3 | 1 12\n15: 1 | 14\n17: 14 2 | 1 7\n" +
		"23: 25 1 | 22 14\n28: 16 1\n4: 1 1\n20: 14 14 | 1 15\n3: 5 14 | 16 1\n27: 1 6 | 14 18\n" +
		"14: \"b\"\n21: 14 1 | 1 14\n25: 1 1 | 1 14\n22: 14 14\n8: 42\n26: 14 22 | 1 20\n18: 15 15\n" +
		"7: 14 5 | 1 21\n24: 14 1\n\nabbbbbabbbaaaababbaabbbbabababbbabbbbbbabaaaa\nbbabbbbaabaabba\n" +
		"babbbbaabbbbbabbbbbbaabaaabaaa\naaabbbbbbaaaabaababaabababbabaaabbababababaaa\n" +
		"bbbbbbbaaaabbbbaaabbabaaa\nbbbababbbbaaaaaaaabbababaaababaabab\nababaaaaaabaaab\n" +
		"ababaaaaabbbaba\nbaabbaaaabbaaaababbaababb\nabbbbabbbbaaaababbbbbbaaaababb\n" +
		"aaaaabbaabaaaaababaa\naaaabbaaaabbaaa\naaaabbaabbaaaaaaabbbabbbaaabbaabaab\n" +
		"babaaabbbaaabaababbaabababaaab\naabbbbbaabbbaaaaaabbbbbababaaaaabbaaabba\n"
	jigsawTiles := "Tile 2311:\n..##.#..#.\n##..#.....\n#...##..#.\n####.#...#\n##.##.###.\n##...#.###\n" +
		".#.#.#..##\n..#....#..\n###...#.#.\n..###..###\n\nTile 1951:\n#.##...##.\n#.####...#\n.....#..##\n" +
		"#...######\n.##.#....#\n.###.#####\n###.##.##.\n.###....#.\n..#.#..#.#\n#...##.#..\n\nTile 1171:\n" +
		"####...##.\n#..##.#..#\n##.#..#.#.\n.###.####.\n..###.####\n.##....##.\n.#...####.\n#.##.####.\n" +
		"####..#...\n.....##...\n\nTile 1427:\n###.##.#..\n.#..#.##..\n.#.##.#..#\n#.#.#.##.#\n" +
		"....#...##\n...##..##.\n...#.#####\n.#.####.#.\n..#..###.#\n..##.#..#.\n\nTile 1489:\n" +
		"##.#.#....\n..##...#..\n.##..##...\n..#...#...\n#####...#.\n#..#.#.#.#\n...#.#.#..\n##.#...##.\n" +
		"..##.##.##\n###.##.#..\n\nTile 2473:\n#....####.\n#..#.##...\n#.##..#...\n######.#.#\n.#...#.#.#\n" +
		".#########\n.###.#..#.\n########.#\n##...##.#.\n..###.#.#.\n\nTile 2971:\n..#.#....#\n#...###...\n" +
		"#.#.###...\n##.##..#..\n.#####..##\n.#..####.#\n#..#.#..#.\n..####.###\n..#.#.###.\n...#.#.#.#\n\n" +
		"Tile 2729:\n...#.#.#.#\n####.#....\n..#.#.....\n....#..#.#\n.##..##.#.\n.#.####...\n####.#.#..\n" +
		"##.####...\n##..#.##..\n#.##...##.\n\nTile 3079:\n#.#.#####.\n.#..######\n..#.......\n######....\n" +
		"####.#..#.\n.#...#.##.\n#.#####.##\n..#.###...\n..#.......\n..#.###...\n"
	foodAllergens := "mxmxvkd kfcds sqjhc nhms (contains dairy, fish)\ntrh fvjkl sbzzf mxmxvkd (contains dairy)\n" +
		"sqjhc fvjkl (contains soy)\nsqjhc mxmxvkd sbzzf (contains fish)\n"
	crabCombat := "Player 1:\n9\n2\n6\n3\n1\n\nPlayer 2:\n5\n8\n4\n7\n10"
	crabCups := "389125467\n"
	tiles := "sesenwnenenewseeswwswswwnenewsewsw\nneeenesenwnwwswnenewnwwsewnenwseswesw\n" +
		"seswneswswsenwwnwse\nnwnwneseeswswnenewneswwnewseswneseene\nswweswneswnenwsewnwneneseenw\n" +
		"eesenwseswswnenwswnwnwsewwnwsene\nsewnenenenesenwsewnenwwwse\nwenwwweseeeweswwwnwwe\n" +
		"wsweesenenewnwwnwsenewsenwwsesesenwne\nneeswseenwwswnwswswnw\nnenwswwsewswnenenewsenwsenwnesesenew\n" +
		"enewnwewneswsewnwswenweswnenwsenwsw\nsweneswneswneneenwnewenewwneswswnese\n" +
		"swwesenesewenwneswnwwneseswwne\nenesenwswwswneneswsenwnewswseenwsese\n" +
		"wnwnesenesenenwwnenwsewesewsesesew\nnenewswnwewswnenesenwnesewesw\n" +
		"eneswnwswnwsenenwnwnwwseeswneewsenese\nneswnwewnwnwseenwseesewsenwsweewe\nwseweeenwnesenwwwswnew\n"
	encryptionKeys := "5764801\n17807724\n"

	cases := []struct {
		day    string
		args   []string
		input  string
		output string
	}{
		{"1", nil, expenses, "Part 1: 514579\nPart 2: 241861950\n"},
		{"2", nil, passwords, "Part 1: 2\nPart 2: 1\n"},
		{"3", nil, trees, "Part 1: 7\nPart 2: 336\n"},
		{"4", nil, passports, "Part 1: 2\nPart 2: 2\n"},
		{"5", nil, seatsWithGap, "Part 1: 361\nPart 2: 359\n"},
		{"6", nil, groups, "Part 1: 11\nPart 2: 6\n"},
		{"7", nil, bags, "Part 1: 4\nPart 2: 32\n"},
		{"8", nil, bootCode, "Part 1: 5\nPart 2: 8\n"},
		{"9", []string{"5"}, jumps, "Part 1: 127\nPart 2: 62\n"},
		{"10", nil, adapters, "Part 1: 220\nPart 2: 19208\n"},
		{"11", nil, seatLayout, "Part 1: 37\nPart 2: 26\n"},
		{"12", nil, navigation, "Part 1: 25\nPart 2: 286\n"},
		{"13", nil, buses, "Part 1: 295\nPart 2: 1068781\n"},
		{"14", nil, maskB, "Part 1: 51\nPart 2: 208\n"},
		{"15", nil, memoryGame, "Part 1: 436\nPart 2: 175594\n"},
		{"16", nil, ticketNotesB, "Part 1: 0\nPart 2: 10\n"},
		{"17", nil, seaMonsters, "Part 1: 112\nPart 2: 848\n"},
		{"18", nil, bitmask, "Part 1: 26457\nPart 2: 694173\n"},
		{"19", nil, messagesLoop, "Part 1: 3\nPart 2: 12\n"},
		{"20", nil, jigsawTiles, "Part 1: 20899048083289\nPart 2: 273\n"},
		{"21", nil, foodAllergens, "Part 1: 5\nPart 2: mxmxvkd,sqjhc,fvjkl\n"},
		{"22", nil, crabCombat, "Part 1: 306\nPart 2: 291\n"},
		{"23", nil, crabCups, "Part 1: 67384529\nPart 2: 149245887792\n"},
		{"24", nil, tiles, "Part 1: 10\nPart 2: 2208\n"},
		{"25", nil, encryptionKeys, "Part 1: 14897079\n"},
	}
	for _, tc := range cases {
		t.Run("day"+tc.day, func(t *testing.T) {
			command := exec.CommandContext(ctx, binary, append([]string{tc.day}, tc.args...)...)
			command.Stdin = strings.NewReader(tc.input)

			output, err := command.CombinedOutput()

			require.NoError(t, err, string(output))
			require.Equal(t, tc.output, string(output))
		})
	}
	t.Run("unknown day fails", func(t *testing.T) {
		command := exec.CommandContext(ctx, binary, "99")

		err := command.Run()

		require.Error(t, err)
	})
}
