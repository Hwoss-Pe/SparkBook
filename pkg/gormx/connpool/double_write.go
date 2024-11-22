package connpool

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"gorm.io/gorm"
)

// TxBeginner实现了ConnPool
// 这个是不停机迁移的一些操作
// 正常不停机迁移一般是需要有四个状态，并且以主库和从库分别为主的设计，保证数据一致性
var errUnknownPattern = errors.New("未知的双写模式")

const (
	PatternSrcOnly  = "src_only"
	PatternSrcFirst = "src_first"
	PatternDstFirst = "dst_first"
	PatternDstOnly  = "dst_only"
)

// DoubleWritePool 对应的启动连接池的ConnPoolBeginner
type DoubleWritePool struct {
	//用的是开源的一个拓展库
	pattern *atomicx.Value[string]
	src     gorm.ConnPool
	dst     gorm.ConnPool
}

// PrepareContext 这个做预处理编译sql的时候才会使用，一般不会去显式调用
func (d *DoubleWritePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	panic(errUnknownPattern)
}

func (d *DoubleWritePool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	//执行更新上下文时候
	switch d.pattern.Load() {
	case PatternSrcOnly:
		return d.src.ExecContext(ctx, query, args)
	case PatternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err == nil {
			_, err1 := d.dst.ExecContext(ctx, query, args...)
			if err1 != nil {
				// 这边要记录日志
				//并且要通知修复数据
			}
		}
		return res, err
	case PatternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err == nil {
			_, err1 := d.src.ExecContext(ctx, query, args...)
			if err1 != nil {
				// 这边要记录日志
				// 并且要通知修复数据
			}
		}
		return res, err
	case PatternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	default:
		return nil, errUnknownPattern
	}
}

// QueryContext 下面的查询都是一样，查询来说实际上前两种状态都是操作一个数据库
func (d *DoubleWritePool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern.Load() {
	case PatternSrcFirst, PatternSrcOnly:
		return d.src.QueryContext(ctx, query, args...)
	case PatternDstFirst, PatternDstOnly:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWritePool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern.Load() {
	case PatternSrcFirst, PatternSrcOnly:
		return d.src.QueryRowContext(ctx, query, args...)
	case PatternDstFirst, PatternDstOnly:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		// 因为返回值里面咩有 error，只能 panic 掉
		panic(errUnknownPattern)
	}
}

// BeginTx 使用事务的时候会自动调用，或者手动调用
func (d *DoubleWritePool) BeginTx(ctx context.Context, opts *sql.TxOptions) (gorm.ConnPool, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly:
		//	只有一个源库初始状态
		tx, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriteTx{
			pattern: pattern,
			src:     tx,
		}, err
	case PatternSrcFirst:
		return d.startTwoTx(d.src, d.dst, pattern, ctx, opts)
	case PatternDstFirst:
		return d.startTwoTx(d.dst, d.src, pattern, ctx, opts)
	case PatternDstOnly:
		tx, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriteTx{
			pattern: pattern,
			src:     tx,
		}, err
	default:
		return nil, errUnknownPattern
	}
}

func NewDoubleWritePool(srcDB *gorm.DB, dst *gorm.DB) *DoubleWritePool {
	return &DoubleWritePool{
		src:     srcDB.ConnPool,
		dst:     dst.ConnPool,
		pattern: atomicx.NewValueOf(PatternSrcOnly)}
}

// ChangePattern 支持动态修改双写模式
func (d *DoubleWritePool) ChangePattern(pattern string) {
	d.pattern.Store(pattern)
}

func (d *DoubleWritePool) startTwoTx(first gorm.ConnPool, second gorm.ConnPool,
	pattern string, ctx context.Context, opts *sql.TxOptions) (*DoubleWriteTx, error) {
	//这里就要开始思考先进来的主，后的是从，并且事务控制双写
	src, err := first.(gorm.TxBeginner).BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	dst, err := second.(gorm.TxBeginner).BeginTx(ctx, opts)
	//某个双写失败应该让源库回滚
	if err != nil {
		_ = src.Rollback()
		//记录日志
	}
	return &DoubleWriteTx{src: src, dst: dst, pattern: pattern}, nil
}

// DoubleWriteTx 去实现底层的gorm.ConnPool,这种可以手动控制事务的提交和回滚，但是容易忘记不推荐
type DoubleWriteTx struct {
	pattern string
	src     *sql.Tx
	dst     *sql.Tx
}

// 实现Tx.Committer 手动控制它的commit和rollback，比如可以去支配具体哪个数据库进行回滚提交
// 策略：两个库都会提交，但是报错只返回主库，也就是说不关心目的库
// 需要显式调用tx事务的提交/回滚

func (d *DoubleWriteTx) Commit() error {
	switch d.pattern {
	case PatternSrcFirst:

		err := d.src.Commit()
		if d.dst != nil {
			err1 := d.dst.Commit()
			if err1 != nil {
				// 记录日志
			}
		}
		return err
	case PatternSrcOnly:
		return d.src.Commit()
	case PatternDstFirst:
		err := d.dst.Commit()
		if d.src != nil {
			err1 := d.src.Commit()
			if err1 != nil {
				// 记录日志
			}
		}

		return err
	case PatternDstOnly:
		return d.dst.Commit()
	default:
		return errUnknownPattern
	}
}

func (d *DoubleWriteTx) Rollback() error {
	switch d.pattern {
	case PatternSrcFirst:
		err := d.src.Rollback()
		if d.dst != nil {
			err1 := d.dst.Rollback()
			if err1 != nil {
				// 记录日志
			}
		}
		return err
	case PatternSrcOnly:
		return d.src.Rollback()
	case PatternDstOnly:
		return d.dst.Rollback()
	case PatternDstFirst:
		err := d.dst.Rollback()
		if d.src != nil {
			err1 := d.src.Rollback()
			if err1 != nil {
				// 记录日志
			}
		}
		return err
	default:
		return errUnknownPattern
	}
}

func (d *DoubleWriteTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	switch d.pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.PrepareContext(ctx, query)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.PrepareContext(ctx, query)
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWriteTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case PatternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err == nil {
			_, err1 := d.dst.ExecContext(ctx, query, args...)
			if err1 != nil {
				// 这边要记录日志
				// 并且要通知修复数据
			}
		}
		return res, err
	case PatternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err == nil {
			_, err1 := d.src.ExecContext(ctx, query, args...)
			if err1 != nil {
				// 这边要记录日志
				// 并且要通知修复数据
			}
		}
		return res, err
	case PatternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWriteTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern {
	case PatternSrcFirst, PatternSrcOnly:
		return d.src.QueryContext(ctx, query, args...)
	case PatternDstFirst, PatternDstOnly:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWriteTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern {
	case PatternSrcFirst, PatternSrcOnly:
		return d.src.QueryRowContext(ctx, query, args...)
	case PatternDstFirst, PatternDstOnly:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		// 因为返回值里面咩有 error，只能 panic 掉
		panic(errUnknownPattern)
	}
}
