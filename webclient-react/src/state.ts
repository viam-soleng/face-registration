import { useEffect, useRef, useState } from 'react';
import type { RobotClient, StreamClient, CameraClient, VisionClient } from '@viamrobotics/sdk';
import {
  getRobotClient,
  getStreamClient,
  type RobotCredentials,
  getCameraClient,
  getIdentificationClient,
} from './client.js';

export const DISCONNECTED = 'disconnected';
export const CONNECTING = 'connecting';
export const DISCONNECTING = 'disconnecting';
export const CONNECTED = 'connected';

interface ClientStateDisconnected {
  status: typeof DISCONNECTED;
  error?: unknown;
}

interface ClientStateTransitioning {
  status: typeof CONNECTING | typeof DISCONNECTING;
}

interface ClientStateConnected {
  status: typeof CONNECTED;
  client: RobotClient;
  streamClient: StreamClient;
  cameraClient: CameraClient;
  identificationClient: VisionClient;
}

type ClientState =
  | ClientStateDisconnected
  | ClientStateTransitioning
  | ClientStateConnected;

export type ClientStatus = ClientState['status'];

export interface Store {
  status: ClientStatus;
  client?: RobotClient;
  streamClient?: StreamClient;
  cameraClient?: CameraClient;
  identificationClient?: VisionClient;
  connectOrDisconnect: (credentials: RobotCredentials) => unknown;
}

export const useStore = (): Store => {
  const [state, setState] = useState<ClientState>({ status: DISCONNECTED });

  if (state.status === DISCONNECTED && state.error) {
    console.warn('Connection error', state.error);
  }

  const connectOrDisconnect = (credentials: RobotCredentials): void => {
    if (state.status === DISCONNECTED) {
      setState({ status: CONNECTING });

      getRobotClient(credentials)
        .then((client) => {
          const streamClient = getStreamClient(client);
          const cameraClient = getCameraClient(client);
          const identificationClient = getIdentificationClient(client);
          setState({ status: CONNECTED, client, streamClient, cameraClient, identificationClient });
        })
        .catch((error: unknown) => setState({ status: DISCONNECTED, error }));
    } else if (state.status === CONNECTED) {
      setState({ status: DISCONNECTING });

      state.client
        .disconnect()
        .then(() => setState({ status: DISCONNECTED }))
        .catch((error: unknown) => setState({ status: DISCONNECTED, error }));
    }
  };

  return {
    connectOrDisconnect,
    status: state.status,
    client: state.status === CONNECTED ? state.client : undefined,
    streamClient: state.status === CONNECTED ? state.streamClient : undefined,
    cameraClient: state.status === CONNECTED ? state.cameraClient : undefined,
    identificationClient: state.status === CONNECTED ? state.identificationClient : undefined,
  };
};

export const useStream = (
  streamClient: StreamClient | undefined,
  cameraName: string
): MediaStream | undefined => {
  const okToConnectRef = useRef(true);
  const [stream, setStream] = useState<MediaStream | undefined>();

  useEffect(() => {
    if (streamClient && okToConnectRef.current) {
      okToConnectRef.current = false;

      streamClient
        .getStream(cameraName)
        .then((mediaStream) => setStream(mediaStream))
        .catch((error: unknown) => {
          console.warn(`Unable to connect to camera ${cameraName}`, error);
        });

      return () => {
        okToConnectRef.current = true;

        streamClient.remove(cameraName).catch((error: unknown) => {
          console.warn(`Unable to disconnect to camera ${cameraName}`, error);
        });
      };
    }

    return undefined;
  }, [streamClient, cameraName]);

  return stream;
};
