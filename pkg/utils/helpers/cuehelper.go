package helpers

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"fmt"
)

var NotFoundError = fmt.Errorf("not found")

type CueHelper struct {
	file     string
	instance cue.Value
}

//// filllabel 是模板中的 struct 名称。 如果没有会panic
//// data 是一个cue.value ,, 一般是 上传到到content article 中的文件
//func (c *CueHelper) FillData(filllabel string, data cue.Value) *CueHelper {
//	fillPath := c.instance.LookupPath(cue.ParsePath(filllabel))
//	if !fillPath.Exists() {
//		panic("无法找到FillData对应的label")
//	}
//
//	dataCvBytes, _ := data.MarshalJSON()
//
//	c.instance = c.instance.FillPath(cue.ParsePath(filllabel),
//		c.instance.Context().CompileBytes(dataCvBytes))
//
//	if c.instance.Err() != nil {
//		panic(c.instance.Err())
//	}
//	return c
//}

//func (c *CueHelper) MustParse(path string) string {
//	ret, err := c.Parse(path)
//	if err != nil {
//		panic(err)
//	}
//	return ret
//}
//func (c *CueHelper) GetInstance() cue.Value {
//	return c.instance
//}
//
//// 只支持 Int64
//func (c *CueHelper) GetValueForInt(path string) int64 {
//	getPath := c.instance.LookupPath(cue.ParsePath(path))
//	if !getPath.Exists() {
//		return 0
//	}
//	i, err := c.instance.LookupPath(cue.ParsePath(path)).Int64()
//	if err != nil {
//		return 0
//	}
//	return i
//}
//
//// 给Value设置值 。 思路是自己 实现的。 没有经过大量测试性能
//// value 目前只支持 string
//func (c *CueHelper) SetValue(path string, value interface{}) error {
//	b, err := c.instance.MarshalJSON()
//	if err != nil {
//		return err
//	}
//
//	getpath := gjson.Get(string(b), path)
//	if !getpath.Exists() {
//		return NotFoundError
//	}
//	newString, err := sjson.Set(string(b), path, value)
//	if err != nil {
//		return err
//	}
//	c.instance = c.instance.Context().CompileString(newString)
//	if c.instance.Err() != nil {
//		return c.instance.Err()
//	}
//	return nil
//}


func (c *CueHelper) Parse(path string) (string, error) {
	v := c.instance.LookupPath(cue.ParsePath(path))
	if !v.Exists() {
		return "", NotFoundError
	}
	return v.String()
}

// NewCueHelper
func NewCueHelper(file string) *CueHelper {
	v := MustLoadFileInstance(file)
	return &CueHelper{file: file, instance: v}
}

// NewCueHelperFromInstance 直接传cue.value进行初始化
func NewCueHelperFromInstance(file string, v cue.Value) *CueHelper {
	return &CueHelper{file: file, instance: v}
}

// MustLoadFileInstance 获取 cue 实例
func MustLoadFileInstance(filepath string) cue.Value {
	cv, err := LoadFileInstance(filepath)
	if err != nil {
		panic(err)
	}
	return cv
}

// LoadFileInstance 获取cue实例
func LoadFileInstance(filepath string) (cue.Value, error) {
	insts := load.Instances([]string{filepath}, nil)
	if len(insts) != 1 {
		return cue.Value{}, fmt.Errorf("load instance error:%s\n", filepath)
	}
	cc := cuecontext.New()
	cv := cc.BuildInstance(insts[0])
	if cv.Err() != nil {
		return cue.Value{}, cv.Err()
	}
	return cv, nil
}

// MustEncode 通用的转化对象
func MustEncode(object interface{}) cue.Value {
	cc := cuecontext.New()
	v := cc.Encode(object)
	if v.Err() != nil {
		panic(v.Err())
	}
	return v
}
