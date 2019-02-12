package p1

import "fmt"

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

	//case 1
	//mpt.Insert("a", "apple")
	//mpt.Insert("p", "papaya")
	//mpt.Insert("abc", "hahaha")
	//fmt.Println("	Result:",mpt.Get( "a"))
	//fmt.Println("	Result:",mpt.Get("p"))
	//fmt.Println("	Result:",mpt.Get("abc"))

	//case 2
	//mpt.Insert("p", "papaya")
	//mpt.Insert("aaaaa", "5a")
	//mpt.Insert("aa", "new")
	//fmt.Println("	Result:",mpt.Get( "p"))
	//fmt.Println("	Result:",mpt.Get("aaaaa"))
	//fmt.Println("	Result:",mpt.Get("aa"))

	//case 3
	//mpt.Insert("a", "apple")
	//mpt.Insert("b", "banana")
	//mpt.Insert("p", "papaya")
	////fmt.Println("	Result:",mpt.Get( "a"))
	////fmt.Println("	Result:",mpt.Get("b"))
	////fmt.Println("	Result:",mpt.Get("p"))

	//case 4
	//mpt.Insert("c", "cat")
	//mpt.Insert("aa", "ant")
	//mpt.Insert("ap", "orange")
	//fmt.Println("	Result:",mpt.Get( "p"))
	//fmt.Println("	Result:",mpt.Get("aa"))
	//fmt.Println("	Result:",mpt.Get("ap"))

	//fmt.Println("	Result:",mpt.Get("a"))

	//case 5
	//mpt.Insert("p", "apple")
	//mpt.Insert("aaaa", "banana")
	//mpt.Insert("aaaap", "orange")
	//mpt.Insert("aa", "new")
	//fmt.Println("	Result:",mpt.Get( "p"))
	//fmt.Println("	Result:",mpt.Get("aaaa"))
	//fmt.Println("	Result:",mpt.Get("aaaap"))
	//fmt.Println("	Result:",mpt.Get("aa"))

	//case 6
	//mpt.Insert("a", "10")
	//mpt.Insert("b", "20")
	//mpt.Insert("t", "30")
	//mpt.Insert("ab", "40")
	//fmt.Println("	Result:",mpt.Get( "a"))
	//fmt.Println("	Result:",mpt.Get("b"))
	//fmt.Println("	Result:",mpt.Get("t"))
	//fmt.Println("	Result:",mpt.Get("ab"))

	//case 7
	//mpt.Insert("a", "10")
	//mpt.Insert("b", "20")
	//mpt.Insert("ab", "30")
	//fmt.Println("	Result:",mpt.Get( "a"))
	//fmt.Println("	Result:",mpt.Get("b"))
	//fmt.Println("	Result:",mpt.Get("ab"))

	//case 8
	//mpt.Insert("a", "10")
	//mpt.Insert("p", "20")
	//mpt.Insert("ab", "30")
	//fmt.Println("	Result:",mpt.Get( "a"))
	//fmt.Println("	Result:",mpt.Get("p"))
	//fmt.Println("	Result:",mpt.Get("ab"))

	//case 9
	//mpt.Insert("p", "10")
	//mpt.Insert("a", "20")
	//mpt.Insert("ab", "30")
	//mpt.Insert("b", "40")
	//fmt.Println(mpt.Get( "p"))
	//fmt.Println(mpt.Get("a"))
	//fmt.Println(mpt.Get("ab"))
	//fmt.Println(mpt.Get("b"))

	//case 10
	mpt.Insert("a", "10")
	mpt.Insert("p", "20")
	fmt.Println(mpt.Get( "a"))
	fmt.Println(mpt.Get("b"))

	// map len
	fmt.Println(mpt.Order_nodes())
	fmt.Println("Map len", len(mpt.db))
}

//func Test_Create_Get() {
//	db := make(map[string]Node)
//	root := ""
//	mpt := MerklePatriciaTrie{db, root}
//	error := mpt.CreateTestMpt()
//	if error != nil {
//		fmt.Println("test1: ", mpt.Get("do"))
//		fmt.Println("test2: ", mpt.Get("dog"))
//		fmt.Println("test3: ", mpt.Get("doge"))
//		fmt.Println("test4: ", mpt.Get("horse"))
//	}
//}