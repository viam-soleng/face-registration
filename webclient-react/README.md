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
  "agent_config": {
    "subsystems": {
      "viam-agent": {
        "release_channel": "stable",
        "pin_version": "",
        "pin_url": "",
        "disable_subsystem": false
      },
      "viam-server": {
        "pin_version": "",
        "pin_url": "",
        "disable_subsystem": false,
        "release_channel": "stable"
      },
      "agent-provisioning": {
        "release_channel": "stable",
        "pin_version": "",
        "pin_url": "",
        "disable_subsystem": false
      },
      "agent-syscfg": {
        "pin_url": "",
        "disable_subsystem": false,
        "release_channel": "stable",
        "pin_version": ""
      }
    }
  },
  "components": [
    {
      "namespace": "rdk",
      "type": "camera",
      "model": "webcam",
      "attributes": {},
      "name": "camera",
      "depends_on": []
    },
    {
      "name": "tf-face-detect",
      "model": "transform",
      "type": "camera",
      "namespace": "rdk",
      "attributes": {
        "source": "camera",
        "pipeline": [
          {
            "type": "detections",
            "attributes": {
              "detector_name": "vis-face-detect",
              "confidence_threshold": 0.5,
              "valid_labels": [
                "face"
              ]
            }
          }
        ]
      },
      "depends_on": []
    },
    {
      "name": "tf-identification",
      "model": "transform",
      "type": "camera",
      "namespace": "rdk",
      "attributes": {
        "source": "camera",
        "pipeline": [
          {
            "type": "detections",
            "attributes": {
              "detector_name": "vis-identification",
              "confidence_threshold": 0.5
            }
          }
        ]
      },
      "depends_on": []
    },
    {
      "name": "face-camera",
      "model": "viam-soleng:camera:face-camera",
      "type": "camera",
      "namespace": "rdk",
      "attributes": {
        "padding": 30,
        "path": "/Users/username/Downloads/faces",
        "camera": "camera",
        "confidence": 0.5,
        "detector": "vis-face-detect",
        "labels": [
          "face"
        ]
      },
      "depends_on": []
    }
  ],
  "services": [
    {
      "type": "mlmodel",
      "model": "viam-labs:mlmodel:onnx-cpu",
      "attributes": {
        "model_path": "${packages.face-detector-onnx}/face_detector_640.onnx",
        "label_path": "${packages.face-detector-onnx}/face_labels.txt",
        "num_threads": 1,
        "package_reference": "viam-soleng/face-detector-onnx"
      },
      "name": "ml-face-detect",
      "namespace": "rdk"
    },
    {
      "name": "vis-face-detect",
      "type": "vision",
      "model": "mlmodel",
      "attributes": {
        "mlmodel_name": "ml-face-detect",
        "remap_input_names": {
          "input": "image"
        },
        "remap_output_names": {
          "boxes": "location",
          "scores": "score"
        },
        "xmin_ymin_xmax_ymax_order": [
          0,
          1,
          2,
          3
        ]
      }
    },
    {
      "name": "vis-identification",
      "type": "vision",
      "namespace": "rdk",
      "model": "viam:vision:deepface-identification",
      "attributes": {
        "camera_name": "camera",
        "picture_directory": "/Users/username/Downloads/faces"
      }
    }
  ],
  "modules": [
    {
      "version": "0.0.1",
      "type": "registry",
      "name": "viam-soleng_face-registration",
      "module_id": "viam-soleng:face-registration"
    },
    {
      "type": "registry",
      "name": "viam_deepface-identification",
      "module_id": "viam:deepface-identification",
      "version": "0.2.4"
    },
    {
      "module_id": "viam-labs:onnx-cpu",
      "version": "0.1.2",
      "type": "registry",
      "name": "viam-labs_onnx-cpu"
    }
  ],
  "packages": [
    {
      "version": "latest",
      "name": "face-detector-onnx",
      "package": "viam-soleng/face-detector-onnx",
      "type": "ml_model"
    }
  ]
}
```


