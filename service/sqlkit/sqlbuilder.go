package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
)

type SQLBuilder struct {
	inner     squirrel.StatementBuilderType
	ModelMeta *ModelMeta
}

// Select 默认取modelmeta中的columns，并装饰引号；fields中不装饰，因为可能存在表达式
func (b SQLBuilder) _select(fields ...string) SelectBuilder {
	return SelectBuilder{
		inner:     b.inner.Select(fields...),
		ModelMeta: b.ModelMeta,
	}
}

func (b SQLBuilder) Select(fields ...string) SelectBuilder {
	if len(fields) == 0 {
		return b.SelectWithout()
	} else {
		return b._select(fields...)
	}
}

// SelectWithout 在modelmeta columns中去掉指定的字段
func (b SQLBuilder) SelectWithout(fields ...string) SelectBuilder {
	return b.SelectPrefix("", fields...)
}

// SelectPrefix 在modelmeta的字段前增加prefix
func (b SQLBuilder) SelectPrefix(prefix string, without ...string) SelectBuilder {
	if b.ModelMeta == nil {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	return b._select(b.ModelMeta.GetColumnsWithPrefix(prefix, without...)...)
}

func (b SQLBuilder) Update() UpdateBuilder {
	if b.ModelMeta == nil {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	return UpdateBuilder{
		inner:     b.inner.Update(b.ModelMeta.GetTableName()),
		ModelMeta: b.ModelMeta,
	}
}

func (b SQLBuilder) Delete() DeleteBuilder {
	if b.ModelMeta == nil {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	return DeleteBuilder{
		inner:     b.inner.Delete(b.ModelMeta.GetTableName()),
		ModelMeta: b.ModelMeta,
	}
}

func (b SQLBuilder) Insert() InsertBuilder {
	if b.ModelMeta == nil {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	return InsertBuilder{
		inner:     b.inner.Insert(b.ModelMeta.GetTableName()),
		ModelMeta: b.ModelMeta,
	}
}

func (b SQLBuilder) Replace() InsertBuilder {
	if b.ModelMeta == nil {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	return InsertBuilder{
		inner:     b.inner.Replace(b.ModelMeta.GetTableName()),
		ModelMeta: b.ModelMeta,
	}
}
