package tpls

import (
	 appsv1 "k8s.io/api/apps/v1"
)

//  主对象
generic:{}

apps_v1_deployments_input:{
	 name: string
	 namespace: string | *"default"
	 image: string
	 replicas: number
	 imagePullPolicy: string | *"IfNotPresent"
	 containerPort: 80
}

apps_v1_deployments: appsv1.#Deployment & {
	 apiVersion: "apps/v1"
	 kind: "Deployment"
	 metadata: {
	 	name: apps_v1_deployments_input.name
	 	namespace: apps_v1_deployments_input.namespace
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
      selector: matchLabels: app: "test"
      replicas: apps_v1_deployments_input.replicas
      template: {
       metadata: labels: app: "test"
       spec: containers: [{
        name:            apps_v1_deployments_input.name+"-container"
        image:           apps_v1_deployments_input.image
        imagePullPolicy: apps_v1_deployments_input.imagePullPolicy
        ports: [{
         containerPort: apps_v1_deployments_input.containerPort
        }]
       }]
      }
     }
}
