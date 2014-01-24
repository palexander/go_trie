package main

import (
	"fmt"
	// "math"
	"bufio"
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
)

type TrieNode struct {
	// WordId   []int
	Children []*TrieNode
	Delim    bool
	Rune     rune
}

func (tn *TrieNode) AddChild(r rune) *TrieNode {
	node, exists := tn.FindChild(r)
	if !exists {
		node = &TrieNode{Rune: r}
		tn.Children = append(tn.Children, node)
	}
	return node
}

func (tn *TrieNode) FindChild(r rune) (*TrieNode, bool) {
	var foundNode *TrieNode
	found := false
	for _, node := range tn.Children {
		if node.Rune == r {
			found = true
			foundNode = node
			break
		}
	}
	return foundNode, found
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	node := TrieNode{}
	t := (Trie{&node})
	return &t
}

func NewTrieFile(path string) *Trie {
	trie := NewTrie()
	file, _ := os.Open(path)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id_word := strings.Split(scanner.Text(), "\t")
		id, _ := strconv.Atoi(id_word[0])
		trie.AddWord(id_word[1], id)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return trie
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

	// Delete first entry
	chars = append(chars[:0], chars[0+1:]...)

	trieNode, exists := parent.FindChild(currentRune)
	if !exists {
		trieNode = parent.AddChild(currentRune)
	}

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
		child, exists = child.FindChild(currentRune)
		if !exists {
			break
		}
	}
	return exists, child
}

func BytesToGB(bb float64) string {
	return fmt.Sprintf("%.3f GB", bb/(1024.0*1024.0*1024.0))
}

func Stop() {
	log.Println("stopped - enter anything to exit")
	in := bufio.NewReader(os.Stdin)
	in.ReadString('\n')
}

func GoRuntimeStats() {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	log.Println("Memory Acquired: ", BytesToGB(float64(m.Sys)))
	log.Println("Memory Used    : ", BytesToGB(float64(m.Alloc)))
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
	trie.AddWord("ab", 0)
	trie.AddWord("ab", 0)
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
		// fmt.Println(count)
		if count == 50000000 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// time.Sleep(time.Minute * 5)
	time.Sleep(1)

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
	fmt.Println(trie.IsPrefix("ab"))
	fmt.Println(trie.IsPrefix("the"))
	fmt.Println(trie.IsWord("the best ever"))
	fmt.Println(trie.IsWord("terror"))
	fmt.Println(trie.IsWord("annie"))
	fmt.Println(trie.IsWord("annie oakley"))
	fmt.Println(trie.IsWord("hello π"))
	fmt.Println(trie.IsWord("Љ"))

	fmt.Println("THE LIES")
	fmt.Println(trie.IsWord("the"))
	fmt.Println(trie.IsPrefix("blah"))
	fmt.Println(trie.IsWord("Ђ"))

	GoRuntimeStats()
	runtime.GC()
	GoRuntimeStats()
	Stop()

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
