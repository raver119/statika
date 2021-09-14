import { describe, test, beforeAll, expect, jest } from "@jest/globals";
import fs from "fs";
import path from "path";
import os from "os";
import { v4 as uuid } from "uuid";
import { authenticationBean, AuthenticationBean } from "../classes/api/AuthenticationBean";
import { authorizeUpload } from "./helpers";
import { Statika } from "../Statika";
import { testCoordinates } from "../classes/system/EndpointsCoordinates";
import { makeRandomString, randomInt } from "../Utilities";

jest.setTimeout(30000);

const UPLOAD_KEY = process.env.UPLOAD_KEY ?? "TEST_UPLOAD_KEY";
const TEST_BUCKET = "test_bucket";

const host = process.env.API_NODE ?? "127.0.0.1";
const port = process.env.API_PORT ?? "9191";

let bean: AuthenticationBean;
beforeAll(async () => {
  const uploadToken = await authorizeUpload(UPLOAD_KEY, TEST_BUCKET);
  bean = authenticationBean(uploadToken, TEST_BUCKET);
});

describe("Large uploads", () => {
  const api = Statika(testCoordinates(host, port));

  test("Single large file", async () => {
    const folderName = createFolderWithFiles(1, 1000000);

    const files = fs.readdirSync(folderName);
    const promises = files
      .map((f) => ({ filename: f, content: fs.readFileSync(path.join(folderName, f)) }))
      .map((f) => api.storage.uploadFile(bean, f.filename, f.content));

    await expect(Promise.all(promises)).resolves.toBeDefined();
  });

  test("Multiple files", async () => {
    const num = 10;
    const folderName = createFolderWithFiles(num);

    const files = fs.readdirSync(folderName);
    const promises = files
      .map((f) => ({ filename: f, content: fs.readFileSync(path.join(folderName, f)) }))
      .map((f) => api.storage.uploadFile(bean, f.filename, f.content));

    expect(files.length).toBe(num);
    await expect(Promise.all(promises)).resolves.toBeDefined();
  });
});

// this function creates
function createFolderWithFiles(numFiles: number, minFileSize = 100, maxFileSize = 1000000): string {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "multiple-"));

  for (let i = 0; i < numFiles; i++) {
    const fileName = path.join(tmp, `${uuid()}.txt`);
    const fileSize = randomInt(minFileSize, maxFileSize);
    fs.writeFileSync(fileName, makeRandomString(fileSize));
  }

  return tmp;
}
