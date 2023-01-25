package main

import (
	"fmt"
)

/*
* hashTable struct
* points - an array of ints containing receipt point values
* ids - an array of strings containing reciept ids
* size - maximum number of elements of the array
* elements - current number of elements in the array
*/
type hashTable struct{
    points []int
    ids    []string
    size   int
    elements int
}

/*
* Initializes a hashTable of an input size
*
* arguments
* size - the size that the hashTable should be initialized to
*/
func Init(size int) *hashTable {
    table := &hashTable{}
    table.size = size
    table.elements = 0
    table.points = make([]int, int(size))
    table.ids = make([]string, int(size))
    return table
}

/*
* inserts a new key value pair into the hash table
*
* arguments
* table - the hashTable to be inserted into
* key - the id associated with the request
* value - the points associated with the receipt
* tableAddr a pointer to the pointer containing the hashTables location
* it is changed in the case of a resize
*/
func (table *hashTable) insert(key string, value int, tableAddr **hashTable){
	fmt.Println("insert got: ", value)
	index := hashFunc(key, table.size)
	done := 0
	for done != 1 {
        fmt.Println("index: ", index)
	if table.ids[index] != "" {
		index += 1
		if index == table.size {
			index -= table.size
		}
	} else {
		done = 1
		table.ids[index] = key
		table.points[index] = value
		table.elements += 1
		fmt.Println("just inserted into index: ", index)
		fmt.Println("inserted id: ", table.ids[index])
		fmt.Println("inserted points: ", table.points[index])
	
		if (table.elements << 1) >= table.size {
			*tableAddr = table.resize()
		}
	}
    }
}

/*
* removes an entry associated with input key from the hashTable
* note - this function was unused and is therefore untested in the receipt-processor implementation
*
* arguments
* key - id associated with the request
*/
func (table *hashTable) remove(key string){
	index := hashFunc(key, table.size)
	done := 0
	for done != 1 {
	if table.ids[index] == key {  //feels wrong to use == with a string but according to internet this works
		done = 1
		table.ids[index] = ""
		table.elements -= 1
	} else {
		index += 1
		if index == table.size {
			index -= table.size
		}
	}
     }
}

/*
* returns a new hashTable with 16 times the size as the input hashTable 
* but containing all the same key value pairs
*/
func (table *hashTable) resize() *hashTable{
	newTable := Init(table.size << 4)
	i := 0
	for i < table.size {
		if table.ids[i] != "" {
			newTable.insert(table.ids[i], table.points[i], nil)
		}
		i += 1
	}
	return newTable
}

/*
* hashFunc takes a key and size and returns an index for the hashTable based on it's
* function
* 
* arguments
* key - a string that contains the ID of the receipt
* size - the size of the hashTable, determines the limit of how large a number
* the hash function can return
*/
func hashFunc(key string, size int) int{
	index := 1
	i := 0
	// not a great hash function imo expected value is 0 so might tend toward lower numbers
	// increasing collision chance
	for i < len(key) {
//		index += index * 7 + int(key[i])
		index += index + int(key[i])
		i += 1
	}
	// might be better to absolute value each individual addition
	if index < 0 {
           index = -index
	}
	// modulo size so returned value is within desired range
	return index % size
}
/*
* get returns the point value associated with given key if it exists in the hash table
* returns -1 on failure to find the receipt of the givne id 
*
* arguments
* key - string of the id associated with the request
* table - the hashtable to be searched
*/
func (table *hashTable) get(key string) int{
	index := hashFunc(key, table.size)
	i := 0

	for i < table.size{
		if table.ids[index] == key {

			fmt.Println("found a match at index: ", index)
			fmt.Println("points value found: ", table.points[index])
			return table.points[index]
		}
		i += 1
		index += 1
		if index == table.size {
			index -= table.size
		}
	}
	return -1
}
