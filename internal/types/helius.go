package types

type TransactionSubscribePayload struct {
	Jsonrpc string                       `json:"jsonrpc"`
	ID      int                          `json:"id"`
	Method  string                       `json:"method"`
	Params  []TransactionSubscribeParams `json:"params"`
}

type TransactionSubscribeParams struct {
	Vote                           bool     `json:"vote"`
	Failed                         bool     `json:"failed"`
	Signature                      string   `json:"signature,omitempty"`
	AccountInclude                 []string `json:"accountInclude"`
	AccountExclude                 []string `json:"accountExclude"`
	AccountRequired                []string `json:"accountRequired"`
	Commitment                     string   `json:"commitment"`
	Encoding                       string   `json:"encoding"`
	TransactionDetails             string   `json:"transactionDetails,omitempty"`
	ShowRewards                    bool     `json:"showRewards"`
	MaxSupportedTransactionVersion int      `json:"maxSupportedTransactionVersion"`
}

type TransactionSubscribeResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  int64  `json:"result"`
	ID      int    `json:"id"`
}

type TransactionNotification struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Subscription int64 `json:"subscription"`
		Result       struct {
			Transaction struct {
				Transaction []string `json:"transaction"`
				Meta        struct {
					Err    any `json:"err"`
					Status struct {
						Ok any `json:"Ok"`
					} `json:"status"`
					Fee               int     `json:"fee"`
					PreBalances       []int64 `json:"preBalances"`
					PostBalances      []int64 `json:"postBalances"`
					InnerInstructions []struct {
						Index        int `json:"index"`
						Instructions []struct {
							Accounts       []int  `json:"accounts"`
							Data           string `json:"data"`
							ProgramIDIndex int    `json:"programIdIndex"`
							StackHeight    int    `json:"stackHeight"`
						} `json:"instructions"`
					} `json:"innerInstructions"`
					LogMessages      []string `json:"logMessages"`
					PreTokenBalances []struct {
						AccountIndex  int    `json:"accountIndex"`
						Mint          string `json:"mint"`
						Owner         string `json:"owner"`
						ProgramID     string `json:"programId"`
						UITokenAmount struct {
							Amount         string  `json:"amount"`
							Decimals       int     `json:"decimals"`
							UIAmount       float64 `json:"uiAmount"`
							UIAmountString string  `json:"uiAmountString"`
						} `json:"uiTokenAmount"`
					} `json:"preTokenBalances"`
					PostTokenBalances []struct {
						AccountIndex  int    `json:"accountIndex"`
						Mint          string `json:"mint"`
						Owner         string `json:"owner"`
						ProgramID     string `json:"programId"`
						UITokenAmount struct {
							Amount         string  `json:"amount"`
							Decimals       int     `json:"decimals"`
							UIAmount       float64 `json:"uiAmount"`
							UIAmountString string  `json:"uiAmountString"`
						} `json:"uiTokenAmount"`
					} `json:"postTokenBalances"`
					Rewards         any `json:"rewards"`
					LoadedAddresses struct {
						Writable []any `json:"writable"`
						Readonly []any `json:"readonly"`
					} `json:"loadedAddresses"`
					ComputeUnitsConsumed int `json:"computeUnitsConsumed"`
				} `json:"meta"`
			} `json:"transaction"`
			Signature string `json:"signature"`
		} `json:"result"`
	} `json:"params"`
}

type GetAssetResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Interface string `json:"interface"`
		ID        string `json:"id"`
		Content   struct {
			Schema  string `json:"$schema"`
			JSONURI string `json:"json_uri"`
			Files   []struct {
				URI    string `json:"uri"`
				CdnURI string `json:"cdn_uri"`
				Mime   string `json:"mime"`
			} `json:"files"`
			Metadata struct {
				Description   string `json:"description"`
				Name          string `json:"name"`
				Symbol        string `json:"symbol"`
				TokenStandard string `json:"token_standard"`
			} `json:"metadata"`
			Links struct {
				Image string `json:"image"`
			} `json:"links"`
		} `json:"content"`
		Authorities []struct {
			Address string   `json:"address"`
			Scopes  []string `json:"scopes"`
		} `json:"authorities"`
		Compression struct {
			Eligible    bool   `json:"eligible"`
			Compressed  bool   `json:"compressed"`
			DataHash    string `json:"data_hash"`
			CreatorHash string `json:"creator_hash"`
			AssetHash   string `json:"asset_hash"`
			Tree        string `json:"tree"`
			Seq         int    `json:"seq"`
			LeafID      int    `json:"leaf_id"`
		} `json:"compression"`
		Grouping []any `json:"grouping"`
		Royalty  struct {
			RoyaltyModel        string  `json:"royalty_model"`
			Target              any     `json:"target"`
			Percent             float64 `json:"percent"`
			BasisPoints         int     `json:"basis_points"`
			PrimarySaleHappened bool    `json:"primary_sale_happened"`
			Locked              bool    `json:"locked"`
		} `json:"royalty"`
		Creators  []any `json:"creators"`
		Ownership struct {
			Frozen         bool   `json:"frozen"`
			Delegated      bool   `json:"delegated"`
			Delegate       any    `json:"delegate"`
			OwnershipModel string `json:"ownership_model"`
			Owner          string `json:"owner"`
		} `json:"ownership"`
		Supply    any  `json:"supply"`
		Mutable   bool `json:"mutable"`
		Burnt     bool `json:"burnt"`
		TokenInfo struct {
			Supply       int64  `json:"supply"`
			Decimals     int    `json:"decimals"`
			TokenProgram string `json:"token_system"`
		} `json:"token_info"`
	} `json:"result"`
	ID string `json:"id"`
}
