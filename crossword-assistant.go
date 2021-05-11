//+build js, wasm

package main

import (
	"sort"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/Tom-Johnston/mamba/dawg"
)

type letterCount struct {
	letter byte
	count  int
}

//AnagramSearcher searches a dawg for words which have the same multiset of letters as the target. Any letters with the value blank are assumed to be wildcards and match any letter.
type PanagramSearcher struct {
	counts       []letterCount
	blanks       int
	blank        byte
	targetLength int
	currPath     []byte
}

//AllowStep checks if taking the step b is valid from the current position.
func (p PanagramSearcher) AllowStep(b byte) bool {
	if p.targetLength <= len(p.currPath) {
		return false
	}
	if p.blanks > 0 {
		return true
	}
	for i := range p.counts {
		if p.counts[i].letter == b && p.counts[i].count > 0 {
			return true
		}
	}
	return false
}

//Step takes the step b from the current position. This modifies the searcher.
func (p *PanagramSearcher) Step(b byte) {
	for i := range p.counts {
		if p.counts[i].letter == b && p.counts[i].count > 0 {
			p.counts[i].count--
			p.currPath = append(p.currPath, b)
			return
		}
	}
	p.blanks--
	p.currPath = append(p.currPath, p.blank)
}

//Backstep undoes the last step from the searcher.This modifies the searcher.
func (p *PanagramSearcher) Backstep() {
	pathElem := p.currPath[len(p.currPath)-1]
	p.currPath = p.currPath[:len(p.currPath)-1]
	if pathElem == p.blank {
		p.blanks++
		return
	}
	for i := range p.counts {
		if p.counts[i].letter == pathElem {
			p.counts[i].count++
			return
		}
	}
}

//AllowWord checks if the current position is allowed as a matching word.
func (p PanagramSearcher) AllowWord() bool {
	return true
}

//Chosen notifies the searcher that the current position has been chosen as a valid word.
func (p PanagramSearcher) Chosen() {}

//NewAnagramSearcher returns an anagram searcher to find anagrams of a word in a dawg with blanks.
func NewPanagramSearcher(anagram []byte, blank byte) *PanagramSearcher {
	tmp := make([]byte, len(anagram))
	copy(tmp, anagram)
	sort.Slice(tmp, func(i, j int) bool { return anagram[i] < anagram[j] })
	counts := make([]letterCount, 0)
	blanks := 0
	for i, l := range tmp {
		if l == blank {
			blanks++
			continue
		}
		if i > 1 && tmp[i-1] == tmp[i] {
			counts[len(counts)-1].count++
		} else {
			counts = append(counts, letterCount{letter: l, count: 1})
		}
	}

	currPath := make([]byte, 0, len(anagram))

	return &PanagramSearcher{counts: counts, blanks: blanks, blank: blank, targetLength: len(anagram), currPath: currPath}
}

//LimitSearcher limits the search to only find the first n solutions.
type LimitSearcher struct {
	seen  int
	limit int
}

//AllowStep checks if taking the step b is valid from the current position.
func (l LimitSearcher) AllowStep(b byte) bool {
	return l.seen < l.limit
}

//Step takes the step b from the current position. This modifies the searcher.
func (l *LimitSearcher) Step(b byte) {
}

//Backstep undoes the last step from the searcher.This modifies the searcher.
func (l *LimitSearcher) Backstep() {
}

//AllowWord checks if the current position is allowed as a matching word.
func (l LimitSearcher) AllowWord() bool {
	return l.seen < l.limit
}

//Chosen notifies the searcher that the current position has been chosen as a valid word.
func (l *LimitSearcher) Chosen() {
	l.seen++
}

type SubwordSearcher struct {
	word            []byte
	currentPostions []int
	blank           byte
}

func (s SubwordSearcher) AllowStep(b byte) bool {
	currentPosition := -1
	if len(s.currentPostions) > 0 {
		currentPosition = s.currentPostions[len(s.currentPostions)-1]
	}

	for i := currentPosition + 1; i < len(s.word); i++ {
		if s.word[i] == b || s.word[i] == s.blank {
			return true
		}
	}
	return false
}

func (s *SubwordSearcher) Step(b byte) {
	currentPosition := -1
	if len(s.currentPostions) > 0 {
		currentPosition = s.currentPostions[len(s.currentPostions)-1]
	}
	for i := currentPosition + 1; i < len(s.word); i++ {
		if s.word[i] == b || s.word[i] == s.blank {
			s.currentPostions = append(s.currentPostions, i)
			return
		}
	}
}

func (s *SubwordSearcher) Backstep() {
	s.currentPostions = s.currentPostions[:len(s.currentPostions)-1]
}

func (s SubwordSearcher) AllowWord() bool {
	return true
}

func (s SubwordSearcher) Chosen() {}

func NewSubwordSearcher(word []byte, blank byte) *SubwordSearcher {
	return &SubwordSearcher{word: word, currentPostions: []int{}, blank: blank}
}

//NewLimitSearcher returns a limit searcher for searching for the first n solutions.
func NewLimitSearcher(n int) *LimitSearcher {
	return &LimitSearcher{seen: 1, limit: n}
}

func main() {
	ds := make([]*dawg.Dawg, 0)
	for i := range dGobs {
		d := new(dawg.Dawg)
		if dGobs[i] == nil {
			ds = append(ds, nil)
			continue
		}

		err := d.GobDecode(dGobs[i])
		if err != nil {
			panic(err)
		}
		ds = append(ds, d)
	}

	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 {
			return "error - nothing to search for"
		}
		jsSearchers := args[0]
		length := jsSearchers.Length()
		if length == 0 {
			return "error - nothing to search for"
		}

		searchers := make([]dawg.Searcher, 0, length+1)
		for i := 0; i < length; i++ {
			jsSrch := jsSearchers.Index(i)
			jsType := jsSrch.Get("type")
			if jsType.Type() != js.TypeString {
				return "error - no type string found"
			}
			jsTypeString := jsType.String()
			jsContent := jsSrch.Get("content")
			if jsContent.Type() != js.TypeString {
				return "error - no content string found"
			}
			jsContentString := jsContent.String()
			if len(jsContentString) == 0 {
				continue
			}
			if jsTypeString == "anagram" {
				srch := NewPanagramSearcher([]byte(jsContentString), 63)
				searchers = append(searchers, srch)
				continue
			}
			if jsTypeString == "pattern" {
				srch := dawg.NewPatternSearcher([]byte(jsContentString), 63)
				searchers = append(searchers, srch)
				continue
			}
			return "error - unknown query type"
		}
		if len(searchers) == 0 {
			return "error - nothing to search for"
		}
		searchers = append(searchers, NewLimitSearcher(250))
		output := make(map[string]interface{})
		for i := len(ds) - 1; i >= 0; i-- {
			d := ds[i]
			if d == nil {
				continue
			}
			w, ids := d.Search(searchers...)
			r := make([]interface{}, 0)
			for j, k := range ids {
				pre := ""
				//Check if this id has extra data.
				index := sort.SearchInts(indices[i], k)
				if index < len(indices[i]) && indices[i][index] == k {
					pre = words[i][index]
				}

				splits := strings.Split(pre, ";")
				for _, s := range splits {
					if s == "" {
						r = append(r, string(w[j]))
					} else {
						r = append(r, s)
					}
				}
			}
			if len(r) > 0 {
				output[strconv.Itoa(i)] = r
			}
		}
		return output

	})
	js.Global().Set("search", cb)
	select {}
}
