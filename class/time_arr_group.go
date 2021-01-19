package class

import (
	"github.com/spf13/cast"
	"sort"
	"time"
)

//
// 两个时间一组的数组
//
type TimeArrGroup [][]time.Time

type TimePoint struct {
	Time time.Time
	Id   string
	Type int32 // 1,2; 11,12-剔除的开始结束点
}
type TimePointList []*TimePoint

func (l TimePointList) Len() int           { return len(l) }
func (l TimePointList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l TimePointList) Less(i, j int) bool { return l[i].Time.Unix() < l[j].Time.Unix() }

// 累计时间长，s
func (th TimeArrGroup) Sum() int64 {
	var all int64 = 0
	for _, t := range th {
		if len(t) != 2 {
			continue
		}
		dif := t[1].Unix() - t[0].Unix()
		if dif <= 0 {
			continue
		}
		all += dif
	}
	return all
}

// typeFlag=10表示是剔除
func _group2TimeArr(a, b TimeArrGroup, typeFlag int32) TimePointList {
	list := make(TimePointList, 0, len(a)*2+len(b)*2)
	for i, e := range a {
		if len(e) != 2 {
			continue
		}
		list = append(list, &TimePoint{
			Time: e[0],
			Id:   "a" + cast.ToString(i),
			Type: 1,
		}, &TimePoint{
			Time: e[1],
			Id:   "a" + cast.ToString(i),
			Type: 2,
		})
	}
	for i, e := range b {
		if len(e) != 2 {
			continue
		}
		list = append(list, &TimePoint{
			Time: e[0],
			Id:   "b" + cast.ToString(i),
			Type: 1 + typeFlag,
		}, &TimePoint{
			Time: e[1],
			Id:   "b" + cast.ToString(i),
			Type: 2 + typeFlag,
		})
	}
	sort.Sort(list)
	return list
}

// 合并，当前的时间数组和参数的时间数组合并。
func (th TimeArrGroup) Merge(obj TimeArrGroup) TimeArrGroup {
	list := _group2TimeArr(th, obj, 0)
	ret := TimeArrGroup{}
	temp := map[string]bool{}
	var startTemp time.Time
	for _, e := range list {
		if startTemp.IsZero() && e.Type != 1 {
			continue
		}
		if e.Type > 2 || e.Type < 1 {
			continue
		}
		if startTemp.IsZero() {
			startTemp = e.Time
			temp[e.Id] = true
			continue
		}
		switch e.Type {
		case 1:
			temp[e.Id] = true
			if startTemp.IsZero() {
				startTemp = e.Time
			}
		case 2:
			if _, ok := temp[e.Id]; ok {
				delete(temp, e.Id)
				if len(temp) == 0 {
					// 封包
					ret = append(ret, []time.Time{startTemp, e.Time})
					startTemp = time.Time{}
				}
			}
		}
	}
	return ret
}

// 剔除，当前的时间数组剔除参数的时间数组范围
func (th TimeArrGroup) Eliminate(obj TimeArrGroup) TimeArrGroup {
	list := _group2TimeArr(th, obj, 10)
	ret := TimeArrGroup{}
	// type=1/2的临时存放，合并项或是被截断的项
	temp1 := map[string]bool{}
	// type=11/12 的临时存放
	temp2 := map[string]bool{}
	var startTemp time.Time
	for i, e := range list {
		if i == 0 && e.Type != 1 && e.Type != 11 {
			continue
		}
		switch e.Type {
		case 1:
			temp1[e.Id] = true
			if len(temp2) == 0 && startTemp.IsZero() {
				startTemp = e.Time
			}
		case 11:
			temp2[e.Id] = true
			// 被截断的情况
			if len(temp1) > 0 {
				ret = append(ret, []time.Time{startTemp, e.Time})
			}
			startTemp = time.Time{}
		case 2:
			if _, ok := temp1[e.Id]; ok {
				delete(temp1, e.Id)
				if len(temp1) == 0 && len(temp2) == 0 {
					// 封包
					ret = append(ret, []time.Time{startTemp, e.Time})
					startTemp = time.Time{}
				}
			}
		case 12:
			if _, ok := temp2[e.Id]; ok {
				delete(temp2, e.Id)
				if len(temp1) > 0 {
					// 存在截断
					startTemp = e.Time
				}
			}
		}
	}
	return ret
}
