/**
 * @jest-environment node
 */

import {Statika, pickDefined, testCoordinates, UploadResponse} from "../src";
import {authorizeUpload, httpGet} from "./helpers";
import {beforeAll, test, expect} from "@jest/globals"
import 'whatwg-fetch'

const UPLOAD_KEY = "TEST_UPLOAD_KEY"
const TEST_BUCKET = "test_bucket"

const host = pickDefined(process.env.API_NODE, "127.0.0.1")
const port = pickDefined(process.env.API_PORT, "8080")

const enc = new TextEncoder()
let uploadToken: string

beforeAll(async () => {
    uploadToken = await authorizeUpload(UPLOAD_KEY, TEST_BUCKET)
})

test("Upload.test_upload_1", async () => {
    let inst = new Statika(testCoordinates(host, port), uploadToken, TEST_BUCKET)

    const buffer = enc.encode("test content")
    const response = await inst.uploadFile("filename.txt", buffer).then(resp => {
        expect(resp.filename).toBe(`/${TEST_BUCKET}/filename.txt`)
        return resp as UploadResponse
    })

    // lets check the file is actually stored
    await expect(httpGet(`${testCoordinates(host, port).toString()}${response.filename}`)).resolves.toBe("test content")
})

test("Upload.test_upload_2", async () => {
    let inst = new Statika(testCoordinates(host, port), uploadToken, TEST_BUCKET)

    const buffer = enc.encode("test content")

    // empty filename is forbidden
    await expect(inst.uploadFile("", buffer)).rejects.toThrow()
})

test("Upload.test_upload_3", async () => {
    let inst = new Statika(testCoordinates(host, port), "bad token", TEST_BUCKET)

    const buffer = enc.encode("test content")

    // bad token, not authorized
    await expect(inst.uploadFile("filename.txt", buffer)).rejects.toThrow()
})

test("Upload.test_upload_4", async () => {
    let inst = new Statika(testCoordinates(host, port), uploadToken, "another bucket")

    const buffer = enc.encode("test content")

    // not your bucket not authorized
    await expect(inst.uploadFile("filename.txt", buffer)).rejects.toThrow()
})