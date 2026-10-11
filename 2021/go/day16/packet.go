package main

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

var errInvalidTransmission = errors.New("invalid transmission")

// Packet type IDs. Type ID 4 marks a literal value; every other ID is an operator.
const (
	typeSum uint64 = iota
	typeProduct
	typeMinimum
	typeMaximum
	typeLiteral
	typeGreaterThan
	typeLessThan
	typeEqualTo
)

// bitReader is a bit cursor over the hexadecimal digits of a transmission.
type bitReader struct {
	nibbles  []byte
	position int
}

func newBitReader(input string) (*bitReader, error) {
	digits := strings.TrimSpace(input)
	if digits == "" {
		return nil, fmt.Errorf("%w: it is empty", errInvalidTransmission)
	}
	nibbles := make([]byte, 0, len(digits))
	for _, digit := range digits {
		value, err := strconv.ParseUint(string(digit), 16, 8)
		if err != nil {
			return nil, fmt.Errorf("hexadecimal digit %q: %w", digit, err)
		}
		nibbles = append(nibbles, byte(value))
	}
	return &bitReader{nibbles: nibbles}, nil
}

func (r *bitReader) length() int {
	return len(r.nibbles) * 4
}

func (r *bitReader) bit(index int) uint64 {
	return uint64(r.nibbles[index/4]>>(3-index%4)) & 1
}

func (r *bitReader) take(count int) (uint64, error) {
	end := r.position + count
	if end > r.length() {
		return 0, fmt.Errorf("%w: it ends inside a packet", errInvalidTransmission)
	}
	var value uint64
	for ; r.position < end; r.position++ {
		value = value<<1 | r.bit(r.position)
	}
	return value, nil
}

// restIsPadding reports whether only the zero bits padding the final nibble remain.
func (r *bitReader) restIsPadding() bool {
	for index := r.position; index < r.length(); index++ {
		if r.bit(index) != 0 {
			return false
		}
	}
	return true
}

// packet is a decoded packet: a literal value when typeID is typeLiteral,
// otherwise an operator over its sub-packets.
type packet struct {
	version  uint64
	typeID   uint64
	literal  uint64
	operands []packet
}

func parsePacket(r *bitReader) (packet, error) {
	version, err := r.take(3)
	if err != nil {
		return packet{}, err
	}
	typeID, err := r.take(3)
	if err != nil {
		return packet{}, err
	}
	decoded := packet{version: version, typeID: typeID}
	if typeID == typeLiteral {
		decoded.literal, err = parseLiteral(r)
	} else {
		decoded.operands, err = parseOperands(r)
	}
	if err != nil {
		return packet{}, err
	}
	return decoded, nil
}

// parseLiteral reads five-bit groups, each but the last flagged with a leading bit.
func parseLiteral(r *bitReader) (uint64, error) {
	var value uint64
	for {
		group, err := r.take(5)
		if err != nil {
			return 0, err
		}
		if value > math.MaxUint64/16 {
			return 0, fmt.Errorf("%w: a literal value overflows", errInvalidTransmission)
		}
		value = value<<4 | group&0b1111
		if group&0b1_0000 == 0 {
			return value, nil
		}
	}
}

func parseOperands(r *bitReader) ([]packet, error) {
	lengthType, err := r.take(1)
	if err != nil {
		return nil, err
	}
	if lengthType == 1 {
		count, err := r.take(11)
		if err != nil {
			return nil, err
		}
		operands := make([]packet, 0, count)
		for range count {
			operand, err := parsePacket(r)
			if err != nil {
				return nil, err
			}
			operands = append(operands, operand)
		}
		return operands, nil
	}
	length, err := r.take(15)
	if err != nil {
		return nil, err
	}
	end := r.position + int(length)
	if end > r.length() {
		return nil, fmt.Errorf("%w: sub-packets run past it", errInvalidTransmission)
	}
	var operands []packet
	for r.position < end {
		operand, err := parsePacket(r)
		if err != nil {
			return nil, err
		}
		operands = append(operands, operand)
	}
	if r.position != end {
		return nil, fmt.Errorf("%w: sub-packets overrun their declared length", errInvalidTransmission)
	}
	return operands, nil
}

func (p packet) versionSum() uint64 {
	sum := p.version
	for _, operand := range p.operands {
		sum += operand.versionSum()
	}
	return sum
}

func (p packet) value() (uint64, error) {
	if p.typeID == typeLiteral {
		return p.literal, nil
	}
	values := make([]uint64, len(p.operands))
	for index, operand := range p.operands {
		value, err := operand.value()
		if err != nil {
			return 0, err
		}
		values[index] = value
	}
	return apply(p.typeID, values)
}

func apply(typeID uint64, values []uint64) (uint64, error) {
	switch typeID {
	case typeSum:
		var total uint64
		for _, value := range values {
			if total > math.MaxUint64-value {
				return 0, fmt.Errorf("%w: a sum overflows", errInvalidTransmission)
			}
			total += value
		}
		return total, nil
	case typeProduct:
		total := uint64(1)
		for _, value := range values {
			if value != 0 && total > math.MaxUint64/value {
				return 0, fmt.Errorf("%w: a product overflows", errInvalidTransmission)
			}
			total *= value
		}
		return total, nil
	case typeMinimum, typeMaximum:
		if len(values) == 0 {
			return 0, fmt.Errorf("%w: an extremum needs at least one sub-packet", errInvalidTransmission)
		}
		if typeID == typeMinimum {
			return slices.Min(values), nil
		}
		return slices.Max(values), nil
	case typeGreaterThan, typeLessThan, typeEqualTo:
		if len(values) != 2 {
			return 0, fmt.Errorf("%w: a comparison needs exactly two sub-packets", errInvalidTransmission)
		}
		left, right := values[0], values[1]
		switch typeID {
		case typeGreaterThan:
			return boolValue(left > right), nil
		case typeLessThan:
			return boolValue(left < right), nil
		default:
			return boolValue(left == right), nil
		}
	default:
		return 0, fmt.Errorf("%w: unknown packet type %d", errInvalidTransmission, typeID)
	}
}

func boolValue(condition bool) uint64 {
	if condition {
		return 1
	}
	return 0
}

// solve returns the sum of every packet version and the value of the outermost packet.
func solve(input string) (part1, part2 uint64, err error) {
	reader, err := newBitReader(input)
	if err != nil {
		return 0, 0, err
	}
	outermost, err := parsePacket(reader)
	if err != nil {
		return 0, 0, err
	}
	if !reader.restIsPadding() {
		return 0, 0, fmt.Errorf("%w: it has data after the outermost packet", errInvalidTransmission)
	}
	value, err := outermost.value()
	if err != nil {
		return 0, 0, err
	}
	return outermost.versionSum(), value, nil
}
