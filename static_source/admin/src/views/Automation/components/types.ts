import {ApiTrigger} from "@/api/stub";

export interface Trigger extends ApiTrigger {
    timePluginOptions?: string;
    systemPluginOptions?: string;
    alexaPluginOptions?: number;
}
