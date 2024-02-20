package testutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

type IndexInfo struct {
	Table        string
	NonUnique    int    `gorm:"column:Non_unique"`
	KeyName      string `gorm:"column:Key_name"`
	SeqInIndex   int    `gorm:"column:Seq_in_index"`
	ColumnName   string `gorm:"column:Column_name"`
	Collation    string
	Cardinality  int
	Null         string
	IndexType    string `gorm:"column:Index_type"`
	Comment      string
	IndexComment string `gorm:"column:Index_comment"`
}

func GetIndexInfo(t *testing.T, db *gorm.DB, tableName, indexName string) (IndexInfo, error) {
	// check location indexes
	var indexes []IndexInfo
	err := db.Raw(fmt.Sprintf("SHOW INDEXES FROM %s WHERE Key_name=?", tableName), indexName).Scan(&indexes).Error
	assert.NoError(t, err)
	assert.Len(t, indexes, 1)
	return indexes[0], nil
}
