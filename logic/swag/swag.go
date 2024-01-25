package swag

type Swag struct {
	File      string // swagger.json文件路径
	ProjectId string // 项目id
}

func NewSwag(file string, projectId string) *Swag {
	return &Swag{
		File:      file,
		ProjectId: projectId,
	}
}
