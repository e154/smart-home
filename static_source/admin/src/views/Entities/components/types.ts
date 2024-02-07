import {
    ApiAttribute,
    ApiMetric,
    ApiUpdateEntityRequestAction,
    ApiUpdateEntityRequestState,
    ApiArea,
    ApiImage, ApiScript, ApiTypes
} from "@/api/stub";

export interface Plugin {
    name: string;
    version: string;
    enabled: boolean;
    system: boolean;
    actor?: boolean;
    settings: Attribute[];
    triggers?: boolean;
    actors?: boolean;
    actorCustomAttrs?: boolean;
    actorAttrs?: Attribute[];
    actorCustomActions?: boolean;
    actorActions?: EntityAction[];
    actorCustomStates?: boolean;
    actorStates?: EntityState[];
    actorCustomSetts?: boolean;
    actorSetts?: Attribute[];
    setts?: Attribute[];
}

export interface Parent {
    id?: string;
}

export interface EntityAction {
    name?: string;
    description?: string;
    icon?: string;
    image?: ApiImage;
    imageId?: number;
    script?: ApiScript;
    scriptId?: number;
    type?: string;
}

export interface EntityState {
    name?: string;
    description?: string;
    icon?: string;
    image?: ApiImage;
    imageId?: number;
    style?: string;
}

export interface Attribute {
    name?: string;
    type?: ApiTypes;
    int?: number;
    string?: string;
    bool?: boolean;
    float?: number;
    array?: Attribute[];
    map?: Record<string, Attribute>;
    time?: string;
    imageUrl?: string;
    point?: string;
    encrypted?: string;
}

export interface Entity {
    id?: string;
    pluginName?: string;
    plugin?: Plugin,
    description?: string;
    area?: ApiArea;
    areaId?: number;
    icon?: string;
    image?: ApiImage;
    imageId?: number;
    autoLoad?: boolean;
    restoreState?: boolean;
    isLoaded?: boolean;
    parent?: Parent;
    parentId?: string;
    actions?: EntityAction[];
    states?: EntityState[];
    attributes?: Record<string, Attribute>;
    settings?: Record<string, Attribute>;
    scriptIds?: number[];
    tags?: string[];
    scripts?: ApiScript[];
    metrics?: ApiMetric[];
}


export class EntityAttribute implements Attribute {
    constructor(name: string) {
        this.name = name
        this.type = ApiTypes.STRING
        this.string = ''
    }

    name: string;
    type: ApiTypes;
    int?: number;
    string: string;
    bool?: boolean;
    float?: number;
    array?: ApiAttribute[];
}
