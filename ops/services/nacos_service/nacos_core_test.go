package nacos_service

import (
	"fmt"
	"testing"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 7/20/21 11:12 AM
 * @Desc:
 */

func Test_GetConfig(t *testing.T) {
	nacos, err := NewNacosClient("127.0.0.1:8849","nacos", "nacos")
	if err != nil {
		t.Fatal(err)
	}
	nsList, err := nacos.GetNamespace()
	fmt.Println(err, nsList)
	//data, configType, err := nacos.GetConfig("dev", "com.fulu.ops", "mysql")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(configType)
	//t.Log(data)
	//err = nacos.CopyConfig("dev", "com.fulu.ops", "mysql", "prod", "com.fulu.ops1", "mysql")
	//if err != nil {
	//	t.Fatal(err)
	//}
}