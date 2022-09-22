package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ra111eo/TestTaskITgo/pb"
	"github.com/ra111eo/TestTaskITgo/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const deadlineTime = 5 //in seconds

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("Dial server %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Can not dial server: ", err)
	}
	ewalletClient := pb.NewEwalletServiceClient(conn)
	for {
		fmt.Printf("Please enter your command :")
		command := Scan()
		command = command[0 : len(command)-1] //delete \n from string
		s := strings.Split(command, " ")
		if s[0] == "ewallet" {
			switch s[1] {
			case "create":
				err = letsGoCreate(ewalletClient)
				if err != nil {
					fmt.Printf("Something goes wrong: %v\n", err)
				}
			case "send":
				var trans pb.Transaction
				trans.From = s[2]
				trans.To = s[3]
				trans.Amount, err = strconv.ParseFloat(s[4], 64)
				if err != nil {
					fmt.Println("Invalid amount try again")
				}
				err = letsGoSend(&trans, ewalletClient)
				if err != nil {
					fmt.Printf("Something goes wrong: %v\n", err)
				}
			case "getlast":
				err = letsGoGetlast(ewalletClient)
				if err != nil {
					fmt.Printf("Something goes wrong: %v\n", err)
				}
				fmt.Println("Invalid command, please try again")
			}
		} else {
			fmt.Println("Invalid command, please try again")
		}

	}
}

func Scan() string {
	in := bufio.NewReader(os.Stdin)
	str, err := in.ReadString('\n')
	if err != nil {
		return ""
	}
	return str
}

func letsGoCreate(ewalletClient pb.EwalletServiceClient) error {
	req := &pb.CreateEwalletRequest{
		Ewallet: sample.NewEwallet(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), deadlineTime*time.Second)
	defer cancel()
	res, err := ewalletClient.CreateEwallet(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Print("ewallet already exists")
		} else {
			log.Fatal("can not create ewallet: ", err)
		}
		return err
	}
	log.Printf("created ewallet with id: %s", res.Id)
	return nil
}

func letsGoSend(trans *pb.Transaction, ewalletClient pb.EwalletServiceClient) error {
	req := &pb.SendRequest{
		Transaction: trans,
	}
	ctx, cancel := context.WithTimeout(context.Background(), deadlineTime*time.Second)
	defer cancel()
	res, err := ewalletClient.Send(ctx, req)
	if err != nil || res.Amount != trans.Amount {
		log.Printf("transaction error: %v", err)
		return err
	}
	fmt.Println("transactions succsessful")
	return nil
}

func letsGoGetlast(ewalletClient pb.EwalletServiceClient) error {
	req := &pb.GetlastRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), deadlineTime*time.Second)
	defer cancel()
	res, err := ewalletClient.Getlast(ctx, req)
	if err != nil {
		log.Println(err)
	} else {
		jsonRes, err := json.Marshal(res)
		if err != nil {
			log.Fatal("Marshaling error: ", err)
		} else {
			fmt.Println(string(jsonRes))
		}
	}
	return nil
}
