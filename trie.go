package main

import (
	"fmt"
	// "math"
	"bufio"
	"flag"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
)

type Trie struct {
	root *TrieNode
}

type TrieHolder map[rune]*TrieNode

type TrieNode struct {
	WordId   []int
	Children TrieHolder
	Parent   *TrieNode
	Delim    bool
}

func NewTrie() *Trie {
	node := TrieNode{Children: make(TrieHolder)}
	t := (Trie{&node})
	return &t
}

func (t Trie) AddWord(word string, id int) {
	word = strings.ToLower(word)
	splitWord := strings.Split(word, "")
	t.BuildTree(splitWord, id, t.root)
}

func (t Trie) BuildTree(chars []string, id int, parent *TrieNode) {
	if len(chars) == 0 {
		parent.Delim = true
		// parent.WordId = append(parent.WordId, id)
		return
	}

	// Get char as rune
	currentRune, _, _, _ := strconv.UnquoteChar(chars[0], 0)
	// fmt.Println(currentRune)

	// Delete first entry
	chars = append(chars[:0], chars[0+1:]...)

	trieNode, exists := parent.Children[currentRune]
	if !exists {
		children := make(TrieHolder)
		trieNode = &TrieNode{Parent: parent, Children: children}
	}

	parent.Children[currentRune] = trieNode
	t.BuildTree(chars, id, trieNode)
}

func (t Trie) RemoveWord() {

}

func (t Trie) IsWord(word string) (bool, *TrieNode) {
	exists, node := t.IsPrefix(word)
	delim := false
	if exists {
		delim = (node.Delim == true)
	}
	return delim, node
}

func (t Trie) IsPrefix(prefix string) (bool, *TrieNode) {
	prefix = strings.ToLower(prefix)
	splitPrefix := strings.Split(prefix, "")
	var child *TrieNode
	var exists bool
	child = t.root
	for _, char := range splitPrefix {
		currentRune, _, _, _ := strconv.UnquoteChar(char, 0)
		child, exists = child.Children[currentRune]
		if !exists {
			break
		}
	}
	return exists, child
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Println(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	trie := NewTrie()
	trie.AddWord("the best ever", 0)
	trie.AddWord("terror", 1)
	trie.AddWord("terror", 7)
	trie.AddWord("terrorist", 2)
	trie.AddWord("annie", 3)
	trie.AddWord("annie oakley", 4)
	trie.AddWord("hello π", 5)
	trie.AddWord("Љ", 6)

	file, _ := os.Open("/Users/palexand/tmp/dictionary.txt")
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		id_word := strings.Split(scanner.Text(), "\t")
		id, _ := strconv.Atoi(id_word[0])
		trie.AddWord(id_word[1], id)
		count += 1
		fmt.Println(count)
		if count == 500000 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			fmt.Println(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	fmt.Println("THE TRUTH")
	fmt.Println(trie.IsPrefix("the"))
	fmt.Println(trie.IsWord("the best ever"))
	fmt.Println(trie.IsWord("terror"))
	fmt.Println(trie.IsWord("annie"))
	fmt.Println(trie.IsWord("annie oakley"))
	fmt.Println(trie.IsWord("hello π"))
	fmt.Println(trie.IsWord("Љ"))
	fmt.Println(trie.IsWord("Hygiène_de_vie"))

	fmt.Println("THE LIES")
	fmt.Println(trie.IsWord("the"))
	fmt.Println(trie.IsPrefix("blah"))
	fmt.Println(trie.IsWord("Ђ"))

	// This results in almost 5GB memory
	// var holder [][math.MaxInt8]TrieNode
	// holder = make([][math.MaxInt8]TrieNode, 1000000)
	// for i := 0; i < 1000000; i++ {
	// 	var newArray [math.MaxInt8]TrieNode
	// 	holder[i] = newArray
	// }
	// fmt.Println("TrieNode Arrays created")
	// for true {
	// 	fmt.Println("waiting")
	// }
}
