import { CameraClient, Struct, VisionClient } from "@viamrobotics/sdk";
import React, { ReactNode, useState } from "react";

export interface SelfieCameraProps {
  cameraClient?: CameraClient;
  identificationClient?: VisionClient;
  children?: ReactNode;
}

export const SelfieCamera = (props: SelfieCameraProps): JSX.Element => {
  const { cameraClient, identificationClient } = props;
  const [name, setName] = useState("");

  async function addFace() {
    if (cameraClient != undefined && identificationClient != undefined) {
      if (name === "") {
        alert("Please provide a name!");
      } else {
        let result;
        try {
          await cameraClient?.doCommand(
            Struct.fromJson({ command: "add_face", name: name })
          );
          //await cameraClient?.doCommand({ "command": "add_face", "name": name });
          // TODO: waiting for merge of PR15
          await identificationClient?.doCommand(
            Struct.fromJson({ command: "recompute_embeddings" })
          );
          alert(name + " has been registered!");
        } catch (err: any) {
          console.log(JSON.stringify(err));
          alert(err.cause);
        }
      }
    } else {
      alert(
        "Camera or IdentificationClient undefined!" +
          cameraClient +
          identificationClient
      );
    }
  }

  async function removeFace() {
    if (cameraClient != undefined && identificationClient != undefined) {
      if (name === "") {
        alert("Please provide a name!");
      } else {
        try {
          let cmd = Struct.fromJson({ command: "remove_face", name: name });
          await cameraClient?.doCommand(
            Struct.fromJson({ command: "remove_face", name: name })
          );
          // TODO: waiting for merge of PR15
          await identificationClient?.doCommand(
            Struct.fromJson({ command: "recompute_embeddings" })
          );
          alert(name + " has been removed!");
        } catch (err: any) {
          console.log(JSON.stringify(err));
          alert(err.cause);
        }
      }
    } else {
      alert(
        "Camera or IdentificationClient undefined!" +
          cameraClient +
          identificationClient
      );
    }
  }

  return (
    <div className="flex flex-row p-4 w-96">
      <input
        name="name"
        className=" w-32 border-solid border-2 border-black"
        value={name}
        onChange={(e) => setName(e.target.value)}
      />
      <button
        onClick={async () => {
          await addFace();
        }}
        className=" w-32 border-solid border-2 border-black"
      >
        {"Add Face"}
      </button>
      <button
        onClick={async () => {
          await removeFace();
        }}
        className=" w-32 border-solid border-2 border-black"
      >
        {"Remove Face"}
      </button>
    </div>
  );
};
