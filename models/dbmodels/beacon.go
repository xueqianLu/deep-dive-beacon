package dbmodels

import (
	"gorm.io/gorm"
)

type BeaconBlock struct {
	gorm.Model
	SlotNumber  uint64 `gorm:"uniqueIndex;not null" json:"slot_number"` // 槽位号
	EpochNumber uint64 `gorm:"index;not null" json:"epoch_number"`      // Epoch号

	// 验证者信息
	ProposerIndex uint64 `gorm:"not null" json:"proposer_index"`               // 验证者索引
	ParentRoot    string `gorm:"type:varchar(66);not null" json:"parent_root"` // 父区块根哈希
	StateRoot     string `gorm:"type:varchar(66);not null" json:"state_root"`  // 状态根哈希

	// RANDAO相关
	RandaoReveal string `gorm:"type:varchar(194);not null" json:"randao_reveal"` // RANDAO揭示

	// Graffiti
	Graffiti string `gorm:"type:varchar(66);not null" json:"graffiti"` // Graffiti数据

	// Eth1相关信息
	Eth1BlockHash    string `gorm:"type:varchar(66)" json:"eth1_block_hash"`   // Eth1区块哈希
	Eth1DepositRoot  string `gorm:"type:varchar(66)" json:"eth1_deposit_root"` // Eth1存款根
	Eth1DepositCount uint64 `json:"eth1_deposit_count"`                        // Eth1存款计数

	// 签名
	Signature string `gorm:"type:varchar(194);not null" json:"signature"` // 区块签名

	// Slashing信息
	ProposerSlashed uint `gorm:"default:0" json:"proposer_slashed"` // 提议者被slash数量
	AttesterSlashed uint `gorm:"default:0" json:"attester_slashed"` // 证明者被slash数量
}

type BeaconAttestation struct {
	gorm.Model
	SlotNumber      uint64 `gorm:"index;not null" json:"slot_number"`                  // 槽位号
	AttestIndex     int    `gorm:"not null" json:"attest_index"`                       // 在该slot中的证明索引
	AggregationBits string `gorm:"type:text;not null" json:"aggregation_bits"`         // 聚合位图
	BeaconBlockRoot string `gorm:"type:varchar(66);not null" json:"beacon_block_root"` // 关联的区块根哈希
	CommitteeIndex  uint64 `gorm:"not null" json:"committee_index"`
	SourceEpoch     uint64 `gorm:"not null" json:"source_epoch"`
	SourceRoot      string `gorm:"type:varchar(66);not null" json:"source_root"`
	TargetEpoch     uint64 `gorm:"not null" json:"target_epoch"`
	TargetRoot      string `gorm:"type:varchar(66);not null" json:"target_root"`
	Signature       string `gorm:"type:varchar(194);not null" json:"signature"`
}
