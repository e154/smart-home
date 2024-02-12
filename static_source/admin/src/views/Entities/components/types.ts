import {ApiArea, ApiAttribute, ApiImage, ApiMetric, ApiScript} from "@/api/stub";

export interface Plugin {
  name: string;
  version: string;
  enabled: boolean;
  system: boolean;
  actor?: boolean;
  settings: ApiAttribute[];
  triggers?: boolean;
  actors?: boolean;
  actorCustomAttrs?: boolean;
  actorAttrs?: ApiAttribute[];
  actorCustomActions?: boolean;
  actorActions?: EntityAction[];
  actorCustomStates?: boolean;
  actorStates?: EntityState[];
  actorCustomSetts?: boolean;
  actorSetts?: ApiAttribute[];
  setts?: ApiAttribute[];
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
  attributes?: Record<string, ApiAttribute>;
  settings?: Record<string, ApiAttribute>;
  scriptIds?: number[];
  tags?: string[];
  scripts?: ApiScript[];
  metrics?: ApiMetric[];
}


