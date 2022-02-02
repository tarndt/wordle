package main

import (
	"fmt"
	"testing"
)

func TestGames(t *testing.T) {
	tcases := [][][2]string{
		[][2]string{ //223
			[2]string{"sales", "nnnpn"},
			[2]string{"deice", "nynnn"},
			[2]string{"perry", "yyyny"},
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
		[][2]string{ //228
			[2]string{"sales", "pnnnn"},
			[2]string{"moist", "yyyyy"},
			[2]string{"moist", ""},
		},
	}

	for i, tcase := range tcases {
		letterc := uint(len(tcase[0][0]))
		dict, dicLetterc, err := loadDictionary(defDictPath, letterc, defExludePunc)
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
				println("rec", rec)

				if result == "" {
					if remaining := guesser.PossibleMatches(); remaining > 1 {
						t.Fatalf("%d remaining possible matches (%q recomended) rather than just %q", remaining, rec, guess)
					}
					if rec != guess {
						t.Fatalf("Computed guess was %q, but %q was expected", rec, guess)
					}
					break
				}

				if rec != guess {
					t.Logf("Warning: Recomendation was %q but guess was %q", rec, guess)
				}
				println("guessing", guess)
				if err = guesser.Narrow(guess, result); err != nil {
					t.Fatalf("Narrow failed: %s", err)
				}
			}
		})
	}
}
