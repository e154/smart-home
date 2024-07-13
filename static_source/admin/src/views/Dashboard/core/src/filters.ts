import {parseTime} from '@/utils';
import {Resolve, scriptService} from "@/views/Dashboard/core";

export const ApplyFilter = (value: any, filter: string): any => {
  if (value == undefined || filter == undefined) {
    return value;
  }

  let args = filter.split('::');
  if (args.length > 1) {
    filter = args[0];
    args.shift();
  } else {
    args = []
  }

  switch (filter) {
    case 'secToTime':
      return secToTime(value, ...args);
    case 'secToCounter':
      return secToCounter(value, ...args);
    case 'formatdate':
      return formatdate(value, ...args);
    case 'formatBytes':
      return formatBytes(value, ...args);
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
    case 'upperCase':
      return upperCase(value, ...args);
    case 'render':
      return render(value, ...args);
    case 'script':
      return scriptService.evalScript(value, ...args);
    default:
      console.warn(`unknown filter "${filter}"!`);
      return value;
  }
}

const render = (obj: any, ...args: string[]): string => {
  let token = '';
  if (args[0]) {
    token = args[0];
  }
  return Resolve(token, JSON.parse(obj));
}

//DEPRECATED
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

function secToCounter(value: string, ...args: string[]): string {
  const seconds = parseInt(value);

  const mSeconds = seconds * 1000

  const epoch = new Date(0);
  const delta = new Date(epoch.getTime() + mSeconds);

  // const years = delta.getYear() - epoch.getYear();
  const months = delta.getUTCMonth() - epoch.getUTCMonth();
  const days = delta.getUTCDate() - epoch.getUTCDate();
  const hours = delta.getUTCHours() - epoch.getUTCHours();
  const minutes = delta.getUTCMinutes() - epoch.getUTCMinutes();

  if (args && args.length) {
    let result = '';
    for (let i = 0; i < args.length; i++) {
      for (let j = 0; j < args[i].length; j++) {
        // console.log('e ', args[i].charAt(j))
        switch (args[i].charAt(j)) {
          // case 'y':
          //   result += String(years).padStart(2, '0');
          //   break;
          case 'M':
            result += String(months).padStart(2, '0');
            break;
          case 'd':
            result += String(days).padStart(2, '0');
            break;
          case 'h':
            result += String(hours).padStart(2, '0');
            break;
          case 'm':
            result += String(minutes).padStart(2, '0');
            break;
          default:
            result += args[i].charAt(j)
        }
      }
    }
    return result
  } else {
    const formattedMonths = String(months).padStart(2, '0');
    const formattedDays = String(days).padStart(2, '0');
    const formattedHours = String(hours).padStart(2, '0');
    const formattedMinutes = String(minutes).padStart(2, '0');

    return `${formattedMonths}:${formattedDays}:${formattedHours}:${formattedMinutes}`;
  }
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
    function (txt) {
      return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
    }
  );
}

function upperCase(value: string, ...args: string[]): string {
  return value.toUpperCase()
}

export function formatBytes(value: string, ...args: string[]): number | string {
  let decimals = 2;
  const bytes = parseInt(value);

  if (args && args.length > 0) {
    decimals = parseInt(args[0])
  }

  if (bytes === 0) {
    return '0 Bytes';
  }

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

export interface Filter {
  name: string;
  description: string;
  example: string;
  args?: string[];
}

export const Filters: Filter[] = [
  {
    name: 'ccToTitleCase',
    description: 'The ccToTitleCase function takes a string in camelCase format and converts it to a string where each word begins with a capital letter.',
    example: 'string|ccToTitleCase',
  },
  {
    name: 'formatdate',
    description: 'The formatdate function takes a date in a specific format and converts it to another specified format. Input format YYYY-MM-DD hh:mm:ss',
    example: 'datetime|formatdate::{d}.{m}.{y}',
  },
  {
    name: 'formatBytes',
    description: 'The formatBytes function takes the file size in bytes and converts it to a human-readable format such as kilobytes, megabytes, or gigabytes, depending on the file size.',
    example: 'number|formatBytes::2',
  },
  {
    name: 'getDayOfWeek',
    description: 'The getDayOfWeek function takes a date and converts it to a day of the week. args: \'long\' | \'short\' | \'narrow\' ',
    example: 'datetime|getDayOfWeek::short',
  },
  {
    name: 'render',
    description: 'The render function takes a string in json format and tries to find the specified argument value in it',
    example: 'string|render::new_state.attribute',
  },
  {
    name: 'scToTitleCase',
    description: 'The snakeCaseStringToTitleCase function takes a string in snake_case format and converts it to a string where each word begins with a capital letter.',
    example: 'string|scToTitleCase',
  },
  {
    name: 'script',
    description: 'The script function calls the script with the name from the argument and passes the value as the callable parameter.',
    example: 'value|script::123',
  },
  {
    name: 'secToTime',
    description: 'DEPRECATED! The secToTime function takes the number of seconds as an argument and converts this value to a time format in hours, minutes, and seconds.',
    example: 'uptime_total|secToTime::H::m',
    args: ["H", "d", "m"],
  },
  {
    name: 'secToCounter',
    description: 'The secToCounter function takes the number of seconds as an argument and converts this value into a time counter format that can be used for counting down.',
    example: 'uptime_total|secToCounter::M::d::h',
    args: ["M", "d", "h", "m"],
  },
  {
    name: 'seconds',
    description: 'The seconds function takes nanoseconds and converts them to seconds.',
    example: 'number|seconds',
  },
  {
    name: 'toTitleCase',
    description: 'The toTitleCase function takes a string and converts the first letter of each word in the string to uppercase and the remaining letters to lowercase.',
    example: 'string|toTitleCase',
  },
  {
    name: 'toFixed',
    description: 'The toFixed function takes a number and number of decimal places as arguments, and returns a number with the specified number of decimal places.',
    example: 'number|toFixed::2',
  },
  {
    name: 'upperCase',
    description: 'The UpperCase function takes a string and converts it into a string with each letter capitalized.',
    example: 'string|upperCase',
  }
]
