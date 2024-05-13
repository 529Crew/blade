package sol

import (
	"context"

	"github.com/529Crew/blade/internal/client"
	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
)

func ResolveAddressTables(tx *solana.Transaction) (*solana.Transaction, error) {
	if !tx.Message.IsVersioned() {
		return tx, nil
	}

	tblKeys := tx.Message.GetAddressTableLookups().GetTableIDs()
	if len(tblKeys) == 0 {
		return tx, nil
	}

	numLookups := tx.Message.GetAddressTableLookups().NumLookups()
	if numLookups == 0 {
		return tx, nil
	}

	resolutions := make(map[solana.PublicKey]solana.PublicKeySlice)
	for _, key := range tblKeys {
		info, err := client.GetUtil().GetAccountInfo(
			context.Background(),
			key,
		)
		if err != nil {
			return nil, err
		}

		tableContent, err := addresslookuptable.DecodeAddressLookupTableState(info.GetBinary())
		if err != nil {
			return nil, err
		}

		resolutions[key] = tableContent.Addresses
	}

	err := tx.Message.SetAddressTables(resolutions)
	if err != nil {
		return nil, err
	}

	err = tx.Message.ResolveLookups()
	if err != nil {
		return nil, err
	}

	return tx, nil
}
