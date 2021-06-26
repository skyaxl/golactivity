package example

import (
	"context"
	"fmt"
)

type FakeService struct {
}
type Request struct {
	Name string `json:"name,omitempty"`
	OK   bool
}

//Execute test
//@draw
func (fs *FakeService) Execute(ctx context.Context, request Request) string {

	array1 := [][]int{{1}}
	for i := 0; i < len(array1); i++ {
		array2 := array1[i]
		if j := len(array2); j == 0 {
			return fmt.Sprint(j)
		}
	}

	// b := true

	// if (request.OK == b || ctx == nil) && request.Name != "nullable" {
	// 	fs.ExecuteNone(ctx, request)
	// 	return "no"
	// }

	// if 2*2 == 4 {
	// 	return "4"
	// }

	// if !request.OK {
	// 	return "none"
	// }
	// if request.Name == "" {
	// 	return "name"
	// }

	// fs.ExecuteNone(ctx, request)

	// validator := func() bool {
	// 	return true
	// }
	// validator()

	// switch request.Name {
	// case "a":
	// 	{
	// 		return request.Name + request.Name
	// 	}
	// case "go":
	// 	{
	// 		c1 := make(chan string)
	// 		go func() {
	// 			c1 <- "test"
	// 		}()
	// 		return <-c1
	// 	}
	// default:
	// 	{
	// 		return request.Name
	// 	}
	// }
	return ""
}

func (fs *FakeService) Two() {

}

//ExecuteNone
func (fs *FakeService) ExecuteNone(ctx context.Context, request Request) string {
	return "nil"
}
