// (c) 2021, AXIA Systems, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package api

import (
	"encoding/json"

	"github.com/axiacoin/axia-network-v2/ids"
	"github.com/axiacoin/axia-network-v2-magellan/models"
)

const (
	AVMName     = "avm"
	SwapChainAlias = "swap"
	PVMName     = "pvm"
	CoreChainAlias = "core"
)

func newIndexResponse(networkID uint32, swapChainID ids.ID, axcAssetID ids.ID) ([]byte, error) {
	return json.Marshal(&struct {
		NetworkID uint32                      `json:"network_id"`
		Chains    map[string]models.ChainInfo `json:"chains"`
	}{
		NetworkID: networkID,
		Chains: map[string]models.ChainInfo{
			swapChainID.String(): {
				VM:          AVMName,
				Alias:       SwapChainAlias,
				NetworkID:   networkID,
				AXCAssetID: models.StringID(axcAssetID.String()),
				ID:          models.StringID(swapChainID.String()),
			},
			ids.Empty.String(): {
				VM:          PVMName,
				Alias:       CoreChainAlias,
				NetworkID:   networkID,
				AXCAssetID: models.StringID(axcAssetID.String()),
				ID:          models.StringID(ids.Empty.String()),
			},
		},
	})
}

func newLegacyIndexResponse(networkID uint32, swapChainID ids.ID, axcAssetID ids.ID) ([]byte, error) {
	return json.Marshal(&models.ChainInfo{
		VM:          AVMName,
		NetworkID:   networkID,
		Alias:       SwapChainAlias,
		AXCAssetID: models.StringID(axcAssetID.String()),
		ID:          models.StringID(swapChainID.String()),
	})
}
