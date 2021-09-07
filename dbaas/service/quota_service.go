package service

import (
	"DBaas/models"
	"errors"
	"github.com/go-xorm/xorm"
)

type QuotaService interface {
	EditQuota(quota *models.ApiQuota) error
	ApiList() ([]*models.ApiQuota, error)
	ApiUsageList() ([]*models.ApiQuotaView, error)
}

type quotaService struct {
	engine *xorm.Engine
}

func (qs *quotaService) ApiUsageList() (ret []*models.ApiQuotaView, err error) {
	quotas := make([]models.ApiQuota, 0)
	err = qs.engine.Find(&quotas)
	if err != nil {
		return nil, err
	}
	ret = make([]*models.ApiQuotaView, len(quotas))
	for i := range quotas {
		usage, err := getApiUsage(quotas[i].Path, qs.engine)
		if err != nil {
			return nil, err
		}
		ret[i] = quotas[i].ToUsageView(usage)
	}
	return
}

func (qs *quotaService) ApiList() (ret []*models.ApiQuota, err error) {
	ret = make([]*models.ApiQuota, 0)
	err = qs.engine.Find(&ret)
	if err != nil {
		return nil, err
	}
	return
}

func (qs *quotaService) EditQuota(quota *models.ApiQuota) error {
	usage, err := getApiUsage(quota.Path, qs.engine)
	if err != nil {
		return err
	}
	if quota.Cpu < usage.Cpu || quota.Memory < usage.Memory || quota.Storage < usage.Storage {
		return errors.New("the set value is less than the current usage")
	}
	_, err = qs.engine.Where("id = ?", quota.Id).Omit("path").Update(quota)
	return err
}

func NewQuotaService(db *xorm.Engine) QuotaService {
	return &quotaService{db}
}
