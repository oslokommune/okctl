package core

import (
	"fmt"
	"io"

	"github.com/mishudark/errors"

	"github.com/google/uuid"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

const (
	defaultStateObjectName = "state.db"
	keyID                  = "LockID"
	keyDigest              = "Digest"
)

type remoteStateService struct {
	objectAPI api.ObjectStorageService
	kvAPI     api.KeyValueStoreService
}

// Upload knows how to upload a state.db file to a remote location
func (r *remoteStateService) Upload(clusterID api.ID, reader io.Reader) error {
	bucketName := generateStateDBBucketName(clusterID.ClusterName)

	_, err := r.objectAPI.CreateBucket(api.CreateBucketOpts{
		ClusterID:  clusterID,
		BucketName: bucketName,
		Private:    true,
	})
	if err != nil {
		return fmt.Errorf("creating bucket: %w", err)
	}

	err = r.objectAPI.PutObject(api.PutObjectOpts{
		ClusterID:  clusterID,
		BucketName: bucketName,
		Path:       defaultStateObjectName,
		Content:    reader,
	})
	if err != nil {
		return fmt.Errorf("putting object: %w", err)
	}

	return nil
}

// Download knows how to download a state.db file from a remote location
func (r *remoteStateService) Download(clusterID api.ID) (io.Reader, error) {
	object, err := r.objectAPI.GetObject(api.GetObjectOpts{
		ClusterID:  clusterID,
		BucketName: generateStateDBBucketName(clusterID.ClusterName),
		Path:       defaultStateObjectName,
	})
	if err != nil {
		return nil, errors.E(err, "getting object")
	}

	return object, nil
}

// AcquireStateLock knows how to activate a locking mechanism preventing others from mutating state
func (r *remoteStateService) AcquireStateLock(clusterID api.ID) error {
	storeName := generateStateLockStoreName(clusterID)

	err := r.kvAPI.CreateStore(api.CreateStoreOpts{
		ClusterID: clusterID,
		Name:      storeName,
		Keys:      []string{keyID, keyDigest},
	})
	if err != nil {
		return fmt.Errorf("creating store: %w", err)
	}

	isLocked, err := r.hasLock(clusterID)
	if err != nil {
		return fmt.Errorf("checking for existing lock: %w", err)
	}

	if isLocked {
		return errors.New("state is locked")
	}

	err = r.kvAPI.InsertItem(api.InsertItemOpts{
		ClusterID: clusterID,
		TableName: storeName,
		Item: api.KeyValueStoreItem{
			Fields: map[string]interface{}{
				keyID:     generateStateLockID(clusterID),
				keyDigest: uuid.New().String(),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("inserting item: %w", err)
	}

	return nil
}

// ReleaseStateLock knows how to deactivate a locking mechanism allowing others to mutate state
func (r *remoteStateService) ReleaseStateLock(clusterID api.ID) error {
	err := r.kvAPI.RemoveItem(api.DeleteItemOpts{
		ClusterID: clusterID,
		TableName: generateStateLockStoreName(clusterID),
		Field:     keyID,
		Key:       generateStateLockID(clusterID),
	})
	if err != nil {
		return fmt.Errorf("removing item: %w", err)
	}

	return nil
}

// Purge knows how to remove both remote state and related locking mechanisms
func (r *remoteStateService) Purge(clusterID api.ID) error {
	stateDBBucketName := generateStateDBBucketName(clusterID.ClusterName)

	err := r.objectAPI.EmptyBucket(api.EmptyBucketOpts{BucketName: stateDBBucketName})
	if err != nil {
		return fmt.Errorf("emptying state bucket: %w", err)
	}

	err = r.objectAPI.DeleteBucket(api.DeleteBucketOpts{
		ClusterID:  clusterID,
		BucketName: stateDBBucketName,
	})
	if err != nil {
		return fmt.Errorf("deleting object store: %w", err)
	}

	err = r.kvAPI.DeleteStore(api.DeleteStoreOpts{
		ClusterID: clusterID,
		Name:      generateStateLockStoreName(clusterID),
	})
	if err != nil {
		return fmt.Errorf("deleting key / value store: %w", err)
	}

	return nil
}

func (r *remoteStateService) hasLock(clusterID api.ID) (bool, error) {
	_, err := r.kvAPI.GetString(api.GetStringOpts{
		ClusterID: clusterID,
		TableName: generateStateLockStoreName(clusterID),
		Selector: api.ItemSelector{
			Key:   keyID,
			Value: generateStateLockID(clusterID),
		},
		Field: keyDigest,
	})
	if err == nil {
		return true, nil
	}

	if errors.IsKind(err, errors.NotExist) {
		return false, nil
	}

	return false, fmt.Errorf("fetching string: %w", err)
}

func generateStateLockStoreName(clusterID api.ID) api.StoreName {
	return api.StoreName(fmt.Sprintf("%s-state-lock", clusterID.ClusterName))
}

func generateStateLockID(clusterID api.ID) string {
	return fmt.Sprintf("okctl/%s/state.db", clusterID.ClusterName)
}

// NewRemoteStateService returns an initialized RemoteStateService
func NewRemoteStateService(kvAPI api.KeyValueStoreService, objectAPI api.ObjectStorageService) client.RemoteStateService {
	return &remoteStateService{
		kvAPI:     kvAPI,
		objectAPI: objectAPI,
	}
}

func generateStateDBBucketName(clusterName string) string {
	return fmt.Sprintf("okctl-%s-meta", clusterName)
}
