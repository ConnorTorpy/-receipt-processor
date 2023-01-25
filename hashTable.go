package main


type hashTable struct{
    points []int
    ids    []string
    size   int
}

func Init(size int) *hashTable {
    table := &hashTable{}
    table.size = size
    table.points = make([]int, int(size))
    table.ids = make([]string, int(size))
    return table
}

func insert(){

}

func remove(){

}

func resize(){

}
