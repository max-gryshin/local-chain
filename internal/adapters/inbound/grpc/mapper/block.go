package mapper

import (
	grpcPkg "local-chain/transport/gen/transport"

	"local-chain/internal/types"
)

type BlockMapper struct{}

func NewBlockMapper() *BlockMapper {
	return &BlockMapper{}
}

func (bm *BlockMapper) BlockToRpc(block *types.Block) *grpcPkg.Block {
	return &grpcPkg.Block{
		Timestamp:    block.Timestamp,
		PreviousHash: block.PrevHash,
		Hash:         block.Hash,
	}
}

func (bm *BlockMapper) BlocksToRpc(blocks types.Blocks) []*grpcPkg.Block {
	rpcBlocks := make([]*grpcPkg.Block, 0, len(blocks))
	for _, block := range blocks {
		rpcBlocks = append(rpcBlocks, bm.BlockToRpc(block))
	}
	return rpcBlocks
}
