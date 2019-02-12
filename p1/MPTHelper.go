package p1

import (
	"cs686_blockchain_P1_Go/stack"
	"fmt"
)

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (mpt *MerklePatriciaTrie) updateParents(node_stack *stack.Stack, hash_child string) {
	if node_stack.IsEmpty() {
		mpt.root = hash_child
		return
	}
	// while the stack is not empty, update the hash value of its children
	for !node_stack.IsEmpty() {
		parent := node_stack.Pop().(ParentNodeRef)
		old_hash_parent_node := parent.hash_node
		parent_node := mpt.db[old_hash_parent_node]
		fmt.Println("Updating parent.............", parent_node)
		// if the parent is a branch node
		fmt.Println("Len of map",len(mpt.db))
		if parent_node.node_type == 1 {
			branch_value_index := parent.index
			parent_node.branch_value[branch_value_index] = hash_child
		} else { // if the parent is an extension node
			parent_node.flag_value.value = hash_child
		}
		hash_parent_node := parent_node.hash_node()
		// delete old parent node
		delete(mpt.db, old_hash_parent_node)
		// update mpt db
		mpt.db[hash_parent_node] = parent_node
		// update child
		hash_child = hash_parent_node
		if node_stack.IsEmpty() {
			mpt.root = hash_parent_node
		}
	}
}