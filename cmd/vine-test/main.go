package main

import "github.com/umizu/vine"

func main() {
	v := vine.New()
	v.Start(":9999")
}
