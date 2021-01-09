package main

import "fmt"

func main() {
	employee := map[string]string{}

	employee["name"] = "knrt10"
	employee["age"] = "23"
	employee["tag"] = "Cloud Infrastructure Engineer"
	employee["location"] = "India"

	for key, value := range employee {
		fmt.Printf("%s:%s, ", key, value)
	}
}
