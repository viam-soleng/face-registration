# Viam Face Registration Camera

A camera module which will allow you to take a picture of your face which will be stored under a configurable path on the computer the viam service runs. Previously registered faces can also be removed. This module works well in combination with the face recognition module [viam-face-identification](https://github.com/viam-labs/viam-face-identification).

### Component Configuration

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
### Attributes
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `path` | string | Yes | - | Directory path where face images will be stored |
| `camera` | string | Yes | - | Name of the camera component (e.g., "webcam") |
| `confidence` | float | No | 0.8 | Confidence threshold for face detection (0.0-1.0) |
| `detector` | string | Yes | - | Name of the face detection vision service |
| `labels` | array | No | ["person"] | List of valid labels for face detection |
| `padding` | int | No | 30 | Pixel padding around detected faces |

### Do Command 
**Add a face:**
```json
{
  "command": "add_face", 
  "name": "Your Name / ID etc."
}
```
**Remove a face:**
```json
{
  "command": "remove_face", 
  "name": "Your Name / ID etc."
}
```
### Web Client React

The [webclient-react](./webclient-react) folder contains a sample webapplication which allows capturing an image of a face and store it in a folder where the Viam server runs. This image can then be used by the [viam-face-identification](https://github.com/viam-labs/viam-face-identification) module to identify people.
A full viam machine configuration is provided in the folder.


### Common Issues

- Ensure the specified path exists and is writable
- Verify the face detection service is properly configured
- Check that the camera component is functioning