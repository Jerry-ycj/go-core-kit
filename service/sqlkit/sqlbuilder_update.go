package sqlkit

import (
	"github.com/Masterminds/squirrel"
)

type UpdateBuilder struct {
	inner     squirrel.UpdateBuilder
	modelMeta ModelMeta
}

func (b UpdateBuilder) Sql() (string, []interface{}) {
	return b.inner.MustSql()
}
func (b UpdateBuilder) ToSql() (string, []interface{}, error) {
	return b.inner.ToSql()
}

// SQL methods

// Prefix 在 sql 前写入语句
func (b UpdateBuilder) Prefix(sql string, args ...interface{}) UpdateBuilder {
	b.inner = b.inner.Prefix(sql, args...)
	return b
}
func (b UpdateBuilder) PrefixExpr(expr UpdateBuilder) UpdateBuilder {
	b.inner = b.inner.PrefixExpr(expr)
	return b
}
func (b UpdateBuilder) Suffix(sql string, args ...interface{}) UpdateBuilder {
	b.inner = b.inner.Suffix(sql, args...)
	return b
}
func (b UpdateBuilder) SuffixExpr(expr UpdateBuilder) UpdateBuilder {
	b.inner = b.inner.SuffixExpr(expr)
	return b
}

func (b UpdateBuilder) Set(column string, value interface{}) UpdateBuilder {
	b.inner = b.inner.Set(b.modelMeta.escapeName(column), value)
	return b
}
func (b UpdateBuilder) Where(pred interface{}, args ...interface{}) UpdateBuilder {
	b.inner = b.inner.Where(pred, args...)
	return b
}
func (b UpdateBuilder) FromSelect(from SelectBuilder, alias string) UpdateBuilder {
	b.inner = b.inner.FromSelect(from.inner, alias)
	return b
}