import {ApiAttribute, ApiPluginOptionsResultEntityAction, ApiPluginOptionsResultEntityState} from "@/api/stub";

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
