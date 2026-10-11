package main

import (
	"fmt"
	"strconv"
	"strings"
)

var errInvalidStep = fmt.Errorf("invalid step")

func steps(input string) []string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, ",")
}

func hash(step string) uint32 {
	var value uint32
	for _, b := range []byte(step) {
		value = (value + uint32(b)) * 17 % 256
	}
	return value
}

type lens struct {
	label       string
	focalLength int
}

type instruction struct {
	label  string
	remove bool
	lens   lens
}

func parseInstruction(step string) (instruction, error) {
	if label, ok := strings.CutSuffix(step, "-"); ok {
		if label == "" {
			return instruction{}, fmt.Errorf("%w: %q", errInvalidStep, step)
		}
		return instruction{label: label, remove: true}, nil
	}
	if label, focal, ok := strings.Cut(step, "="); ok {
		if label == "" {
			return instruction{}, fmt.Errorf("%w: %q", errInvalidStep, step)
		}
		focalLength, err := strconv.Atoi(focal)
		if err != nil {
			return instruction{}, fmt.Errorf("%w: %q: %v", errInvalidStep, step, err)
		}
		return instruction{label: label, lens: lens{label: label, focalLength: focalLength}}, nil
	}
	return instruction{}, fmt.Errorf("%w: %q", errInvalidStep, step)
}

// hashMap models the 256-box lens configuration from the initialization manual.
type hashMap struct {
	boxes [256][]lens
}

func (m *hashMap) apply(ins instruction) {
	box := &m.boxes[hash(ins.label)]
	index := -1
	for i, l := range *box {
		if l.label == ins.label {
			index = i
			break
		}
	}
	if ins.remove {
		if index >= 0 {
			*box = append((*box)[:index], (*box)[index+1:]...)
		}
		return
	}
	if index >= 0 {
		(*box)[index] = ins.lens
	} else {
		*box = append(*box, ins.lens)
	}
}

func (m *hashMap) focusingPower() int {
	total := 0
	for boxIndex, box := range m.boxes {
		for slot, l := range box {
			total += (boxIndex + 1) * (slot + 1) * l.focalLength
		}
	}
	return total
}
