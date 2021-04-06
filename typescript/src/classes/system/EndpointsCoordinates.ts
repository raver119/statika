
export type EndpointSchemaType = "http" | "https" | "same"

/*
    This class contains information about 2 endpoints:
    1) Auth proxy endpoint
    2) Storage endpoint
 */
export const coordinates = (schema: EndpointSchemaType, host: string, port: number|string) => {
    return {
        schema: schema,
        hostname: host,
        port: port,

        toString() :string {
            if (schema === "same") {
                // same FQDN + relative path will be used
                return ""
            } else {
                let fport = port === 80 || port === "80" ? "" : `:${port}`
                return `${schema}://${host}${fport}`
            }
        },
    }
}

export type EndpointsCoordinates = ReturnType<typeof coordinates>

export const testCoordinates = (host: string, port: number|string) => {
    return coordinates("http", host, port)
}

export const sameCoordinates = () => {
    return coordinates("same", undefined, undefined)
}

