/*
Copyright 2018 Pressinfra SRL

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mysqlcluster

import (
	"github.com/presslabs/controller-util/syncer"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/presslabs/mysql-operator/pkg/internal/mysqlcluster"
)

// NewHealthyReplicasSVCSyncer returns a service syncer for healthy replicas service
func NewHealthyReplicasSVCSyncer(c client.Client, scheme *runtime.Scheme, cluster *mysqlcluster.MysqlCluster) syncer.Interface {
	obj := &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.GetNameForResource(mysqlcluster.HealthyReplicasService),
			Namespace: cluster.Namespace,
		},
	}

	return syncer.NewObjectSyncer("HealthyReplicasSVC", cluster.Unwrap(), obj, c, scheme, func(in runtime.Object) error {
		out := in.(*core.Service)

		// set service labels
		out.Labels = cluster.GetLabels()
		out.Labels["mysql.presslabs.org/service-type"] = "ready-replicas"

		// set selectors for healthy replica (non-master) mysql pods only
		out.Spec.Selector = cluster.GetSelectorLabels()
		out.Spec.Selector["role"] = "replica"
		out.Spec.Selector["healthy"] = "yes"

		if len(out.Spec.Ports) != 2 {
			out.Spec.Ports = make([]core.ServicePort, 2)
		}

		out.Spec.Ports[0].Name = MysqlPortName
		out.Spec.Ports[0].Port = MysqlPort
		out.Spec.Ports[0].TargetPort = TargetPort
		out.Spec.Ports[0].Protocol = core.ProtocolTCP

		out.Spec.Ports[1].Name = SidecarServerPortName
		out.Spec.Ports[1].Port = SidecarServerPort
		out.Spec.Ports[1].Protocol = core.ProtocolTCP

		return nil
	})

}
