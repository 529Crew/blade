package pump_tx

import (
	"context"
	"math/big"
	"math/rand"

	"github.com/529Crew/blade/idls/pump"
	"github.com/529Crew/blade/internal/client"
	"github.com/529Crew/blade/internal/config"
	"github.com/529Crew/blade/internal/constants"
	"github.com/529Crew/blade/internal/logger"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/rpc"
)

func Buy(buyInst *pump.Buy, solSpent float64, tokenBalance float64) error {
	cfg := config.Get()

	instructions := []solana.Instruction{
		computebudget.NewSetComputeUnitLimitInstruction(200_420).Build(),
		computebudget.NewSetComputeUnitPriceInstruction(config.Get().BuyPriorityFee + uint64(rand.Intn(1000))).Build(),
	}

	mint := buyInst.GetMintAccount().PublicKey.String()

	logger.Log.Printf("[BUY] attempting to buy mint %s\n", mint)

	associatedUser, _, err := solana.FindAssociatedTokenAddress(constants.Wallet.PublicKey(), buyInst.GetMintAccount().PublicKey)
	if err != nil {
		logger.Log.Printf("[BUY] error buying mint %s: %v", mint, err)
		return err
	}

	_, err = client.GetUtil().GetAccountInfo(context.Background(), associatedUser)
	if err != nil {
		if err.Error() == "not found" {
			instructions = append(instructions,
				associatedtokenaccount.NewCreateInstruction(
					constants.Wallet.PublicKey(),
					constants.Wallet.PublicKey(),
					buyInst.GetMintAccount().PublicKey,
				).Build(),
			)
		} else {
			logger.Log.Printf("[BUY] error buying mint %s: %v", mint, err)
			return err
		}
	}

	/* get quote using bundled buy amount and sol spent */

	sol_buy_amount := new(big.Int).SetUint64(cfg.BuyAmount) // amount to buy in sol

	// initial_real_token_reserves := new(big.Int).SetUint64(793_100_000_000_000)
	initial_virtual_token_reserves := new(big.Int).SetUint64(1_073_000_000_000_000)
	initial_virtual_sol_reserves := new(big.Int).SetUint64(30_000_000_000)

	tokens_bundle_bought := new(big.Int).SetUint64(uint64(tokenBalance))
	sol_spent := new(big.Int).SetUint64(uint64(solSpent))

	// real_token_reserves := new(big.Int).Sub(initial_real_token_reserves, tokens_bundle_bought)
	virtual_token_reserves := new(big.Int).Sub(initial_virtual_token_reserves, tokens_bundle_bought)
	virtual_sol_reserves := new(big.Int).Add(initial_virtual_sol_reserves, sol_spent)

	mul_result := new(big.Int).Mul(virtual_sol_reserves, virtual_token_reserves)
	add_result := new(big.Int).Add(virtual_sol_reserves, sol_buy_amount)
	div_result := new(big.Int).Quo(mul_result, add_result)
	add_result_2 := new(big.Int).Add(div_result, new(big.Int).SetUint64(1))
	tokens_out := new(big.Int).Sub(virtual_token_reserves, add_result_2)

	sol_buy_amount_percent := new(big.Int).Mul(sol_buy_amount, new(big.Int).SetInt64(cfg.BuySlippage))
	sol_slippage := new(big.Int).Quo(sol_buy_amount_percent, new(big.Int).SetInt64(100))
	sol_quote_with_slippage := new(big.Int).Add(sol_buy_amount, sol_slippage)

	logger.Log.Printf("[BUY] quote %d max lamports for %d tokens with %d%% slippage", sol_quote_with_slippage.Int64(), tokens_out, cfg.BuySlippage)

	/* generate jito tip inst */
	logger.Log.Printf("[BUY] using %d lamports jito tip", cfg.JitoTip)
	tipInst, err := client.GetJito().GenerateTipRandomAccountInstruction(cfg.JitoTip, constants.Wallet.PublicKey())
	if err != nil {
		logger.Log.Printf("[BUY] error buying mint %s: %v", mint, err)
		return err
	}

	instructions = append(instructions,
		/* jito tip inst */
		tipInst,
		/* pump fun buy inst */
		pump.NewBuyInstruction(
			tokens_out.Uint64(),
			sol_quote_with_slippage.Uint64(),
			buyInst.GetGlobalAccount().PublicKey,
			buyInst.GetFeeRecipientAccount().PublicKey,
			buyInst.GetMintAccount().PublicKey,
			buyInst.GetBondingCurveAccount().PublicKey,
			buyInst.GetAssociatedBondingCurveAccount().PublicKey,
			associatedUser,
			constants.Wallet.PublicKey(),
			buyInst.GetSystemProgramAccount().PublicKey,
			buyInst.GetTokenProgramAccount().PublicKey,
			buyInst.GetRentAccount().PublicKey,
			buyInst.GetEventAuthorityAccount().PublicKey,
			buyInst.GetProgramAccount().PublicKey,
		).Build(),
	)

	block, err := client.Get().GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		logger.Log.Printf("[BUY] error buying mint %s: %v", mint, err)
		return err
	}

	tx, err := solana.NewTransaction(
		instructions,
		block.Value.Blockhash,
		solana.TransactionPayer(constants.Wallet.PublicKey()),
	)
	if err != nil {
		logger.Log.Printf("[BUY] error buying mint %s: %v", mint, err)
		return err
	}

	if _, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if constants.Wallet.PublicKey().Equals(key) {
				return &constants.Wallet
			}
			return nil
		},
	); err != nil {
		logger.Log.Printf("[BUY] error buying mint %s: %v", mint, err)
		return err
	}

	logger.Log.Printf("[BUY] submitting jito bundle to buy mint %s\n", mint)
	resp, err := client.GetJito().BroadcastBundle([]*solana.Transaction{tx})
	if err != nil {
		logger.Log.Printf("error buying %s: %v", mint, err)
		return err
	}
	logger.Log.Printf("[BUY] submitted jito bundle to buy mint %s: %s\n", mint, resp.Uuid)

	return nil
}
