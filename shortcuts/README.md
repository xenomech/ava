# Ava shortcuts

Eight shortcuts that drive Ava from Siri, one per action. Each is plain HTTP requests
authenticated by a personal access token, so there is no login step and no password on
your phone.

| Shortcut | What it does |
| --- | --- |
| Ava · List devices | Prints your devices, so you can copy their ids. Run this first. |
| Ava · List scenes | Prints a room's scenes and the targets they hold. |
| Ava · Light on | Powers one device on. |
| Ava · Light off | Powers one device off. |
| Ava · Brightness 50 | Sets one device to 50%. Duplicate it for other levels. |
| Ava · Room on | Powers three devices on, one request each. |
| Ava · Room off | Powers three devices off. |
| Ava · Scene | Sets warmth and brightness across two devices. |

## Setting up

**1. Make a token.** In the app, Settings → Tokens. Give it `Devices · view` and
`Devices · change`; add `Scenes · view` if you want *List scenes* to work. Copy the value
— it is shown once.

**2. Import the signed files** from `signed/`. They are signed with Apple's own tool, so
they import normally, with no untrusted-shortcuts toggle.

**3. Fill in the blanks.** Every shortcut has placeholders to replace:

- `ava_pat_PASTE_YOUR_TOKEN_HERE` in the Authorization header
- `PASTE_DEVICE_UUID`, or `PASTE_DEVICE_1_UUID` and friends, in the URL
- `PASTE_ROOM_UUID` in *List scenes*

Run *Ava · List devices* first: it prints the ids everything else needs. Set up one
shortcut fully, then duplicate it for the others so the token is pasted once.

**4. Put them in a folder.** Shortcuts folders are not part of the file, so make one
called Ava in the app and drag them in. The `Ava · ` prefix keeps them together in lists
either way.

Rename each shortcut to whatever you want to say. Siri uses the name, so "Kitchen on"
beats "Ava · Light on".

## Notes

**A room is several requests, not one.** There is a batch endpoint, `/devices/apply`,
but its `targets` field is an array, and a Shortcuts JSON body field sends text as a
string — which the API rejects. One request per device avoids the problem entirely and
keeps every action the same simple shape.

**A scene is written out, not referenced.** These shortcuts set the traits directly
rather than reading a saved scene, because reading one means fetching a list and picking
from it. If you edit the scene in the app, update the shortcut to match — or use
*List scenes* to see the current targets.

**Types matter.** `power` is a Boolean, `brightness` and `color_temp` are Numbers. If you
add a field by hand, set the type in the JSON body editor rather than leaving it as text.

**Building one by hand** is a single action, if you would rather not import anything:

- **Get Contents of URL**
- URL `https://api-stage-df53.up.railway.app/api/v1/devices/<device id>/command`
- Method **POST**
- Headers: `Authorization` = `Bearer ava_pat_...`
- Request Body **JSON**: `trait` (Text) = `power`, `value` (**Boolean**) = true
