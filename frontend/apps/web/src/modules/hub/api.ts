import { hubDto, type ActivateHubRequest, type HubDto } from "@ava/contracts";
import { z } from "zod";

import { request } from "@/config/http/request";

export function listHubs(signal?: AbortSignal): Promise<HubDto[]> {
  return request({ url: "/hubs", schema: z.array(hubDto), signal });
}

export function activateHub(body: ActivateHubRequest): Promise<HubDto> {
  return request({ url: "/hubs/activate", method: "post", body, schema: hubDto });
}

export function removeHub(hubID: string): Promise<void> {
  return request({ url: `/hubs/${hubID}`, method: "delete" });
}
