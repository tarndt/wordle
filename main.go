package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/tarndt/wordle/conf"
	"github.com/tarndt/wordle/solve"
)

func main() {
	dict, letterc, err := solve.LoadDictionary(conf.MustGetArgs())
	if err != nil {
		log.Fatalf("Could not load dictionary: %s", err)
	}

	guesser := solve.NewGuesser(dict, letterc)
	usr := bufio.NewScanner(os.Stdin)
	for {
		possibleMatches := guesser.PossibleMatches()
		fmt.Printf("%d possible remaining matches of %d letters.\n", possibleMatches, letterc)
		recomendation := guesser.Guess()
		fmt.Printf("Recommended guess is: %q\n\n", recomendation)

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
			if solve.ValidResult(guess) {
				fmt.Println("You entered a result, not word. Try again.")
				guess = ""
				continue
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
