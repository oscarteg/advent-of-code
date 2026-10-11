use crate::day07::custom_error::AocError;
use std::cmp::Ordering;

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
pub enum Card {
    Two = 0,
    Three = 1,
    Four = 2,
    Five = 3,
    Six = 4,
    Seven = 5,
    Eight = 6,
    Nine = 7,
    Ten = 8,
    Jack = 9,
    Queen = 10,
    King = 11,
    Ace = 12,
}

impl Card {
    pub fn from_char(c: char) -> Result<Self, AocError> {
        match c {
            '2' => Ok(Card::Two),
            '3' => Ok(Card::Three),
            '4' => Ok(Card::Four),
            '5' => Ok(Card::Five),
            '6' => Ok(Card::Six),
            '7' => Ok(Card::Seven),
            '8' => Ok(Card::Eight),
            '9' => Ok(Card::Nine),
            'T' => Ok(Card::Ten),
            'J' => Ok(Card::Jack),
            'Q' => Ok(Card::Queen),
            'K' => Ok(Card::King),
            'A' => Ok(Card::Ace),
            _ => Err(AocError::InvalidInput("invalid card")),
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
pub enum Joker {
    Joker = 0,
    Two = 1,
    Three = 2,
    Four = 3,
    Five = 4,
    Six = 5,
    Seven = 6,
    Eight = 7,
    Nine = 8,
    Ten = 9,
    Queen = 10,
    King = 11,
    Ace = 12,
}

impl Joker {
    pub fn from_char(c: char) -> Result<Self, AocError> {
        match c {
            '2' => Ok(Joker::Two),
            '3' => Ok(Joker::Three),
            '4' => Ok(Joker::Four),
            '5' => Ok(Joker::Five),
            '6' => Ok(Joker::Six),
            '7' => Ok(Joker::Seven),
            '8' => Ok(Joker::Eight),
            '9' => Ok(Joker::Nine),
            'T' => Ok(Joker::Ten),
            'J' => Ok(Joker::Joker),
            'Q' => Ok(Joker::Queen),
            'K' => Ok(Joker::King),
            'A' => Ok(Joker::Ace),
            _ => Err(AocError::InvalidInput("invalid card")),
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
pub enum HandType {
    HighCard = 0,
    OnePair = 1,
    TwoPair = 2,
    ThreeOfAKind = 3,
    FullHouse = 4,
    FourOfAKind = 5,
    FiveOfAKind = 6,
}

pub fn classify_hand(cards: &[Card; 5]) -> HandType {
    let mut counts = [0; 13];
    for card in cards {
        counts[*card as usize] += 1;
    }
    let mut frequencies: Vec<_> = counts.iter().filter(|&&c| c > 0).copied().collect();
    frequencies.sort_by_key(|&c| std::cmp::Reverse(c));

    match frequencies.as_slice() {
        [5] => HandType::FiveOfAKind,
        [4, 1] => HandType::FourOfAKind,
        [3, 2] => HandType::FullHouse,
        [3, 1, 1] => HandType::ThreeOfAKind,
        [2, 2, 1] => HandType::TwoPair,
        [2, 1, 1, 1] => HandType::OnePair,
        _ => HandType::HighCard,
    }
}

pub fn classify_hand_joker(cards: &[Joker; 5]) -> HandType {
    let mut counts = [0; 13];
    let mut joker_count = 0;
    for card in cards {
        if *card == Joker::Joker {
            joker_count += 1;
        } else {
            counts[*card as usize] += 1;
        }
    }

    if joker_count == 5 {
        return HandType::FiveOfAKind;
    }

    let mut frequencies: Vec<_> = counts.iter().filter(|&&c| c > 0).copied().collect();
    frequencies.sort_by_key(|&c| std::cmp::Reverse(c));

    if !frequencies.is_empty() {
        frequencies[0] += joker_count;
        frequencies.sort_by_key(|&c| std::cmp::Reverse(c));
    } else {
        frequencies.push(joker_count);
    }

    match frequencies.as_slice() {
        [5, ..] => HandType::FiveOfAKind,
        [4, 1, ..] => HandType::FourOfAKind,
        [3, 2, ..] => HandType::FullHouse,
        [3, 1, 1, ..] => HandType::ThreeOfAKind,
        [2, 2, 1, ..] => HandType::TwoPair,
        [2, 1, 1, 1, ..] => HandType::OnePair,
        _ => HandType::HighCard,
    }
}

#[derive(Debug, Clone)]
pub struct Hand {
    pub cards: [Card; 5],
    pub hand_type: HandType,
    pub bid: u32,
}

impl Hand {
    pub fn parse(line: &str) -> Result<Self, AocError> {
        let parts: Vec<&str> = line.split_whitespace().collect();
        if parts.len() != 2 {
            return Err(AocError::InvalidHandFormat);
        }

        let cards_str = parts[0];
        if cards_str.len() != 5 {
            return Err(AocError::InvalidHandFormat);
        }

        let mut cards = [Card::Two; 5];
        for (i, c) in cards_str.chars().enumerate() {
            cards[i] = Card::from_char(c)?;
        }

        let hand_type = classify_hand(&cards);
        let bid = parts[1].parse()?;

        Ok(Hand {
            cards,
            hand_type,
            bid,
        })
    }
}

impl Ord for Hand {
    fn cmp(&self, other: &Self) -> Ordering {
        if self.hand_type != other.hand_type {
            return self.hand_type.cmp(&other.hand_type);
        }
        self.cards.cmp(&other.cards)
    }
}

impl PartialOrd for Hand {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl Eq for Hand {}

impl PartialEq for Hand {
    fn eq(&self, other: &Self) -> bool {
        self.hand_type == other.hand_type && self.cards == other.cards
    }
}

#[derive(Debug, Clone)]
pub struct HandJoker {
    pub cards: [Joker; 5],
    pub hand_type: HandType,
    pub bid: u32,
}

impl HandJoker {
    pub fn parse(line: &str) -> Result<Self, AocError> {
        let parts: Vec<&str> = line.split_whitespace().collect();
        if parts.len() != 2 {
            return Err(AocError::InvalidHandFormat);
        }

        let cards_str = parts[0];
        if cards_str.len() != 5 {
            return Err(AocError::InvalidHandFormat);
        }

        let mut cards = [Joker::Two; 5];
        for (i, c) in cards_str.chars().enumerate() {
            cards[i] = Joker::from_char(c)?;
        }

        let hand_type = classify_hand_joker(&cards);
        let bid = parts[1].parse()?;

        Ok(HandJoker {
            cards,
            hand_type,
            bid,
        })
    }
}

impl Ord for HandJoker {
    fn cmp(&self, other: &Self) -> Ordering {
        if self.hand_type != other.hand_type {
            return self.hand_type.cmp(&other.hand_type);
        }
        self.cards.cmp(&other.cards)
    }
}

impl PartialOrd for HandJoker {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl Eq for HandJoker {}

impl PartialEq for HandJoker {
    fn eq(&self, other: &Self) -> bool {
        self.hand_type == other.hand_type && self.cards == other.cards
    }
}
