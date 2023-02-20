import React, { ChangeEvent } from 'react';
import { InlineField, Switch, Input, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { AmqpDataSourceOptions, AmqpSecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<AmqpDataSourceOptions> {}


/*
			SetHost("localhost").
			SetPort(5552).
			SetVHost("/").
			SetUser("guest").
			SetPassword("guest").
			SetMaxProducersPerClient(1).
			SetMaxConsumersPerClient(1).
			IsTLS(false).
			// SetTLSConfig(&tls.Config{}).
			SetRequestedHeartbeat(60 * time.Second).
			SetRequestedMaxFrameSize(1048576).
			SetWriteBuffer(8192).
			SetReadBuffer(65536).
			SetNoDelay(false),
*/

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const onHostChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      path: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };
  /*const onPortChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      path: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };
  const onUsernameChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      path: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };
  const onTlsConnectionChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      path: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  const onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        password: event.target.value,
      },
    });
  };*/

  const onResetPassword = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        password: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        password: '',
      },
    });
  };

  const { jsonData, secureJsonFields } = options;
  const secureJsonData = (options.secureJsonData || {}) as AmqpSecureJsonData;

  const labelWidth = 20;
  const inputWidth = 40;
  const switchWidth = 10;

  const amqpDefaultHost = "localhost";
  const amqpDefaultPort = 5672;
  const amqpDefaultTls = false;
  const amqpDefaultUser = "guest";
  const amqpDefaultPassword = "guest"

  return (
    <div className="gf-form-group">
      <InlineField label="Host" labelWidth={labelWidth} tooltip="Host of the AMQP server">
        <Input
          onChange={onHostChange}
          value={jsonData.host || amqpDefaultHost}
          placeholder={amqpDefaultHost}
          width={inputWidth}
        />
      </InlineField>
      <InlineField label="Port" labelWidth={labelWidth} tooltip="Port of the AMQP server">
        <Input
          onChange={onHostChange}
          value={jsonData.port || amqpDefaultPort}
          placeholder={amqpDefaultPort.toString()}
          width={inputWidth}
        />
      </InlineField>
      <InlineField label="Username" labelWidth={labelWidth} tooltip="Username to connect to the AMQP server">
        <Input
          onChange={onHostChange}
          value={jsonData.username || amqpDefaultUser}
          placeholder={amqpDefaultUser}
          width={inputWidth}
        />
      </InlineField>
      <InlineField label="TlsConnection" labelWidth={labelWidth} tooltip="Should use TLS to connect to the AMQP server">
        <Switch
          onChange={onHostChange}
          value={jsonData.tlsConnection || amqpDefaultTls}
          width={switchWidth}
        />
      </InlineField>

      <InlineField label="Password" labelWidth={labelWidth} tooltip="Password to connect to the AMQP server">
        <SecretInput
          isConfigured={!!secureJsonFields.password}
          value={secureJsonData?.password ?? amqpDefaultPassword}
          placeholder={amqpDefaultPassword}
          width={inputWidth}
          onReset={onResetPassword}
          onChange={onHostChange}
        />
      </InlineField>
    </div>
  );
}
