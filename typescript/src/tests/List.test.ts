/**
 * @jest-environment node
 */

import { Statika, pickDefined, testCoordinates, AuthenticationBean, authenticationBean, fileEntry } from "../index";
import { authorizeUpload } from "./helpers";
import { beforeAll, test, expect } from "@jest/globals";
import { v4 as uuid } from "uuid";
import "whatwg-fetch";

const UPLOAD_KEY = process.env.UPLOAD_KEY ?? "TEST_UPLOAD_KEY";
const TEST_BUCKET = uuid();

const host = pickDefined(process.env.API_NODE, "127.0.0.1");
const port = pickDefined(process.env.API_PORT, "9191");

const enc = new TextEncoder();
let bean: AuthenticationBean;

beforeAll(async () => {
  const uploadToken = await authorizeUpload(UPLOAD_KEY, TEST_BUCKET);
  bean = authenticationBean(uploadToken, TEST_BUCKET);
});

test("List.test_bucket_list_1", async () => {
  let s = Statika(testCoordinates(host, port));

  const buffer = enc.encode("test content");
  await s.storage.uploadFile(bean, "file1.txt", buffer);
  await s.storage.uploadFile(bean, "file2.txt", buffer);

  const list = await s.storage.listFiles(bean);
  expect(list.bucket).toStrictEqual(TEST_BUCKET);
  expect(list.files).toStrictEqual([fileEntry("file1.txt"), fileEntry("file2.txt")]);
});
