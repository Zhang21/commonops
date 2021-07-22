package nacos

import (
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/services/nacos_service"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/chujieyang/commonops/ops/value_objects/nacos"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 7/20/21 3:07 PM
 * @Desc:
 */

func IPostNacos(c *gin.Context) {
	var req nacos.NacosServer
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	_, err = models.AddNewNacosServer(req.Alias, req.EndPoint, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IPostNacosConfig(c *gin.Context) {
	var req nacos.CreateNacosConfigReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosInfo, err := models.GetNacosInfoById(req.Id)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err = nacosClient.PublishConfig(req.Namespace, req.DataId, req.Group, req.Content, req.ConfigType); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IPostNacosConfigCopy(c *gin.Context) {
	var req nacos.CreateNacosConfigCopyReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosInfo, err := models.GetNacosInfoById(req.Id)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err = nacosClient.CopyConfig(req.SrcNamespace, req.SrcDataId, req.SrcGroup, req.DstNamespace, req.DstDataId, req.DstGroup); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IPutNacosConfig(c *gin.Context) {
	var req nacos.UpdateNacosConfigReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosInfo, err := models.GetNacosInfoById(req.Id)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err = nacosClient.PublishConfig(req.Namespace, req.DataId, req.Group, req.Content, req.ConfigType); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IDeleteNacosConfig(c *gin.Context) {
	var req nacos.DeleteNacosConfigReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosInfo, err := models.GetNacosInfoById(req.Id)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err = nacosClient.DeleteConfig(req.Namespace, req.DataId, req.Group); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IGetNacosList(c *gin.Context) {
	nacosList, err := models.GetNacosList()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nacosList})
}

func IGetNacosNamespaceList(c *gin.Context) {
	var req nacos.NacosNsListReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosInfo, err := models.GetNacosInfoById(req.Id)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nsList, err := nacosClient.GetNamespace()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nsList})
}

func IGetNacosConfigList(c *gin.Context) {
	var req nacos.NacosConfigsReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosInfo, err := models.GetNacosInfoById(req.Id)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	configsList, err := nacosClient.GetNsConfigs(req.Namespace, req.Page, req.Size)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: configsList})
}