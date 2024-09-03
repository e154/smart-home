import {startRecording, stopRecording} from "@/components/Stt";

class STT {
    _isRecording: boolean;

    constructor() {
        this._isRecording = false;
    }

    startRecording = () => {
        this._isRecording = true;
        startRecording()
    }

    stopRecording = () => {
        this._isRecording = false;
        stopRecording()
    }

    toggleRecording = () => {
        if (!this._isRecording) {
            this.startRecording()
        } else {
            this.stopRecording()
        }
    }

    eventBusHandler(eventName: string, args: any[]) {
        switch (eventName) {
            case 'SttStartRecording':
                this.startRecording()
                break;
            case 'SttStopRecording':
                this.stopRecording()
                break;
            case 'SttToggleRecording':
                this.toggleRecording()
                break;

        }
    }
}

export const stt = new STT();
