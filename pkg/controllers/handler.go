package controllers

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/pkg/strings"
	"fmt"
	"github.com/myoperator/common_operator/pkg/apis/generic/v1alpha1"
	"github.com/myoperator/common_operator/pkg/utils/helpers"
	"github.com/myoperator/common_operator/pkg/common"
)

// 这些符号 一律要替换成 _
var replaces = []string{
	"\\", "/", ".",
}

// ConvertToCueFile 通过GVR来获取cue文件名
// v1/pods or apps/v1/deployments or v1.pods  v1\pods v1/pods
// 返回值：第一个是名称，第二个是文件名
func ConvertToCueFile(gvr string) (string, string) {
	//  替换v1_pods     apps_v1_deployments
	for _, r := range replaces {
		gvr = strings.Replace(gvr, r, "_", -1)
	}
	var cuePath string
	gvrList := strings.Split(gvr, "_")

	// 区分 core 与其他资源组的区别
	if len(gvrList) == 2 {
		cuePath = fmt.Sprintf("core_%s_%s", gvrList[0], gvrList[1])
	} else if len(gvrList) == 3 {
		cuePath = fmt.Sprintf("%s_%s_%s", gvrList[0], gvrList[1], gvrList[2])
	}
	return cuePath, fmt.Sprintf("%s/%s.cue", common.CueTplRoot, cuePath)
}

// handler 处理调协逻辑
// 1. 根据crd的字段获取cue模版对象
// 2. 根据template传入内容，渲染cue模版
// 3. 转为k8s apply可使用的对象
// 4. k8s apply
func handler(g *v1alpha1.Generic) error {
	// 根据 crd中定义的gvr解析cue模板名称和对应的文件名
	cueName, cueFile := ConvertToCueFile(g.Spec.Gvr)
	if cueName == "" || cueFile == "" {
		return fmt.Errorf("wrong gvr or config option")
	}

	// 根据template传入内容，渲染cue模版
	inst := helpers.MustLoadFileInstance(cueFile)

	filldCV := inst.Fill(map[string]interface{}{
		"generic":          inst.Context().Encode(g),
		cueName + "_input": inst.Context().Encode(g.Spec.Template),
	})
	if filldCV.Err() != nil {
		return filldCV.Err()
	}

	// 转为k8s apply可使用的对象
	jsonBytes, err := filldCV.LookupPath(cue.ParsePath(cueName)).MarshalJSON()
	if err != nil {
		return err
	}

	// apply操作
	_, err = helpers.K8sApply(jsonBytes, common.RestConfig, common.K8sRestMapper)
	if err != nil {
		return err
	}

	return nil
}
