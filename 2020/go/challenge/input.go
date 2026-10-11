package challenge

import (
	"bufio"
	"io"
	"strings"
)

type Input struct {
	scanner *bufio.Scanner

	lines chan string
}

func FromLiteral(input string) *Input {
	return newInputFromReader(strings.NewReader(input), nil)
}

func newInputFromReader(r io.Reader, c io.Closer) *Input {
	result := &Input{
		scanner: bufio.NewScanner(r),
		lines:   make(chan string),
	}

	go func() {
		defer func() {
			if c != nil {
				_ = c.Close()
			}
		}()

		for result.scanner.Scan() {
			result.lines <- result.scanner.Text()
		}

		close(result.lines)
	}()

	return result
}

func (c *Input) Lines() <-chan string {
	return c.lines
}

func (c *Input) LinesSlice() (result []string) {
	for line := range c.Lines() {
		result = append(result, line)
	}

	return
}

func (c *Input) ToGrid() *Grid {
	lines := c.LinesSlice()

	return &Grid{
		tiles:  lines,
		width:  len(lines[0]),
		height: len(lines),
	}
}
