package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	defDictPath    = "/usr/share/dict/words"
	defLetterCount = 5
	defExludePunc  = true
)

func main() {
	dict, letterc, err := loadDictionary(mustGetArgs())
	if err != nil {
		log.Fatalf("Could not load dictionary: %s", err)
	}

	guesser := newGuesser(dict, letterc)
	usr := bufio.NewScanner(os.Stdin)
	for {
		possibleMatches := guesser.PossibleMatches()
		fmt.Printf("%d possible remaining matches of %d letters.\n", possibleMatches, letterc)
		recomendation := guesser.Guess()
		fmt.Printf("Recomended guess is: %q\n\n", recomendation)

		switch possibleMatches {
		case 1:
			fmt.Println("Our work here is done!")
			return
		case 0:
			fmt.Println("Something went wrong! Check your results for typos, " +
				"the dictonary is missing something, or there is a bug :)")
			return
		}

		var guess string
		for len(guess) != int(letterc) {
			fmt.Printf("Please enter actual guess (%d letters) or enter for %q: ", letterc, recomendation)
			if usr.Scan() {
				guess = usr.Text()
			}
			if err := usr.Err(); err != nil {
				log.Fatalf("Coult not read user input: %s", err)
			}
			if guess == "" {
				guess = recomendation
			}
		}

		var result string
		for len(result) != int(letterc) {
			fmt.Printf("Please enter result of guess (%d y|n|p): ", letterc)
			if usr.Scan() {
				result = usr.Text()
			}
			if err := usr.Err(); err != nil {
				log.Fatalf("Coult not read user input: %s", err)
			}
		}

		if err := guesser.Narrow(guess, result); err != nil {
			fmt.Printf("Invalid guess or result: %s. Try again!\n", err)
		}
	}

}

func mustGetArgs() (dictPath string, letterc uint, excludePunc bool) {
	var help bool
	flag.UintVar(&letterc, "letters", defLetterCount, "Number of letters in word")
	flag.StringVar(&dictPath, "dict-file", defDictPath, "Path of dictionary file")
	flag.BoolVar(&excludePunc, "no-punc", defExludePunc, "Should words with puncuation be excluded")
	flag.BoolVar(&help, "help", false, "Show usage and exit")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	return
}

type Dictionary map[string]struct{}

func loadDictionary(dictPath string, letterc uint, excludePunc bool) (Dictionary, uint, error) {
	fin, err := os.Open(dictPath)
	if err != nil {
		return nil, 0, fmt.Errorf("Could not open dictionary file %q: %w", dictPath, err)
	}
	defer fin.Close()

	keep := func(word []byte, excludePunc bool) bool {
		for _, char := range word {
			switch {
			case char < 'a', char > 'z', excludePunc && char == '\'':
				return false
			}
		}
		return true
	}

	count := int(letterc)
	dict := make(Dictionary, 6000)
	scan := bufio.NewScanner(fin)
	for scan.Scan() {
		word := scan.Bytes()
		if len(word) != count {
			continue
		}

		word = bytes.ToLower(word)
		if keep(word, excludePunc) {
			dict[string(word)] = struct{}{}
		}
	}
	if err = scan.Err(); err != nil {
		return nil, 0, fmt.Errorf("Could not read dictionary file %q: %w", dictPath, err)
	}

	return dict, letterc, nil
}

type Guesser struct {
	letterc uint
	opts    []map[byte]uint
	freq    map[byte]uint
	dict    Dictionary
}

func newGuesser(dict Dictionary, letterc uint) *Guesser {
	g := &Guesser{
		letterc: letterc,
		opts:    make([]map[byte]uint, letterc),
		dict:    dict,
	}

	return g
}

func (g *Guesser) Guess() string {
	g.analyze()

	var cannidate string
	var highscore, freqscore uint

	for word, _ := range g.dict {
		var score, freq uint
		for pos, letter := range word {
			char := byte(letter)
			score += g.opts[pos][char]
			freq += g.freq[char]
		}

		if score > highscore || cannidate == "" {
			cannidate, highscore, freqscore = word, score, freq
		} else if score == highscore {
			if freq > freqscore {
				cannidate, highscore, freqscore = word, score, freq
			}
		}
	}

	return cannidate
}

func (g *Guesser) Narrow(guess string, result string) error {
	switch {
	case len(guess) != int(g.letterc):
		return fmt.Errorf("guess: %q, was not %d letters long", guess, g.letterc)
	case len(guess) != len(result):
		return fmt.Errorf("guess: %q, and result %q are not the same length", guess, result)
	}

	result = strings.ToLower(result)
	seen := make(map[byte]uint)
	nos := make(map[byte]uint)
	for i, letter := range guess {
		char := byte(letter)
		switch result[i] {
		case 'y':
			seen[char]++
			for word, _ := range g.dict {
				if word[i] != char {
					delete(g.dict, word)
				}
			}

		case 'n':
			nos[char]++
			for word, _ := range g.dict {
				if word[i] == char {
					delete(g.dict, word)
				}
			}

		case 'p':
			seen[char]++
			for word, _ := range g.dict {
				if !strings.ContainsRune(word, letter) || word[i] == char {
					delete(g.dict, word)
				}
			}

		default:
			return fmt.Errorf("Result should be a series of y (yes match), "+
				" n (no match), p (partial mach, wrong position), but %q was found", result[i])
		}
	}

	for char, _ := range nos {
		for word, _ := range g.dict {
			count := uint(strings.Count(word, string(char)))
			if count < 1 {
				continue
			}

			if seen[char] < 1 || count > seen[char] {
				delete(g.dict, word)
			}
		}
	}

	return nil
}

func (g *Guesser) PossibleMatches() uint {
	return uint(len(g.dict))
}

func (g *Guesser) analyze() {
	for i := range g.opts {
		g.opts[i] = make(map[byte]uint)
	}
	g.freq = make(map[byte]uint)

	for word, _ := range g.dict {
		for pos, letter := range word {
			char := byte(letter)
			g.opts[pos][char]++
			g.freq[char]++
		}
	}
}
