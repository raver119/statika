/**
 * @jest-environment node
 */

import {beforeAll, test, expect} from "@jest/globals"
import {Statika, responseOk} from "../src";
import {authorizeUpload} from "./helpers";

const UPLOAD_KEY = "TEST UPLOAD KET"
const TEST_BUCKET = "test_bucket"
let uploadToken: string

beforeAll(async () => {
    uploadToken = await authorizeUpload(UPLOAD_KEY, TEST_BUCKET)
})

test("Authentication.test_login_1", async () => {
    let inst = new Statika(uploadToken, TEST_BUCKET)

    await expect(inst.ping()).resolves.toStrictEqual(responseOk())
})