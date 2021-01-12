package main

import (
	"fmt"
	"github.com/meshplus/hyperbench/cmd"
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
