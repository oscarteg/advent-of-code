package day6

import "fmt"

type group struct {
	people []map[rune]struct{}
}

func (g group) anyAnswered() map[rune]struct{} {
	result := make(map[rune]struct{})
	for _, person := range g.people {
		for letter := range person {
			result[letter] = struct{}{}
		}
	}
	return result
}

func (g group) allAnswered() map[rune]struct{} {
	result := make(map[rune]struct{})
	for letter := range g.people[0] {
		result[letter] = struct{}{}
	}
	for _, person := range g.people[1:] {
		for letter := range result {
			if _, ok := person[letter]; !ok {
				delete(result, letter)
			}
		}
	}
	return result
}

func parsePerson(line string) (map[rune]struct{}, error) {
	answers := make(map[rune]struct{})
	for _, letter := range line {
		if letter < 'a' || letter > 'z' {
			return nil, fmt.Errorf("invalid answer letter %q", letter)
		}
		answers[letter] = struct{}{}
	}
	return answers, nil
}

func parseGroups(lines []string) ([]group, error) {
	var groups []group
	var current group
	for _, line := range lines {
		if line == "" {
			if len(current.people) == 0 {
				return nil, fmt.Errorf("empty group")
			}
			groups = append(groups, current)
			current = group{}
			continue
		}
		person, err := parsePerson(line)
		if err != nil {
			return nil, err
		}
		current.people = append(current.people, person)
	}
	if len(current.people) == 0 {
		return nil, fmt.Errorf("empty group")
	}
	groups = append(groups, current)
	return groups, nil
}
