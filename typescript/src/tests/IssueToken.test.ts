import { describe, test, expect } from "@jest/globals";
import { Statika } from "../Statika";
import { testCoordinates } from "../classes/system/EndpointsCoordinates";
import { v4 as uuid } from "uuid";
import { pickDefined } from "../Utilities";
import { isAuthenticationBean } from "../classes/api/AuthenticationBean";

const UPLOAD_KEY = process.env.UPLOAD_KEY ?? "TEST_UPLOAD_KEY";
const TEST_BUCKET = uuid();

const host = pickDefined(process.env.API_NODE, "127.0.0.1");
const port = pickDefined(process.env.API_PORT, "9191");

describe("Issue token", () => {
  const api = Statika(testCoordinates(host, port));

  test("successful", async () => {
    const bean = await api.system.issueToken(UPLOAD_KEY, TEST_BUCKET);
    expect(isAuthenticationBean(bean)).toBeTruthy();
  });

  test("failed", async () => {
    await expect(api.system.issueToken("bad upload key", TEST_BUCKET)).rejects.toThrow();
  });
});
