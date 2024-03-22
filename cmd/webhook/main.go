package main

import "github.com/selectel/external-dns-webhook/cmd/webhook/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
