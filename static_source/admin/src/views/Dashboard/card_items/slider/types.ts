export enum OrientationType {
    horizontal = 'horizontal',
    vertical = 'vertical',
    circular = 'circular',
}

export interface ItemPayloadSlider {
    height?: number;
    color?: string;
    trackColor?: string;
    min?: number;
    max?: number;
    step?: number;
    orientation?: OrientationType;
    attribute?: string;
    action?: string;
    tooltip?: boolean;
}
