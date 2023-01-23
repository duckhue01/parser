package parsing

import (
	"fmt"

	"github.com/duckhue01/lexer"
)

func main() {
	l := lexer.New("asdas", nil)
	l.NextToken()

	fmt.Println("hello world!")
}
