package sample

import (
	"github.com/ra111eo/TestTaskITgo/pb"
)

func NewEwallet() *pb.Ewallet {
	wallet := &pb.Ewallet{
		Id:      randomID(),
		Balance: randomFloat64(1000.0, 3000.0),
		Get:     []*pb.Transaction{},
		Send:    []*pb.Transaction{},
	}
	return wallet
}
