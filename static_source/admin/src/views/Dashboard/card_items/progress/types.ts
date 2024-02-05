import {CompareProp} from "@/views/Dashboard/core";

export interface PropgressProp extends CompareProp {
    color?: string
}

export interface ItemPayloadProgress {
    items: PropgressProp[]
    type: string
    textInside: boolean
    showText: boolean
    strokeWidth: number
    width: number
    value: string
    color: string
}
