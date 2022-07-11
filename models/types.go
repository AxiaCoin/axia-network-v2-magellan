// (c) 2021, Axia Systems, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package models

type AXChainType uint16

var (
	AXChainIn     AXChainType = 1
	AXchainOut    AXChainType = 2
	AXChainImport AXChainType = 1
	AXChainExport AXChainType = 2

	OutputTypesSECP2556K1Transfer OutputType = 7
	OutputTypesSECP2556K1Mint     OutputType = 6
	OutputTypesNFTMint            OutputType = 10
	OutputTypesNFTTransfer        OutputType = 11
	OutputTypesAtomicExportTx     OutputType = 0xFFFFFFF1
	OutputTypesAtomicImportTx     OutputType = 0xFFFFFFF2

	BlockTypeProposal BlockType = 0x0
	BlockTypeAbort    BlockType = 0x1
	BlockTypeCommit   BlockType = 0x2
	BlockTypeStandard BlockType = 0x3
	BlockTypeAtomic   BlockType = 0x4

	TransactionTypeBase               TransactionType = 0x0
	TransactionTypeCreateAsset        TransactionType = 0x1
	TransactionTypeOperation          TransactionType = 0x2
	TransactionTypeAVMImport          TransactionType = 0x3
	TransactionTypeAVMExport          TransactionType = 0x4
	TransactionTypeAddValidator       TransactionType = 0xc
	TransactionTypeAddAllychainValidator TransactionType = 0xd
	TransactionTypeAddNominator       TransactionType = 0xe
	TransactionTypeCreateChain        TransactionType = 0xf
	TransactionTypeCreateAllychain       TransactionType = 0x10
	TransactionTypePVMImport          TransactionType = 0x11
	TransactionTypePVMExport          TransactionType = 0x12
	TransactionTypeAdvanceTime        TransactionType = 0x13
	TransactionTypeRewardValidator    TransactionType = 0x14

	ResultTypeTransaction SearchResultType = "transaction"
	ResultTypeAsset       SearchResultType = "asset"
	ResultTypeAddress     SearchResultType = "address"
	ResultTypeOutput      SearchResultType = "output"

	TypeUnknown = "unknown"
)

// BlockType represents a sub class of Block.
type BlockType uint16

// TransactionType represents a sub class of Transaction.
type TransactionType uint16

func (t TransactionType) String() string {
	switch t {
	case TransactionTypeBase:
		return "base"

		// AVM
	case TransactionTypeCreateAsset:
		return "create_asset"
	case TransactionTypeOperation:
		return "operation"
	case TransactionTypeAVMImport:
		return "import"
	case TransactionTypeAVMExport:
		return "export"

		// PVM
	case TransactionTypeAddValidator:
		return "add_validator"
	case TransactionTypeAddAllychainValidator:
		return "add_allychain_validator"
	case TransactionTypeAddNominator:
		return "add_nominator"
	case TransactionTypeCreateChain:
		return "create_chain"
	case TransactionTypeCreateAllychain:
		return "create_allychain"
	case TransactionTypePVMImport:
		return "pvm_import"
	case TransactionTypePVMExport:
		return "pvm_export"
	case TransactionTypeAdvanceTime:
		return "advance_time"
	case TransactionTypeRewardValidator:
		return "reward_validator"
	default:
		return TypeUnknown
	}
}

// OutputType represents a sub class of Output.
type OutputType uint32

func (t OutputType) String() string {
	switch t {
	case OutputTypesSECP2556K1Transfer:
		return "secp256k1_transfer"
	case OutputTypesSECP2556K1Mint:
		return "secp256k1_mint"
	case OutputTypesNFTTransfer:
		return "nft_transfer"
	case OutputTypesNFTMint:
		return "nft_mint"
	case OutputTypesAtomicExportTx:
		return "atomic_export"
	case OutputTypesAtomicImportTx:
		return "atomic_import"
	default:
		return TypeUnknown
	}
}

// SearchResultType is the type for an object found from a search query.
type SearchResultType string
