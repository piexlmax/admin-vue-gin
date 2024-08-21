package biz_apphub

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/biz_apphub"
	biz_apphubReq "github.com/flipped-aurora/gin-vue-admin/server/model/biz_apphub/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/compress"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/httpx"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/oss"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
)

type BizAppHubService struct{}

type FileSrc struct {
	SrcPath     string `json:"srcPath"`     //文件中的原始地址
	OssFullPath string `json:"OssFullPath"` //替换到上传oss后的地址
	OssPath     string `json:"ossPath"`     //上传oss的地址
	LocalPath   string `json:"localPath"`   //本地所处的本地地址
}

func FileGetAll(path string) (files []string, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // 打印错误信息
			return err       // 可以选择返回错误，但这里我们选择继续遍历
		}
		if !info.IsDir() {
			files = append(files, path) // 将文件路径添加到切片中
		}
		return nil
	})

	fmt.Println(files, err)
	return files, nil
}

func (bizAppHubService *BizAppHubService) ParseFileSrc(unZipPath string) (fileList []FileSrc, err error) {
	fmt.Println(unZipPath)
	unZipPath = strings.ReplaceAll(unZipPath, "\\", "/")
	refPath := unZipPath
	unZipPath = "./" + unZipPath //./soft/beiluo/json-tool/json_conv/v1.0/dist
	files, err := FileGetAll(unZipPath)
	if err != nil {
		return nil, err
	}
	var filesSrc []FileSrc
	for _, file := range files {
		filePath := strings.ReplaceAll(file, "\\", "/")
		srcPath := strings.ReplaceAll(filePath, refPath, "")
		filesSrc = append(filesSrc, FileSrc{
			SrcPath:     srcPath,
			LocalPath:   "./" + filePath,
			OssFullPath: "http://cdn.geeleo.com/" + filePath,
			OssPath:     filePath,
		})
	}

	return filesSrc, nil
}

// FileSrcReplace 替换文件中的地址用oss的地址
func (bizAppHubService *BizAppHubService) FileSrcReplace(fileList []FileSrc) error {
	fileReplace := make(map[string]string)
	for _, file := range fileList {
		fileReplace[file.SrcPath] = file.OssFullPath
	}

	for _, src := range fileList {
		fileBytes, err := os.ReadFile(src.LocalPath)
		if err != nil {
			return err
		}
		fileContent := string(fileBytes)
		for k, v := range fileReplace { //把文件中的本地地址替换成oss地址
			fileContent = strings.ReplaceAll(fileContent, k, v)
		}

		err = os.Remove(src.LocalPath)
		if err != nil {
			return err
		}
		create, err := os.Create(src.LocalPath)
		if err != nil {
			return err
		}
		_, err = create.Write([]byte(fileContent))
		if err != nil {
			return err
		}
		create.Close()
	}
	for _, file := range fileList {
		store := oss.NewDefaultQiNiu()
		_, err := store.UploadLocalFile(file.LocalPath, file.OssPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bizAppHubService *BizAppHubService) Decode(req biz_apphub.BizAppHub) error {

	absPath := "./soft"
	absPath = strings.Join([]string{absPath, "beiluo", req.AppCode, req.Version}, "/")
	err := os.MkdirAll(absPath, 0755)
	if err != nil {
		return err
	}
	url := "http://cdn.geeleo.com/" + req.OssPath

	path := strings.Split(req.OssPath, "/")
	fileName := path[len(path)-1]
	err = httpx.DownloadFile(url, absPath+"/"+fileName)

	unZipPath, err := compress.DeCompress(filepath.Join(absPath, fileName), absPath)
	if err != nil {
		return err
	}

	list, err := bizAppHubService.ParseFileSrc(unZipPath)
	if err != nil {
		return err
	}

	fmt.Println(list)
	err = bizAppHubService.FileSrcReplace(list)
	if err != nil {
		return err
	}
	return nil
}

// CreateBizAppHub 创建biz_apphub记录
// Author [piexlmax](https://github.com/piexlmax)
func (bizAppHubService *BizAppHubService) CreateBizAppHub(bizAppHub *biz_apphub.BizAppHub) (err error) {
	err = bizAppHubService.Decode(*bizAppHub)
	if err != nil {
		return err
	}
	err = global.GVA_DB.Create(bizAppHub).Error
	return err
}

// DeleteBizAppHub 删除biz_apphub记录
// Author [piexlmax](https://github.com/piexlmax)
func (bizAppHubService *BizAppHubService) DeleteBizAppHub(ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&biz_apphub.BizAppHub{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&biz_apphub.BizAppHub{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteBizAppHubByIds 批量删除biz_apphub记录
// Author [piexlmax](https://github.com/piexlmax)
func (bizAppHubService *BizAppHubService) DeleteBizAppHubByIds(IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&biz_apphub.BizAppHub{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&biz_apphub.BizAppHub{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateBizAppHub 更新biz_apphub记录
// Author [piexlmax](https://github.com/piexlmax)
func (bizAppHubService *BizAppHubService) UpdateBizAppHub(bizAppHub biz_apphub.BizAppHub) (err error) {
	err = global.GVA_DB.Model(&biz_apphub.BizAppHub{}).Where("id = ?", bizAppHub.ID).Updates(&bizAppHub).Error
	return err
}

// GetBizAppHub 根据ID获取biz_apphub记录
// Author [piexlmax](https://github.com/piexlmax)
func (bizAppHubService *BizAppHubService) GetBizAppHub(ID string) (bizAppHub biz_apphub.BizAppHub, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&bizAppHub).Error
	return
}

// GetBizAppHubInfoList 分页获取biz_apphub记录
// Author [piexlmax](https://github.com/piexlmax)
func (bizAppHubService *BizAppHubService) GetBizAppHubInfoList(info biz_apphubReq.BizAppHubSearch) (list []biz_apphub.BizAppHub, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&biz_apphub.BizAppHub{})
	var bizAppHubs []biz_apphub.BizAppHub
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.AppName != "" {
		db = db.Where("app_name = ?", info.AppName)
	}
	if info.AppCode != "" {
		db = db.Where("app_code = ?", info.AppCode)
	}
	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}
	if info.Desc != "" {
		db = db.Where("desc LIKE ?", "%"+info.Desc+"%")
	}
	if info.Classify != "" {
		db = db.Where("classify = ?", info.Classify)
	}
	if info.Version != "" {
		db = db.Where("version = ?", info.Version)
	}
	if info.Mode != "" {
		db = db.Where("mode = ?", info.Mode)
	}
	if info.Tags != "" {
		db = db.Where("tags LIKE ?", "%"+info.Tags+"%")
	}
	if info.Video != "" {
		db = db.Where("video = ?", info.Video)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["version"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&bizAppHubs).Error
	return bizAppHubs, total, err
}
