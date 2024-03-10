import {EventStateChange} from "@/api/types";
import {eventBus} from "@/views/Dashboard/core";
import { isEqual } from 'lodash-es'
import stream from "@/api/stream";
import {UUID} from "uuid-generator-ts";

class StateService {
    private _lastEvents: Map<string, EventStateChange> = {} as Map<string, EventStateChange>;

    lastEvent(entityId?: string): EventStateChange | undefined {
        if (!entityId) {
            return;
        }
        if (!this._lastEvents[entityId]) {
            this._lastEvents[entityId] = {} as EventStateChange
            this.requestCurrentState(entityId)
            return undefined
        }
        return this._lastEvents[entityId];
    }

    requestCurrentState = (entityId?: string) => {
        if (!entityId) {
            return;
        }
        if (this._lastEvents[entityId] && !isEqual(this._lastEvents[entityId],{})) {
            this.onStateChanged(this._lastEvents[entityId])
            return;
        }
        if (isEqual(this._lastEvents[entityId],{})) {
            return;
        }
        if (!this._lastEvents[entityId]) {
            this._lastEvents[entityId] = {} as EventStateChange;
        }
        // console.log('requestCurrentState', entityId);
        stream.send({
            id: UUID.createUUID(),
            query: 'event_get_last_state',
            body: btoa(JSON.stringify({entity_id: entityId}))
        });
    }

    onStateChanged = (event: EventStateChange)=> {
        // console.log('onStateChanged', event);
        if (event.entity_id && this._lastEvents[event.entity_id]) {
            this._lastEvents[event.entity_id] = event;
        }
        eventBus.emit('stateChanged', event)
    }
}

export const stateService = new StateService();

export function requestCurrentState(entityId?: string) {
    stateService.requestCurrentState(entityId)
}

export function lastEvent(entityId?: string) {
    stateService.lastEvent(entityId)
}
