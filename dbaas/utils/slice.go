package utils

/*
通用的查询切片是否包含某元素, 传入fn指定匹配方法, fn入参是切片的下标
返回int为切片的下标
*/
func SliceExist(len int, fn func(int) bool) (bool, int) {
	for i := 0; i < len; i++ {
		if fn(i) {
			return true, i
		}
	}
	return false, -1
}

/*
查询int切片是否包含某元素, 返回int为匹配值
*/
func SliceExistInt(slice []int, target int) (exist bool, value int) {
	exist, index := SliceExist(len(slice), func(i int) bool { return slice[i] == target })
	if exist {
		value = slice[index]
	}
	return
}

/*
查询string切片是否包含某元素, 返回string为匹配值
*/
func SliceExistString(slice []string, target string) (exist bool, value string) {
	exist, index := SliceExist(len(slice), func(i int) bool { return slice[i] == target })
	if exist {
		value = slice[index]
	}
	return
}
