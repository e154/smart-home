<script>
import { deepFlat } from "@daybrush/utils";
import { useKeycon } from "vue-keycon";
import Selecto from "vue3-selecto";
import Moveable, { MoveableTargetGroupsType } from "vue3-moveable";
import { GroupManager, TargetList } from "@moveable/helper";
import { ref, onMounted } from "vue";

export default {
  components: { Moveable, Selecto },
  setup() {
    const { isKeydown: isCommand } = useKeycon({ keys: "meta" });
    const { isKeydown: isShift } = useKeycon({ keys: "shift" });
    const groupManagerRef = ref(undefiend);
    const targets = ref([]);
    const moveableRef = ref(null);
    const selectoRef = ref(null);
    const cubes = [];
    for (let i = 0; i < 30; ++i) {
      cubes.push(i);
    }
    const setSelectedTargets = (nextTargetes) => {
      selectoRef.value.setSelectedTargets(deepFlat(nextTargetes));
      targets.value = nextTargetes;
    };
    const onClickGroup = e => {
      if (!e.moveableTarget) {
        setSelectedTargets([]);
        return;
      }
      if (e.isDouble) {
        const childs = groupManagerRef.value.selectSubChilds(
            targets.value,
            e.moveableTarget
        );
        setSelectedTargets(childs.targets());
        return;
      }
      if (e.isTrusted) {
        selectoRef.value.clickTarget(e.inputEvent, e.moveableTarget);
      }
    };
    const onDrag = e => {
      e.target.style.transform = e.transform;
    };
    const onRenderGroup = e => {
      e.events.forEach(ev => {
        ev.target.style.cssText += ev.cssText;
      });
    };
    const onDragStart = e => {
      const moveable = moveableRef.value;
      const target = e.inputEvent.target;
      const flatted = targets.value.flat(3);
      if (moveable.isMoveableElement(target)
          || flatted.some(t => t === target || t.contains(target))
      ) {
        e.stop();
      }
    };
    const onSelectEnd = e => {
      const {
        isDragStartEnd,
        isClick,
        added,
        removed,
        inputEvent,
      } = e;
      const moveable = moveableRef.value;
      if (isDragStartEnd) {
        inputEvent.preventDefault();
        moveable.waitToChangeTarget().then(() => {
          moveable.dragStart(inputEvent);
        });
      }
      const groupManager = groupManagerRef.value;
      let nextChilds;
      if (isDragStartEnd || isClick) {
        if (isCommand.value) {
          nextChilds = groupManager.selectSingleChilds(
              targets.value,
              added,
              removed
          );
        } else {
          nextChilds = groupManager.selectCompletedChilds(
              targets.value,
              added,
              removed,
              isShift.value
          );
        }
      } else {
        nextChilds = groupManager.selectSameDepthChilds(
            targets.value,
            added,
            removed
        );
      }
      e.currentTarget.setSelectedTargets(nextChilds.flatten());
      setSelectedTargets(nextChilds.targets());
    };
    onMounted(() => {
      const elements = selectoRef.value.getSelectableElements();
      groupManagerRef.value = new GroupManager([
        [[elements[0], elements[1]], elements[2]],
        [elements[5], elements[6], elements[7]],
      ], elements);
    });
    return {
      moveableRef,
      targets,
      onClickGroup,
      onDrag,
      onRenderGroup,
      selectoRef,
      window,
      onDragStart,
      onSelectEnd,
      cubes
    };
  }
};
</script>
<template>
  <div class="root">
    <div class="container">
      <Moveable
          ref="moveableRef"
          :draggable="true"
          :rotatable="true"
          :scalable="true"
          :target="targets"
          @clickGroup="onClickGroup"
          @drag="onDrag"
          @renderGroup="onRenderGroup"
      ></Moveable>
      <Selecto
          ref="selectoRef"
          :dragContainer="window"
          :selectableTargets="['.selecto-area .cube']"
          :hitRate="0"
          :selectByClick="true"
          :selectFromInside="false"
          :toggleContinueSelect="['shift']"
          :ratio="0"
          @dragStart="onDragStart"
          @selectEnd="onSelectEnd"
      ></Selecto>
      <p>[[0, 1], 2] is group</p>
      <p>[5, 6, 7] is group</p>
      <div class="elements selecto-area">
        <div
            class="cube"
            :key="i"
            v-for="i in cubes"
        >{{i}}</div>
      </div>
      <div class="empty elements"></div>
    </div>
  </div>
</template>
