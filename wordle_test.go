package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestGames(t *testing.T) {
	const warnDifGuess = false

	for i, tcase := range games {
		letterc := uint(len(tcase[0][0]))
		dict, dicLetterc, err := loadDictionary(defDictPath, letterc)
		if err != nil {
			t.Fatalf("Could not load dictionary from %q (default config): %s", defDictPath, err)
		} else if letterc != dicLetterc {
			t.Fatalf("Dictionary returned letter count of %d rather than %d as expected", dicLetterc, letterc)
		}
		guesser := newGuesser(dict, letterc)

		t.Run(fmt.Sprintf("game-%d", i), func(t *testing.T) {
			for _, pair := range tcase {
				guess, result := pair[0], pair[1]
				rec := guesser.Guess()

				if result == "" {
					if remaining := guesser.PossibleMatches(); remaining > 1 {
						t.Fatalf("%d remaining possible matches (%q recomended) rather than just %q", remaining, rec, guess)
					}
					if rec != guess {
						t.Fatalf("Computed guess was %q, but %q was expected", rec, guess)
					}
					break
				}

				if warnDifGuess && rec != guess {
					t.Logf("Warning: Recomendation was %q but guess was %q", rec, guess)
				}

				if err = guesser.Narrow(guess, result); err != nil {
					t.Fatalf("Narrow failed: %s", err)
				}
			}
		})
	}
}

var games = [][][2]string{
	[][2]string{ //1
		[2]string{"sores", "nnppn"},
		[2]string{"trace", "ppnnp"},
		[2]string{"retry", "yypnn"},
		[2]string{"refit", "yynny"},
		[2]string{"rebut", ""},
	},
	[][2]string{ //1
		[2]string{"sales", "nnnpn"},
		[2]string{"deice", "nynnn"},
		[2]string{"terry", "pypnn"},
		[2]string{"rebut", ""},
	},
	[][2]string{ //2
		[2]string{"sores", "ynnnp"},
		[2]string{"slash", "ynnyn"},
		[2]string{"sissy", "yyyyy"},
		[2]string{"sissy", ""},
	},
	[][2]string{ //222
		[2]string{"sores", "nynnn"},
		[2]string{"loony", "nynyn"},
		[2]string{"found", "nyyyn"},
		[2]string{"mount", "yyyyy"},
	},
	[][2]string{ //222
		[2]string{"sales", "nnnnn"},
		[2]string{"corny", "nynyn"},
		[2]string{"mound", "yyyyn"},
		[2]string{"mount", ""},
	},
	[][2]string{ //223
		[2]string{"sales", "nnnpn"},
		[2]string{"deice", "nynnn"},
		[2]string{"perry", "yyyny"},
		[2]string{"perky", ""},
	},
	[][2]string{ //223
		[2]string{"sores", "nnypn"},
		[2]string{"eerie", "nyynn"},
		[2]string{"berry", "nyyny"},
		[2]string{"nerdy", "nyyny"},
		[2]string{"perky", "yyyyy"},
		[2]string{"perky", ""},
	},
	[][2]string{ //224
		[2]string{"sales", "nnpnn"},
		[2]string{"clint", "ypnnn"},
		[2]string{"coyly", "yynyn"},
		[2]string{"could", ""},
	},
	[][2]string{ //225
		[2]string{"sales", "nnnnn"},
		[2]string{"bonny", "nnnyn"},
		[2]string{"wring", "yynyy"},
		[2]string{"wrung", ""},
	},
	[][2]string{ //227
		[2]string{"sores", "ppnpn"},
		[2]string{"chose", "nyyyy"},
		[2]string{"those", "yyyyy"},
		[2]string{"those", ""},
	},
	[][2]string{ //228
		[2]string{"sales", "pnnnn"},
		[2]string{"moist", "yyyyy"},
		[2]string{"moist", ""},
	},
	[][2]string{ //229
		[2]string{"sores", "ynpnn"},
		[2]string{"shark", "yyyyn"},
		[2]string{"sharp", "yyyyn"},
		[2]string{"shard", ""},
	},
}

func TestPlayEntireDictionary(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	dict, letterc, err := loadDictionary(defDictPath, defLetterCount)
	if err != nil {
		t.Fatalf("Could not load dictionary from %q (default config): %s", defDictPath, err)
	}

	//[]scorefunc{scoreBasic, scoreWFreq, scoreBasicSquared, scoreWeightOptLeft, scoreWeightOptElim}
	for i, scorer := range []scorefunc{scoreBasic} {
		t.Run(fmt.Sprintf("scorer-%d", i), func(t *testing.T) {
			totalGuesses := 0
			for word := range dict {
				gameDict := dict.Clone()
				guesser, guessCount, guess := newGuesser(gameDict, letterc), 0, ""
				guesser.score = scorer

				for {
					if _, containsWord := gameDict[word]; !containsWord {
						t.Fatalf("%q was removed from dictionary!", word)
					}

					if guesser.PossibleMatches() < 1 {
						t.Fatalf("Failed to solve: %q", word)
					}

					guessCount++
					if guess = guesser.Guess(); guess == word {
						break
					}

					result := judge(guess, word)
					//println("guess:", guess, "word:", word, "result:", result)
					if strings.Count(result, "y") == int(letterc) {
						break
					}

					if err = guesser.Narrow(guess, result); err != nil {
						t.Fatalf("Narrow failed: %s", err)
					}
				}
				totalGuesses += guessCount
			}
			t.Logf("%.2f turns per game\n", float64(totalGuesses)/float64(len(dict)))
		})
	}
}

func judge(guess, word string) string {
	result := make([]byte, len(word))
	matchChars := make(map[byte]int)

	for i := range guess {
		guessChar := guess[i]

		switch {
		case guessChar == word[i]:
			matchChars[guessChar]++
			result[i] = 'y'

		case !strings.Contains(word, string(guessChar)):
			result[i] = 'n'
		}
	}

	for i, cur := range result {
		if cur != 0 {
			continue
		}

		guessChar := guess[i]
		matchCount := matchChars[guessChar]
		if matchCount < strings.Count(word, string(guessChar)) {
			matchChars[guessChar]++
			result[i] = 'p'
			continue
		}

		result[i] = 'n'
	}

	return string(result)
}
