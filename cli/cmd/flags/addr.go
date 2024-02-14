package flags

import "fmt"

func GetAddr() string {
	var addr string
	if Expose {
		addr = fmt.Sprintf(":%d", Port)
	} else {
		addr = fmt.Sprintf("127.0.0.1:%d", Port)
	}
	return addr
}
