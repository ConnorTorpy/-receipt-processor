I ran this inside a linux vm on my windows machine. Couldn't get docker working on my computer since my windows version doesn't
support hyper-v.

first you should download the master branch of this github into a directory of your choosing (mine is called receipt-processor)

According to the github you should already have a go environment so the main other thing you will need that you may not have is
gin. You can use go mod to handle the dependency or you can download it yourself at https://github.com/gin-gonic/gin.

Then you should compile the project, I did so by running go build . inside the receipt-processor directory

In the linux environment that I was in I then opened two terminals and ran ./receipt-processor in one and initiated get
and post reqeusts from the other.

to initiate a post request in linux use the following command (you may change the port if you wish)
curl localhost:8080/receipts/process --include --header "Content-Type: application/json" -d @receipts/receiptX.json
where X is the receipt number you wish to post.
Feel free to have a look at receipts/expecteds.txt for more information on the test cases

to initiate a get requst in linux use the following command
curl localhost:8080/receipts/{id}/points
note the id will be returned after running the post request
