package endpoints


import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"fmt"
	// "github.com/dgraph-io/badger"
	"log"

	"HackathonBlockchain/services"
)


type Set struct {
	NewBlockchainEndpoint endpoint.Endpoint
	PrintBlockchainEndpoint endpoint.Endpoint
	AddBlockEndpoint endpoint.Endpoint
}

func New(svc services.Service) Set {
	var newBlockchainEndpoint endpoint.Endpoint
	{
		newBlockchainEndpoint = MakeNewBlockchainEndpoint(svc)
	}
	var printBlockchainEndpoint endpoint.Endpoint
	{
		printBlockchainEndpoint = MakePrintBlockchainEndpoint(svc)
	}
	var addBlockEndpoint endpoint.Endpoint
	{
		addBlockEndpoint = MakeAddBlockEndpoint(svc)
	}
	return Set {
		NewBlockchainEndpoint: newBlockchainEndpoint,
		PrintBlockchainEndpoint: printBlockchainEndpoint,
		AddBlockEndpoint: addBlockEndpoint,
	}
}

func (s Set) NewBlockchain(ctx context.Context, Dataset []byte, baseline float64) ([]byte, error){
	resp, _ := s.NewBlockchainEndpoint(ctx, NewBlockchainRequest{Dataset, baseline})
	response := resp.(NewBlockchainResponse)
	log.Print("Endpoint: %+v",response)
	return response.Blockchain, response.Err
}

func MakeNewBlockchainEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(NewBlockchainRequest)
		bc, err := s.NewBlockchain(ctx, services.InitData{req.Dataset,req.Baseline})
		bcData,_ := json.Marshal(bc.Tip)
		return NewBlockchainResponse{Blockchain: bcData, Err: err}, nil
	}
}


type NewBlockchainRequest struct {
	Dataset []byte `json:"dataset"`
	Baseline float64 `json:"baseline"`
}

type NewBlockchainResponse struct {
	Blockchain []byte `json:"blockchain"`
	Err error `json:"err,omitempty"`
}


func (s Set) PrintBlockchain(ctx context.Context) ([]byte, error){
	resp, err := s.PrintBlockchainEndpoint(ctx, PrintBlockchainRequest{})
	if err != nil {
		return  nil, err
	}
	response := resp.(PrintBlockchainResponse)
	fmt.Printf("Response is %+v", response)
	return response.BlockchainIter, response.Err
}

func MakePrintBlockchainEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		bci, err := s.PrintBlockchain(ctx)
		bciData, _ := json.Marshal(bci.Blocks)
		return PrintBlockchainResponse{BlockchainIter: bciData, Err: err}, nil
	}
}

type PrintBlockchainRequest struct {
}

type PrintBlockchainResponse struct {
	BlockchainIter []byte 
	Err error 
}

func (s Set) AddBlock(ctx context.Context, md services.ModelData) (error){
	resp, err := s.AddBlockEndpoint(ctx, AddBlockRequest{Md:md})
	if err != nil {
		return err
	}
	response,_ := resp.(AddBlockResponse)
	return response.Err
}

func MakeAddBlockEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddBlockRequest)
		err = s.AddBlock(ctx, req.Md)
		return AddBlockResponse{Err: err}, nil
	}
}

type AddBlockRequest struct {
	Md services.ModelData `json:"modelData"`
}

type AddBlockResponse struct {
	Err error `json:"err,omitempty"`
}