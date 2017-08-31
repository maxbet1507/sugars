package sugars

import "fmt"

func ExampleOnetime() {
	f1 := Onetime(func() {
		fmt.Println("f1")
	})
	f2 := Onetime(func() {
		fmt.Println("f2")
	})

	fmt.Println(1)
	f1()
	fmt.Println(2)
	f1()
	fmt.Println(3)
	f2()
	fmt.Println(4)
	f2()
	fmt.Println(5)
	f1()
	fmt.Println(6)
	f2()

	// Output:
	// 1
	// f1
	// 2
	// 3
	// f2
	// 4
	// 5
	// 6
}
