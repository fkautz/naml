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

package codify

import (
	"bytes"
	"fmt"
	"github.com/kris-nova/logger"
	corev1 "k8s.io/api/core/v1"
	"text/template"
	"time"
)

type ConfigMap struct {
	i *corev1.ConfigMap
}

func NewConfigMap(obj *corev1.ConfigMap) *ConfigMap {
	obj.ObjectMeta = cleanObjectMeta(obj.ObjectMeta)
	return &ConfigMap{
		i: obj,
	}
}

func (k ConfigMap) Install() string {
	l := Literal(k.i)
	install := fmt.Sprintf(`
	{{ .Name }}ConfigMap := %s

	_, err = client.CoreV1().ConfigMaps("{{ .Namespace }}").Create(context.TODO(), {{ .Name }}ConfigMap, v1.CreateOptions{})
	if err != nil {
		return err
	}
`, l)

	tpl := template.New(fmt.Sprintf("%s", time.Now().String()))
	tpl.Parse(install)
	buf := &bytes.Buffer{}
	k.i.Name = varName(k.i.Name)
	err := tpl.Execute(buf, k.i)
	if err != nil {
		logger.Debug(err.Error())
	}
	return alias(buf.String(), "corev1")
}

func (k ConfigMap) Uninstall() string {
	uninstall := `
	err = client.CoreV1().ConfigMaps("{{ .Namespace }}").Delete(context.TODO(), "{{ .Name }}", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
 `
	tpl := template.New(fmt.Sprintf("%s", time.Now().String()))
	tpl.Parse(uninstall)
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, k.i)
	if err != nil {
		logger.Debug(err.Error())
	}
	return buf.String()
}
