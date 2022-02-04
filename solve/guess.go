package solve

import (
	"fmt"
	"strings"
)

type Guesser struct {
	letterc uint
	opts    []map[byte]uint
	freq    map[byte]uint
	dict    Dictionary
	score   scorefunc
}

func NewGuesser(dict Dictionary, letterc uint) *Guesser {
	g := &Guesser{
		letterc: letterc,
		opts:    make([]map[byte]uint, letterc),
		dict:    dict,
		score:   scoreBasic,
	}

	return g
}

func (g *Guesser) PossibleMatches() uint {
	return uint(len(g.dict))
}

var errNotResult = fmt.Errorf("Result should be a series of y (yes match), n (no match), p (partial mach, wrong position)")

func (g *Guesser) Narrow(guess string, result string) error {
	switch {
	case len(guess) != int(g.letterc):
		return fmt.Errorf("guess: %q, was not %d letters long", guess, g.letterc)
	case len(guess) != len(result):
		return fmt.Errorf("guess: %q, and result %q are not the same length", guess, result)
	case !ValidResult(result):
		return errNotResult
	}

	result = strings.ToLower(result)
	seen := make(map[byte]uint)
	nos := make(map[byte]uint)
	for i, letter := range guess {
		char := byte(letter)
		switch result[i] {
		case 'y':
			seen[char]++
			for word := range g.dict {
				if word[i] != char {
					delete(g.dict, word)
				}
			}

		case 'n':
			nos[char]++
			for word := range g.dict {
				if word[i] == char {
					delete(g.dict, word)
				}
			}

		case 'p':
			seen[char]++
			for word := range g.dict {
				if !strings.ContainsRune(word, letter) || word[i] == char {
					delete(g.dict, word)
				}
			}

		default:
			return fmt.Errorf("%q was unexpected; %w", char, errNotResult)
		}
	}

	//Remove all the non-straight forward no results; those words with the no
	// letter in a different spot than the no letter in the guess
	for char := range nos {
		for word := range g.dict {
			count := uint(strings.Count(word, string(char)))
			if count < 1 {
				continue
			}

			if seenCount := seen[char]; seenCount < 1 || count > seenCount {
				delete(g.dict, word)
			}
		}
	}

	//If we got multiple y or p for the same char, ensure the words all have enough
	for char, seenCount := range seen {
		for word := range g.dict {
			if uint(strings.Count(word, string(char))) < seenCount {
				delete(g.dict, word)
			}
		}
	}

	return nil
}

func (g *Guesser) Guess() string {
	g.analyze()

	var cannidate string
	var highscore, freqscore uint

	for word := range g.dict {
		var score, freq uint
		for pos, letter := range word {
			char := byte(letter)
			score += g.score(g.opts, g.freq, pos, char)
			freq += g.freq[char]
		}

		if score > highscore || cannidate == "" {
			cannidate, highscore, freqscore = word, score, freq
		} else if score == highscore && freq > freqscore {
			cannidate, highscore, freqscore = word, score, freq
		}
	}

	//Manual optimization for turn 1 now... generalization in progress
	if len(g.dict) > 500 && cannidate == "sores" {
		if _, contains := g.dict["sales"]; contains {
			return "sales"
		}
	}

	return cannidate
}

func (g *Guesser) analyze() {
	for i := range g.opts {
		g.opts[i] = make(map[byte]uint)
	}
	g.freq = make(map[byte]uint)

	for word := range g.dict {
		for pos, letter := range word {
			char := byte(letter)
			g.opts[pos][char]++
			g.freq[char]++
		}
	}
}

func ValidResult(result string) bool {
	if result == "" {
		return false
	}

	for _, char := range result {
		switch char {
		case 'n', 'p', 'y':
			continue
		default:
			return false
		}
	}
	return true
}
