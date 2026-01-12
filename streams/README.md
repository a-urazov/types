# Stream Package

This package provides a simple implementation of a stream processing API for Go collections.

## Usage

```go
package main

import (
	"fmt"
	"github.com/your-repo/types/streams"
)

func main() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := streams.Of(data).
		Filter(func(i int) bool {
			return i%2 == 0
		}).
		Map(func(i int) int {
			return i * 2
		}).
		Collect()

	fmt.Println(result) // Output: [4 8 12 16 20]
}
```
