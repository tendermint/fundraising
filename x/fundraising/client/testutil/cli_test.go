//go:build norace
// +build norace

package testutil

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestTxCmdTestSuite(t *testing.T) {
	suite.Run(t, new(TxCmdTestSuite))
}

func TestQueryCmdTestSuite(t *testing.T) {
	suite.Run(t, new(QueryCmdTestSuite))
}
