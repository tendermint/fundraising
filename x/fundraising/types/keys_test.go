package types_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

type keysTestSuite struct {
	suite.Suite
}

func TestKeysTestSuite(t *testing.T) {
	suite.Run(t, new(keysTestSuite))
}

func (s *keysTestSuite) TestGetAuctionKey() {
	s.Require().Equal([]byte{0x21, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, types.GetAuctionKey(0))
	s.Require().Equal([]byte{0x21, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9}, types.GetAuctionKey(9))
	s.Require().Equal([]byte{0x21, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa}, types.GetAuctionKey(10))
}
