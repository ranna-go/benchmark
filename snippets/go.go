// Implementation of basic bubble sort algorithm
// Reference: https://en.wikipedia.org/wiki/Bubble_sort
package main

import "fmt"

func bubbleSort(arr []int) []int {
	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < len(arr)-1; i++ {
			if arr[i+1] < arr[i] {
				arr[i+1], arr[i] = arr[i], arr[i+1]
				swapped = true
			}
		}
	}
	return arr
}

func main() {
	arr := []int{506, 48, 123, 79, 0, 362, 951, 500, 0}
	arr = bubbleSort(arr)
	fmt.Println(arr)
}
