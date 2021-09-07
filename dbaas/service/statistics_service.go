package service

import (
	"DBaas/models"
	"fmt"
	"github.com/go-xorm/xorm"
	"time"
)

func InitStatistics(engine *xorm.Engine) {
	statistics = &statisticsService{engine: engine}
	Statistics = statistics
}

var Statistics StatisticsService
var statistics *statisticsService

type StatisticsService interface {
	ClusterList(page, pageSize int, search string, t string, replicas int, sort string, desc bool) ([]models.StatisticsCluster, map[string]int, error)
}

type statisticsService struct {
	engine *xorm.Engine
}

func clusterKey(name string, id int) string {
	return fmt.Sprintf("%v-%v", name, id)
}

func (ss *statisticsService) clusterExist(name string, id int) (models.StatisticsCluster, bool, error) {
	m := models.StatisticsCluster{Cluster: clusterKey(name, id)}
	exist, err := ss.engine.Exist(&m)
	return m, exist, err
}

func (ss *statisticsService) ClusterDeploy(name string, id, replicas int, from string) {
	m, exist, err := ss.clusterExist(name, id)
	if err != nil || exist {
		return
	}
	m.DeployStart = time.Now()
	m.Replicas = replicas
	m.Type = from
	_, _ = ss.engine.Insert(&m)
}

func (ss *statisticsService) ClusterComplete(name string, id int) {
	m := models.StatisticsCluster{Cluster: clusterKey(name, id)}
	exist, err := ss.engine.Cols("deploy_end").Get(&m)
	if err != nil || !exist || !m.DeployEnd.IsZero() {
		return
	}
	m.DeployEnd = time.Now()
	_, _ = ss.engine.Where("cluster = ?", m.Cluster).Cols("deploy_end").Update(&m)
}

func (ss *statisticsService) ClusterTimeout(name string, id int) {
	m, exist, err := ss.clusterExist(name, id)
	if err != nil || !exist {
		return
	}
	m.DeployTimeout = true
	_, _ = ss.engine.Where("cluster = ?", m.Cluster).Cols("deploy_timeout").Update(&m)
}

func (ss *statisticsService) ClusterDelete(name string, id int) {
	m, exist, err := ss.clusterExist(name, id)
	if err != nil || !exist {
		return
	}
	m.DeleteAt = time.Now()
	_, _ = ss.engine.Where("cluster = ?", m.Cluster).Cols("delete_at").Update(&m)
}

func (ss *statisticsService) ClusterList(page, pageSize int, search string, t string, replicas int, sort string, desc bool) ([]models.StatisticsCluster, map[string]int, error) {
	var session = ss.engine.Where("deploy_end is not null")
	summarySQL := "select count(*), floor(avg(extract(epoch FROM (sc.deploy_end - sc.deploy_start)))) average, min(extract(epoch FROM (sc.deploy_end - sc.deploy_start))) shortest, max(extract(epoch FROM (sc.deploy_end - sc.deploy_start))) longest from statistics_cluster sc where deploy_end is not null"
	summaryArgs := make([]interface{}, 0)
	if search != "" {
		search = "%" + search + "%"
		session.And("cluster like ?", search)
		summarySQL += " and cluster like ?"
		summaryArgs = append(summaryArgs, search)
	}
	if t != "" {
		session.And("type = ?", t)
		summarySQL += " and type = ?"
		summaryArgs = append(summaryArgs, t)
	}
	if replicas > 0 {
		session.And("replicas = ?", replicas)
		summarySQL += " and replicas = ?"
		summaryArgs = append(summaryArgs, replicas)
	}
	if sort != "" {
		if desc {
			session.Desc(sort)
		} else {
			session.OrderBy(sort)
		}
	}
	ret := make([]models.StatisticsCluster, 0)
	err := session.Limit(pageSize, pageSize*(page-1)).Find(&ret)
	if err != nil {
		return nil, nil, err
	}
	summary := map[string]int{}
	if len(ret) > 0 {
		for i := range ret {
			ret[i].Format()
		}
		_, err = session.SQL(summarySQL, summaryArgs...).Get(&summary)
		if err != nil {
			return nil, nil, err
		}
	}
	return ret, summary, nil
}
