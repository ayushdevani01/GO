package main

import "fmt"

func Add(a int, b int) int {
	return a + b
}

func Devide(a int, b int) (int, error) {

	if b == 0 {
		return 0, fmt.Errorf("Cannot devide by zero")
	}
	return a / b, nil
}
