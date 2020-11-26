/**
 * @jest-environment node
 */

import {beforeAll, test, expect} from "@jest/globals"
import {Statika, responseOk, pickDefined, testCoordinates} from "../src";
import {authorizeUpload} from "./helpers";

const UPLOAD_KEY = "TEST_UPLOAD_KEY"
const TEST_BUCKET = "test_bucket"

const host = pickDefined(process.env.API_NODE, "127.0.0.1")
const port = pickDefined(process.env.API_PORT, "8080")

let uploadToken: string

beforeAll(async () => {
    uploadToken = await authorizeUpload(UPLOAD_KEY, TEST_BUCKET)
})

test("Authentication.test_login_1", async () => {
    let inst = new Statika(testCoordinates(host, port), uploadToken, TEST_BUCKET)

    await expect(inst.ping()).resolves.toStrictEqual({...responseOk()})
})

test("Authentication.test_login_2", async () => {
    let inst = new Statika(testCoordinates(host, port), "bad token", TEST_BUCKET)

    await expect(inst.ping()).rejects.toThrow()
})

test("Authentication.test_login_3", async () => {
    let inst = new Statika(testCoordinates(host, port), uploadToken, "RANDOM BUCKET NAME")

    // this test must pass since ping has nothing to do with actual storage
    await expect(inst.ping()).resolves.toStrictEqual({...responseOk()})
})