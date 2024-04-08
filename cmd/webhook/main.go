package main

import "github.com/selectel/external-dns-selectel-webhook/cmd/webhook/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
