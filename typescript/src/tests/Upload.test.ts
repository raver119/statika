/**
 * @jest-environment node
 */

import { Statika, testCoordinates, UploadResponse, AuthenticationBean, authenticationBean } from "../index";
import { authorizeUpload, httpGet } from "./helpers";
import { beforeAll, test, expect } from "@jest/globals";
import "whatwg-fetch";

const UPLOAD_KEY = process.env.UPLOAD_KEY ?? "TEST_UPLOAD_KEY";
const TEST_BUCKET = "test_bucket";

const host = process.env.API_NODE ?? "127.0.0.1";
const port = process.env.API_PORT ?? "9191";

const enc = new TextEncoder();
let bean: AuthenticationBean;
let badBucketBean: AuthenticationBean;
const badTokenBean = authenticationBean("bad token", TEST_BUCKET);

beforeAll(async () => {
  const uploadToken = await authorizeUpload(UPLOAD_KEY, TEST_BUCKET);
  bean = authenticationBean(uploadToken, TEST_BUCKET);
  badBucketBean = authenticationBean(uploadToken, "RANDOM BUCKET NAME");
});

test("Upload.test_upload_1", async () => {
  let inst = Statika(testCoordinates(host, port));

  const buffer = enc.encode("test content");
  const response = await inst.storage.uploadFile(bean, "filename.txt", buffer).then((resp) => {
    expect(resp.filename).toBe(`/${TEST_BUCKET}/filename.txt`);
    return resp as UploadResponse;
  });

  // lets check the file is actually stored
  await expect(httpGet(`${testCoordinates(host, port).toString()}${response.filename}`)).resolves.toBe("test content");
});

test("Upload.test_upload_2", async () => {
  let inst = Statika(testCoordinates(host, port));

  const buffer = enc.encode("test content");

  // empty filename is forbidden
  await expect(inst.storage.uploadFile(bean, "", buffer)).rejects.toThrow();
});

test("Upload.test_upload_3", async () => {
  let inst = Statika(testCoordinates(host, port));

  const buffer = enc.encode("test content");

  // bad token, not authorized
  await expect(inst.storage.uploadFile(badTokenBean, "filename.txt", buffer)).rejects.toThrow();
});

test("Upload.test_upload_4", async () => {
  let inst = Statika(testCoordinates(host, port));

  const buffer = enc.encode("test content");

  // not your bucket not authorized
  await expect(inst.storage.uploadFile(badBucketBean, "filename.txt", buffer)).rejects.toThrow();
});
