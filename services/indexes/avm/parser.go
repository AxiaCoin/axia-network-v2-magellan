// (c) 2021, AXIA Systems, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"github.com/axiacoin/axia-network-v2/codec"
	"github.com/axiacoin/axia-network-v2/vms/avm"
)

func parseTx(c codec.Manager, bytes []byte) (*avm.Tx, error) {
	tx := &avm.Tx{}
	ver, err := c.Unmarshal(bytes, tx)
	if err != nil {
		return nil, err
	}
	unsignedBytes, err := c.Marshal(ver, &tx.UnsignedTx)
	if err != nil {
		return nil, err
	}

	tx.Initialize(unsignedBytes, bytes)
	return tx, nil
}
