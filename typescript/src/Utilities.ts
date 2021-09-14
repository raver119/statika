import { MetaType } from "./Statika";

export function pickDefined<T extends any>(val: T | undefined, def: T): T {
  return val === undefined ? def : val;
}

export const objectifyMeta = (extras: MetaType): {} => {
  const obj = {};

  for (let v of extras.entries()) {
    obj[v[0]] = v[1];
  }

  return obj;
};

export const stringifyMeta = (extras: MetaType): string => {
  return JSON.stringify(objectifyMeta(extras));
};

export const materializeMeta = (meta: any): MetaType => {
  if (meta.has !== undefined && meta.set !== undefined) return meta;

  const clone = new Map<string, string>();

  Object.keys(meta).forEach((key) => {
    clone.set(key, meta[key]);
  });

  return clone;
};

export const randomFloat = (min: number, max: number): number => {
  if (max < min) throw new Error(`Min <${min}> shouldn't be > Max <${max}>`);

  if (max === min) return max;

  max = max - 1;
  return Math.random() * (max - min) + min;
};

export const randomInt = (min: number, max: number): number => {
  return Math.round(randomFloat(min, max));
};

export function makeRandomString(length: number): string {
  const charSet = "abcdedfghijklmnopqrstuvwzyz01234567890\nABSC JKW";

  let b: string = "";
  for (let i = 0; i < length; i++) {
    const random = randomInt(0, charSet.length);
    b = `${b}${charSet[random]}`;
  }

  return b;
}
