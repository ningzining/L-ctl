package logic

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/ningzining/L-ctl/util/caseutil"
	"github.com/ningzining/L-ctl/util/httputil"
	"github.com/ningzining/L-ctl/util/pathutil"
	"github.com/ningzining/L-ctl/util/templateutil"
	"path/filepath"
)

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

const (
	CamelCase     = "lCtl"  // 驼峰命名
	UnderLineCase = "l_ctl" // 下划线命名
)

// Generate 生成repo文件
func (r *Repo) Generate(dir, tableName, style string) error {
	// 获取生成文件的路径
	filePath, err := generateFilePath(dir, tableName, style)
	if err != nil {
		return err
	}

	// 判断目标文件是否存在
	exist, err := pathutil.Exist(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("生成模板失败,%s\n", err))
	}
	if exist {
		return errors.New(fmt.Sprintf("当前路径:`%s`已存在目标文件,生成失败,请选择另外的路径\n", filePath))
	}

	// 创建文件夹
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return errors.New(fmt.Sprintf("获取绝对路径失败: %s\n", err))
	}
	if err = pathutil.MkdirIfNotExist(dirAbs); err != nil {
		return err
	}

	// 获取数据并生成模板文件
	data := genRepoTemplateData(filePath, tableName)
	if err = createRepoTemplate(filePath, data); err != nil {
		return err
	}

	color.Green("文件生成成功: %s", filePath)
	return nil
}

// 获取生成目标文件的路径
func generateFilePath(dirPath, tableName, style string) (string, error) {
	var fileName string
	switch style {
	case UnderLineCase:
		fileName = fmt.Sprintf("%s.go", caseutil.ToUnderLineCase(tableName))
	case CamelCase:
		fileName = fmt.Sprintf("%s.go", caseutil.ToCamelCase(tableName, false))
	default:
		fileName = fmt.Sprintf("%s.go", caseutil.ToUnderLineCase(tableName))
	}
	filePath := filepath.Join(dirPath, fileName)
	abs, err := filepath.Abs(filePath)
	if err != nil {
		return "", nil
	}
	return abs, nil
}

// 创建文件
func createRepoTemplate(filePath string, data map[string]any) error {
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
	repoTemplate, err := templateutil.GetRepoTemplatePath()
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
func saveByLocalTemplate(savePath string, data map[string]any) error {
	// 获取本地模板文件的路径
	templatePath, err := templateutil.GetRepoTemplatePath()
	if err != nil {
		return err
	}
	// 通过本地文件保存模板
	if err = templateutil.SaveTemplateByLocal(templatePath, savePath, data); err != nil {
		return err
	}
	return nil
}

// 生成模板所需要的data数据
func genRepoTemplateData(filePath, tableName string) map[string]any {
	data := map[string]any{
		"Name":      caseutil.ToCamelCase(tableName, true),
		"TableName": tableName,
		"pkg":       filepath.Base(filePath),
	}
	return data
}
