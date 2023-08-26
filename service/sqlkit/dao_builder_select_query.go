package sqlkit

import (
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/spf13/cast"
)

// 链式查询的

func (dao SelectDao[T]) QueryRows() *sqlx.Rows {
	return dao.QueryRaw(dao.Sql())
}

func (dao SelectDao[T]) One() *T {
	d := dao
	// 取未删除的
	if !dao.ignoreLogicDel {
		d = dao.whereNLogicDel()
	}
	rows := dao.QueryRaw(d.Sql())
	defer rows.Close()
	for rows.Next() {
		m := new(T)
		err := rows.StructScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		if d.Cascade != nil {
			d.Cascade(m)
		}
		return m
	}
	return nil
}

func (dao SelectDao[T]) List() []*T {
	d := dao
	if !dao.ignoreLogicDel {
		d = dao.whereNLogicDel()
	}
	rows := dao.QueryRaw(d.Sql())
	list := make([]*T, 0, 5)
	defer rows.Close()
	for rows.Next() {
		m := new(T)
		err := rows.StructScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	if dao.Cascade != nil {
		for i := range list {
			dao.Cascade(list[i])
		}
	}
	return list
}

func (dao SelectDao[T]) OneMap() map[string]any {
	d := dao
	if !dao.ignoreLogicDel {
		d = dao.whereNLogicDel()
	}
	rows := dao.QueryRaw(d.Sql())
	defer rows.Close()
	for rows.Next() {
		m := map[string]any{}
		err := rows.MapScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		return m
	}
	return nil
}

func (dao SelectDao[T]) OneString() string {
	d := dao
	if !dao.ignoreLogicDel {
		d = dao.whereNLogicDel()
	}
	rows := dao.QueryRaw(d.Sql())
	defer rows.Close()
	for rows.Next() {
		ret, err := rows.SliceScan()
		if err != nil {
			panic(exception.New(err.Error()))
		}
		return cast.ToString(ret[0])
	}
	return ""
}

func (dao SelectDao[T]) ListMap() []map[string]any {
	d := dao
	if !dao.ignoreLogicDel {
		d = dao.whereNLogicDel()
	}
	rows := dao.QueryRaw(d.Sql())
	defer rows.Close()
	list := make([]map[string]any, 0, 5)
	for rows.Next() {
		m := map[string]any{}
		err := rows.MapScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	return list
}

func (dao SelectDao[T]) ListString() []string {
	d := dao
	if !dao.ignoreLogicDel {
		d = dao.whereNLogicDel()
	}
	rows := dao.QueryRaw(d.Sql())
	defer rows.Close()
	list := make([]string, 0, 5)
	defer rows.Close()
	for rows.Next() {
		ret, err := rows.SliceScan()
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, ret[0].(string))
	}
	return list
}

// todo
//func (dao SelectDao[T]) Page()
