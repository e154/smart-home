import {useEmitt} from "@/hooks/web/useEmitt";
import {ApiImage} from "@/api/stub";

const {emitter} = useEmitt()

export interface point {
    x: number;
    y: number;
}

export interface JoystickAction {
    entityId?: string
    action?: string
    tags?: string[]
    areaID?: number
}

export interface ItemPayloadJoystick {
    stickImage?: ApiImage;
    startAction?: JoystickAction;
    endAction?: JoystickAction;
    startTimeout?: number;
    endTimeout?: number;
}

export class JoystickController {
    dragStart: point | null;
    _value: point;
    active: boolean;
    touchId: any;

    constructor(_stick: any, maxDistance, deadzone, currentID: string) {

        // location from which drag begins, used to calculate offsets
        this.dragStart = null;

        // track touch identifier in case multiple joysticks present
        this.touchId = null;

        this.active = false;
        this._value = {x: 0, y: 0};

        const handleDown = (event) => {
            this.active = true;

            // all drag movements are instantaneous
            _stick._value.style.transition = '0s';

            // touch event fired before mouse event; prevent redundant mouse event from firing
            event.preventDefault();

            if (event.changedTouches)
                this.dragStart = {x: event.changedTouches[0].clientX, y: event.changedTouches[0].clientY};
            else
                this.dragStart = {x: event.clientX, y: event.clientY};

            // if this is a touch event, keep track of which one
            if (event.changedTouches)
                this.touchId = event.changedTouches[0].identifier;
        }

        const handleMove = (event) => {
            if (!this.active) return;

            // if this is a touch event, make sure it is the right one
            // also handle multiple simultaneous touchmove events
            let touchmoveId = null;
            if (event.changedTouches) {
                for (let i = 0; i < event.changedTouches.length; i++) {
                    if (this.touchId == event.changedTouches[i].identifier) {
                        touchmoveId = i;
                        event.clientX = event.changedTouches[i].clientX;
                        event.clientY = event.changedTouches[i].clientY;
                    }
                }

                if (touchmoveId == null) return;
            }

            const xDiff = event.clientX - this.dragStart?.x;
            const yDiff = event.clientY - this.dragStart?.y;
            const angle = Math.atan2(yDiff, xDiff);
            const distance = Math.min(maxDistance, Math.hypot(xDiff, yDiff));
            const xPosition = distance * Math.cos(angle);
            const yPosition = distance * Math.sin(angle);

            // move stick image to new position
            _stick._value.style.transform = `translate3d(${xPosition}px, ${yPosition}px, 0px)`;

            // deadzone adjustment
            const distance2 = (distance < deadzone) ? 0 : maxDistance / (maxDistance - deadzone) * (distance - deadzone);
            const xPosition2 = distance2 * Math.cos(angle);
            const yPosition2 = distance2 * Math.sin(angle);
            const xPercent = parseFloat((xPosition2 / maxDistance).toFixed(4));
            const yPercent = parseFloat((yPosition2 / maxDistance).toFixed(4));

            this._value = {x: xPercent, y: yPercent * -1};
            emitter.emit('updateValue', {
                id: currentID,
                value: this._value,
            })
        }

        const handleUp = (event) => {
            if (!this.active) return;

            // if this is a touch event, make sure it is the right one
            if (event.changedTouches && this.touchId != event.changedTouches[0].identifier) return;

            // transition the joystick position back to center
            _stick._value.style.transition = '.2s';
            _stick._value.style.transform = `translate3d(0px, 0px, 0px)`;

            // reset everything
            this._value = {x: 0, y: 0};
            emitter.emit('updateValue', {
                id: currentID,
                value: this._value,
            })

            this.touchId = null;
            this.active = false;
        }

        _stick._value.addEventListener('mousedown', handleDown, {passive: false});
        _stick._value.addEventListener('touchstart', handleDown, {passive: false});
        document.addEventListener('mousemove', handleMove, {passive: false});
        document.addEventListener('touchmove', handleMove, {passive: false});
        document.addEventListener('mouseup', handleUp, {passive: false});
        document.addEventListener('touchend', handleUp, {passive: false});
    }

}
