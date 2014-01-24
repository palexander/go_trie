package go_trie

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
)

const (
	DEBUG = false
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func TestSimpleTrie(t *testing.T) {
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
	trie.AddWord("ab", 1)
	trie.AddWord("the best ever", 2)
	trie.AddWord("terror", 3)
	trie.AddWord("terror", 4)
	trie.AddWord("terrorist", 5)
	trie.AddWord("annie", 6)
	trie.AddWord("annie oakley", 7)
	trie.AddWord("hello π", 8)
	trie.AddWord("Љ", 9)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			fmt.Println(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	var test bool

	test, _ = trie.IsPrefix("ab")
	if test == false {
		t.Error("Missing prefix")
	}
	test, _ = trie.IsPrefix("the")
	if test == false {
		t.Error("Missing word")
	}
	test, _ = trie.IsWord("the best ever")
	if test == false {
		t.Error("Missing word")
	}
	test, _ = trie.IsWord("terror")
	if test == false {
		t.Error("Missing word")
	}
	test, _ = trie.IsWord("annie")
	if test == false {
		t.Error("Missing word")
	}
	test, _ = trie.IsWord("annie oakley")
	if test == false {
		t.Error("Missing word")
	}
	test, _ = trie.IsWord("hello π")
	if test == false {
		t.Error("Missing word")
	}
	test, _ = trie.IsWord("Љ")
	if test == false {
		t.Error("Missing word")
	}

	test, _ = trie.IsPrefix("blah")
	if test == true {
		t.Error("Prefix exists that shouldn't")
	}
	test, _ = trie.IsWord("the")
	if test == true {
		t.Error("Word exists that shouldn't")
	}
	test, _ = trie.IsWord("Ђ")
	if test == true {
		t.Error("Word exists that shouldn't")
	}

	if DEBUG {
		GoRuntimeStats()
		runtime.GC()
		GoRuntimeStats()
		Stop()
	}
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
