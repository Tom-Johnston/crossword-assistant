package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Tom-Johnston/mamba/dawg"
)

var toReplace = []rune{'À', 'Á', 'Â', 'Ã', 'Ä', 'Å', 'Ç', 'È', 'É', 'Ê', 'Ë', 'Ì', 'Í', 'Î', 'Ï', 'Ð', 'Ñ', 'Ò', 'Ó', 'Ô', 'Õ', 'Ö', 'Ø', 'Ù', 'Ú', 'Û', 'Ü', 'Ý', 'ß', 'à', 'á', 'â', 'ã', 'ä', 'å', 'ç', 'è', 'é', 'ê', 'ë', 'ì', 'í', 'î', 'ï', 'ñ', 'ò', 'ó', 'ô', 'õ', 'ö', 'ø', 'ù', 'ú', 'û', 'ü', 'ý', 'ÿ', 'Ā', 'ā', 'Ă', 'ă', 'Ą', 'ą', 'Ć', 'ć', 'Ĉ', 'ĉ', 'Ċ', 'ċ', 'Č', 'č', 'Ď', 'ď', 'Đ', 'đ', 'Ē', 'ē', 'Ĕ', 'ĕ', 'Ė', 'ė', 'Ę', 'ę', 'Ě', 'ě', 'Ĝ', 'ĝ', 'Ğ', 'ğ', 'Ġ', 'ġ', 'Ģ', 'ģ', 'Ĥ', 'ĥ', 'Ħ', 'ħ', 'Ĩ', 'ĩ', 'Ī', 'ī', 'Ĭ', 'ĭ', 'Į', 'į', 'İ', 'ı', 'Ĵ', 'ĵ', 'Ķ', 'ķ', 'Ĺ', 'ĺ', 'Ļ', 'ļ', 'Ľ', 'ľ', 'Ŀ', 'ŀ', 'Ł', 'ł', 'Ń', 'ń', 'Ņ', 'ņ', 'Ň', 'ň', 'ŉ', 'Ō', 'ō', 'Ŏ', 'ŏ', 'Ő', 'ő', 'Ŕ', 'ŕ', 'Ŗ', 'ŗ', 'Ř', 'ř', 'Ś', 'ś', 'Ŝ', 'ŝ', 'Ş', 'ş', 'Š', 'š', 'Ţ', 'ţ', 'Ť', 'ť', 'Ŧ', 'ŧ', 'Ũ', 'ũ', 'Ū', 'ū', 'Ŭ', 'ŭ', 'Ů', 'ů', 'Ű', 'ű', 'Ų', 'ų', 'Ŵ', 'ŵ', 'Ŷ', 'ŷ', 'Ÿ', 'Ź', 'ź', 'Ż', 'ż', 'Ž', 'ž', 'ſ', 'ƒ', 'Ơ', 'ơ', 'Ư', 'ư', 'Ǎ', 'ǎ', 'Ǐ', 'ǐ', 'Ǒ', 'ǒ', 'Ǔ', 'ǔ', 'Ǖ', 'ǖ', 'Ǘ', 'ǘ', 'Ǚ', 'ǚ', 'Ǜ', 'ǜ', 'Ǻ', 'ǻ', 'Ǿ', 'ǿ'}
var replaceWith = []byte{'A', 'A', 'A', 'A', 'A', 'A', 'C', 'E', 'E', 'E', 'E', 'I', 'I', 'I', 'I', 'D', 'N', 'O', 'O', 'O', 'O', 'O', 'O', 'U', 'U', 'U', 'U', 'Y', 'S', 'A', 'A', 'A', 'A', 'A', 'A', 'C', 'E', 'E', 'E', 'E', 'I', 'I', 'I', 'I', 'N', 'O', 'O', 'O', 'O', 'O', 'O', 'U', 'U', 'U', 'U', 'Y', 'Y', 'A', 'A', 'A', 'A', 'A', 'A', 'C', 'C', 'C', 'C', 'C', 'C', 'C', 'C', 'D', 'D', 'D', 'D', 'E', 'E', 'E', 'E', 'E', 'E', 'E', 'E', 'E', 'E', 'G', 'G', 'G', 'G', 'G', 'G', 'G', 'G', 'H', 'H', 'H', 'H', 'I', 'I', 'I', 'I', 'I', 'I', 'I', 'I', 'I', 'I', 'J', 'J', 'K', 'K', 'L', 'L', 'L', 'L', 'L', 'L', 'L', 'L', 'L', 'L', 'N', 'N', 'N', 'N', 'N', 'N', 'N', 'O', 'O', 'O', 'O', 'O', 'O', 'R', 'R', 'R', 'R', 'R', 'R', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'T', 'T', 'T', 'T', 'T', 'T', 'U', 'U', 'U', 'U', 'U', 'U', 'U', 'U', 'U', 'U', 'U', 'U', 'W', 'W', 'Y', 'Y', 'Y', 'Z', 'Z', 'Z', 'Z', 'Z', 'Z', 'S', 'F', 'O', 'O', 'U', 'U', 'A', 'A', 'I', 'I', 'O', 'O', 'U', 'U', 'U', 'U', 'U', 'U', 'U', 'U', 'U', 'U', 'A', 'A', 'O', 'O'}

type word struct {
	presentations []string
	normalised    []byte
}

func main() {
	start := time.Now()
	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
	scanner.Split(bufio.ScanLines)
	counter := 0
	wordList := []word{}

	longestWord := 0
	for scanner.Scan() {
		s := scanner.Text()
		normalised, err := normaliseASCII(s)
		if err != nil {
			fmt.Fprintln(os.Stderr, s)
			continue
		}

		if len(normalised) > longestWord {
			longestWord = len(normalised)
		}

		presentation := strings.ToUpper(s)
		if presentation == string(normalised) {
			presentation = ""
		}
		wordList = append(wordList, word{presentations: []string{presentation}, normalised: normalised})

		counter++
	}
	fmt.Fprintf(os.Stderr, "Normalised %v words in %v.\n", counter, time.Since(start))
	step := time.Now()
	less := func(i, j int) bool {
		return bytes.Compare(wordList[i].normalised, wordList[j].normalised) < 0
	}
	sort.SliceStable(wordList, less)
	fmt.Fprintf(os.Stderr, "Sorted in %v.\n", time.Since(step))
	step = time.Now()
	//Count the number of repeats.
	numberOfRepeats := 0
	index := 0
listLoop:
	for i := 1; i < len(wordList); i++ {
		if bytes.Compare(wordList[index].normalised, wordList[i].normalised) == 0 {
			for _, v := range wordList[index].presentations {
				if v == wordList[i].presentations[0] {
					numberOfRepeats++
					continue listLoop
				}
			}
			wordList[index].presentations = append(wordList[index].presentations, wordList[i].presentations[0])
			numberOfRepeats++
			continue
		}
		index++
		wordList[index] = wordList[i]
	}
	wordList = wordList[:index+1]
	fmt.Fprintf(os.Stderr, "Collapsed %v repeats in %v.\n", numberOfRepeats, time.Since(step))
	
	words := make([][][]byte, longestWord+1)
	presentations := make([][]string, longestWord+1)
	indices := make([][]int, longestWord+1)
	for i := range indices {
		indices[i] = []int{}
	}

	for _, v := range wordList {
		n := len(v.normalised)
		words[n] = append(words[n], v.normalised)
		pre := v.presentations[0]
		for i := 1; i < len(v.presentations); i++ {
			pre += ";"
			pre += v.presentations[i]
		}

		if pre != "" {
			presentations[n] = append(presentations[n], pre)
			indices[n] = append(indices[n], len(words[n])-1)
		}
	}

	file, err := os.Create("words.go")
	if err != nil{
		panic(err)
	}

	_, err = fmt.Fprintln(file, "package main")
	if err != nil{
		panic(err)
	}

	dGobs := make([][]byte, longestWord+1)

	for i, w := range words {
		step = time.Now()
		if len(w) == 0 {
			continue
		}
		d, err := dawg.New(w)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(os.Stderr, "Length of each word:", i, "Number of words: ", len(w), "Took:", time.Since(step))
		dGob, err := d.GobEncode()
		if err != nil {
			panic(err)
		}
		
		dGobs[i] = dGob
	}

	_, err = fmt.Fprintf(file, "var dGobs = %#v \n", dGobs)
	if err != nil{
		panic(err)
	}
	_, err = fmt.Fprintf(file, "var words = %#v\n", presentations)
	if err != nil{
		panic(err)
	}
	_, err = fmt.Fprintf(file, "var indices = %#v\n", indices)
	if err != nil{
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "Took %v.\n", time.Since(start))
}

func normaliseASCII(str string) ([]byte, error) {
	s := []rune(str)
	b := make([]byte, 0, len(s))
runeLoop:
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'a' <= c && c <= 'z' {
			c -= 'a' - 'A'
			b = append(b, byte(c))
		} else if 'A' <= c && c <= 'Z' {
			b = append(b, byte(c))
		} else if c == 39 {
			continue
		} else {
			for j := range toReplace {
				if c == toReplace[j] {
					b = append(b, replaceWith[j])
					continue runeLoop
				}
			}
			return nil, errors.New("Help")
		}
	}
	return b, nil
}
