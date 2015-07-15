package aflShared

import (
	"log"

	"github.com/go-errors/errors"
)

//HandleError prints out calls stacks for returned errors
func HandleError(err error) {

	if err != nil {
		log.Println(err)
		log.Println(err.(*errors.Error).ErrorStack())
	}
}

//KellyCriterion calculation
func KellyCriterion(price float32, percentage float32) float32 {
	return ((percentage*price - 1) / (price - 1))
}
