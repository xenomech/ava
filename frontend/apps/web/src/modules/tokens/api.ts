import {
  apiTokenDto,
  createdApiTokenDto,
  type ApiTokenDto,
  type CreateApiTokenRequest,
  type CreatedApiTokenDto,
} from "@ava/contracts";
import { z } from "zod";

import { request } from "@/config/http/request";

export function listTokens(signal?: AbortSignal): Promise<ApiTokenDto[]> {
  return request({ url: "/tokens", schema: z.array(apiTokenDto), signal });
}

export function createToken(body: CreateApiTokenRequest): Promise<CreatedApiTokenDto> {
  return request({ url: "/tokens", method: "post", body, schema: createdApiTokenDto });
}

export function revokeToken(tokenId: string): Promise<void> {
  return request({ url: `/tokens/${tokenId}/revoke`, method: "post" });
}

export function deleteToken(tokenId: string): Promise<void> {
  return request({ url: `/tokens/${tokenId}`, method: "delete" });
}
