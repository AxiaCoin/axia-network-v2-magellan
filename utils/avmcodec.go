// (c) 2021, AXIA Systems, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"errors"

	"github.com/axiacoin/axia/codec"
	"github.com/axiacoin/axia/genesis"
	"github.com/axiacoin/axia/utils/constants"
	"github.com/axiacoin/axia/vms/avm"
	"github.com/axiacoin/axia/vms/nftfx"
	"github.com/axiacoin/axia/vms/platformvm"
	"github.com/axiacoin/axia/vms/secp256k1fx"
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
