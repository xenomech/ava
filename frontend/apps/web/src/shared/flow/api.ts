import { flowStateDto, type FlowStateDto } from "@ava/contracts";

import { request } from "@/config/http/request";

export function getFlow(flowType: string, signal?: AbortSignal): Promise<FlowStateDto> {
  return request({ url: `/flows/${flowType}`, schema: flowStateDto, signal });
}

export function submitStep(flowType: string, stepId: string, data: unknown): Promise<FlowStateDto> {
  return request({
    url: `/flows/${flowType}/steps/${stepId}`,
    method: "put",
    body: { data },
    schema: flowStateDto,
  });
}

export function skipStep(flowType: string): Promise<FlowStateDto> {
  return request({ url: `/flows/${flowType}/skip`, method: "post", schema: flowStateDto });
}

export function goBack(flowType: string): Promise<FlowStateDto> {
  return request({ url: `/flows/${flowType}/back`, method: "post", schema: flowStateDto });
}
