package main

import (
	"bytes"
	"encoding/binary"
)

const HEADER = 4

// | type | nkeys |  pointers  |   offsets  | key-values | unused |
// |  2B  |   2B  | nkeys * 8B | nkeys * 2B |     ...    |        |

// | klen | vlen | key | val |
// |  2B  |  2B  | ... | ... |

const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

type BNode []byte

type BTree struct {
	root uint64
	get func(uint64) []byte
	new func([]byte) uint64
	del func(uint64)
}

const (
	BNODE_NODE = 1
	BNODE_LEAF = 2
)

func (node BNode) btype() uint16 {
	// LittleEndian is used because it's default endian in x86, x86-64 architecture, which are the most popular in modern computers
	return binary.LittleEndian.Uint16(node[0:2])
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

func (node BNode) setHeader(btype, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

func (node BNode) getPtr(idx uint16) uint64 {
	// assert(idx < node.nkeys())
	pos := HEADER + 8 * idx
	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) setPtr(idx uint16, val uint64)

func offsetPos(node BNode, idx uint16) uint16 {
	return HEADER + 8 * node.nkeys() + 2*(idx - 1)
}

func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node[offsetPos(node, idx):])
}

func (node BNode) kvPos(idx uint16) uint16 {
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16)  []byte {
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:klen]
}

func (node BNode) getVal(idx uint16) []byte

func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}

func nodeLookupLE(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	found := uint16(0)
	for i := uint16(1); i < nkeys; i++ {
		comp := bytes.Compare(node.getKey(i), key)
		if comp <= 0 {
			found = i
		}
		if comp >= 0 {
			break
		}
	}
	return found
}

