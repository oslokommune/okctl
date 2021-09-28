package run

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/oslokommune/okctl/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	appsv1 "k8s.io/api/apps/v1"
)

const (
	disableTCPEarlyDemuxKey = "DISABLE_TCP_EARLY_DEMUX"
	vpcCNIInitContainerName = "aws-vpc-cni-init"
)

func findTCPEarlyDemuxIndexes(ctx context.Context, initContainerIndex, envVarIndex *int) kube.ApplyFn {
	return func(clientSet kubernetes.Interface, config *rest.Config) (interface{}, error) {
		result, err := clientSet.AppsV1().DaemonSets("kube-system").Get(ctx, "aws-node", metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("getting aws-node data: %w", err)
		}

		containerIndex := getVpcCniInitContainerIndex(result)
		disableTCPEarlyDemuxVarIndex := getDisableTCPEarlyDemuxEnvIndex(result, containerIndex)

		*initContainerIndex = containerIndex
		*envVarIndex = disableTCPEarlyDemuxVarIndex

		return nil, nil
	}
}

func disableEarlyDemuxPatchApplier(ctx context.Context, rawPatch []byte) kube.ApplyFn {
	return func(clientSet kubernetes.Interface, config *rest.Config) (interface{}, error) {
		_, err := clientSet.AppsV1().DaemonSets("kube-system").Patch(
			ctx,
			"aws-node",
			types.JSONPatchType,
			rawPatch,
			metav1.PatchOptions{},
		)
		if err != nil {
			return nil, fmt.Errorf("patching aws-node: %w", err)
		}

		return nil, nil
	}
}

func getVpcCniInitContainerIndex(ds *appsv1.DaemonSet) int {
	for index, container := range ds.Spec.Template.Spec.InitContainers {
		if container.Name == vpcCNIInitContainerName {
			return index
		}
	}

	return -1
}

func getDisableTCPEarlyDemuxEnvIndex(ds *appsv1.DaemonSet, initContainerIndex int) int {
	for index, envVar := range ds.Spec.Template.Spec.InitContainers[initContainerIndex].Env {
		if envVar.Name == disableTCPEarlyDemuxKey {
			return index
		}
	}

	return -1
}

func generateRawDisableEarlyDemuxPatch(initContainerIndex, disableTCPEarlyDemuxVarIndex int) ([]byte, error) {
	patch := []map[string]interface{}{
		{
			"op": "replace",
			"path": fmt.Sprintf(
				"/spec/template/spec/initContainers/%d/env/%d",
				initContainerIndex,
				disableTCPEarlyDemuxVarIndex,
			),
			"value": map[string]string{
				"name":  disableTCPEarlyDemuxKey,
				"value": "true",
			},
		},
	}

	patchData, err := json.Marshal(patch)
	if err != nil {
		return nil, fmt.Errorf("marshalling demux json patch: %w", err)
	}

	return patchData, nil
}
