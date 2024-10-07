package connpool

import (
	"context"
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"sync/atomic"
)

const (
	PattenDependBase   = "depend_base"
	PattenDependTarget = "depend_target"
	PattenOnlyBase     = "only_base"
	PattenOnlyTarget   = "only_target"
)

// sql语句执行会进这

type DoubleWritePool struct {
	base   gorm.ConnPool
	target gorm.ConnPool
	patten atomic.Value
}

func NewDoubleWritePool(base, target gorm.ConnPool) *DoubleWritePool {
	var p atomic.Value
	p.Store(PattenDependBase)
	return &DoubleWritePool{
		base:   base,
		target: target,
		patten: p,
	}
}

func (d *DoubleWritePool) UpdatePatten(patten string) {
	d.patten.Store(patten)
}

// PrepareContext prepare进入此方法
func (d *DoubleWritePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	panic("不支持")
}

// ExecContext 非查询语句进入此方法
func (d *DoubleWritePool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.patten.Load().(string) {
	case PattenDependBase:
		res, err := d.base.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		_, err = d.target.ExecContext(ctx, query, args...)
		if err != nil {
			// log
		}
		return res, err

	case PattenDependTarget:
		res, err := d.target.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		_, err = d.base.ExecContext(ctx, query, args...)
		if err != nil {
			// log
		}
		return res, err

	case PattenOnlyBase:
		return d.base.ExecContext(ctx, query, args...)

	case PattenOnlyTarget:
		return d.target.ExecContext(ctx, query, args...)

	default:
		return nil, errors.New("patten error")

	}
}

// QueryContext 查询语句进入此方法
func (d *DoubleWritePool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.patten.Load().(string) {
	case PattenDependBase, PattenOnlyBase:
		return d.base.QueryContext(ctx, query, args...)

	case PattenDependTarget, PattenOnlyTarget:
		return d.target.QueryContext(ctx, query, args...)

	default:
		//return nil, errors.New("patten error")
		panic("patten error")

	}
}

func (d *DoubleWritePool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.patten.Load().(string) {
	case PattenDependBase, PattenOnlyBase:
		return d.base.QueryRowContext(ctx, query, args...)

	case PattenDependTarget, PattenOnlyTarget:
		return d.target.QueryRowContext(ctx, query, args...)

	default:
		// 构建不出错误
		panic("patten error")

	}
}

func (d *DoubleWritePool) BeginTx(ctx context.Context, opts *sql.TxOptions) (gorm.ConnPool, error) {
	switch d.patten.Load().(string) {
	case PattenOnlyBase:
		tx, err := d.base.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWritePoolTx{
			base:   tx,
			patten: PattenOnlyBase,
		}, err
	case PattenDependBase:
		baseTx, err := d.base.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		targetTx, err := d.target.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			// log
		}

		// 这里可能赋值给nil，即事务开启失败
		return &DoubleWritePoolTx{
			base:   baseTx,
			target: targetTx,
			patten: PattenDependBase,
		}, nil
	case PattenOnlyTarget:
		tx, err := d.target.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWritePoolTx{
			target: tx,
			patten: PattenOnlyTarget,
		}, err
	case PattenDependTarget:
		targetTx, err := d.target.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		baseTx, err := d.base.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			// log
		}

		// 这里可能赋值给nil，即事务开启失败
		return &DoubleWritePoolTx{
			base:   baseTx,
			target: targetTx,
			patten: PattenDependBase,
		}, nil
	default:
		return nil, errors.New("patten error")
	}
}

type DoubleWritePoolTx struct {
	base   *sql.Tx
	target *sql.Tx
	patten string // 无并发问题，直接string
}

func (d *DoubleWritePoolTx) Commit() error {
	switch d.patten {
	case PattenOnlyBase:
		return d.base.Commit()
	case PattenDependBase:
		err := d.base.Commit()
		if err != nil { // 事务提交失败，直接返回
			return err
		}
		if d.target != nil {
			err = d.target.Commit()
			if err != nil {
				// 这里挂了就只能记录日志了
			}
		}

		return nil
	case PattenOnlyTarget:
		return d.target.Commit()
	case PattenDependTarget:
		err := d.target.Commit()
		if err != nil {
			return err
		}
		if d.base != nil {
			err = d.base.Commit()
			if err != nil {
			}
		}

		return nil
	default:
		return errors.New("patten error")
	}
}

func (d *DoubleWritePoolTx) Rollback() error {
	switch d.patten {
	case PattenOnlyBase:
		return d.base.Rollback()
	case PattenDependBase:
		err := d.base.Rollback()
		if err != nil {
			return err
		}
		if d.target != nil {
			err = d.target.Rollback()
			if err != nil {
			}
		}

		return nil
	case PattenOnlyTarget:
		return d.target.Rollback()
	case PattenDependTarget:
		err := d.target.Rollback()
		if err != nil {
			return err
		}
		if d.base != nil {
			err = d.base.Rollback()
			if err != nil {
			}
		}

		return nil
	default:
		return errors.New("patten error")
	}
}

func (d *DoubleWritePoolTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	panic("implement me")
}

func (d *DoubleWritePoolTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.patten {
	case PattenDependBase:
		res, err := d.base.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		if d.target == nil {
			return res, err
		}

		_, err = d.target.ExecContext(ctx, query, args...)
		if err != nil {
			// log
		}
		return res, err

	case PattenDependTarget:
		res, err := d.target.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		if d.base == nil {
			return res, err
		}

		_, err = d.base.ExecContext(ctx, query, args...)
		if err != nil {
			// log
		}
		return res, err

	case PattenOnlyBase:
		return d.base.ExecContext(ctx, query, args...)

	case PattenOnlyTarget:
		return d.target.ExecContext(ctx, query, args...)

	default:
		return nil, errors.New("patten error")

	}
}

func (d *DoubleWritePoolTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.patten {
	case PattenDependBase, PattenOnlyBase:
		return d.base.QueryContext(ctx, query, args...)

	case PattenDependTarget, PattenOnlyTarget:
		return d.target.QueryContext(ctx, query, args...)

	default:
		//return nil, errors.New("patten error")
		panic("patten error")

	}
}

func (d *DoubleWritePoolTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.patten {
	case PattenDependBase, PattenOnlyBase:
		return d.base.QueryRowContext(ctx, query, args...)

	case PattenDependTarget, PattenOnlyTarget:
		return d.target.QueryRowContext(ctx, query, args...)

	default:
		// 构建不出错误
		panic("patten error")

	}
}
