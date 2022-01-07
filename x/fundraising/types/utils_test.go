package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestMustParseRFC3339(t *testing.T) {
	normalCase := "9999-12-31T00:00:00Z"
	normalRes, err := time.Parse(time.RFC3339, normalCase)
	require.NoError(t, err)
	errorCase := "9999-12-31T00:00:00_ErrorCase"
	_, err = time.Parse(time.RFC3339, errorCase)
	require.PanicsWithError(t, err.Error(), func() { types.MustParseRFC3339(errorCase) })
	require.Equal(t, normalRes, types.MustParseRFC3339(normalCase))
}

func TestDeriveAddress(t *testing.T) {
	testCases := []struct {
		addressType     types.AddressType
		moduleName      string
		name            string
		expectedAddress string
	}{
		{
			types.ReserveAddressType,
			types.ModuleName,
			"SellingReserveAcc|1",
			"cosmos18xzvtd72y9j8xyf8a36z5jjhth7qgtcwhh8lz7yee3tvxqn6ll5quh78zq",
		},
		{
			types.ReserveAddressType,
			types.ModuleName,
			"PayingReserveAcc|1",
			"cosmos18permjyqvk5flft8ey9egr7hd4ry8tauqt4f9mg9knn4vvtkry9sujucrl",
		},
		{
			types.ReserveAddressType,
			types.ModuleName,
			"VestingReserveAcc|1",
			"cosmos1gukaqt783nhz79uhcqklsty7lc7jfyy8scn5ke4x7v0m3rkpt4dst7y4l3",
		},
		{
			types.AddressType20Bytes,
			"",
			"fee_collector",
			"cosmos17xpfvakm2amg962yls6f84z3kell8c5lserqta",
		},
		{
			types.AddressType32Bytes,
			"farming",
			"GravityDEXFarmingBudget",
			"cosmos1228ryjucdpdv3t87rxle0ew76a56ulvnfst0hq0sscd3nafgjpqqkcxcky",
		},
		{
			types.AddressType20Bytes,
			types.ModuleName,
			"",
			"cosmos1vh7g0ypukt8xyxm3zlf8f2t4sjnzxe63pe3cap",
		},
		{
			types.AddressType20Bytes,
			types.ModuleName,
			"test1",
			"cosmos1n7h778sm85f0x6h76nlrcd57eza702m6gskhhv",
		},
		{
			types.AddressType32Bytes,
			types.ModuleName,
			"test1",
			"cosmos1zrwtzgxy5urtfwp5r9t0qkeuynh78k7z2047sqafrx9hg8x4rq0qjspx0y",
		},
		{
			types.AddressType32Bytes,
			"test2",
			"",
			"cosmos1v9ejakp386det8xftkvvazvqud43v3p5mmjdpnuzy3gw84h4dwxsfn6dly",
		},
		{
			types.AddressType32Bytes,
			"test2",
			"test2",
			"cosmos1qmsgyd6yu06uryqtw7t6lg7ua5ll7s3ej828fcqfakrphppug4xqcx7w45",
		},
		{
			types.AddressType20Bytes,
			"",
			"test2",
			"cosmos1vqcr4c3tnxyxr08rk28n8mkphe6c5gfuk5eh34",
		},
		{
			types.AddressType20Bytes,
			"test2",
			"",
			"cosmos1vqcr4c3tnxyxr08rk28n8mkphe6c5gfuk5eh34",
		},
		{
			types.AddressType20Bytes,
			"test2",
			"test2",
			"cosmos15642je7gk5lxugnqx3evj3jgfjdjv3q0nx6wn7",
		},
		{
			3,
			"test2",
			"invalidAddressType",
			"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedAddress, types.DeriveAddress(tc.addressType, tc.moduleName, tc.name).String())
		})
	}
}
