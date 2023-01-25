package main


type hashTable struct{
    receipts []receipt
    size int
}

func Init(size int) *hashTable {
    table := &hashTable{}
    table.size = size
    table.receipts = make([]receipt, int(size))
    return table
}

func insert(){

}

func remove(){

}

func resize(){

}
