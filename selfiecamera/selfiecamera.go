package selfiecamera

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"sync"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/rimage/transform"
	"go.viam.com/rdk/services/vision"

	"go.viam.com/rdk/gostream"
	"go.viam.com/utils"
)

var Model = resource.ModelNamespace("viam-soleng").WithFamily("camera").WithModel("selfie-camera")

type Config struct {
	Camera string
	Vision string
	Path   string
}

func (cfg *Config) Validate(path string) ([]string, error) {
	if cfg.Camera == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "camera")
	}

	if cfg.Vision == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "vision")
	}

	if cfg.Path == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "path")
	}

	return []string{cfg.Camera, cfg.Vision}, nil
}

func init() {
	resource.RegisterComponent(camera.API, Model, resource.Registration[camera.Camera, *Config]{
		Constructor: func(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (camera.Camera, error) {
			newConf, err := resource.NativeConfig[*Config](conf)
			if err != nil {
				return nil, err
			}
			fc := &selfieCamera{name: conf.ResourceName(), conf: newConf, logger: logger}
			fc.cam, err = camera.FromDependencies(deps, newConf.Camera)
			if err != nil {
				return nil, err
			}
			fc.vis, err = vision.FromDependencies(deps, newConf.Vision)
			if err != nil {
				return nil, err
			}
			fc.path = newConf.Path
			return fc, nil
		},
	})
}

type selfieCamera struct {
	resource.AlwaysRebuild
	resource.TriviallyCloseable

	name   resource.Name
	conf   *Config
	logger logging.Logger

	cam  camera.Camera
	vis  vision.Service
	path string

	image image.Image

	mu sync.Mutex
}

func (sc *selfieCamera) Name() resource.Name {
	return sc.name
}

func (sc *selfieCamera) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	name, ok := cmd["name"].(string)
	if ok {
		image, err := sc.takeSelfie(name)
		if err != nil {
			return nil, err
		} else {
			return map[string]any{"image": image}, nil
		}
	} else {
		return nil, errors.New(`"name" value must be string`)
	}
}

func (sc *selfieCamera) Images(ctx context.Context) ([]camera.NamedImage, resource.ResponseMetadata, error) {
	images, meta, err := sc.cam.Images(ctx)
	if err != nil {
		return images, meta, err
	}
	return images, meta, nil
}

func (sc *selfieCamera) Stream(ctx context.Context, errHandlers ...gostream.ErrorHandler) (gostream.VideoStream, error) {
	camStream, err := sc.cam.Stream(ctx, errHandlers...)
	if err != nil {
		return nil, err
	}
	image, release, err := camStream.Next(ctx)
	if err != nil {
		return nil, err
	}
	defer release()
	sc.image = image
	return sourceStream{camStream, sc}, nil
}

type sourceStream struct {
	cameraStream gostream.VideoStream
	fc           *selfieCamera
}

func (sc sourceStream) Next(ctx context.Context) (image.Image, func(), error) {
	return sc.cameraStream.Next(ctx)
}

func (sc sourceStream) Close(ctx context.Context) error {
	return sc.cameraStream.Close(ctx)
}

func (sc *selfieCamera) NextPointCloud(ctx context.Context) (pointcloud.PointCloud, error) {
	return nil, fmt.Errorf("Selfie-Camera doesn't support pointclouds")
}

func (sc *selfieCamera) Properties(ctx context.Context) (camera.Properties, error) {
	p, err := sc.cam.Properties(ctx)
	if err == nil {
		p.SupportsPCD = false
	}
	return p, err
}

func (fc *selfieCamera) Projector(ctx context.Context) (transform.Projector, error) {
	return fc.cam.Projector(ctx)
}

func (sc *selfieCamera) takeSelfie(name string) (image.Image, error) {
	sc.logger.Infof("And the name is: %s", name)
	// get image from camera

	// Get bounding box from vision service

	// Crop Face

	// Store cropped image under path
	if sc.image != nil {
		err := saveImage(sc.image, name, sc.path)
		if err != nil {
			return nil, err
		} else {
			return sc.image, nil
		}
	}
	return nil, errors.New("image buffer empty, activate camera first")

}

// Saves images to a path on disk
func saveImage(image image.Image, name string, path string) error {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)
	if err != nil {
		return err
	}
	digest := sha256.New()
	digest.Write(buf.Bytes())
	hash := digest.Sum(nil)
	f, err := os.Create(fmt.Sprintf("%s/%s_%x.jpg", path, name, hash))
	if err != nil {
		return err
	}
	defer f.Close()
	opt := jpeg.Options{
		Quality: 90,
	}
	jpeg.Encode(f, image, &opt)
	return nil
}
