package main

import (
	"fmt"
	"github.com/wjoseperez20/zenwallet/pkg/auth"
)

func main() {
	fmt.Println(auth.GenerateRandomKey())
}
