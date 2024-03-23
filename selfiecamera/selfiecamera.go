package selfiecamera

import (
	"context"
	"image"
	"sync"

	"github.com/pkg/errors"
	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/gostream"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/rimage/transform"
	"go.viam.com/rdk/services/vision"
)

var errUnimplemented = errors.New("unimplemented")
var Model = resource.NewModel("viam-soleng", "camera", "selfiecamera")
var PrettyName = "Viam selfie camera"
var Description = "A Viam camera component module allowing people to take a selfie of their face"

func init() {
	resource.RegisterComponent(
		camera.API,
		Model,
		resource.Registration[camera.Camera, *Config]{
			Constructor: newCamera,
		})
}

func newCamera(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (camera.Camera, error) {
	logger.Debugf("Starting %s %s", PrettyName)

	newConf, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return nil, err
	}
	of := &selfieCamera{name: conf.ResourceName(), conf: newConf, logger: logger}
	of.srcCamera, err = camera.FromDependencies(deps, newConf.SrcCamera)
	if err != nil {
		return nil, err
	}

	camera := selfieCamera{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
		mu:     sync.RWMutex{},
		done:   make(chan bool),
	}

	/*
		if err := camera.Reconfigure(ctx, deps, conf); err != nil {
			return nil, err
		}
	*/
	return &camera, nil
}

type Config struct {
	SrcCamera          string  `json:"src_camera"`
	Detector           string  `json:"detector_service"`
	DetectorConfidence float64 `json:"detector_confidence"`
	BBoxPadding        int     `json:"padding"`
	Path               string  `json:"path"`
}

func (cfg *Config) Validate(path string) ([]string, error) {
	return []string{cfg.SrcCamera}, nil
}

type selfieCamera struct {
	resource.Named
	resource.AlwaysRebuild

	name               resource.Name
	conf               *Config
	logger             logging.Logger
	srcCamera          camera.Camera
	detector           vision.Service
	detectorConfidence float64
	bboxPadding        int
	path               string
	mu                 sync.RWMutex
	cancelCtx          context.Context
	cancelFunc         func()
	done               chan bool
}

/*
// Reconfigure reconfigures with new settings.
func (sc *selfieCamera) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.logger.Debugf("Reconfiguring %s", PrettyName)
	// In case the module has changed name
	sc.Named = conf.ResourceName().AsNamed()
	newConf, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return err
	}
	// Get the camera
	sc.srcCamera, err = camera.FromDependencies(deps, newConf.SrcCamera)
	if err != nil {
		return errors.Wrapf(err, "unable to get source camera %v for image sourcing...", newConf.Detector)
	}
	sc.logger.Debug("**** Reconfigured ****")
	return nil
}
*/

// Images implements camera.Camera.
func (sc *selfieCamera) Images(ctx context.Context) ([]camera.NamedImage, resource.ResponseMetadata, error) {
	panic("unimplemented")
}

// Name implements camera.Camera.
// Subtle: this method shadows the method (Named).Name of selfieCamera.Named.
func (sc *selfieCamera) Name() resource.Name {
	panic("unimplemented")
}

// NextPointCloud implements camera.Camera.
func (sc *selfieCamera) NextPointCloud(ctx context.Context) (pointcloud.PointCloud, error) {
	panic("unimplemented")
}

// Projector implements camera.Camera.
func (sc *selfieCamera) Projector(ctx context.Context) (transform.Projector, error) {
	panic("unimplemented")
}

// Properties implements camera.Camera.
func (sc *selfieCamera) Properties(ctx context.Context) (camera.Properties, error) {
	panic("unimplemented")
}

// Stream implements camera.Camera.
func (sc *selfieCamera) Stream(ctx context.Context, errHandlers ...gostream.ErrorHandler) (gostream.MediaStream[image.Image], error) {

	// gets the stream from a camera
	stream, err := sc.srcCamera.Stream(context.Background(), errHandlers...)
	if err != nil {
		return nil, err
	}
	defer stream.Close(ctx)
	return srcCamStream{stream}, nil
}

// DoCommand can be implemented to extend functionality but returns unimplemented currently.
func (sc *selfieCamera) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, errUnimplemented
}

// The close method is executed when the component is shut down
func (sc *selfieCamera) Close(ctx context.Context) error {
	sc.logger.Debugf("Shutting down %s", PrettyName)
	return nil
}

type srcCamStream struct {
	cameraStream gostream.VideoStream
}

// Close implements gostream.MediaStream.
func (sc srcCamStream) Close(ctx context.Context) error {
	panic("unimplemented")
}

// Next implements gostream.MediaStream.
func (scs srcCamStream) Next(ctx context.Context) (image.Image, func(), error) {
	// Get next camera img
	img, release, err := scs.cameraStream.Next(ctx)
	if err != nil {
		return nil, nil, err
	}
	// return raw image
	return img, release, nil
}
