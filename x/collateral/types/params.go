package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"strings"
)

// Parameter keys
var (
	KeyAllowedAssets          = []byte("AllowedAssets")
	KeyStableDenom            = []byte("StableDenom")
	KeyWithdrawalsEnabled     = []byte("WithdrawalsEnabled")
	KeyPriceAdjustmentEnabled = []byte("PriceAdjustmentEnabled")
	KeyMinCollateralRatio     = []byte("MinCollateralRatio")
	KeyMaxUserAssets          = []byte("MaxUserAssets")
	KeyTrackTokenPrefix       = []byte("TrackTokenPrefix")
	KeyOracleModule           = []byte("OracleModule")
)

var _ paramstypes.ParamSet = &Params{}

func DefaultParams() Params {
	return Params{
		AllowedAssets: []*AllowedAsset{
			{Denom: "wBTC"},
			{Denom: "wETC"},
			{Denom: "wSOL"},
		},
		StableDenom:            "USDC",
		WithdrawalsEnabled:     true,
		PriceAdjustmentEnabled: false,
		MinCollateralRatio:     "0.5",
		MaxUserAssets:          "3",
		TrackTokenPrefix:       "t",
	}
}

func (p Params) Validate() error {
	return nil
}

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() types.ParamSetPairs {
	return types.ParamSetPairs{
		types.NewParamSetPair(KeyAllowedAssets, &p.AllowedAssets, validateSupportedTokens),
		types.NewParamSetPair(KeyStableDenom, &p.StableDenom, validateStableToken),
		types.NewParamSetPair(KeyWithdrawalsEnabled, &p.WithdrawalsEnabled, validatePriceOrWithdrawalEnabled),
		types.NewParamSetPair(KeyPriceAdjustmentEnabled, &p.PriceAdjustmentEnabled, validatePriceOrWithdrawalEnabled),
		types.NewParamSetPair(KeyMinCollateralRatio, &p.MinCollateralRatio, validateRatio),
		types.NewParamSetPair(KeyMaxUserAssets, &p.MaxUserAssets, validateMaxUserAssets),
		types.NewParamSetPair(KeyTrackTokenPrefix, &p.TrackTokenPrefix, validateKeyTrackTokenPrefix),
		types.NewParamSetPair(KeyOracleModule, &p.XOracleAddr, validateKeyOracleAddr),
	}
}

func validateKeyOracleAddr(i interface{}) error {
	addr, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	_, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return fmt.Errorf("invalid oracle address: %w", err)
	}

	return nil
}

func validateKeyTrackTokenPrefix(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 || len(v) > 3 {
		return fmt.Errorf("TrackTokenPrefix should has 1 or 3 length, got %d", len(v))
	}

	return nil
}

func validateMaxUserAssets(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	maxAssets, ok := osmomath.NewIntFromString(v)
	if !ok {
		return fmt.Errorf("invalid MaxUserAssets value: %s", v)
	}

	if !maxAssets.IsPositive() {
		return fmt.Errorf("MaxUserAssets must be positive")
	}

	return nil
}

func validatePriceOrWithdrawalEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateRatio(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	rewardRatio, err := osmomath.NewDecFromStr(v)
	if err != nil {
		return fmt.Errorf("invalid ratio: %v", err)
	}

	if !rewardRatio.IsPositive() {
		return fmt.Errorf("ration must be positive")
	}

	return nil
}

func IsAllowedToken(token string, allowedTokens []*AllowedAsset) bool {
	for _, t := range allowedTokens {
		if t.Denom == token {
			return true
		}
	}
	return false
}

func validateStableToken(i interface{}) error {
	t, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(t) == "" {
		return fmt.Errorf("supported token cannot be blank")
	}
	if err := sdk.ValidateDenom(t); err != nil {
		return err
	}

	return nil
}

func validateSupportedTokens(i interface{}) error {
	tokens, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(tokens) == 0 {
		return fmt.Errorf("supported tokens cannot be empty")
	}

	for _, t := range tokens {
		if strings.TrimSpace(t) == "" {
			return fmt.Errorf("supported token cannot be blank")
		}
		if err := sdk.ValidateDenom(t); err != nil {
			return err
		}
	}

	return nil
}
