export function cssToHue(css: string): number {
  const rgb = toRgb(css.trim());
  if (!rgb) return 0;

  const [r, g, b] = rgb.map((v) => v / 255) as [number, number, number];
  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  const span = max - min;

  if (span === 0) return 0;

  const hue =
    max === r ? ((g - b) / span) % 6 : max === g ? (b - r) / span + 2 : (r - g) / span + 4;

  return (((hue * 60) % 360) + 360) % 360;
}

function toRgb(css: string): [number, number, number] | null {
  const hsl = /^hsl\(\s*([\d.]+)/.exec(css);
  if (hsl) return hueOnly(Number(hsl[1]));

  const hex = /^#([0-9a-f]{6})$/i.exec(css);
  if (hex) {
    const n = Number.parseInt(hex[1]!, 16);
    return [(n >> 16) & 255, (n >> 8) & 255, n & 255];
  }

  const rgb = /^rgba?\(\s*([\d.]+)[\s,]+([\d.]+)[\s,]+([\d.]+)/.exec(css);
  if (rgb) return [Number(rgb[1]), Number(rgb[2]), Number(rgb[3])];

  return null;
}

function hueOnly(hue: number): [number, number, number] {
  const h = (((hue % 360) + 360) % 360) / 60;
  const x = Math.round(255 * (1 - Math.abs((h % 2) - 1)));

  const table: [number, number, number][] = [
    [255, x, 0],
    [x, 255, 0],
    [0, 255, x],
    [0, x, 255],
    [x, 0, 255],
    [255, 0, x],
  ];

  return table[Math.floor(h) % 6]!;
}
