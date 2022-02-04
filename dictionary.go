package main

import (
	"bufio"
	"fmt"
	"os"
)

type Dictionary map[string]struct{}

func loadDictionary(dictPath string, letterc uint) (Dictionary, uint, error) {
	fin, err := os.Open(dictPath)
	if err != nil {
		return nil, 0, fmt.Errorf("Could not open dictionary file %q: %w", dictPath, err)
	}
	defer fin.Close()

	keep := func(word []byte) bool {
		for _, char := range word {
			if char < 'a' || char > 'z' {
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

		if keep(word) {
			dict[string(word)] = struct{}{}
		}
	}
	if err = scan.Err(); err != nil {
		return nil, 0, fmt.Errorf("Could not read dictionary file %q: %w", dictPath, err)
	}

	return dict, letterc, nil
}

func (dict Dictionary) Clone() Dictionary {
	clone := make(Dictionary, len(dict))
	for key := range dict {
		clone[key] = struct{}{}
	}
	return clone
}
