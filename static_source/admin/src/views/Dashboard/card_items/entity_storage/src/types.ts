export interface Column {
    name: string
    attribute: string
    filter?: string
    width?: number
    sortable?: boolean;
}

export interface ItemPayloadEntityStorage {
    entityIds: string[]
    filter: boolean
    columns: Column[]
    eventName?: string
    showPopup?: boolean
}
