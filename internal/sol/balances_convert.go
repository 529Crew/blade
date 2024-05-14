package sol

import (
	"github.com/529Crew/blade/internal/types"
	"github.com/weeaa/jito-go/proto"
)

func ConvertPreAndPostBalances(preBalances []uint64, postBalances []uint64) ([]int64, []int64) {
	preBalancesNew := []int64{}
	for _, preBalance := range preBalances {
		preBalancesNew = append(preBalancesNew, int64(preBalance))
	}

	postBalancesNew := []int64{}
	for _, postBalance := range postBalances {
		postBalancesNew = append(postBalancesNew, int64(postBalance))
	}

	return preBalancesNew, postBalancesNew
}

func ConvertPreAndPostTokenBalances(preTokenBalances []*proto.TokenBalance, postTokenBalances []*proto.TokenBalance) ([]types.TokenBalance, []types.TokenBalance) {
	preTokenBalancesNew := []types.TokenBalance{}
	for _, preTokenBalance := range preTokenBalances {
		preTokenBalancesNew = append(preTokenBalancesNew, types.TokenBalance{
			Mint:  preTokenBalance.Mint,
			Owner: preTokenBalance.Owner,
			UITokenAmount: types.UiTokenAmount{
				UIAmount: preTokenBalance.UiTokenAmount.UiAmount,
			},
		})
	}

	postTokenBalancesNew := []types.TokenBalance{}
	for _, postTokenBalance := range postTokenBalances {
		postTokenBalancesNew = append(postTokenBalancesNew, types.TokenBalance{
			Mint:  postTokenBalance.Mint,
			Owner: postTokenBalance.Owner,
			UITokenAmount: types.UiTokenAmount{
				UIAmount: postTokenBalance.UiTokenAmount.UiAmount,
			},
		})
	}

	return preTokenBalancesNew, postTokenBalancesNew
}
