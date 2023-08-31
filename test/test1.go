package main

import (
	"cuelang.org/go/cue"
	"fmt"
	"github.com/myoperator/common_operator/pkg/apis/generic/v1alpha1"
	"github.com/myoperator/common_operator/pkg/controllers"
	"github.com/myoperator/common_operator/pkg/utils/helpers"
)

func main() {

	g := &v1alpha1.Generic{
		Spec: v1alpha1.GenericSpec{
			Gvr: "v1/pods",
			Template: map[string]interface{}{
				"name":  "testpod",
				"image": "nginx:1.18-alpine",
			},
		},
	}
	g.APIVersion = "api.practice.com/v1alpha1"
	g.Kind = "Generic"
	g.Name = "xxxoo"
	g.UID = "xxxx"

	cueName, cueFilePath := controllers.ConvertToCueFile(g.Spec.Gvr)
	inst := helpers.MustLoadFileInstance(cueFilePath)
	//filldCV := inst.FillPath(cue.ParsePath(cueName+"_input"),
	//	inst.Context().Encode(g.Spec.Template))
	filldCV := inst.Fill(map[string]interface{}{
		"generic":          inst.Context().Encode(g),
		cueName + "_input": inst.Context().Encode(g.Spec.Template),
	})

	//filldCV := inst.FillPath(cue.ParsePath(cueName+"_input"), inst.Context().Encode(g.Spec.Template))
	jsonBytes, err := filldCV.LookupPath(cue.ParsePath(cueName)).MarshalJSON()

	fmt.Println(err)
	fmt.Println(string(jsonBytes))

}
