package service

import (
	"DBaas/models"
	"DBaas/x/response"
	"errors"
	"github.com/go-xorm/xorm"
)

type ImageService interface {
	List(page int, pageSize int, key string) ([]models.Images, int64, error)
	Add(image models.Images) error
	InitAdd(images []models.Images) error
	Update(image models.Images) error
	Delete(id int) error
	Param(id int) ([]models.Defaultparameters, error)
	// 获取镜像类型或类别, imageType为空时返回Type集合, 不为空返回Category集合
	GetImageType(imageType string) []string
}

type imagesService struct {
	Engine *xorm.Engine
	cs     CommonService
}

func NewImageService(db *xorm.Engine, cs CommonService) ImageService {
	return &imagesService{
		Engine: db,
		cs:     cs,
	}
}

func getImageAddress(engine *xorm.Engine) string {
	address := models.Sysparameter{ParamKey: "harbor_address"}
	_, _ = engine.Cols("param_value").Get(&address)
	return address.ParamValue
}

func (is *imagesService) GetImageType(imageType string) []string {
	types := make([]models.ImageType, 0)
	images := make([]models.Images, 0)
	_ = is.Engine.Distinct("image_type_id").Find(&images)
	// type去重用
	distinctType := make(map[string]struct{})
	getType := imageType == ""
	if getType {
		_ = is.Engine.Find(&types)
	} else {
		_ = is.Engine.Where("type=?", imageType).Find(&types)
	}
	var result []string
	for _, imageType := range types {
		if imageType.Type == "Operator" {
			// 镜像列表中已经添加过的镜像类型直接跳过, 仅限Operator类型
			for i := range images {
				if images[i].ImageTypeId == imageType.Id {
					goto End
				}
			}
		}
		if getType {
			if _, has := distinctType[imageType.Type]; !has {
				result = append(result, imageType.Type)
				distinctType[imageType.Type] = struct{}{}
			}
		} else {
			result = append(result, imageType.Category)
		}
	End:
	}
	return result
}

func (is *imagesService) List(page int, pageSize int, key string) ([]models.Images, int64, error) {
	typeMap := make(map[int64]models.ImageType)
	err := is.Engine.Find(&typeMap)
	imagesList := make([]models.Images, 0)
	count, err := is.Engine.Where("image_name like ?", "%"+key+"%").Count(new(models.Images))
	err = is.Engine.Where("image_name like ?", "%"+key+"%").Limit(pageSize, pageSize*(page-1)).Desc("id").Find(&imagesList)
	address := getImageAddress(is.Engine)
	for i, image := range imagesList {
		// 地址拼接名称和版本
		imagesList[i].Address = address + "/" + image.ImageName + ":" + image.Version
		// 设置类型和类别
		imagesList[i].Type = typeMap[int64(image.ImageTypeId)].Type
		imagesList[i].Category = typeMap[int64(image.ImageTypeId)].Category
	}
	return imagesList, count, err
}

// 根据镜像的类型和类别查询ImageTypeId
func (is *imagesService) getImageTypeId(image models.Images) int {
	imageType := models.ImageType{Type: image.Type, Category: image.Category}
	_, _ = is.Engine.Cols("id").Get(&imageType)
	return imageType.Id
}

func (is *imagesService) Add(image models.Images) error {
	session := is.Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	image.SetStatus(getImageAddress(is.Engine))
	image.ImageTypeId = is.getImageTypeId(image)
	if _, err := session.Insert(image); err != nil {
		_ = session.Rollback()
		return err
	}
	return session.Commit()
}

func (is *imagesService) InitAdd(images []models.Images) error {
	session := is.Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	for _, image := range images {
		image.SetStatus(getImageAddress(is.Engine))
		image.ImageTypeId = is.getImageTypeId(image)
		_, err := session.Insert(image)
		if err != nil {
			_ = session.Rollback()
			return err
		}
	}
	return session.Commit()
}

func (is *imagesService) Update(image models.Images) error {
	session := is.Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	// 检查镜像状态
	var checkImage models.Images
	_, err := session.ID(image.Id).Cols("image_name").Get(&checkImage)
	if err != nil {
		return err
	}
	checkImage.Version = image.Version
	checkImage.SetStatus(getImageAddress(is.Engine))
	image.Status = checkImage.Status
	image.ImageTypeId = is.getImageTypeId(image)
	_, err = session.ID(image.Id).Cols("version", "description", "status", "image_type_id").Update(image)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	return session.Commit()
}

func (is *imagesService) Delete(id int) (err error) {
	imageType := models.ImageType{}
	_, err = is.Engine.
		Cols("type").
		Join("LEFT OUTER", "images", "images.image_type_id = image_type.id").
		Where("images.id = ?", id).Get(&imageType)
	if err != nil {
		return
	}
	if imageType.Type == "Operator" {
		status, _ := models.GetConfigBool("operator@status", is.Engine)
		if status {
			return errors.New(response.ErrorImageOccupied.En)
		}
	} else if imageType.Type == "Mysql" {
		exist, _ := is.Engine.Exist(&models.ClusterInstance{ImageId: id})
		if exist {
			return errors.New(response.ErrorImageOccupied.En)
		}
	}
	_, err = is.Engine.Id(id).Delete(new(models.Images))
	return err
}

func getImageParam(imageId int, engine *xorm.Engine) ([]models.Defaultparameters, error) {
	defaultParams := make([]models.Defaultparameters, 0)
	err := engine.SQL(
		"select dp.* from defaultparameters dp inner join images i on dp.image_type_id = i.image_type_id where i.id = ?",
		imageId).Find(&defaultParams)
	return defaultParams, err
}

func (is *imagesService) Param(id int) ([]models.Defaultparameters, error) {
	return getImageParam(id, is.Engine)
}
