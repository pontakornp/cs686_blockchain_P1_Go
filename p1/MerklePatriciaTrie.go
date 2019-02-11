package p1

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
)

type Flag_value struct {
	encoded_prefix []uint8 // shared nibble(s) for ext node or key for leaf node
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

// check if encoded array is a leaf
func isLeafNode(encoded_arr []uint8) bool {
	prefix := encoded_arr[0] / 16
	if prefix == 0 || prefix == 1 {
		return false
	}
	return true
}

func ConvertStringToHexArray(str string) []uint8 {
	hex_array := []uint8{}
	for i := 0; i < len(str); i++ {
		hex_array = append(hex_array, str[i]/16)
		hex_array = append(hex_array, str[i]%16)
	}
	return hex_array
}

func (mpt *MerklePatriciaTrie) Get(key string) string {
	//hex_key, err := hex.DecodeString(key)

	hex_key := ConvertStringToHexArray(key)

	if(len(mpt.db) == 0 || mpt.root == "") {
		return ""
	}
	hash_node := mpt.root
	fmt.Println(hash_node)
	for hash_node != "" {
		node := mpt.db[hash_node]
		fmt.Println(node)
		node_type := node.node_type
		if isEmpty(node) || node_type == 0 { // null node
			fmt.Println("node type is 0 or empty")
			return ""
		} else if node_type == 1 { // branch node
			fmt.Println("branch")
			// if hex_key is empty string check if value exists
			// if yes, return value,
			// if not, return empty string
			if hex_key == nil || len(hex_key) == 0 {
				tempValue := node.branch_value[len(node.branch_value) - 1]
				if tempValue != "" {
					return tempValue
				} else {
					fmt.Println("node is empty")
					return ""
				}
			}
			// update hash_node
			tempValue := node.branch_value[hex_key[0]]
			if tempValue != "" {
				hash_node = tempValue
			} else {
				fmt.Println("tempVal is empty")
				return ""
			}
			// if hex_key has one character left, update hex_key to nil
			// if hex_key has more than one character, remove the first ele
			if len(hex_key) == 1 {
				hex_key = nil
			} else {
				hex_key = hex_key[1:]
			}

		} else { // node_type == 2, ext or leaf node
			fmt.Println("leaf or ext")
			encoded_arr := node.flag_value.encoded_prefix // encoded_prefix
			decoded_arr := compact_decode(encoded_arr) // decode ascii prefix to hex string
			boo := isLeafNode(encoded_arr) // if leaf, true else false
			if boo { // leaf node
				if len(hex_key) == len(decoded_arr) {
					for i := 0; i < len(decoded_arr); i++ {
						if hex_key[i] != decoded_arr[i] {
							return ""
						}
					}
					return node.flag_value.value
				}
				return ""
			} else { // extension node
				// if hex_key length is less than key of the node, return empty string
				if len(hex_key) < len(decoded_arr) {
					return ""
				}
				// if hex_key length is equal to or more than the key of the node
				// loop through each character of node key
				// if any character does not match, return empty string
				for i := 0; i < len(decoded_arr); i++ {
					if hex_key[i] != decoded_arr[i] {
						return ""
					}
				}
				//if the remaining key length is equal to zero, then set hex_key to nil
				remaining_len := len(hex_key) - len(decoded_arr)
				if remaining_len == 0 {
					hex_key = nil
				} else { // if the remaining key length is more than zero
					hex_key = hex_key[len(decoded_arr):]
				}
				hash_node = node.flag_value.value
				// if value of the next hash node is empty then return empty string
				if hash_node == "" {
					return ""
				}
			}
		}
	}
	fmt.Println("return node empty")
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	// TODO
	// if root is empty, create a leaf and insert
	// if root is not empty, perform an operation according to each node type
	//node_stack := stack.New()
	//temp_val_stack := stack.New()
	if key == "" {
		return
	}
	path_arr := ConvertStringToHexArray(key)
	// if root is empty, create a leaf and insert
	// if root is not empty, perform an operation according to each node type
	hash_node := mpt.root

	// case when root is empty
	if(hash_node == "") {
		fmt.Println("root is empty")
		leaf_node := newLeafNode(path_arr, new_value)
		hash_leaf_node := leaf_node.hash_node()
		mpt.db[hash_leaf_node] = leaf_node
		mpt.root = hash_leaf_node
		return
	}
	fmt.Println("root is not empty")
	node := mpt.db[hash_node]
	for !isEmpty(node) {
		node_type := node.node_type
		//for()
		if node_type == 0 {

		} else if node_type == 1 {
			//if remain_prefix == "" {
			//
			//} else if len(remain_prefix) <= 2 { // case where the next node is a leaf node
			//	index := remain_prefix[0]
			//	new_node_type := 2
			//	value := new_value
			//	nibble := []uint8{}
			//	if len(remain_prefix) == 2 {
			//		nibble = append(nibble, remain_prefix[1])
			//	}
			//	//new leaf node
			//	new_node := newNode(new_node_type, value, nibble)
			//	node.branch_value[index] = new_value
			//	node_stack.Push(new_node)
			//} else { // case where the next node is not leaf node
			//
			//}
		} else { // node_type == 2
			encoded_prefix := node.flag_value.encoded_prefix
			nibble_arr := compact_decode(encoded_prefix)
			if isLeafNode(encoded_prefix) { // if leaf node
				match_arr := []uint8{}
				min_len := min(len(path_arr), len(nibble_arr))
				for i := 0; i < min_len; i++ {
					if path_arr[i] == nibble_arr[i] {
						match_arr = append(match_arr, path_arr[i])
					} else {
						break
					}
				}
				// case 1: nibbles and path do not match
				match_len := len(match_arr)
				if match_len == 0 {
					fmt.Println("No match")
					nibble_value := node.flag_value.value
					leaf_path_prefix := []uint8{}
					if(len(path_arr) > 1) {
						leaf_path_prefix = path_arr[1:]
					}
					leaf_nibble_prefix := []uint8{}
					if(len(nibble_arr) > 1) {
						leaf_nibble_prefix = nibble_arr[1:]
					}
					leaf_path_node := newLeafNode(leaf_path_prefix, new_value)
					leaf_nibble_node := newLeafNode(leaf_nibble_prefix, nibble_value)
					hash_leaf_path_node := leaf_path_node.hash_node()
					hash_leaf_nibble_node := leaf_nibble_node.hash_node()
					// create 1 branch node
					branch_value := [17]string{}
					branch_value[path_arr[0]] = hash_leaf_path_node
					fmt.Println("leaf1:",path_arr[0], new_value)
					branch_value[nibble_arr[0]] = hash_leaf_nibble_node
					fmt.Println("leaf2:",nibble_arr[0], nibble_value)
					branch_node := newBranchNode(branch_value, "")
					hash_branch_node := branch_node.hash_node()
					// delete the unwanted node
					delete(mpt.db, hash_node)
					// add all nodes to db
					mpt.db[hash_branch_node] = branch_node
					mpt.db[hash_leaf_path_node] = leaf_path_node
					mpt.db[hash_leaf_nibble_node] = leaf_nibble_node
					// update root
					mpt.root = hash_branch_node
					return
				} else if len(path_arr) == match_len && len(nibble_arr) == match_len { // case 2: complete match
					fmt.Println("Complete match")
					//????
				} else if len(path_arr) - match_len >= 1 && len(nibble_arr) - match_len >= 1 {
					fmt.Println("Partial match with extra nibble and extra path")
					path_arr = path_arr[match_len:]
					nibble_arr = nibble_arr[match_len:]
					//for i := 0; i < len(nibble_arr); i++ {
					//	fmt.Println("Leaf Nibbbbbbb: ", nibble_arr[i])
					//}
					nibble_value := node.flag_value.value
					// create 2 leaf nodes
					leaf_path_prefix := []uint8{}
					if(len(path_arr) > 1) {
						leaf_path_prefix = path_arr[1:]
					}
					leaf_nibble_prefix := []uint8{}
					if(len(nibble_arr) > 1) {
						leaf_nibble_prefix = nibble_arr[1:]
					}
					leaf_path_node := newLeafNode(leaf_path_prefix, new_value)
					//for i := 0; i < len(leaf_path_prefix); i++ {
					//	fmt.Println("Leaf Path Prefix: ", leaf_path_prefix[i])
					//}
					leaf_nibble_node := newLeafNode(leaf_nibble_prefix, nibble_value)
					//for i := 0; i < len(leaf_nibble_prefix); i++ {
					//	fmt.Println("Leaf Nibble Prefix: ", leaf_nibble_prefix[i])
					//}
					hash_leaf_path_node := leaf_path_node.hash_node()
					hash_leaf_nibble_node := leaf_nibble_node.hash_node()
					// create 1 branch node
					branch_value := [17]string{}
					branch_value[path_arr[0]] = hash_leaf_path_node
					fmt.Println("leaf1:",path_arr[0], new_value)
					branch_value[nibble_arr[0]] = hash_leaf_nibble_node
					fmt.Println("leaf2:",nibble_arr[0], nibble_value)
					branch_node := newBranchNode(branch_value, "")
					hash_branch_node := branch_node.hash_node()
					// create 1 extension node
					ext_node := newExtNode(match_arr, hash_branch_node)
					hash_ext_node := ext_node.hash_node()
					// delete the unwanted node
					delete(mpt.db, hash_node)
					// add all nodes to db
					mpt.db[hash_ext_node] = ext_node
					mpt.db[hash_branch_node] = branch_node
					mpt.db[hash_leaf_path_node] = leaf_path_node
					mpt.db[hash_leaf_nibble_node] = leaf_nibble_node
					// update root
					mpt.root = hash_ext_node
					return
				} else if len(path_arr) - match_len == 0 && len(nibble_arr) - match_len >= 1 {
					nibble_arr = nibble_arr[match_len:]
					//for i := 0; i < len(nibble_arr); i++ {
					//	fmt.Println("Leaf Nibbbbbbb: ", nibble_arr[i])
					//}
					nibble_value := node.flag_value.value
					// create 1 leaf nodes
					leaf_nibble_prefix := []uint8{}
					if(len(nibble_arr) > 1) {
						leaf_nibble_prefix = nibble_arr[1:]
					}
					leaf_nibble_node := newLeafNode(leaf_nibble_prefix, nibble_value)
					hash_leaf_nibble_node := leaf_nibble_node.hash_node()
					// create 1 branch node
					branch_value := [17]string{}
					branch_value[len(branch_value) - 1] = new_value
					fmt.Println("branch value:", nibble_value)
					branch_value[nibble_arr[0]] = hash_leaf_nibble_node
					fmt.Println("leaf2:",nibble_arr[0], nibble_value)
					branch_node := newBranchNode(branch_value, "")
					hash_branch_node := branch_node.hash_node()
					// create 1 extension node
					ext_node := newExtNode(match_arr, hash_branch_node)
					hash_ext_node := ext_node.hash_node()
					// delete the unwanted node
					delete(mpt.db, hash_node)
					// add all nodes to db
					mpt.db[hash_ext_node] = ext_node
					mpt.db[hash_branch_node] = branch_node
					mpt.db[hash_leaf_nibble_node] = leaf_nibble_node
					// update root
					mpt.root = hash_ext_node
					return
				} else if len(path_arr) - match_len >= 1 && len(nibble_arr) - match_len == 0 {
					path_arr = path_arr[match_len:]
					nibble_value := node.flag_value.value
					// create 1 leaf nodes
					leaf_path_prefix := []uint8{}
					if(len(path_arr) > 1) {
						leaf_path_prefix = path_arr[1:]
					}
					leaf_path_node := newLeafNode(leaf_path_prefix, new_value)
					hash_leaf_path_node := leaf_path_node.hash_node()
					// create 1 branch node
					branch_value := [17]string{}
					fmt.Println("branch value:", nibble_value)
					branch_value[path_arr[0]] = hash_leaf_path_node
					fmt.Println("leaf2:",nibble_arr[0], nibble_value)
					branch_node := newBranchNode(branch_value, nibble_value)
					hash_branch_node := branch_node.hash_node()
					// create 1 extension node
					ext_node := newExtNode(match_arr, hash_branch_node)
					hash_ext_node := ext_node.hash_node()
					// delete the unwanted node
					delete(mpt.db, hash_node)
					// add all nodes to db
					mpt.db[hash_ext_node] = ext_node
					mpt.db[hash_branch_node] = branch_node
					mpt.db[hash_leaf_path_node] = leaf_path_node
					// update root
					mpt.root = hash_ext_node
				} else {
					fmt.Println("check other cases")
				}
			} else { // if extension node

			}
		}
	}
}

func newBranchNode(branch_value [17]string, value string) Node {
	if(value != "") {
		branch_value[len(branch_value) - 1] = value
	}
	flag_value := Flag_value{[]uint8{},""}
	node := Node {
		1,
		branch_value,
		flag_value,
	}
	return node
}

func newExtNode(prefix []uint8, value string) Node {
	encoded_prefix := compact_encode(prefix)
	flag := Flag_value {
		encoded_prefix,
		value,
	}
	node := Node {
		2,
		[17]string{},
		flag,
	}
	return node
}

func newLeafNode(prefix []uint8, value string) Node {
	prefix = append(prefix, 16)
	encoded_prefix := compact_encode(prefix)
	flag := Flag_value {
		encoded_prefix,
		value,
	}
	node := Node {
		2,
		[17]string{},
		flag,
	}
	return node
}

func newNode(node_type int, nibble []uint8, value string) Node {
	if node_type == 1 {
		branch_value := [17]string{}
		branch_value[16] = value
		flag_value := Flag_value{
			[]uint8{},
			"",
		}
		node := Node {
			1,
			branch_value,
			flag_value,
		}
		return node
	} else if node_type == 2 {
		encoded_prefix := compact_encode(nibble)
		flag := Flag_value {
			encoded_prefix,
			value,
		}
		node := Node {
			2,
			[17]string{},
			flag,
		}
		return node
	}
	return Node {0,[17]string{}, Flag_value{[]uint8{},""}}
}

func (mpt *MerklePatriciaTrie) Delete(key string) {
	// TODO
}

// encode hex_array to ascii
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

func Test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
	fmt.Println("", compact_encode([]uint8{2, 6, 3, 16}))
	fmt.Println("", compact_decode(compact_encode([]uint8{5, 16})))
	//fmt.Println("", compact_decode(compact_encode([]uint8{2, 6, 3, 16})))
}

func Test_Get() {
	db := make(map[string]Node)
	root := Node{2, [17]string{}, Flag_value{encoded_prefix: compact_encode([]uint8{6,1,6,2,16}), value: "apple"}}
	//root := Node{2, [17]string{}, Flag_value{encoded_prefix: compact_encode([]uint8{6,1,16}), value: "apple"}}
	hash_root := root.hash_node()
	db[hash_root] = root
	mpt := MerklePatriciaTrie{db, hash_root}
	fmt.Println(mpt.Get("ab"))
}

func Test_ConvertStringToHexArray() {
	fmt.Println(ConvertStringToHexArray("a"))
}

func Test_Insert() {
	db := make(map[string]Node)
	root := ""
	mpt := MerklePatriciaTrie{db, root}
	mpt.Insert("p", "pap")
	//mpt.Insert("b", "banana")
	mpt.Insert("aaaaa", "5a")
	//fmt.Println("root", mpt.root)
	fmt.Println("Result:",mpt.Get( "p"))
	fmt.Println("Result:",mpt.Get("aaaaa"))
}

func Test_Create_Get() {
	db := make(map[string]Node)
	root := ""
	mpt := MerklePatriciaTrie{db, root}
	error := mpt.CreateTestMpt()
	if error != nil {
		fmt.Println("test1: ", mpt.Get("do"))
		fmt.Println("test2: ", mpt.Get("dog"))
		fmt.Println("test3: ", mpt.Get("doge"))
		fmt.Println("test4: ", mpt.Get("horse"))
	}
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