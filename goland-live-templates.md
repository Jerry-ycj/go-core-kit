
## action-init
```
func Init(router *router.Router) {
	tag := "$tname$"
	r := router.Group("/rest/$tag$")
	r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Post("/$name$", $name$).Tag(tag).Summary("$summary$").Param($name$Params{})
	}
}

type $name$Params struct{
    Phone    string `description:"手机号" default:"" trim:"true"`
	Pwd      string `validate:"required"`
}
func $name$(ctx *context.Context){
    params := $name$Params{}
	ctx.BindForm(&params)
	
    ctx.JsonSuccess(nil)
}
```

## action
```
type $name$Params struct{
    Phone    string `description:"手机号" default:"" trim:"true"`
	Pwd      string `validate:"required"`
}
func $name$(ctx *context.Context){
    params := $name$Params{}
	ctx.BindForm(&params)
	
    ctx.JsonSuccess(nil)
}
```

## bean-extend
```
func (th *$struct$) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	th.Id = cast.ToInt32(value)
	return nil
}
func (th *$struct$) Value() (driver.Value, error) {
	if th == nil {
		return nil, nil
	}
	return th.Id, nil
}
```

## dao
```
/// auto template
type Dao struct {
	sqlkit.Dao
}
const (
	ResultDefault byte = iota
	ResultNone
)
func New(schema string, tx ...*sqlkit.Dao) *Dao{
	dao:=&Dao{}
	dao.NewHelper(schema,tx...)
	return dao
}
func (dao *Dao) cascade(obj *$bean$) {
	switch dao.ResultType {
	case ResultDefault:
		// todo 注意校验nil
		// todo 如果没有级联，此函数删除
	}
}
func (dao *Dao) scan(sql string, args []interface{}) []$bean$ {
	rows := dao.Query(sql, args...)
	list := make([]$bean$,0,5)
	defer rows.Close()
	for rows.Next() {
		m := $bean${}
		err := rows.StructScan(&m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	for i := range list{
		dao.cascade(&list[i])
	}
	return list
}
func (dao *Dao) scanOne(sql string, args []interface{}) *$bean$ {
	rows := dao.Query(sql, args...)
	for rows.Next() {
		m := $bean${}
		err := rows.StructScan(&m)
		rows.Close()
		if err != nil {
			panic(exception.New(err.Error()))
		}
		dao.cascade(&m)
		return &m
	}
	return nil
}
////
```

## dao-demo
```
func (dao *Dao) FindById(id int32) *$bean$ {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("$name$")).Where("id=?",id).MustSql()
	return dao.scanOne(sql, args)
}

func (dao *Dao) ListAll() []$bean$ {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("$name$")).OrderBy("id").MustSql()
	return dao.scan(sql, args)
}
```

## exception
```
panic(exception.New("$msg$"))
```