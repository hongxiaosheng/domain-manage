package controllers

import (
	"cmit.com/crd/domain-manage/api/v1alpha1"
	"context"
	"github.com/go-logr/logr"
	"go.etcd.io/etcd/clientv3"
	"istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cli *clientv3.Client
)

func GetValue(Domains []v1alpha1.Domains, reqLogger logr.Logger, istioClient *versioned.Clientset) {
	reqLogger.Info("istio  --- Get -------操作开始-----------")
	dms := Domains

	//取值
	for i := 0; i < len(dms); i++ {
		ns := dms[i].Namespace
		reqLogger.Info("istio ", "ns: ", ns)
		// 获取namespace下面的所有VS
		vsList, err := istioClient.NetworkingV1alpha3().VirtualServices(ns).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			return
		}
		// 遍历输出所有VS 的相关信息
		for i := range vsList.Items {
			vs := vsList.Items[i]
			reqLogger.Info("istio vs", "VirtualService Hosts:", vs.Spec.Hosts, " VirtualService Gateway:", vs.Spec.Gateways)

		}
	}
}

func DelValue(Domains []v1alpha1.Domains, reqLogger logr.Logger, istioClient *versioned.Clientset) {
	dms := Domains
	reqLogger.Info("-----istio DomainCancel begin-----")

	//取值
	for i := 0; i < len(dms); i++ {
		ns := dms[i].Namespace
		domain := dms[i].Domain
		option := dms[i].Option
		reqLogger.Info("istio ", "ns: ", ns)
		// 获取namespace下面的所有VS
		vsList, err := istioClient.NetworkingV1alpha3().VirtualServices(ns).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			reqLogger.Error(err, err.Error())
			return
		}
		reqLogger.Info("istio ", "vsList: ", vsList)
		// 遍历输出所有VS 的相关信息
		for j := range vsList.Items {
			vs := vsList.Items[j]
			reqLogger.Info("istio vs", "VirtualService Hosts:", vs.Spec.Hosts, " VirtualService Gateway:", vs.Spec.Gateways)
			hosts := vs.Spec.Hosts
			for k := 0; k < len(hosts); k++ {
				if hosts[k] == domain && "del" == option {
					istioClient.NetworkingV1alpha3().VirtualServices(ns).Delete(context.TODO(), vs.Name, v1.DeleteOptions{})
					reqLogger.Info("istio ", domain, " del success")
				}
			}
		}
	}
	reqLogger.Info("-----istio DomainCancel end-----")
	return
}
