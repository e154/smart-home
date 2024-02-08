export enum playerType {
  onvifMse = 'onvifMse',
  youtube = 'youtube',
}

export interface ItemPayloadVideo {
  playerType?: playerType;
  attribute?: string;
}
