package main

import (
	"cuelang.org/go/cue"
	"fmt"
	"github.com/myoperator/common_operator/pkg/apis/generic/v1alpha1"
	"github.com/myoperator/common_operator/pkg/controllers"
	"github.com/myoperator/common_operator/pkg/utils/helpers"
)

func main() {

	//g := &v1alpha1.Generic{
	//	Spec: v1alpha1.GenericSpec{
	//		Gvr: "v1/pods",
	//		Template: map[string]interface{}{
	//			"name": "pods-test",
	//			"namespace": "kube-system",
	//		},
	//	},
	//}
	//g := &v1alpha1.Generic{
	//	Spec: v1alpha1.GenericSpec{
	//		Gvr: "v1/services",
	//		Template: map[string]interface{}{
	//			"name": "pods-test",
	//			"namespace": "kube-system",
	//		},
	//	},
	//}
	g := &v1alpha1.Generic{
		Spec: v1alpha1.GenericSpec{
			Gvr: "apps/v1/deployments",
			Template: map[string]interface{}{
				"name":      "pods-test",
				"namespace": "kube-system",
				"image":     "test",
				"replicas":  3,
			},
		},
	}
	cueName, cueFilePath := controllers.ConvertToCueFile(g.Spec.Gvr)

	inst := helpers.MustLoadFileInstance(cueFilePath)
	filldCV := inst.Fill(map[string]interface{}{
		"generic":          inst.Context().Encode(g),
		cueName + "_input": inst.Context().Encode(g.Spec.Template),
	})
	// 模版没渲染的内容
	fmt.Println(filldCV.LookupPath(cue.ParsePath(cueName)))
	// 转为json内容
	jsonBytes, _ := filldCV.LookupPath(cue.ParsePath(cueName)).MarshalJSON()
	fmt.Println(string(jsonBytes))

}
