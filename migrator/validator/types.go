package validator

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"post/migrator"
	"post/migrator/events"
	"reflect"
	"time"
)

// Validator T需要实现Entity接口
type Validator[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB

	p         events.InconsistentProducer
	direction string
}

func NewValidator[T migrator.Entity](base, target *gorm.DB, p events.InconsistentProducer, direction string) *Validator[T] {
	// 可在此处开一个goroutine监控负载情况
	return &Validator[T]{
		base:      base,
		target:    target,
		p:         p,
		direction: direction,
	}
}

// Validate 仅适用同构数据库
// todo 注意，此处校验完后，会存在多余数据，即base硬删除了数据，但是target没发现
// todo 可以采用慢启动的方式，对比count的数量，不一致再遍历找到多余的数据
func (v *Validator[T]) Validate(ctx context.Context, utime int64) {
	//utime := time.Now().UnixMilli() // 需要外部传入，即开始同步的时间
	for offset := 0; ; offset++ {
		// base, 实现了Entity的struct
		var base T
		ctxSon, cancel := context.WithTimeout(ctx, time.Second)
		//err := v.base.WithContext(ctxSon).Offset(offset).Order("id").First(&base).Error
		err := v.base.Where("utime < ?", utime).Offset(offset).First(&base).Error
		cancel()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return
			}
			//panic(err)
			// 监控，数据库不正常错误
			continue
		}

		// target, 实现了Entity的struct
		var target T
		ctxSon, cancel = context.WithTimeout(ctx, time.Second)
		err = v.target.WithContext(ctxSon).Where("id = ?", base.GetID()).First(&target).Error
		cancel()
		if err != nil {
			// target 缺失数据
			if errors.Is(err, gorm.ErrRecordNotFound) {
				v.sendEvent(ctx, base.GetID(), events.InconsistentEventTypeNoExist)
			}

			// todo
			// 监控，数据库不正常错误
		}

		//// 简单直接的版本
		//if !base.CompareWith(target) {
		//
		//}

		// 查看base是否实现CompareWith接口
		var baseAny any = base // 因为断言需要interface{}类型
		if b, ok := baseAny.(interface {
			CompareWith(e migrator.Entity) bool
		}); ok {
			if !b.CompareWith(target) {
				v.sendEvent(ctx, base.GetID(), events.InconsistentEventTypeNoEquals)
			}
		} else {
			if !reflect.DeepEqual(base, target) {

			}
		}
	}
}

func (v *Validator[T]) sendEvent(ctx context.Context, id int64, eventType string) {
	event := &events.InconsistentEvent{
		Direction: v.direction,
		ID:        id,
		Type:      eventType,
	}
	ctxSon, cancel := context.WithTimeout(ctx, time.Second)
	err := v.p.InconsistentEventProducer(ctxSon, event)
	cancel()
	if err != nil {
		// kafka发送失败，监控此处，只能人工处理
	}
}
