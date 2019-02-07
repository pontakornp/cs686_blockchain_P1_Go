package p1

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	//"reflect"
)

type Flag_value struct {
	encoded_prefix []uint8
	value string
}

type Node struct {
	node_type int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value Flag_value
}

type MerklePatriciaTrie struct {
	db map[string]Node
	root string
}

func (mpt *MerklePatriciaTrie) Get(key string) string {
	// TODO
	return ""
}

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	// TODO
}

func (mpt *MerklePatriciaTrie) Delete(key string) {
	// TODO
}

func compact_encode(hex_array []uint8) []uint8 {
	term := false
	if hex_array[len(hex_array) - 1] == 16 {
		term = true
	}
	if term {
		hex_array = hex_array[:len(hex_array) - 1]
	}
	var odd_len int = len(hex_array) % 2
	var flags uint8 = uint8(2 * odd_len + 1)
	if(odd_len == 1) {
		hex_array  = append([]uint8{flags}, hex_array...)
	}
	o := []uint8{}
	for i := 0; i < len(hex_array); i+=2 {
		o = append(o, 16 * hex_array[i] + hex_array[i + 1])
	}
	return o
}

// If Leaf, ignore 16 at the end
func compact_decode(encoded_arr []uint8) []uint8 {
	// TODO
	hex_array := []uint8{}
	for i := 0; i < len(encoded_arr); i++ {
		hex_array = append(hex_array, encoded_arr[i]/16)
		hex_array = append(hex_array, encoded_arr[i]%16)
	}
	return []uint8{}
}

func Test_compact_encode() {
	//fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	//fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	//fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	//fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
	fmt.Println("", compact_encode([]uint8{2, 6, 3}))
}

func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value
	}

	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}