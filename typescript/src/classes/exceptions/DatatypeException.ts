
export class DatatypeException extends Error {
    constructor(expected: string, obj: any) {
        super(`Expected <${expected}>, but got <${JSON.stringify(obj)}> instead`);
    }
}