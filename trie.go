package main

import (
	"fmt"
	// "math"
	"bufio"
	"log"
	"os"
	"runtime"
	// "runtime/pprof"
	"github.com/davecheney/profile"
	"strconv"
	"strings"
	"time"
)

type TrieNode struct {
	// WordId   []int
	Children []*TrieNode
	Delim    bool
	Rune     []rune
}

func (tn *TrieNode) AddChild(r []rune) *TrieNode {
	node := &TrieNode{Rune: r}
	tn.Children = append(tn.Children, node)
	return node
}

func (tn *TrieNode) FindChild(word []rune) (*TrieNode, bool) {
	var foundNode *TrieNode
	found := false
	child, exists, endPos := tn.FindChildPrefix(word)
	if exists && (endPos+1 == len(word)) {
		foundNode = child
		found = true
	}
	return foundNode, found
}

func (tn *TrieNode) FindChildPrefix(word []rune) (*TrieNode, bool, int) {
	if len(word) == 0 {
		// fmt.Println("No length word in find child prefix")
		return nil, false, 0
	}

	var foundNode *TrieNode
	found := false

	// Find child with same first letter
	for _, node := range tn.Children {
		if len(node.Rune) > 0 && node.Rune[0] == word[0] {
			foundNode = node
			found = true
		}
	}

	lastIndex := -1
	if found {
		for index, char := range word {
			if index >= len(foundNode.Rune) || char != foundNode.Rune[index] {
				break
			}
			lastIndex += 1
		}
	}

	return foundNode, found, lastIndex
}

func (tn *TrieNode) FindDeepestPrefix(word []rune) (*TrieNode, []rune, bool, int) {
	foundNode, found, endPos := tn.FindChildPrefix(word)

	newFoundNode := foundNode
	newFound := found
	newEndPos := endPos
	newWord := word
	for len(newWord) > 0 && newFound && len(foundNode.Children) > 0 {
		newWord = newWord[newEndPos+1:]
		newFoundNode, newFound, newEndPos = foundNode.FindChildPrefix(newWord)
		if newFound {
			foundNode = newFoundNode
			endPos = newEndPos
			word = newWord
		}
	}

	return foundNode, word, found, endPos
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	node := TrieNode{}
	t := (Trie{&node})
	return &t
}

func (t Trie) AddWord(word string, id int) {
	word = strings.ToLower(word)
	// splitWord := strings.Split(word, "")
	// t.BuildTree(splitWord, id, t.root)
	t.AddRadixWord([]rune(word), id, t.root)
}

func Overlap(word1 []rune, word2 []rune) int {
	lastIndex := -1
	for index, char := range word1 {
		if index >= len(word2) || char != word2[index] {
			break
		}
		lastIndex += 1
	}
	return lastIndex
}

func (t Trie) AddRadixWord(word []rune, id int, parent *TrieNode) {
	// We hit the endo the line, make the parent a word delim
	if len(word) == 0 {
		parent.Delim = true
		return
	}

	// Easy case
	if len(parent.Children) == 0 {
		node := parent.AddChild(word)
		node.Delim = true
		return
	}

	child, exists := parent.FindChild(word)
	// child1, lastWordFragment1, _, endPos1 := parent.FindDeepestPrefix(word)
	// fmt.Println("\n\nWord:" + string(word))
	// fmt.Println(exists)
	// fmt.Println(child)
	// fmt.Println(child1)
	// fmt.Println(string(lastWordFragment1))
	// fmt.Println(endPos1)
	if exists {
		child.Delim = true
	} else {
		child, lastWordFragment, exists, _ := parent.FindDeepestPrefix(word)
		if exists {
			overlap := Overlap(lastWordFragment, child.Rune)
			newChildRune := lastWordFragment[overlap+1:]
			movedChildRune := child.Rune[overlap+1:]
			// fmt.Println("Moving runes")
			// fmt.Println("Overlap:")
			// fmt.Println(overlap)
			// fmt.Println("word:" + string(lastWordFragment))
			// fmt.Println("child.Rune:" + string(child.Rune))
			child.Rune = child.Rune[:overlap+1]
			// fmt.Println("altered child.Rune:" + string(child.Rune))
			_, newExists := child.FindChild(newChildRune)
			// fmt.Println(newExists)
			if len(newChildRune) > 0 && !newExists {
				// fmt.Println("newChildRune:" + string(newChildRune))
				newChild := child.AddChild(newChildRune)
				newChild.Delim = true
			}
			_, movedExists := child.FindChild(movedChildRune)
			if len(movedChildRune) > 0 && !movedExists {
				// fmt.Println("movedChildRune:" + string(movedChildRune))
				movedChild := child.AddChild(movedChildRune)
				movedChild.Delim = child.Delim
				child.Delim = false
			}
		} else {
			node := parent.AddChild(word)
			node.Delim = true
		}
	}
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
	// prefix = strings.ToLower(prefix)
	// splitPrefix := strings.Split(prefix, "")
	// var child *TrieNode
	// var exists bool
	// child = t.root
	// for _, char := range splitPrefix {
	// 	currentRune, _, _, _ := strconv.UnquoteChar(char, 0)
	// 	child, exists = child.FindChild(currentRune)
	// 	if !exists {
	// 		break
	// 	}
	// }
	// return exists, child

	return false, &TrieNode{}
}

func main() {
	defer profile.Start(profile.CPUProfile).Stop()

	trie := NewTrie()
	count := 0
	for count < 1 {
		trie.AddWord("apple", 0)
		trie.AddWord("at", 0)
		trie.AddWord("apple", 0)
		trie.AddWord("art", 0)
		trie.AddWord("application", 0)
		trie.AddWord("abacus", 0)
		trie.AddWord("algebra", 0)
		trie.AddWord("baby", 0)
		trie.AddWord("broken", 0)
		trie.AddWord("belly", 0)
		trie.AddWord("the best ever", 0)
		trie.AddWord("terror", 1)
		trie.AddWord("terror", 7)
		trie.AddWord("terrorist", 2)
		trie.AddWord("annie", 3)
		trie.AddWord("annie oakley", 4)
		trie.AddWord("hello π", 5)
		trie.AddWord("Љ", 6)
		count += 1
		if count%1000 == 0 {
			fmt.Println(count)
		}
	}

	file, _ := os.Open("/Users/palexand/tmp/dictionary.txt")
	scanner := bufio.NewScanner(file)
	count = 0
	for scanner.Scan() {
		id_word := strings.Split(scanner.Text(), "\t")
		id, _ := strconv.Atoi(id_word[0])
		trie.AddWord(id_word[1], id)
		count += 1
		if count%10000 == 0 {
			fmt.Println(count)
		}
		if count == 10000000 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// time.Sleep(time.Minute * 5)
	time.Sleep(1)

	fmt.Println(len(trie.root.Children))
	for _, node := range trie.root.Children {
		fmt.Println("\nfirst")
		fmt.Println(node)
		fmt.Println(string(node.Rune))
		for _, node1 := range node.Children {
			fmt.Println("second")
			fmt.Println(node1)
			fmt.Println(string(node1.Rune))
			for _, node2 := range node1.Children {
				fmt.Println("third")
				fmt.Println(node2)
				fmt.Println(string(node2.Rune))
			}
		}
	}

	// fmt.Println("THE TRUTH")
	// fmt.Println(trie.IsPrefix("ab"))
	// fmt.Println(trie.IsPrefix("the"))
	// fmt.Println(trie.IsWord("the best ever"))
	// fmt.Println(trie.IsWord("terror"))
	// fmt.Println(trie.IsWord("annie"))
	// fmt.Println(trie.IsWord("annie oakley"))
	// fmt.Println(trie.IsWord("hello π"))
	// fmt.Println(trie.IsWord("Љ"))

	// fmt.Println("THE LIES")
	// fmt.Println(trie.IsWord("the"))
	// fmt.Println(trie.IsPrefix("blah"))
	// fmt.Println(trie.IsWord("Ђ"))

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
