"""Regenerates the Ava mark: a lowercase Gabarito SemiBold 'a' on the app's dark tile.

Writes public/logo.png (the source the PWA icons are generated from) and public/favicon.svg
(the same letter as an outline, so it needs no font at runtime). Run `pnpm generate-pwa-assets`
afterwards to rebuild the PNG set and favicon.ico from the new logo.png.

    pip install pillow fonttools
    python3 scripts/make-mark.py
"""

import io
import pathlib
import urllib.request

from fontTools.pens.boundsPen import BoundsPen
from fontTools.pens.svgPathPen import SVGPathPen
from fontTools.pens.transformPen import TransformPen
from fontTools.ttLib import TTFont
from PIL import Image, ImageDraw, ImageFont

# Gabarito SemiBold, the weight the Wordmark component uses. Asking as an old browser gets TTF.
FONT_CSS = "https://fonts.googleapis.com/css2?family=Gabarito:wght@600"
FONT_UA = "Mozilla/5.0 (Linux; U; Android 2.3.7; en-us) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"

LETTER = "a"
TILE = "#0b0b0d"
INK = "#ffffff"
RADIUS = 0.22
MASTER = 1024

# The letter fills more of the tile in the SVG, which stays crisp where the bitmap would blur.
PNG_GLYPH = 0.60
SVG_GLYPH = 0.62
SVG_VIEW = 32

PUBLIC = pathlib.Path(__file__).resolve().parent.parent / "public"


def fetch_font() -> bytes:
    css = urllib.request.urlopen(
        urllib.request.Request(FONT_CSS, headers={"User-Agent": FONT_UA})
    ).read().decode()

    start = css.index("url(") + 4
    return urllib.request.urlopen(css[start : css.index(")", start)]).read()


def write_png(font_bytes: bytes) -> None:
    """Draws the tile supersampled, so the corner radius and the bowl come down clean."""
    scale = 4
    size = MASTER * scale
    image = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(image)
    draw.rounded_rectangle([0, 0, size - 1, size - 1], radius=size * RADIUS, fill=TILE)

    # Size and centre the letter on its own ink, not on the font's line box.
    probe = ImageFont.truetype(io.BytesIO(font_bytes), 200)
    left, top, right, bottom = probe.getbbox(LETTER)
    points = round(200 * size * PNG_GLYPH / max(right - left, bottom - top))

    font = ImageFont.truetype(io.BytesIO(font_bytes), points)
    left, top, right, bottom = font.getbbox(LETTER)
    draw.text(
        ((size - (right - left)) / 2 - left, (size - (bottom - top)) / 2 - top),
        LETTER,
        font=font,
        fill=INK,
    )

    image.resize((MASTER, MASTER), Image.LANCZOS).save(PUBLIC / "logo.png")


def write_svg(font_bytes: bytes) -> None:
    font = TTFont(io.BytesIO(font_bytes))
    glyphs = font.getGlyphSet()
    glyph = glyphs[font.getBestCmap()[ord(LETTER)]]

    bounds = BoundsPen(glyphs)
    glyph.draw(bounds)
    x0, y0, x1, y1 = bounds.bounds
    scale = SVG_VIEW * SVG_GLYPH / max(x1 - x0, y1 - y0)

    # Fonts measure upwards and SVG measures down, so the y axis flips here.
    dx = (SVG_VIEW - (x1 - x0) * scale) / 2 - x0 * scale
    dy = (SVG_VIEW - (y1 - y0) * scale) / 2 + y1 * scale

    pen = SVGPathPen(glyphs, ntos=lambda v: f"{v:.2f}".rstrip("0").rstrip("."))
    glyph.draw(TransformPen(pen, (scale, 0, 0, -scale, dx, dy)))

    (PUBLIC / "favicon.svg").write_text(
        f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {SVG_VIEW} {SVG_VIEW}">'
        f'<rect width="{SVG_VIEW}" height="{SVG_VIEW}" rx="{SVG_VIEW * RADIUS:.2f}" fill="{TILE}"/>'
        f'<path d="{pen.getCommands()}" fill="{INK}"/>'
        "</svg>\n"
    )


if __name__ == "__main__":
    font_bytes = fetch_font()
    write_png(font_bytes)
    write_svg(font_bytes)
    print(f"wrote {PUBLIC / 'logo.png'} and {PUBLIC / 'favicon.svg'}")
    print("now run: pnpm generate-pwa-assets")
