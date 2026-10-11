package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Card int

const (
	Two Card = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func cardFromRune(r rune) (Card, error) {
	switch r {
	case '2':
		return Two, nil
	case '3':
		return Three, nil
	case '4':
		return Four, nil
	case '5':
		return Five, nil
	case '6':
		return Six, nil
	case '7':
		return Seven, nil
	case '8':
		return Eight, nil
	case '9':
		return Nine, nil
	case 'T':
		return Ten, nil
	case 'J':
		return Jack, nil
	case 'Q':
		return Queen, nil
	case 'K':
		return King, nil
	case 'A':
		return Ace, nil
	default:
		return 0, fmt.Errorf("invalid card: %c", r)
	}
}

type Joker int

const (
	JokerCard Joker = iota
	JokerTwo
	JokerThree
	JokerFour
	JokerFive
	JokerSix
	JokerSeven
	JokerEight
	JokerNine
	JokerTen
	JokerQueen
	JokerKing
	JokerAce
)

func jokerFromRune(r rune) (Joker, error) {
	switch r {
	case '2':
		return JokerTwo, nil
	case '3':
		return JokerThree, nil
	case '4':
		return JokerFour, nil
	case '5':
		return JokerFive, nil
	case '6':
		return JokerSix, nil
	case '7':
		return JokerSeven, nil
	case '8':
		return JokerEight, nil
	case '9':
		return JokerNine, nil
	case 'T':
		return JokerTen, nil
	case 'J':
		return JokerCard, nil
	case 'Q':
		return JokerQueen, nil
	case 'K':
		return JokerKing, nil
	case 'A':
		return JokerAce, nil
	default:
		return 0, fmt.Errorf("invalid card: %c", r)
	}
}

type HandType int

const (
	HighCard HandType = iota
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

func classifyHand(cards [5]Card) HandType {
	counts := make(map[Card]int)
	for _, c := range cards {
		counts[c]++
	}

	var freqs []int
	for _, count := range counts {
		freqs = append(freqs, count)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(freqs)))

	switch {
	case len(freqs) == 1:
		return FiveOfAKind
	case len(freqs) == 2 && freqs[0] == 4:
		return FourOfAKind
	case len(freqs) == 2 && freqs[0] == 3:
		return FullHouse
	case len(freqs) == 3 && freqs[0] == 3:
		return ThreeOfAKind
	case len(freqs) == 3 && freqs[0] == 2:
		return TwoPair
	case len(freqs) == 4:
		return OnePair
	default:
		return HighCard
	}
}

func classifyHandJoker(cards [5]Joker) HandType {
	counts := make(map[Joker]int)
	jokerCount := 0
	for _, c := range cards {
		if c == JokerCard {
			jokerCount++
		} else {
			counts[c]++
		}
	}

	if jokerCount == 5 {
		return FiveOfAKind
	}

	var freqs []int
	for _, count := range counts {
		freqs = append(freqs, count)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(freqs)))

	if len(freqs) == 0 {
		freqs = []int{jokerCount}
	} else {
		freqs[0] += jokerCount
		sort.Sort(sort.Reverse(sort.IntSlice(freqs)))
	}

	switch {
	case freqs[0] == 5:
		return FiveOfAKind
	case freqs[0] == 4:
		return FourOfAKind
	case freqs[0] == 3 && len(freqs) > 1 && freqs[1] == 2:
		return FullHouse
	case freqs[0] == 3:
		return ThreeOfAKind
	case freqs[0] == 2 && len(freqs) > 1 && freqs[1] == 2:
		return TwoPair
	case freqs[0] == 2:
		return OnePair
	default:
		return HighCard
	}
}

type Hand struct {
	cards    [5]Card
	handType HandType
	bid      int
}

func parseHand(line string) (Hand, error) {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return Hand{}, fmt.Errorf("invalid hand format")
	}

	cardsStr := parts[0]
	if len(cardsStr) != 5 {
		return Hand{}, fmt.Errorf("invalid hand format")
	}

	var cards [5]Card
	for i, r := range cardsStr {
		c, err := cardFromRune(r)
		if err != nil {
			return Hand{}, err
		}
		cards[i] = c
	}

	bid, err := strconv.Atoi(parts[1])
	if err != nil {
		return Hand{}, fmt.Errorf("invalid bid: %w", err)
	}

	handType := classifyHand(cards)
	return Hand{cards: cards, handType: handType, bid: bid}, nil
}

type HandJoker struct {
	cards    [5]Joker
	handType HandType
	bid      int
}

func parseHandJoker(line string) (HandJoker, error) {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return HandJoker{}, fmt.Errorf("invalid hand format")
	}

	cardsStr := parts[0]
	if len(cardsStr) != 5 {
		return HandJoker{}, fmt.Errorf("invalid hand format")
	}

	var cards [5]Joker
	for i, r := range cardsStr {
		c, err := jokerFromRune(r)
		if err != nil {
			return HandJoker{}, err
		}
		cards[i] = c
	}

	bid, err := strconv.Atoi(parts[1])
	if err != nil {
		return HandJoker{}, fmt.Errorf("invalid bid: %w", err)
	}

	handType := classifyHandJoker(cards)
	return HandJoker{cards: cards, handType: handType, bid: bid}, nil
}

type Hands []Hand

func (h Hands) Len() int      { return len(h) }
func (h Hands) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h Hands) Less(i, j int) bool {
	if h[i].handType != h[j].handType {
		return h[i].handType < h[j].handType
	}
	for k := 0; k < 5; k++ {
		if h[i].cards[k] != h[j].cards[k] {
			return h[i].cards[k] < h[j].cards[k]
		}
	}
	return false
}

type HandsJoker []HandJoker

func (h HandsJoker) Len() int      { return len(h) }
func (h HandsJoker) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h HandsJoker) Less(i, j int) bool {
	if h[i].handType != h[j].handType {
		return h[i].handType < h[j].handType
	}
	for k := 0; k < 5; k++ {
		if h[i].cards[k] != h[j].cards[k] {
			return h[i].cards[k] < h[j].cards[k]
		}
	}
	return false
}

func solve(input string) (int, int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	var hands Hands
	for _, line := range lines {
		if line == "" {
			continue
		}
		hand, err := parseHand(line)
		if err != nil {
			return 0, 0, err
		}
		hands = append(hands, hand)
	}

	sort.Sort(hands)
	part1 := 0
	for i, hand := range hands {
		part1 += hand.bid * (i + 1)
	}

	var handsJoker HandsJoker
	for _, line := range lines {
		if line == "" {
			continue
		}
		hand, err := parseHandJoker(line)
		if err != nil {
			return 0, 0, err
		}
		handsJoker = append(handsJoker, hand)
	}

	sort.Sort(handsJoker)
	part2 := 0
	for i, hand := range handsJoker {
		part2 += hand.bid * (i + 1)
	}

	return part1, part2, nil
}
