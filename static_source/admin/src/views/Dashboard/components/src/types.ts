import {ApiImage} from "@/api/stub";

export interface KeysProp {
    keys?: Map<number, string>;
    entity?: { id?: string };
    entityId?: string;
    action?: string;
    areaId?: number;
    tags?: string[];
}

export interface FrameItem {
    x: number;
    y: number;
    width: number;
    height: number;
    base64?: string;
    canvasData?: any;
}

export interface FrameProp {
    image?: ApiImage;
    items: Map<string, FrameItem>;
}
