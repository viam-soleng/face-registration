import { VideoStream } from "./components/video-stream.js";
import { ConnectForm } from "./components/connect-form.js";
import { useStore, useStream } from "./state.js";
import React from "react";
import { SelfieCamera } from "./components/face-camera.js";

export const App = (): JSX.Element => {
  const {
    status,
    connectOrDisconnect,
    streamClient,
    cameraClient,
    identificationClient,
  } = useStore();
  const stream = useStream(streamClient, "face-camera");

  return (
    <>
      <ConnectForm status={status} onSubmit={connectOrDisconnect} />
      {status === "connected" ? (
        <>
          <SelfieCamera
            cameraClient={cameraClient}
            identificationClient={identificationClient}
          ></SelfieCamera>
          <VideoStream stream={stream}></VideoStream>
        </>
      ) : (
        <></>
      )}
    </>
  );
};
