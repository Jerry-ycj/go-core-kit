package privilegedao

import (
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

type Dao struct {
	sqlkit.Dao[model.PrivilegeConstant]
}

func New(ds ...*sqlkit.DataSource) Dao {
	return Dao{sqlkit.New[model.PrivilegeConstant](ds...)}
}

func (dao Dao) ListPrivileges() []*model.PrivilegeConstant {
	builder := dao.Builder().Select().OrderBy("sort")
	return dao.QueryList(builder)
}
