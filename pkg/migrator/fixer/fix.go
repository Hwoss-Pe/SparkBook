package fixer

import (
	"Webook/pkg/migrator"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OverrideFixer[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB
	//这里只能支持同结构迁移
	columns []string
}

func NewOverrideFixer[T migrator.Entity](base *gorm.DB, target *gorm.DB) (*OverrideFixer[T], error) {
	//这里去查一下目的表
	var t T
	rows, err := target.Model(&t).Limit(1).Rows()
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	return &OverrideFixer[T]{
		base:    base,
		target:  target,
		columns: columns}, nil
}

// Fix 根据id进行搜寻修复数据 以源库为准的校验
func (o *OverrideFixer[T]) Fix(ctx context.Context, id int64) error {
	var src T
	err := o.base.WithContext(ctx).Where("id = ?", id).First(&src).Error
	switch err {
	case nil:
		//	找到数据就校验目的库,进行upsert
		return o.target.Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(o.columns),
		}).Create(&src).Error
	case gorm.ErrRecordNotFound:
		//以源库为准的校验
		return o.target.Delete("id =?", id).Error
	default:
		return err
	}
}
