
/*
    This class contains information about 2 endpoints:
    1) Auth proxy endpoint
    2) Storage endpoint
 */
export interface EndpointsCoordinates {
    authentication: Coordinates
    storage: Coordinates
}

class Coordinates {
    schema: "http" | "https" | "same"
    hostname: string
    port: number

    toString() :string {
        return ""
    }
}