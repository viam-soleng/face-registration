import {
  useState,
  type ChangeEventHandler,
  type FormEventHandler,
} from 'react';

import {
  DISCONNECTED,
  CONNECTING,
  DISCONNECTING,
  CONNECTED,
  type ClientStatus,
} from '../state.js';

import type { RobotCredentials } from '../client.js';
import React from 'react';

export interface ConnectFormProps {
  status: ClientStatus;
  onSubmit: (credentials: RobotCredentials) => unknown;
}

const DISABLED_BY_STATUS = {
  [DISCONNECTED]: false,
  [CONNECTING]: true,
  [DISCONNECTING]: true,
  [CONNECTED]: false,
};

const BUTTON_TEXT_BY_STATUS = {
  [DISCONNECTED]: 'Connect',
  [CONNECTING]: 'Connecting...',
  [DISCONNECTING]: 'Disconnecting...',
  [CONNECTED]: 'Disconnect',
};

const INITIAL_HOSTNAME = import.meta.env.VITE_ROBOT_HOSTNAME ?? '';
const INITIAL_KEY_ID = import.meta.env.VITE_ROBOT_KEY_ID ?? '';
const INITIAL_KEY_VALUE = import.meta.env.VITE_ROBOT_KEY_VALUE ?? '';

export const ConnectForm = (props: ConnectFormProps): JSX.Element => {
  const { status, onSubmit } = props;
  const [hostname, setHostname] = useState(INITIAL_HOSTNAME);
  const [key_id, setKeyID] = useState(INITIAL_KEY_ID);
  const [key_value, setKeyValue] = useState(INITIAL_KEY_VALUE);
  const disabled = DISABLED_BY_STATUS[status];
  const buttonText = BUTTON_TEXT_BY_STATUS[status];

  const handleHost: ChangeEventHandler<HTMLInputElement> = (event) => {
    setHostname(event.target.value);
  };
  const handleSecret: ChangeEventHandler<HTMLInputElement> = (event) => {
    setKeyValue(event.target.value);
  };
  const handleSubmit: FormEventHandler = (event) => {
    onSubmit({ hostname, key_id, key_value });
    event.preventDefault();
  };

  if (status === "connected") {
    return (
      <form className="flex flex-col p-4 w-96" onSubmit={handleSubmit}>
        <button
          type="submit"
          disabled={disabled}
          className=" w-32 border-solid border-2 border-black"
        >
          {buttonText}
        </button>
      </form>
    )
  } else {
    return (
      <form className="flex flex-col p-4 w-96" onSubmit={handleSubmit}>
        <label className="flex flex-col mb-1">
          Remote Address
          <input
            type="text"
            className="px-1 border-solid border-2 border-black"
            value={hostname}
            onChange={handleHost}
            disabled={disabled}
          />
        </label>
        <label className="flex flex-col mb-6">
          Location Secret
          <input
            type="password"
            className="px-1 border-solid border-2 border-black"
            value={key_value}
            onChange={handleSecret}
            disabled={disabled}
          />
        </label>
        <button
          type="submit"
          disabled={disabled}
          className=" w-32 border-solid border-2 border-black"
        >
          {buttonText}
        </button>
      </form>
    );
  }
};
