package tekton

import (
	"os"
	"testing"

	tektonconfig "g.hz.netease.com/horizon/pkg/config/tekton"

	"github.com/stretchr/testify/assert"
)

func TestNewTekton(t *testing.T) {
	tektonConfig := &tektonconfig.Tekton{
		Kubeconfig: "/",
	}
	tekton, err := NewTekton(tektonConfig)
	assert.Nil(t, tekton)
	assert.NotNil(t, err)

	tektonConfig = &tektonconfig.Tekton{
		Kubeconfig: "",
	}
	tekton, err = NewTekton(tektonConfig)
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	// 通过这两个环境变量判断是否在k8s集群中运行
	if len(host) == 0 || len(port) == 0 {
		assert.Nil(t, tekton)
		assert.NotNil(t, err)
	} else {
		assert.NotNil(t, tekton)
		assert.Nil(t, err)
	}
}
