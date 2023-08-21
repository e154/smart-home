import {ApiDashboardTab, ApiEntity} from "@/api/stub";


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
