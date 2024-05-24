package main

import (
	"fmt"
	"time"
)

func main() {

	a := time.Now()
	fmt.Println(a)
	l, _ := time.LoadLocation("Asia/Almaty")
	fmt.Println(a.In(l))
	fmt.Printf("%v", a.String())
}
