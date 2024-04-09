import { CameraClient, VisionClient } from '@viamrobotics/sdk';
import React, { ReactNode, useState } from 'react';

export interface SelfieCameraProps {
    cameraClient?: CameraClient;
    identificationClient?: VisionClient;
    children?: ReactNode;
}

export const SelfieCamera = (props: SelfieCameraProps): JSX.Element => {

    const { cameraClient, identificationClient } = props;
    const [name, setName] = useState('');

    async function takeSelfie() {
        if ((cameraClient != undefined) && (identificationClient != undefined)) {
            if (name === "") {
                alert("Please provide a name!")
            } else {
                let result;
                try {
                    await cameraClient?.doCommand({ "name": name });
                    await identificationClient?.doCommand({"command": "recompute_embeddings"})
                    alert(name + " has been registered!");
                }
                catch (err: any) {
                    console.log(JSON.stringify(err));
                    alert(err.metadata.headersMap["grpc-message"]);
                }
            }
        } else {
            alert("Camera or IdentificationClient undefined!" + cameraClient + identificationClient);
        }

    }

    return (
        <div className="flex flex-row p-4 w-96">
            <input
                name="name"
                className=" w-32 border-solid border-2 border-black"
                value={name}
                onChange={e => setName(e.target.value)}
            />
            <button
                onClick={async () => { await takeSelfie(); }}
                className=" w-32 border-solid border-2 border-black"
            >
                {"Take Selfie"}
            </button>
        </div>
    );
};
