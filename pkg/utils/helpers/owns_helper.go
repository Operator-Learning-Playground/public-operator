package helpers

import (
	"github.com/myoperator/common_operator/pkg/common"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// scanTplFileToGVRs 扫描cue模版文件名，并转为k8s的gvr
func scanTplFileToGVRs() []schema.GroupVersionResource {
	var ret []schema.GroupVersionResource
	_ = filepath.Walk(common.CueTplRoot, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			// 解析字段与拼接字段
			fName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			fList := strings.Split(fName, "_")
			if len(fList) == 3 {
				var g, v, r = fList[0], fList[1], fList[2]
				if g == "core" {
					g = ""
				}
				ret = append(ret, schema.GroupVersionResource{
					Group: g, Version: v, Resource: r,
				})
			}
		}
		return nil
	})
	return ret

}

// ScanTplFileToObjects 返回获取到的runtime.Object对象
func ScanTplFileToObjects() []runtime.Object {
	gvrs := scanTplFileToGVRs() // 得到cue文件所对应的gvr
	var gvks []schema.GroupVersionKind
	var objs []runtime.Object
	for _, gvr := range gvrs {
		list, err := common.K8sRestMapper.KindsFor(gvr)
		if err == nil {
			gvks = append(gvks, list...)
		}
	}
	allTypes := common.GlobalScheme.AllKnownTypes()
	for k, t := range allTypes {
		for _, gvk := range gvks {
			if k == gvk {
				if obj, ok := reflect.New(t).Interface().(runtime.Object); ok {
					objs = append(objs, obj)
				}
			}
		}
	}
	return objs
}
