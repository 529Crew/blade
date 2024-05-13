package pump_monitor_hooks

import (
	"fmt"

	"github.com/529Crew/blade/idls/pump"
	"github.com/529Crew/blade/internal/sol"
	"github.com/529Crew/blade/internal/util"
	"github.com/gagliardetto/solana-go"
)

func ParseCreateAndBuy(tx *solana.Transaction, sig string) error {
	transaction, err := sol.ResolveAddressTables(tx)
	if err != nil {
		return err
	}
	tx = transaction

	parseCreateInst := func(inst *solana.CompiledInstruction) error {
		instAccs, err := inst.ResolveInstructionAccounts(&transaction.Message)
		if err != nil {
			return err
		}

		instruction, err := pump.DecodeInstruction(instAccs, inst.Data)
		if err != nil {
			return err
		}

		createInst, ok := instruction.Impl.(*pump.Create)
		if !ok {
			return fmt.Errorf("error casting instruction to create: %v", err)
		}
		util.PrettyPrint(createInst)

		return nil
	}

	parseBuyInst := func(inst *solana.CompiledInstruction) error {
		instAccs, err := inst.ResolveInstructionAccounts(&transaction.Message)
		if err != nil {
			return err
		}

		instruction, err := pump.DecodeInstruction(instAccs, inst.Data)
		if err != nil {
			return err
		}

		buyInst, ok := instruction.Impl.(*pump.Buy)
		if !ok {
			return fmt.Errorf("error casting instruction to buy: %v", err)
		}
		util.PrettyPrint(buyInst)

		return nil
	}

	/* search for create or buy inst */
	for _, inst := range tx.Message.Instructions {
		if len([]byte(inst.Data)) < 8 {
			continue
		}
		var discriminator [8]byte
		copy(discriminator[:], []byte(inst.Data)[:8])

		switch {
		case discriminator == [8]byte{24, 30, 200, 40, 5, 28, 7, 119}: /* create */
			parseCreateInst(&inst)
		case discriminator == [8]byte{102, 6, 61, 18, 1, 218, 235, 234}: /* buy */
			parseBuyInst(&inst)
		}
	}

	return nil
}
