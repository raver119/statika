/**
 * @jest-environment node
 */

import {beforeAll, test, expect} from "@jest/globals"
import {Statika, responseOk, testCoordinates, AuthenticationBean, authenticationBean} from "../index";
import {authorizeUpload} from "./helpers";

const UPLOAD_KEY = process.env.UPLOAD_KEY ?? "TEST_UPLOAD_KEY"
const TEST_BUCKET = "test_bucket"

const host = process.env.API_NODE ?? "127.0.0.1"
const port = process.env.API_PORT ??  "9191"

let uploadToken: string
let bean: AuthenticationBean
let badBucketBean: AuthenticationBean
const badTokenBean = authenticationBean("bad token", TEST_BUCKET)

beforeAll(async () => {
    uploadToken = await authorizeUpload(UPLOAD_KEY, TEST_BUCKET)
    bean = authenticationBean(uploadToken, TEST_BUCKET)
    badBucketBean = authenticationBean(uploadToken, "RANDOM BUCKET NAME")
})

test("Authentication.test_login_1", async () => {
    let inst = Statika(testCoordinates(host, port))

    await expect(inst.system.ping(bean)).resolves.toStrictEqual({...responseOk()})
})

test("Authentication.test_login_2", async () => {
    let inst = Statika(testCoordinates(host, port))


    await expect(inst.system.ping(badTokenBean)).rejects.toThrow()
})

test("Authentication.test_login_3", async () => {
    let inst = Statika(testCoordinates(host, port))

    // this test must pass since ping has nothing to do with actual storage
    await expect(inst.system.ping(badBucketBean)).resolves.toStrictEqual({...responseOk()})
})