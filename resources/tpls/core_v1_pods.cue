package tpls

import (
	 corev1 "k8s.io/api/core/v1"
)

//  主对象
generic:{}

core_v1_pods_input:{
	 name: string
	 namespace: string | *"default"
	 image: string
}

core_v1_pods: corev1.#Pod & {
	 apiVersion: "v1"
	 kind: "Pod"
	 metadata: {
	 	name: core_v1_pods_input.name
	 	namespace: core_v1_pods_input.namespace
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
			containers: [{
					name:  core_v1_pods_input.name+"-container"
					image: core_v1_pods_input.image
			}]
   }
}
