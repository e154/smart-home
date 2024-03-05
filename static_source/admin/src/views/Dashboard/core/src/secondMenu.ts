import {Core} from "./core";
import {useBus} from "./bus"
import {EventContextMenu} from "@/views/Dashboard/core";
import {MenuItem} from "@imengyu/vue3-context-menu/lib/ContextMenuDefine";
import ContextMenu from "@imengyu/vue3-context-menu";
import {useAppStore} from "@/store/modules/app";
import {useI18n} from "@/hooks/web/useI18n";
import {CardItemList, ItemsType} from "@/views/Dashboard/card_items";

const {t} = useI18n()
const appStore = useAppStore()
const {emit} = useBus()

export class SecondMenu {
  private core: Core;

  constructor(core: Core) {
    this.core = core;

    useBus({
      name: 'eventContextMenu',
      callback: (event) => this.contextMenu(event),
    })
  }

  start = () => {
  }
  shutdown = () => {
  }

  private genItemMenu = (cardId: number) => {

    const getChildren = (list: ItemsType[]): MenuItem[] => {
      const menuItem: MenuItem[] = []
      for (const item of list) {
        const _item: MenuItem = {
          label: item.label,
        }
        if (item.children && item.children.length) {
          _item.children = getChildren(item.children)
        } else {
          _item.onClick = () => {
            this.core.createCardItem(cardId, item.value)
          }
        }
        menuItem.push(_item)
      }
      return menuItem
    }

    return getChildren(CardItemList)
  }

  private buildMenu = (event: EventContextMenu): MenuItem[] => {

    const items: MenuItem[] = []

    if (event?.cardId) {
      items.push(...[
        {
          label: t('dashboard.addCardItem'),
          children: this.genItemMenu(event?.cardId)
        } as MenuItem,
      ])
    }

    if (this.core.getActiveTab) {
      items.push(
        {
          label: t('dashboard.addNewCard'),
          onClick: () => {
            this.core.createCard();
          }
        })
    }

    items.push({
      label: t('dashboard.addNewTab'),
      onClick: () => {
        this.core.createTab()
      }
    })

      items.push(
          {
            divided: 'self',
          } as MenuItem)

    // check current tab
    const tabs: MenuItem[] = [
      {
        label: t('main.import'),
        onClick: () => {
          emit('showTabImportDialog', true)
        }
      },
    ]
    if (this.core.getActiveTab) {
      tabs.push({
        label: t('main.export'),
        onClick: () => {
          emit('showTabExportDialog', true)
        }
      })
    }
    if (tabs.length) {
      items.push({
        label: t('dashboard.tabs'),
        children: tabs,
      })
    }

    // check current card
    const cards: MenuItem[] = []
    if (this.core.getActiveTab) {
      cards.push(...[
        {
          label: t('main.import'),
          onClick: () => {
            emit('showCardImportDialog', true)
          }
        },
      ])
    }
    if (event?.cardId) {
      cards.push(...[{
        label: t('main.export'),
        onClick: () => {
          emit('showCardExportDialog', event.cardId)
        }
      },
        {
          label: t('main.duplicate'),
          onClick: () => {
            this.core.copyCard(event?.cardId)
          }
        }
      ])
    }
    if (cards.length) {
      items.push({
        label: t('dashboard.cards'),
        children: cards,
      })
    }

    // check current card item
    const cardItems: MenuItem[] = []
    if (event?.cardId) {
      cardItems.push({
        label: t('main.import'),
        onClick: () => {
          emit('showCardItemImportDialog', event?.cardId)
        }
      })
    }
    if (event?.cardItemId) {
      cardItems.push(...[
        {
          label: t('main.export'),
          onClick: () => {
            emit('showCardItemExportDialog', event.cardItemId)
          }
        },
        {
          label: t('main.duplicate'),
          onClick: () => {
            this.core.copyCardItem(event?.cardItemId)
          }
        }
      ])
    }
    if (cardItems.length) {
      items.push({
        label: t('dashboard.cardItem'),
        children: cardItems,
      })
    }

    // divider
    if (this.core.getActiveTab || event.tabId || event.cardId || event.cardItemId) {
      items.push({
        divided: 'self',
      })
    }

    // remove card item
    if (event?.cardItemId) {
      items.push({
        label: t('dashboard.removeCardItem'),
        onClick: () => {
          this.core.removeCardItemById(event.cardItemId)
        }
      })
    }

    // remove card
    if (event?.cardId) {
      items.push({
        label: t('dashboard.removeCard'),
        onClick: () => {
          this.core.removeCard(event.cardId);
        }
      })
    }

    // remove tab
    if (this.core.getActiveTab) {
      items.push({
        label: t('dashboard.removeTab'),
        onClick: () => {
          this.core.removeTab();
        }
      })
    }

    return items;
  }

  private contextMenu = (event: EventContextMenu) => {

    const theme = appStore.isDark ? 'dark' : 'light'
    ContextMenu.showContextMenu({
      x: event.event.x,
      y: event.event.y,
      theme: 'flat ' + theme,
      zIndex: 9999,
      // minWidth: 230,
      items: this.buildMenu(event)
    });
  }
}

