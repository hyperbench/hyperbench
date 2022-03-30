package main

import (
	"fmt"
	"github.com/hyperbench/hyperbench/cmd"
)

func main() {

	err := cmd.InitCmd(debug)
	if err != nil {
		fmt.Println("cmd init fail: ", err)
		return
	}

	err = cmd.GetRootCmd().Execute()
	if err != nil {
		fmt.Println("cmd execute fail: ", err)
	}
}
