package tpls

import (
	 corev1 "k8s.io/api/core/v1"
)

//  主对象
generic:{}

core_v1_services_input:{
	 name: string
	 namespace: string | *"default"
	 type: string | *"ClusterIP"
	 port: int
	 targetPort: int
}

core_v1_services: corev1.#Service & {
	 apiVersion: "v1"
	 kind: "Service"
	 metadata: {
	 	name: core_v1_services_input.name
	 	namespace: core_v1_services_input.namespace
	 	ownerReferences: [
    	 		{
    				apiVersion: generic.apiVersion
    				kind: generic.kind
    				name: generic.metadata.name
    				uid: generic.metadata.uid
    				controller: true
    	 		}
		]
	 }
	 spec: {
			type: core_v1_services_input.type
         ports: [{
          port:       core_v1_services_input.port
          targetPort: core_v1_services_input.targetPort
         }]
         selector: { //service通过selector和pod建立关联
          app: "test"
         }
   }
}
