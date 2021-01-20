import {Communicator} from "./Communicator";
import {ApiResponse} from "../entities/ApiResponse";
import {MetaType} from "../../Statika";


export const metaApi = (communicator: Communicator) => {
    return {
        updateMetaInfo(fileName: string, metaInfo: MetaType = undefined) :Promise<ApiResponse> {
            return undefined
        },

        getMetaInfo(fileName: string) :Promise<MetaType> {
            return undefined
        },

        deleteMetaInfo(fileName: string) :Promise<ApiResponse> {
            return undefined
        },
    }
}

export type MetaApi = ReturnType<typeof metaApi>