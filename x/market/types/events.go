package types

// Market module event types
const (
	EventSwap = "swap"

	AttributeKeyOffer     = "offer"
	AttributeKeyTrader    = "trader"
	AttributeKeyRecipient = "recipient"
	AttributeKeySwapCoin  = "swap_coin"
	AttributeKeySwapFee   = "swap_fee"

	AttributeValueCategory = ModuleName

	EventUpdateParams = "update_params"

	AttributePreviousTaxReceiver = "previous_tax_receiver"
	AttributeNewTaxReceiver      = "new_tax_receiver"
)
