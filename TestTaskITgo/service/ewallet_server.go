package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/ra111eo/TestTaskITgo/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EwalletServer struct {
	Store EwalletStore
}

func NewEwalletServer(store EwalletStore) *EwalletServer {
	return &EwalletServer{store}
}

// CreateEwallet
func (server *EwalletServer) CreateEwallet(
	ctx context.Context,
	req *pb.CreateEwalletRequest,
) (*pb.CreateEwalletResponse, error) {
	ewallet := req.GetEwallet()
	log.Printf("recieve a create-ewallet request with id %s", ewallet.Id)
	if len(ewallet.Id) > 0 {
		_, err := uuid.Parse(ewallet.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "ewallet ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := GetNewId()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "can not generate a new ewallet ID: %v", err)
		}
		ewallet.Id = id
	}

	if ctx.Err() == context.Canceled {
		log.Print("request is canceled")
		return nil, status.Error(codes.Canceled, "request is canceled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	err := server.Store.Save(ewallet)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "can not save ewallet to store: %v", err)
	}

	log.Printf("saved ewallet with id: %s", ewallet.Id)

	res := &pb.CreateEwalletResponse{
		Id: ewallet.Id,
	}
	return res, nil
}

func (server *EwalletServer) Send(ctx context.Context,
	req *pb.SendRequest,
) (*pb.SendResponse, error) {
	transaction := req.GetTransaction()
	log.Printf("recieve a transaction from %s to %s with amount %f", transaction.From, transaction.To, transaction.Amount)
	err := server.Store.Send(transaction)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not get transaction: %v", err)
	}
	if ctx.Err() == context.Canceled {
		log.Print("request is canceled")
		return nil, status.Error(codes.Canceled, "request is canceled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	res := &pb.SendResponse{
		Amount: transaction.Amount,
	}
	return res, nil
}

func (server *EwalletServer) Getlast(
	ctx context.Context, req *pb.GetlastRequest,
) (*pb.GetlastResponse, error) {
	if ctx.Err() == context.Canceled {
		log.Print("request is canceled")
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	res := &pb.GetlastResponse{}
	res.JsonData = string(server.Store.Getlast())
	return res, nil
}

/*
func (server *EwalletServer) Update(
	ctx context.Context,
	req *pb.UpdateStoreRequest,
) (*pb.UpdateStoreResponse, error) {
	req.GetKey()
	log.Printf("++++ IM IN UPDATESTORE FUNC")
	if ctx.Err() == context.Canceled {
		log.Print("request is canceled")
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	err := server.Store.UpdateDBStore()
	if err != nil {
		log.Fatal("update failed")
		return nil, err
	}
	res := &pb.UpdateStoreResponse{}
	return res, nil
}
*/
