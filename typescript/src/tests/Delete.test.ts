/**
 * @jest-environment node
 */

import { Statika, testCoordinates, UploadResponse, AuthenticationBean, authenticationBean } from "../index";
import { authorizeUpload, httpGet } from "./helpers";
import { beforeAll, test, expect } from "@jest/globals";
import { HttpStatusCode } from "../classes/entities/HttpStatusCode";

const UPLOAD_KEY = process.env.UPLOAD_KEY ?? "TEST_UPLOAD_KEY";
const TEST_BUCKET = "test_bucket";
const EVIL_BUCKET = "evil_bucket";

const host = process.env.API_NODE ?? "127.0.0.1";
const port = process.env.API_PORT ?? "9191";

const enc = new TextEncoder();
let goodBean: AuthenticationBean;
let evilBean: AuthenticationBean;

beforeAll(async () => {
  const uploadToken = await authorizeUpload(UPLOAD_KEY, TEST_BUCKET);
  const evilToken = await authorizeUpload(UPLOAD_KEY, EVIL_BUCKET);

  goodBean = authenticationBean(uploadToken, TEST_BUCKET);
  evilBean = authenticationBean(evilToken, EVIL_BUCKET);
});

test("Delete.test_delete_1", async () => {
  let inst = Statika(testCoordinates(host, port));

  const buffer = enc.encode("test content");
  let response = await inst.storage.uploadFile(goodBean, "fileToDelete.txt", buffer).then((resp) => {
    expect(resp.filename).toBe(`/${TEST_BUCKET}/fileToDelete.txt`);
    return resp;
  });

  await inst.storage.deleteFile(goodBean, response.filename).then((resp) => {
    expect(resp.statusCode).toBe(HttpStatusCode.OK);
  });

  // file must be gone now
  await expect(httpGet(`${testCoordinates(host, port).toString()}${response.filename}`)).rejects.toThrow(
    /API response code: 404/
  );
});

test("Delete.test_delete_2", async () => {
  let inst = Statika(testCoordinates(host, port));
  let evil = Statika(testCoordinates(host, port));

  const buffer = enc.encode("test content");
  let response = await inst.storage.uploadFile(goodBean, "AnotherFileToDelete.txt", buffer).then((resp) => {
    expect(resp.filename).toBe(`/${TEST_BUCKET}/AnotherFileToDelete.txt`);
    return resp as UploadResponse;
  });

  await expect(evil.storage.deleteFile(evilBean, response.filename)).rejects.toThrow(/API response code: 401/);
});
