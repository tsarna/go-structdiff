package structdiff_test

import (
	"fmt"
	"time"

	"github.com/tsarna/go-structdiff"
)

func ExampleDiffStructs() {
	type User struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	oldUser := User{Name: "John", Age: 30, Email: "john@old.com"}
	newUser := User{Name: "John", Age: 31, Email: "john@new.com"}

	diff, _ := structdiff.DiffStructs(oldUser, newUser)
	fmt.Printf("Changes: %+v\n", diff)
	// Output: Changes: map[age:31 email:john@new.com]
}

func ExampleApplyToStruct() {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var user User
	patch := map[string]any{
		"name": "Alice",
		"age":  "25", // string converted to int
	}

	err := structdiff.ApplyToStruct(&user, patch)
	if err != nil {
		panic(err)
	}

	fmt.Printf("User: %+v\n", user)
	// Output: User: {Name:Alice Age:25}
}

func ExampleToMap() {
	type Config struct {
		Host     string  `json:"host"`
		Port     int     `json:"port"`
		Password *string `json:"password,omitempty"`
		Debug    bool    `json:"debug"`
	}

	config := Config{Host: "localhost", Port: 8080, Debug: false}
	m := structdiff.ToMap(config)

	fmt.Printf("Map: %+v\n", m)
	// Output: Map: map[debug:false host:localhost port:8080]
}

func ExampleDiffStructs_nested() {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type Employee struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	old := Employee{
		Name:    "Alice",
		Address: Address{Street: "123 Main St", City: "NYC"},
	}

	new := Employee{
		Name:    "Alice",
		Address: Address{Street: "456 Oak Ave", City: "NYC"},
	}

	diff, _ := structdiff.DiffStructs(old, new)
	fmt.Printf("Changes: %+v\n", diff)
	// Output: Changes: map[address:map[street:456 Oak Ave]]
}

func ExampleApplyToMap() {
	original := map[string]any{"x": 1, "y": 2}
	patch := map[string]any{"y": 3, "z": 4, "x": nil} // x deleted

	result := structdiff.ApplyToMap(original, patch)
	fmt.Printf("Result: %+v\n", result)
	// Output: Result: map[y:3 z:4]
}

func ExampleDiffStructs_roundTrip() {
	type Person struct {
		Name string    `json:"name"`
		Born time.Time `json:"born"`
	}

	alice := Person{
		Name: "Alice",
		Born: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	bob := Person{
		Name: "Bob",
		Born: time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC),
	}

	// Compute diff
	diff, _ := structdiff.DiffStructs(alice, bob)

	// Apply diff to transform alice into bob
	result := alice
	err := structdiff.ApplyToStruct(&result, diff)
	if err != nil {
		panic(err)
	}

	// Verify the transformation worked
	fmt.Printf("Original: %s, born %s\n", alice.Name, alice.Born.Format("2006-01-02"))
	fmt.Printf("Result: %s, born %s\n", result.Name, result.Born.Format("2006-01-02"))
	fmt.Printf("Matches target: %t\n", result.Name == bob.Name && result.Born.Equal(bob.Born))
	// Output: Original: Alice, born 1990-01-01
	// Result: Bob, born 1985-05-15
	// Matches target: true
}
