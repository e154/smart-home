import {parseTime} from '@/utils';

export function ApplyFilter(value: any, filter: string): any {
  if (value == undefined || filter == undefined) {
    return value;
  }

  const args = filter.split('::');
  if (args.length > 1) {
    filter = args[0];
    args.splice(0, 1);
  }

  switch (filter) {
    case 'secToTime':
      return secToTime(value, ...args);
    case 'formatdate':
      return formatdate(value, ...args);
    case 'formatBytes':
      return formatBytes(value, 2);
    case 'seconds':
      return seconds(value, ...args);
    case 'getDayOfWeek':
      return getDayOfWeek(value, ...args);
    case 'toFixed':
      return toFixed(value, ...args);
    case 'scToTitleCase':
      return snakeCaseStringToTitleCase(value, ...args);
    case 'ccToTitleCase':
      return camelCaseStringToTitleCase(value, ...args);
    case 'toTitleCase':
      return toTitleCase(value, ...args);
    default:
      console.warn(`unknown filter "${filter}"!`);
      return value;
  }
}

function secToTime(value: string, ...args: string[]): string {
  const num = parseInt(value);
  const days = Math.floor(num / (24 * 3600));
  const hours = Math.floor((num - days * 24 * 3600) / 3600);
  const minutes = Math.floor(num % 3600 / 60);

  let d = 'Days';
  let h = 'Hours';
  let m = 'Minutes';
  if (args[0]) {
    d = args[0];
  }
  if (args[1]) {
    h = args[1];
  }
  if (args[2]) {
    m = args[2];
  }
  let result = minutes + `${m}`;
  if (h) {
    result = hours + `${h}:` + result;
  }
  if (d) {
    result = days + `${d}:` + result;
  }
  return result;
}

function seconds(value: string, ...args: string[]) {
  // nanoseconds
  const num = parseInt(value);
  return Math.floor(num / 100000) / 10000;
}

function formatdate(value: string, ...args: string[]): string {
  if (args && args.length > 0 && args[0] != 'formatdate') {
    return parseTime(value, ...args) || '[date not valid]';
  }
  return parseTime(value) || '[date not valid]';
}

function getDayOfWeek(value: string, ...args: string[]): string {
  const date = new Date(value);
  let weekday: 'long' | 'short' | 'narrow' = 'long';
  if (args && args.length > 0) {
    weekday = args[0] as 'long' | 'short' | 'narrow';
  }
  return date.toLocaleString(
    'en-US', {weekday: weekday}
  );
}

function toFixed(value: string, ...args: string[]): string {
  const date = parseFloat(value);
  if (args && args.length > 0) {
    return date.toFixed(parseInt(args[0]));
  }
  return date.toString();
}

function snakeCaseStringToTitleCase(value: string, ...args: string[]): string {
  return value.replace(/^_*(.)|_+(.)/g, (s, c, d) => c ? c.toUpperCase() : ' ' + d.toUpperCase());
}

function camelCaseStringToTitleCase(value: string, ...args: string[]): string {
  const result = value.replace(/([A-Z])/g, ' $1');
  return result.charAt(0).toUpperCase() + result.slice(1);
}

function toTitleCase(value: string, ...args: string[]): string {
  return value.replace(
    /\w\S*/g,
    function(txt) {
      return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
    }
  );
}

function formatBytes(value: string, decimals = 2): number | string {
  const bytes = parseInt(value);

  if (bytes === 0) {
    return '0 Bytes';
  }

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

export interface filterInfo {
  name: string;
  description: string;
  args?: string[];
}

export function filterList(): filterInfo[] {
  const info: filterInfo[] = [
    {
      name: 'secToTime',
      description: 'convert seconds to SS:MM:DD string',
      args: []
    },
    {
      name: 'formatdate',
      description: 'format raw string date',
      args: []
    },
    {
      name: 'formatBytes',
      description: 'convert bytes to Gb, Mb',
      args: []
    },
    {
      name: 'seconds',
      description: 'format nanoseconds to seconds string',
      args: []
    },
    {
      name: 'getDayOfWeek',
      description: 'Get the Day of the Week',
      args: ['long', 'short', 'narrow']
    },
    {
      name: 'toFixed',
      description: 'Function rounds to the desired decimal point',
      args: ['1', '2', '3']
    },
    {
      name: 'scToTitleCase',
      description: 'snake case string to title case',
      args: []
    },
    {
      name: 'ccToTitleCase',
      description: 'camel case string to title case',
      args: []
    },
    {
      name: 'toTitleCase',
      description: 'string to title case',
      args: []
    }
  ];
  return info;
}
