import {
  memberDto,
  tenantDto,
  tenantSummaryDto,
  type InviteMemberRequest,
  type MemberDto,
  type TenantDto,
  type TenantRole,
  type TenantSummaryDto,
} from "@ava/contracts";
import { z } from "zod";

import { request } from "@/config/http/request";

export function listMyTenants(signal?: AbortSignal): Promise<TenantSummaryDto[]> {
  return request({ url: "/tenants", schema: z.array(tenantSummaryDto), signal });
}

export function getCurrentTenant(signal?: AbortSignal): Promise<TenantDto> {
  return request({ url: "/tenants/current", schema: tenantDto, signal });
}

export function updateCurrentTenant(name: string): Promise<TenantDto> {
  return request({ url: "/tenants/current", method: "patch", body: { name }, schema: tenantDto });
}

export function listMembers(signal?: AbortSignal): Promise<MemberDto[]> {
  return request({ url: "/tenants/current/members", schema: z.array(memberDto), signal });
}

export function inviteMember(body: InviteMemberRequest): Promise<void> {
  return request({ url: "/tenants/current/invitations", method: "post", body });
}

export function updateMemberRole(userId: string, role: TenantRole): Promise<MemberDto> {
  return request({
    url: `/tenants/current/members/${userId}`,
    method: "patch",
    body: { role },
    schema: memberDto,
  });
}

export function removeMember(userId: string): Promise<void> {
  return request({ url: `/tenants/current/members/${userId}`, method: "delete" });
}

export function createTenant(name: string, slug: string): Promise<TenantDto> {
  return request({ url: "/tenants", method: "post", body: { name, slug }, schema: tenantDto });
}
