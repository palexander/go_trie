package go_trie

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TrieNode struct {
	WordId   []int
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
		parent.WordId = append(parent.WordId, id)
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
	//  currentRune, _, _, _ := strconv.UnquoteChar(char, 0)
	//  child, exists = child.FindChild(currentRune)
	//  if !exists {
	//    break
	//  }
	// }
	// return exists, child

	return false, &TrieNode{}
}
