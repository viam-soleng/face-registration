import {
  createRobotClient,
  StreamClient,
  CameraClient,
  type RobotClient,
  VisionClient,
} from '@viamrobotics/sdk';

export interface RobotCredentials {
  hostname: string;
  key_id: string;
  key_value: string;
}

/**
 * Given a set of credentials, get a robot client.
 *
 * @param credentials Robot URL and location secret
 * @returns A connected client
 */
export const getRobotClient = async (
  credentials: RobotCredentials
): Promise<RobotClient> => {

  const { hostname, key_id, key_value } = credentials;

  return createRobotClient({
    host:hostname,
    credential: {
      type: 'api-key' ,
      payload: key_value,
    } ,
    authEntity: key_id,
    signalingAddress: 'https://app.viam.com:443',
  });
};

/**
 * StreamClient factory
 *
 * @param client A connected RobotClient
 * @returns A connected stream client
 */
export const getStreamClient = (client: RobotClient): StreamClient => {
  return new StreamClient(client);
};

/**
 * BaseClient factory
 *
 * @param client A connected RobotClient
 * @returns A connected camera client
 */
export const getCameraClient = (client: RobotClient): CameraClient => {
  return new CameraClient(client, 'face-camera');
};

/**
 * BaseClient factory
 *
 * @param client A connected RobotClient
 * @returns A connected vision client
 */
export const getIdentificationClient = (client: RobotClient): VisionClient => {
  return new VisionClient(client, 'vis-identification');
};