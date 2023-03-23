package example

import (
	"github.com/KevinZonda/batchq"
	"time"
)

func main() {
	q := batchq.NewBatchQEasy(10, 10*time.Second)
	q.Start()

}
