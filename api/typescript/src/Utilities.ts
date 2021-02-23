import {MetaType} from "./Statika";

export function pickDefined<T extends any>(val: T|undefined, def: T) :T {
    return val === undefined ? def : val
}

export const objectifyMeta = (extras: MetaType) : {} => {
    const obj = {}

    for (let v of extras.entries()) {
        obj[v[0]] = v[1]
    }

    return obj
}

export const stringifyMeta = (extras: MetaType) :string => {
    return JSON.stringify(objectifyMeta(extras))
}

export const materializeMeta = (meta: any) :MetaType => {
    if (meta.has !== undefined && meta.set !== undefined)
        return meta

    const clone = new Map<string, string>()

    Object.keys(meta).forEach(key => {
        clone.set(key, meta[key])
    })

    return clone
}