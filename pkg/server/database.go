package server

import (
	"github.com/mrth1995/go-mockva/pkg/account"
	"github.com/mrth1995/go-mockva/pkg/accounttransaction"
)

type Entity struct {
	Entity interface{}
}

func (s *Server) RegisterEntities() []Entity {
	return []Entity{
		{Entity: account.Account{}},
		{Entity: accounttransaction.AccountTransaction{}},
	}
}
