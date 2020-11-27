
export {EndpointsCoordinates, testCoordinates, sameCoordinates, coordinates} from "./classes/system/EndpointsCoordinates"

export {AuthenticationResponse, isAuthenticationResponse} from "./classes/comms/AuthenticationResponse"
export {UploadResponse, isUploadResponse} from "./classes/comms/UploadResponse"
export {ApiResponse, responseError, responseOk, responseUnauthorized} from "./classes/comms/ApiResponse"
export {UploadRequest} from "./classes/comms/UploadRequest"

export {HttpException} from "./classes/exceptions/HttpException"
export {AuthenticationException} from "./classes/exceptions/AuthenticationException"

export {pickDefined} from "./Utilities"

export {Statika} from "./Statika"