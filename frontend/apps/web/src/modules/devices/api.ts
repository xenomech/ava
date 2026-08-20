import { deviceDto, type DeviceDto, type SendCommandRequest, type UpdateDeviceRequest } from "@ava/contracts";
import { z } from "zod";

import { request } from "@/config/http/request";

export function listDevices(signal?: AbortSignal): Promise<DeviceDto[]> {
  return request({ url: "/devices", schema: z.array(deviceDto), signal });
}

export function updateDevice(deviceID: string, body: UpdateDeviceRequest): Promise<DeviceDto> {
  return request({ url: `/devices/${deviceID}`, method: "patch", body, schema: deviceDto });
}

export function sendCommand(deviceID: string, body: SendCommandRequest): Promise<void> {
  return request({ url: `/devices/${deviceID}/command`, method: "post", body });
}
