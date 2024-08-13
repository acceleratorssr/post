package fixer

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"post/migrator"
	"post/migrator/events"
)

type Fixer[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB
}

func NewFixerDependBase[T migrator.Entity](base, target *gorm.DB) *Fixer[T] {
	return &Fixer[T]{
		base:   base,
		target: target,
	}
}

func NewFixerDependTarget[T migrator.Entity](base, target *gorm.DB) *Fixer[T] {
	return &Fixer[T]{
		base:   target,
		target: base,
	}
}

// Fix 也可改批量
func (f *Fixer[T]) Fix(ctx context.Context, event events.InconsistentEvent) error {
	switch event.Type {
	case events.InconsistentEventTypeNoEquals:
		var t T
		err := f.base.WithContext(ctx).Where("id = ?", event.ID).First(&t).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return f.target.WithContext(ctx).Delete(&t).Error
			}
			return err
		}

		return f.target.WithContext(ctx).Updates(&t).Error

	case events.InconsistentEventTypeNoExist:
		var t T
		err := f.base.WithContext(ctx).Where("id = ?", event.ID).First(&t).Error
		if err != nil {
			// 注意此处的err也可能是没找到数据，即又硬删除了
			return err
		}
		return f.target.WithContext(ctx).Create(&t).Error
	}
	return errors.New("unknown event type")
}
