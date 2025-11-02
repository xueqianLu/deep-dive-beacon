package types

type BlockScanTask struct {
	Start uint64 `json:"start"`
}

type DirectScanTask struct {
	Start uint64 `json:"start"`
	End   uint64 `json:"end"`
}
type DeployConfig struct {
	BlockScan  *BlockScanTask    `json:"block_scan"`
	DirectScan []*DirectScanTask `json:"direct_scan"`
}
