package sol

import (
	"github.com/gagliardetto/solana-go"
	"github.com/weeaa/jito-go/proto"
)

func ConvertTx(tx *proto.Transaction) *solana.Transaction {
	/* convert message */
	newMessage := solana.Message{}

	/* account keys */
	newMessage.AccountKeys = []solana.PublicKey{}
	for _, byteKey := range tx.Message.AccountKeys {
		newMessage.AccountKeys = append(newMessage.AccountKeys, solana.PublicKeyFromBytes(byteKey))
	}

	/* header */
	newMessage.Header = solana.MessageHeader{
		NumRequiredSignatures:       uint8(tx.Message.Header.NumRequiredSignatures),
		NumReadonlySignedAccounts:   uint8(tx.Message.Header.NumReadonlySignedAccounts),
		NumReadonlyUnsignedAccounts: uint8(tx.Message.Header.NumReadonlyUnsignedAccounts),
	}

	/* recent blockhash */
	newMessage.RecentBlockhash = solana.HashFromBytes(tx.Message.RecentBlockhash)

	/* instructions */
	newMessage.Instructions = []solana.CompiledInstruction{}
	for _, inst := range tx.Message.Instructions {
		accounts := []uint16{}
		for _, acc := range inst.Accounts {
			accounts = append(accounts, uint16(acc))
		}

		newMessage.Instructions = append(
			newMessage.Instructions,
			solana.CompiledInstruction{
				ProgramIDIndex: uint16(inst.ProgramIdIndex),
				Accounts:       accounts,
				Data:           inst.Data,
			},
		)
	}

	/* address table lookups */
	newMessage.AddressTableLookups = []solana.MessageAddressTableLookup{}
	for _, lookup := range tx.Message.AddressTableLookups {
		newMessage.AddressTableLookups = append(
			newMessage.AddressTableLookups,
			solana.MessageAddressTableLookup{
				AccountKey:      solana.PublicKey(lookup.AccountKey),
				WritableIndexes: lookup.WritableIndexes,
				ReadonlyIndexes: lookup.ReadonlyIndexes,
			},
		)
	}

	/* convert signatures */
	newSignatures := []solana.Signature{}
	for _, sig := range tx.Signatures {
		newSignatures = append(newSignatures, solana.SignatureFromBytes(sig))
	}

	/* set message version */
	if tx.Message.Versioned {
		newMessage.SetVersion(solana.MessageVersionV0)
	}

	return &solana.Transaction{
		Signatures: newSignatures,
		Message:    newMessage,
	}
}
