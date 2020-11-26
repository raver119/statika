
export function pickDefined<T extends any>(val: T|undefined, def: T) :T {
    return val === undefined ? def : val
}