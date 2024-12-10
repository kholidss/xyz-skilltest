package presentation

type (
	ReqTransactionCreditUser struct {
		MerchantID string `json:"merchant_id,omitempty"`
		AssetName  string `json:"asset_name,omitempty"`
		Tenor      int    `json:"tenor,omitempty"`
		OTRAmount  int    `json:"otr_amount,omitempty"`
	}

	RespTransactionCreditUser struct {
		Merchant struct {
			ID   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"merchant,omitempty"`
		AssetName            string `json:"asset_name,omitempty"`
		Tenor                int    `json:"tenor,omitempty"`
		OTRAmount            int    `json:"otr_amount"`
		FeeAmount            int    `json:"fee_amount"`
		InterestAmount       int    `json:"interest_amount"`
		InterestPercentage   int    `json:"interest_percentage"`
		InstallmentAmount    int    `json:"installment_amount"`
		RemainingLimitAmount int    `json:"remaining_limit_amount"`
	}
)
