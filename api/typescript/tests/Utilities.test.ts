/**
 * @jest-environment node
 */

import {test, expect} from "@jest/globals"
import {materializeMeta, metaInfo, pair, stringifyMeta} from "../src";

test("Utilities.test_meta_conversion", () => {
    const meta = metaInfo([pair("k1", "v1"),  pair("k2", "v2")])

    const json = stringifyMeta(meta)
    const restored = materializeMeta(JSON.parse(json))

    expect(restored).toStrictEqual(meta)
    expect(stringifyMeta(restored)).toStrictEqual(json)
    expect(meta.get("k2")).toStrictEqual(restored.get("k2"))
})