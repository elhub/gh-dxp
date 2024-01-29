package main

import (
	"github.com/michaeloa/gh-devxp/pkg/cmd"
)

var (
	version = "0.1.0"
)

func main() {
	cmd.Execute(version)
	/*
		fmt.Println("hi world, this is the gh-devxp extension!")
		client, err := api.DefaultRESTClient()
		if err != nil {
			fmt.Println(err)
			return
		}
		response := struct{ Login string }{}
		err = client.Get("user", &response)
		if err != nil {
			fmt.Println(err)
			return
		}
	*/
	//fmt.Printf("running as %s\n", response.Login)
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
