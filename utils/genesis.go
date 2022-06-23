// (c) 2021, AXIA Systems, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"github.com/axiacoin/axia-network-v2/genesis"
	"github.com/axiacoin/axia-network-v2/ids"
	"github.com/axiacoin/axia-network-v2/utils/constants"
	"github.com/axiacoin/axia-network-v2/vms/platformvm"
)

type GenesisContainer struct {
	NetworkID       uint32
	SwapChainGenesisTx *platformvm.Tx
	SwapChainID        ids.ID
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

	gc.SwapChainGenesisTx, err = genesis.VMGenesis(gc.GenesisBytes, constants.AVMID)
	if err != nil {
		return nil, err
	}

	gc.SwapChainID = gc.SwapChainGenesisTx.ID()
	return gc, nil
}
