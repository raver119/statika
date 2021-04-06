import {FileEntry} from "./FileEntry";

export const listResponse = (bucket: string, files: FileEntry[] = []) => ({
    bucket: bucket,
    files: files,
})

export type ListResponse = ReturnType<typeof listResponse>
