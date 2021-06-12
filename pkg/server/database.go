package server

import (
	accountModel "github.com/mrth1995/go-mockva/pkg/account/model"
	accountTrxModel "github.com/mrth1995/go-mockva/pkg/accounttransaction/model"
)

type Entity struct {
	Entity interface{}
}

func (s *Server) RegisterEntities() []Entity {
	return []Entity{
		{Entity: accountModel.Account{}},
		{Entity: accountTrxModel.AccountTransaction{}},
	}
}
