package redis

import (
	"errors"
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"go.uber.org/zap"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v20/ingest"
	"github.com/osmosis-labs/osmosis/v20/ingest/sqs/domain"
	"github.com/osmosis-labs/osmosis/v20/ingest/sqs/log"
	"github.com/osmosis-labs/osmosis/v20/ingest/sqs/pools/common"
	"github.com/osmosis-labs/osmosis/v20/x/concentrated-liquidity/client/queryproto"
	concentratedtypes "github.com/osmosis-labs/osmosis/v20/x/concentrated-liquidity/types"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v20/x/poolmanager/types"
)

// poolIngester is an ingester for pools.
// It implements ingest.Ingester.
// It reads all pools from the state and writes them to the pools repository.
// As part of that, it instruments each pool with chain native balances and
// OSMO based TVL.
// NOTE:
// - TVL is calculated using spot price. TODO: use TWAP (https://app.clickup.com/t/86a182835)
// - TVL does not account for token precision. TODO: use assetlist for pulling token precision data
// (https://app.clickup.com/t/86a18287v)
// - If error in TVL calculation, TVL is set to the value that could be computed and the pool struct
// has a flag to indicate that there was an error in TVL calculation.
type poolIngester struct {
	poolsRepository    domain.PoolsRepository
	routerRepository   domain.RouterRepository
	tokensUseCase      domain.TokensUsecase
	repositoryManager  domain.TxManager
	gammKeeper         common.PoolKeeper
	concentratedKeeper common.ConcentratedKeeper
	cosmWasmKeeper     common.CosmWasmPoolKeeper
	bankKeeper         common.BankKeeper
	protorevKeeper     common.ProtorevKeeper
	poolManagerKeeper  common.PoolManagerKeeper
	logger             log.Logger
}

// denomRoutingInfo encapsulates the routing information for a pool.
// It has a pool ID of the pool that is paired with OSMO.
// It has a spot price from that pool with OSMO as the base asset.
type denomRoutingInfo struct {
	PoolID uint64
	Price  osmomath.BigDec
}

const (
	UOSMO          = "uosmo"
	uosmoPrecision = 6

	noTokenPrecisionErrorFmtStr = "error getting token precision %s"
	spotPriceErrorFmtStr        = "error calculating spot price for denom %s, %s"

	noTotalValueLockedError = ""
)

var uosmoPrecisionBigDec = osmomath.NewBigDec(uosmoPrecision)

// NewPoolIngester returns a new pool ingester.
func NewPoolIngester(poolsRepository domain.PoolsRepository, routerRepository domain.RouterRepository, tokensUseCase domain.TokensUsecase, repositoryManager domain.TxManager, gammKeeper common.PoolKeeper, concentratedKeeper common.ConcentratedKeeper, cosmwasmKeeper common.CosmWasmPoolKeeper, bankKeeper common.BankKeeper, protorevKeeper common.ProtorevKeeper, poolManagerKeeper common.PoolManagerKeeper) ingest.AtomicIngester {
	return &poolIngester{
		poolsRepository:    poolsRepository,
		routerRepository:   routerRepository,
		tokensUseCase:      tokensUseCase,
		repositoryManager:  repositoryManager,
		gammKeeper:         gammKeeper,
		concentratedKeeper: concentratedKeeper,
		cosmWasmKeeper:     cosmwasmKeeper,
		bankKeeper:         bankKeeper,
		protorevKeeper:     protorevKeeper,
		poolManagerKeeper:  poolManagerKeeper,
	}
}

// ProcessBlock implements ingest.Ingester.
func (pi *poolIngester) ProcessBlock(ctx sdk.Context, tx domain.Tx) error {
	return pi.processPoolState(ctx, tx)
}

var _ ingest.AtomicIngester = &poolIngester{}

// processPoolState processes the pool state. an
func (pi *poolIngester) processPoolState(ctx sdk.Context, tx domain.Tx) error {
	goCtx := sdk.WrapSDKContext(ctx)

	// TODO: can be cached
	tokenPrecisionMap, err := pi.tokensUseCase.GetDenomPrecisions(goCtx)
	if err != nil {
		return err
	}

	// Create a map from denom to routable pool ID.
	denomToRoutablePoolIDMap := make(map[string]denomRoutingInfo)

	// CFMM pools

	cfmmPools, err := pi.gammKeeper.GetPools(ctx)
	if err != nil {
		return err
	}

	denomPairToTakerFeeMap := make(map[domain.DenomPair]osmomath.Dec, 0)

	// Parse CFMM pool to the standard SQS types.
	cfmmPoolsParsed := make([]domain.PoolI, 0, len(cfmmPools))
	for _, pool := range cfmmPools {

		// Parse CFMM pool to the standard SQS types.
		pool, err := pi.convertPool(ctx, pool, denomToRoutablePoolIDMap, denomPairToTakerFeeMap, tokenPrecisionMap)
		if err != nil {
			return err
		}

		cfmmPoolsParsed = append(cfmmPoolsParsed, pool)
	}

	// Concentrated pools

	concentratedPools, err := pi.concentratedKeeper.GetPools(ctx)
	if err != nil {
		return err
	}

	concentratedPoolsParsed := make([]domain.PoolI, 0, len(concentratedPools))
	for _, pool := range concentratedPools {
		// Parse concentrated pool to the standard SQS types.
		pool, err := pi.convertPool(ctx, pool, denomToRoutablePoolIDMap, denomPairToTakerFeeMap, tokenPrecisionMap)
		if err != nil {
			return err
		}

		concentratedPoolsParsed = append(concentratedPoolsParsed, pool)
	}

	// CosmWasm pools

	cosmWasmPools, err := pi.cosmWasmKeeper.GetPoolsWithWasmKeeper(ctx)
	if err != nil {
		return err
	}

	cosmWasmPoolsParsed := make([]domain.PoolI, 0, len(cosmWasmPools))
	for _, pool := range cosmWasmPools {
		// Parse cosmwasm pool to the standard SQS types.
		pool, err := pi.convertPool(ctx, pool, denomToRoutablePoolIDMap, denomPairToTakerFeeMap, tokenPrecisionMap)
		if err != nil {
			return err
		}

		cosmWasmPoolsParsed = append(cosmWasmPoolsParsed, pool)
	}

	pi.logger.Info("ingesting pools to Redis", zap.Int64("height", ctx.BlockHeight()), zap.Int("num_cfmm", len(cfmmPools)), zap.Int("num_concentrated", len(concentratedPools)), zap.Int("num_cosmwasm", len(cosmWasmPools)))

	err = pi.poolsRepository.StorePools(goCtx, tx, cfmmPoolsParsed, concentratedPoolsParsed, cosmWasmPoolsParsed)
	if err != nil {
		return err
	}

	// persist taker fees
	err = pi.persistTakerFees(ctx, tx, denomPairToTakerFeeMap)
	if err != nil {
		return err
	}

	return nil
}

// convertPool converts a pool to the standard SQS pool type.
// It instruments the pool with chain native balances and OSMO based TVL.
// If error occurs in TVL estimation, it is silently skipped and the error flag
// set to true in the pool model.
// Note:
// - TVL is calculated using spot price. TODO: use TWAP (https://app.clickup.com/t/86a182835)
// - TVL does not account for token precision. TODO: use assetlist for pulling token precision data
// (https://app.clickup.com/t/86a18287v)
func (pi *poolIngester) convertPool(
	ctx sdk.Context,
	pool poolmanagertypes.PoolI,
	denomToRoutingInfoMap map[string]denomRoutingInfo,
	denomPairToTakerFeeMap domain.TakerFeeMap,
	tokenPrecisionMap map[string]int,
) (domain.PoolI, error) {
	balances := pi.bankKeeper.GetAllBalances(ctx, pool.GetAddress())

	osmoPoolTVL := osmomath.ZeroInt()

	poolDenoms := pool.GetPoolDenoms(ctx)
	poolDenomsMap := map[string]struct{}{}

	// Convert pool denoms to map
	for _, poolDenom := range poolDenoms {
		poolDenomsMap[poolDenom] = struct{}{}
	}

	spreadFactor := pool.GetSpreadFactor(ctx)

	// Note that this must follow the call to GetPoolDenoms() and GetSpreadFactor.
	// Otherwise, the CosmWasmPool model panics.
	pool = pool.AsSerializablePool()

	var errorInTVLStr string
	for _, balance := range balances {

		// Note that there are edge cases where gamm shares or some random
		// garbage tokens are in the balance that do not belong to the pool.
		// A mainnet example is pool ID 2 with the following extra denoms:
		// ibc/65BCD5909ED3D9E6223529017BC828ECBECCBE3F63D444EC44CE7412EF8C82D6
		// ibc/778F0504E33BBB66D0950FE12E29BA81C258ED0A10CCEF9CB0096BA9E22C5D61
		// As a result, we skilently skip them
		// TODO: cover with test
		_, exists := poolDenomsMap[balance.Denom]
		if !exists {
			continue
		}

		if balance.Denom == UOSMO {
			osmoPoolTVL = osmoPoolTVL.Add(balance.Amount)
			continue
		}

		// Check if routable poolID already exists for the denom
		routingInfo, ok := denomToRoutingInfoMap[balance.Denom]
		if !ok {
			poolForDenomPair, err := pi.protorevKeeper.GetPoolForDenomPair(ctx, UOSMO, balance.Denom)
			if err != nil {
				pi.logger.Debug("error getting OSMO-based pool", zap.String("denom", balance.Denom), zap.Error(err))
				errorInTVLStr = err.Error()
				continue
			}

			basePrecison, ok := tokenPrecisionMap[balance.Denom]
			if !ok {
				errorInTVLStr = fmt.Sprintf(noTokenPrecisionErrorFmtStr, balance.Denom)
				pi.logger.Debug(errorInTVLStr)
				continue
			}

			uosmoBaseAssetSpotPrice, err := pi.poolManagerKeeper.RouteCalculateSpotPrice(ctx, poolForDenomPair, balance.Denom, UOSMO)
			if err != nil {
				errorInTVLStr = fmt.Sprintf(spotPriceErrorFmtStr, balance.Denom, err)
				pi.logger.Debug(errorInTVLStr)
				continue
			}

			// Scale on-chain spot price to the correct token precision.
			precisionMultiplier := osmomath.NewBigDec(int64(basePrecison)).Quo(uosmoPrecisionBigDec)

			uosmoBaseAssetSpotPrice = uosmoBaseAssetSpotPrice.Mul(precisionMultiplier)

			routingInfo = denomRoutingInfo{
				PoolID: poolForDenomPair,
				Price:  uosmoBaseAssetSpotPrice,
			}
		}

		osmoPoolTVL = osmoPoolTVL.Add(osmomath.NewBigDecFromBigInt(balance.Amount.BigInt()).MulMut(routingInfo.Price).Dec().TruncateInt())
	}

	// Get pool denoms. Although these can be inferred from balances, this is safer.
	// If we used balances, for pools with no liquidity, we would not be able to get the denoms.
	denoms, err := pi.poolManagerKeeper.RouteGetPoolDenoms(ctx, pool.GetId())
	if err != nil {
		return nil, err
	}

	// Sort denoms for consistent ordering.
	denoms = sort.StringSlice(denoms)

	// Mutates denomPairToTakerFeeMap with the taker fee for every uniquer denom pair in the denoms list.
	err = retrieveTakerFeeToMapIfNotExists(ctx, denoms, denomPairToTakerFeeMap, pi.poolManagerKeeper)
	if err != nil {
		return nil, err
	}

	// Get the tick model for concentrated pools
	var tickModel *domain.TickModel

	// For CL pools, get the tick data
	if pool.GetType() == poolmanagertypes.Concentrated {
		tickData, currentTickIndex, err := pi.concentratedKeeper.GetTickLiquidityForFullRange(ctx, pool.GetId())
		// If there is no error, we set the tick model
		if err == nil {
			tickModel = &domain.TickModel{
				Ticks:            tickData,
				CurrentTickIndex: currentTickIndex,
			}
			// If there is no liquidity, we set the tick model to nil and update no liquidity flag
		} else if err != nil && errors.Is(err, concentratedtypes.RanOutOfTicksForPoolError{PoolId: pool.GetId()}) {
			tickModel = &domain.TickModel{
				Ticks:            []queryproto.LiquidityDepthWithRange{},
				CurrentTickIndex: -1,
				HasNoLiquidity:   true,
			}

			// On any other error, we return the error
		} else {
			return nil, err
		}
	}

	return &domain.PoolWrapper{
		ChainModel: pool,
		SQSModel: domain.SQSPool{
			TotalValueLockedUSDC:  osmoPoolTVL,
			TotalValueLockedError: errorInTVLStr,
			Balances:              balances,
			PoolDenoms:            denoms,
			SpreadFactor:          spreadFactor,
		},
		TickModel: tickModel,
	}, nil
}

// persistTakerFees persists all taker fees to the router repository.
func (pi *poolIngester) persistTakerFees(ctx sdk.Context, tx domain.Tx, takerFeeMap domain.TakerFeeMap) error {
	for denomPair, takerFee := range takerFeeMap {
		err := pi.routerRepository.SetTakerFee(sdk.WrapSDKContext(ctx), tx, denomPair.Denom0, denomPair.Denom1, takerFee)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetLogger implements ingest.AtomicIngester.
func (pi *poolIngester) SetLogger(logger log.Logger) {
	pi.logger = logger
}