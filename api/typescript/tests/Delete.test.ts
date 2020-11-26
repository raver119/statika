/**
 * @jest-environment node
 */

import {Statika, pickDefined, testCoordinates, UploadResponse} from "../src";
import {authorizeUpload, httpGet} from "./helpers";
import {beforeAll, test, expect} from "@jest/globals"
import {HttpStatusCode} from "../src/classes/comms/HttpStatusCode";

const UPLOAD_KEY = "TEST_UPLOAD_KEY"
const TEST_BUCKET = "test_bucket"
const EVIL_BUCKET = "evil_bucket"

const host = pickDefined(process.env.API_NODE, "127.0.0.1")
const port = pickDefined(process.env.API_PORT, "8080")

const enc = new TextEncoder()
let uploadToken: string
let evilToken: string

beforeAll(async () => {
    uploadToken = await authorizeUpload(UPLOAD_KEY, TEST_BUCKET)
    evilToken = await authorizeUpload(UPLOAD_KEY, EVIL_BUCKET)
})

test("Delete.test_delete_1", async () => {
    let inst = new Statika(testCoordinates(host, port), uploadToken, TEST_BUCKET)

    const buffer = enc.encode("test content")
    let response = await inst.uploadFile("fileToDelete.txt", buffer).then(resp => {
        expect(resp.filename).toBe(`/${TEST_BUCKET}/fileToDelete.txt`)
        return resp as UploadResponse
    })


    await inst.deleteFile(response.filename).then(resp => {
        expect(resp.statusCode).toBe(HttpStatusCode.OK)
    })

    // file must be gone now
    await expect(httpGet(`${testCoordinates(host, port).toString()}${response.filename}`)).rejects.toThrow(/API response code: 404/)
})

test("Delete.test_delete_2", async () => {
    let inst = new Statika(testCoordinates(host, port), uploadToken, TEST_BUCKET)
    let evil = new Statika(testCoordinates(host, port), evilToken, EVIL_BUCKET)

    const buffer = enc.encode("test content")
    let response = await inst.uploadFile("AnotherFileToDelete.txt", buffer).then(resp => {
        expect(resp.filename).toBe(`/${TEST_BUCKET}/AnotherFileToDelete.txt`)
        return resp as UploadResponse
    })


    await expect(evil.deleteFile(response.filename)).rejects.toThrow(/API response code: 401/)
})