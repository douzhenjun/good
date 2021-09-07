/**
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author:  zhangwei
 * @Date: 2020/11/16 10:30
 * @LastEditors: zhangwei
 * @LastEditTime: 2020/11/16 10:30
 **/
package service

import (
	"DBaas/utils"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
	"strconv"

	//"github.com/kataras/iris/v12"
	"DBaas/models"
	//"DBaas/utils"
)

//  用户服务接口定义
type UserService interface {
	SaveUser(user *models.User) bool
	SelectOne(id int) (models.User, string)
	//SelectOneByIp(ip string) (models.Host, bool)
	SelectOneByName(name string) (models.User, bool)
	DeleteUser(userId int) bool
	DeleteUserAndUserinst(userId int, operUsername string) (bool, string, string)
	ModifyUser(user models.User, id int, updateRemark bool) error
	ListUser(limit int, offset int, key string) ([]models.User, error)
	ListUserAll(key string) ([]models.User, error)
	GetUserCount(key string) (int64, error)
	GetStorageCount(userId int) (int64, error)
	GetClusterInstanceCount(userId int) (int64, error)
	GetClusterCpuMemStorage(userId int) (int64, int64, int64, error)
	GetClustersByUser(userId int) ([]models.ClusterInstance, error)
	GetPvBySc(scId int) ([]models.PersistentVolume, error)
	GetAllStorageCount() (int64, error)
	GetUseBackup(userId int) (int, error)
	SelectIdByTag(userTag string) (int, error)
}

//  创建主机服务的接口
func NewUserService(db *xorm.Engine) UserService {
	return &userService{
		Engine: db,
	}
}

//  主机服务结构体
type userService struct {
	Engine *xorm.Engine
}

func (us *userService) SelectIdByTag(userTag string) (int, error) {
	u := models.User{UserTag: userTag}
	exist, err := us.Engine.Cols("id").Get(&u)
	if !exist {
		return 0, fmt.Errorf("not found user by tag: %v, error: %v", userTag, err)
	}
	return u.Id, err
}

func (us *userService) GetUseBackup(userId int) (int, error) {
	return GetUseBackup(userId, us.Engine)
}

//  保存用户信息  新增用户
func (us *userService) SaveUser(user *models.User) bool {
	_, err := us.Engine.Insert(user)
	utils.LoggerError(err)
	return err == nil
}

func (us *userService) DeleteUser(userId int) bool {
	_, err := us.Engine.Where(" id = ? ", userId).Delete(new(models.User))
	_, _ = us.Engine.Where("user_id = ?", userId).Delete(new(models.ScUser))
	utils.LoggerError(err)
	return err == nil
}

func (us *userService) DeleteUserAndUserinst(userId int, operUsername string) (bool, string, string) {
	capricornService, conn := NewCapricornService()
	defer CloseGrpc(conn)
	oldUser, errM := us.SelectOne(userId)
	if errM != "" {
		return true, errM, errM
	}
	userInstIdString := strconv.Itoa(oldUser.ZdcpId)
	_, errorMsgEn, errorMsgZh := capricornService.DeleteUserResources(userInstIdString, operUsername)
	if errorMsgEn != "" && errorMsgZh != "" {
		return false, errorMsgEn, errorMsgZh
	}
	isSuccess := us.DeleteUser(userId)
	if !isSuccess {
		return false, "Failed to delete user information in dbaas", "删除DBaaS中用户信息失败"
	}
	return true, "", ""
}

func (us *userService) ModifyUser(user models.User, id int, updateRemark bool) error {
	session := us.Engine.Id(id)
	if updateRemark {
		session.MustCols("remark")
	}
	_, err := session.Update(&user)
	return err
}

func (us *userService) SelectOne(id int) (models.User, string) {
	var user models.User
	_, err := us.Engine.Where(" id = ? ", id).Get(&user)
	if err != nil {
		utils.LoggerError(err)
		return user, err.Error()
	}
	return user, ""
}

func (us *userService) SelectOneByName(name string) (models.User, bool) {
	var user models.User
	_, err := us.Engine.Where(" user_name = ? ", name).Get(&user)
	if err != nil {
		capricornService, conn := NewCapricornService()
		defer CloseGrpc(conn)
		userIdString := strconv.Itoa(user.ZdcpId)

		hostInfo := make(map[string]interface{}, 0)
		userList, ErrorMsgEn, ErrorMsgZh := capricornService.GetUserResources(userIdString, "", "")

		if ErrorMsgEn != "" && ErrorMsgZh != "" {
			iris.New().Logger().Error(ErrorMsgEn)
		}
		if len(userList) > 0 {
			hostInfo = userList[0]
			if _, ok := hostInfo["username"]; ok {
				if user.UserName != hostInfo["username"] {
					updateUser := models.User{UserName: hostInfo["username"].(string)}
					err = us.ModifyUser(updateUser, user.Id, false)
					if err != nil {
						utils.LoggerInfo("同步capricorn模块用户信息")
					}
				}
			}
		}
		utils.LoggerError(err)
	}
	return user, err == nil
}

//
//func (cs *userService) SelectOneByIp(ip string) (models.Host, bool) {
//	var host models.Host
//	_, err := cs.Engine.Where(" ip = ? ", ip).Get(&host)
//	if err != nil {
//		utils.LoggerError(err)
//	}
//	return host, err == nil
//}
//
func (us *userService) ListUser(limit int, offset int, key string) ([]models.User, error) {
	userList := make([]models.User, 0)
	err := us.Engine.Where(" user_name like ? ", "%"+key+"%").Or(" remarks like ? ", "%"+key+"%").Limit(limit, offset).OrderBy("-id").Find(&userList)
	utils.LoggerError(err)
	return userList, err
}

//
func (us *userService) ListUserAll(key string) ([]models.User, error) {
	userList := make([]models.User, 0)
	err := us.Engine.Where(" user_name like ? ", "%"+key+"%").Or(" remarks like ? ", "%"+key+"%").OrderBy("-id").Find(&userList)
	utils.LoggerError(err)
	return userList, err
}

/**
 * 获取用户总数量
 */
func (us *userService) GetUserCount(key string) (int64, error) {
	count, err := us.Engine.Where(" user_name like ? ", "%"+key+"%").Or(" remarks like ? ", "%"+key+"%").Count(new(models.User))
	if err != nil {
		return 0, err
	}
	return count, nil
}

/**
 * 获取用户下sc总数量/所有sc数量
 */
func (us *userService) GetStorageCount(userId int) (int64, error) {
	var count int64
	var err error
	if userId != 0 {
		count, err = us.Engine.Where(" user_id = ? ", userId).Count(new(models.ScUser))
	} else {
		count, err = us.Engine.Count(new(models.ScUser))
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

/**
 * 获取sc总数量
 */
func (us *userService) GetAllStorageCount() (int64, error) {
	count, err := us.Engine.Count(new(models.Sc))
	if err != nil {
		return 0, err
	}
	return count, nil
}

/**
 * 获取用户下ClusterInstance总数量/所有ClusterInstance数量
 */
func (us *userService) GetClusterInstanceCount(userId int) (int64, error) {
	var count int64
	var err error
	if userId != 0 {
		count, err = us.Engine.Where(" user_id = ? ", userId).Count(new(models.ClusterInstance))
	} else {
		count, err = us.Engine.Count(new(models.ClusterInstance))
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

/**
 * 获取用户下ClusterInstance已用CPU MEM
 */
func (us *userService) GetClusterCpuMemStorage(userId int) (int64, int64, int64, error) {
	cpuTotal, err := us.Engine.Where(" user_id = ? ", userId).SumInt(new(models.ClusterInstance), "limit_cpu")
	memTotal, err := us.Engine.Where(" user_id = ? ", userId).SumInt(new(models.ClusterInstance), "limit_mem")
	storageTotal, err := us.Engine.Where(" user_id = ? ", userId).SumInt(new(models.ClusterInstance), "storage")
	if err != nil {
		utils.LoggerError(err)
		return 0, 0, 0, err
	}
	return cpuTotal, memTotal, storageTotal, nil
}

/**
 * 获取sc下pv
 */
func (us *userService) GetPvBySc(scId int) ([]models.PersistentVolume, error) {

	pvList := make([]models.PersistentVolume, 0)
	err := us.Engine.Where(" sc_id = ? ", scId).Find(&pvList)
	utils.LoggerError(err)

	return pvList, nil
}

/**
 * 获取用户下cluster
 */
func (us *userService) GetClustersByUser(userId int) ([]models.ClusterInstance, error) {
	clList := make([]models.ClusterInstance, 0)
	err := us.Engine.Where("user_id = ?", userId).Omit("yaml_text").Find(&clList)
	return clList, err
}
