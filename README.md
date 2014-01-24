# go_trie

A trie is an ordered tree that is commonly used for searching a dictionary of terms.
Tries allow for fast prefix search and string matching and differ from binary trees
in that each decendent of a root node share their parents as a prefix.

## Memory

Simple tries, where each character is stored as a node, aren't particularly memory efficient
as Go structs take 64 bytes at a minimum. Representing the entire UTF-8 character space requires
using runes (equivalent to an int16).

Using a radix trie, where multiple characters are stored per node, provides for improved
memory usage. For example, the structure might look like:

```
- a
-- pp
--- le
--- lication
```

instead of:

```
- a
-- p
--- p
---- l
----- e
---- i
----- c
------ a
------- t
-------- i
--------- o
---------- n
```

## TODO

Radix trie implementation (WIP in radix_trie branch)
