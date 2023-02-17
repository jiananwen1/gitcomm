package gitcomm

import (
	"fmt"
	"strings"
	"testing"
)

func Test_gitWithOutput(t *testing.T) {
	out, err := gitWithOutput("log", "-1")
	if err != nil {
		panic(err)
	}
	//fmt.Println(out)

	lists := strings.Split(out, "\n")
	fmt.Println(len(lists))
	for _, each := range lists {
		if strings.Contains(each, "--story") {
			fmt.Println(each)

			parts := strings.Split(each, "=")
			fmt.Println(parts[1])
		}

	}
}
