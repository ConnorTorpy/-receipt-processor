package main


type hashTable struct{
    receipts []receipt
    size int
}

func init(size int) *hashTable {
    table := &hashTable{}
    table.size = size
    table.receipts = [size]receipt
}

func insert(){

}

func remove(){

}

func resize(){

}
