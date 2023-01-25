package main


import (
    "fmt"
    "github.com/gin-gonic/gin"
    "encoding/hex"
    "strings"
    "strconv"
    "math"
)

/*
* receipt struct
* Retailer: string containing the name of the retailer
* PurchaseDate: String containing the date in the format yyyy-mm-dd
* PurchaseTime: String containing the time of the purchase in format hh:mm
* Items: an array of item structs
* Total: contains the total in the form d+.dd
*/
type receipt struct{
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []item `json:"items"`
	Total        string `json:"total"`
}
/*
* item struct
* ShortDescription: string containing a description of item
* Price: contains item price in form d+.dd
*/
type item struct{
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

/*
* id struct
* ID the id given in the form hhhh-hhhh-hhhh-hhhh-hhhh where h means hexadecimal
*/
type id struct{
	ID string `json:"id" uri:"id"`
}

// hashtable that stores ids and associated points values
var table *hashTable

/*
* returns a Json object containing the points value of the requested id
* path: /receipts/{id}/points
* 
* arguments
* context - a gin context variable that is used to extract the id from the path
*/
func getPoints(context *gin.Context){
	var idIn id

	err := context.BindUri(&idIn)
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

/*
* postReceipt takes a Json receipt object as an argument through the context variable
* then generates an ID and associated point value and stores it within the hashtable
*
* arguments
* context - a gin context variable that is used to extract the id from the path
*/
func postReceipt (context *gin.Context){
    var newReceipt receipt
    
    err := context.BindJSON(&newReceipt)
    if err != nil {
    context.String(400, "The receipt is invalid, failed to bind\n")
    }
    
    ID := getId(&newReceipt)
    points := 0

    // get points from retailer's name
    // note I chose not to upper or lower case or trim this as I believe
    // it would make it less efficient as those functions would be running these types of compares
    // anyways
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
    TotalLen := len(newReceipt.Total)

    // cents should be an integer value associated with the final two digits of newReceipt.Total
    cents, err := strconv.Atoi(newReceipt.Total[TotalLen-2:TotalLen])
    if err != nil {	    
        context.String(400, "The receipt is invalid, total failed to read\n")
    }
    // I chose not to do cents % 25 for the second part to avoid the division operator
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
                context.String(400, "The receipt is invalid, failed to read an items price\n")
	    }
            points += int(math.Ceil(price * 0.2))
	}
        i += 1
    }
    // points for every 2 items
    // pulled it out of the loop for increased efficiency
    points += (i >> 1) * 5

    // points from odd date
    // note that ascii numbers that are odd have odd ascii values so this does work
    if int(newReceipt.PurchaseDate[9]) & 0x1 != 0 {
        points += 6
    }

    //points from time of purchase
    hour, err := strconv.Atoi(newReceipt.PurchaseTime[0:2])
    if err != nil {
        context.String(400, "The receipt is invalid, did not get numbers in first two digits of time\n")    
    }
    // Wasn't sure If I should include 2:00 technically any time with 2:00 in the purchaseTime
    // is some number of seconds after 2pm if the prompt means exactly 2pm, not hard to change if this is not
    // what was intended by the problem
    if hour >= 14 && hour < 16 {
        points += 10
    }

    // Right now there is nothing to handle identical ID's being generated which could result in incorrect behavior
    // currently this is not tested
    table.insert(ID, points, &table) 
    context.JSON(200, gin.H{"id": ID})
}

/*
* Generates an id of the form hhhh-hhhh-hhhh-hhhh-hhhh where h is a hexadecimal character
* based on characteristics of the input receipt
*
*/
func getId(newReceipt *receipt) string{

	// Tried to do something that somewhat emulated the example id's 
	// (5 sections with hex characters separated by -) but I don't think that this generation
	// method is ideal since it can have ID collisions which result in undesired behavior

	// this first bit is here because the original method does not function if there aren't at least 4 characters
	// in the first item's short description and the retailer field
	RetailerComponent := ""
	DescriptionComponent := ""
	if len(newReceipt.Retailer) >= 4 {
		RetailerComponent = hex.EncodeToString([]byte(newReceipt.Retailer))[0:4]
	} else {
		RetailerComponent = hex.EncodeToString([]byte(newReceipt.Retailer))[0:1]
                RetailerComponent = RetailerComponent + RetailerComponent + RetailerComponent + RetailerComponent
	}
	if len(newReceipt.Items[0].ShortDescription) >= 4 {
		DescriptionComponent = hex.EncodeToString([]byte(newReceipt.Items[0].ShortDescription))[0:4]
	} else {
		DescriptionComponent = hex.EncodeToString([]byte(newReceipt.Items[0].ShortDescription))[0:1]
                DescriptionComponent = DescriptionComponent + DescriptionComponent + DescriptionComponent + 
		DescriptionComponent
	}

	// concatenates the 5 sections
	returnValue := hex.EncodeToString([]byte(newReceipt.Total))[0:4] + "-" +
	    DescriptionComponent + "-" +
      	    hex.EncodeToString([]byte(newReceipt.PurchaseTime))[0:4] + "-" +
	    hex.EncodeToString([]byte(newReceipt.PurchaseDate))[0:4] + "-" +
	    RetailerComponent

return returnValue
}

func main() {
    // Initialize the hashTable
    table = Init(10)
    // Initialize the gin router
    router := gin.Default()
    // tie /receipts/{id}/points to the getPoints function
    router.GET("/receipts/:id/*points", getPoints)
    // tie /receipts/process to the postReceipt function
    router.POST("/receipts/process", postReceipt)
    // set trusted proxies to nil for safety
    router.SetTrustedProxies(nil)
    // run the router
    router.Run("localhost:8080")
}

