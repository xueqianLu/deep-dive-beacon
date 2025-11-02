package dbmodels

import "time"

type ScanTask struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	TaskType   string `gorm:"type:varchar(50);index"`
	LastNumber uint64
	Enabled    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type DirectlyScanTask struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	TaskType   string `gorm:"type:varchar(50);index"`
	Start      uint64
	End        uint64
	LastNumber uint64
	Enabled    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
