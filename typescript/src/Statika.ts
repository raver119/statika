import { EndpointsCoordinates } from "./classes/system/EndpointsCoordinates";
import "whatwg-fetch";
import { communicator } from "./classes/api/Communicator";
import { MetaApi, metaApi } from "./classes/api/MetaApi";
import { StorageApi, storageApi } from "./classes/api/StorageApi";
import { systemApi, SystemApi } from "./classes/api/SystemApi";

// TODO: replace this type with proper semi-defined type here and in the backend
export type MetaType = Map<string, string>;

export const metaInfo = (values: { k: string; v: string }[]): MetaType => {
  const m = new Map<string, string>();
  values.forEach((p) => m.set(p.k, p.v));
  return m;
};

export const pair = (k: string, v: string) => ({
  k: k,
  v: v,
});

export interface StatikaApi {
  meta: MetaApi;
  storage: StorageApi;
  system: SystemApi;
}

/**
 * This function creates an API instance
 * @param coords - address of the Statika backend server
 * @constructor
 */
export const Statika = (coords: EndpointsCoordinates): StatikaApi => {
  const comm = communicator(coords);

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
    system: systemApi(comm),
  };
};
