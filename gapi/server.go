package gapi

import (
	"fmt"

	db "simplebank.com/db/sqlgen"
	"simplebank.com/pb"
	"simplebank.com/token"
	"simplebank.com/utils"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	token, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: token, config: config}

	return server, nil
}
