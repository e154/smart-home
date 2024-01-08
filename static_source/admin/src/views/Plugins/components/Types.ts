import {ApiAttribute, ApiPluginOptionsResultEntityAction, ApiPluginOptionsResultEntityState} from "@/api/stub";
import {parseTime} from "@/utils";

export interface Plugin {
    name: string;
    version: string;
    enabled: boolean;
    system: boolean;
    actor?: boolean;
    triggers?: boolean;
    actors?: boolean;
    actorCustomAttrs?: boolean;
    actorCustomActions?: boolean;
    actorCustomStates?: boolean;
    actorCustomSetts?: boolean;
    actorAttrs?: Record<string, ApiAttribute>;
    actorActions?: ApiPluginOptionsResultEntityAction[];
    actorStates?: ApiPluginOptionsResultEntityState[];
    actorSetts?: Record<string, ApiAttribute>;
    setts?: Record<string, ApiAttribute>;
}

export const getUrl = (imageUrl: string | undefined): string => {
    if (!imageUrl) {
        return '';
    }
    return import.meta.env.VITE_API_BASEPATH + imageUrl;
}

export const getValue = (attr: ApiAttribute): any => {
    switch (attr.type) {
        case 'STRING':
            return attr.string;
        case 'INT':
            return attr.int;
        case 'FLOAT':
            return attr.float;
        case 'ARRAY':
            return attr.array;
        case 'BOOL':
            return attr.bool;
        case 'TIME':
            return parseTime(attr.time);
        case 'MAP':
            return attr.map;
        case 'IMAGE':
            return getUrl(attr.imageUrl);
    }
}
