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
	ID string `json:"id" uri:"id"`
}

var table *hashTable

func getPoints(context *gin.Context){
	var idIn id
	err := context.BindUri(&idIn)
	if err != nil {
        context.String(404, "Id format invalid\n")
        }
	fmt.Println(idIn.ID)
        points := table.get(idIn.ID)
	fmt.Println("get got this: ", points)
        if points == -1 {
        context.String(404, "No receipt found for that id\n")
	} else {
	context.JSON(200, gin.H{"points": points})
	}
}

func postReceipt (context *gin.Context){
    var newReceipt receipt
    fmt.Println("Posting Receipt")
    
    err := context.BindJSON(&newReceipt)
    if err != nil {
    context.String(400, "The receipt is invalid, failed to bind\n")
    }
    
    fmt.Println("Receipt Bound")
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
    
    fmt.Println("Name Points Received: ", points)

    // get points from total
    TotalLen := len(newReceipt.Total)
    cents, err := strconv.Atoi(newReceipt.Total[TotalLen-2:TotalLen])
    fmt.Println("cents: ", cents)
    if err != nil {	    
        context.String(400, "The receipt is invalid, total failed to read\n")
    }
    if(cents == 0) {
        points += 75
    } else if cents == 25 || cents == 50 || cents == 75 {
        points += 25
    }
    fmt.Println("Total Points Received: ", points)

    // get points from items
    i = 0
    for i < len(newReceipt.Items) {
        if len(strings.Trim(newReceipt.Items[i].ShortDescription, " ")) % 3 == 0 {
            price, err := strconv.ParseFloat(newReceipt.Items[i].Price, 32)
	    if err != nil {
                context.String(400, "The receipt is invalid, failed to read an items price\n")
	    }

            fmt.Println("item with multiple of 3 length found")
            points += int(price * 0.2) + 1
	}
        i += 1
    }
    // points for every 2 items
    points += (i >> 1) * 5
    fmt.Println("Item Points Received: ", points)

    // points from odd date
    if int(newReceipt.PurchaseDate[9]) & 0x1 != 0 {
        points += 6
    }
    fmt.Println("Odd Date Points Received: ", points)

    //points from time of purchase
    hour, err := strconv.Atoi(newReceipt.PurchaseTime[0:1])
    if err != nil {
        context.String(400, "The receipt is invalid, did not get numbers in first two digits of time\n")    
    }
    if hour >= 14 && hour < 16 {
        points += 6
    }
    fmt.Println("Purchase Time Points Received: ", points)

    //TODO add something to handle identical id attempts in insert method
    table.insert(ID, points, &table) 

    fmt.Println("Returning ID")
    context.JSON(200, gin.H{"id": ID})
}

func getId(newReceipt *receipt) string{

	returnValue := hex.EncodeToString([]byte(newReceipt.Total))[0:4] + "-" +
	hex.EncodeToString([]byte(newReceipt.Items[0].ShortDescription))[0:4] + "-" +
	hex.EncodeToString([]byte(newReceipt.PurchaseTime))[0:4] + "-" +
	hex.EncodeToString([]byte(newReceipt.PurchaseDate))[0:4] + "-" +
	hex.EncodeToString([]byte(newReceipt.Retailer))[0:4]

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

