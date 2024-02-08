import {ApiDashboardTab, ApiEntity} from "@/api/stub";
import {comparisonType} from "@/views/Dashboard/core/core";


export interface Area {
    id: number;
    name: string;
    description?: string;
}

export interface Dashboard {
    id: number;
    name: string;
    description?: string;
    enabled: boolean;
    areaId?: number;
    area?: Area;
    tabs?: ApiDashboardTab[];
    entities?: Map<string, ApiEntity>;
    createdAt?: string;
    updatedAt?: string;
}

export interface CompareProp {
    key: string;
    comparison: comparisonType;
    value: string;
    entity?: { id?: string };
    entityId?: string;
}
