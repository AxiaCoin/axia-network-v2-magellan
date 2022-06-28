package params

import (
	"testing"

	"github.com/axiacoin/axia-network-v2/ids"
	"github.com/axiacoin/axia-network-v2/utils/hashing"
)

func TestForValueChainID(t *testing.T) {
	res := ForValueChainID(nil, nil)
	if res != nil {
		t.Error("ForValueChainID failed")
	}
	temcoreChain, _ := ids.ToID(hashing.ComputeHash256([]byte("tid1")))
	res = ForValueChainID(&temcoreChain, nil)
	if len(res) != 1 || res[0] != temcoreChain.String() {
		t.Error("ForValueChainID failed")
	}
	res = ForValueChainID(&temcoreChain, []string{})
	if len(res) != 1 || res[0] != temcoreChain.String() {
		t.Error("ForValueChainID failed")
	}
	res = ForValueChainID(&temcoreChain, []string{temcoreChain.String()})
	if len(res) != 1 || res[0] != temcoreChain.String() {
		t.Error("ForValueChainID failed")
	}
	temcoreChain2, _ := ids.ToID(hashing.ComputeHash256([]byte("tid2")))
	if temcoreChain.String() == temcoreChain2.String() {
		t.Error("toId failed")
	}
	res = ForValueChainID(&temcoreChain, []string{temcoreChain2.String()})
	if len(res) != 2 || res[0] != temcoreChain.String() || res[1] != temcoreChain2.String() {
		t.Error("ForValueChainID failed")
	}
}
