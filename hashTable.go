package main

import (
	"fmt"
)

type hashTable struct{
    points []int
    ids    []string
    size   int
    elements int
}

func Init(size int) *hashTable {
    table := &hashTable{}
    table.size = size
    table.elements = 0
    table.points = make([]int, int(size))
    table.ids = make([]string, int(size))
    return table
}

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

func (table *hashTable) remove(key string){
	index := hashFunc(key, table.size)
	done := 0
	for done != 1 {
	if table.ids[index] == key {  //feels wrong to use == with a string but according to internet this works i guess
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

func hashFunc(key string, size int) int{
	index := 1
	i := 0 
	for i < len(key) {
//		index += index * 7 + int(key[i])
		index += index + int(key[i])
		i += 1
	}
	if index < 0 {
           index = -index
	}
	return index % size
}

func (table *hashTable) get(key string) int{
	index := hashFunc(key, table.size)
	i := 0

	fmt.Println("get first try is index: ", index)
	fmt.Println("get given ID is: ", key)
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
