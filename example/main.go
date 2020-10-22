package example

import (
	"context"
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

	b := true

	if (request.OK == b || ctx == nil) && request.Name != "nullable" {
		fs.ExecuteNone(ctx, request)
		return "no"
	}

	// if i := len([]int{}); i == 0 {
	// 	return fmt.Sprint(i)
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

//ExecuteNone
func (fs *FakeService) ExecuteNone(ctx context.Context, request Request) string {
	return "nil"
}
