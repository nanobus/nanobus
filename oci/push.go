package oci

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/errcode"
)

func Push(reference, base string, fileRefs []string, dryRun bool) error {
	store := file.New("")
	defer store.Close()

	ctx := context.Background()
	fileAnnotations := map[string]map[string]string{}

	descs, err := loadFiles(ctx, store, fileAnnotations, base, fileRefs, false)
	if err != nil {
		return err
	}

	artifactType := "application/vnd.nanobus.iota.v1+json"
	annotations := map[string]string{}
	// subjectManifest := []byte(`{"layers":[]}`)
	// subjectDesc := ocispec.Descriptor{
	// 	MediaType:    ocispec.MediaTypeImageManifest,
	// 	Digest:       digest.FromBytes(subjectManifest),
	// 	Size:         int64(len(subjectManifest)),
	// 	ArtifactType: artifactType,
	// 	Annotations:  annotations,
	// }

	// test Pack
	packOpts := oras.PackOptions{
		// PackImageManifest:   true,
		// Subject:             &subjectDesc,
		ManifestAnnotations: annotations,
	}

	pack := func() (ocispec.Descriptor, error) {
		root, err := oras.Pack(ctx, store, artifactType, descs, packOpts)
		if err != nil {
			return ocispec.Descriptor{}, err
		}
		if err = store.Tag(ctx, root, root.Digest.String()); err != nil {
			return ocispec.Descriptor{}, err
		}
		return root, nil
	}

	if dryRun {
		return nil
	}

	// prepare push
	dst, err := getRepository(reference)
	if err != nil {
		return err
	}

	copyOptions := oras.DefaultCopyOptions
	//copyOptions.Concurrency = opts.concurrency
	//updateDisplayOption(&copyOptions.CopyGraphOptions, store, packOpts.Verbose)
	copy := func(root ocispec.Descriptor) error {
		if tag := dst.Reference.Reference; tag == "" {
			err = oras.CopyGraph(ctx, store, dst, root, copyOptions.CopyGraphOptions)
		} else {
			_, err = oras.Copy(ctx, store, root.Digest.String(), dst, tag, copyOptions)
		}
		return err
	}

	// Push
	root, err := pushArtifact(dst, pack, &packOpts, copy, &copyOptions.CopyGraphOptions, false)
	if err != nil {
		return err
	}

	fmt.Println("Digest:", root.Digest)

	return nil
}

type packFunc func() (ocispec.Descriptor, error)
type copyFunc func(desc ocispec.Descriptor) error

func pushArtifact(dst *remote.Repository, pack packFunc, packOpts *oras.PackOptions, copy copyFunc, copyOpts *oras.CopyGraphOptions, verbose bool) (ocispec.Descriptor, error) {
	root, err := pack()
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	copyRootAttempted := false
	preCopy := copyOpts.PreCopy
	copyOpts.PreCopy = func(ctx context.Context, desc ocispec.Descriptor) error {
		if content.Equal(root, desc) {
			// copyRootAttempted helps track whether the returned error is
			// generated from copying root.
			copyRootAttempted = true
		}
		if preCopy != nil {
			return preCopy(ctx, desc)
		}
		return nil
	}

	// push
	if err = copy(root); err == nil {
		return root, nil
	}

	if !copyRootAttempted || root.MediaType != ocispec.MediaTypeArtifactManifest ||
		!isManifestUnsupported(err) {
		return ocispec.Descriptor{}, err
	}

	if err := PrintStatus(root, "Fallback ", verbose); err != nil {
		return ocispec.Descriptor{}, err
	}
	dst.SetReferrersCapability(false)
	packOpts.PackImageManifest = true
	root, err = pack()
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	copyOpts.FindSuccessors = func(ctx context.Context, fetcher content.Fetcher, node ocispec.Descriptor) ([]ocispec.Descriptor, error) {
		if content.Equal(node, root) {
			// skip non-config
			content, err := content.FetchAll(ctx, fetcher, root)
			if err != nil {
				return nil, err
			}
			var manifest ocispec.Manifest
			if err := json.Unmarshal(content, &manifest); err != nil {
				return nil, err
			}
			return []ocispec.Descriptor{manifest.Config}, nil
		}

		// config has no successors
		return nil, nil
	}
	if err = copy(root); err != nil {
		return ocispec.Descriptor{}, err
	}
	return root, nil
}

func isManifestUnsupported(err error) bool {
	var errResp *errcode.ErrorResponse
	if !errors.As(err, &errResp) || errResp.StatusCode != http.StatusBadRequest {
		return false
	}

	var errCode errcode.Error
	if !errors.As(errResp, &errCode) {
		return false
	}

	// As of November 2022, ECR is known to return UNSUPPORTED error when
	// putting an OCI artifact manifest.
	switch errCode.Code {
	case errcode.ErrorCodeManifestInvalid, errcode.ErrorCodeUnsupported:
		return true
	}
	return false
}
