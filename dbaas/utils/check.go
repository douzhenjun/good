package utils

/*
MustMap 检查map是否包含keys
*/
func MustMap(m map[string]interface{}, keys ...string) bool {
	for i := range keys {
		if _, ok := m[keys[i]]; !ok {
			return false
		}
	}
	return true
}

/*
MustInt 入参Int集合不能有0值
*/
func MustInt(values ...int) bool {
	for i := range values {
		if values[i] == 0 {
			return false
		}
	}
	return true
}

/*
StringLength 检测字符串长度在范围内
 */
func StringLength(s string, min, max int) bool {
	l := len(s)
	if l < min {
		return false
	}
	if l > max {
		return false
	}
	return true
}

/*
ReadMapString 安全读取map的值, 不存在或类型不匹配返回空值
*/
func ReadMapString(m map[string]interface{}, key string) string {
	return ReadMapStringDef(m, key, "")
}

func ReadMapStringDef(m map[string]interface{}, key string, def string) (vs string) {
	if v, ok := m[key]; ok {
		if vs, ok = v.(string); ok {
			return vs
		}
	}
	return def
}

func ReadMapInt(m map[string]interface{}, key string) int {
	return ReadMapIntDef(m, key, 0)
}

func ReadMapIntDef(m map[string]interface{}, key string, def int) int {
	if v, ok := m[key]; ok {
		if vi, ok := v.(int); ok {
			return vi
		}
		if vf, ok := v.(float64); ok {
			return int(vf)
		}
	}
	return def
}
