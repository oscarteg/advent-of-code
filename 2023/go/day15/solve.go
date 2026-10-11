package main

func solve(input string) (part1, part2 int, err error) {
	for _, step := range steps(input) {
		part1 += int(hash(step))
	}

	m := &hashMap{}
	for _, step := range steps(input) {
		ins, parseErr := parseInstruction(step)
		if parseErr != nil {
			return 0, 0, parseErr
		}
		m.apply(ins)
	}
	part2 = m.focusingPower()

	return part1, part2, nil
}
