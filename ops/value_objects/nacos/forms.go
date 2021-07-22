package nacos

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 7/20/21 3:09 PM
 * @Desc:
 */

type NacosServer struct {
	Alias    string `json:"alias" form:"alias" required:"true"`
	EndPoint string    `json:"endpoint" form:"endpoint" required:"true"`
	Username string `json:"username" form:"username" required:"true"`
	Password string `json:"password" form:"password" required:"true"`
}

type NacosNsListReq struct {
	Id string    `json:"id" form:"id" required:"true"`
}

type NacosConfigsReq struct {
	Id string    `json:"id" form:"id" required:"true"`
	Namespace string    `json:"namespace" form:"namespace" required:"true"`
	Page int    `json:"page" form:"page" required:"true"`
	Size int    `json:"size" form:"size" required:"true"`
}

type CreateNacosConfigReq struct {
	Id string  `json:"id" form:"id" required:"true"`
	Namespace string  `json:"namespace" form:"namespace" required:"true"`
	DataId string  `json:"dataId" form:"dataId" required:"true"`
	Group string  `json:"group" form:"group" required:"true"`
	Content string  `json:"content" form:"content" required:"true"`
	ConfigType string  `json:"configType" form:"configType" required:"true"`
}

type UpdateNacosConfigReq struct {
	Id string  `json:"id" form:"id" required:"true"`
	ConfigId string  `json:"configId" form:"configId" required:"true"`
	Namespace string  `json:"namespace" form:"namespace" required:"true"`
	DataId string  `json:"dataId" form:"dataId" required:"true"`
	Group string  `json:"group" form:"group" required:"true"`
	Content string  `json:"content" form:"content" required:"true"`
	ConfigType string  `json:"configType" form:"configType" required:"true"`
}

type DeleteNacosConfigReq struct {
	Id string  `json:"id" form:"id" required:"true"`
	Namespace string  `json:"namespace" form:"namespace" required:"true"`
	DataId string  `json:"dataId" form:"dataId" required:"true"`
	Group string  `json:"group" form:"group" required:"true"`
}

type CreateNacosConfigCopyReq struct {
	Id string  `json:"id" form:"id" required:"true"`
	SrcNamespace string  `json:"srcNamespace" form:"srcNamespace" required:"true"`
	SrcDataId string  `json:"srcDataId" form:"srcDataId" required:"true"`
	SrcGroup string  `json:"srcGroup" form:"srcGroup" required:"true"`
	DstNamespace string  `json:"dstNamespace" form:"dstNamespace" required:"true"`
	DstDataId string  `json:"dstDataId" form:"dstDataId" required:"true"`
	DstGroup string  `json:"dstGroup" form:"dstGroup" required:"true"`
}