package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/osmosis-labs/osmosis/osmoutils/osmocli"
	"github.com/osmosis-labs/osmosis/v27/x/market/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := osmocli.TxIndexCmd(types.ModuleName)
	osmocli.AddTxCmd(cmd, NewSwapCmd)
	osmocli.AddTxCmd(cmd, CmdUpdateParams)

	return cmd
}

// NewSwapCmd will create and send a MsgSwap
func NewSwapCmd() (*osmocli.TxCliDesc, *types.MsgSwap) {
	return &osmocli.TxCliDesc{
		Use:     "swap",
		NumArgs: 3,
		//Args:  cobra.RangeArgs(2, 3),
		Short: "Atomically swap currencies at their target exchange rate",
		Long: strings.TrimSpace(`
   Swap the offer-coin to the ask-denom currency at the oracle's effective exchange rate.

   $ symphonyd market swap symphony1fr2x4cdvka7yfs8q9gqh0gzmh4hkmktpqwqj63 1000stake note

   The to-address can be specified. A default to-address is trader.

   $ symphonyd market swap "symphony1..." "1000stake" "note"
   `),
		ParseAndBuildMsg: NewSwapMsg,
	}, &types.MsgSwap{}
}

func NewSwapMsg(clientCtx client.Context, args []string, fs *flag.FlagSet) (sdk.Msg, error) {
	offerCoinStr := args[1]
	offerCoin, err := sdk.ParseCoinNormalized(offerCoinStr)
	if err != nil {
		return nil, err
	}

	askDenom := args[2]
	fromAddress := clientCtx.GetFromAddress()

	var msg sdk.Msg
	if len(args) == 3 {
		toAddress, err := sdk.AccAddressFromBech32(args[0])
		if err != nil {
			return nil, err
		}

		innerMsg := types.NewMsgSwapSend(fromAddress, toAddress, offerCoin, askDenom)
		if err = innerMsg.ValidateBasic(); err != nil {
			return nil, err
		}
		msg = innerMsg
	} else {
		innerMsg := types.NewMsgSwap(fromAddress, offerCoin, askDenom)
		if err = innerMsg.ValidateBasic(); err != nil {
			return nil, err
		}
		msg = innerMsg
	}
	return msg, nil
}

// CmdUpdateParams will create and send a MsgUpdateParams
func CmdUpdateParams() (*osmocli.TxCliDesc, *types.MsgUpdateParams) {
	return &osmocli.TxCliDesc{
		Use:     "update-params [tax_receiver]",
		NumArgs: 1,
		Short:   "Update module parameters",
		Long: strings.TrimSpace(`
   Update module parameters

   $ symphonyd tx market update-params symphonyaddrr...
   `),
		ParseAndBuildMsg: NewUpdateParamsMsg,
	}, &types.MsgUpdateParams{}
}

func NewUpdateParamsMsg(clientCtx client.Context, args []string, flags *flag.FlagSet) (sdk.Msg, error) {
	fromAddress := clientCtx.GetFromAddress()

	var msg sdk.Msg
	if len(args) == 1 {
		taxReceiverAddr, err := sdk.AccAddressFromBech32(args[0])
		if err != nil {
			return nil, err
		}

		innerMsg := types.NewMsgUpdateParams(fromAddress, taxReceiverAddr)
		if err = innerMsg.ValidateBasic(); err != nil {
			return nil, err
		}
		msg = innerMsg
	} else {
		return nil, fmt.Errorf("tax_receiver not found")
	}
	return msg, nil
}
