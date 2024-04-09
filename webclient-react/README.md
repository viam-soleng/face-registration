# Selfie-Camera Webclient-React

This example project allows you to take pictures of your face and store them in a local folder on the device where the Viam server runs.

## Setup

First, follow the setup instructions for the repository in `CONTRIBUTING.md`. Then, install development dependencies for the demo and launch a dev server.

```shell
npm install
```

The connection hostname and secret fields can be pre-filled from a `.env` file in the `webclient-react` directory. You have to set these before running npm start. 

```ini
# ../webclient-react/.env
VITE_ROBOT_HOSTNAME=<-YOUR HOSTNAME->
VITE_ROBOT_KEY_ID=<-YOUR KEY ID->
VITE_ROBOT_KEY_VALUE=<- YOUR KEY VALUE->
```

```shell
npm start
```

### Viam Machine Configuration

```json

{
  "components": [
    {
      "name": "selfiecamera",
      "namespace": "rdk",
      "type": "camera",
      "model": "viam-soleng:camera:selfie-camera",
      "attributes": {
        "confidence": 0.5,
        "detector": "vision-face",
        "labels": [
          "face"
        ],
        "padding": 30,
        "path": "/Users/Username/Documents/GitHub/selfie-camera/selfies",
        "camera": "camera"
      },
      "depends_on": [
        "camera",
        "vision-face"
      ]
    },
    {
      "name": "camera",
      "namespace": "rdk",
      "type": "camera",
      "model": "webcam",
      "attributes": {
        "video_path": "FDF90FEB-59E5-4FCF-AABD-DA03C4E19BFB"
      }
    },
    {
      "name": "camera-transform",
      "namespace": "rdk",
      "type": "camera",
      "model": "transform",
      "attributes": {
        "pipeline": [
          {
            "type": "detections",
            "attributes": {
              "detector_name": "vision-face",
              "valid_labels": [
                "face"
              ],
              "confidence_threshold": 0.8
            }
          }
        ],
        "source": "selfiecamera"
      }
    },
    {
      "name": "camera-face-identification",
      "namespace": "rdk",
      "type": "camera",
      "model": "transform",
      "attributes": {
        "pipeline": [
          {
            "attributes": {
              "confidence_threshold": 0.5,
              "detector_name": "vision-deepface"
            },
            "type": "detections"
          }
        ],
        "source": "camera"
      }
    }
  ],
  "services": [
    {
      "name": "vision-face",
      "namespace": "rdk",
      "type": "vision",
      "model": "mlmodel",
      "attributes": {
        "xmin_ymin_xmax_ymax_order": [
          0,
          1,
          2,
          3
        ],
        "mlmodel_name": "model-face",
        "remap_input_names": {
          "input": "image"
        },
        "remap_output_names": {
          "scores": "score",
          "boxes": "location"
        }
      }
    },
    {
      "name": "model-face",
      "namespace": "rdk",
      "type": "mlmodel",
      "model": "viam-labs:mlmodel:onnx-cpu",
      "attributes": {
        "label_path": "${packages.face-detector-onnx}/face_labels.txt",
        "model_path": "${packages.face-detector-onnx}/face_detector_640.onnx",
        "num_threads": 1,
        "package_reference": "viam-soleng/face-detector-onnx"
      }
    }
  ],
  "modules": [
    {
      "type": "registry",
      "name": "viam-soleng_selfie-camera",
      "module_id": "viam-soleng:selfie-camera",
      "version": "0.0.3"
    }
  ]
}
```


