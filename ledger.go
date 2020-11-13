package main

import (
	"flag"
	"fmt"
)

func main() {
	var from = flag.String("from", "", "bucket from which the amount is taken")
	var to = flag.String("to", "", "bucket into which the amount is deposited")
	var date = flag.String("date", "", "date of transaction")
	var amt = flag.Uint("amt", 0, "amount in cents of the transaction")
	flag.Parse()

	fmt.Printf("value of flags: from=%s,to=%s,date=%s,amt=%d\n", *from, *to, *date, *amt)
}
