// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"errors"

	"github.com/axiacoin/avalanchego/codec"
	"github.com/axiacoin/avalanchego/genesis"
	"github.com/axiacoin/avalanchego/utils/constants"
	"github.com/axiacoin/avalanchego/vms/avm"
	"github.com/axiacoin/avalanchego/vms/nftfx"
	"github.com/axiacoin/avalanchego/vms/platformvm"
	"github.com/axiacoin/avalanchego/vms/secp256k1fx"
)

var (
	ErrIncorrectGenesisChainTxType = errors.New("incorrect genesis chain tx type")
)

const MaxCodecSize = 100_000_000

func NewAVMCodec(networkID uint32) (codec.Manager, error) {
	genesisBytes, _, err := genesis.FromConfig(genesis.GetConfig(networkID))
	if err != nil {
		return nil, nil
	}

	g, err := genesis.VMGenesis(genesisBytes, constants.AVMID)
	if err != nil {
		return nil, err
	}

	createChainTx, ok := g.UnsignedTx.(*platformvm.UnsignedCreateChainTx)
	if !ok {
		return nil, ErrIncorrectGenesisChainTxType
	}

	var (
		fxIDs = createChainTx.FxIDs
		fxs   = make([]avm.Fx, 0, len(fxIDs))
	)

	for _, fxID := range fxIDs {
		switch {
		case fxID == secp256k1fx.ID:
			fxs = append(fxs, &secp256k1fx.Fx{})
		case fxID == nftfx.ID:
			fxs = append(fxs, &nftfx.Fx{})
		default:
			// return nil, fmt.Errorf("Unknown FxID: %s", fxID)
		}
	}

	_, codec, err := avm.NewCodecs(fxs)
	codec.SetMaxSize(MaxCodecSize)
	return codec, err
}
