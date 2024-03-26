package selfiecamera

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"slices"
	"sort"
	"sync"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/rimage/transform"
	"go.viam.com/rdk/services/vision"
	"go.viam.com/rdk/vision/objectdetection"

	"go.viam.com/rdk/gostream"
	"go.viam.com/utils"
)

var Model = resource.ModelNamespace("viam-soleng").WithFamily("camera").WithModel("selfie-camera")

type Config struct {
	Camera     string
	Detector   string
	Confidence float64
	Path       string
	Labels     []string
	Padding    int
}

func (cfg *Config) Validate(path string) ([]string, error) {
	if cfg.Camera == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "camera")
	}

	if cfg.Detector == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "detector")
	}

	if cfg.Path == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "path")
	}

	return []string{cfg.Camera, cfg.Detector}, nil
}

func init() {
	resource.RegisterComponent(camera.API, Model, resource.Registration[camera.Camera, *Config]{
		Constructor: func(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (camera.Camera, error) {
			newConf, err := resource.NativeConfig[*Config](conf)
			if err != nil {
				return nil, err
			}
			fc := &selfieCamera{name: conf.ResourceName(), conf: newConf, logger: logger}
			fc.camera, err = camera.FromDependencies(deps, newConf.Camera)
			if err != nil {
				return nil, err
			}
			fc.detector, err = vision.FromDependencies(deps, newConf.Detector)
			if err != nil {
				return nil, err
			}
			fc.conf.Confidence = newConf.Confidence
			fc.labels = newConf.Labels
			fc.padding = newConf.Padding
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

	camera     camera.Camera
	detector   vision.Service
	confidence float64
	labels     []string
	padding    int
	path       string

	image image.Image

	mu sync.Mutex
}

func (sc *selfieCamera) Name() resource.Name {
	return sc.name
}

func (sc *selfieCamera) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	name, ok := cmd["name"].(string)
	if ok {
		image, err := sc.takeSelfie(ctx, name)
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
	images, meta, err := sc.camera.Images(ctx)
	if err != nil {
		return images, meta, err
	}
	return images, meta, nil
}

func (sc *selfieCamera) Stream(ctx context.Context, errHandlers ...gostream.ErrorHandler) (gostream.VideoStream, error) {
	camStream, err := sc.camera.Stream(ctx, errHandlers...)
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
	p, err := sc.camera.Properties(ctx)
	if err == nil {
		p.SupportsPCD = false
	}
	return p, err
}

func (fc *selfieCamera) Projector(ctx context.Context) (transform.Projector, error) {
	return fc.camera.Projector(ctx)
}

func (sc *selfieCamera) takeSelfie(ctx context.Context, name string) (image.Image, error) {
	// Get bounding box from vision service
	detections, err := sc.detectFace(ctx, sc.image)
	if err != nil {
		return nil, err
	}
	if len(detections) == 0 {
		return nil, errors.New("no face detected")
	}
	// Crop image
	croppedImage, err := cropImage(sc.image, detections[0], sc.padding)
	if err != nil {
		return nil, err
	}
	// Store cropped image under path
	if croppedImage != nil {
		err := saveImage(croppedImage, name, sc.path)
		if err != nil {
			return nil, err
		} else {
			return sc.image, nil
		}
	}
	return nil, errors.New("image buffer empty, activate camera stream")

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

// Take an input image, detect objects, crop the image down to the detected bounding box and
// hand over to classifier for more accurate classifications
func (sc *selfieCamera) detectFace(ctx context.Context, img image.Image) ([]objectdetection.Detection, error) {
	// Get detections from the provided Image
	detections, err := sc.detector.Detections(ctx, img, nil)
	if err != nil {
		return nil, err
	}
	// Filter detections by detector confidence level and valid labels settings
	filterFunc := func(detection objectdetection.Detection) bool {
		return (detection.Score() >= sc.confidence) && (slices.Contains(sc.labels, detection.Label()) || len(sc.labels) == 0)
	}
	detections = filter(detections, filterFunc)

	// Sort filtered detections based upon score
	sort.Slice(detections, func(i, j int) bool {
		return detections[i].Score() > detections[j].Score()
	})
	return detections, nil
}

// Generic helper function to filter slices
func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// Crops images based upon bounding box and padding
func cropImage(img image.Image, detection objectdetection.Detection, padding int) (image.Image, error) {
	// The cropping operation is done by creating a new image of the size of the rectangle
	// and drawing the relevant part of the original image onto the new image.
	// Increase/decrease bounding box according to detection border setting
	rectangle := image.Rect(
		detection.BoundingBox().Min.X-padding,
		detection.BoundingBox().Min.Y-padding,
		detection.BoundingBox().Max.X+padding,
		detection.BoundingBox().Max.Y+padding)

	cropped := image.NewRGBA(rectangle.Bounds())
	draw.Draw(cropped, rectangle.Bounds(), img, rectangle.Min, draw.Src)
	return cropped, nil
}
