/*


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

package controllers

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"os"
	"time"

	"github.com/go-logr/logr"
	v3 "github.com/velann21/etcd/clientv3"
	"github.com/velann21/etcd/pkg/transport"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	statefulv1 "statefulops/api/v1"
)

// ApplicationStateReconciler reconciles a ApplicationState object
type ApplicationStateReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

const (
	ClientPEM = "/Users/singaravelannandakumar/github/src/statefulops/controllers/server.cert"
	ClientKey = "/Users/singaravelannandakumar/github/src/statefulops/controllers/server.key"
	TrustedCA = "/Users/singaravelannandakumar/github/src/statefulops/controllers/ca.cert"
)

func New(url []string) (*v3.Client, error) {
	tlsInfo := transport.TLSInfo{
		CertFile:      ClientPEM,
		KeyFile:       ClientKey,
		TrustedCAFile: TrustedCA,
	}

	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}

	cfg := v3.Config{
		Endpoints:   url,
		DialTimeout: 10 * time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
		TLS:         tlsConfig,
	}

	c, err := v3.New(cfg)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("couldnt connect to etcd")
	}
	return c, nil
}

// +kubebuilder:rbac:groups=stateful.myoperator.state,resources=applicationstates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stateful.myoperator.state,resources=applicationstates/status,verbs=get;update;patch
func (r *ApplicationStateReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("applicationstate", req.NamespacedName)

	// your logic here

	fmt.Println("Inside the reconcil loop")
	dir, err := os.Getwd()
	fmt.Println(dir)
	etcdClient, err := New([]string{"192.168.64.5:2379"})
	if err != nil {

		fmt.Println(err, "Erroro")

	}

	resp, err := etcdClient.Put(context.Background(), "/registry/pods/testkey", "Sample testasasass")
	if err != nil {

		fmt.Println(err, "Put Error")

	}

	fmt.Println(resp)

	getresp, err := etcdClient.Get(context.Background(), "/registry/pods/testkey")
	if err != nil {

		fmt.Println(err, "Get Error")

	}

	for _, v := range getresp.Kvs {
		fmt.Println("Key is", string(v.Key))
		fmt.Println("value is", string(v.Value))
		fmt.Println("KV and v done")

	}

	watchChan := etcdClient.Watch(context.Background(), "/registry/pods/", v3.WithPrefix())
	go func(chn v3.WatchChan) {
		for wresp := range chn {
			for _, ev := range wresp.Events {
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}(watchChan)

	return ctrl.Result{}, nil
}

func (r *ApplicationStateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&statefulv1.ApplicationState{}).
		Complete(r)
}
