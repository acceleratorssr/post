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

	utime int64
}

func NewValidator[T migrator.Entity](base, target *gorm.DB, p events.InconsistentProducer, direction string, utime int64) *Validator[T] {
	// 可在此处开一个goroutine监控负载情况
	return &Validator[T]{
		base:      base,
		target:    target,
		p:         p,
		direction: direction,
		utime:     utime,
	}
}

// Validate
// todo 仅适用同构数据库，可改造offset并行
//
//	注意，此处校验完后，会存在多余数据，即base硬删除了数据，但是target没发现
//	可以采用慢启动的方式，对比count的数量，不一致再遍历找到多余的数据
//
// todo 此处的sql可以优化下，可以将两库都按照时间顺序排序，然后对比，均批量取数据
// 在主从的情况下，校验可以通过从库进行，发现不一致再查主库
func (v *Validator[T]) Validate(ctx context.Context, msTimeout int64, limit int) {
	//utime := time.Now().UnixMilli() // 需要外部传入，即开始同步的时间
	// base, 实现了Entity的struct
	base := make([]T, 0, limit)

	var target T
	var null T
	// 查看T类型是否实现 CompareWith 接口
	// 另一种方法是直接将CompareWith作为方法写入Entity接口中，因为T类型一定要实现Entity接口
	var targetAny any = target // 因为断言需要interface{}类型
	var fn func(index int)

	if t, ok := targetAny.(interface {
		CompareWith(e migrator.Entity) bool
	}); ok {
		fn = func(index int) {
			if !t.CompareWith(target) {
				v.sendEvent(ctx, base[index].GetID(), events.InconsistentEventTypeNoEquals)
			}
		}
	} else {
		fn = func(index int) {
			if !reflect.DeepEqual(base, target) {
				v.sendEvent(ctx, base[index].GetID(), events.InconsistentEventTypeNoEquals)
			}
		}
	}

	for offset := 0; ; offset += len(base) {
		//base = make([]T, 0, limit)
		base = base[:0] // 清空切片，但保留容量不变

		ctxSon, cancel := context.WithTimeout(ctx, time.Duration(msTimeout))
		//err := v.base.WithContext(ctxSon).Offset(offset).Order("id").First(&base).Error
		err := v.base.Where("utime < ?", v.utime).Offset(offset).Limit(limit).Find(&base).Error // utime记得加索引
		cancel()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return //校验完成
			}
			if errors.Is(err, context.DeadlineExceeded) {
				// log
				// 监控，超时了
			}
			if errors.Is(err, context.Canceled) {
				return
			}
			// 监控，数据库不正常错误
			continue
		}

		for i := 0; i < len(base); i++ {
			// target, 实现了Entity的struct
			ctxSon, cancel = context.WithTimeout(ctx, 200*time.Millisecond)
			err = v.target.WithContext(ctxSon).Where("id = ?", base[i].GetID()).First(&target).Error
			cancel()
			if err != nil {
				// target 缺失数据
				if errors.Is(err, gorm.ErrRecordNotFound) {
					v.sendEvent(ctx, base[i].GetID(), events.InconsistentEventTypeNoExist)
				}
				if errors.Is(err, context.DeadlineExceeded) {
					// log
					// 监控，超时了
				}
				if errors.Is(err, context.Canceled) {
					return
				}

				// todo
				// 监控，数据库不正常错误
			}

			fn(i)
			target = null
		}

		// 校验完成后改为调stop结束
		//if limit > len(base) {
		//	return
		//}
		if len(base) == 0 {
			offset = 0
		}

		// 也可以通过switch：context.Canceled
		select {
		case <-ctx.Done():
			return
		default:

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
