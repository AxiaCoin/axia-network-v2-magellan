// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"github.com/axiacoin/avalanchego/genesis"
	"github.com/axiacoin/avalanchego/ids"
	"github.com/axiacoin/avalanchego/utils/constants"
	"github.com/axiacoin/avalanchego/vms/platformvm"
)

type GenesisContainer struct {
	NetworkID       uint32
	XChainGenesisTx *platformvm.Tx
	XChainID        ids.ID
	AvaxAssetID     ids.ID
	GenesisBytes    []byte
}

func NewGenesisContainer(networkID uint32) (*GenesisContainer, error) {
	gc := &GenesisContainer{NetworkID: networkID}
	var err error
	gc.GenesisBytes, gc.AvaxAssetID, err = genesis.FromConfig(genesis.GetConfig(gc.NetworkID))
	if err != nil {
		return nil, err
	}

	gc.XChainGenesisTx, err = genesis.VMGenesis(gc.GenesisBytes, constants.AVMID)
	if err != nil {
		return nil, err
	}

	gc.XChainID = gc.XChainGenesisTx.ID()
	return gc, nil
}
