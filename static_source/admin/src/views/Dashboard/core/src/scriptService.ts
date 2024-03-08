import api from "@/api/api";
import {Cache} from "@/views/Dashboard/core";
import stream from "@/api/stream";
import {UUID} from "uuid-generator-ts";
import {EventUpdatedScriptModel} from "@/api/types";

const uuid = new UUID()

export class ScriptService {
  private cache: Cache;
  readonly currentID: string = uuid.getDashFreeUUID();
  private isStarted = false;

  constructor() {
    this.cache = new Cache()
  }

  public start() {
    if (this.isStarted) return;
    this.isStarted = true;
    stream.subscribe('event_updated_script_model', this.currentID, this.eventHandler)
  }

  public shutdown() {
    if (!this.isStarted) return;
    this.isStarted = false;
    stream.unsubscribe('event_updated_script_model', this.currentID)
  }

  async fetchScript(scriptId: number) {
    try {
      api.v1.scriptServiceGetCompiledScriptById(scriptId)
        .then((res) => {
          if (!res.data) {
            return
          }
          this.cache.push(scriptId + '', res.data)
        })
    } catch (e) {
      return
    }
  }

  eventHandler = (event: EventUpdatedScriptModel) => {
    if (!this.cache.exist(event?.script_id + '')) return
    this.fetchScript(event.script_id)
  }

   evalScript = async (value: string, ...args: string[]): string => {
    if (!args || args.length == 0) {
      return `[${value}::${args}]`
    }

    const scriptId = parseInt(args[0]);

    if (this.cache.get(scriptId + '')) {
      return window.eval.call(window, `(${this.cache.get(scriptId + '')})`)(value);
    }

    try {
      const res = await api.v1.scriptServiceGetCompiledScriptById(scriptId)
      if (!res?.data) {
        return '[SCRIPT]'
      }

      this.cache.push(scriptId + '', res.data)

      return window.eval.call(window, `(${res.data})`)(value);
    } catch (e) {
      return '[SCRIPT]'
    }
  }
}

export const scriptService = new ScriptService()
