package types

type BlockScanTask struct {
	Start uint64 `json:"start"`
}
type DeployConfig struct {
	BlockScan *BlockScanTask `json:"block_scan"`
}
