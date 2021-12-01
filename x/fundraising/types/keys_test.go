package types_test

import (
	fmt "fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

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

func (s *keysTestSuite) TestGetSequenceIndexKey() {
	testCases := []struct {
		auctionID uint64
		sequence  uint64
		expected  []byte
	}{
		{
			uint64(5),
			uint64(10),
			[]byte{0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa},
		},
		{
			uint64(2),
			uint64(10),
			[]byte{0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa},
		},
		{
			uint64(3),
			uint64(5),
			[]byte{0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5},
		},
	}

	for _, tc := range testCases {
		key := types.GetSequenceIndexKey(tc.auctionID, tc.sequence)
		s.Require().Equal(tc.expected, key)

		auctionID, sequence := types.ParseSequenceIndexKey(key)
		s.Require().Equal(tc.auctionID, auctionID)
		s.Require().Equal(tc.sequence, sequence)
	}
}

func (s *keysTestSuite) TestGetBidKey() {
	// TODO: not implemented yet
	key := []byte{49, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 10}

	fmt.Println("key: ", key)
	// key -> [49 0 0 0 0 0 0 0 5 0 0 0 0 0 0 0 10]
	// key[1:] -> [0 0 0 0 0 0 0 5 0 0 0 0 0 0 0 10]
	// sdk.BigEndianToUint64(key[1:] -> 5

	fmt.Println(sdk.BigEndianToUint64(key[1+8:]))

}

func (s *keysTestSuite) TestGetBidderKey() {
	// TODO: not implemented yet
}
