import { Chip, Device, DeviceHalo, cn } from "@ava/ui";
import { emitsLight } from "@ava/contracts";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "@tanstack/react-router";
import { useState } from "react";

import { hubQueries } from "@/modules/hub";
import { Loader } from "@/shared/components/loader";
import { DeviceControls } from "../components/device-controls";
import { OnAPlug, deviceColor, deviceKind, deviceLevel } from "../components/device-stage";
import { Missing, NoDevices } from "../components/empty-state";
import { HubOfflineNotice } from "../components/hub-notice";
import { useDevices } from "../use-devices";

/**
 * A device that is not in any room.
 *
 * Anything with a room opens inside that room instead, so this is the one case
 * with no surface to sit on. It is the same controls in a plain frame rather
 * than a second design: put the device in a room and it moves to the room page
 * for good.
 */
export function LooseDevicePage() {
  const { deviceId } = useParams({ from: "/_protected/devices/$deviceId" });
  const { devices, isPending } = useDevices();
  const hubs = useQuery(hubQueries.list());
  const [dragging, setDragging] = useState<number | null>(null);

  if (isPending) return <Loader label="Loading device" />;

  const device = devices.find((entry) => entry.id === deviceId);

  if (!device) {
    return devices.length === 0 ? (
      <NoDevices hasHub={(hubs.data ?? []).length > 0} />
    ) : (
      <Missing
        title="That device is gone"
        detail="It may have been removed, or it belongs to a hub you no longer have."
      />
    );
  }

  const hub = (hubs.data ?? []).find((entry) => entry.id === device.hub_id);
  const hubOffline = hub !== undefined && !hub.online;
  const offline = device.status === "offline" || hubOffline;

  const level = dragging ?? deviceLevel(device);
  const color = deviceColor(device);

  return (
    <div
      className={cn(
        "grid h-full",
        hubOffline ? "grid-rows-[auto_minmax(0,1fr)]" : "grid-rows-[minmax(0,1fr)]",
      )}
    >
      {hubOffline ? <HubOfflineNotice name={hub.name} /> : null}

      <div className="min-h-0 overflow-y-auto lg:grid lg:grid-cols-[minmax(0,1fr)_340px] lg:overflow-hidden">
        <main className="grid min-h-0 grid-rows-[34vh_auto] p-5 pt-11 sm:p-6 md:pt-6 lg:grid-rows-[minmax(0,1fr)_auto]">
          <div
            className="relative grid min-h-0 place-items-center"
            style={{ "--level": level, "--lit": color } as React.CSSProperties}
          >
            {emitsLight(deviceKind(device)) ? <DeviceHalo className="w-[46%]" /> : null}
            {device.kind === "plug" && device.appliance ? (
              <OnAPlug className="absolute right-0 top-0" />
            ) : null}
            <Device
              kind={deviceKind(device)}
              level={level}
              color={color}
              className="h-[58%] max-h-[380px]"
            />
          </div>

          <div className="pt-4">
            <h1 className="text-display font-semibold">{device.name}</h1>
            <p className="mt-1.5 flex items-center gap-2 text-small text-muted">
              Not in a room
              {offline ? (
                <Chip tone="warning">{hubOffline ? "Hub offline" : "Offline"}</Chip>
              ) : null}
            </p>
          </div>
        </main>

        <aside className="border-t border-border p-5 lg:min-h-0 lg:overflow-y-auto lg:border-l lg:border-t-0">
          <DeviceControls device={device} offline={offline} onLevelChange={setDragging} />
        </aside>
      </div>
    </div>
  );
}
