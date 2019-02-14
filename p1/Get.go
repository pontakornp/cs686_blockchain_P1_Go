package p1

import (
	"cs686_blockchain_P1_Go/stack"
	"errors"
	"fmt"
)

func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	hex_key := ConvertStringToHexArray(key)
	if(len(mpt.db) == 0 || mpt.root == "") {
		return "", errors.New("No value found, root is empty")
	}
	hash_node := mpt.root
	//fmt.Println(hash_node)
	for hash_node != "" {
		node := mpt.db[hash_node]
		//fmt.Println(node)
		node_type := node.node_type
		if isEmpty(node) || node_type == 0 { // null node
			return "", errors.New("No value found, node type is 0 or empty")
		} else if node_type == 1 { // branch node
			fmt.Println("branch")
			// if hex_key is empty string check if value exists
			// if yes, return value,
			// if not, return empty string
			if hex_key == nil || len(hex_key) == 0 {
				tempValue := node.branch_value[len(node.branch_value) - 1]
				if tempValue != "" {
					return tempValue, nil
				} else {
					return "", errors.New("No value found, node is empty")
				}
			}
			// update hash_node
			tempValue := node.branch_value[hex_key[0]]
			if tempValue != "" {
				hash_node = tempValue
			} else {
				return "", errors.New("No value found, branch at specific index is empty")
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
							return "", errors.New("No value found, leaf node does not match")
						}
					}
					return node.flag_value.value, nil
				}
				return "", errors.New("No value found, leaf node does not match")
			} else { // extension node
				// if hex_key length is less than key of the node, return empty string
				if len(hex_key) < len(decoded_arr) {
					return "", errors.New("No value found, extension nibble not match")
				}
				// if hex_key length is equal to or more than the key of the node
				// loop through each character of node key
				// if any character does not match, return empty string
				for i := 0; i < len(decoded_arr); i++ {
					if hex_key[i] != decoded_arr[i] {
						return "", errors.New("No value found, extension nibble not match")
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
					return "", errors.New("No value found, value of the next hash is empty")
				}
			}
		}
	}
	return "", errors.New("No value found")
}

func (mpt *MerklePatriciaTrie) GetStack(key string) (*stack.Stack, error) {
	hex_key := ConvertStringToHexArray(key)
	node_stack := stack.New()
	if(len(mpt.db) == 0 || mpt.root == "") {
		return nil, errors.New("No value found, root is empty")
	}
	hash_node := mpt.root
	//fmt.Println(hash_node)
	for hash_node != "" {
		node := mpt.db[hash_node]
		//fmt.Println(node)
		node_type := node.node_type
		if isEmpty(node) || node_type == 0 { // null node
			return nil, errors.New("No value found, node type is 0 or empty")
		} else if node_type == 1 { // branch node
			fmt.Println("branch")
			// if hex_key is empty string check if value exists
			// if yes, return value,
			// if not, return empty string
			if hex_key == nil || len(hex_key) == 0 {
				tempValue := node.branch_value[16]
				if tempValue != "" {
					ref := ParentNodeRef{hash_node,16}
					node_stack.Push(ref)
					return node_stack, nil
				} else {
					return nil, errors.New("No value found, node is empty")
				}
			}
			// update hash_node
			tempValue := node.branch_value[hex_key[0]]
			if tempValue != "" {
				ref := ParentNodeRef{hash_node,hex_key[0]}
				node_stack.Push(ref)
				hash_node = tempValue
			} else {
				return nil, errors.New("No value found, branch at specific index is empty")
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
							return nil, errors.New("No value found, leaf node does not match")
						}
					}
					ref := ParentNodeRef{hash_node,17}
					node_stack.Push(ref)
					return node_stack, nil
				}
				return nil, errors.New("No value found, leaf node does not match")
			} else { // extension node
				// if hex_key length is less than key of the node, return empty string
				if len(hex_key) < len(decoded_arr) {
					return node_stack, errors.New("No value found, extension nibble not match")
				}
				// if hex_key length is equal to or more than the key of the node
				// loop through each character of node key
				// if any character does not match, return empty string
				for i := 0; i < len(decoded_arr); i++ {
					if hex_key[i] != decoded_arr[i] {
						return node_stack, errors.New("No value found, extension nibble not match")
					}
				}
				//if the remaining key length is equal to zero, then set hex_key to nil
				remaining_len := len(hex_key) - len(decoded_arr)
				if remaining_len == 0 {
					hex_key = nil
				} else { // if the remaining key length is more than zero
					hex_key = hex_key[len(decoded_arr):]
				}
				ref := ParentNodeRef{hash_node,17}
				node_stack.Push(ref)
				hash_node = node.flag_value.value
				// if value of the next hash node is empty then return empty string
				if hash_node == "" {
					return nil, errors.New("No value found, value of the next hash is empty")
				}
			}
		}
	}
	return node_stack, errors.New("No value found")
}
