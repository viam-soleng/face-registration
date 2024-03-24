package selfiecamera

import (
	"context"
	"fmt"
	"image"
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
	Camera  string
	Vision  string
	Objects map[string]float64
}

func (cfg *Config) Validate(path string) ([]string, error) {
	if cfg.Camera == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "camera")
	}

	if cfg.Vision == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "vision")
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

	cam camera.Camera
	vis vision.Service

	mu sync.Mutex
}

func (fc *selfieCamera) Name() resource.Name {
	return fc.name
}

func (fc *selfieCamera) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, resource.ErrDoUnimplemented
}

func (fc *selfieCamera) Images(ctx context.Context) ([]camera.NamedImage, resource.ResponseMetadata, error) {
	images, meta, err := fc.cam.Images(ctx)
	if err != nil {
		return images, meta, err
	}
	return images, meta, nil
}

func (fc *selfieCamera) Stream(ctx context.Context, errHandlers ...gostream.ErrorHandler) (gostream.VideoStream, error) {
	camStream, err := fc.cam.Stream(ctx, errHandlers...)
	if err != nil {
		return nil, err
	}

	return sourceStream{camStream, fc}, nil
}

type sourceStream struct {
	cameraStream gostream.VideoStream
	fc           *selfieCamera
}

func (fs sourceStream) Next(ctx context.Context) (image.Image, func(), error) {
	return fs.cameraStream.Next(ctx)
}

func (fs sourceStream) Close(ctx context.Context) error {
	return fs.cameraStream.Close(ctx)
}

func (fc *selfieCamera) NextPointCloud(ctx context.Context) (pointcloud.PointCloud, error) {
	return nil, fmt.Errorf("selfieCamera doesn't support pointclouds")
}

func (fc *selfieCamera) Properties(ctx context.Context) (camera.Properties, error) {
	p, err := fc.cam.Properties(ctx)
	if err == nil {
		p.SupportsPCD = false
	}
	return p, err
}

func (fc *selfieCamera) Projector(ctx context.Context) (transform.Projector, error) {
	return fc.cam.Projector(ctx)
}
