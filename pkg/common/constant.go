package common

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

// CueTplRoot 模板根目录
const CueTplRoot = "./resources/tpls"

// 管理器初始化的全局变量
var (
	K8sRestMapper meta.RESTMapper
	RestConfig    *rest.Config
	GlobalScheme  *runtime.Scheme
)
