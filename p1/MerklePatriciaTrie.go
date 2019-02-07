package p1

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"log"
	"reflect"
)

type Flag_value struct {
	encoded_prefix []uint8
	value string // hash node or value of leaf
}

type Node struct {
	node_type int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value Flag_value
}

func isEmpty(node Node) bool {
	return reflect.DeepEqual(node, nil)
}

type MerklePatriciaTrie struct {
	db map[string]Node
	root string
}

func (mpt *MerklePatriciaTrie) Get(key string) string {
	// TODO
	hex_key, err := hex.DecodeString(key)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	if(len(mpt.db) == 0 || mpt.root == "") {
		return ""
	}
	hash_node := mpt.root
	for hash_node != "" {
		node := mpt.db[hash_node]
		if isEmpty(node) {
			return ""
		}
		node_type := node.node_type
		if node_type == 0 { // null node
			return ""
		} else if node_type == 1 { // branch node
			if len(hex_key) == 1 && node.branch_value[len(node.branch_value) - 1] != ""{
				return node.branch_value[len(node.branch_value) - 1]
			}
			hash_node = node.branch_value[hex_key[0]]
			if hash_node == "" {
				return ""
			}
			hex_key = hex_key[1:]
		} else { // node_type == 2, ext or leaf node
			encoded_arr := node.flag_value.encoded_prefix
			decoded_arr := compact_decode(encoded_arr)
			boo := isLeafNode(encoded_arr)
			if boo { // leaf node
				if len(hex_key) == len(decoded_arr) {
					for i := 0; i < len(decoded_arr); i++ {
						if hex_key[i] != decoded_arr[i] {
							return ""
						}
					}
					return node.flag_value.value
				}
			} else { // extension node
				if len(hex_key) < len(decoded_arr) {
					return ""
				}
				for i := 0; i < len(decoded_arr); i++ {
					if hex_key[i] != decoded_arr[i] {
						return ""
					}
				}
				hex_key = hex_key[len(hex_key) - len(decoded_arr):]
				hash_node = node.flag_value.value
			}
		}
	}
	return ""
}

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	// TODO
}

func (mpt *MerklePatriciaTrie) Delete(key string) {
	// TODO
}

func compact_encode(hex_array []uint8) []uint8 {
	term := 0
	if hex_array[len(hex_array) - 1] == 16 {
		term = 1
	}
	if term == 1 {
		hex_array = hex_array[:len(hex_array) - 1]
	}
	var odd_len int = len(hex_array) % 2
	var flags uint8 = uint8(2 * term + odd_len)
	if odd_len == 1 {
		hex_array  = append([]uint8{flags}, hex_array...)
	} else {
		hex_array = append([]uint8{flags, 0}, hex_array...)
	}
	o := []uint8{}
	for i := 0; i < len(hex_array); i+=2 {
		o = append(o, 16 * hex_array[i] + hex_array[i + 1])
	}
	return o
}

// If Leaf, ignore 16 at the end
func compact_decode(encoded_arr []uint8) []uint8 {
	hex_array := []uint8{}
	for i := 0; i < len(encoded_arr); i++ {
		hex_array = append(hex_array, encoded_arr[i] / 16)
		hex_array = append(hex_array, encoded_arr[i] % 16)
	}
	if hex_array[0] == 0 || hex_array[0] == 2 {
		hex_array = hex_array[2:]
	} else {
		hex_array = hex_array[1:]
	}
	return hex_array
}

func isLeafNode(encoded_arr []uint8) bool {
	prefix := encoded_arr[0] / 16
	if prefix == 0 || prefix == 1 {
		return false
	}
	return true
}

func Test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
	fmt.Println("", compact_encode([]uint8{2, 6, 3, 16}))
	//fmt.Println("", compact_decode(compact_encode([]uint8{2, 6, 3, 16})))
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