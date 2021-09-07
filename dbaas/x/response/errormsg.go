package response

// 请求状态码
const (
	CodeSuccess = 0  // 请求成功
	CodeFail    = 1  // 请求失败
	CodeNoAuth  = -1 // 未登录 | 没有权限
)

var (
	EmptyMsg              = &errMsg{}
	ErrorParameter        = &errMsg{En: "Parameter error", Zh: "参数错误"}
	ErrorUpdate           = &errMsg{En: "update error", Zh: "修改失败"}
	ErrorImageOccupied    = &errMsg{En: "This image is occupied!", Zh: "此镜像被占用无法删除！"}
	ErrorImageExist       = &errMsg{En: "The image already exists!", Zh: "镜像已存在！"}
	ErrorBackupMax        = &errMsg{En: "the number of user backups exceeds the maximum!", Zh: "备份数量已达到上限！"}
	ErrorOperatorNoChange = &errMsg{En: "no image changes, no deployment required", Zh: "没有镜像更改, 无需部署"}
)
