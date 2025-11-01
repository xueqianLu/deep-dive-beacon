package dbmodels

import "time"

type Eth1BlockHeader struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	BlockHash  string `gorm:"type:varchar(66);uniqueIndex"`
	ParentHash string `gorm:"type:varchar(66);index"`
	Number     uint64 `gorm:"index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Eth1BlockBody struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Eth1Transaction struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Eth1Receipt struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Eth1Contract struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Eth1ContractTransaction struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Eth1Log struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
