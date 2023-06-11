export default {
  route: {
    dashboard: 'Dashboard',
    dashboards: 'Dashboards',
    dashboardEdit: 'Edit',
    dashboardList: 'Dashboards',
    development: 'Development',
    documentation: 'Documentation',
    errorPages: 'Error Pages',
    page401: '401',
    page404: '404',
    errorLog: 'Error Log',
    i18n: 'I18n',
    externalLink: 'External Link',
    profile: 'Profile',
    scripts: 'Scripts',
    scriptList: 'Scripts',
    scriptNew: 'New script',
    scriptEdit: 'Edit script',
    plugins: 'Plugins',
    pluginList: 'Plugins',
    areas: 'Areas',
    areaList: 'Areas',
    areaNew: 'New Area',
    areaEdit: 'Edit area',
    entities: 'Entities',
    entityList: 'Entities',
    entityNew: 'New entity',
    entityEdit: 'Edit entity',
    automation: 'Automation',
    taskList: 'Automation',
    taskNew: 'New automation',
    taskEdit: 'Edit automation',
    zigbee2mqtt: 'Zigbee2mqtt',
    bridgeList: 'Zigbee2mqtt',
    bridgeNew: 'New zigbee2mqtt',
    bridgeEdit: 'Edit zigbee2mqtt',
    logs: 'Logs',
    logList: 'List',
    images: 'Images',
    imageList: 'List',
    swagger: 'Swagger',
    swaggerList: 'Swagger',
    users: 'Users',
    userList: 'Users',
    maps: 'Maps',
    mapList: 'List',
    mqtt: 'Mqtt',
    mqttList: 'List',
    etc: 'etc',
    variables: 'Variables',
    backupList: 'Backups',
    messageDelivery: 'Message Delivery'
  },
  navbar: {
    logOut: 'Log Out',
    dashboard: 'Dashboard',
    github: 'Github',
    theme: 'Theme',
    size: 'Global Size',
    profile: 'Profile'
  },
  login: {
    title: 'Login Form',
    logIn: 'Login',
    username: 'Username',
    password: 'Password',
    any: 'any',
    thirdparty: 'Or connect with',
    thirdpartyTips: 'Can not be simulated on local, so please combine you own business simulation! ! !'
  },
  table: {
    dynamicTips1: 'Fixed header, sorted by header order',
    dynamicTips2: 'Not fixed header, sorted by click order',
    dragTips1: 'The default order',
    dragTips2: 'The after dragging order',
    title: 'Title',
    importance: 'Importance',
    type: 'Type',
    remark: 'Remark',
    search: 'Search',
    add: 'Add',
    export: 'Export',
    reviewer: 'Reviewer',
    id: 'ID',
    date: 'Date',
    author: 'Author',
    readings: 'Readings',
    status: 'Status',
    actions: 'Actions',
    edit: 'Edit',
    publish: 'Publish',
    draft: 'Draft',
    delete: 'Delete',
    cancel: 'Cancel',
    confirm: 'Confirm'
  },
  scripts: {
    addNew: 'Add script',
    newScript: 'New script',
    view: 'View',
    edit: 'Edit',
    name: 'Name',
    language: 'Language',
    description: 'Description',
    table: {
      id: 'ID',
      name: 'Name',
      lang: 'Lang',
      description: 'Description',
      createdAt: 'Created at',
      updatedAt: 'Updated at'
    }
  },
  plugins: {
    table: {
      name: 'Name',
      version: 'Version',
      enabled: 'Enabled',
      system: 'System'
    }
  },
  variables: {
    addNew: 'Add variable',
    table: {
      name: 'Name',
      value: 'Value'
    }
  },
  areas: {
    addNew: 'Add area',
    edit: 'Edit',
    name: 'Name',
    description: 'Description'
  },
  entities: {
    addNew: 'Add entity',
    addAction: 'Add action',
    addState: 'Add state',
    name: 'Name',
    description: 'Description',
    callAction: 'Call Action',
    setState: 'Set State',
    addAttribute: 'Add attribute',
    loadFromPlugin: 'Load from plugin',
    table: {
      id: 'Id',
      name: 'Name',
      image: 'Image',
      icon: 'Icon',
      type: 'Type',
      value: 'Value',
      parent: 'Parent',
      operations: 'Operations',
      restart: 'Restart',
      pluginName: 'Plugin name',
      autoLoad: 'Auto Load',
      area: 'Area',
      script: 'Script',
      scripts: 'Scripts',
      description: 'Description',
      createdAt: 'Created at',
      updatedAt: 'Updated at'
    },
    metric: {
      name: 'Name',
      description: 'Description',
      type: 'Type',
      ranges: 'Ranges',
      addProp: 'Add prop',
      color: 'Color',
      label: 'Label',
      addMetric: 'Add metric',
      translate: 'Translate',
      list: 'Matric list'
    }
  },
  automation: {
    addTrigger: 'Add trigger',
    addCondition: 'Add condition',
    addAction: 'Add action',
    addNew: 'Add task',
    trigger: {
      pluginOptions: 'Plugin options'
    },
    table: {
      id: 'Id',
      name: 'Name',
      description: 'Description',
      enabled: 'Enabled',
      condition: 'Condition',
      area: 'Area',
      pluginName: 'Plugin',
      script: 'Script',
      entity: 'Entity',
      createdAt: 'Created at',
      updatedAt: 'Updated at'
    }
  },
  message: {
    attributes: 'Attributes',
    table: {
      type: 'Message type',
    }
  },
  message_delivery: {
    table: {
      id: 'Id',
      attributes: 'Attributes',
      status: 'Status',
      createdAt: 'Created at',
      updatedAt: 'Updated at'
    }
  },
  zigbee2mqtt: {
    addNew: 'Add bridge',
    table: {
      id: 'Id',
      name: 'Name',
      login: 'Login',
      password: 'Password',
      model: 'Model',
      status: 'Status',
      description: 'Description',
      permitJoin: 'Permit join',
      createdAt: 'Created at',
      updatedAt: 'Updated at'
    }
  },
  log: {
    table: {
      id: 'Id',
      level: 'Level',
      owner: 'Owner',
      body: 'Body',
      createdAt: 'Created at'
    }
  },
  backup: {
    addNew: 'Create new backup',
    table: {
      name: 'Name',
      restore: 'Restore'
    }
  },
  dashboard: {
    card: {
      'frontend-version': {
        name: 'Frontend'
      }
    },
    addNew: 'Add dashboard',
    table: {
      id: 'Id',
      name: 'Name',
      edit: 'Edit',
      operations: 'Operations',
      description: 'Description',
      enabled: 'Active',
      createdAt: 'Created at',
      updatedAt: 'Updated at'
    },
    editor: {
      getEvent: 'Get event',
      html: 'html',
      addProp: 'Add prop',
      please_add_tab: 'Please add tab',
      currentTab: 'Current tab',
      tabList: 'Tab list',
      tabSettings: 'Tab settings',
      addTab: 'Add tab',
      name: 'Name',
      text: 'Text',
      description: 'Description',
      area: 'Area',
      defaultText: 'Default text',
      asButton: 'As button',
      addAction: 'Add action',
      action: 'Action',
      icon: 'Icon',
      enabled: 'Enabled',
      dragEnabled: 'Drag Enabled',
      gap: 'Gap',
      columnWidth: 'Column Width',
      background: 'Background color',
      addNewCard: 'Add new card',
      addCardItem: 'Add card item',
      currentCard: 'Current card',
      cardList: 'Card list',
      title: 'Title',
      height: 'Height',
      width: 'Width',
      hidden: 'Hidden',
      disabled: 'Disabled',
      frozen: 'Frozen',
      type: 'Type',
      items: 'Items',
      cardItems: 'Card items',
      itemDetail: 'Item detail',
      itemList: 'Item list',
      dashboardSettings: 'Dashboard settings',
      cardSettings: 'Card settings',
      pleaseSelectType: 'Please select type',
      buttonOptions: 'Button options',
      selectAction: 'Select action',
      mainOptions: 'Main options',
      size: 'Size',
      entity: 'Entity',
      image: 'Image',
      value: 'Value',
      comparison: 'Comparison',
      selectStatus: 'Select status',
      showOn: 'Show on',
      hideOn: 'Hide on',
      tokens: 'Tokens',
      defaultImage: 'Default image',
      imageOptions: 'Image options',
      stateOptions: 'State options',
      addNewProp: 'Add new prop',
      textOptions: 'Text options',
      eventstateJSONobject: 'Event state JSON object',
      round: 'Round',
      attrField: 'Attribute field',
      chart: {
        type: 'Chart type',
        entity_metric: 'Entity metric',
        metric_props: 'Metric props',
        borderWidth: 'Border width',
        xAxis: 'Show X axis',
        yAxis: 'Show Y axis',
        legend: 'Show legend',
        range: 'Range',
        filter: 'Filter'
      }
    },
    mainDashboard: 'main dashboard',
    devDashboard: 'development dashboard'
  },
  settings: {
    sidebarTextTheme: 'Sidebar text theme',
    fixedHeader: 'Fixed header',
    showSidebarLogo: 'Show sidebar logo',
    showTagsView: 'Show tags view',
    theme: 'Theme',
    title: 'Title'
  },
  entityStorage: {
    table: {
      operations: 'Operations',
      nothing: 'Nothing',
      state: 'State',
      attributes: 'Attributes',
      createdAt: 'Created at',
      updatedAt: 'Updated at',
      entityId: 'Entity Id'
    }
  },
  users: {
    addNew: 'Add user',
    table: {
      id: 'ID',
      nickname: 'Nick',
      email: 'Email',
      status: 'Status',
      firstName: 'First Name',
      lastName: 'Last Name',
      lang: 'Lang',
      image: 'Image',
      role: 'Role',
      password: 'Password',
      passwordRepeat: 'Password Repeat',
      createdAt: 'Created at',
      updatedAt: 'Updated at'
    }
  },
  roles: {
    addNew: 'Add role',
    table: {
      name: 'Name',
      description: 'Description',
      createdAt: 'Created at',
      updatedAt: 'Updated at'
    }
  },
  main: {
    create: 'Create',
    cancel: 'Cancel',
    edit: 'Edit',
    copy: 'Copy',
    reload: 'Reload',
    remove: 'Remove',
    update: 'Update',
    save: 'Save',
    exec: 'Exec',
    load_from_server: 'Load from server',
    restart: 'Restart',
    call: 'Call',
    ok: 'OK',
    export: 'Export',
    import: 'Import',
    no: 'No',
    'are_you_sure_to_do_want_this?': 'Are you sure to do want this?'
  }
};
