package init_controller

import (
	"github.com/myoperator/common_operator/pkg/apis/generic/v1alpha1"
	"github.com/myoperator/common_operator/pkg/common"
	"github.com/myoperator/common_operator/pkg/controllers"
	"github.com/myoperator/common_operator/pkg/utils/helpers"
	"log"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

// InitManager 初始化控制器管理器
func InitManager() {

	logf.SetLogger(zap.New())
	mgr, err := manager.New(K8sRestConfig(),
		manager.Options{
			Logger: logf.Log.WithName("public-operator"),
	})

	if err != nil {
		log.Fatal("unable to set up manager:", err.Error())
	}

	// FIXME: 搬到init()中
	common.K8sRestMapper = mgr.GetRESTMapper()
	common.RestConfig = mgr.GetConfig()
	common.GlobalScheme = mgr.GetScheme() // 保存Scheme 到全局对象

	// 2. ++ 注册进入序列化表
	err = v1alpha1.SchemeBuilder.AddToScheme(mgr.GetScheme())
	if err != nil {
		mgr.GetLogger().Error(err, "unable add schema")
		os.Exit(1)
	}

	// 3. 控制器相关
	genericCtl := controllers.NewGenericController(
		mgr.GetEventRecorderFor("public-operator"),
	)

	// 4. 获取目前需要Owns的资源对象
	ownsObjects := helpers.ScanTplFileToObjects()

	bdr := builder.ControllerManagedBy(mgr).
		For(&v1alpha1.Generic{})

		// 这是另一个方法，外部手动加入，但不使用，使用自动注册Owns的方式
		//Owns(&v1.Pod{})
		//Watches(&source.Kind{Type: &v1.Pod{}},
		//	handler.Funcs{
		//		UpdateFunc: genericCtl.OnUpdatePodHandler,
		//		DeleteFunc: genericCtl.OnDeletePodHandler,
		//	},
		//)

	for _, obj := range ownsObjects {
		// 这里要特别注意，需要断言为client.Object对象
		bdr = bdr.Owns(obj.(client.Object))
	}

	if err = bdr.Complete(genericCtl); err != nil {
		mgr.GetLogger().Error(err, "unable to create manager")
		os.Exit(1)
	}

	// 5. 启动controller管理器
	if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
		mgr.GetLogger().Error(err, "unable to start manager")
	}
}
