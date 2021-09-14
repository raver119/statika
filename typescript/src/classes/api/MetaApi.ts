import { Communicator } from "./Communicator";
import { ApiResponse } from "../entities/ApiResponse";
import { MetaType } from "../../Statika";
import { AuthenticationBean } from "./AuthenticationBean";
import { materializeMeta, objectifyMeta } from "../../Utilities";

export const metaApi = (communicator: Communicator) => {
  return {
    updateMetaInfo(bean: AuthenticationBean, fileName: string, metaInfo: MetaType): Promise<ApiResponse> {
      // serialize map
      return communicator.post(bean, `/meta/${bean.bucket}/${fileName}`, objectifyMeta(metaInfo));
    },

    getMetaInfo(bean: AuthenticationBean, fileName: string): Promise<MetaType> {
      // deserialize map
      return communicator.get(bean, `/meta/${bean.bucket}/${fileName}`).then((data) => materializeMeta(data));
    },

    deleteMetaInfo(bean: AuthenticationBean, fileName: string): Promise<ApiResponse> {
      return communicator.delete(bean, `/meta/${bean.bucket}/${fileName}`);
    },
  };
};

export type MetaApi = ReturnType<typeof metaApi>;
