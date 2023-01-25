package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
)


type receipt struct{
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []item `json:"items"`
	Total        string `json:"total"`
}

type item struct{
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}


func getPoints(context *gin.Context){
// returning a good request looks like this
//    context.JSON(200, gin.H{"points": "numpoints"})

//returning a bad request looks like this
    context.String(404, "No receipt found for that id\n")
}

func postReceipt (context *gin.Context){
    var newReceipt receipt
    
    err := context.BindJSON(&newReceipt)
    if err != nil {
    context.String(400, "The receipt is invalid\n")
    }





// returning a good request looks like this
//    context.JSON(200, gin.H{"ID": "fealfiushdflzoidsufe"})

//returning a bad request looks like this
    context.String(400, "The receipt is invalid\n")
}

func main() {
    fmt.Println("Hello World")
    router := gin.Default()
    router.GET("/receipts/:id/*points", getPoints)
    router.POST("/receipts/process", postReceipt)
    router.SetTrustedProxies(nil)
    router.Run("localhost:8080")
}

