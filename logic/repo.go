package logic

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/ningzining/L-ctl/util/caseutil"
	"github.com/ningzining/L-ctl/util/httputil"
	"github.com/ningzining/L-ctl/util/pathutil"
	"github.com/ningzining/L-ctl/util/templateutil"
	"net/url"
)

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

// Generate 生成repo文件
func (r *Repo) Generate(dirPath string, tableName string) error {
	fileName := fmt.Sprintf("%s.go", caseutil.ToCamelCase(tableName, false))
	filePath, err := url.JoinPath(dirPath, fileName)
	if err != nil {
		return err
	}
	// 判断目标文件是否已经存在
	exist, err := pathutil.Exist(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("生成模板失败,%s\n", err.Error()))
	}
	if exist {
		return errors.New(fmt.Sprintf("当前路径:`%s`已存在目标文件,生成失败,请选择另外的路径\n", filePath))
	}
	// 创建文件夹
	if err = pathutil.Mkdir(dirPath); err != nil {
		return err
	}

	// 新建文件
	m := make(map[string]interface{})
	m["Name"] = caseutil.ToCamelCase(tableName, true)
	m["TableName"] = tableName
	if err = createFile(filePath, m); err != nil {
		return err
	}
	color.Green("file is generated success at: %s", filePath)
	return nil
}

// 创建文件
func createFile(filePath string, data interface{}) error {
	init, err := isInit()
	if err != nil {
		return err
	}
	if !init {
		// 未初始化过模板，则使用github上面的模板进行渲染
		err := saveByOriginTemplate(filePath, data)
		if err != nil {
			return err
		}
	}
	// 如果初始化过模板，则使用本地的文件进行渲染创建模板文件
	err = saveByLocalTemplate(filePath, data)
	if err != nil {
		return err
	}
	return nil
}

// 判断是否进行了初始化
func isInit() (bool, error) {
	repoTemplate, err := templateutil.GetLocalRepoTemplate()
	if err != nil {
		return false, err
	}
	exist, err := pathutil.Exist(repoTemplate)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}
	return true, nil
}

// 通过github源文件创建模板
func saveByOriginTemplate(savePath string, data interface{}) error {
	// 通过http的get请求获取模板的字节数组
	templateData, err := httputil.Get(templateutil.TemplateRepoUrl)
	if err != nil {
		return errors.New(fmt.Sprintf("http请求异常,请检查网络,推荐先初始化模板到本地后进行操作\n%s\n", err.Error()))
	}
	// 通过字节数组保存模板
	if err = templateutil.SaveTemplateByData(templateData, savePath, data); err != nil {
		return err
	}
	return nil
}

// 通过本地文件创建模板
func saveByLocalTemplate(savePath string, data interface{}) error {
	// 获取本地模板文件的路径
	repoTemplatePath, err := templateutil.GetLocalRepoTemplate()
	if err != nil {
		return err
	}
	// 通过本地文件保存模板
	if err = templateutil.SaveTemplateByLocal(repoTemplatePath, savePath, data); err != nil {
		return err
	}
	return nil
}
