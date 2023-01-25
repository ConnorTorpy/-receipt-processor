package main


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
	index := hashFunc(key)
	done := 0
	for done != 1 {
	if table.ids[index] != "" {
		index += 1
		if index == 100 {
			index -= 100
		}
	} else {
		done = 1
		table.ids[index] = key
		table.points[index] = value
		table.elements +=1
		if (table.elements << 1) >= table.size {
			*tableAddr = table.resize()
		}
	}
    }
}

func (table *hashTable) remove(key string){
	index := hashFunc(key)
	done := 0
	for done != 1 {
	if table.ids[index] == key {  //feels wrong to use == with a string but according to internet this works i guess
		done = 1
		table.ids[index] = ""
		table.elements -= 1
	} else {
		index += 1
		if index == 100 {
			index -= 100
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

func hashFunc(key string) int{
	index := 1
	i := 0 
	for i < len(key) {
		index += index * 7 + int(key[i])
	}
	return index
}

func (table *hashTable) get(key string) int{
	index := hashFunc(key)
	i := 0
	for i < table.size{
		if table.ids[index] == key {
			return table.points[i]
		}
		i += 1
		index += 1
		if index == 100 {
			index -= 100
		}
	}
	return -1
}
