
/*
    This class contains information about 2 endpoints:
    1) Auth proxy endpoint
    2) Storage endpoint
 */
export class EndpointsCoordinates {
    constructor(public schema: "http" | "https" | "same",
                public hostname: string|undefined,
                public port: number|string|undefined) {
    }

    toString() :string {
        if (this.schema === "same") {
            // same FQDN + relative path will be used
            return ""
        } else {
            let port = this.port === 80 ? "" : `:${this.port}`
            return `${this.schema}://${this.hostname}${port}`
        }
    }
}

export const testCoordinates = (host: string, port: number|string) => {
    return new EndpointsCoordinates("http", host, port)
}

export const sameCoordinates = () => {
    return new EndpointsCoordinates("same", undefined, undefined)
}