import {EndpointsCoordinates} from "./classes/system/EndpointsCoordinates";
import 'whatwg-fetch'
import {communicator} from "./classes/api/Communicator";
import {MetaApi, metaApi} from "./classes/api/MetaApi";
import {StorageApi, storageApi} from "./classes/api/StorageApi";
import {systemApi, SystemApi} from "./classes/api/SystemApi";


export type MetaType = Map<string, string>|undefined

export interface StatikaApi {
    meta: MetaApi,
    storage: StorageApi
    system: SystemApi
}

/**
 * This function creates an API instance
 * @param coords - address of the Statika backend server
 * @constructor
 */
export const Statika = (coords: EndpointsCoordinates) :StatikaApi  => {
    const comm = communicator(coords)

    return {
        /*
            APIs related to object meta information
         */
        meta: metaApi(comm),

        /*
            Actual storage API: file upload, delete, listing etc
         */
        storage: storageApi(comm),

        /*
            System functionality
         */
        system: systemApi(comm)
    }
}
