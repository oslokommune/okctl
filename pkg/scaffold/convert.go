// Package resources knows how to produce and handle Kubernetes resource
package scaffold

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"sigs.k8s.io/yaml"

	v1 "k8s.io/api/core/v1"
)

var (
	emptyMatcher, _  = regexp.Compile("^.*: (null|{})$")
	statusMatcher, _ = regexp.Compile("^\\s*?status*:$")
)

// volumesAsBytes knows how to convert a Kubernetes PersistentVolumeClaim to a byte array
func volumesAsBytes(volumes []*v1.PersistentVolumeClaim) ([]byte, error) {
	var writer bytes.Buffer

	for index := range volumes {
		err := writeCleanResource(&writer, volumes[index])
		if err != nil {
			return nil, fmt.Errorf("writing volume to buffer: %w", err)
		}
	}

	return writer.Bytes(), nil
}

// resourceAsBytes knows how to convert a Kubernetes resource to a byte array
func resourceAsBytes(data interface{}) ([]byte, error) {
	var writer bytes.Buffer

	err := writeCleanResource(&writer, data)
	if err != nil {
		return nil, fmt.Errorf("writing resource definition to buffer: %w", err)
	}

	return writer.Bytes(), nil
}

func writeResource(w io.Writer, resource interface{}) error {
	serializedResource, err := yaml.Marshal(resource)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s\n---\n", serializedResource)
	if err != nil {
		return err
	}

	return nil
}

func writeCleanResource(w io.Writer, resource interface{}) error {
	var buf bytes.Buffer

	err := writeResource(&buf, resource)
	if err != nil {
		return err
	}

	result, err := cleanResources(buf)
	if err != nil {
		return err
	}

	_, err = w.Write(result)
	if err != nil {
		return err
	}

	return nil
}

func cleanResources(buf bytes.Buffer) ([]byte, error) {
	content, err := ioutil.ReadAll(&buf)
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer
	for _, item := range strings.Split(string(content), "\n") {
		if !emptyMatcher.MatchString(item) && !statusMatcher.MatchString(item) {
			result.Write([]byte(item + "\n"))
		}
	}

	return result.Bytes(), nil
}
