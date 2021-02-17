// Package storageclass provides a storage class creator and applier
package storageclass

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// StorageClass contains the state for creating a k8s storage class
type StorageClass struct {
	Name                 string
	Parameters           *EBSParameters
	Annotations          map[string]string
	ReclaimPolicy        corev1.PersistentVolumeReclaimPolicy
	AllowVolumeExpansion bool
	VolumeBindingMode    v1.VolumeBindingMode
	Ctx                  context.Context
}

// VolumeType enumerates the storage types allowed by EBS
// - https://github.com/kubernetes-sigs/aws-ebs-csi-driver#createvolume-parameters
type VolumeType string

// nolint: golint
const (
	IO1VolumeType VolumeType = "io1"
	IO2VolumeType VolumeType = "io2"
	GP2VolumeType VolumeType = "gp2"
	GP3VolumeType VolumeType = "gp3"
	SC1VolumeType VolumeType = "sc1"
	ST1VolumeType VolumeType = "st1"
)

// FileSystemType enumerates the types available for a
// preformatted volume
// - https://github.com/kubernetes-sigs/aws-ebs-csi-driver#createvolume-parameters
type FileSystemType string

// nolint: golint
const (
	XFSFileSystemType  FileSystemType = "xfs"
	EXT2FileSystemType FileSystemType = "ext2"
	EXT3FileSystemType FileSystemType = "ext3"
	EXT4FileSystemType FileSystemType = "ext4"
)

// EncryptedType enumerates the options
// for encrypting a volume
// - https://github.com/kubernetes-sigs/aws-ebs-csi-driver#createvolume-parameters
type EncryptedType string

// nolint: golint
const (
	EncryptedTrue  EncryptedType = "true"
	EncryptedFalse EncryptedType = "false"
)

// EBSParameters contains the parameters that can be set
// to modify the profile of the storage class
type EBSParameters struct {
	Type       VolumeType
	IOPSPerGB  string // Required with io1 or io2 volume types
	IOPS       string // Only applicable for gp3
	Throughput string // Only applicable for gp3
	Encrypted  EncryptedType
	FSType     FileSystemType
}

// Validate the parameters
func (p *EBSParameters) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.IOPSPerGB, validation.When(p.Type == IO1VolumeType || p.Type == IO2VolumeType, validation.Required).Else(validation.Empty)),
		validation.Field(&p.IOPS, validation.When(p.Type == GP3VolumeType, validation.Required).Else(validation.Empty)),
		validation.Field(&p.Throughput, validation.When(p.Type == GP3VolumeType, validation.Required).Else(validation.Empty)),
	)
}

// ToMap produces a parameters map
func (p *EBSParameters) ToMap() map[string]string {
	params := map[string]string{}

	if len(p.Type) > 0 {
		params["type"] = string(p.Type)
	}

	if len(p.FSType) > 0 {
		params["csi.storage.k8s.io/fsType"] = string(p.FSType)
	}

	if len(p.IOPS) > 0 {
		params["iops"] = p.IOPS
	}

	if len(p.IOPSPerGB) > 0 {
		params["iopsPerGB"] = p.IOPSPerGB
	}

	if len(p.Throughput) > 0 {
		params["throughput"] = p.Throughput
	}

	if len(p.Encrypted) > 0 {
		params["encrypted"] = string(p.Encrypted)
	}

	return params
}

// NewEBSParameters returns the AWS defaults
func NewEBSParameters() *EBSParameters {
	return &EBSParameters{
		Type:       GP3VolumeType,
		IOPS:       "3000",
		Throughput: "125",
		Encrypted:  EncryptedTrue,
		FSType:     EXT4FileSystemType,
	}
}

// New returns an initialised storage class creator
func New(name string, parameters *EBSParameters, annotations map[string]string) (*StorageClass, error) {
	err := parameters.Validate()
	if err != nil {
		return nil, err
	}

	return &StorageClass{
		Name:                 name,
		Parameters:           parameters,
		Annotations:          annotations,
		ReclaimPolicy:        corev1.PersistentVolumeReclaimDelete,
		AllowVolumeExpansion: true,
		VolumeBindingMode:    v1.VolumeBindingWaitForFirstConsumer,
		Ctx:                  context.Background(),
	}, nil
}

// DeleteStorageClass deletes a storage class
func (n *StorageClass) DeleteStorageClass(_ kubernetes.Interface, config *rest.Config) (interface{}, error) {
	client := kubernetes.NewForConfigOrDie(config)

	ns, err := client.StorageV1().StorageClasses().List(n.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	found := false

	for _, item := range ns.Items {
		if item.Name == n.Name {
			found = true
		}
	}

	if !found {
		return nil, nil
	}

	policy := metav1.DeletePropagationForeground

	return nil, client.StorageV1().StorageClasses().Delete(n.Ctx, n.Name, metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
}

// CreateStorageClass creates a storage Class
func (n *StorageClass) CreateStorageClass(clientset kubernetes.Interface, _ *rest.Config) (interface{}, error) {
	client := clientset.StorageV1().StorageClasses()

	classes, err := client.List(n.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, ns := range classes.Items {
		if ns.Name == n.Name {
			r, err := client.Get(n.Ctx, ns.Name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			return r, nil
		}
	}

	return client.Create(n.Ctx, n.StorageClassManifest(), metav1.CreateOptions{})
}

// StorageClassManifest returns the storage class manifest
func (n *StorageClass) StorageClassManifest() *v1.StorageClass {
	return &v1.StorageClass{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StorageClass",
			APIVersion: "storage.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        n.Name,
			Annotations: n.Annotations,
		},
		Provisioner:          "ebs.csi.aws.com",
		Parameters:           n.Parameters.ToMap(),
		ReclaimPolicy:        &n.ReclaimPolicy,
		AllowVolumeExpansion: &n.AllowVolumeExpansion,
		VolumeBindingMode:    &n.VolumeBindingMode,
	}
}

// DefaultStorageClassAnnotation returns the annotation that can be used
// to indicate that a storage class should be the default one.
func DefaultStorageClassAnnotation() map[string]string {
	return map[string]string{
		"storageclass.kubernetes.io/is-default-class": "true",
	}
}
