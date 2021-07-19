package cron

import (
	"fmt"
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/opslog"
	"strings"
	"time"

	"context"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/robfig/cron"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func init() {
	opslog.Info().Println("cron init exec ...")
	c := cron.New()
	registerCronJobs(c, "0 0 23 * *", SyncAliYunEcsData)
	registerCronJobs(c, "0 10 23 * *", SyncAliYunRdsData)
	registerCronJobs(c, "0 20 23 * *", SyncAliYunKvData)
	registerCronJobs(c, "0 30 23 * *", SyncAliYunSlbData)
	registerCronJobs(c, "*/15 * * * *", collectK8sMetricsData)
	c.Start()
}

func registerCronJobs(c *cron.Cron, spec string, cmd func()) {
	if err := c.AddFunc(spec, cmd); err != nil {
		opslog.Error().Println(err)
	}
}

/**
同步阿里云的 ecs 服务器数据
*/
func SyncAliYunEcsData() {
	opslog.Info().Println("定时任务 [同步ecs数据]")
	err := service.GetEcsService().EcsDiffCacheClean()
	if err != nil {
		opslog.Error().Println(err)
	}
	accountList := models.GetAllCloudAccounts()
	for _, account := range accountList {
		accountId := account.Id
		accessKey := account.Key
		accessSecret := account.Secret
		if accessKey == "" || accessKey == "未知" {
			continue
		}
		opslog.Info().Printf("同步云账号信息：%s - %s \n", account.Name, account.Key)
		data := service.GetEcsService().GetEcsList(accessKey, accessSecret)
		for _, v := range data {
			ecs := models.Ecs{
				ImageId:                 v.ImageId,
				InstanceType:            v.InstanceType,
				InstanceNetworkType:     v.InstanceNetworkType,
				LocalStorageAmount:      fmt.Sprintf("%d", v.LocalStorageAmount),
				InstanceChargeType:      v.InstanceChargeType,
				ClusterId:               v.ClusterId,
				InstanceName:            v.InstanceName,
				StartTime:               utils.ConvertUtcTimeToLocal(v.StartTime, "2006-01-02T15:04Z"),
				ZoneId:                  v.ZoneId,
				InternetChargeType:      v.InternetChargeType,
				InternetMaxBandwidthIn:  fmt.Sprintf("%d", v.InternetMaxBandwidthIn),
				HostName:                v.HostName,
				Cpu:                     fmt.Sprintf("%d", v.Cpu),
				Status:                  v.Status,
				OSName:                  v.OSName,
				OSNameEn:                v.OSNameEn,
				SerialNumber:            v.SerialNumber,
				RegionId:                v.RegionId,
				InternetMaxBandwidthOut: fmt.Sprintf("%d", v.InternetMaxBandwidthOut),
				ResourceGroupId:         v.ResourceGroupId,
				InstanceTypeFamily:      v.InstanceTypeFamily,
				InstanceId:              v.InstanceId,
				DeploymentSetId:         v.DeploymentSetId,
				Description:             v.Description,
				ExpiredTime:             utils.ConvertUtcTimeToLocal(v.ExpiredTime, "2006-01-02T15:04Z"),
				OSType:                  v.OSType,
				Memory:                  fmt.Sprintf("%d", v.Memory),
				CreationTime:            utils.ConvertUtcTimeToLocal(v.CreationTime, "2006-01-02T15:04Z"),
				LocalStorageCapacity:    fmt.Sprintf("%d", v.LocalStorageCapacity),
				InnerIpAddress:          strings.Join(v.InnerIpAddress.IpAddress, ","),
				PublicIpAddress:         strings.Join(v.PublicIpAddress.IpAddress, ","),
				PrivateIpAddress:        strings.Join(v.VpcAttributes.PrivateIpAddress.IpAddress, ","),
				DataStatus:              1,
			}
			models.SaveOrUpdateEcs(v.InstanceId, accountId, ecs)
			opslog.Info().Printf("[cron jobs - SyncAliYunEcsData]: %s, %s, %s \n ", v.HostName,
				v.InstanceId, v.InstanceName)
		}
	}
}

/**
同步阿里云的 Rds 服务器数据
*/
func SyncAliYunRdsData() {
	opslog.Info().Println("定时任务 [同步rds数据]")
	err := service.GetRdsService().RdsDiffCacheClean()
	if err != nil {
		opslog.Error().Println(err)
	}
	accountList := models.GetAllCloudAccounts()
	for _, account := range accountList {
		accountId := account.Id
		accessKey := account.Key
		accessSecret := account.Secret
		if accessKey == "" || accessKey == "未知" {
			continue
		}
		opslog.Info().Printf("[SyncAliYunRdsData] 同步云账号信息：%s, %s \n", account.Name, account.Key)
		data := service.GetRdsService().GetRdsList(accessKey, accessSecret)
		for _, v := range data {
			opslog.Info().Printf("[SyncAliYunRdsData] rds: %s \n", v.DBInstance.DBInstanceId)
			attribute := v.DBAttribute
			rds := models.Rds{
				DataStatus:            1,
				InsId:                 v.DBInstance.InsId,
				DBInstanceId:          v.DBInstance.DBInstanceId,
				DBInstanceDescription: v.DBInstance.DBInstanceDescription,
				PayType:               v.DBInstance.PayType,
				DBInstanceType:        v.DBInstance.DBInstanceType,
				RegionId:              v.DBInstance.RegionId,
				ExpireTime:            utils.ConvertUtcTimeToLocal(v.DBInstance.ExpireTime, "2006-01-02T15:04:05Z"),
				DestroyTime:           utils.ConvertUtcTimeToLocal(v.DBInstance.DestroyTime, "2006-01-02T15:04:05Z"),
				DBInstanceStatus:      v.DBInstance.DBInstanceStatus,
				Engine:                v.DBInstance.Engine,
				DBInstanceNetType:     v.DBInstance.DBInstanceNetType,
				ConnectionMode:        v.DBInstance.ConnectionMode,
				LockMode:              v.DBInstance.LockMode,
				Category:              v.DBInstance.Category,
				DBInstanceStorageType: v.DBInstance.DBInstanceStorageType,
				DBInstanceClass:       v.DBInstance.DBInstanceClass,
				InstanceNetworkType:   v.DBInstance.InstanceNetworkType,
				VpcCloudInstanceId:    v.DBInstance.VpcCloudInstanceId,
				LockReason:            v.DBInstance.LockReason,
				ZoneId:                v.DBInstance.ZoneId,
				MutriORsignle:         v.DBInstance.MutriORsignle,
				CreateTime:            utils.ConvertUtcTimeToLocal(v.DBInstance.CreateTime, "2006-01-02T15:04:05Z"),
				EngineVersion:         v.DBInstance.EngineVersion,
				GuardDBInstanceId:     v.DBInstance.GuardDBInstanceId,
				TempDBInstanceId:      v.DBInstance.TempDBInstanceId,
				MasterInstanceId:      v.DBInstance.MasterInstanceId,
				VpcId:                 v.DBInstance.VpcId,
				VSwitchId:             v.DBInstance.VSwitchId,
				ReplicateId:           v.DBInstance.ReplicateId,
				ResourceGroupId:       v.DBInstance.ResourceGroupId,
				ConnectionString:      attribute.ConnectionString,
				Port:                  attribute.Port,
				DBInstanceMemory:      int(attribute.DBInstanceMemory),
				DBInstanceStorage:     attribute.DBInstanceStorage,
			}
			isSuccess := models.SaveOrUpdateRds(v.DBInstance.DBInstanceId, accountId, rds)
			opslog.Info().Printf("%s - %s - %s - %s - %s \n ", v.DBInstance.DBInstanceId,
				v.DBInstance.DBInstanceDescription, v.DBInstance.DBInstanceStatus,
				attribute.ConnectionString, isSuccess)
		}
	}
}

/**
同步阿里云的 KvStore 服务器数据
*/
func SyncAliYunKvData() {
	opslog.Info().Println("[SyncAliYunKvData] 定时任务 [同步kv数据]")
	err := service.GetKvService().KvDiffCacheClean()
	if err != nil {
		opslog.Error().Println(err)
	}
	accountList := models.GetAllCloudAccounts()
	for _, account := range accountList {
		accountId := account.Id
		accessKey := account.Key
		accessSecret := account.Secret
		if accessKey == "" || accessKey == "未知" {
			continue
		}
		opslog.Info().Printf("[SyncAliYunKvData] 同步云账号信息：%s, %s \n", account.Name, account.Key)
		data := service.GetKvService().GetKvList(accessKey, accessSecret)
		for _, v := range data {
			kv := models.Kv{
				DataStatus:          1,
				InstanceClass:       v.InstanceClass,
				PackageType:         v.PackageType,
				ChargeType:          v.ChargeType,
				ConnectionDomain:    v.ConnectionDomain,
				SearchKey:           v.SearchKey,
				CreateTime:          utils.ConvertUtcTimeToLocal(v.CreateTime, "2006-01-02T15:04:05Z"),
				HasRenewChangeOrder: fmt.Sprintf("%v", v.HasRenewChangeOrder),
				InstanceType:        v.InstanceType,
				DestroyTime:         utils.ConvertUtcTimeToLocal(v.DestroyTime, "2006-01-02T15:04:05Z"),
				RegionId:            v.RegionId,
				PrivateIp:           v.PrivateIp,
				InstanceId:          v.InstanceId,
				InstanceStatus:      v.InstanceStatus,
				Bandwidth:           int(v.Bandwidth),
				NetworkType:         v.NetworkType,
				VpcId:               v.VpcId,
				NodeType:            v.NodeType,
				Connections:         int(v.Connections),
				ArchitectureType:    v.ArchitectureType,
				ReplacateId:         v.ReplacateId,
				EngineVersion:       v.EngineVersion,
				Capacity:            int(v.Capacity),
				VSwitchId:           v.VSwitchId,
				InstanceName:        v.InstanceName,
				Port:                int(v.Port),
				ZoneId:              v.ZoneId,
				EndTime:             utils.ConvertUtcTimeToLocal(v.EndTime, "2006-01-02T15:04:05Z"),
				QPS:                 int(v.QPS),
				UserName:            v.UserName,
				Config:              v.Config,
				IsRds:               v.IsRds,
				ConnectionMode:      v.ConnectionMode,
			}
			isSuccess := models.SaveOrUpdateKv(v.InstanceId, accountId, kv)
			opslog.Info().Printf("%s - %s - %s - %s \n", v.InstanceId, v.InstanceName, v.InstanceClass, isSuccess)
		}
	}
}

/**
同步阿里云的 SLB 数据
*/
func SyncAliYunSlbData() {
	opslog.Info().Println("[SyncAliYunSlbData] 定时任务 [同步slb数据]")
	err := service.GetSlbService().SlbDiffCacheClean()
	if err != nil {
		opslog.Error().Println(err)
	}
	accountList := models.GetAllCloudAccounts()
	for _, account := range accountList {
		accountId := account.Id
		accessKey := account.Key
		accessSecret := account.Secret
		if accessKey == "" || accessKey == "未知" {
			continue
		}
		opslog.Info().Printf("[SyncAliYunSlbData] 同步云账号信息：%s, %s \n", account.Name, account.Key)
		data := service.GetSlbService().GetSlbList(accessKey, accessSecret)
		for _, v := range data {
			slb := models.Slb{
				DataStatus:         1,
				Count:              0,
				SlaveZoneId:        v.SlaveZoneId,
				LoadBalancerStatus: v.LoadBalancerStatus,
				VSwitchId:          v.VSwitchId,
				MasterZoneId:       v.MasterZoneId,
				PayType:            v.PayType,
				RegionIdAlias:      v.RegionIdAlias,
				CreateTime:         utils.ConvertUtcTimeToLocal(v.CreateTime, "2006-01-02T15:04Z"),
				Address:            v.Address,
				LoadBalancerId:     v.LoadBalancerId,
				AddressIPVersion:   v.AddressIPVersion,
				RegionId:           v.RegionId,
				ResourceGroupId:    v.ResourceGroupId,
				LoadBalancerName:   v.LoadBalancerName,
				InternetChargeType: v.InternetChargeType,
				AddressType:        v.AddressType,
				VpcId:              v.VpcId,
				NetworkType:        v.NetworkType,
			}
			isSuccess := models.SaveOrUpdateSlb(v.LoadBalancerId, accountId, slb)
			opslog.Info().Printf("%s - %s - %s - %s \n", v.LoadBalancerId, v.LoadBalancerName,
				v.LoadBalancerStatus, isSuccess)
		}
	}
}

func DiffDataCheck() {
	type Item struct {
		InstanceId string
	}
	var resList []Item

	// ecs data check
	database.Mysql().Raw("select instance_id from ecs inner join cloud_account as ca " +
		"on ecs.cloud_account_id = ca.id where ca.cloud_type = '阿里云' and ecs.data_status = 1").Scan(&resList)
	for _, item := range resList {
		var count = 0
		database.Mysql().Raw("select count(*) from diff_caches where "+
			"type = 'ecs' and instance_id = ? limit 1 ", item.InstanceId).Count(&count)
		if count == 0 {
			err := database.Mysql().Exec("update ecs set data_status = 0 where "+
				"instance_id = ?", item.InstanceId).Error
			if err != nil {
				opslog.Error().Println(err)
			}
		}
	}

	// rds data check
	database.Mysql().Raw("select db_instance_id as instance_id from rds inner join cloud_account as ca " +
		"on rds.cloud_account_id = ca.id where ca.cloud_type = '阿里云' and rds.data_status = 1").Scan(&resList)
	for _, item := range resList {
		var count = 0
		database.Mysql().Raw("select count(*) from diff_caches where "+
			"type = 'rds' and instance_id = ? limit 1 ", item.InstanceId).Count(&count)
		if count == 0 {
			err := database.Mysql().Exec("update rds set data_status = 0 where "+
				"db_instance_id = ?", item.InstanceId).Error
			if err != nil {
				opslog.Error().Println(err)
			}
		}
	}

	// redis data check
	database.Mysql().Raw("select instance_id from kv inner join cloud_account as ca " +
		"on kv.cloud_account_id = ca.id where ca.cloud_type = '阿里云' and kv.data_status = 1").Scan(&resList)
	for _, item := range resList {
		var count = 0
		database.Mysql().Raw("select count(*) from diff_caches where "+
			"type = 'kv' and instance_id = ? limit 1 ", item.InstanceId).Count(&count)
		if count == 0 {
			err := database.Mysql().Exec("update kv set data_status = 0 where "+
				"instance_id = ?", item.InstanceId).Error
			if err != nil {
				opslog.Error().Println(err)
			}
		}
	}

	// slb data check
	database.Mysql().Raw("select load_balancer_id as instance_id from slb inner join cloud_account as ca " +
		"on slb.cloud_account_id = ca.id where ca.cloud_type = '阿里云' and slb.data_status = 1").Scan(&resList)
	for _, item := range resList {
		var count = 0
		database.Mysql().Raw("select count(*) from diff_caches where "+
			"type = 'slb' and instance_id = ? limit 1 ", item.InstanceId).Count(&count)
		if count == 0 {
			err := database.Mysql().Exec("update slb set data_status = 0 where "+
				"load_balancer_id = ?", item.InstanceId).Error
			if err != nil {
				opslog.Error().Println(err)
			}
		}
	}
}

func collectK8sMetricsData() {
	cleanTime := time.Duration(5) * time.Minute
	if err := models.CleanDatabase(&cleanTime); err != nil {
		opslog.Error().Println(err)
	}
	var k8sList []models.K8s
	if err := database.Mysql().Raw("select * from k8s where data_status > 0 ").Scan(&k8sList).Error; err != nil {
		opslog.Error().Println(err)
		return
	}
	for _, k8sInfo := range k8sList {
		config := &rest.Config{
			Host:        k8sInfo.ApiServer,
			BearerToken: k8sInfo.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true,
			},
		}
		k8sClients, err := kubernetes.NewForConfig(config)
		if err != nil {
			opslog.Error().Println(err)
		}
		namespaceList, err := k8sClients.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		clients, err := metrics.NewForConfig(config)
		if err != nil {
			opslog.Error().Println(err)
		}
		nodeMetrics, err := clients.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
		if err := models.UpdateNodeMetrics(k8sInfo.ClusterId, nodeMetrics); err != nil {
			opslog.Error().Println(err)
		}
		for _, namespace := range namespaceList.Items {
			podMetrics, _ := clients.MetricsV1beta1().PodMetricses(namespace.Name).List(context.TODO(), metav1.ListOptions{})
			if err := models.UpdatePodMetrics(k8sInfo.ClusterId, podMetrics); err != nil {
				opslog.Error().Println(err)
			}
		}
	}

}
