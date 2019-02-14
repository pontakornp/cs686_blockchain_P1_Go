package p1

import (
	"cs686_blockchain_P1_Go/stack"
	"fmt"
)

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	// if root is empty, create a leaf and insert
	// if root is not empty, perform an operation according to each node type
	node_stack := stack.New()
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
	for hash_node != "" {
		node := mpt.db[hash_node]
		node_type := node.node_type
		if node_type == 0 {

		} else if node_type == 1 {
			// case where no more values in the path
			if len(path_arr) == 0 {
				fmt.Println("Branch, empty path")
				// insert the value at last index of branch_value
				node.branch_value[16] = new_value
				// hash the branch
				hash_branch_node := node.hash_node()
				// delete the branch from db
				delete(mpt.db, hash_node)
				// add branch node to db
				mpt.db[hash_branch_node] = node
				// update parents
				mpt.updateParents(node_stack, hash_branch_node)
				return
			} else {
				branch_prefix := path_arr[0]
				leaf_path_prefix := []uint8{}
				if len(path_arr) > 1 {
					leaf_path_prefix = path_arr[1:]
				}
				// case where first value in the path is empty, create leaf node
				if node.branch_value[branch_prefix] == "" {
					fmt.Println("Branch, contains path, empty branch index")
					leaf_path_node := newLeafNode(leaf_path_prefix, new_value)
					// hash leaf node
					hash_leaf_path_node := leaf_path_node.hash_node()
					// add leaf to the branch node
					node.branch_value[branch_prefix] = hash_leaf_path_node
					// hash branch node
					hash_branch_node := node.hash_node()
					// delete the branch from db
					delete(mpt.db, hash_node)
					// add all nodes to db
					mpt.db[hash_leaf_path_node] = leaf_path_node
					mpt.db[hash_branch_node] = node
					// update parents
					mpt.updateParents(node_stack, hash_branch_node)
					return
				} else { // case when first value in the path is not empty, traverse
					fmt.Println("Branch, contains path, branch index has value")
					parent := ParentNodeRef{hash_node, branch_prefix}
					fmt.Println("branch prefix", branch_prefix)
					// store parent in stack
					node_stack.Push(parent)
					// update the path
					path_arr = leaf_path_prefix
					// traverse
					hash_node = node.branch_value[branch_prefix]
				}
			}
		} else { // node_type == 2
			encoded_prefix := node.flag_value.encoded_prefix
			nibble_arr := compact_decode(encoded_prefix)
			match_arr := []uint8{}
			min_len := min(len(path_arr), len(nibble_arr))
			for i := 0; i < min_len; i++ {
				if path_arr[i] == nibble_arr[i] {
					match_arr = append(match_arr, path_arr[i])
				} else {
					break
				}
			}
			match_len := len(match_arr)
			if isLeafNode(encoded_prefix) { // if leaf node
				// case 1: complete match
				if len(path_arr) == match_len && len(nibble_arr) == match_len {
					fmt.Println("Leaf, Complete match")
					// if new value equal to the current leaf value, return
					if node.flag_value.value == new_value {
						return
					}
					//update the value to the new one
					node.flag_value.value = new_value
					//hash the node
					hash_leaf_node := node.hash_node()
					//delete the old leaf node
					delete(mpt.db, hash_node)
					//update mpt db
					mpt.db[hash_leaf_node] = node
					//update parents
					mpt.updateParents(node_stack, hash_leaf_node)
					return
				} else if match_len == 0 { // case 2: no match
					fmt.Println("Leaf, No match")
					nibble_value := node.flag_value.value
					//leaf_path_prefix := []uint8{}
					//leaf_nibble_prefix := []uint8{}
					branch_value := [17]string{}
					branch_node := newEmptyNode()
					hash_branch_node := ""
					if len(path_arr) != 0  && len(nibble_arr) != 0 {
						leaf_path_node, branch_path_index := mpt.TraverseLeafHelper(path_arr, new_value)
						leaf_nibble_node, branch_nibble_index := mpt.TraverseLeafHelper(nibble_arr, nibble_value)
						hash_leaf_path_node := leaf_path_node.hash_node()
						hash_leaf_nibble_node := leaf_nibble_node.hash_node()
						// delete the unwanted node
						delete(mpt.db, hash_node)
						//update leaves
						mpt.db[hash_leaf_nibble_node] = leaf_nibble_node
						mpt.db[hash_leaf_path_node] = leaf_path_node
						// create 1 branch node
						branch_value[branch_path_index] = hash_leaf_path_node
						branch_value[branch_nibble_index] = hash_leaf_nibble_node
						branch_node = newBranchNode(branch_value, "")
						hash_branch_node = branch_node.hash_node()
					} else {
						branch_node_val := ""
						if len(path_arr) == 0 {
							leaf_nibble_node, branch_nibble_index := mpt.TraverseLeafHelper(nibble_arr, nibble_value)
							hash_leaf_nibble_node := leaf_nibble_node.hash_node()
							//branch_value[nibble_arr[0]] = hash_leaf_nibble_node
							branch_value[branch_nibble_index] = hash_leaf_nibble_node
							branch_node_val = new_value
							fmt.Println("branch_path:", new_value)
							// delete the unwanted node
							delete(mpt.db, hash_node)
							// update mpt db leaf
							mpt.db[hash_leaf_nibble_node] = leaf_nibble_node
						} else if len(nibble_arr) == 0 {
							leaf_path_node, branch_path_index := mpt.TraverseLeafHelper(path_arr, new_value)
							hash_leaf_path_node := leaf_path_node.hash_node()
							branch_value[branch_path_index] = hash_leaf_path_node
							branch_node_val = nibble_value
							fmt.Println("branch_nibble:", nibble_value)
							// delete the unwanted node
							delete(mpt.db, hash_node)
							// update mpt db leaf
							mpt.db[hash_leaf_path_node] = leaf_path_node
						} else {
							fmt.Println("Bug Found")
							return
						}
						// create 1 branch node
						branch_node = newBranchNode(branch_value, branch_node_val)
						hash_branch_node = branch_node.hash_node()
					}
					// add all nodes to db
					mpt.db[hash_branch_node] = branch_node
					// update parents
					mpt.updateParents(node_stack, hash_branch_node)
					return
				} else if len(path_arr) - match_len >= 1 && len(nibble_arr) - match_len >= 1 { // case 3: partial match with extra nibble and extra path
					fmt.Println("Leaf, Partial match with extra nibble and extra path")
					path_arr = path_arr[match_len:]
					nibble_arr = nibble_arr[match_len:]
					nibble_value := node.flag_value.value
					// create 2 leaf nodes
					leaf_path_node, branch_path_index := mpt.TraverseLeafHelper(path_arr, new_value)
					leaf_nibble_node, branch_nibble_index := mpt.TraverseLeafHelper(nibble_arr, nibble_value)
					hash_leaf_path_node := leaf_path_node.hash_node()
					hash_leaf_nibble_node := leaf_nibble_node.hash_node()
					// create 1 branch node
					branch_value := [17]string{}
					branch_value[branch_path_index] = hash_leaf_path_node
					branch_value[branch_nibble_index] = hash_leaf_nibble_node
					branch_node := newBranchNode(branch_value, "")
					hash_branch_node := branch_node.hash_node()
					// create 1 extension node
					ext_node := newExtNode(match_arr, hash_branch_node)
					hash_ext_node := ext_node.hash_node()
					// delete the unwanted node
					delete(mpt.db, hash_node)
					// add all nodes to db
					mpt.db[hash_leaf_nibble_node] = leaf_nibble_node
					mpt.db[hash_leaf_path_node] = leaf_path_node
					mpt.db[hash_branch_node] = branch_node
					mpt.db[hash_ext_node] = ext_node
					// update parents
					mpt.updateParents(node_stack, hash_ext_node)
					return
				} else if len(path_arr) - match_len == 0 && len(nibble_arr) - match_len >= 1 { // case 4: partial match with extra nibble only
					fmt.Println("Leaf, Partial match with extra nibble only")
					nibble_arr = nibble_arr[match_len:]
					nibble_value := node.flag_value.value
					// create 1 leaf nodes
					leaf_nibble_node, branch_nibble_index := mpt.TraverseLeafHelper(nibble_arr, nibble_value)
					hash_leaf_nibble_node := leaf_nibble_node.hash_node()
					// create 1 branch node
					branch_value := [17]string{}
					fmt.Println("branch child value:", new_value)
					branch_value[branch_nibble_index] = hash_leaf_nibble_node
					//fmt.Println("leaf2:",nibble_arr[0], nibble_value)
					branch_node := newBranchNode(branch_value, new_value)
					hash_branch_node := branch_node.hash_node()
					// create 1 extension node
					ext_node := newExtNode(match_arr, hash_branch_node)
					hash_ext_node := ext_node.hash_node()
					// delete the unwanted node
					delete(mpt.db, hash_node)
					// add all nodes to db
					mpt.db[hash_leaf_nibble_node] = leaf_nibble_node
					mpt.db[hash_branch_node] = branch_node
					mpt.db[hash_ext_node] = ext_node
					// update parents
					mpt.updateParents(node_stack, hash_ext_node)
					return
				} else if len(path_arr) - match_len >= 1 && len(nibble_arr) - match_len == 0 { // case 5: partial match with extra path only
					fmt.Println("Leaf, Partial match with extra path only")
					path_arr = path_arr[match_len:]
					nibble_value := node.flag_value.value
					// create 1 leaf nodes
					leaf_path_node, branch_path_index := mpt.TraverseLeafHelper(path_arr, new_value)
					hash_leaf_path_node := leaf_path_node.hash_node()
					// create 1 branch node
					branch_value := [17]string{}
					fmt.Println("branch child value:", new_value)
					branch_value[branch_path_index] = hash_leaf_path_node
					//fmt.Println("leaf2:",nibble_arr[0], nibble_value)
					branch_node := newBranchNode(branch_value, nibble_value)
					hash_branch_node := branch_node.hash_node()
					// create 1 extension node
					ext_node := newExtNode(match_arr, hash_branch_node)
					hash_ext_node := ext_node.hash_node()
					// delete the unwanted node
					delete(mpt.db, hash_node)
					// add all nodes to db
					mpt.db[hash_leaf_path_node] = leaf_path_node
					mpt.db[hash_branch_node] = branch_node
					mpt.db[hash_ext_node] = ext_node
					// update parents
					mpt.updateParents(node_stack, hash_ext_node)
					return
				} else {
					fmt.Println("check other cases")
					return
				}
			} else { // if extension node
				if match_len == 0 { // case 1: no match
					fmt.Println("Ext, No match")
					// create leaf node, put path node in
					leaf_path_prefix := []uint8{}
					// get branch path prefix (first index)
					branch_path_prefix := path_arr[0]
					if(len(path_arr) > 1) {
						leaf_path_prefix = path_arr[1:]
					}
					leaf_path_node := newLeafNode(leaf_path_prefix, new_value)
					hash_leaf_path_node := leaf_path_node.hash_node()
					// get branch nibble prefix (first index)
					branch_nibble_prefix := nibble_arr[0]
					// create extension node if there's extra nibble left follows the branch node
					hash_ext_nibble_node := ""
					ext_nibble_node := newEmptyNode()
					if(len(nibble_arr) > 1) {
						ext_nibble_prefix := nibble_arr[1:]
						ext_nibble_node = newExtNode(ext_nibble_prefix, node.flag_value.value)
						hash_ext_nibble_node = ext_nibble_node.hash_node()
					}
					// create branch node, put hash of path and nibble in
					branch_value := [17]string{}
					// put hash path node in branch
					branch_value[branch_path_prefix] = hash_leaf_path_node
					// put hash nibble node in branch
					if(hash_ext_nibble_node != "") {
						branch_value[branch_nibble_prefix] = hash_ext_nibble_node
					} else {
						branch_value[branch_nibble_prefix] = node.flag_value.value
					}
					branch_node := newBranchNode(branch_value, "")
					hash_branch_node := branch_node.hash_node()
					// update db
					delete(mpt.db, hash_node)
					if(hash_ext_nibble_node != "") {
						mpt.db[hash_ext_nibble_node] = ext_nibble_node
					}
					mpt.db[hash_leaf_path_node] = leaf_path_node
					mpt.db[hash_branch_node] = branch_node
					// update parent
					mpt.updateParents(node_stack, hash_branch_node)
					return
				} else if len(path_arr) == match_len && len(nibble_arr) == match_len { // case 2: complete match
					fmt.Println("Ext, Complete match")
					// put parent in the stack
					parent := ParentNodeRef{hash_node, 17}
					node_stack.Push(parent)
					// update path
					path_arr = []uint8{}
					// traverse down the trie
					hash_node = node.flag_value.value
				} else if len(path_arr) - match_len >= 1 && len(nibble_arr) - match_len >= 1 { // case 3: partial match with extra nibble and extra path
					fmt.Println("Ext, Partial match with extra nibble and extra path")
					// store extension prefix
					ext_node_prefix := match_arr
					// remove extension node prefix
					remain_path_prefix := path_arr[match_len:]
					remain_nibble_prefix := nibble_arr[match_len:]
					// store branch path prefix
					branch_path_prefix := remain_path_prefix[0]
					// store branch nibble prefix
					branch_nibble_prefix := remain_nibble_prefix[0]
					// store leaf path prefix
					leaf_path_prefix := []uint8{}
					if len(remain_path_prefix) > 1 {
						leaf_path_prefix = remain_path_prefix[1:]
					}
					// create leaf path node
					leaf_path_node := newLeafNode(leaf_path_prefix, new_value)
					hash_leaf_path_node := leaf_path_node.hash_node()
					// store nibble prefix(es)
					// if extra nibble > 1, create extra extension node
					hash_ext_nibble_node := ""
					ext_nibble_node := newEmptyNode()
					if len(remain_nibble_prefix) > 1 {
						ext_nibble_prefix := remain_nibble_prefix[1:]
						ext_nibble_node = newExtNode(ext_nibble_prefix, node.flag_value.value)
						hash_ext_nibble_node = ext_nibble_node.hash_node()
					}
					// create branch
					branch_value := [17]string{}
					branch_value[branch_path_prefix] = hash_leaf_path_node
					// put hash children node in branch node
					if hash_ext_nibble_node != "" {
						branch_value[branch_nibble_prefix] = hash_ext_nibble_node
					} else {
						branch_value[branch_nibble_prefix] = node.flag_value.value
					}
					branch_node := newBranchNode(branch_value, "")
					hash_branch_node := branch_node.hash_node()
					// create extension node with the match prefix and put branch node in extension node
					ext_node := newExtNode(ext_node_prefix, hash_branch_node)
					hash_ext_node := ext_node.hash_node()
					// delete old extension node
					delete(mpt.db, hash_node)
					// update mpt db
					mpt.db[hash_leaf_path_node] = leaf_path_node
					if(hash_ext_nibble_node != "") {
						mpt.db[hash_ext_nibble_node] = ext_nibble_node
					}
					mpt.db[hash_branch_node] = branch_node
					mpt.db[hash_ext_node] = ext_node
					// update parents
					mpt.updateParents(node_stack, hash_ext_node)
					return
				} else if len(path_arr) - match_len == 0 && len(nibble_arr) - match_len >= 1 { // case 4: partial match with extra nibble only
					fmt.Println("Ext, Partial match with extra nibble only")
					// store extension prefix
					ext_node_prefix := match_arr
					// remove extension node prefix
					remain_prefix := nibble_arr[match_len:]
					// store branch prefix
					branch_prefix := remain_prefix[0]
					// store extension nibble prefix, create ext nibble node, if remain_prefix > 1
					hash_ext_nibble_node := ""
					ext_nibble_node := newEmptyNode()
					if len(remain_prefix) > 1 {
						ext_nibble_prefix := remain_prefix[1:]
						ext_nibble_node = newExtNode(ext_nibble_prefix, node.flag_value.value)
						hash_ext_nibble_node = ext_nibble_node.hash_node()
					}
					// create branch node
					branch_value := [17]string{}
					if hash_ext_nibble_node != "" {
						branch_value[branch_prefix] = hash_ext_nibble_node
					} else {
						branch_value[branch_prefix] = node.flag_value.value
					}
					branch_node := newBranchNode(branch_value, new_value)
					hash_branch_node := branch_node.hash_node()
					// create extension node
					ext_node := newExtNode(ext_node_prefix, hash_branch_node)
					hash_ext_node := ext_node.hash_node()
					// delete old node
					delete(mpt.db, hash_node)
					// update mpt db
					if hash_ext_nibble_node != "" {
						mpt.db[hash_ext_nibble_node] = ext_nibble_node
					}
					mpt.db[hash_branch_node] = branch_node
					mpt.db[hash_ext_node] = ext_node
					// update parent
					mpt.updateParents(node_stack, hash_ext_node)
					return
				} else if len(path_arr) - match_len >= 1 && len(nibble_arr) - match_len == 0 { // case 5: partial match with extra path only
					fmt.Println("Ext, Partial match with extra path only")
					//store in stack
					parent := ParentNodeRef{hash_node, 17}
					node_stack.Push(parent)
					// update path
					path_arr = path_arr[match_len:]
					//traverse down
					hash_node = node.flag_value.value
				}
			}
		}
	}
}

func (mpt *MerklePatriciaTrie) TraverseLeafHelper(path_or_nibble_arr []uint8, leaf_value string) (Node, uint8) {
	leaf_prefix := []uint8{}
	if(len(path_or_nibble_arr) > 1) {
		leaf_prefix = path_or_nibble_arr[1:]
	}
	//leaf_prefix := path_or_nibble_arr[1:]
	leaf_node := newLeafNode(leaf_prefix, leaf_value)
	branch_index := path_or_nibble_arr[0]
	fmt.Println("leaf_value:",branch_index, leaf_value)
	return leaf_node, branch_index
}