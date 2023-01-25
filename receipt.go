package main


import (
    "fmt"
    "github.com/gin-gonic/gin"
    "encoding/hex"
    "strings"
    "strconv"
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

type id struct{
	ID string `json:"id"`
}

var table *hashTable

func getPoints(context *gin.Context){
	var idIn id
	err := context.BindJSON(&idIn)
	if err != nil {
        context.String(404, "Id format invalid\n")
	}
        points := table.get(idIn.ID)
        if points == -1 {
        context.String(404, "No receipt found for that id\n")
	} else {
	context.JSON(200, gin.H{"points": points})
	}
}

func postReceipt (context *gin.Context){
    var newReceipt receipt
    
    err := context.BindJSON(&newReceipt)
    if err != nil {
    context.String(400, "The receipt is invalid\n")
    }
    ID := getId(&newReceipt)
    points := 0

    // get points from retailer's name
    i := 0
    for i < len(newReceipt.Retailer) {
        if ((int(newReceipt.Retailer[i]) >= 48 && int(newReceipt.Retailer[i]) <= 57) ||
	     (int(newReceipt.Retailer[i]) >= 65 && int(newReceipt.Retailer[i]) <= 90) ||
	     (int(newReceipt.Retailer[i]) >= 97 && int(newReceipt.Retailer[i]) <= 122)) {
	         points += 1
	     }
	     i += 1
    }

    // get points from total
    cents, err := strconv.Atoi(newReceipt.Total[len(newReceipt.Total)-2:len(newReceipt.Total)-1])
    if err != nil {	    
        context.String(400, "The receipt is invalid\n")
    }
    if(cents == 0) {
        points += 75
    } else if cents == 25 || cents == 50 || cents == 75 {
        points += 25
    }

    // get points from items
    i = 0
    for i < len(newReceipt.Items) {
        if len(strings.Trim(newReceipt.Items[i].ShortDescription, " ")) % 3 == 0 {
            price, err := strconv.ParseFloat(newReceipt.Items[i].Price, 32)
	    if err != nil {
                context.String(400, "The receipt is invalid\n")
	    }
            points += -int(price * -.2)
	}
        i += 1
    }
    // points for every 2 items
    points += i >> 1

    // points from odd date
    if int(newReceipt.PurchaseDate[9]) & 0x1 != 0 {
        points += 6
    }

    //points from time of purchase
    hour, err := strconv.Atoi(newReceipt.PurchaseTime[0:1])
    if err != nil {
        context.String(400, "The receipt is invalid\n")    
    }
    if hour >= 14 && hour < 16 {
        points += 6
    }

    //TODO add something to handle identical id attempts in insert method
    table.insert(ID, points, &table) 

    context.JSON(200, gin.H{"id": ID})
}

func getId(newReceipt *receipt) string{

returnValue := hex.EncodeToString([]byte(newReceipt.Total)) + "-" +
               hex.EncodeToString([]byte(newReceipt.Items[0].ShortDescription)) + "-" +
	       hex.EncodeToString([]byte(newReceipt.PurchaseTime)) + "-" +
	       hex.EncodeToString([]byte(newReceipt.PurchaseDate)) + "-" +
	       hex.EncodeToString([]byte(newReceipt.Retailer))

return returnValue
}

func main() {
    table = Init(10)
    fmt.Println("Hello World")
    router := gin.Default()
    router.GET("/receipts/:id/*points", getPoints)
    router.POST("/receipts/process", postReceipt)
    router.SetTrustedProxies(nil)
    router.Run("localhost:8080")
}

