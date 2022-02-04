package main

type optmap []map[byte]uint
type freqmap map[byte]uint
type scorefunc func(optmap, freqmap, int, byte) uint

func scoreBasic(opts optmap, _ freqmap, pos int, char byte) uint {
	return opts[pos][char]
}

func scoreWFreq(opts optmap, freq freqmap, pos int, char byte) uint {
	return opts[pos][char]*2 + freq[char]
}

func scoreWeightOptLeft(opts optmap, _ freqmap, pos int, char byte) uint {
	return opts[pos][char] * uint(len(opts[pos]))
}

func scoreWeightOptElim(opts optmap, _ freqmap, pos int, char byte) uint {
	return opts[pos][char] * uint(26-len(opts[pos]))
}

func scoreBasicSquared(opts optmap, freq freqmap, pos int, char byte) uint {
	score := scoreBasic(opts, freq, pos, char)
	return score * score
}
