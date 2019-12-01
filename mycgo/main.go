package main

/*
#include <stdio.h>
*/
import (
	"C"
	"fmt"
)

func Random() int {
	return int(C.random())
}

func main() {
	fmt.Println(Random())
}
