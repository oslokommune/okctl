package hooks

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/oslokommune/okctl/pkg/logging"

	"github.com/spf13/afero"

	"github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/asdine/storm/v3/codec/json"
	merrors "github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	clientServicesErrFormat = "acquiring client services: %w"
	localStatePathErrFormat = "acquiring local state path: %w"
)

var (
	// ErrImmutable indicates no change is possible
	ErrImmutable = errors.New("immutable")
	// ErrNotFound indicates something is missing
	ErrNotFound = errors.New("not found")
	// ErrNotInitialized indicates the state database has yet to be initialized
	ErrNotInitialized = errors.New("not initialized")
)

// DownloadState downloads state from remote storage or initializes it
// Requires the okctl struct to be initialized.
//
// writable: set this to true if you are going to mutate the state
//
// Remember to acquire a lock when state can potentially be mutated
func DownloadState(o *okctl.Okctl, writable bool) RunEer {
	return func(_ *cobra.Command, _ []string) error {
		if o.DB == nil {
			return ErrNotInitialized
		}

		services, err := o.ClientServices(o.StateHandlers(o.StateNodes()))
		if err != nil {
			return fmt.Errorf(clientServicesErrFormat, err)
		}

		localStateDBPath, err := getLocalStatePath(o)
		if err != nil {
			return fmt.Errorf(localStatePathErrFormat, err)
		}

		state, err := services.RemoteState.Download(metadataAsClusterID(o.Declaration.Metadata))
		if err != nil {
			if !merrors.IsKind(err, merrors.NotExist) {
				return fmt.Errorf("acquiring remote state: %w", err)
			}

			o.Logger.Debug("state database not found, creating.")

			state, err = initializeStateDB(o.FileSystem)
			if err != nil {
				return fmt.Errorf("initializing state database: %w", err)
			}
		}

		err = o.FileSystem.WriteReader(localStateDBPath, state)
		if err != nil {
			return fmt.Errorf("storing state: %w", err)
		}

		o.DB.SetWritable(writable)
		o.DB.SetDatabaseFilePath(localStateDBPath)

		return nil
	}
}

// UploadState uploads local state to remote storage
// Requires the okctl struct to be initialized.
//
// Remember to release the lock after uploading the state.
func UploadState(o *okctl.Okctl) RunEer {
	log := logging.GetLogger(logComponent, "UploadState")

	return func(_ *cobra.Command, _ []string) error {
		log.Debug("starting")

		if o.DB == nil {
			return ErrNotInitialized
		}

		if !o.DB.IsWritable() {
			return ErrImmutable
		}

		services, err := o.ClientServices(o.StateHandlers(o.StateNodes()))
		if err != nil {
			return fmt.Errorf(clientServicesErrFormat, err)
		}

		localStateDBPath, err := getLocalStatePath(o)
		if err != nil {
			return fmt.Errorf(localStatePathErrFormat, err)
		}

		f, err := o.FileSystem.Open(localStateDBPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return ErrNotFound
			}

			return fmt.Errorf("opening local state: %w", err)
		}

		log.Debug("uploading state")

		err = services.RemoteState.Upload(metadataAsClusterID(o.Declaration.Metadata), f)
		if err != nil {
			return fmt.Errorf("uploading state: %w", err)
		}

		err = f.Close()
		if err != nil {
			return fmt.Errorf("closing local state: %w", err)
		}

		return nil
	}
}

// ClearLocalState cleans up the localState directory ensuring no leftover local state
func ClearLocalState(o *okctl.Okctl) RunEer {
	return func(_ *cobra.Command, _ []string) error {
		localStateDBPath, err := getLocalStatePath(o)
		if err != nil {
			return fmt.Errorf(localStatePathErrFormat, err)
		}

		err = o.FileSystem.RemoveAll(path.Dir(localStateDBPath))
		if err != nil {
			return fmt.Errorf("removing temp state dir: %w", err)
		}

		return nil
	}
}

// AcquireStateLock prevents state to be mutated by anyone other than the client in possession of the lock
func AcquireStateLock(o *okctl.Okctl) RunEer {
	return func(_ *cobra.Command, _ []string) error {
		services, err := o.ClientServices(o.StateHandlers(o.StateNodes()))
		if err != nil {
			return fmt.Errorf(clientServicesErrFormat, err)
		}

		err = services.RemoteState.AcquireStateLock(metadataAsClusterID(o.Declaration.Metadata))
		if err != nil {
			return fmt.Errorf("acquiring lock: %w", err)
		}

		return nil
	}
}

// ReleaseStateLock releases the state mutation prevention allowing other clients to mutate state
func ReleaseStateLock(o *okctl.Okctl) RunEer {
	return func(_ *cobra.Command, _ []string) error {
		services, err := o.ClientServices(o.StateHandlers(o.StateNodes()))
		if err != nil {
			return fmt.Errorf(clientServicesErrFormat, err)
		}

		err = services.RemoteState.ReleaseStateLock(metadataAsClusterID(o.Declaration.Metadata))
		if err != nil {
			return fmt.Errorf("releasing lock: %w", err)
		}

		return nil
	}
}

// PurgeRemoteState removes all traces of remote state suitable for a delete environment kind of operation
func PurgeRemoteState(o *okctl.Okctl) RunEer {
	log := logging.GetLogger(logComponent, "PurgeRemoteState")

	return func(_ *cobra.Command, _ []string) error {
		log.Debug("running hook")

		services, err := o.ClientServices(o.StateHandlers(o.StateNodes()))
		if err != nil {
			return fmt.Errorf(clientServicesErrFormat, err)
		}

		localStateDBPath, err := getLocalStatePath(o)
		if err != nil {
			return fmt.Errorf(localStatePathErrFormat, err)
		}

		okctlDir, err := o.GetUserDataDir()
		if err != nil {
			return fmt.Errorf("acquiring okctl directory: %w", err)
		}

		log.Debug("backing up state database")

		err = backupStateDatabase(o.FileSystem, okctlDir, o.Declaration.Metadata.Name, localStateDBPath)
		if err != nil {
			return fmt.Errorf("backing up local database: %w", err)
		}

		log.Debug("purging state")

		err = services.RemoteState.Purge(metadataAsClusterID(o.Declaration.Metadata))
		if err != nil {
			return fmt.Errorf("purging remote state: %w", err)
		}

		return nil
	}
}

// VerifyClusterExistsInState ensures we have a cluster in the state before continuing
// This could happen if maintenance mode was not run after upgrading
// to 0.0.80 and we don't have a state.db in S3 - see #491
func VerifyClusterExistsInState(o *okctl.Okctl) RunEer {
	return func(_ *cobra.Command, _ []string) error {
		handlers := o.StateHandlers(o.StateNodes())

		_, err := handlers.Cluster.GetCluster(o.Declaration.Metadata.Name)
		if err != nil {
			o.Logger.Debug("Could not find cluster in current state: ", err)

			if errors.Is(err, storm.ErrNotFound) {
				return fmt.Errorf("getting existing cluster %s: %w", o.Declaration.Metadata.Name, err)
			}

			return fmt.Errorf("getting cluster from state: %w", err)
		}

		return nil
	}
}

func backupStateDatabase(fs *afero.Afero, rootDir string, clusterName string, dbPath string) error {
	stateBackupPath := path.Join(
		rootDir,
		".backups",
		fmt.Sprintf("%s-%s-state.db", clusterName, time.Now().Format(time.RFC3339)),
	)

	err := fs.MkdirAll(path.Dir(stateBackupPath), 0o700)
	if err != nil {
		return fmt.Errorf("creating state db backup dir: %w", err)
	}

	f, err := fs.Open(dbPath)
	if err != nil {
		return fmt.Errorf("opening local state database: %w", err)
	}

	defer func() {
		_ = f.Close()
	}()

	err = fs.WriteReader(stateBackupPath, f)
	if err != nil {
		return fmt.Errorf("writing state database backup: %w", err)
	}

	return nil
}

func initializeStateDB(fs *afero.Afero) (io.Reader, error) {
	baseDir := fs.GetTempDir("okctl")
	tempDBPath := path.Join(baseDir, "fresh.db")

	db, err := storm.Open(tempDBPath, storm.Codec(json.Codec))
	if err != nil {
		return nil, err
	}

	err = db.Close()
	if err != nil {
		return nil, fmt.Errorf("closing database: %w", err)
	}

	rawDB, err := fs.ReadFile(tempDBPath)
	if err != nil {
		return nil, fmt.Errorf("buffering database: %w", err)
	}

	return bytes.NewBuffer(rawDB), nil
}

func metadataAsClusterID(metadata v1alpha1.ClusterMeta) api.ID {
	return api.ID{
		Region:       metadata.Region,
		AWSAccountID: metadata.AccountID,
		ClusterName:  metadata.Name,
	}
}

func getLocalStatePath(o *okctl.Okctl) (string, error) {
	dataDir, err := o.GetUserDataDir()
	if err != nil {
		return "", fmt.Errorf("acquiring user data dir: %w", err)
	}

	dir := path.Join(dataDir, "localState", o.Declaration.Metadata.Name)

	err = o.FileSystem.MkdirAll(dir, 0o700)
	if err != nil {
		return "", fmt.Errorf("creating temp state folder: %w", err)
	}

	return path.Join(dir, constant.DefaultStormDBName), nil
}
