package connpool

import (
	"context"
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"sync/atomic"
)

const (
	pattenDependBase   = "depend_base"
	pattenDependTarget = "depend_target"
	pattenOnlyBase     = "only_base"
	pattenOnlyTarget   = "only_target"
)

// sql语句执行会进这

type DoubleWritePool struct {
	base   gorm.ConnPool
	target gorm.ConnPool
	patten atomic.Value
}

// PrepareContext prepare进入此方法
func (d *DoubleWritePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	panic("不支持 ")
}

// ExecContext 非查询语句进入此方法
func (d *DoubleWritePool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.patten.Load().(string) {
	case pattenDependBase:
		res, err := d.base.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		_, err = d.target.ExecContext(ctx, query, args...)
		if err != nil {
			// log
		}
		return res, err

	case pattenDependTarget:
		res, err := d.target.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		_, err = d.base.ExecContext(ctx, query, args...)
		if err != nil {
			// log
		}
		return res, err

	case pattenOnlyBase:
		return d.base.ExecContext(ctx, query, args...)

	case pattenOnlyTarget:
		return d.target.ExecContext(ctx, query, args...)

	default:
		return nil, errors.New("patten error")

	}
}

// QueryContext 查询语句进入此方法
func (d *DoubleWritePool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.patten.Load().(string) {
	case pattenDependBase, pattenOnlyBase:
		return d.base.QueryContext(ctx, query, args...)

	case pattenDependTarget, pattenOnlyTarget:
		return d.target.QueryContext(ctx, query, args...)

	default:
		//return nil, errors.New("patten error")
		panic("patten error")

	}
}

func (d *DoubleWritePool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.patten.Load().(string) {
	case pattenDependBase, pattenOnlyBase:
		return d.base.QueryRowContext(ctx, query, args...)

	case pattenDependTarget, pattenOnlyTarget:
		return d.target.QueryRowContext(ctx, query, args...)

	default:
		// 构建不出错误
		panic("patten error")

	}
}

func (d *DoubleWritePool) BeginTx(ctx context.Context, opts *sql.TxOptions) (gorm.ConnPool, error) {
	switch d.patten.Load().(string) {
	case pattenOnlyBase:
		tx, err := d.base.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWritePoolTx{
			base:   tx,
			patten: pattenOnlyBase,
		}, err
	case pattenDependBase:
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
			patten: pattenDependBase,
		}, nil
	case pattenOnlyTarget:
		tx, err := d.target.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWritePoolTx{
			target: tx,
			patten: pattenOnlyTarget,
		}, err
	case pattenDependTarget:
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
			patten: pattenDependBase,
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
	case pattenOnlyBase:
		return d.base.Commit()
	case pattenDependBase:
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
	case pattenOnlyTarget:
		return d.target.Commit()
	case pattenDependTarget:
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
	case pattenOnlyBase:
		return d.base.Rollback()
	case pattenDependBase:
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
	case pattenOnlyTarget:
		return d.target.Rollback()
	case pattenDependTarget:
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
	case pattenDependBase:
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

	case pattenDependTarget:
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

	case pattenOnlyBase:
		return d.base.ExecContext(ctx, query, args...)

	case pattenOnlyTarget:
		return d.target.ExecContext(ctx, query, args...)

	default:
		return nil, errors.New("patten error")

	}
}

func (d *DoubleWritePoolTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.patten {
	case pattenDependBase, pattenOnlyBase:
		return d.base.QueryContext(ctx, query, args...)

	case pattenDependTarget, pattenOnlyTarget:
		return d.target.QueryContext(ctx, query, args...)

	default:
		//return nil, errors.New("patten error")
		panic("patten error")

	}
}

func (d *DoubleWritePoolTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.patten {
	case pattenDependBase, pattenOnlyBase:
		return d.base.QueryRowContext(ctx, query, args...)

	case pattenDependTarget, pattenOnlyTarget:
		return d.target.QueryRowContext(ctx, query, args...)

	default:
		// 构建不出错误
		panic("patten error")

	}
}
