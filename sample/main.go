//
// Copyright © 2021 Kris Nóva <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//   ███╗   ██╗ █████╗ ███╗   ███╗██╗
//   ████╗  ██║██╔══██╗████╗ ████║██║
//   ██╔██╗ ██║███████║██╔████╔██║██║
//   ██║╚██╗██║██╔══██║██║╚██╔╝██║██║
//   ██║ ╚████║██║  ██║██║ ╚═╝ ██║███████╗
//   ╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝
//

package main

import (
	"context"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kris-nova/logger"
	"github.com/kris-nova/naml"
	"k8s.io/client-go/kubernetes"
)

// main is the main entry point for your CLI application
func main() {

	naml.Register(New("beeps", "Basic busybox application.", "default"))
	naml.Register(New("boops", "Just another busybox application.", "default"))
	naml.Register(New("meeps-meeps", "Just another busybox app - like the others.", "default"))

	// Run the default CLI tooling
	err := naml.RunCommandLine()
	if err != nil {
		logger.Critical("%v", err)
		os.Exit(1)
	}
}

// NAMLApp is used for testing and debugging
type NAMLApp struct {
	description string
	meta        *metav1.ObjectMeta
}

func New(name, description, namespace string) *NAMLApp {
	return &NAMLApp{
		meta: &metav1.ObjectMeta{
			Name:            name,
			ResourceVersion: "1.0.0",
			Namespace:       namespace,
		},
		description: description,
	}
}

func (n *NAMLApp) Install(client *kubernetes.Clientset) error {
	deployment := naml.BusyboxDeployment(n.Meta().Name)
	_, err := client.AppsV1().Deployments(n.Meta().Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	return err
}

func (n *NAMLApp) Uninstall(client *kubernetes.Clientset) error {
	return client.AppsV1().Deployments(n.Meta().Namespace).Delete(context.TODO(), n.Meta().Name, metav1.DeleteOptions{})
}

func (n *NAMLApp) Description() string {
	return n.description
}

func (n *NAMLApp) Meta() *metav1.ObjectMeta {
	return n.meta
}
