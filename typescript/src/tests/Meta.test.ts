/**
 * @jest-environment node
 */

import {
    Statika,
    testCoordinates,
    AuthenticationBean,
    authenticationBean,
    metaInfo, pair
} from "../index";
import {authorizeUpload} from "./helpers";
import {beforeAll, test, expect} from "@jest/globals"
import {v4 as uuid} from "uuid";
import 'whatwg-fetch'

const UPLOAD_KEY =  process.env.UPLOAD_KEY ?? "TEST_UPLOAD_KEY"
const TEST_BUCKET = uuid()

const host = process.env.API_NODE ?? "127.0.0.1"
const port = process.env.API_PORT ?? "9191"

let bean: AuthenticationBean

beforeAll(async () => {
    const uploadToken = await authorizeUpload(UPLOAD_KEY, TEST_BUCKET)
    bean = authenticationBean(uploadToken, TEST_BUCKET)
})

test("Meta.test_meta_crd", async () => {
    let s = Statika(testCoordinates(host, port))
    const fileName = "random.file"
    const meta = metaInfo([pair("k1", "v1"),  pair("k2", "v2")])

    await expect(s.meta.updateMetaInfo(bean, fileName, meta)).resolves.toBeDefined()

    const restored = await s.meta.getMetaInfo(bean, fileName)
    expect(restored).toStrictEqual(meta)

    await expect(s.meta.deleteMetaInfo(bean, fileName)).resolves.toBeDefined()

    const empty =  await s.meta.getMetaInfo(bean, fileName)
    expect(empty).toStrictEqual(metaInfo([]))
})