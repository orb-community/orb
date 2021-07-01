export interface LoraMetadata {
    type: string;
    lora: {
      dev_eui?: string,
      app_id?: string,
    };
    channel_id: string;
}

export interface LoraDevice {
  name?: string;
  id?: string;
  key?: string;
  appID?: string;
  devEUI?: string;
  seen?: string;
  messages?: string;
  metadata?: LoraMetadata;
}
