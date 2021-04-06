
export const fileEntry = (fileName: string) => ({
    filename: fileName
})

export type FileEntry = ReturnType<typeof fileEntry>