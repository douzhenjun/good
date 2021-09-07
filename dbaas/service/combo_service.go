package service

import (
	"DBaas/models"
	"DBaas/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"strconv"
	"strings"
)

type ComboService interface {
	List(page, pageSize, userId, clusterId int, key string) ([]models.Combo, int, error)
	Add(combo models.Combo, userIdStr string) error
	Edit(combo models.Combo) error
	Delete(comboId int) error
	User(comboId int, userIdStr string) error

	TagList() ([]models.ComboTag, error)
	TagAdd(name string) error
	TagDelete(tagId int) error
}

type comboService struct {
	engine *xorm.Engine
}

func (cs *comboService) TagList() ([]models.ComboTag, error) {
	list := make([]models.ComboTag, 0)
	return list, cs.engine.OrderBy("preset, id").Find(&list)
}

func (cs *comboService) TagAdd(name string) error {
	if len(name) > 20 {
		return errors.New("name is too long")
	}
	_, err := cs.engine.Insert(&models.ComboTag{Name: name})
	return err
}

func (cs *comboService) TagDelete(tagId int) error {
	if tagId <= 0 {
		return errors.New("tag id must > 0")
	}
	_, err := cs.engine.Where("id = ?", tagId).And("preset = false").Delete(new(models.ComboTag))
	return err
}

func (cs *comboService) Delete(comboId int) error {
	if comboId <= 0 {
		return errors.New("combo id must > 0")
	}
	_, err := cs.engine.ID(comboId).Delete(new(models.Combo))
	if err != nil {
		return err
	}
	_, _ = cs.engine.Where("combo_id = ?", comboId).Delete(new(models.ComboUser))
	return nil
}

func (cs *comboService) Edit(combo models.Combo) error {
	_, err := matchSc(combo.Copy, 0, cs.engine)
	if err != nil {
		return err
	}
	_, err = cs.engine.ID(combo.Id).MustCols("write_iops", "read_iops", "write_bps", "read_bps").Update(&combo)
	return err
}

func (cs *comboService) User(comboId int, userIdStr string) error {
	if comboId <= 0 {
		return errors.New("combo id must > 0")
	}
	// userId为-1时表示ALL，为空时表示无租户，有租户时以逗号分割
	if userIdStr == "-1" {
		_, _ = cs.engine.Where("combo_id = ?", comboId).Delete(new(models.ComboUser))
	} else {
		dbUser := make([]models.ComboUser, 0)
		err := cs.engine.Where("combo_id = ?", comboId).Find(&dbUser)
		if err != nil {
			return err
		}
		dbUserM := map[int]int{}
		for i := range dbUser {
			dbUserM[dbUser[i].UserId] = i
		}
		userStrList := strings.Split(userIdStr, ",")
		insertList := make([]models.ComboUser, 0)
		for i := range userStrList {
			if len(userStrList[i]) == 0 {
				continue
			}
			var userId, err = strconv.Atoi(userStrList[i])
			if err != nil {
				utils.LoggerError(err)
				continue
			}
			if _, ok := dbUserM[userId]; ok {
				dbUserM[userId] = -1
				continue
			}
			insertList = append(insertList, models.ComboUser{ComboId: comboId, UserId: userId})
		}
		if len(insertList) > 0 {
			_, err := cs.engine.Insert(&insertList)
			if err != nil {
				return err
			}
		}
		for _, v := range dbUserM {
			if v != -1 {
				_, _ = cs.engine.ID(dbUser[v].Id).Delete(new(models.ComboUser))
			}
		}
	}
	_, err := cs.engine.ID(comboId).Cols("assign_all").Update(&models.Combo{AssignAll: userIdStr == "-1"})
	return err
}

func (cs *comboService) Add(combo models.Combo, userIdStr string) error {
	_, err := matchSc(combo.Copy, 0, cs.engine)
	if err != nil {
		return err
	}
	combo.AssignAll = userIdStr == "-1"
	_, err = cs.engine.Insert(&combo)
	if err != nil || combo.AssignAll {
		return err
	}
	userIds := strings.Split(userIdStr, ",")
	insert := make([]models.ComboUser, len(userIds))
	for i := range userIds {
		userId, err := strconv.Atoi(userIds[i])
		utils.LoggerError(err)
		insert[i] = models.ComboUser{ComboId: combo.Id, UserId: userId}
	}
	if len(insert) > 0 {
		_, _ = cs.engine.Insert(&insert)
	}
	return err
}

func (cs *comboService) List(page, pageSize, userId, clusterId int, key string) ([]models.Combo, int, error) {
	session := cs.engine.Where("name like ?", "%"+key+"%").Desc("id")
	if userId > 0 {
		session.
			Join("LEFT OUTER", "combo_user", "combo_user.combo_id = combo.id").
			And("combo.assign_all = true OR combo_user.user_id = ?", userId)
	}
	var list = make([]models.Combo, 0)
	count, err := pageFind(page, pageSize, &list, session, new(models.Combo))
	if err != nil {
		return nil, 0, err
	}
	// 实例Sc节点数, 当clusterId大于0时, 用于验证套餐的可用性（套餐节点数 <= 实例Sc节点数）
	var clusterScNodeCount int
	if clusterId > 0 {
		cluster := models.ClusterInstance{Id: clusterId}
		exist, err := cs.engine.Cols("sc_name").Get(&cluster)
		if !exist {
			return nil, 0, fmt.Errorf("not found cluster %v, error: %v", clusterId, err)
		}
		sc := models.Sc{Name: cluster.ScName}
		exist, err = cs.engine.Cols("sc_type", "node_num", "name", "id").Get(&sc)
		if !exist {
			return nil, 0, fmt.Errorf("not found sc %v, error: %v", clusterId, err)
		}
		sc.CheckNodeNum(cs.engine)
		clusterScNodeCount = sc.NodeNum
	}
	if userId > 0 {
		cpu, mem, storage := getUserSurplus(userId, cs.engine)
		for i := range list {
			list[i].Available = list[i].Cpu <= cpu && list[i].Mem <= mem && list[i].Storage <= storage
			if list[i].Available && clusterId > 0 {
				list[i].Available = list[i].Copy <= clusterScNodeCount
			}
		}
	} else {
		for i := range list {
			list[i].Available = true
			if list[i].AssignAll {
				list[i].UserList = utils.RawJson("-1")
			} else {
				var userIds = make([]models.User, 0)
				err = cs.engine.SQL("select u.id, u.user_name from \"user\" u inner join combo_user cu on u.id = cu.user_id where cu.combo_id = ?", list[i].Id).Find(&userIds)
				utils.LoggerError(err)
				s, _ := json.Marshal(userIds)
				list[i].UserList = utils.RawJson(utils.Bytes2str(s))
			}
		}
	}
	return list, int(count), err
}

func NewComboService(db *xorm.Engine) ComboService {
	return &comboService{
		engine: db,
	}
}
