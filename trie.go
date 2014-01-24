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
		parent.WordId = append(parent.WordId, id)
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
