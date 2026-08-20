import { Slider, Tabs, TabsContent, TabsList, TabsTrigger } from "@ava/ui";

import { cssToHue } from "@/shared/lib/color";
import { useLiveSlider } from "../use-live-slider";

const KELVIN_STEP = 50;

const SWATCHES = [
  { css: "#ff5a3c", name: "Ember" },
  { css: "#ffb347", name: "Amber" },
  { css: "#ffe28a", name: "Candle" },
  { css: "#b8ff6b", name: "Lime" },
  { css: "#4bffb5", name: "Mint" },
  { css: "#4bd6ff", name: "Lagoon" },
  { css: "#6b8cff", name: "Dusk" },
  { css: "#b46bff", name: "Violet" },
  { css: "#ff6bd6", name: "Blossom" },
  { css: "#ffffff", name: "Cool white" },
] as const;

const KELVIN_RAMP =
  "linear-gradient(to right, #ff9233, #ffbe7a, #ffe4c4, #ffffff, #d6e6ff, #a8c8ff)";
const HUE_RAMP =
  "linear-gradient(to right,#ff2d2d,#ffd400 17%,#4bff4b 33%,#00e5ff 50%,#2d6bff 67%,#b44bff 83%,#ff2d8a)";

export function ColorControl({
  color,
  kelvin,
  kelvinMin,
  kelvinMax,
  showColor = false,
  disabled = false,
  onWhitePreview,
  onWhite,
  onColor,
}: {
  color: string;
  kelvin: number | null;
  kelvinMin: number;
  kelvinMax: number;
  showColor?: boolean;
  disabled?: boolean;
  onWhitePreview: (kelvin: number) => void;
  onWhite: (kelvin: number) => void;
  onColor: (color: string) => void;
}) {
  const white = useLiveSlider(kelvin ?? kelvinMin, onWhitePreview, onWhite);
  const hue = (cssToHue(color) / 360) * 100;

  return (
    <Tabs defaultValue={showColor && kelvin === null ? "color" : "white"} className="grid gap-2.5">
      {showColor ? (
        <TabsList className="w-full">
          <TabsTrigger value="white">White</TabsTrigger>
          <TabsTrigger value="color">Colour</TabsTrigger>
        </TabsList>
      ) : null}

      <TabsContent value="white" className="grid gap-2">
        <Slider
          value={[white.value]}
          min={kelvinMin}
          max={kelvinMax}
          step={KELVIN_STEP}
          variant="marker"
          disabled={disabled}
          aria-label="White temperature"
          aria-valuetext={`${white.value} kelvin`}
          onValueChange={([v]) => white.change(v ?? kelvinMin)}
          onValueCommit={([v]) => white.release(v ?? kelvinMin)}
          style={{ background: KELVIN_RAMP }}
        />
        <p className="font-mono text-caption text-subtle tabular">{white.value}K</p>
      </TabsContent>

      <TabsContent value="color" className={showColor ? "grid gap-2.5" : "hidden"}>
        <Slider
          value={[hue]}
          max={100}
          step={1}
          variant="marker"
          disabled={disabled}
          aria-label="Hue"
          onValueCommit={([v]) => onColor(`hsl(${Math.round(((v ?? 0) / 100) * 360)} 92% 62%)`)}
          style={{ background: HUE_RAMP }}
        />

        <div className="grid grid-cols-5 gap-1.5">
          {SWATCHES.map((s) => (
            <button
              key={s.css}
              type="button"
              disabled={disabled}
              aria-label={s.name}
              aria-pressed={color === s.css}
              onClick={() => onColor(s.css)}
              style={{ background: s.css }}
              className="aspect-square rounded-sm border border-border transition-transform duration-150 ease-out active:scale-95 disabled:opacity-40 aria-pressed:ring-2 aria-pressed:ring-fg aria-pressed:ring-offset-2 aria-pressed:ring-offset-surface"
            />
          ))}
        </div>
      </TabsContent>
    </Tabs>
  );
}
