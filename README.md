# Viam Selfie Camera

A camera module which will allow you to take a picture of your face to be stored under a configurable path for further use by for example a face recognition module such as [viam-face-identification](https://github.com/viam-labs/viam-face-identification).

## Component Configuration

```json
{
  "path": "/<- YOUR PATH->/",
  "camera": "camera",
  "confidence": 0.8,
  "detector": "<- YOUR FACE DETECTION VISION SERVICE",
  "labels": [
    "valid labels"
  ],
  "padding": 30
}
```
## Do Command Input

```json
{
  "name": "Maybe Your Name"
}
```
## Web Client React

The [webclient-react]() folder contains a sample webapplication which allows capturing an image of a face and store it in a folder where the Viam server runs. This image can then be used by the [viam-face-identification](https://github.com/viam-labs/viam-face-identification) module to identify people.
A full viam machine configuration is provided in the folder.
