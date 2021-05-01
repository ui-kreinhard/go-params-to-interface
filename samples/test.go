package main

import "fmt"
import "strings"

// SayHello says Hello
func SayHello(greetings []string) {
	fmt.Println(joinStrings(greetings))
}

// joinStrings joins strings
func joinStrings(words []string) string {
	return strings.Join(words, ", ")
}

func add(a, b int, c int) int {
	return a + b
}

type Test struct {
	
}

func (t *Test) add(a int, b int, c int, d int) {
	
}