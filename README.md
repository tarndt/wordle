# wordle
Simple wordle solver

Lets see it in action:
```
$ ./wordle 
5914 possible remaining matches of 5 letters.
Recomended guess is: "sales"

Please enter actual guess (5 letters) or enter for "sales": 
Please enter result of guess (5 y|n|p): pnnpn
210 possible remaining matches of 5 letters.
Recomended guess is: "beets"

Please enter actual guess (5 letters) or enter for "beets": 
Please enter result of guess (5 y|n|p): npnpp
8 possible remaining matches of 5 letters.
Recomended guess is: "crest"

Please enter actual guess (5 letters) or enter for "crest": 
Please enter result of guess (5 y|n|p): nnpyp
1 possible remaining matches of 5 letters.
Recomended guess is: "those"

Our work here is done!
```

Usage:
```
$ ./wordle -help
  -dict-file string
    	Path of dictionary file (default "/usr/share/dict/words")
  -help
    	Show usage and exit
  -letters uint
    	Number of letters in word (default 5)
  -no-punc
    	Should words with puncuation be excluded (default true)
```